(engine-manifest)=

# Engine manifest

Each engine is accompanied by an `engine.yaml` manifest file identifying its
runtime, supported models, default configurations, and required computing devices.
The following sections describe the accepted fields and their validation
requirements.

## `name`

**Type:** string. **Required:** yes.

A non-empty, unique name for the engine. It must match the name of the directory
containing the `engine.yaml` file.

## `summary`

**Type:** string. **Required:** yes.

A non-empty, short description of the engine, displayed in the engines table.
The maximum length is 56 bytes (56 characters for ASCII text).

## `description`

**Type:** string. **Required:** no.

A longer description of the engine.

## `vendor`

**Type:** string. **Required:** yes.

The non-empty name of the responsible {spellexception}`organization`.

## `experimental`

**Type:** boolean. **Required:** no. **Default:** `false`.

Indicates whether the engine is experimental. Engines with
`experimental: true` are excluded from automatic selection. Omitting the field
or setting it to `false` makes the engine non-experimental.

## `devices`

**Type:** mapping. **Required:** no.

Required computing devices, grouped under two optional keys containing lists of
device mappings:

- `allof`: every listed device requirement must be satisfied.
- `anyof`: at least one listed device requirement must be satisfied when the list
  is non-empty.

When both lists are provided, both conditions must be satisfied. For example, an
engine requiring an `amd64` CPU and either an NVIDIA or AMD GPU would list the
CPU under `allof` and the alternative GPUs under `anyof`.

Each entry has an optional string field `type` with values `cpu`, `gpu`, or `npu`.
An omitted or empty `type` matches peripherals without restricting the device
type. The remaining fields depend on the device type and bus. See
{ref}`device-specific properties <device-specific-fields>` below.

## `runtime`

**Type:** string. **Required:** no.

The name of the {ref}`runtime manifest <runtime-manifest>` used by the engine,
for example `llamacpp-cuda`.

## `model`

**Type:** mapping. **Required:** no.

The models supported by the engine:

- `default` (string): the model name used by default when selecting the engine.
- `options` (list of strings): the model names available for the engine.

These names reference {ref}`model manifests <model-manifest>`.

## `configurations`

**Type:** mapping. **Required:** no.

Default engine configurations, with string keys and scalar values (strings,
numbers, or booleans). Nested mappings and lists are not supported.

(device-specific-fields)=
## Device-specific fields

Device entries may also have the following device-specific fields:

### CPUs

`architecture`: CPU architecture in Debian nomenclature (`amd64` or `arm64`).
This field is mandatory for `type: cpu`. The remaining fields are optional and
architecture-specific:

| Architecture | Field | Meaning |
| --- | --- | --- |
| `amd64` | `manufacturer-id` | Manufacturer string reported by the `CPUID` instruction, such as `GenuineIntel` or `AuthenticAMD`. |
| `amd64` | `flags` | List of required CPU flags. |
| `arm64` | `implementer-id` | Implementer ID from the Main ID register, as a hexadecimal number. |
| `arm64` | `part-number` | Part number from the Main ID register, as a hexadecimal number. |
| `arm64` | `features` | List of required CPU features. |

CPU entries do not accept `bus`, PCI identifiers, or `snap-connections`.

(pci-peripherals)=
### PCI peripherals

- `bus`: `pci`. This is the default when `bus` is omitted for a non-CPU device.
- `vendor-id`: PCI vendor ID as a hexadecimal number, for example `0x10de`.
- `device-id`: PCI device ID as a hexadecimal number. Matched only when
  `vendor-id` is also specified and matches the device vendor.
- `snap-connections`: a list of snap plug names that must be connected for the
  device to be considered compatible, for example `[intel-npu, npu-libs]`.

The identifiers and `snap-connections` are optional.

### GPUs

GPU entries use the {ref}`PCI peripheral <pci-peripherals>` fields and may also
specify:

- `vram`: minimum required video RAM as a size string. Values without a suffix
  are in bytes. The suffixes `M` and `G` denote multiples of 1024 squared and
  1024 cubed bytes, respectively, for example `512M` or `4G`.
- `compute-capability`: a single version-constraint string for NVIDIA GPUs
  (`vendor-id: 0x10de`), not a YAML list. Constraints use
  [Masterminds semver syntax](https://github.com/Masterminds/semver#checking-version-constraints).
  For example, `">=6.0, <7.0"` requires both comparisons to match, while
  `"=5.3 || >=6.2"` accepts either alternative.
- `microarchitecture`: the required GPU microarchitecture string, matched
  exactly, for example `gfx1030` for AMD (`vendor-id: 0x1002`).

If a requested GPU property is not reported by the host, that requirement is not
satisfied.

### NPUs

NPU entries support `bus: pci` (the default) or `bus: fastrpc`.

PCI entries accept the {ref}`PCI peripheral <pci-peripherals>` fields. FastRPC
entries accept only `type`, `bus`, and `snap-connections`, not PCI identifiers.
FastRPC also permits an omitted or empty `type`, but not `cpu` or `gpu`.

### USB devices

`usb` is a recognized bus value, but USB device validation and matching are not
implemented. Entries with `bus: usb` fail validation.

### Compatibility diagnostics

`compatibility-issues` is a list of diagnostic strings populated by the selector
when device requirements are not met. It can appear in serialized results, but
must not be supplied in an authored engine manifest.

## Example YAML serialization of an engine manifest

```yaml
# engines/nvidia-gpu/engine.yaml
name: nvidia-gpu
summary: CUDA engine for NVIDIA GPUs
description: Runs a supported model on an NVIDIA GPU.
vendor: Engine Ltd
experimental: false

devices:
  # At least one CPU architecture must match.
  anyof:
    - type: cpu
      architecture: amd64
    - type: cpu
      architecture: arm64
  # The GPU is always required.
  allof:
    - type: gpu
      bus: pci
      vendor-id: 0x10de
      vram: 4G
      compute-capability: "=5.3 || >=6.2"

runtime: llamacpp-cuda

model:
  default: smollm2-135m
  options:
    - smollm2-135m

configurations:
  sleep-idle-seconds: 600
  min-context-size: 4096
```
