#!/usr/bin/env bash
selections=$(curl -sL https://www.toptal.com/developers/gitignore/api/list\?format\=lines | fzf -m --height=80% \
    --prompt='▶ ' --pointer='→' \
    --border=sharp \
    --preview='curl -sL https://www.toptal.com/developers/gitignore/api/{}' \
    --preview-window='45%,border-sharp' \
    --prompt='gitignore ▶ ' \
    --bind='ctrl-r:reload(curl -sL https://www.toptal.com/developers/gitignore/api/list\?format\=lines)' \
    --bind='ctrl-p:toggle-preview' \
    --header '
--------------------------------------------------------------
* Tab/Shift-Tab:       mark multiple items
* ENTER:               append the selected to .gitignore file
* Ctrl-r:              refresh the list
* Ctrl-p:              toggle preview
* Ctrl-q:              exit
* Shift-up/Shift-down: scroll the preview
--------------------------------------------------------------
')

if [[ ${#selections[@]} == 0 ]]; then
    echo "▶ Nothing selected"
    return 0
fi

# allow multi-select
touch $PWD/.gitignore
while IFS= read -r s; do
    curl -sL https://www.toptal.com/developers/gitignore/api/$s >>$PWD/.gitignore
    echo "▶ Appended: $s"
done <<<"$selections"
