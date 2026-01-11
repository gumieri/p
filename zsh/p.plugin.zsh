# p.plugin.zsh - Zsh integration for p (projects)
# https://github.com/gumieri/p
#
# This plugin provides a Ctrl+P keybinding to quickly navigate to projects
# using fzf for fuzzy filtering.
#
# Requirements:
#   - p: The projects management CLI (this project)
#   - fzf: Command-line fuzzy finder (https://github.com/junegunn/fzf)
#
# Configuration:
#   - PROJECTS_PATH: Base directory for projects (default: $HOME/Projects)

if (( ! $+commands[p] )); then
    return
fi

function open-project-path {
    if (( ! $+commands[fzf] )); then
        echo -e "\n\033[31mError: fzf is required but not installed.\033[0m"
        zle reset-prompt
        return 1
    fi

    local project=$(p | fzf --height 40% --layout=reverse --border --prompt="Project > " --scheme=path)

    if [[ -n "$project" ]]; then
        local base_path="${PROJECTS_PATH:-$HOME/Projects}"

        if [[ -d "$base_path/$project" ]]; then
            cd "$base_path/$project"
            echo -e "\n\033[32m-> $project\033[0m"
        else
            echo -e "\n\033[31mError: Directory not found: $base_path/$project\033[0m"
        fi
    fi

    zle reset-prompt
}

zle -N open-project-path
bindkey '^p' open-project-path
