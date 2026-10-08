# Engine auto selection

An inference snap includes engines targeting a range of specific and generic hardware.
During installation, automatic selection chooses the most appropriate engine for the host system.

## Engine selection process

During installation, or when a user manually calls `<inference snap> use-engine --auto`, the following steps are performed:

1. A summary of available compute hardware is made using the `lscompute` library.
2. The list of available engines is filtered to remove any engines whose hardware requirements are not satisfied.
3. The remaining engines are sorted by how specifically they target and match the available hardware.

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

The device definition in the engine manifest has a number of optional properties.
Some of them can be generic like CPU architecture, while other properties can be specific like PCI device ID.

In general, defining more hardware properties makes an engine's requirements more specific.
Matching more specific hardware requirements generally results in a higher score.
Compatible engines are ranked using these scores, with higher scores indicating a stronger preference for an engine based on its match to the host hardware.
The score is not a measurement or guarantee of performance.

The highest-scoring eligible engine is selected as the most appropriate engine for the host system.
