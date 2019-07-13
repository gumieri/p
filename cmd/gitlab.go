package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

func gitlabRun(cmd *cobra.Command, args []string) {
	gl := gitlab.NewClient(nil, viper.GetString("gitlab_token"))
	gl.SetBaseURL("http://" + viper.GetString("gitlab_url") + "/api/v4")
	projects, _, err := gl.Projects.ListProjects(&gitlab.ListProjectsOptions{Archived: gitlab.Bool(false)})
	if err != nil {
		log.Fatal(err)
	}
	for _, project := range projects {
		fmt.Println(project.PathWithNamespace)
	}
}

// gitlabCmd represents the gitlab command
var gitlabCmd = &cobra.Command{
	Use:   "gitlab",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: gitlabRun,
}

func init() {
	rootCmd.AddCommand(gitlabCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gitlabCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gitlabCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
