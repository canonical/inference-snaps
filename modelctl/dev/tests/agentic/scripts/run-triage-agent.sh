#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# ---------------------------------------------------------------------------
# Phase 2: run the triage agent and produce snap-triage-report.json.
#
# This step compares the error findings in snap-test-report.json against the
# open GitHub issues and classifies each as NEW or DUPLICATE. It reuses the
# workshop already launched by run-test-agent.sh and reads the test report the agent
# wrote to /tmp inside the workshop.
#
# It only makes sense to run this when the test report contains error
# findings; the workflow gates it on run-test-agent.sh's error_count output.
# ---------------------------------------------------------------------------


if [[ ! -f snap-test-report.json ]]; then
    echo "ERROR: snap-test-report.json not found; run-test-agent.sh first." >&2
    exit 1
fi

workshop exec --env OPENROUTER_API_KEY="$OPENROUTER_API_KEY" --env OPENROUTER_MODEL="$OPENROUTER_MODEL" --env ISSUE_REPO="$ISSUE_REPO" --env GITHUB_SERVER_URL="${GITHUB_SERVER_URL:-}" --env GITHUB_REPOSITORY="${GITHUB_REPOSITORY:-}" --env GITHUB_RUN_ID="${GITHUB_RUN_ID:-}" -- \
    opencode run --auto --log-level ERROR \
    --model "openrouter/$OPENROUTER_MODEL" \
    "$(cat "$SCRIPT_DIR/TRIAGE.md")" || true

TRIAGE_JSON=$(workshop exec -- sh -c 'cat /tmp/snap-triage-report.json 2>/dev/null' || true)
if [[ -z "$TRIAGE_JSON" ]] || ! echo "$TRIAGE_JSON" | jq . > /dev/null 2>&1; then
    echo "ERROR: triage agent did not produce a valid /tmp/snap-triage-report.json" >&2
    exit 1
fi

echo "$TRIAGE_JSON" > snap-triage-report.json
jq . snap-triage-report.json

NEW_COUNT=$(jq '.new_count // 0' snap-triage-report.json)
DUP_COUNT=$(jq '.duplicate_count // 0' snap-triage-report.json)
echo ""
echo "New issues : $NEW_COUNT"
echo "Duplicates : $DUP_COUNT"

