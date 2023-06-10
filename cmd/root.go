/*
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

package main

import (
	"os"

	"github.com/imdario/mergo"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	version = ""
)

type GlobalOptions struct {
	Verbose    bool
	Version    bool
	ConfigFile string
}

var (
	globalOpts = GlobalOptions{}
	rootCmd    = &cobra.Command{
		Use:     "go-cli-template",
		Version: version,
	}
)

func main() {
	rootCmd.PersistentFlags().BoolVar(&globalOpts.Verbose, "verbose", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVar(&globalOpts.Version, "version", false, "version")
	rootCmd.PersistentFlags().StringVarP(&globalOpts.ConfigFile, "file", "f", "", "YAML Config File")

	rootCmd.AddCommand(&cobra.Command{Use: "completion", Hidden: true})
	cobra.EnableCommandSorting = false

	lo.Must0(rootCmd.Execute())
}

func ParseConfig[T any](globalOpts GlobalOptions, opts T) (T, error) {
	if globalOpts.ConfigFile == "" {
		return opts, nil
	}
	configBytes, err := os.ReadFile(globalOpts.ConfigFile)
	if err != nil {
		return opts, err
	}
	var parsedCreateOpts T
	if err := yaml.Unmarshal(configBytes, &parsedCreateOpts); err != nil {
		return opts, err
	}
	if err := mergo.Merge(&opts, parsedCreateOpts, mergo.WithOverride); err != nil {
		return opts, err
	}
	return opts, nil
}
