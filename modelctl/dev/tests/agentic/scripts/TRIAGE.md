# Issue Triage Agent Instructions

Your job is to compare the error findings from a snap test run against the currently open
issues on the `$ISSUE_REPO` GitHub repository (defaults to `canonical/inference-snaps`),
and report which findings are already tracked and which are new.

## Steps

1. **Read the test report.**

   ```
   cat /tmp/snap-test-report.json
   ```

   Extract only the findings with `"severity": "error"`. If there are none, print
   "No error findings — nothing to triage." and stop.

2. **Fetch the bug report issue template** from the repository so your suggested bodies
   always match the current template structure.

   ```
   curl -s "https://raw.githubusercontent.com/$ISSUE_REPO/main/.github/ISSUE_TEMPLATE/bug_report.yaml"
   ```

   Read the `body` fields to understand what sections are expected and in what order.
   Use the `label`, `description`, and `value` of each field to guide how you fill in
   the content for each finding.

3. **Fetch open issues from GitHub** (the repository is public; no token required).

   Fetch only **open** issues — these are the duplicate targets:

   ```
   curl -s "https://api.github.com/repos/$ISSUE_REPO/issues?state=open&per_page=100"
   ```

   If you get back fewer than 100 issues, one page is enough. If the response contains
   exactly 100 items, fetch the next page by appending `&page=2`, and continue until a
   page returns fewer than 100 items.

   Collect all fetched issues before proceeding. For each issue keep: number, title, body
   (first 500 chars is enough), html_url, and the list of label names (`labels[].name`).
   The label names matter for reporting cross-snap duplicates: issues are labelled with a
   `snap/<snap name>` label identifying which snap they were detected on.

4. **Triage each error finding.**

   Match on **root cause**, not surface wording. Two findings describe the same problem
   when they share the same underlying cause, even if the affected command, error text,
   or snap revision differ. In particular, a missing/unconnected interface (e.g.
   `hardware-observe` not auto-connecting) is a single root cause that can surface as
   several different command failures — treat all of them as the same problem.

   4a. **Collapse same-run duplicates first.** Before comparing against GitHub, group the
   run's own error findings by root cause. If several findings share one cause (same
   missing interface, same confinement denial, same crashing component), keep the
   clearest one as the representative and mark the others as duplicates *of that
   representative finding* (Status `DUPLICATE`, and note in the Reason that it duplicates
   another finding in this same run). Do not file multiple GitHub issues for one cause.

   4b. **Compare each representative finding against the collected open issues.** Decide
   whether an existing issue describes the same underlying problem closely enough that
   filing a new issue would be a duplicate. Use the finding's `title`, `description`,
   `observed`, and `expected` fields for comparison. A match is "close enough" when the
   existing issue is clearly about the **same failure mode / root cause in the same snap
   component** — for example both are caused by the same interface not being connected —
   even if it was reported on an earlier revision. A shared root cause outweighs
   differences in the exact command or error string. Do **not** match on vague keyword
   overlap alone.

5. **Print the triage report** in this exact format so it is easy to read in CI logs:

   ```
   === ISSUE TRIAGE REPORT ===

   Finding 1: <finding title>
   Status   : DUPLICATE
   Issue    : #<number> — <title>
   URL      : <html_url>
   Reason   : <one sentence explaining why this is the same root cause; if it duplicates
              another finding in this same run, say so instead of citing a GitHub issue>
   ---
   Suggested title:
   <A concise GitHub issue title>

   Suggested body:
   <Issue body filled in using the section structure from the fetched bug_report.yaml
   template. Populate each field using the finding and report data:
   - Bug description    ← finding.description, then observed vs expected behaviour;
                          if report.environment.ci_run_url is non-null, append a line:
                          "Detected by CI run: <ci_run_url>"
   - To reproduce       ← finding.reproduction
   - Snap version       ← report.environment.snap_list_output (as a shell code block)
   - System information ← report.environment.show_machine_output (wrapped in the
                          <details> block the template specifies)>

   ---

   Finding 2: <finding title>
   Status   : NEW
   ---
   Suggested title:
   <A concise GitHub issue title>

   Suggested body:
   <same structure as above>

   ---
   ```

   Print one block per error finding. Use `DUPLICATE` or `NEW` for Status.
   Always include the suggested title and body regardless of status — for
   duplicates this makes it easy to update the existing issue with fresh reproduction
   steps and environment details.
   Do not skip any error finding.

6. Write the triage results as JSON to `/tmp/snap-triage-report.json` using this schema:

   ```json
   {
     "findings_triaged": 2,
     "new_count": 1,
     "duplicate_count": 1,
     "results": [
         {
           "finding_title": "...",
           "status": "NEW",
           "suggested_title": "...",
           "suggested_body": "<issue body using sections from the fetched bug_report.yaml template>"
         },
          {
            "finding_title": "...",
            "status": "DUPLICATE",
            "duplicate_of_number": 42,
            "duplicate_of_url": "https://github.com/canonical/inference-snaps/issues/42",
            "duplicate_of_labels": ["bot", "snap/other-snap"],
            "duplicate_of_finding": null,
            "duplicate_reason": "...",
            "suggested_title": "...",
            "suggested_body": "<issue body using sections from the fetched bug_report.yaml template>"
          }
      ]
    }
    ```

   Notes on the duplicate fields:
   - `duplicate_of_number`/`duplicate_of_url` reference the matched **open** GitHub issue.
   - `duplicate_of_labels` is the list of label names currently on the matched GitHub issue
     (copy them verbatim). This lets the workflow detect when the duplicate was originally
     detected on a **different** snap (its labels contain a `snap/<other snap>` label but
     not this run's `snap/<snap name>` label) so it can add this snap's label to the shared
     issue. Set it to `null` for same-run duplicates.
   - `duplicate_of_finding` is the `finding_title` of another finding **in this same run**
     when this finding was collapsed as a same-run duplicate (step 4a); in that case
     `duplicate_of_number`/`duplicate_of_url`/`duplicate_of_labels` may be null.

