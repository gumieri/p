package cmd

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/adrg/xdg"
	"github.com/gumieri/typist"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var t = typist.New(&typist.Config{})

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
	projects := ""

	cache := ""
	useCache := !viper.GetBool("no-cache")
	cachePath := path.Join(xdg.CacheHome, "p")
	cacheFilePath := path.Join(cachePath, "projects")

	if useCache {
		cacheB, err := ioutil.ReadFile(cacheFilePath)
		t.Must(err)
		cache = string(cacheB)
		t.Outln(cache)
	}

	t.Must(filepath.Walk(projectsPath, func(cwd string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if _, err := os.Stat(path.Join(cwd, ".git")); !os.IsNotExist(err) {
				project := cwd[len(projectsPath)+1:]
				projects += "\n" + project
				if !useCache || !strings.Contains(cache, project) {
					t.Outln(project)
				}

				return filepath.SkipDir
			}
		}
		return nil
	}))

	os.MkdirAll(cachePath, 0700)
	ioutil.WriteFile(cacheFilePath, []byte(projects), 0600)
}

func persistentPreRun(cmd *cobra.Command, args []string) {
	t.Config.Quiet = quiet
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
	t.Must(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().Bool("no-cache", false, "perform the search with no cache data")
	viper.BindPFlag("no-cache", rootCmd.Flags().Lookup("no-cache"))

	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.p.yaml or $XDG_CONFIG_HOME/p/config.yaml)")

	rootCmd.PersistentFlags().String("projects_path", "", "The root path of all projects")
	viper.BindPFlag("projects_path", rootCmd.PersistentFlags().Lookup("projects_path"))

	rootCmd.PersistentFlags().String("gitlab_url", "", "Base URL for the Gitlab API")
	viper.BindPFlag("gitlab_url", rootCmd.PersistentFlags().Lookup("gitlab_url"))

	rootCmd.PersistentFlags().String("gitlab_token", "", "User's token for consuming the Gitlab API")
	viper.BindPFlag("gitlab_token", rootCmd.PersistentFlags().Lookup("gitlab_token"))

	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Disable standard data output")
	viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))
}

func initConfig() {
	var err error

	config := viper.GetString("config")
	if config == "" {
		home, err = homedir.Dir()
		t.Must(err)

		viper.SetConfigName("config")
		viper.AddConfigPath(path.Join(home, ".p"))
		viper.AddConfigPath(path.Join(xdg.ConfigHome, "p"))
	} else {
		viper.SetConfigFile(config)
	}

	viper.AutomaticEnv()

	viper.ReadInConfig()
}
