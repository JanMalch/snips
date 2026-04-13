//go:build !windows

package config

const fzfOptPreviewLabel = "[[ -n {} ]] && printf \" %s \" {1}"
const fzfOptListLabel = `if [[ -z $FZF_QUERY ]]; then
  echo " $FZF_MATCH_COUNT snippets "
else
  echo " $FZF_MATCH_COUNT snippets for '$FZF_QUERY' "
fi`
