/**
 * Teleport
 * Copyright (C) 2023  Gravitational, Inc.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

import { ChildProcess, exec, fork, spawn } from 'node:child_process';
import fs from 'node:fs/promises';
import path from 'node:path';
import { promisify } from 'node:util';

import { ChannelCredentials } from '@grpc/grpc-js';
import { GrpcTransport } from '@protobuf-ts/grpc-transport';
import {
  app,
  dialog,
  ipcMain,
  Menu,
  MenuItemConstructorOptions,
  nativeTheme,
  shell,
} from 'electron';

import { AbortError } from 'shared/utils/error';

import Logger from 'teleterm/logger';
import { getAssetPath } from 'teleterm/mainProcess/runtimeSettings';
import {
  ChildProcessAddresses,
  MainProcessClient,
  MainProcessIpc,
  RendererIpc,
  TERMINATE_MESSAGE,
} from 'teleterm/mainProcess/types';
import {
  TSH_AUTOUPDATE_ENV_VAR,
  TSH_AUTOUPDATE_OFF,
} from 'teleterm/node/tshAutoupdate';
import { AppUpdater, AppUpdaterStorage } from 'teleterm/services/appUpdater';
import { subscribeToFileStorageEvents } from 'teleterm/services/fileStorage';
import * as grpcCreds from 'teleterm/services/grpcCredentials';
import {
  createFileLoggerService,
  KeepLastChunks,
  LoggerColor,
} from 'teleterm/services/logger';
import {
  AutoUpdateClient,
  createAutoUpdateClient,
  createTshdClient,
  TshdClient,
} from 'teleterm/services/tshd';
import { loggingInterceptor } from 'teleterm/services/tshd/interceptors';
import { staticConfig } from 'teleterm/staticConfig';
import { FileStorage, RuntimeSettings } from 'teleterm/types';
import { RootClusterUri } from 'teleterm/ui/uri';

import {
  ConfigService,
  subscribeToConfigServiceEvents,
} from '../services/config';
import { downloadAgent, FileDownloader, verifyAgent } from './agentDownloader';
import { AgentRunner } from './agentRunner';
import { subscribeToTabContextMenuEvent } from './contextMenus/tabContextMenu';
import { subscribeToTerminalContextMenuEvent } from './contextMenus/terminalContextMenu';
import {
  createAgentConfigFile,
  generateAgentConfigPaths,
  getAgentsDir,
  isAgentConfigFileCreated,
  removeAgentDirectory,
  type CreateAgentConfigFileArgs,
} from './createAgentConfigFile';
import { ResolveError, resolveNetworkAddress } from './resolveNetworkAddress';
import { terminateWithTimeout } from './terminateWithTimeout';
import { WindowsManager } from './windowsManager';

type Options = {
  settings: RuntimeSettings;
  logger: Logger;
  configService: ConfigService;
  appStateFileStorage: FileStorage;
  configFileStorage: FileStorage;
  windowsManager: WindowsManager;
};

export default class MainProcess {
  readonly settings: RuntimeSettings;
  private readonly logger: Logger;
  private readonly configService: ConfigService;
  private tshdProcess: ChildProcess;
  private tshdProcessLastLogs: KeepLastChunks<string>;
  private sharedProcess: ChildProcess;
  private sharedProcessLastLogs: KeepLastChunks<string>;
  private appStateFileStorage: FileStorage;
  private configFileStorage: FileStorage;
  private resolvedChildProcessAddresses: Promise<ChildProcessAddresses>;
  private windowsManager: WindowsManager;
  // this function can be safely called concurrently
  private downloadAgentShared = sharePromise(() =>
    downloadAgent(
      new FileDownloader(this.windowsManager.getWindow()),
      this.settings,
      process.env
    )
  );
  private readonly agentRunner: AgentRunner;
  private tshdClients: Promise<{
    terminalService: TshdClient;
    autoUpdateService: AutoUpdateClient;
  }>;
  private readonly appUpdater: AppUpdater;

  private constructor(opts: Options) {
    this.settings = opts.settings;
    this.logger = opts.logger;
    this.configService = opts.configService;
    this.appStateFileStorage = opts.appStateFileStorage;
    this.configFileStorage = opts.configFileStorage;
    this.windowsManager = opts.windowsManager;
    this.agentRunner = new AgentRunner(
      this.settings,
      path.join(__dirname, 'agentCleanupDaemon.js'),
      (rootClusterUri, state) => {
        const window = this.windowsManager.getWindow();
        if (window.isDestroyed()) {
          return;
        }
        window.webContents.send(
          RendererIpc.ConnectMyComputerAgentUpdate,
          rootClusterUri,
          state
        );
      }
    );

    const getClusterVersions = async () => {
      const { autoUpdateService } = await this.getTshdClients();
      const { response } = await autoUpdateService.getClusterVersions({});
      return response;
    };
    const getDownloadBaseUrl = async () => {
      const { autoUpdateService } = await this.getTshdClients();
      const {
        response: { baseUrl },
      } = await autoUpdateService.getDownloadBaseUrl({});
      return baseUrl;
    };
    this.appUpdater = new AppUpdater(
      makeAppUpdaterStorage(this.appStateFileStorage),
      getClusterVersions,
      getDownloadBaseUrl
    );
  }

  /**
   * create starts necessary child processes such as tsh daemon and the shared process. It also sets
   * up IPC handlers and resolves the network addresses under which the child processes set up gRPC
   * servers.
   *
   * create might throw an error if spawning a child process fails, see initTshd for more details.
   */
  static create(opts: Options) {
    const instance = new MainProcess(opts);
    instance.init();
    return instance;
  }

  async dispose(): Promise<void> {
    this.windowsManager.dispose();
    this.appUpdater.dispose();
    await Promise.all([
      // sending usage events on tshd shutdown has 10-seconds timeout
      terminateWithTimeout(this.tshdProcess, 10_000, () => {
        this.gracefullyKillTshdProcess();
      }),
      terminateWithTimeout(this.sharedProcess, 5_000, process =>
        // process.kill doesn't allow running a cleanup code in the child process
        // on Windows
        process.send(TERMINATE_MESSAGE)
      ),
      this.agentRunner.killAll(),
    ]);
  }

  private init() {
    this.updateAboutPanelIfNeeded();
    this.setAppMenu();
    this.initTshd();
    this.initSharedProcess();
    this.initResolvingChildProcessAddresses();
    this.initIpc();
  }

  async getTshdClients(): Promise<{
    terminalService: TshdClient;
    autoUpdateService: AutoUpdateClient;
  }> {
    if (!this.tshdClients) {
      this.tshdClients = this.resolvedChildProcessAddresses.then(
        ({ tsh: tshdAddress }) =>
          setUpTshdClients({
            runtimeSettings: this.settings,
            tshdAddress,
          })
      );
    }

    return this.tshdClients;
  }

  private initTshd() {
    const { binaryPath, homeDir } = this.settings.tshd;
    this.logger.info(`Starting tsh daemon from ${binaryPath}`);

    // spawn might either fail immediately by throwing an error or cause the error event to be emitted
    // on the process value returned by spawn.
    //
    // Some spawn failures result in an error being thrown immediately such as EPERM on Windows [1].
    // This might be related to the fact that in Node.js, an error event causes an error to be
    // thrown if there are no listeners for the error event itself [2].
    //
    // [1] https://stackoverflow.com/a/42262771
    // [2] https://nodejs.org/docs/latest-v22.x/api/errors.html#error-propagation-and-interception
    this.tshdProcess = spawn(
      binaryPath,
      ['daemon', 'start', ...this.getTshdFlags()],
      {
        stdio: 'pipe', // stdio must be set to `pipe` as the gRPC server address is read from stdout
        windowsHide: true,
        env: {
          ...process.env,
          TELEPORT_HOME: homeDir,
          [TSH_AUTOUPDATE_ENV_VAR]: TSH_AUTOUPDATE_OFF,
        },
      }
    );

    this.logProcessExitAndError('tshd', this.tshdProcess);

    const loggerService = createFileLoggerService({
      dev: this.settings.dev,
      dir: this.settings.logsDir,
      name: TSHD_LOGGER_NAME,
      loggerNameColor: LoggerColor.Cyan,
      passThroughMode: true,
    });

    this.tshdProcessLastLogs = new KeepLastChunks(NO_OF_LAST_LOGS_KEPT);
    loggerService.pipeProcessOutputIntoLogger(
      this.tshdProcess,
      this.tshdProcessLastLogs
    );
  }

  private getTshdFlags(): string[] {
    const settings = this.settings;
    const agentsDir = getAgentsDir(settings.userDataDir);

    const flags = [
      // grpc-js requires us to pass localhost:port for TCP connections,
      // for tshd we have to specify the protocol as well.
      `--addr=${settings.tshd.requestedNetworkAddress}`,
      `--certs-dir=${settings.certsDir}`,
      `--prehog-addr=${staticConfig.prehogAddress}`,
      `--kubeconfigs-dir=${settings.kubeConfigsDir}`,
      `--agents-dir=${agentsDir}`,
      `--installation-id=${settings.installationId}`,
      `--add-keys-to-agent=${this.configService.get('sshAgent.addKeysToAgent').value}`,
    ];

    if (this.configService.get('hardwareKeyAgent.enabled').value) {
      flags.unshift('--hardware-key-agent');
    }
    if (settings.insecure) {
      flags.unshift('--insecure');
    }
    if (settings.debug) {
      flags.unshift('--debug');
    }

    return flags;
  }

  private initSharedProcess() {
    this.sharedProcess = fork(
      path.join(__dirname, 'sharedProcess.js'),
      [`--runtimeSettingsJson=${JSON.stringify(this.settings)}`],
      {
        stdio: 'pipe', // stdio must be set to `pipe` as the gRPC server address is read from stdout
      }
    );

    this.logProcessExitAndError('shared process', this.sharedProcess);

    const loggerService = createFileLoggerService({
      dev: this.settings.dev,
      dir: this.settings.logsDir,
      name: SHARED_PROCESS_LOGGER_NAME,
      loggerNameColor: LoggerColor.Yellow,
      passThroughMode: true,
    });

    this.sharedProcessLastLogs = new KeepLastChunks(NO_OF_LAST_LOGS_KEPT);
    loggerService.pipeProcessOutputIntoLogger(
      this.sharedProcess,
      this.sharedProcessLastLogs
    );
  }

  private initResolvingChildProcessAddresses(): void {
    this.resolvedChildProcessAddresses = Promise.all([
      resolveNetworkAddress(
        this.settings.tshd.requestedNetworkAddress,
        this.tshdProcess
      ).catch(
        rewrapResolveError(
          this.logger,
          this.settings,
          'the tsh daemon',
          TSHD_LOGGER_NAME,
          this.tshdProcessLastLogs
        )
      ),
      resolveNetworkAddress(
        this.settings.sharedProcess.requestedNetworkAddress,
        this.sharedProcess
      ).catch(
        rewrapResolveError(
          this.logger,
          this.settings,
          'the shared helper process',
          SHARED_PROCESS_LOGGER_NAME,
          this.sharedProcessLastLogs
        )
      ),
    ]).then(([tsh, shared]) => ({ tsh, shared }));
  }

  private initIpc() {
    ipcMain.on(MainProcessIpc.GetRuntimeSettings, event => {
      event.returnValue = this.settings;
    });

    ipcMain.on('main-process-should-use-dark-colors', event => {
      event.returnValue = nativeTheme.shouldUseDarkColors;
    });

    ipcMain.handle('main-process-get-resolved-child-process-addresses', () => {
      return this.resolvedChildProcessAddresses;
    });

    // the handler can remove a single kube config file or entire directory for given cluster
    ipcMain.handle(
      'main-process-remove-kube-config',
      (
        _,
        options: {
          relativePath: string;
          isDirectory?: boolean;
        }
      ) => {
        const { kubeConfigsDir } = this.settings;
        const filePath = path.join(kubeConfigsDir, options.relativePath);
        const isOutOfRoot = filePath.indexOf(kubeConfigsDir) !== 0;

        if (isOutOfRoot) {
          return Promise.reject('Invalid path');
        }
        return fs
          .rm(filePath, { recursive: !!options.isDirectory })
          .catch(error => {
            if (error.code !== 'ENOENT') {
              throw error;
            }
          });
      }
    );

    ipcMain.handle('main-process-show-file-save-dialog', (_, filePath) =>
      dialog.showSaveDialog({
        defaultPath: path.basename(filePath),
      })
    );

    ipcMain.handle(
      MainProcessIpc.SaveTextToFile,
      async (
        _,
        {
          text,
          defaultBasename,
        }: Parameters<MainProcessClient['saveTextToFile']>[0]
      ): ReturnType<MainProcessClient['saveTextToFile']> => {
        const { canceled, filePath } = await dialog.showSaveDialog({
          // Don't trust the renderer and make sure defaultBasename is indeed a basename only.
          // defaultPath accepts different kinds of paths. For security reasons, the renderer should
          // not be able to influence _where_ the file is saved, only how the file is named.
          defaultPath: path.basename(defaultBasename),
          properties: [
            'createDirectory', // macOS only
            'showOverwriteConfirmation', // Linux only.
          ],
        });

        if (canceled) {
          return { canceled };
        }

        try {
          await fs.writeFile(filePath, text, {
            // Overwrite file.
            flag: 'w',
          });
        } catch (error) {
          // Log the original error on this side of the context bridge.
          this.logger.error(`Could not save text to "${filePath}"`, error);
          throw error;
        }

        return { canceled };
      }
    );

    ipcMain.handle(
      MainProcessIpc.ForceFocusWindow,
      async (
        _,
        args?: Parameters<MainProcessClient['forceFocusWindow']>[0]
      ) => {
        if (args?.wait && args.signal?.aborted) {
          return;
        }

        this.windowsManager.forceFocusWindow();

        if (args?.wait) {
          await this.windowsManager.waitForWindowFocus(args.signal);
        }
      }
    );

    // Used in the `tsh install` command on macOS to make the bundled tsh available in PATH.
    // Returns true if tsh got successfully installed, false if the user closed the osascript
    // prompt. Throws an error when osascript fails.
    ipcMain.handle('main-process-symlink-tsh-macos', async () => {
      const source = this.settings.tshd.binaryPath;
      const target = '/usr/local/bin/tsh';
      const prompt =
        'Teleport Connect wants to create a symlink for tsh in /usr/local/bin.';
      const command = `osascript -e "do shell script \\"mkdir -p /usr/local/bin && ln -sf '${source}' '${target}'\\" with prompt \\"${prompt}\\" with administrator privileges"`;

      try {
        await promisify(exec)(command);
        this.logger.info(`Created the symlink to ${source} under ${target}`);
        return true;
      } catch (error) {
        // Ignore the error if the user canceled the prompt.
        // https://developer.apple.com/library/archive/documentation/AppleScript/Conceptual/AppleScriptLangGuide/reference/ASLR_error_codes.html#//apple_ref/doc/uid/TP40000983-CH220-SW2
        if (error instanceof Error && error.message.includes('-128')) {
          return false;
        }
        this.logger.error(error);
        throw error;
      }
    });

    ipcMain.handle('main-process-remove-tsh-symlink-macos', async () => {
      const target = '/usr/local/bin/tsh';
      const prompt =
        'Teleport Connect wants to remove a symlink for tsh from /usr/local/bin.';
      const command = `osascript -e "do shell script \\"rm '${target}'\\" with prompt \\"${prompt}\\" with administrator privileges"`;

      try {
        await promisify(exec)(command);
        this.logger.info(`Removed the symlink under ${target}`);
        return true;
      } catch (error) {
        // Ignore the error if the user canceled the prompt.
        // https://developer.apple.com/library/archive/documentation/AppleScript/Conceptual/AppleScriptLangGuide/reference/ASLR_error_codes.html#//apple_ref/doc/uid/TP40000983-CH220-SW2
        if (error instanceof Error && error.message.includes('-128')) {
          return false;
        }
        this.logger.error(error);
        throw error;
      }
    });

    ipcMain.handle('main-process-open-config-file', async () => {
      const path = this.configFileStorage.getFilePath();
      await shell.openPath(path);
      return path;
    });

    ipcMain.handle(MainProcessIpc.DownloadConnectMyComputerAgent, () =>
      this.downloadAgentShared()
    );

    ipcMain.handle(MainProcessIpc.VerifyConnectMyComputerAgent, async () => {
      await verifyAgent(this.settings.agentBinaryPath);
    });

    ipcMain.handle(
      'main-process-connect-my-computer-create-agent-config-file',
      (_, args: CreateAgentConfigFileArgs) =>
        createAgentConfigFile(this.settings, {
          proxy: args.proxy,
          token: args.token,
          rootClusterUri: args.rootClusterUri,
          username: args.username,
        })
    );

    ipcMain.handle(
      'main-process-connect-my-computer-is-agent-config-file-created',
      async (
        _,
        args: {
          rootClusterUri: RootClusterUri;
        }
      ) => isAgentConfigFileCreated(this.settings, args.rootClusterUri)
    );

    ipcMain.handle(
      'main-process-connect-my-computer-kill-agent',
      async (
        _,
        args: {
          rootClusterUri: RootClusterUri;
        }
      ) => {
        await this.agentRunner.kill(args.rootClusterUri);
      }
    );

    ipcMain.handle(
      'main-process-connect-my-computer-remove-agent-directory',
      (
        _,
        args: {
          rootClusterUri: RootClusterUri;
        }
      ) => removeAgentDirectory(this.settings, args.rootClusterUri)
    );

    ipcMain.handle(MainProcessIpc.TryRemoveConnectMyComputerAgentBinary, () =>
      this.agentRunner.tryRemoveAgentBinary()
    );

    ipcMain.handle(
      'main-process-connect-my-computer-run-agent',
      async (
        _,
        args: {
          rootClusterUri: RootClusterUri;
        }
      ) => {
        await this.agentRunner.start(args.rootClusterUri);
      }
    );

    ipcMain.on(
      'main-process-connect-my-computer-get-agent-state',
      (
        event,
        args: {
          rootClusterUri: RootClusterUri;
        }
      ) => {
        event.returnValue = this.agentRunner.getState(args.rootClusterUri);
      }
    );

    ipcMain.on(
      'main-process-connect-my-computer-get-agent-logs',
      (
        event,
        args: {
          rootClusterUri: RootClusterUri;
        }
      ) => {
        event.returnValue = this.agentRunner.getLogs(args.rootClusterUri);
      }
    );

    ipcMain.handle(
      'main-process-open-agent-logs-directory',
      async (
        _,
        args: {
          rootClusterUri: RootClusterUri;
        }
      ) => {
        const { logsDirectory } = generateAgentConfigPaths(
          this.settings,
          args.rootClusterUri
        );
        const error = await shell.openPath(logsDirectory);
        if (error) {
          throw new Error(error);
        }
      }
    );

    ipcMain.handle(
      MainProcessIpc.SelectDirectoryForDesktopSession,
      async (_, args: { desktopUri: string; login: string }) => {
        const value = await dialog.showOpenDialog({
          properties: ['openDirectory'],
        });
        if (value.canceled) {
          throw new AbortError();
        }
        if (value.filePaths.length !== 1) {
          throw new Error('No directory selected.');
        }

        const [dirPath] = value.filePaths;
        const { terminalService } = await this.getTshdClients();
        await terminalService.setSharedDirectoryForDesktopSession({
          desktopUri: args.desktopUri,
          login: args.login,
          path: dirPath,
        });

        return path.basename(dirPath);
      }
    );

    subscribeToTerminalContextMenuEvent(this.configService);
    subscribeToTabContextMenuEvent(
      this.settings.availableShells,
      this.configService
    );
    subscribeToConfigServiceEvents(this.configService);
    subscribeToFileStorageEvents(this.appStateFileStorage);
  }

  private setAppMenu() {
    const isMac = this.settings.platform === 'darwin';
    const commonHelpTemplate: MenuItemConstructorOptions[] = [
      { label: 'Open Documentation', click: openDocsUrl },
      {
        label: 'Open Logs Directory',
        click: () => openLogsDirectory(this.settings),
      },
    ];

    // Enable actions like reload or toggle dev tools only in debug mode.
    const viewMenuTemplate: MenuItemConstructorOptions = this.settings.debug
      ? { role: 'viewMenu' }
      : {
          label: 'View',
          submenu: [
            { role: 'resetZoom' },
            { role: 'zoomIn' },
            { role: 'zoomOut' },
            { type: 'separator' },
            { role: 'togglefullscreen' },
          ],
        };

    const macTemplate: MenuItemConstructorOptions[] = [
      { role: 'appMenu' },
      { role: 'editMenu' },
      viewMenuTemplate,
      {
        label: 'Window',
        submenu: [{ role: 'minimize' }, { role: 'zoom' }],
      },
      {
        role: 'help',
        submenu: commonHelpTemplate,
      },
    ];

    const otherTemplate: MenuItemConstructorOptions[] = [
      { role: 'fileMenu' },
      { role: 'editMenu' },
      viewMenuTemplate,
      { role: 'windowMenu' },
      {
        role: 'help',
        submenu: [
          ...commonHelpTemplate,
          { type: 'separator' },
          { role: 'about' },
        ],
      },
    ];

    const menu = Menu.buildFromTemplate(isMac ? macTemplate : otherTemplate);
    Menu.setApplicationMenu(menu);
  }

  private updateAboutPanelIfNeeded(): void {
    // On Windows and Linux default menu does not show copyrights.
    if (
      this.settings.platform === 'linux' ||
      this.settings.platform === 'win32'
    ) {
      app.setAboutPanelOptions({
        applicationName: app.getName(),
        applicationVersion: app.getVersion(),
        iconPath: getAssetPath('icon-linux/512x512.png'), //.ico is not supported
        copyright: `Copyright © ${new Date().getFullYear()} Gravitational, Inc.`,
      });
    }
  }

  /**
   * On Windows, where POSIX signals do not exist, the only way to gracefully
   * kill a process is to send Ctrl-Break to its console. This task is done by
   * `tsh daemon stop` program. On Unix, the standard `SIGTERM` signal is sent.
   */
  private gracefullyKillTshdProcess() {
    if (this.settings.platform !== 'win32') {
      this.tshdProcess.kill('SIGTERM');
      return;
    }

    const logger = new Logger('Daemon stop');
    const daemonStop = spawn(
      this.settings.tshd.binaryPath,
      ['daemon', 'stop', `--pid=${this.tshdProcess.pid}`],
      {
        windowsHide: true,
        timeout: 2_000,
      }
    );
    daemonStop.on('error', error => {
      logger.error('daemon stop process failed to start', error);
    });
    daemonStop.stderr.setEncoding('utf-8');
    daemonStop.stderr.on('data', logger.error);
  }

  private logProcessExitAndError(
    processName: string,
    childProcess: ChildProcess
  ) {
    childProcess.on('error', error => {
      this.logger.error(`${processName} failed to start`, error);
    });

    childProcess.once('exit', (code, signal) => {
      const codeOrSignal = [
        // code can be 0, so we cannot just check it the same way as the signal.
        code != null && `code ${code}`,
        signal && `signal ${signal}`,
      ]
        .filter(Boolean)
        .join(' ');

      this.logger.info(`${processName} exited with ${codeOrSignal}`);
    });
  }
}

const TSHD_LOGGER_NAME = 'tshd';
const SHARED_PROCESS_LOGGER_NAME = 'shared';
const DOCS_URL =
  'https://goteleport.com/docs/connect-your-client/teleport-connect/';

function openDocsUrl() {
  shell.openExternal(DOCS_URL);
}

function openLogsDirectory(settings: RuntimeSettings) {
  shell.openPath(settings.logsDir);
}

/** Shares promise returned from `promiseFn` across multiple concurrent callers. */
function sharePromise<T>(promiseFn: () => Promise<T>): () => Promise<T> {
  let pending: Promise<T> | undefined = undefined;

  return () => {
    if (!pending) {
      pending = promiseFn();
      pending.finally(() => {
        pending = undefined;
      });
    }
    return pending;
  };
}

/**
 * Sets up the gRPC clients for tsh daemon used in the main process.
 */
async function setUpTshdClients({
  runtimeSettings,
  tshdAddress,
}: {
  runtimeSettings: RuntimeSettings;
  tshdAddress: string;
}): Promise<{
  terminalService: TshdClient;
  autoUpdateService: AutoUpdateClient;
}> {
  const creds = await createGrpcCredentials(runtimeSettings);
  const transport = new GrpcTransport({
    host: tshdAddress,
    channelCredentials: creds,
    interceptors: [loggingInterceptor(new Logger('tshd'))],
  });
  return {
    terminalService: createTshdClient(transport),
    autoUpdateService: createAutoUpdateClient(transport),
  };
}

async function createGrpcCredentials(
  runtimeSettings: RuntimeSettings
): Promise<ChannelCredentials> {
  if (!grpcCreds.shouldEncryptConnection(runtimeSettings)) {
    return grpcCreds.createInsecureClientCredentials();
  }

  const { certsDir } = runtimeSettings;
  const [mainProcessKeyPair, tshdCert] = await Promise.all([
    grpcCreds.generateAndSaveGrpcCert(
      certsDir,
      grpcCreds.GrpcCertName.MainProcess
    ),
    grpcCreds.readGrpcCert(certsDir, grpcCreds.GrpcCertName.Tshd),
    // tsh daemon expects both certs to be created before accepting connections. So even though the
    // main process does not use the cert of the renderer process, it must still wait for the cert
    // to be saved to disk.
    grpcCreds.readGrpcCert(certsDir, grpcCreds.GrpcCertName.Renderer),
  ]);

  return grpcCreds.createClientCredentials(mainProcessKeyPair, tshdCert);
}

// The number of lines was chosen by looking at logs from the shared process when the glibc version
// is too old and making sure that we'd have been able to see the actual issue in the error dialog.
// See the PR description for the logs: https://github.com/gravitational/teleport/pull/38724
const NO_OF_LAST_LOGS_KEPT = 25;

function rewrapResolveError(
  logger: Logger,
  runtimeSettings: RuntimeSettings,
  processName: string,
  processLoggerName: string,
  processLastLogs: KeepLastChunks<string>
) {
  return (error: unknown) => {
    if (!(error instanceof ResolveError)) {
      throw error;
    }

    // Log the original error for full address.
    logger.error(error);

    const logPath = path.join(
      runtimeSettings.logsDir,
      `${processLoggerName}.log`
    );
    // TODO(ravicious): It'd have been ideal to get the logs from the file instead of keeping last n
    // lines in memory.
    // We tried to use winston.Logger.prototype.query for that but it kept returning no results,
    // even with structured logging turned on.
    const lastLogs = processLastLogs.getChunks().join('\n');

    throw new Error(
      `Could not communicate with ${processName}.\n\n` +
        `Last logs from ${logPath}:\n${lastLogs}`
    );
  };
}

const APP_UPDATER_STATE_KEY = 'appUpdater';
function makeAppUpdaterStorage(fs: FileStorage): AppUpdaterStorage {
  return {
    get: () => {
      const state = fs.get(APP_UPDATER_STATE_KEY) as object;
      return { managingClusterUri: '', ...state };
    },
    put: value => {
      const state = fs.get(APP_UPDATER_STATE_KEY) as object;
      fs.put(APP_UPDATER_STATE_KEY, { ...state, ...value });
    },
  };
}
