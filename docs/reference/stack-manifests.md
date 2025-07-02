# Stack manifests

Each stack is accompanied by a manifest file outlining the stack’s components 
and configurations. The manifest also describes the stack’s system requirements. 
It is required during the selection process in order to pick the most 
{spellexception}`optimized` stack for the host machine.

The stack manifest should include the following attributes:

`name`  
A unique name for the stack.   

`description`   
the description of the stack  

`vendor`   
the name of the responsible {spellexception}`organization`  

`grade`  
indicates the stability and compliance of the stack. Only stable stacks 
get auto selected. (supported values: `stable` and `devel`)  

`devices`  
lists of required computing devices. It’s possible to indicate if more than one
device is required, or if there are multiple options, using the `any` and `all`
keywords. For example, a stack requiring a generic CPU, combined with either an
Nvidia GPU with compute capability 7.0 or an AMD GPU with 8GB vRAM should list 
the CPU under `all` and GPUs under `any`.  
  - `type` - type of the device (cpu, gpu, npu, nil)  
  - Remaining field types depend on the type. See other device-specific 
  properties below.  
  
`memory`   
required system memory to load the model  

`disk-space`   
the total size of stack components plus runtime additions  

`components`  
list of snap components required by the stack  

`configurations`  
default snap configurations  
  - `engine` - name of engine snap component (mandatory)  
  - `model` - one of:  
    * Name of model snap component  
    * Path to local directory or file containing the model  
    * nil \- indicates that a separate model is not needed, i.e., when one is 
    already embedded in the engine (e.g. [llamafile](https://github.com/Mozilla-Ocho/llamafile))

Device entries may also have the following device-specific fields:  

**For CPUs:**  
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

**For PCI peripherals**

* `bus`: `pci`  
* `vendor-id` - PCI vendor ID hex number  
* `device-id` - PCI device ID hex number  
* `device-class` - not used  
* `subvendor-id` - not used  
* `subdevice-id` - not used  
* `programming-interface` - not used

**For GPUs:**

* `bus` - the bus or protocol used by the device (supported values: `pci`, default: `pci`)  
* `vram` - minimum required Video RAM  
* `compute-capability` - applicable to NVIDIA (vendor-id 0x10de) only. A list 
of MAJOR.MINOR version strings prefixed with a comparison operator (==, \<=, \>=, \>, \<[^1]). 
MAJOR and MINOR must both be integers.

When bus is set to PCI, this object inherits all PCI peripheral fields.

**For NPUs:**

* `bus` - the bus or protocol used by the device (supported values: `pci`, default: `pci`)

When bus is set to PCI, this object inherits all PCI peripheral fields.

**TPUs:** not used

* `bus` - the bus or protocol used by the device (supported values: `pci`, `usb`, default: `pci`)

Additional rules:

* The {spellexception}`serialization` omits prefixes such as “required”, 
“minimum”. All listed devices, memory, disk space, and components are to be 
treated as requirements to use the stack.   
* If a stack is supported on a limited set of architectures, it should define
that requirement via a CPU entry.  
* Disk/memory capacities should be given either in gigabytes or megabytes, with
**G** or **M** as suffixes.  
* Device Bus is to uniquely identify the vendor, as Vendor IDs are only unique
within their own subsystem. For example, the Nvidia vendor ID is 0x10DE for PCI
and 0x0955 for USB devices. For graphics accelerators, we only expect pci or
usb as the possible values. For the full list from Linux, 
refer [here](https://github.com/torvalds/linux/blob/158f238aa69d91ad74e535c73f552bd4b025109c/scripts/mod/file2alias.c#L1546-L1599).

Example YAML {spellexception}`serialization` of a stack manifest:

```yaml
  # metadata
  name: stack-template # <vendor>-<id>
  description: This is a stack definition template
  vendor: Stack Ltd
  grade: stable

  # hardware requirements
  devices:
    # at least one is required
    any:
      - type: cpu
        manufacturer-id: AuthenticAMD
        architecture: amd64
      - type: cpu
        manufacturer-id: GenuineIntel
        architecture: amd64
        flags: [avx2, avx512]
    # all are required
    all:
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
    engine: llamacpp
    model: model-q4-k-m-gguf

```

[^1]:  The comparison operators used for Compute Capability are based on 
[PEP](https://peps.python.org/pep-0440/#version-specifiers) 
and [Snapcraft’s Python plugin](https://snapcraft.io/docs/python-plugin). 
When compared to [npm’s semantic versioning](https://docs.npmjs.com/about-semantic-versioning) 
and [Debian package’s Depend syntax](https://www.debian.org/doc/debian-policy/ch-relationships.html), 
all operators except the equal operator match.