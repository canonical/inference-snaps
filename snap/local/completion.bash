# Bash completion for the app named after the snap

# Unset the _init_completion function from the bash-completion package to force
# use of the basic but functional internal implementation.
# Issue: https://github.com/canonical/stack-utils/issues/115
unset -f _init_completion

source <($SNAP/bin/modelctl completion bash)
