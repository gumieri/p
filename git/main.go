package git

import (
	"os/exec"
	"strings"
)

// TopLevelPath return the directory where the `.git` is located
func TopLevelPath() (path string, err error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return
	}

	path = strings.Trim(string(out), "\n")

	return
}

// Push a branch or tag to the specified remote
func Push(remote, reference string, del bool) (err error) {
	parameters := []string{"push", remote, reference}

	if del {
		parameters = append(parameters, "--delete")
	}

	_, err = exec.Command("git", parameters...).Output()
	return
}
