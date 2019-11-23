package cmd

import (
	"io/ioutil"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	gittransport "gopkg.in/src-d/go-git.v4/plumbing/transport"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

var dotGit = regexp.MustCompile(`.git$`)

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "The `git clone` command with power-ups.",
	Long: `To clone git projects into projects path.
It will create the path as $PROJECTS_PATH/HOST/…/GROUPS/PROJECT.`,
	Run: cloneRun,
}

func cloneRun(cmd *cobra.Command, args []string) {
	var u *url.URL
	var err error

	commonURL, err := gittransport.NewEndpoint(args[0])
	t.Must(err)

	u, err = url.Parse(commonURL.String())
	t.Must(err)

	var addresses []*url.URL
	if u.Hostname() == viper.GetString("gitlab_url") {
		protocol := "http"
		if viper.GetBool("gitlab_https") {
			protocol = "https"
		}

		gl := gitlab.NewClient(nil, viper.GetString("gitlab_token"))
		gl.SetBaseURL(protocol + "://" + viper.GetString("gitlab_url") + "/api/v4")
		allProjects, _, err := gl.Projects.ListProjects(&gitlab.ListProjectsOptions{Archived: gitlab.Bool(false)})
		t.Must(err)

		for _, project := range allProjects {
			subprojectPath := path.Join("/", project.PathWithNamespace)

			rel, err := filepath.Rel(u.Path, subprojectPath)
			t.Must(err)

			if len(rel) > 2 && rel[0:2] == ".." { // is not a subdirectory (subgroup)
				continue
			}

			address := new(url.URL)
			address.Scheme = u.Scheme
			address.Opaque = u.Opaque
			address.User = u.User
			address.Host = u.Host
			address.Path = subprojectPath
			address.RawPath = u.RawPath
			address.ForceQuery = u.ForceQuery
			address.RawQuery = u.RawQuery
			address.Fragment = u.Fragment

			addresses = append(addresses, address)
		}
	}

	if len(addresses) == 0 {
		addresses = append(addresses, u)
	}

	var wg sync.WaitGroup
	for _, project := range addresses {
		wg.Add(1)
		go cloneProject(project, &wg)
	}
	wg.Wait()
}

func cloneProject(u *url.URL, wg *sync.WaitGroup) {
	clonePath := path.Join(viper.GetString("projects_path"), u.Hostname(), u.Path)
	if dotGit.MatchString(u.Path) {
		clonePath = clonePath[:len(clonePath)-4]
	}

	t.Infof("%s: cloning into %s.\n", u, clonePath)

	sshKey, err := ioutil.ReadFile(path.Join(home, ".ssh", "id_rsa"))
	signer, err := ssh.ParsePrivateKey([]byte(sshKey))
	_, err = git.PlainClone(clonePath, viper.GetBool("bare"), &git.CloneOptions{
		URL:  u.String(),
		Auth: &gitssh.PublicKeys{User: u.User.Username(), Signer: signer},
	})

	if err == nil {
		t.Infof("%s: completed.\n", u)
	} else {
		t.Errorf("%s: %s\n", u, err)
	}

	wg.Done()
}

func init() {
	rootCmd.AddCommand(cloneCmd)

	cloneCmd.Flags().Bool("bare", false, "Make a bare Git repository.")
}
