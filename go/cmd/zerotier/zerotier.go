/*
 * Copyright (c)2019 ZeroTier, Inc.
 *
 * Use of this software is governed by the Business Source License included
 * in the LICENSE.TXT file in the project's root directory.
 *
 * Change Date: 2023-01-01
 *
 * On the date above, in accordance with the Business Source License, use
 * of this software will be governed by version 2.0 of the Apache License.
 */
/****/

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
	"zerotier/cmd/zerotier/cli"
	"zerotier/pkg/zerotier"
)

func readAuthToken(basePath string) string {
	data, _ := ioutil.ReadFile(path.Join(basePath, "authtoken.secret"))
	if len(data) > 0 {
		return string(data)
	}
	userHome, _ := os.UserHomeDir()
	if len(userHome) > 0 {
		if runtime.GOOS == "darwin" {
			data, _ = ioutil.ReadFile(userHome + "/Library/Application Support/ZeroTier/authtoken.secret")
			if len(data) > 0 {
				return string(data)
			}
			data, _ = ioutil.ReadFile(userHome + "/Library/Application Support/ZeroTier/One/authtoken.secret")
			if len(data) > 0 {
				return string(data)
			}
		}
		data, _ = ioutil.ReadFile(path.Join(userHome, ".zerotierauth"))
		if len(data) > 0 {
			return string(data)
		}
		data, _ = ioutil.ReadFile(path.Join(userHome, ".zeroTierOneAuthToken"))
		if len(data) > 0 {
			return string(data)
		}
	}
	return ""
}

func main() {
	globalOpts := flag.NewFlagSet("global", flag.ContinueOnError)
	hflag := globalOpts.Bool("h", false, "") // support -h to be canonical with other Unix utilities
	jflag := globalOpts.Bool("j", false, "")
	pflag := globalOpts.String("p", "", "")
	tflag := globalOpts.String("t", "", "")
	err := globalOpts.Parse(os.Args[1:])
	if err != nil {
		cli.Help()
		os.Exit(1)
		return
	}
	args := globalOpts.Args()
	if len(args) < 1 || *hflag {
		cli.Help()
		os.Exit(0)
		return
	}
	var cmdArgs []string
	if len(args) > 1 {
		cmdArgs = args[1:]
	}

	basePath := zerotier.PlatformDefaultHomePath
	if len(*pflag) > 0 {
		basePath = *pflag
	}
	var authToken string
	if len(*tflag) > 0 {
		authToken = *tflag
	} else {
		authToken = readAuthToken(basePath)
	}
	if len(authToken) == 0 {
		fmt.Println("FATAL: unable to read API authorization token from service path or user home ('sudo' may be needed)")
		os.Exit(1)
	}
	authToken = strings.TrimSpace(authToken)

	switch args[0] {
	case "help":
		cli.Help()
		os.Exit(0)
	case "version":
		fmt.Printf("%d.%d.%d\n", zerotier.CoreVersionMajor, zerotier.CoreVersionMinor, zerotier.CoreVersionRevision)
		os.Exit(0)
	case "service":
		cli.Service(basePath, authToken, cmdArgs)
	case "status":
		cli.Status(basePath, authToken, cmdArgs, *jflag)
	case "peers":
		cli.Peers(basePath, authToken, cmdArgs)
	case "roots":
		cli.Roots(basePath, authToken, cmdArgs)
	case "addroot":
		cli.AddRoot(basePath, authToken, cmdArgs)
	case "removeroot":
		cli.RemoveRoot(basePath, authToken, cmdArgs)
	case "networks":
		cli.Networks(basePath, authToken, cmdArgs)
	case "join":
		cli.Join(basePath, authToken, cmdArgs)
	case "leave":
		cli.Leave(basePath, authToken, cmdArgs)
	case "show":
		cli.Show(basePath, authToken, cmdArgs)
	case "set":
		cli.Set(basePath, authToken, cmdArgs)
	}

	cli.Help()
	os.Exit(1)
}