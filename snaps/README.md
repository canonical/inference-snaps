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

## Building inference snaps

When cloning the projects as submodules, the `.git` directory will not be present within each submodule.
Instead, the `.git` information is stored in the parent repository and the submodule directory contains a reference to it.
This prevents Snapcraft from accessing the Git directory during the snap build process, and as a result the build fails when extracting version information.

To build the inference snaps it is necessary to convert submodules into standalone Git repositories.
This can be done by running the following command:

```shell
git submodule foreach --recursive '
  REAL_GITDIR=$(git rev-parse --git-dir)
  if [ -f .git ]; then
    rm .git
    cp -r "$REAL_GITDIR" .git
    # Remove the old relative path to the parent working directory
    # so Git defaults to using the current directory instead
    git config --local --unset core.worktree
  fi
'
```

Next, enter the snap subdirectory and run `make` for usage instructions.

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
