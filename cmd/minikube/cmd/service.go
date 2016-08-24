/*
Copyright (C) 2016 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/docker/machine/libmachine"
	"github.com/jimmidyson/minishift/pkg/minikube/cluster"
	"github.com/jimmidyson/minishift/pkg/minikube/constants"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var (
	namespace        string
	serviceURLMode   bool
	serviceURLPrefix string
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service [flags] SERVICE",
	Short: "Gets the URL for the specified service in your local cluster",
	Long:  `Gets the URL for the specified service in your local cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || len(args) > 1 {
			fmt.Fprintln(os.Stderr, "Please specify a service name.")
			os.Exit(1)
		}

		service := args[0]

		api := libmachine.NewClient(constants.Minipath, constants.MakeMiniPath("certs"))
		defer api.Close()
		url, err := cluster.GetServiceURL(api, namespace, service, serviceURLPrefix)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			if _, ok := err.(cluster.MissingNodePortError); !ok {
				fmt.Fprintln(os.Stderr, "Check that minishift is running and that you have specified the correct namespace (-n flag).")
			}
			os.Exit(1)
		}

		if serviceURLMode {
			fmt.Fprintln(os.Stdout, url)
		} else {
			fmt.Fprintln(os.Stdout, "Opening kubernetes service "+namespace+"/"+service+" in default browser...")
			browser.OpenURL(url)
		}
	},
}

func init() {
	serviceCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "The service namespace")
	serviceCmd.Flags().BoolVar(&serviceURLMode, "url", false, "Display the kubernetes service URL in the CLI instead of opening it in the default browser")
	serviceCmd.Flags().StringVar(&serviceURLPrefix, "prefix", "http://", "URL prefix to use when using --url flag")
	RootCmd.AddCommand(serviceCmd)
}
