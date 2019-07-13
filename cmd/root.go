/*
Copyright © 2019 RAFAEL GUMIERI <rgumieri@gmail.com>

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
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/adrg/xdg"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

func rootRun(cmd *cobra.Command, args []string) {
	gl := gitlab.NewClient(nil, viper.GetString("gitlab_token"))
	gl.SetBaseURL(viper.GetString("gitlab_url"))
	projects, _, err := gl.Projects.ListProjects(&gitlab.ListProjectsOptions{Archived: gitlab.Bool(false)})
	if err != nil {
		log.Fatal(err)
	}
	for _, project := range projects {
		fmt.Println(project.PathWithNamespace)
	}
}

var rootCmd = &cobra.Command{
	Use:   "p",
	Short: "Tool for helping the management of git projects",
	Long:  `Collection of helping commands for the management of projects using git.`,
	Run:   rootRun,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.p.yaml or $XDG_CONFIG_HOME/p/config.yaml)")
	rootCmd.PersistentFlags().String("gitlab_url", "", "Base URL for the Gitlab API")
	rootCmd.PersistentFlags().String("gitlab_token", "", "User's token for consuming the Gitlab API")

	viper.BindPFlag("gitlab_url", rootCmd.PersistentFlags().Lookup("gitlab_url"))
	viper.BindPFlag("gitlab_token", rootCmd.PersistentFlags().Lookup("gitlab_token"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile := viper.GetString("config"); cfgFile == "" {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.SetConfigName("config")
		viper.AddConfigPath(path.Join(home, ".p"))
		viper.AddConfigPath(path.Join(xdg.ConfigHome, "p"))
	} else {
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
