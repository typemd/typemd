#!/usr/bin/env bash
# Get details for multiple issues in a single JSON array.
#
# Usage: bash get-issue-details.sh <issue_number> [issue_number ...]
# Output: JSON array of issue objects with number, title, labels, body (truncated to 300 chars)

set -euo pipefail

if [ $# -eq 0 ]; then
  echo "Usage: $0 <issue_number> [issue_number ...]" >&2
  exit 1
fi

echo "["
FIRST=true
for NUM in "$@"; do
  if [ "$FIRST" = true ]; then
    FIRST=false
  else
    echo ","
  fi
  gh issue view "$NUM" --json number,title,body,labels \
    | jq '{
        number: .number,
        title: .title,
        labels: [.labels[].name],
        body: (.body[:300])
      }'
done
echo "]"
