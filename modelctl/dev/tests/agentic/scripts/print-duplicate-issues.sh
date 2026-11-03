#!/usr/bin/env bash
set -euo pipefail

# ---------------------------------------------------------------------------
# Prints a summary of the DUPLICATE triaged findings from
# snap-triage-report.json — each finding and the existing issue it duplicates.
# Read-only: labelling of cross-snap duplicates is done by
# label-duplicate-issues.sh.
# ---------------------------------------------------------------------------

if [[ ! -f snap-triage-report.json ]]; then
    echo "ERROR: snap-triage-report.json not found; run-triage-agent.sh first." >&2
    exit 1
fi

DUP_INDICES=$(jq -r '[.results | to_entries[] | select(.value.status == "DUPLICATE") | .key] | .[]' snap-triage-report.json)

if [[ -z "$DUP_INDICES" ]]; then
    echo "No DUPLICATE findings."
    exit 0
fi

echo "Duplicate finding(s):"
for i in $DUP_INDICES; do
    FINDING_TITLE=$(jq -r ".results[$i].finding_title // .results[$i].suggested_title" snap-triage-report.json)
    DUP_URL=$(jq -r ".results[$i].duplicate_of_url // \"(unknown)\"" snap-triage-report.json)
    DUP_REASON=$(jq -r ".results[$i].duplicate_reason // \"(unknown)\"" snap-triage-report.json)
    echo ""
    echo "$FINDING_TITLE"
    echo "  Duplicate of: $DUP_URL"
    echo "  Reason      : $DUP_REASON"
done


