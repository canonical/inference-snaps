#!/usr/bin/env bash
set -euo pipefail

# ---------------------------------------------------------------------------
# Creates GitHub issues for NEW triaged findings in $ISSUE_REPO and logs a
# summary only (no full bodies). Intended for create mode (create_issues=true);
# in dry-run mode print-new-issues.sh runs instead.
#
# Issues are filed in a *different* repo than the one running this workflow, so
# a fine-grained PAT scoped to $ISSUE_REPO (Issues: read & write) is required —
# the default GITHUB_TOKEN cannot write cross-repo.
# ---------------------------------------------------------------------------


if [[ ! -f snap-triage-report.json ]]; then
    echo "ERROR: snap-triage-report.json not found; run-triage-agent.sh first." >&2
    exit 1
fi

NEW_INDICES=$(jq -r '[.results | to_entries[] | select(.value.status == "NEW") | .key] | .[]' snap-triage-report.json)

if [[ -z "$NEW_INDICES" ]]; then
    echo "No NEW findings to file — nothing to create in $ISSUE_REPO."
    exit 0
fi

if [[ -z "$ISSUE_CREATE_TOKEN" ]]; then
    echo "WARNING: ISSUE_CREATE_TOKEN is empty; cannot create issues in $ISSUE_REPO." >&2
    echo "Set the ISSUE_CREATE_TOKEN secret (fine-grained PAT with Issues: read & write on $ISSUE_REPO)." >&2
    exit 0
fi

echo "Creating GitHub issues in $ISSUE_REPO"
CREATED=0
FAILED=0
for i in $NEW_INDICES; do
    ISSUE_TITLE=$(jq -r ".results[$i].suggested_title" snap-triage-report.json)
    ISSUE_BODY=$(jq -r ".results[$i].suggested_body" snap-triage-report.json)

    # Guard against missing fields: jq -r prints the literal "null" for an
    # absent key, which would otherwise file a malformed issue.
    if [[ -z "$ISSUE_TITLE" || "$ISSUE_TITLE" == "null" ]]; then
        echo "SKIPPED: result[$i] has no suggested_title; not creating an issue." >&2
        FAILED=$((FAILED + 1))
        continue
    fi
    if [[ -z "$ISSUE_BODY" || "$ISSUE_BODY" == "null" ]]; then
        echo "SKIPPED: '$ISSUE_TITLE' has no suggested_body; not creating an issue." >&2
        FAILED=$((FAILED + 1))
        continue
    fi

    if ISSUE_URL=$(GH_TOKEN="$ISSUE_CREATE_TOKEN" gh issue create \
        --repo "$ISSUE_REPO" \
        --title "$ISSUE_TITLE" \
        --body "$ISSUE_BODY" \
        --label "bot" \
        --label "snap/$SNAP_NAME" 2>/tmp/gh-issue-err); then
        echo "Created: $ISSUE_URL  ($ISSUE_TITLE)"
        CREATED=$((CREATED + 1))

        # gh issue create has no --type flag, so patch it via the API.
        # Non-fatal: warn but keep the issue.
        ISSUE_NUMBER=$(basename "$ISSUE_URL")
        if ! GH_TOKEN="$ISSUE_CREATE_TOKEN" gh api --method PATCH \
            "repos/$ISSUE_REPO/issues/$ISSUE_NUMBER" -f type=Bug \
            >/dev/null 2>/tmp/gh-type-err; then
            ERR_MSG=$(jq -r '.message // empty' /tmp/gh-type-err 2>/dev/null)
            ERR_MSG=${ERR_MSG:-$(cat /tmp/gh-type-err 2>/dev/null || echo "unknown error")}
            echo "WARNING: created $ISSUE_URL but failed to set type=Bug — $ERR_MSG" >&2
        fi
    else
        ERR_MSG=$(cat /tmp/gh-issue-err 2>/dev/null || echo "unknown error")
        echo "FAILED: $ISSUE_TITLE — $ERR_MSG" >&2
        FAILED=$((FAILED + 1))
    fi
done

echo ""
echo "Issues created : $CREATED"
echo "Issues failed  : $FAILED"

if [[ "$FAILED" -gt 0 ]]; then
    echo "WARNING: $FAILED issue(s) could not be created in $ISSUE_REPO." >&2
fi

