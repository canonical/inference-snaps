#!/usr/bin/env bash
set -euo pipefail

# ---------------------------------------------------------------------------
# Records this snap on cross-snap duplicates from snap-triage-report.json: if a
# finding duplicates an existing issue first detected on a *different* snap (its
# labels include some snap/<other> label but not snap/$SNAP_NAME), add
# snap/$SNAP_NAME to that issue so the tracker lists every snap affected by the
# same bug.
#
# The workflow only runs this in create mode. Issues live in a *different* repo
# than the one running this workflow, so a fine-grained PAT scoped to
# $ISSUE_REPO (Issues: read & write) is required — the default GITHUB_TOKEN
# cannot write cross-repo.
# ---------------------------------------------------------------------------


if [[ ! -f snap-triage-report.json ]]; then
    echo "ERROR: snap-triage-report.json not found; run-triage-agent.sh first." >&2
    exit 1
fi

CROSS_SNAP_INDICES=$(jq -r --arg snap_label "snap/$SNAP_NAME" '
    [.results | to_entries[]
     | select(.value.status == "DUPLICATE"
              and .value.duplicate_of_number != null
              and ((.value.duplicate_of_labels // []) | any(startswith("snap/")))
              and ((.value.duplicate_of_labels // []) | index($snap_label) | not))
     | .key] | .[]' snap-triage-report.json)

if [[ -z "$CROSS_SNAP_INDICES" ]]; then
    echo "No cross-snap duplicate issues to label."
    exit 0
fi

if [[ -z "$ISSUE_CREATE_TOKEN" ]]; then
    echo "WARNING: ISSUE_CREATE_TOKEN is empty; cannot label cross-snap duplicate issues in $ISSUE_REPO." >&2
    exit 0
fi

echo "Labelling cross-snap duplicate issues in $ISSUE_REPO"
LABELLED=0
LABEL_FAILED=0
for i in $CROSS_SNAP_INDICES; do
    ISSUE_NUMBER=$(jq -r ".results[$i].duplicate_of_number" snap-triage-report.json)
    ISSUE_URL=$(jq -r ".results[$i].duplicate_of_url // \"(unknown)\"" snap-triage-report.json)

    if GH_TOKEN="$ISSUE_CREATE_TOKEN" gh issue edit "$ISSUE_NUMBER" \
        --repo "$ISSUE_REPO" --add-label "snap/$SNAP_NAME" >/dev/null 2>/tmp/gh-label-err; then
        echo "Labelled: $ISSUE_URL  (added snap/$SNAP_NAME)"
        LABELLED=$((LABELLED + 1))
    else
        ERR_MSG=$(cat /tmp/gh-label-err 2>/dev/null || echo "unknown error")
        echo "FAILED: $ISSUE_URL — $ERR_MSG" >&2
        LABEL_FAILED=$((LABEL_FAILED + 1))
    fi
done

echo ""
echo "Issues labelled : $LABELLED"
echo "Labels failed   : $LABEL_FAILED"

if [[ "$LABEL_FAILED" -gt 0 ]]; then
    echo "WARNING: $LABEL_FAILED cross-snap duplicate(s) could not be labelled in $ISSUE_REPO." >&2
fi

