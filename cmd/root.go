package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/adrg/xdg"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// GetProjects already in the projects path
func GetProjects(projectsPath string) (projects []string, err error) {
	err = filepath.Walk(projectsPath, func(cwd string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if _, err := os.Stat(path.Join(cwd, ".git")); !os.IsNotExist(err) {
				projects = append(projects, cwd[len(projectsPath)+1:])
				return filepath.SkipDir
			}
		}
		return nil
	})

	return
}

func rootRun(cmd *cobra.Command, args []string) {
	projectsPath := viper.GetString("projects_path")

	err := filepath.Walk(projectsPath, func(cwd string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if _, err := os.Stat(path.Join(cwd, ".git")); !os.IsNotExist(err) {
				fmt.Println(cwd[len(projectsPath)+1:])
				return filepath.SkipDir
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
	rootCmd.PersistentFlags().String("projects_path", "", "The root path of all projects")
	rootCmd.PersistentFlags().String("gitlab_url", "", "Base URL for the Gitlab API")
	rootCmd.PersistentFlags().String("gitlab_token", "", "User's token for consuming the Gitlab API")

	viper.BindPFlag("projects_path", rootCmd.PersistentFlags().Lookup("projects_path"))
	viper.BindPFlag("gitlab_url", rootCmd.PersistentFlags().Lookup("gitlab_url"))
	viper.BindPFlag("gitlab_token", rootCmd.PersistentFlags().Lookup("gitlab_token"))
}

func initConfig() {
	var err error

	if cfgFile := viper.GetString("config"); cfgFile == "" {
		Home, err = homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.SetConfigName("config")
		viper.AddConfigPath(path.Join(Home, ".p"))
		viper.AddConfigPath(path.Join(xdg.ConfigHome, "p"))
	} else {
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
