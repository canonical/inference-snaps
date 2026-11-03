# Bash completion script for modelctl packaged in a snap
#
# To use, in snapcraft.yaml:
# 1. Set apps.<app-name>.completer to the name of this script for the app.
# 2. Add a symlink from <app> to modelctl binary.
#
# For example, snapcraft.yaml:
# name: mymodel
# <...>
# parts:
#   cli:
#     source: <...>
#     plugin: dump
#     override-build: |
#       # For tab completion
#       ln --symbolic ./modelctl bin/mymodel
#       craftctl default
# apps:
#   mymodel: # command named after the snap
#     command: bin/modelctl
#     completer: bin/snap-completer.bash
#

# Unset the _init_completion function from the bash-completion package to force
# use of the basic but functional internal implementation.
# Issue: https://github.com/canonical/inference-snaps/issues/183
unset -f _init_completion

source <($SNAP/bin/modelctl completion bash)
