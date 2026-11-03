# Inference Snaps Dev

This project contains developer tools for building, uploading, and testing Inference Snaps.

## Reusable Workflows
A few Reusable Workflows are available in the `.github/workflows` folder:

- [build-publish-snap](./.github/workflows/build-publish-snap.yaml) - to build and publish a snap with its components
- [reuse-cicd](./.github/workflows/reuse-cicd.yaml) - to run CICD on push events from the default branch
- [reuse-pr-build-test](./.github/workflows/reuse-pr-build-test.yaml) - to build or build+test a snap depending on the PR label that triggered it
- [reuse-triage-failures](./.github/workflows/reuse-triage-failures.yaml) - to triage and retry failed jobs classified as failures

## Actions
A few Actions are available in the `.github/actions` folder:

- [agentic-test](./.github/actions/agentic-test/action.yaml) - to run the agentic snap test via Workshop & OpenCode
- [remove-label](./.github/actions/remove-label/action.yaml) - to remove a label from a PR after the associated workflow has run
- [smoke-test](./.github/actions/smoke-test/action.yaml) - to run smoke tests against an installed snap

## Scripts
This project contains a few scripts to help build, install, and upload the snaps:

- [build.sh](./build.sh) - to build the snap and its components
- [install.sh](./install.sh) - to install the built snap
- [upload.sh](./upload.sh) - to upload the built snap to the Snap Store
- [smoke-test.sh](./smoke-test.sh) - to run smoke tests against an installed snap

Add this repo as a submodule to your snap project:
```shell
git submodule add --branch v2 https://github.com/canonical/inference-snaps-dev dev
```

Use the scripts from the root of the snap repo, e.g.:
```
./dev/upload.sh
```
