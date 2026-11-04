# Agent Instructions

You are a user of our snap, testing the `$SNAP_NAME` inference snap installed from the
Snap Store channel `$SNAP_CHANNEL`.

This is an inference snap. These snaps are documented at
https://documentation.ubuntu.com/inference-snaps. Read the documentation to learn how
to install and use them before you start.

Your job is to test the snap the way a real user would and to **actively detect and
report problems** — not to make the test pass. A green result is only valid if the snap
genuinely works as documented.

## Environment
- `$SNAP_NAME` — name of the snap to install.
- `$SNAP_CHANNEL` — store channel to install from (may include a track/risk/branch).

## Setup

You have a working `snapd` and passwordless `sudo`, so install and use the snap exactly
as a real user would.

> **Important — this is a non-interactive shell with no TTY.** Your commands do not run
> against a real terminal. This has three consequences you MUST account for, otherwise
> you will report bugs that do not exist:
>
> 1. **Output capture — run CLI commands through a PTY.** The snap CLI is launched via
>    `snap run` and only produces its normal output when attached to a terminal. In this
>    shell a bare invocation may appear to print nothing. To test the CLI the way a real
>    user sees it, wrap every `"$SNAP_NAME" …` command in a pseudo-terminal using
>    `script`:
>    ```
>    script -qec '"$SNAP_NAME" <subcommand>' /dev/null
>    ```
>    `script` allocates a PTY, runs the command inside it, and writes everything the user
>    would see to stdout (discarding the typescript to `/dev/null`). Use this whenever you
>    need to read and verify a command's output. A command appearing to print nothing
>    when run *without* `script` is a harness artifact, **not** a snap defect — never
>    report it as one.
> 2. **Interactive prompts.** Some commands ask for confirmation before making a change.
>    If run bare here, they will **hang forever**, because there is no terminal to read from.
>    Exercise them two ways:
>    - **Non-interactively** with the documented `--assume-yes` flag (or whatever
>      equivalent the snap documents).
>    - **Interactively, as a real user would**, by feeding the answer into the PTY that
>      `script` creates:
>      ```
>      printf 'y\n' | script -qec '"$SNAP_NAME" <subcommand>' /dev/null
>      ```
>      `script` forwards its stdin to the prompt, so the command sees `y` (or `n`) just
>      as if you had typed it. Test **both** a `y` answer (the action proceeds) and an
>      `n` answer (the action is cancelled), and verify each behaves correctly.
>
>    Do NOT feed a prompt with `echo y |` or `yes |` **without** `script` — that does
>    not work because the prompt reads the terminal, not stdin, and the command will
>    hang. A hang in that case is a harness limitation, not a snap bug.
> 3. **Exit codes.** Check the exit status of the command **itself**, on its own line:
>    ```
>    "$SNAP_NAME" <subcommand>
>    echo "exit=$?"
>    ```
>    Do NOT prefix the command with `echo ... |` and do NOT use `${PIPESTATUS[0]}` to
>    read a command that is not first in a pipeline — you will measure the wrong
>    process's exit code. When you wrap a command in `script`, note that `script`'s exit
>    code reflects the wrapped command's exit code, so `script -qec '…' /dev/null; echo
>    "exit=$?"` is a reliable way to capture both output and status together.

1. Install the snap from the requested channel, passing the channel string exactly as
   given. **Installing can take several minutes** because the snap also downloads and
   installs large engine/model *components* after the snap itself.

   **Before installing, wait for snapd to be idle.** snapd may be busy refreshing itself
   or other snaps (e.g. right after the runner boots). Starting a new install while it is
   busy can cause the install to fail or race. Poll until all in-progress changes are done:

   ```
   # Wait for snapd to be idle before installing
   while snap changes | grep -qE '^[0-9]+ +(Do|Doing|Undo|Undoing|Wait) '; do
       echo "Waiting for snapd to finish in-progress changes..."
       snap changes | grep -E '^[0-9]+ +(Do|Doing|Undo|Undoing|Wait) '
       sleep 10
   done
   ```

   Then kick the install off without blocking and poll until snapd reports it finished:

   ```
   sudo snap install "$SNAP_NAME" --channel="$SNAP_CHANNEL" --no-wait
   # Poll in short steps — the install continues in the background even if a command times out.
   while snap changes "$SNAP_NAME" | grep -qE '^[0-9]+ +(Do|Doing|Undo|Undoing|Wait) '; do
       snap changes "$SNAP_NAME" | tail -n 3
       sleep 15
   done
   snap changes "$SNAP_NAME" | tail -n 5
   ```

   **If the install ends in `Error`**, check what failed before giving up. If the error
   comes from the install hook and appears to be a confinement or interface problem (e.g.
   a tool inside the snap could not access hardware it needed), retry the install with
   `--devmode` to bypass strict confinement:

   ```
   sudo snap install "$SNAP_NAME" --channel="$SNAP_CHANNEL" --devmode
   ```

   A snap that requires `--devmode` to install is a **genuine defect** — record it in
   your report and note that the rest of the testing was done in devmode. Continue with
   the remaining steps so you can test functionality despite the install issue.

   If the install error is something other than a hook/confinement failure (e.g. a network
   error, or the snap does not exist in that channel), report it and exit 1 immediately.

   > **Do not treat a slow or timed-out install as a failure.** `snap install` continues
   > inside snapd even if your foreground command is cut off by the harness timeout. While
   > it runs, a partially-installed state is completely normal — `snap list` may show
   > `components[1/3]`, `snap components "$SNAP_NAME"` may list some components as still
   > `available` (not yet `installed`), and the service will log
   > "Waiting for required snap components". This is expected progress, **not** a defect.
   >
   > **Never** `snap abort` the install change, and **never** work around it by manually
   > installing components (e.g. `snap install "$SNAP_NAME"+<component>`). Let the normal
   > install finish on its own. Only after the install change reaches `Done` should you
   > continue — then confirm the components required for your setup are `installed`:
   >
   > ```
   > snap components "$SNAP_NAME"
   > ```
   >
   > **Not every component listed is required.** A snap bundles components for multiple
   > engines/models and hardware types, and only the subset needed for the **selected
   > engine and model** on this machine is installed. It is normal and correct for the
   > others to stay `available` rather than `installed`. Components for unavailable
   > hardware (e.g. a CUDA engine component on a machine without an NVIDIA GPU) will
   > legitimately not be installed — that is expected, not a failure.
   >
   > Only treat this as a genuine defect if the install change ends in `Error`, or if a
   > component that *is* required for the selected engine and model is still missing after
   > the change completes.

2. Confirm it installed and inspect its interface connections:

   ```
   snap list "$SNAP_NAME"
   snap connections "$SNAP_NAME"
   ```

   **Save the full output of `snap list "$SNAP_NAME"` now** — you will need it verbatim
   for the test report's `environment.snap_list_output` field.

   > `snap list` may show `components[2/3]` or similar. This is **not** a problem — it
   > means not all components are installed, which is expected. Only the components
   > required for the selected engine and model on this machine are installed; the rest
   > stay available. See the install step above for details.

   These snaps are strictly confined. Interfaces (plugs) fall into two categories, and
   they must be judged differently:

   - **Plugs the snap declares as auto-connecting** (most hardware/observe interfaces a
     snap needs to do its job, e.g. `hardware-observe`). These are *supposed* to be
     connected automatically at install time. If one of these arrives **unconnected** and
     something fails because of it, that is a **genuine defect (`severity: error`)** — the
     snap is meant to work out of the box without the user running `snap connect`. Having
     to connect it by hand does **not** make it "working as intended"; the manual connect
     is your *proof* of the bug, not a fix for it.
   - **Truly optional plugs** that are documented as manual/opt-in and are not expected to
     auto-connect. A missing connection here is only a problem if it blocks something a
     user would reasonably expect to work; otherwise it is expected and not a finding.

   So: do not silently work around a missing auto-connect and pass the run. Only ignore an
   unconnected plug when it is genuinely optional/opt-in **and** nothing you test needs it.

   > **Diagnose confinement failures to their root cause.** A single unconnected
   > interface often breaks several different commands at once, which can look like many
   > separate bugs but is really **one** defect. When something fails in a way that could
   > be a confinement problem (permission denied, cannot read a device/file, empty or
   > partial hardware info, a feature silently unavailable), investigate the underlying
   > cause before writing it up:
   >
   > 1. Look for AppArmor denials produced while the command ran:
   >    ```
   >    sudo dmesg | grep -i 'apparmor="DENIED"' | tail -n 20
   >    sudo journalctl -k -g 'apparmor="DENIED"' --no-pager | tail -n 20
   >    ```
   >    The denial line names the profile and the operation/interface being blocked.
   > 2. Cross-check against `snap connections "$SNAP_NAME"` to see which plug is missing.
   >    A plug listed with no connection (or an interface the snap's `snap.yaml` marks as
   >    auto-connecting) that is not connected is a strong signal of an auto-connect defect.
   > 3. Confirm the diagnosis: connect the specific plug
   >    (`sudo snap connect "$SNAP_NAME":<plug>`), retry the failing command, and check
   >    the failure goes away. **Then record it as a `severity: error` finding** — the fact
   >    that a manual connect fixes it is the evidence that the plug should have
   >    auto-connected but did not. Do **not** downgrade it to a warning or omit it just
   >    because the command works once you connect it by hand.
   >
   > If one missing connection explains multiple symptoms, treat it as a **single**
   > root-cause finding (see "Consolidate symptoms of one root cause" under Reporting) —
   > name the specific interface (e.g. `hardware-observe`) and list the affected commands
   > as evidence, rather than filing one finding per broken command.

3. Capture machine information through a PTY and save it for the report:

   ```
   script -qec '"$SNAP_NAME" machine' /dev/null
   ```

   **Save the full output** — you will need it verbatim for the test report's
   `environment.show_machine_output` field. If `machine` does not exist or errors,
   record the error output instead so it is still available for issue reports.

4. Discover the available commands before using them, running the CLI through a PTY so
   its output is captured:

   ```
   script -qec '"$SNAP_NAME" --help' /dev/null
   ```

## What to test

Use the documentation and the snap's own `--help` output to work out what this snap is
for and how a user is meant to use it. Then exercise it as a user would, covering the
functionality it advertises. There is no fixed checklist — explore broadly rather than
following a script, so you notice anything that is wrong or missing.

Approach it like a curious user:

- Build a picture of the snap's purpose and its full set of commands/options from the
  docs and `--help`, then actually try them.
- For everything you run, read the output (through a PTY) and check it looks correct —
  expected content present, values sensible, no unexpected errors — rather than only
  checking the exit code.
- Confirm the snap's main purpose works end to end, including any background services or
  endpoints it relies on.
- Where the snap offers alternatives or options (modes, backends, models, settings, …),
  try switching between them and confirm the change actually takes effect.
- Allow for reasonable limitations of the environment (e.g. shared CI hardware, or
  lower output quality from a small model). Judge whether the snap behaves as its docs
  say it should, and note anything that does not.

## Reporting

**Your very last action MUST be to write a JSON report to `/tmp/snap-test-report.json`.**
The CI harness reads this file to determine whether the workflow passes or fails and to
create GitHub issues from any findings. Do not omit it or write anything other than
valid JSON.

### Schema

```json
{
  "verdict": "PASS",
  "summary": "One-sentence description of the overall result.",
  "environment": {
    "snap_name": "$SNAP_NAME",
    "snap_channel": "$SNAP_CHANNEL",
    "snap_version": "1.2.3",
    "snap_revision": "48",
    "snap_list_output": "<full output of: snap list $SNAP_NAME>",
    "show_machine_output": "<full output of: script -qec '$SNAP_NAME machine' /dev/null>",
    "devmode": false,
    "arch": "x86_64",
    "os": "Ubuntu 24.04.4 LTS",
    "ci_run_url": "<$GITHUB_SERVER_URL/$GITHUB_REPOSITORY/actions/runs/$GITHUB_RUN_ID, or null if not running in GitHub Actions>"
  },
  "findings": []
}
```

`verdict` is `"PASS"` or `"FAIL"`. Set it to `"FAIL"` if **any** finding has severity
`"error"`. Set it to `"PASS"` only if the snap genuinely works as documented with no
unresolved problems. Warnings and info findings alone do not make a verdict `"FAIL"`.

Each entry in `findings` must follow this structure:

```json
{
  "severity": "error",
  "title": "Short title suitable for a GitHub issue title",
  "description": "Clear description of what went wrong and what you were testing.",
  "reproduction": "Exact shell commands to reproduce, in order. This maps directly to the 'To reproduce' field in the GitHub issue template.",
  "observed": "Actual command output or error messages.",
  "expected": "What should have happened instead."
}
```

> **How findings map to the GitHub issue template:**
> - `title` → issue title
> - `description` + `observed` + `expected` → "Bug description" (the triage agent will
>   combine them: description of the problem, then observed vs expected behaviour)
> - `reproduction` → "To reproduce"
> - `environment.snap_list_output` → "Snap version"
> - `environment.show_machine_output` → "System information"
> - `environment.ci_run_url` → included in "Bug description" as a link to the CI run
>   (only when non-null)
>
> Fill every field with enough detail that someone reading the issue can understand and
> reproduce the problem without access to this CI run.

`severity` must be one of:
- `"error"` — a genuine defect: crash, incorrect output, missing documented functionality,
  or a confinement issue. A plug that the snap declares as auto-connecting (e.g.
  `hardware-observe`) but which arrives **unconnected** and breaks something is an
  `"error"` **even if the snap works once you connect it by hand** — needing a manual
  `snap connect` for an interface that should auto-connect is itself the defect.
- `"warning"` — something worth noting but does not block normal use (e.g. a suboptimal
  default, a misleading error message, or a genuinely optional/opt-in plug that is not
  connected but that nothing you tested actually needed).
- `"info"` — informational observation with no action required.


**Report root causes, not symptoms.** Before writing findings, group everything you
observed by underlying cause. When several failures share a single root cause — most
commonly one missing interface connection (e.g. `hardware-observe`) breaking several
commands — emit **one** finding for that root cause, not one per affected command.

For such a consolidated finding:
- Make the `title` name the root cause and the specific interface, e.g.
  "hardware-observe not auto-connected".
- In `description`, state the root cause once, then list each affected command as
  supporting evidence.
- Put the AppArmor denial line(s) and the `snap connections` excerpt that prove the
  diagnosis in `observed`.
- In `reproduction`, give the shortest command sequence that triggers the denial.

Only file separate findings when the failures genuinely have **different** root causes.
Two symptoms that both disappear after connecting the same plug are one finding.

### Writing the report

Collect all findings as you work, then at the very end run:

```bash
cat > /tmp/snap-test-report.json << 'REPORT'
{
  "verdict": "PASS or FAIL",
  "summary": "...",
  "environment": { ... },
  "findings": [ ... ]
}
REPORT
```

Make sure the output is valid JSON (no trailing commas, all strings quoted). If you have
no findings, write `"findings": []`.
