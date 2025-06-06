//go:build linux

/*
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

package botfs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/user"
	"strconv"
	"sync"
	"syscall"

	"github.com/coreos/go-semver/semver"
	"github.com/gravitational/trace"
	"github.com/joshlf/go-acl"
	"golang.org/x/sys/unix"

	"github.com/gravitational/teleport/lib/utils"
)

// mostACLRead is a permission mode granting readonly access to a file.
const modeACLRead fs.FileMode = 04

// modeACLReadExecute is the permissions mode needed for read on directories.
const modeACLReadExecute fs.FileMode = 05

// modeACLReadWrite is the lower 3 bytes of a UNIX file mode for permission
// bits, i.e. just one r/w/x.
const modeACLReadWrite fs.FileMode = 06

// modeACLReadWriteExecute is the permissions mode needed for full rwx on
// directories.
const modeACLReadWriteExecute fs.FileMode = 07

// modeACLNone is the UNIX file mode for no permissions, used for group and
// world read/write.
const modeACLNone fs.FileMode = 0

// missingSyscallWarning is used to reduce log spam when a syscall is missing.
var missingSyscallWarning sync.Once

// openSecure opens the given path for either reading or writing with the
// RESOLVE_NO_SYMLINKS flag set.
func openSecure(path string, flags OpenFlags) (*os.File, error) {
	var mode uint64
	if flags != ReadFlags {
		// openat2() with a nonzero mode will raise EINVAL unless O_CREATE or
		// O_TMPFILE is set, so only set this in non-read mode.
		mode = uint64(DefaultMode.Perm())
	}

	how := unix.OpenHow{
		Mode:    mode,
		Flags:   uint64(flags.Flags()) | unix.O_CLOEXEC,
		Resolve: unix.RESOLVE_NO_SYMLINKS,
	}

	for {
		fd, err := unix.Openat2(unix.AT_FDCWD, path, &how)
		if errors.Is(err, syscall.EINTR) {
			// Per the stdlib's implementation, EINTR errors should be ignored
			// and retried.
			continue
		} else if err != nil {
			// note: returning the original error here for comparison purposes
			return nil, err
		}

		file := os.NewFile(uintptr(fd), path)
		if file == nil {
			// Probably useless since this implies the fd itself is invalid, but
			// attempt to close the fd anyway.
			_ = unix.Close(fd)
			return nil, os.ErrInvalid
		}

		return file, nil
	}
}

// openSymlinks mode opens the file for read or write using the given symlink
// mode, potentially failing or logging a warning if symlinks can't be
// secured.
func openSymlinksMode(path string, flags OpenFlags, symlinksMode SymlinksMode) (*os.File, error) {
	var file *os.File
	var err error

	switch symlinksMode {
	case SymlinksSecure:
		file, err = openSecure(path, flags)
		if errors.Is(err, unix.ENOSYS) {
			return nil, trace.Errorf("openSecure failed due to missing syscall; configure `symlinks: insecure` for %q", path)
		} else if err != nil {
			return nil, trace.Wrap(err)
		}
	case SymlinksTrySecure:
		file, err = openSecure(path, flags)
		if errors.Is(err, unix.ENOSYS) {
			missingSyscallWarning.Do(func() {
				log.DebugContext(
					context.TODO(),
					"Failed to open file securely due to missing syscall; falling back to regular file handling. Configure `symlinks: insecure` to disable this warning",
					"path", path,
				)
			})

			file, err = openStandard(path, flags)
			if err != nil {
				return nil, trace.Wrap(err)
			}
		} else if err != nil {
			return nil, trace.Wrap(err)
		}
	case SymlinksInsecure:
		file, err = openStandard(path, flags)
		if err != nil {
			return nil, trace.Wrap(err)
		}
	default:
		return nil, trace.BadParameter("invalid symlinks mode %q", symlinksMode)
	}

	return file, nil
}

// createStandard creates an empty file or directory at the given path while
// attempting to prevent symlink attacks.
func createSecure(path string, isDir bool) error {
	if isDir {
		// We can't specify RESOLVE_NO_SYMLINKS for mkdir. This isn't the end
		// of the world, though: if an attacker attempts a symlink attack we'll
		// just open the correct file for read/write later (and error when it
		// doesn't exist).
		if err := os.Mkdir(path, DefaultDirMode); err != nil {
			return trace.Wrap(err)
		}

		return nil
	}

	f, err := openSecure(path, WriteFlags)
	if errors.Is(err, unix.ENOSYS) {
		// bubble up the original error for comparison
		return err
	} else if err != nil {
		return trace.Wrap(err)
	}

	// No writing to do, just close it.
	if err := f.Close(); err != nil {
		log.WarnContext(
			context.TODO(),
			"Failed to close file",
			"path", path,
			"error", err,
		)
	}

	return nil
}

// Create attempts to create the given file or directory with the given
// symlinks mode.
func Create(path string, isDir bool, symlinksMode SymlinksMode) error {
	// Implementation note: paranoid file _creation_ is only really useful for
	// providing an early warning if openat2() / ACLs are unsupported on the
	// host system, as it will catch compatibility issues during `tbot init`.
	// Read() and Write() with Symlinks(Try)Secure are the codepaths that
	// actually prevents symlink attacks.

	switch symlinksMode {
	case SymlinksSecure:
		if err := createSecure(path, isDir); err != nil {
			if errors.Is(err, unix.ENOSYS) {
				return trace.Errorf("createSecure failed due to missing syscall; configure `symlinks: insecure` for %q", path)
			}

			return trace.Wrap(err)
		}
	case SymlinksTrySecure:
		err := createSecure(path, isDir)
		if err == nil {
			// All good, move on.
			return nil
		}

		if !errors.Is(err, unix.ENOSYS) {
			// Something else went wrong, fail.
			return trace.Wrap(err)
		}

		// It's a bit gross to stuff this sync.Once into a global, but
		// hopefully that's forgivable since it just manages a log message.
		missingSyscallWarning.Do(func() {
			log.WarnContext(
				context.TODO(),
				"Failed to create file securely due to missing syscall; falling back to regular file handling. Configure `symlinks: insecure` to disable this warning",
				"path", path,
			)
		})

		return trace.Wrap(createStandard(path, isDir))
	case SymlinksInsecure:
		return trace.Wrap(createStandard(path, isDir))
	default:
		return trace.BadParameter("invalid symlinks mode %q", symlinksMode)
	}

	return nil
}

// Read reads the contents of the given file into memory.
func Read(path string, symlinksMode SymlinksMode) ([]byte, error) {
	file, err := openSymlinksMode(path, ReadFlags, symlinksMode)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return data, nil
}

// Write stores the given data to the file at the given path.
func Write(path string, data []byte, symlinksMode SymlinksMode) error {
	file, err := openSymlinksMode(path, WriteFlags, symlinksMode)
	if err != nil {
		return trace.Wrap(err)
	}

	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return trace.Wrap(err)
	}

	return nil
}

// desiredPerms determines the desired bot permissions for an artifact at
// the given path. Directories require read/exec, files require read/write.
func desiredPerms(path string) (ownerMode fs.FileMode, botAndReaderMode fs.FileMode, err error) {
	stat, err := os.Stat(path)
	if err != nil {
		return 0, 0, trace.Wrap(err)
	}

	botAndReaderMode = modeACLReadWrite
	ownerMode = modeACLReadWrite
	if stat.IsDir() {
		botAndReaderMode = modeACLReadExecute
		ownerMode = modeACLReadWriteExecute
	}

	return
}

// VerifyLegacyACL verifies whether the ACL of the given file allows writes from
// the bot user. Errors may optionally be used as more informational warnings;
// ConfigureLegacyACL can be used to correct them, assuming the user has
// permissions.
func VerifyLegacyACL(path string, opts *ACLOptions) error {
	current, err := acl.Get(path)
	if err != nil {
		return trace.Wrap(err)
	}

	stat, err := os.Stat(path)
	if err != nil {
		return trace.Wrap(err)
	}

	owner, err := GetOwner(stat)
	if err != nil {
		return trace.Wrap(err)
	}

	ownerMode, botAndReaderMode, err := desiredPerms(path)
	if err != nil {
		return trace.Wrap(err)
	}

	// Attempt to determine the max # of expected user tags. We can't know the
	// reader user in all cases because we only know it during `tbot init`, so
	// instead we'll try to determine the upper bound of expected entries here.
	maxExpectedUserTags := 2
	if owner.Uid == opts.BotUser.Uid {
		// This path is owned by the bot user, so at least one user tag will
		// be missing.
		maxExpectedUserTags--
	}

	// Also determine the minimum number of expected user tags. There should
	// generally be at least 1 unless all users are the same.
	minExpectedUserTags := 0
	if owner.Uid != opts.BotUser.Uid {
		minExpectedUserTags++
	}
	if opts.ReaderUser != nil && owner.Uid != opts.ReaderUser.Uid {
		minExpectedUserTags++
	}

	userTagCount := 0
	errors := []error{}

	for _, entry := range current {
		switch entry.Tag {
		case acl.TagUserObj:
			if entry.Perms != ownerMode {
				errors = append(errors, trace.BadParameter("user entry has improper permissions %d", entry.Perms))
			}
		case acl.TagGroupObj:
			if entry.Perms != modeACLNone {
				errors = append(errors, trace.BadParameter("group entry has improper permissions %d", entry.Perms))
			}
		case acl.TagOther:
			if entry.Perms != modeACLNone {
				errors = append(errors, trace.BadParameter("other entry has improper permissions %d", entry.Perms))
			}
		case acl.TagMask:
			if entry.Perms != botAndReaderMode {
				errors = append(errors, trace.BadParameter("mask entry has improper permissions %d", entry.Perms))
			}
		case acl.TagGroup:
			// Group tags currently not allowed.
			errors = append(errors, trace.BadParameter("unexpected group entry found"))
		case acl.TagUser:
			userTagCount++

			// It's only worth checking the qualifiers if we know all expected
			// values in advance. We can't know them at bot runtime in any way
			// that isn't also subject to e.g. an attacker with root / owner
			// access, so we might as well not spam warnings every time we
			// verify the ACLs. We'll have to rely on the tag counter instead.
			if opts.BotUser != nil &&
				opts.ReaderUser != nil &&
				entry.Qualifier != opts.BotUser.Uid &&
				entry.Qualifier != opts.ReaderUser.Uid {
				errors = append(errors, trace.BadParameter("invalid qualifier %q for user entry", entry.Qualifier))
			}

			if entry.Perms != botAndReaderMode {
				errors = append(errors, trace.BadParameter("invalid permissions %q for bot user entry", entry.Perms))
			}
		}
	}

	if userTagCount > maxExpectedUserTags {
		errors = append(errors, trace.BadParameter("too many user tags found"))
	} else if userTagCount < minExpectedUserTags {
		errors = append(errors, trace.BadParameter("too few user tags found"))
	}

	return trace.NewAggregate(errors...)
}

// ConfigureLegacyACL configures ACLs of the given file to allow writes from the bot
// user.
func ConfigureLegacyACL(path string, owner *user.User, opts *ACLOptions) error {
	if owner.Uid == opts.BotUser.Uid && owner.Uid == opts.ReaderUser.Uid {
		// We'll end up with an empty ACL. This isn't technically a problem
		log.WarnContext(
			context.TODO(),
			"The owner, bot, and reader all appear to be the same user. This is an unusual configuration: consider setting `acls: off` in the destination config to remove this warning",
			"username", owner.Username,
		)
	}

	// We fully specify the ACL here to ensure the correct permissions are
	// always set (rather than just appending an "allow" for the bot user).
	// Note: These need to be sorted by tag value.
	ownerMode, botAndReaderMode, err := desiredPerms(path)
	if err != nil {
		return trace.Wrap(err)
	}

	// Note: ACL entries need to be ordered per acl_linux entryPriority
	desiredACL := acl.ACL{
		{
			// Note: Mask does not apply to user object.
			Tag:   acl.TagUserObj,
			Perms: ownerMode,
		},
	}

	// Only add an entry for the bot user if it isn't the owner.
	if owner.Uid != opts.BotUser.Uid {
		desiredACL = append(desiredACL, acl.Entry{
			// Entry to allow the bot to read/write.
			Tag:       acl.TagUser,
			Qualifier: opts.BotUser.Uid,
			Perms:     botAndReaderMode,
		})
	}

	// Only add entry for the reader if it isn't the owner.
	if owner.Uid != opts.ReaderUser.Uid {
		desiredACL = append(desiredACL, acl.Entry{
			// Entry to allow the reader user to read/write.
			Tag:       acl.TagUser,
			Qualifier: opts.ReaderUser.Uid,
			Perms:     botAndReaderMode,
		})
	}

	desiredACL = append(desiredACL,
		acl.Entry{
			Tag:   acl.TagGroupObj,
			Perms: modeACLNone,
		},
		acl.Entry{
			// Mask is the maximum permissions the ACL can grant. This should
			// match the desired bot permissions.
			Tag:   acl.TagMask,
			Perms: botAndReaderMode,
		},
		acl.Entry{
			Tag:   acl.TagOther,
			Perms: modeACLNone,
		},
	)

	// Note: we currently give both the bot and reader read/write to all the
	// files. This is done for simplicity and shouldn't represent a security
	// risk - the bot obviously knows the contents of the destination, and
	// the files being writable to the reader is (maybe arguably) not a
	// security issue.

	log.DebugContext(
		context.TODO(),
		"Configuring legacy ACL on path",
		"path", path,
		"acl", desiredACL,
	)
	return trace.ConvertSystemError(trace.Wrap(acl.Set(path, desiredACL)))
}

// resolveACLReaderSelector attempts to convert an ACL selector into a
// platform-specific acl.Entry that can be applied to a file.
func resolveACLReaderSelector(s *ACLSelector, dir bool) (acl.Entry, error) {
	perm := modeACLRead
	if dir {
		perm = modeACLReadExecute
	}

	switch {
	case s.User != "":
		_, err := strconv.ParseInt(s.User, 10, 32)
		if err == nil {
			// User is a valid number, so return it directly. We don't need to
			// check the particular value as it's explicitly allowed to
			// configure entries for nonexistent users.
			return acl.Entry{
				Tag:       acl.TagUser,
				Qualifier: s.User,
				Perms:     perm,
			}, nil
		}

		user, err := user.Lookup(s.User)
		if err != nil {
			return acl.Entry{}, trace.Wrap(err)
		}

		return acl.Entry{
			Tag:       acl.TagUser,
			Qualifier: user.Uid,
			Perms:     perm,
		}, nil
	case s.Group != "":
		_, err := strconv.ParseInt(s.Group, 10, 32)
		if err == nil {
			// Group is a valid number, so return it directly.
			return acl.Entry{
				Tag:       acl.TagGroup,
				Qualifier: s.Group,
				Perms:     perm,
			}, nil
		}

		group, err := user.LookupGroup(s.Group)
		if err != nil {
			return acl.Entry{}, trace.Wrap(err)
		}

		return acl.Entry{
			Tag:       acl.TagGroup,
			Qualifier: group.Gid,
			Perms:     perm,
		}, nil
	default:
		return acl.Entry{}, trace.BadParameter("unable to resolve ACL selector, user or group must be specified: %+v", s)
	}
}

// aclMaskForSelectors returns an appropriate ACL mask entry given the list of
// selectors.
func aclReaderMask(dir bool) acl.Entry {
	perms := modeACLRead
	if dir {
		perms = modeACLReadExecute
	}

	return acl.Entry{
		Tag:   acl.TagMask,
		Perms: perms,
	}
}

// aclFromReaders builds an ACL for a generic file or directory from the given
// list of reader selectors.
func aclFromReaders(selectors []*ACLSelector, dir bool) (acl.ACL, error) {
	// Note: these entries are out of their required order, but go-acl sorts
	// them before applying.
	ownerPerms := modeACLReadWrite
	if dir {
		ownerPerms = modeACLReadWriteExecute
	}

	desiredACL := acl.ACL{
		{
			Tag:   acl.TagUserObj,
			Perms: ownerPerms,
		},
		{
			Tag:   acl.TagGroupObj,
			Perms: modeACLNone,
		},
		aclReaderMask(dir),
		{
			Tag:   acl.TagOther,
			Perms: modeACLNone,
		},
	}

	for _, selector := range selectors {
		entry, err := resolveACLReaderSelector(selector, dir)
		if err != nil {
			return nil, trace.Wrap(err)
		}

		desiredACL = append(desiredACL, entry)
	}

	return desiredACL, nil
}

// ConfigureACL configures a bot-user-owned ACL at the given path such that it
// can be read by the given list of readers. If the list is empty, appropriate
// non-ACL permissions will be set to ensure only the bot user can read the
// file.
func ConfigureACL(path string, readers []*ACLSelector) error {
	stat, err := os.Stat(path)
	if err != nil {
		return trace.Wrap(err)
	}

	dir := stat.IsDir()

	desiredACL, err := aclFromReaders(readers, dir)
	if err != nil {
		return trace.Wrap(err)
	}

	log.DebugContext(
		context.TODO(),
		"Configuring ACL on path",
		"path", path,
		"acl", desiredACL,
	)
	return trace.ConvertSystemError(trace.Wrap(acl.Set(path, desiredACL)))
}

// permStrings is a list of possible permission string representations,
// duplicated from `acl` as it is private.
var permStrings = []string{
	0: "---",
	1: "--x",
	2: "-w-",
	3: "-wx",
	4: "r--",
	5: "r-x",
	6: "rw-",
	7: "rwx",
}

// formatEntry formats an ACL entry without attempting to resolve UIDs or GIDs
// to ensure consistent results in tests and elsewhere. Inlined from
// acl.Entry.String().
func formatEntry(e acl.Entry) string {
	middle := "::"
	if e.Tag == acl.TagUser || e.Tag == acl.TagGroup {
		middle = ":" + e.Qualifier + ":"
	}
	return fmt.Sprintf("%s%s%s", e.Tag, middle, permStrings[7&e.Perms])
}

// CompareACL compares two ACLs to check if they match. Returns an empty list if
// the two are identity, otherwise returns a list of error messages for each
// issue found.
func CompareACL(expected, candidate acl.ACL) []string {
	expectedSet := utils.NewSet(expected...)
	candidateSet := utils.NewSet(candidate...)

	missing := expectedSet.Clone().Subtract(candidateSet)
	unexpected := candidateSet.Clone().Subtract(expectedSet)

	var issues []string
	for m := range missing {
		issues = append(issues, fmt.Sprintf("missing required entry: %s", formatEntry(m)))
	}
	for u := range unexpected {
		issues = append(issues, fmt.Sprintf("unexpected entry: %s", formatEntry(u)))
	}

	return issues
}

// VerifyACL loads the ACL for the file at the given path and compares it to
// the expected ACL as determined by the given list of reader selectors,
// returning a list of issues found. An empty list (and nil error) implies a
// valid ACL. The path must exist and the user must have permission to `os.Stat`
// it. Errors will be returned if the ACL was unable to be read.
func VerifyACL(path string, readers []*ACLSelector) ([]string, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return []string{}, trace.Wrap(err)
	}

	dir := stat.IsDir()

	expected, err := aclFromReaders(readers, dir)
	if err != nil {
		return []string{}, trace.Wrap(err)
	}

	fsACL, err := acl.Get(path)
	if err != nil {
		return []string{}, trace.Wrap(err)
	}

	issues := CompareACL(expected, fsACL)
	return issues, nil
}

// HasACLSupport determines if this binary / system supports ACLs.
func HasACLSupport() bool {
	// We just assume Linux _can_ support ACLs here, and will test for support
	// at runtime.
	return true
}

// HasSecureWriteSupport determines if `CreateSecure()` should be supported
// on this OS / kernel version. Note that it just checks the kernel version,
// so this should be treated as a fallible hint.
//
// We've encountered this being incorrect in environments where access to the
// kernel is hampered e.g. seccomp/apparmor/container runtimes.
func HasSecureWriteSupport() bool {
	version, err := utils.KernelVersion()
	if err != nil {
		log.InfoContext(
			context.TODO(),
			"Failed to determine kernel version. It will be assumed secure write support is not available",
			"error", err,
		)
		return false
	}
	// kernel release that added support for the openat2() syscall
	if version.LessThan(semver.Version{Major: 5, Minor: 6, Patch: 0}) {
		return false
	}

	return true
}

// GetOwner attempts to retrieve the owner of the given file. This is not
// supported on all platforms and will return a trace.NotImplemented in that
// case.
func GetOwner(fileInfo fs.FileInfo) (*user.User, error) {
	info, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, trace.NotImplemented("Cannot verify file ownership on this platform.")
	}

	user, err := user.LookupId(strconv.Itoa(int(info.Uid)))
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return user, nil
}

// IsOwnedBy checks that the file at the given path is owned by the given user.
// Returns a trace.NotImplemented() on unsupported platforms.
func IsOwnedBy(fileInfo fs.FileInfo, uid int) (bool, error) {
	info, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return false, trace.Errorf("unexpected type of file info on Linux: %T", fileInfo.Sys())
	}

	// Our files are 0600, so don't bother checking gid.
	return int(info.Uid) == uid, nil
}
