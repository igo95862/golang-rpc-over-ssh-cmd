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

package sshrpc

import (
	"io"
	"net/rpc"
	"os"
	"sync"
)

type sshRecieverReadWriteCloser struct {
	stdinPipe  io.WriteCloser
	stdinMutex sync.Mutex

	stdoutPipe  io.ReadCloser
	stdoutMutex sync.Mutex
}

func createSshServer() *sshRecieverReadWriteCloser {
	return &sshRecieverReadWriteCloser{}
}

func (s *sshRecieverReadWriteCloser) Read(p []byte) (int, error) {
	s.stdoutMutex.Lock()
	defer s.stdoutMutex.Unlock()

	return os.Stdin.Read(p)
}

func (s *sshRecieverReadWriteCloser) Write(p []byte) (int, error) {
	s.stdinMutex.Lock()
	defer s.stdinMutex.Unlock()

	return os.Stdout.Write(p)
}

func (s *sshRecieverReadWriteCloser) Close() error {
	return nil
}

// This function should be called on receiving end program.
// It will block indefinitely until SSH channel is closed.
//
// DO NOT USE STDIN AND STDOUT.
//
// STDERR will be linked to the initiator's STDERR and
// can be used for logging.
//
// All types must be registered before invoking this
// function. See net/rpc documentation.
func StartReceiving(server *rpc.Server) {
	server.ServeConn(createSshServer())
}
