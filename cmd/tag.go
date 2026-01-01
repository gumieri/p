package cmd

import (
	"runtime/debug"
	"sort"
	"strconv"
	"strings"

	"github.com/gumieri/p/git"
	version "github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func getVersions(r *gogit.Repository) (versions []*version.Version, err error) {
	tags, err := r.Tags()
	if err != nil {
		return
	}

	versions = make([]*version.Version, 0)
	err = tags.ForEach(func(ref *plumbing.Reference) error {
		raw := ref.Name().Short()

		version, err := version.NewVersion(raw)
		if err != nil {
			return nil
		}

		versions = append(versions, version)
		return nil
	})

	if err != nil {
		return
	}

	sort.Sort(version.Collection(versions))

	return
}

// Segment enumerates a version segment
type Segment int

// Patch refers to the version segment of patch
const Patch Segment = 3

// Minor refers to the version segment of minor
const Minor Segment = 2

// Major refers to the version segment of major
const Major Segment = 1

func increaseVersion(old *version.Version, segment Segment) (*version.Version, error) {
	hasV := old.Original()[0] == 'v'

	segments := old.Segments()

	segments[segment-1]++

	for i := int(segment); i < len(segments); i++ {
		segments[i] = 0
	}

	sSegments := make([]string, len(segments))
	for i, iSegment := range segments {
		sSegments[i] = strconv.Itoa(iSegment)
	}

	var prefix string
	if hasV {
		prefix = "v"
	}

	return version.NewVersion(prefix + strings.Join(sSegments, "."))
}

func createTag(segment Segment) {
	debug.SetGCPercent(-1)

	topPath, err := git.TopLevelPath()
	t.Must(err)

	r, err := gogit.PlainOpen(topPath)
	t.Must(err)

	versions, err := getVersions(r)
	t.Must(err)

	lastVersion := versions[len(versions)-1]

	newVersion, err := increaseVersion(lastVersion, segment)
	t.Must(err)

	ref, err := r.Head()
	t.Must(err)

	_, err = r.CreateTag(newVersion.Original(), ref.Hash(), nil)
	t.Must(err)

	t.Outln(newVersion.Original())

	if push {
		t.Must(git.Push(remote, newVersion.Original(), false))
	}
}

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "The `git tag` command with power-ups.",
	Long:  `To tag git projects with semantic versioning.`,
	Run: func(cmd *cobra.Command, args []string) {
		debug.SetGCPercent(-1)

		topPath, err := git.TopLevelPath()
		t.Must(err)

		r, err := gogit.PlainOpen(topPath)
		t.Must(err)

		versions, err := getVersions(r)
		t.Must(err)

		for _, version := range versions {
			t.Outln(version.Original())
		}
	},
}

var tagPatchCmd = &cobra.Command{
	Use:   "patch",
	Short: "increment a patch to the version tag",
	Run: func(cmd *cobra.Command, args []string) {
		createTag(Patch)
	},
}

var tagMinorCmd = &cobra.Command{
	Use:   "minor",
	Short: "increment a minor to the version tag",
	Run: func(cmd *cobra.Command, args []string) {
		createTag(Minor)
	},
}

var tagMajorCmd = &cobra.Command{
	Use:   "major",
	Short: "increment a major to the version tag",
	Run: func(cmd *cobra.Command, args []string) {
		createTag(Major)
	},
}

var tagLastCmd = &cobra.Command{
	Use:   "last",
	Short: "get / manage the last version tag",
	Run: func(cmd *cobra.Command, args []string) {
		debug.SetGCPercent(-1)

		topPath, err := git.TopLevelPath()
		t.Must(err)

		r, err := gogit.PlainOpen(topPath)
		t.Must(err)

		versions, err := getVersions(r)
		t.Must(err)

		lastVersion := versions[len(versions)-1]

		t.Outln(lastVersion.Original())

		if del {
			t.Must(r.DeleteTag(lastVersion.Original()))
		}

		if push {
			t.Must(git.Push(remote, lastVersion.Original(), del))
		}
	},
}

var del bool
var push bool
var remote string

func init() {
	rootCmd.AddCommand(tagCmd)
	tagCmd.PersistentFlags().BoolVarP(&push, "push", "p", false, "push the changes to the remote")
	tagCmd.PersistentFlags().StringVar(&remote, "remote", "origin", "reference name to the git remote repository")

	tagCmd.AddCommand(tagPatchCmd)
	tagCmd.AddCommand(tagMinorCmd)
	tagCmd.AddCommand(tagMajorCmd)

	tagCmd.AddCommand(tagLastCmd)
	tagLastCmd.Flags().BoolVarP(&del, "delete", "d", false, "delete the tag\ncombined with -p / --push can also delete the remote tag")
}
