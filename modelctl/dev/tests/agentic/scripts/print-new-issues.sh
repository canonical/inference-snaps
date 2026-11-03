#!/usr/bin/env bash
set -euo pipefail

# ---------------------------------------------------------------------------
# Prints each NEW triaged finding (title + suggested body) from
# snap-triage-report.json so it can be copied into a GitHub issue by hand.
#
# Intended for dry-run mode (create_issues=false); when issues are filed
# automatically, create-new-issues.sh runs instead.
# ---------------------------------------------------------------------------


if [[ ! -f snap-triage-report.json ]]; then
    echo "ERROR: snap-triage-report.json not found; run-triage-agent.sh first." >&2
    exit 1
fi

NEW_INDICES=$(jq -r '[.results | to_entries[] | select(.value.status == "NEW") | .key] | .[]' snap-triage-report.json)

if [[ -z "$NEW_INDICES" ]]; then
    echo "No NEW findings — nothing to file in $ISSUE_REPO."
    exit 0
fi

echo "NEW finding(s) were not filed automatically (create_issues=false)."
echo "Copy each into a new issue in $ISSUE_REPO, or re-run with create_issues=true."
for i in $NEW_INDICES; do
    ISSUE_TITLE=$(jq -r ".results[$i].suggested_title" snap-triage-report.json)
    echo ""
    echo "Title: $ISSUE_TITLE"
    echo ""
    echo "Body:"
    echo "------------------------------------------------------------------------"
    jq -r ".results[$i].suggested_body" snap-triage-report.json
    echo "------------------------------------------------------------------------"
done

