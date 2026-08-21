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

## Adding a new snap submodule

To add a new snap repository as submodule, run:
```shell
git submodule add --name <name>-snap --branch main <repository-url> snaps/<name>
```
Where `name` is the repository name without the `-snap` suffix.

<!--
TODO:
Simplify the above command to:
$ git submodule add --branch main <repository-url> snaps/<name>
Defaulting name to the path, and update existing submodules to use this simpler form.
-->
