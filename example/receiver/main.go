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

package main

import (
	"github.com/igo95862/golang-rpc-over-ssh-cmd/sshrpc"
	"net/rpc"
	"os"
)

type TestType struct{}

func (l *TestType) GetHostname(_ struct{}, reply *string) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	*reply = hostname

	return nil
}

func main() {
	test := new(TestType)

	server := rpc.NewServer()
	err := server.Register(test)
	if err != nil {
		panic(err)
	}
	sshrpc.StartReceiving(server)
}
