#!/usr/bin/env bash
set -euo pipefail

# Export variables for subscripts.
export SNAP_NAME="${SNAP_NAME:-smollm2}"
export SNAP_CHANNEL="${SNAP_CHANNEL:-edge}"
export OPENROUTER_API_KEY="${OPENROUTER_API_KEY:-}"
export OPENROUTER_MODEL="${OPENROUTER_MODEL:-deepseek/deepseek-v4-flash}"
export ISSUE_REPO="${ISSUE_REPO:-canonical/inference-snaps}"
export ISSUE_CREATE_TOKEN="${ISSUE_CREATE_TOKEN:-}"

if [[ -z "${OPENROUTER_API_KEY:-}" ]]; then
    echo "ERROR: OPENROUTER_API_KEY is not set in this shell." >&2
    echo "Export it first, e.g.: export OPENROUTER_API_KEY=sk-or-..." >&2
    exit 1
fi

# Launch workshop. In CI it is already running.
if [[ "${GITHUB_ACTIONS:-}" != "true" ]]; then
    echo "Launching workshop instance"
    workshop launch
fi
trap 'echo "::group::Remove workshop instance"; workshop remove || true; echo "::endgroup::"' EXIT

# Run the test agent.
echo "::group::Test agent"
./scripts/run-test-agent.sh
echo "::endgroup::"

VERDICT=$(jq -r '.verdict // "UNKNOWN"' snap-test-report.json)
ERROR_COUNT=$(jq '[.findings[] | select(.severity == "error")] | length' snap-test-report.json)

# If anything was found, triage the results and create issues if requested.
if [[ "$ERROR_COUNT" -gt 0 ]]; then
    echo "::group::Triage agent"
    ./scripts/run-triage-agent.sh
    echo "::endgroup::"

    echo "::group::Duplicate issues"
    ./scripts/print-duplicate-issues.sh
    echo "::endgroup::"

    if [[ "${CREATE_ISSUES:-false}" == "true" ]]; then
        echo "::group::Create issues"
        ./scripts/create-new-issues.sh
        ./scripts/label-duplicate-issues.sh
        echo "::endgroup::"
    else
        echo "::group::New issues"
        ./scripts/print-new-issues.sh
        echo "::endgroup::"
    fi
fi

# Fail the workflow if the verdict is not PASS.
if [[ "$VERDICT" != "PASS" ]]; then
    echo "Test verdict: $VERDICT" >&2
    exit 1
fi
