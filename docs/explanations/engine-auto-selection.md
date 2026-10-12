# Engine auto selection

An inference snap includes engines targeting a range of specific and generic hardware.
During installation, automatic selection chooses the most appropriate engine for the host system.

## Engine selection process

During installation, or when a user manually calls `<inference-snap> use-engine --auto`, the following steps are performed:

1. A summary of available compute hardware is made using the `lscompute` library.
2. The list of available engines is filtered to remove any engines whose hardware requirements are not satisfied.
3. The remaining engines are ranked according to their hardware matches and preferences for compute hardware.

## Filtering step

An engine lists required hardware under the `devices` section in the engine manifest file.
Devices that are always required are listed under `allof`.
Alternative devices are listed under `anyof`.

For example, if an engine always requires an `arm64` CPU, it will be listed under `allof`.
If an engine can run on either an AMD GPU or an NVIDIA GPU, the two GPU devices, with their vendor IDs, will be listed under `anyof`.

During filtering, the device requirements are compared to the list of available compute devices on the host system.
Requirements under `allof` must all be satisfied.
Requirements under `anyof` provide alternatives, at least one of which must be satisfied.
An engine is considered compatible when the host system satisfies its device requirements.
Otherwise it is filtered out and not considered for selection.

Experimental engines are excluded from automatic selection.

## Sorting step

Compatible engines are ranked using a score that reflects how they match the host hardware.
The highest-scoring eligible engine is selected.

An engine manifest can describe broad requirements, such as a CPU architecture, or more specific ones, such as a PCI device ID.
Defining more hardware properties generally makes the requirements more specific.
When these requirements match the host hardware, a more specific match generally earns a higher score.

The score also accounts for hardware preferences, such as favoring discrete GPUs over integrated GPUs.
GPU and NPU engines can therefore rank above CPU-only engines, but no device type is guaranteed to take priority: the overall match determines the ranking.
