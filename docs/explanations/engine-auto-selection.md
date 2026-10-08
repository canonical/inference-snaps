# Engine auto selection

During installation of an inference snap, the most appropriate engine is selected from the set of available engines in the respective snap.
This reference describes the engine auto selection process.

At build time an inference snap is built with a set of engines, targeting a range of specific and generic hardware.
During install time the most appropriate one of these need to be chosen.

## Engine selection process

During installation, or when a user manually calls `<inference snap> use-engine --auto`, the following steps are performed:

1. A summary of available compute hardware is made using the `lscompute` library.
2. The list of available engines is filtered to remove any engines that can not run on the available hardware.
3. The remaining engines are sorted by how specifically they target and match the available hardware.

We'll discuss step 2 and 3 in more detail below.

## Filtering step

An engine lists required hardware under the devices section in the engine manifest file.
Some devices are always required, and listed under `all-of`.
Other devices can have multiple options, and are listed under `any-of`.

For example if an engine always requires an arm64 cpu, it will be listed under `all-of`.
If an engine can run on either an AMD GPU or an NVIDIA GPU, the two GPU devices, with their vendor IDs will be listed under `any-of`.

During filtering, the list of all-of and any-of devices are compared to the list of available compute devices on the host system.
If the host's hardware matches all of the devices listed under `all-of`, and at least one of the devices listed under `any-of`, the engine is considered compatible with the host system.
Otherwise it is filtered out and not considered for selection.

## Sorting step

The device definition in the engine manifest has a number of optional properties.
Some of them can be generic like CPU architecture, while other properties can be specific like PCI device ID.

The more properties are defined, the more specifically this engine targets a specific hardware.
An engine that targets a specific device is considered to perform better than an engine that targets a generic device.

To sort engines from least specific (more generic) to most specific (more targeted), a scoring mechanism is used.
Each property that is listed in the device definition contributes a "weight" to the score.
A higher score therefore indicates that a device definition, or in other words an engine, targets and matches the host system more specifically that an engine with a lower score.

After filtering, all remaining engines are scored, and the engine with the highest score is selected as the most appropriate engine for the host system.

