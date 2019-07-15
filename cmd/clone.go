package cmd

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	gittransport "gopkg.in/src-d/go-git.v4/plumbing/transport"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

var dotGit = regexp.MustCompile(`.git$`)

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: cloneRun,
}

func cloneRun(cmd *cobra.Command, args []string) {
	var u *url.URL
	var err error
	projectsPath := viper.GetString("projects_path")

	commonURL, err := gittransport.NewEndpoint(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	u, err = url.Parse(commonURL.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	clonePath := path.Join(projectsPath, u.Hostname(), u.Path)
	if dotGit.MatchString(u.Path) {
		clonePath = clonePath[:len(clonePath)-4]
	}

	sshKey, err := ioutil.ReadFile(path.Join(Home, ".ssh", "id_rsa"))
	signer, err := ssh.ParsePrivateKey([]byte(sshKey))

	_, err = git.PlainClone(clonePath, viper.GetBool("projects_path"), &git.CloneOptions{
		URL:      args[0],
		Auth:     &gitssh.PublicKeys{User: "git", Signer: signer},
		Progress: os.Stdout,
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(cloneCmd)

	cloneCmd.Flags().Bool("bare", false, "Make a bare Git repository.")
}
