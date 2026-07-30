# Inference snap repositories

The inference snaps are maintained in separate repositories, referenced here as submodules.

Initialize the submodules:
```shell
git submodule update --init --recursive
```

The submodules are updated once a week via Renovate.
To update them earlier to the latest commit on their main branches, run:
```shell
git submodule update --remote
```
