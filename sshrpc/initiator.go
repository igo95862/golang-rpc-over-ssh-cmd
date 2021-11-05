/* SPDX-License-Identifier: Apache-2.0 */
/*
 * Copyright 2021 igo95862
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

// sshrpc package implements the go RPC client and server over the
// SSH binary installed on the system. This means that SSH configuration
// from $HOME/.ssh/config will be used to connect and authenticate.
//
// Upon connecting the initator expects a matching executable to
// be invoked. See example.
//
// Communication happens over STDIN and STDOUT. Make sure that
// receiving program does not write anything to those channels.
//
// STDERR from remote end will be forwarded to local STDERR.
// This can be used for logging.
package sshrpc

import (
	"io"
	"net/rpc"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

type sshReadWriteCloser struct {
	cmd *exec.Cmd

	stdinPipe  io.WriteCloser
	stdinMutex sync.Mutex

	stdoutPipe  io.ReadCloser
	stdoutMutex sync.Mutex
}

func (s *sshReadWriteCloser) Read(p []byte) (int, error) {
	s.stdoutMutex.Lock()
	defer s.stdoutMutex.Unlock()

	return s.stdoutPipe.Read(p)
}

func (s *sshReadWriteCloser) Write(p []byte) (int, error) {
	s.stdinMutex.Lock()
	defer s.stdinMutex.Unlock()

	return s.stdinPipe.Write(p)
}

func (s *sshReadWriteCloser) Close() error {
	if err := s.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		return err
	}

	if err := s.cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func openSshConnection(destination string, arguments ...string) (*sshReadWriteCloser, error) {
	cmd := exec.Command("ssh", append([]string{destination}, arguments...)...)

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	return &sshReadWriteCloser{
		cmd:        cmd,
		stdoutPipe: stdoutPipe,
		stdinPipe:  stdinPipe,
	}, nil
}

// Open new RPC client over SSH connection.
// First argument is the destination for SSH client to connect to.
// Second argument are the arguments to invoke. They should start-up the receiver
// application.
func NewSshRpcClient(destination string, arguments ...string) (*rpc.Client, error) {

	sshConnection, err := openSshConnection(destination, arguments...)
	if err != nil {
		return nil, err
	}

	client := rpc.NewClient(sshConnection)

	return client, nil
}
