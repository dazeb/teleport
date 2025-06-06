/*
 * Teleport
 * Copyright (C) 2024  Gravitational, Inc.
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

package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKubeCredentialsCommand(t *testing.T) {
	testCommand(t, NewKubeCredentialsCommand, []testCommandCase[*KubeCredentialsCommand]{
		{
			name: "success",
			args: []string{
				"credentials",
				"--destination-dir=/tmp",
			},
			assert: func(t *testing.T, got *KubeCredentialsCommand) {
				require.Equal(t, "/tmp", got.DestinationDir)
			},
		},
	})
}
