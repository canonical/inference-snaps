# Engine manifest

Each engine is accompanied by a manifest file outlining the engine’s components
and configurations. The manifest also describes the engine’s system requirements.
It is required during the selection process in order to pick the most
{spellexception}`optimized` engine for the host machine.

The engine manifest should include the following attributes:

## `name`
A unique name for the engine.

## `description`
The description of the engine

## `vendor`
The name of the responsible {spellexception}`organization`

## `grade`
Indicates the stability and compliance of the engine. Only stable engines
get auto selected. (supported values: `stable` and `devel`)

## `devices`
Lists of required computing devices. It’s possible to indicate if more than one
device is required, or if there are multiple options, using the `anyof` and `allof`
keywords. For example, a engine requiring a generic CPU, combined with either an
Nvidia GPU with compute capability 7.0 or an AMD GPU with 8GB vRAM should list
the CPU under `allof` and GPUs under `anyof`.
  - `type` - type of the device (cpu, gpu, npu, nil)
  - Remaining field types depend on the type. See other
  {ref}`device-specific properties <device-specific-fields>` below.

## `memory`
Required system memory to load the model. Memory capacities should be given
either in gigabytes or megabytes, with **G** or **M** as suffixes.

## `disk-space`
The total size of engine components plus runtime additions. Disk capacities
should be given either in gigabytes or megabytes, with **G** or **M** as
suffixes.

## `components`
List of snap components required by the engine

# `configurations`
Default snap configurations
  - `engine` - name of engine snap component (mandatory)
  - `model` - one of:
    * Name of model snap component
    * Path to local directory or file containing the model
    * nil \- indicates that a separate model is not needed, i.e., when one is
    already embedded in the engine (e.g. [llamafile](https://github.com/Mozilla-Ocho/llamafile))

(device-specific-fields)=
## Device-specific fields
Device entries may also have the following device-specific fields:

### CPUs
`architectures`: CPU architecture in Debian nomenclature (amd64, arm64). This
field is  mandatory. The remaining fields are architecture specific.

- amd64:
  * `manufacturer-id` - reported by `CPUID` instruction
  * `flags` - list of required CPU flags
  * `family-id` - not used
* arm64:  
  * `implementer-id` - from Main ID register
  * `part-number` - from Main ID register
  * `features` - not used

(pci-peripherals)=
### PCI peripherals

* `bus`: `pci`
* `vendor-id` - PCI vendor ID hex number.
* `device-id` - PCI device ID hex number
* `device-class` - not used
* `subvendor-id` - not used
* `subdevice-id` - not used
* `programming-interface` - not used

### GPUs

* `bus` - the bus or protocol used by the device (supported values: `pci`, default: `pci`)
* `vram` - minimum required Video RAM
* `compute-capability` - applicable to NVIDIA (vendor-id `0x10de`) only. A list
of MAJOR.MINOR version strings prefixed with a comparison operator (==, \<=, \>=, \>, \<[^1]).
MAJOR and MINOR must both be integers.

When bus is set to PCI, this object inherits all {ref}`PCI peripheral <pci-peripherals>` fields.

### NPUs

* `bus` - the bus or protocol used by the device (supported values: `pci`, default: `pci`)

When bus is set to PCI, this object inherits all {ref}`PCI peripheral <pci-peripherals>` fields.

## Example YAML {spellexception}`serialization` of an engine manifest

```yaml
  # metadata
  name: engine-template # <vendor>-<id>
  description: This is a engine definition template
  vendor: Engine Ltd
  grade: stable

  # hardware requirements
  devices:
    # at least one is required
    anyof:
      - type: cpu
        manufacturer-id: AuthenticAMD
        architecture: amd64
      - type: cpu
        manufacturer-id: GenuineIntel
        architecture: amd64
        flags: [avx2, avx512]
    # all are required
    allof:
      - type: gpu
        vendor-id: "0x10de" # Nvidia Corporation
        # device-id:
        vram: 4G
        compute-capability: ["==5.3", ">=6.2"]
  memory: 4G
  disk-space: 15G

  # required snap components
  components:
    - llamacpp
    - model-q4-k-m-gguf

  # default config
  configurations:
    key: value
```

[^1]:  The comparison operators used for Compute Capability are based on
[PEP](https://peps.python.org/pep-0440/#version-specifiers)
and [Snapcraft’s Python plugin](https://snapcraft.io/docs/python-plugin).
When compared to [npm’s semantic versioning](https://docs.npmjs.com/about-semantic-versioning)
and [Debian package’s Depend syntax](https://www.debian.org/doc/debian-policy/ch-relationships.html),
all operators except the equal operator match.
