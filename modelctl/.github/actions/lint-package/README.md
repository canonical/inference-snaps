# Lint package action

This GitHub Action validates the package by running:

```bash
modelctl debug lint-package <package-dir>
```

## Inputs

| Name | Required | Default | Description |
| --- | --- | --- | --- |
| `package-dir` | No | `.` | Path (relative to `${{ github.workspace }}`) to the snap package source that contains `snap/snapcraft.yaml` as well as `engines`, `models`, and `runtimes` directories |

## Usage

```yaml
jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v7

      - name: Lint package manifests
        uses: canonical/inference-snaps-cli/.github/actions/lint-package@<tag>
        with:
          package-dir: .
```
