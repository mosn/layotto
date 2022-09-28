/*
 * Copyright 2022 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dockercompose

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"gotest.tools/icmd"
)

type Command struct {
	Command    string
	Args       []string
	WorkingDir string
	Env        map[string]string
}

func runDockerComposeCmd(tb testing.TB, args []string, workingDir string,
	envVas map[string]string) error {

	command := Command{
		WorkingDir: workingDir,
		Env:        envVas,
	}

	dockerComposeVersionCmd := icmd.Command("docker", "compose", "version")

	result := icmd.RunCmd(dockerComposeVersionCmd)

	if result.ExitCode == 0 {
		command.Command = "docker"
		command.Args = append([]string{
			"compose", "--project-name", strings.ToLower(tb.Name()),
		}, args...)
	} else {
		command.Command = "docker-compose"
		command.Args = append([]string{
			"--project-name", strings.ToLower(tb.Name()),
		}, args...)
	}

	return runCmd(tb, command)
}

func runCmd(tb testing.TB, command Command) error {
	tb.Logf("Running command %s with args %s", command.Command, command.Args)

	cmd := exec.Command(command.Command, command.Args...)
	cmd.Dir = command.WorkingDir
	cmd.Stdin = os.Stdin
	cmd.Env = formatEnvVars(command)

	return cmd.Run()
}

func formatEnvVars(command Command) []string {
	env := os.Environ()
	for key, value := range command.Env {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	return env
}
