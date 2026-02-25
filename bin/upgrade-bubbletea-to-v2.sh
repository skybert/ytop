#! /usr/bin/env bash

## author: torstein, torstein@skybert.net

set -o errexit
set -o nounset
set -o pipefail

go_mod() {
  go get charm.land/bubbles/v2/table@latest
  go get charm.land/bubbles/v2@latest
  go get charm.land/bubbletea/v2@latest
  go get charm.land/lipgloss/v2@latest
}

imports() {
  local _cwd=
  _cwd="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"


  local sed=sed
  test -x /opt/homebrew/bin/gsed && sed=gsed

  find "${_cwd}/.." -name "*.go" | while read -r f; do
    "${sed}" -i "s#github.com/charmbracelet/bubbletea#charm.land/bubbletea/v2#" "${f}"
    "${sed}" -i 's#github.com/charmbracelet/bubbles#charm.land/bubbles/v2#' "${f}"
    "${sed}" -i 's#tea.KeyMsg#tea.KeyPressMsg#' "${f}"
    "${sed}" -i 's#github.com/charmbracelet/lipgloss#charm.land/lipgloss/v2#' "${f}"
  done

  go mod tidy
}


main() {
  go_mod
  imports
}

main "$@"
