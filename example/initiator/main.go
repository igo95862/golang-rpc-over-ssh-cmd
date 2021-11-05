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
	"flag"
	"github.com/igo95862/golang-rpc-over-ssh-cmd/sshrpc"
	"log"
	"os"
)

func main() {
	flag.Parse()

	localHostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Failed to get local hostname: ", err)
	}
	log.Print("Local hostname: ", localHostname)

	client, err := sshrpc.NewSshRpcClient(flag.Arg(0), flag.Args()[1:]...)
	if err != nil {
		log.Fatal("Failed to open SSH connection: ", err)
	}

	var reply string
	err = client.Call("TestType.GetHostname", struct{}{}, &reply)
	if err != nil {
		log.Fatal("Receiver got error: ", err)
	}
	log.Print("Remote hostname: ", reply)

}
