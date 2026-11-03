# Agentic Snap Test

An AI-driven end-to-end test for inference snaps. Instead of a fixed test script, an
LLM agent installs the snap the way a real user would, reads its documentation and CLI,
exercises its functionality, and writes a structured JSON report. When the report
contains defects, a second agent triages them against the open issues on the tracker
repo, and the pipeline can then file new GitHub issues and label cross-snap duplicates
automatically.

Everything runs inside a **workshop** — an isolated, disposable machine environment that
gives the agent a real `snapd` and passwordless `sudo`, so it can install and drive the
snap exactly as an end user would.

The agents are driven by [OpenCode](https://opencode.ai), pointed at a model served
through [OpenRouter](https://openrouter.ai).


## Calling it from a workflow

A Github action is provided in this repo at `.github/actions/agentic-test`.
The following example workflow invokes the action with the minimum required inputs.
Refer to the action's `action.yaml` for the full list of inputs and their defaults.

```yaml
name: Agentic Snap Test

on:
  workflow_dispatch:

jobs:
  ai-test:
    runs-on: ubuntu-latest
    steps:
      - name: Run test
        uses: canonical/inference-snaps-dev/.github/actions/agentic-test
        with:
          snap-name: smollm2
          snap-channel: edge
          create-issues: false
          openrouter-api-key: ${{ secrets.OPENROUTER_API_KEY }}
          issue-create-token: ${{ secrets.ISSUE_CREATE_TOKEN }}
```

> Issues are filed in a *different* repo than the workflow runs in, so a fine-grained PAT
> with Issues read/write is required — the default token can't write cross-repo.

## Running it locally

`run.sh` reproduces the full flow: launch workshop → test → triage → issues → verdict gate → cleanup.

Prerequisites:

- The `workshop` snap installed and available on `PATH`.
- `opencode`, `jq`, and (for create mode) `gh` installed.
- An OpenRouter API key.

Store your secrets in local files (e.g. `openrouter.key`, `issue-create.token`) and pass
them inline so they are never exported to the session or written to shell history.
Run from this directory (`tests/agentic/`):

```bash
# Dry run: test, triage, and print any new/duplicate findings (no writes to GitHub).
OPENROUTER_API_KEY="$(cat openrouter.key)" \
  SNAP_NAME=smollm2 \
  SNAP_CHANNEL=latest/edge \
  ./run.sh
```

To also file issues and label cross-snap duplicates, provide a token and enable create
mode:

```bash
OPENROUTER_API_KEY="$(cat openrouter.key)" \
  ISSUE_CREATE_TOKEN="$(cat issue-create.token)" \
  ISSUE_REPO=canonical/inference-snaps \
  CREATE_ISSUES=true \
  SNAP_NAME=smollm2 \
  SNAP_CHANNEL=latest/edge \
  ./run.sh
```

Of the variables above, `OPENROUTER_API_KEY` is required and `SNAP_NAME` selects the snap
under test; everything else has a sensible default. See the top of `run.sh` for the full
list of variables and their defaults.

You can also run the phases individually — e.g. `./scripts/run-test-agent.sh` then
`./scripts/run-triage-agent.sh` — which is handy when iterating on `AGENT.md` or `TRIAGE.md`. The
agents write their reports under `/tmp` inside the workshop, which the scripts then copy to
the working directory:

- `snap-test-report.json` — the testing agent's verdict, environment and findings.
- `snap-triage-report.json` — the triage agent's NEW/DUPLICATE classification.
