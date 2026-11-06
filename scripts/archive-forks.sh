#!/usr/bin/env bash
# Hide public forks and old 2018 sample repos from the public profile.
# Requires: gh auth login
set -euo pipefail

gh auth status

echo "Making forks private..."
gh repo list mperuch --fork --limit 200 --json name -q '.[].name' | while read -r name; do
  echo "  private mperuch/$name"
  gh repo edit "mperuch/$name" --visibility private --accept-visibility-change-consequences
done

for name in mysnippets simple-redux-test create-react-app-material-ui; do
  echo "  private old sample mperuch/$name"
  gh repo edit "mperuch/$name" --visibility private --accept-visibility-change-consequences || true
done

echo "Done. Pin quote-engine on github.com/mperuch."
