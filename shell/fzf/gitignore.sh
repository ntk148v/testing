#!/usr/bin/env bash
selections=$(curl -sL https://www.toptal.com/developers/gitignore/api/list\?format\=lines | fzf --height=80% \
    --prompt='▶ ' --pointer='→' \
    --border=sharp \
    --preview='curl -sL https://www.toptal.com/developers/gitignore/api/{}' \
    --preview-window='45%,border-sharp' \
    --prompt='gitignore ▶ ' \
    --bind='ctrl-r:reload(curl -sL https://www.toptal.com/developers/gitignore/api/list\?format\=lines)' \
    --bind='ctrl-p:toggle-preview' \
    --header '
CTRL-A to select all
CTRL-x to deselect all
ENTER to append the selected to .gitignore file
CTRL-r to refresh the list
CTRL-P to toggle preview
')

# allow multi-select
for s in "${selections[@]}"; do
    echo "▶ Selected: $s"
    curl -sL https://www.toptal.com/developers/gitignore/api/$s >>/tmp/.gitignore
done
