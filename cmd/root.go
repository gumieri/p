package cmd

import (
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"

	"github.com/adrg/xdg"
	"github.com/gumieri/typist"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var t = typist.New(&typist.Config{})

// disableGC disables garbage collection for short-lived commands
// This improves performance for commands that exit quickly
func disableGC() {
	debug.SetGCPercent(-1)
}

// GetProjects already in the projects path
func GetProjects(projectsPath string) (projects []string, err error) {
	projects = make([]string, 0, 100)
	err = filepath.Walk(projectsPath, func(cwd string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if _, err := os.Stat(path.Join(cwd, ".git")); !os.IsNotExist(err) {
				rel, err := filepath.Rel(projectsPath, cwd)
				if err != nil {
					return err
				}
				projects = append(projects, rel)
				return filepath.SkipDir
			}
		}
		return nil
	})

	return
}

func fileOrDirExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func rootRun(cmd *cobra.Command, args []string) {
	disableGC()

	projectsPath, err := homedir.Expand(viper.GetString("projects_path"))
	t.Must(err)

	var projectsBuilder strings.Builder
	var cacheMap map[string]bool
	var existingCache string
	useCache := !viper.GetBool("no-cache")
	cachePath := path.Join(xdg.CacheHome, "p")
	cacheFilePath := path.Join(cachePath, "projects")

	if useCache {
		if !fileOrDirExists(cachePath) {
			t.Must(os.MkdirAll(cachePath, 0700))
		}

		cacheMap = make(map[string]bool)
		if fileOrDirExists(cacheFilePath) {
			cacheB, err := os.ReadFile(cacheFilePath)
			t.Must(err)
			existingCache = string(cacheB)
			t.Outln(existingCache)
			for _, line := range strings.Split(existingCache, "\n") {
				line = strings.TrimSpace(line)
				if line != "" {
					cacheMap[line] = true
				}
			}
		}
	}

	t.Must(filepath.Walk(projectsPath, func(cwd string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if fileOrDirExists(path.Join(cwd, ".git")) {
				project, err := filepath.Rel(projectsPath, cwd)
				if err != nil {
					return err
				}
				projectsBuilder.WriteString(project)
				//projectsBuilder.WriteString("\n")
				if !useCache || !cacheMap[project] {
					t.Outln(project)
				}

				return filepath.SkipDir
			}
		}
		return nil
	}))

	if useCache && projectsBuilder.Len() > 0 {
		projects := projectsBuilder.String()
		if !fileOrDirExists(cacheFilePath) || existingCache != projects {
			t.Must(os.WriteFile(cacheFilePath, []byte(projects), 0600))
		}
	}
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

	rootCmd.PersistentFlags().String("projects_path", "~/Projects", "The root path of all projects")
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
