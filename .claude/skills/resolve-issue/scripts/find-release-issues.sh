#!/usr/bin/env bash
# Find the nearest open Release issue and its open sub-issues.
# If no Release issue exists, outputs all open issues.
#
# Usage: bash find-release-issues.sh [owner/repo]
# Output: JSON object with { release: {...} | null, issues: [...] }

set -euo pipefail

REPO="${1:-typemd/typemd}"

# Fetch all open issues with sub-issues in one query
ALL_ISSUES=$(gh api graphql -f query='query($repo: String!, $owner: String!) {
  repository(owner: $owner, name: $repo) {
    issues(first: 50, states: OPEN, orderBy: {field: CREATED_AT, direction: ASC}) {
      nodes {
        number
        title
        labels(first: 5) { nodes { name } }
        subIssues(first: 30) {
          nodes {
            number
            title
            state
            labels(first: 5) { nodes { name } }
          }
        }
      }
    }
  }
}' -f owner="${REPO%%/*}" -f repo="${REPO##*/}")

# Try to find a Release issue (title matches "vX.Y.Z —")
RELEASE=$(echo "$ALL_ISSUES" | jq '
  .data.repository.issues.nodes
  | map(select(.title | test("^v[0-9]+\\.[0-9]+\\.[0-9]+ —")))
  | sort_by(.number)
  | first // null
')

if [ "$RELEASE" != "null" ]; then
  # Extract open sub-issues from the Release issue
  echo "$RELEASE" | jq '{
    release: { number: .number, title: .title },
    issues: [.subIssues.nodes[] | select(.state == "OPEN") | {
      number: .number,
      title: .title,
      labels: [.labels.nodes[].name]
    }]
  }'
else
  # No Release issue — return all open issues
  echo "$ALL_ISSUES" | jq '{
    release: null,
    issues: [.data.repository.issues.nodes[] | {
      number: .number,
      title: .title,
      labels: [.labels.nodes[].name]
    }]
  }'
fi
