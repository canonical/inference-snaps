#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# ---------------------------------------------------------------------------
# Phase 1: run the testing agent and produce snap-test-report.json.
#
# This script launches the workshop (if not already running) and runs the
# OpenCode testing agent inside it. It writes the agent's JSON report to
# snap-test-report.json in the current directory and exposes the error count
# and verdict as GitHub step outputs so later workflow steps (triage, issues,
# verdict gate) can decide whether to run.
#
# It intentionally does NOT remove the workshop or fail the job on a FAIL
# verdict — that is left to later steps so triage/issue-filing can run first.
# ---------------------------------------------------------------------------


workshop exec --env OPENROUTER_API_KEY="$OPENROUTER_API_KEY" --env OPENROUTER_MODEL="$OPENROUTER_MODEL" --env SNAP_NAME="$SNAP_NAME" --env SNAP_CHANNEL="$SNAP_CHANNEL" --env GITHUB_SERVER_URL="${GITHUB_SERVER_URL:-}" --env GITHUB_REPOSITORY="${GITHUB_REPOSITORY:-}" --env GITHUB_RUN_ID="${GITHUB_RUN_ID:-}" -- \
    opencode run --auto --log-level ERROR \
    --model "openrouter/$OPENROUTER_MODEL" \
    "$(cat "$SCRIPT_DIR/AGENT.md")" || true

# Pull the JSON report the agent wrote; fail if it is missing or invalid.
REPORT_JSON=$(workshop exec -- sh -c 'cat /tmp/snap-test-report.json 2>/dev/null' || true)

if [[ -z "$REPORT_JSON" ]]; then
    echo "ERROR: agent did not write /tmp/snap-test-report.json" >&2
    exit 1
fi

if ! echo "$REPORT_JSON" | jq . > /dev/null 2>&1; then
    echo "ERROR: /tmp/snap-test-report.json is not valid JSON:" >&2
    echo "$REPORT_JSON" >&2
    exit 1
fi

echo "$REPORT_JSON" > snap-test-report.json
jq . snap-test-report.json

VERDICT=$(jq -r '.verdict // "UNKNOWN"' snap-test-report.json)
SUMMARY=$(jq -r '.summary // "(no summary)"' snap-test-report.json)
ERROR_COUNT=$(jq '[.findings[] | select(.severity == "error")] | length' snap-test-report.json)

echo ""
echo "Verdict : $VERDICT"
echo "Summary : $SUMMARY"
echo "Errors  : $ERROR_COUNT"

# Expose the verdict and error count as GitHub step outputs so later steps can
# gate on them (e.g. only triage when there are error findings).
if [[ -n "${GITHUB_OUTPUT:-}" ]]; then
    {
        echo "verdict=$VERDICT"
        echo "error_count=$ERROR_COUNT"
    } >> "$GITHUB_OUTPUT"
fi

