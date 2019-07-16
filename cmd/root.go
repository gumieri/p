package cmd

import (
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/adrg/xdg"
	typistPkg "github.com/gumieri/typist"
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

	typist.Must(filepath.Walk(projectsPath, func(cwd string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if _, err := os.Stat(path.Join(cwd, ".git")); !os.IsNotExist(err) {
				typist.Println(cwd[len(projectsPath)+1:])
				return filepath.SkipDir
			}
		}
		return nil
	}))
}

func persistentPreRun(cmd *cobra.Command, args []string) {
	typist = &typistPkg.Typist{
		Quiet: quiet,
		In:    os.Stdin,
		Out:   os.Stdout,
	}
}

var rootCmd = &cobra.Command{
	Use:              "p",
	Short:            "Tool for helping the management of git projects",
	Long:             `Collection of helping commands for the management of projects using git.`,
	Run:              rootRun,
	PersistentPreRun: persistentPreRun,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	typist.Must(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.p.yaml or $XDG_CONFIG_HOME/p/config.yaml)")
	rootCmd.PersistentFlags().String("projects_path", "", "The root path of all projects")
	rootCmd.PersistentFlags().String("gitlab_url", "", "Base URL for the Gitlab API")
	rootCmd.PersistentFlags().String("gitlab_token", "", "User's token for consuming the Gitlab API")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Disable standard data output")

	viper.BindPFlag("projects_path", rootCmd.PersistentFlags().Lookup("projects_path"))
	viper.BindPFlag("gitlab_url", rootCmd.PersistentFlags().Lookup("gitlab_url"))
	viper.BindPFlag("gitlab_token", rootCmd.PersistentFlags().Lookup("gitlab_token"))
	viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))
}

func initConfig() {
	var err error

	config := viper.GetString("config")
	if config == "" {
		home, err = homedir.Dir()
		typist.Must(err)

		viper.SetConfigName("config")
		viper.AddConfigPath(path.Join(home, ".p"))
		viper.AddConfigPath(path.Join(xdg.ConfigHome, "p"))
	} else {
		viper.SetConfigFile(config)
	}

	viper.AutomaticEnv()

	typist.Must(viper.ReadInConfig())
}
