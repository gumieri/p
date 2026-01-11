# p(rojects)
Collection of helping commands for the management of projects using git.

## Configurations
`p` will look for configurations at `$XDG_CONFIG_HOME/p/config.yaml` or `~/.p.yaml`.
The file can be a JSON, TOML, YAML, HCL or envfile. Any configuration can also be set as environment variable.
I recommend to use the projects path configuration as environment variable (`$PROJECTS_PATH`).

A configuration example:
```
projects_path: "~/Projects"
gitlab_token: "wd4qRsNpHci5wtizByap"
gitlab_url: "git.private.local"
gitlab_https: false
```

## Commands
### `p`
List all projects into `$PROJECTS_PATH`

### `p clone`
A `git clone` with power-ups.
It will create the path as `$PROJECTS_PATH/HOST/…/GROUPS/PROJECT`:
```
$ p clone git@github.com:gumieri/p.git
ssh://git@github.com/gumieri/p.git: cloning into /home/rafael@gumieri/Projects/github.com/gumieri/p.
ssh://git@github.com/gumieri/p.git: completed.
```
It is also able to clone multiple projects from gitlab groups and subgroups:
```
$ p clone git@git.private.local:group/subgroup
ssh://git@git.private.local:group/subgroup/project1.git: cloning into /home/rafael@gumieri/Projects/git.private.local/group/subgroup/project1.
ssh://git@git.private.local:group/subgroup/project2.git: cloning into /home/rafael@gumieri/Projects/git.private.local/group/subgroup/project2.
ssh://git@git.private.local:group/subgroup/project1.git: completed.
ssh://git@git.private.local:group/subgroup/project3.git: cloning into /home/rafael@gumieri/Projects/git.private.local/group/subgroup/project3.
ssh://git@git.private.local:group/subgroup/project2.git: completed.
ssh://git@git.private.local:group/subgroup/project3.git: completed.
```

### `p tag`
To tag git projects with semantic versioning.
It will respect if it has a `v` prefix or not.

Listing all tags (semantic ordered):
```
$ p tag
0.0.0
```
Creating new tags, and pushing it:
```
$ p tag patch
0.0.1
$ p tag minor
0.1.1
$ p tag major
1.1.1
$ p tag patch --push
1.1.2
$ p tag major --remote elsewhere --push
2.1.2
```
It also has a shortcut to handle the higher (semantic ordered) tag:
```
$ p tag last
1.1.1
$ p tag last --delete
$ p tag last
0.1.1
```

A simple combination:
```
$ p tag patch && git push origin `p tag last`
```

## Zsh Integration

For quick project navigation with `Ctrl+P`, you can load the included Zsh plugin.

### Requirements
- [fzf](https://github.com/junegunn/fzf) - Command-line fuzzy finder

### Using Sheldon
Add to your `~/.config/sheldon/plugins.toml`:
```toml
[plugins.p]
github = "gumieri/p"
dir = "zsh"
```

### Manual Installation
Source the plugin in your `~/.zshrc`:
```zsh
source /path/to/p/zsh/p.plugin.zsh
```

### Usage
Press `Ctrl+P` to open the project selector. Use fzf to filter and select a project, then press Enter to navigate to it.
