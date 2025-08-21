# Publish snaps

This is a Github Action that can be used to publish snap packages to the Snap Store.

Unlike [`canonical/action-publish`](https://github.com/canonical/action-publish), this action also uploads the component snaps defined in `snap/snapcraft.yaml`.

## Basic Usage

If you wish to define your job inline, you can use the following step:

```yaml
jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Build snap
        uses: canonical/action-build@v1

      - name: Publish snap
        uses: canonical/stack-utils/.github/actions/snap-publish
        with:
          store-credentials: ${{ secrets.STORE_LOGIN }}
```

## Store login

This action requires a Snap Store login secret saved in Github as secret. Export a token as `secret.json` and add its content as a repository secret named `STORE_LOGIN`:

```bash
snapcraft export-login \
    --snaps "<snap-name>" \
    --channels "latest/edge/*" \
    --acls package_access,package_push,package_update,package_release \
    --expires 2026-08-15T00:00:00Z \
    secret.json
```

## API

### Inputs

| Name | Description | Default | Required |
|---|---|---|---|
| `store-credentials` | Snap Store credentials to publish the snap. | `""` | true |

## How it works

**Publish step**  
Uses a bundled `upload.sh` script to parse `snap/snapcraft.yaml` and validate snap and components to upload them all to the Snap Store using `snapcraft upload --release=<channel>`

  - On branch/tag builds: channel is defined to `latest/edge/<short-sha>`  
  - On pull request builds: channel is defined to `latest/edge/pr-<number>`
