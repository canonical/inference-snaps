# Stack manifests

Each stack is accompanied by a manifest file outlining the stack’s components and configurations. The manifest also describes the stack’s system requirements. It is required during the selection process in order to pick the most optimized stack for the host machine.

The stack manifest should include the following attributes:

* `name` - A unique name for the stack.   
* `description` - the description of the stack  
* **vendor** \- the name of the responsible organization  
* **grade** \- indicates the stability and compliance of the stack. Only stable stacks get auto selected. (supported values: `stable` and `devel`)  
* **devices** \- lists of required computing devices. It’s possible to indicate if more than one device is required, or if there are multiple options, using the `any` and `all` keywords. For example, a stack requiring a generic CPU, combined with either an Nvidia GPU with compute capability 7.0 or an AMD GPU with 8GB vRAM should list the CPU under `all` and GPUs under `any`.  
  * **type** \- type of the device (cpu, gpu, npu, nil)  
  * Remaining field types depend on the type. See other device-specific properties below.  
* **memory** \- required system memory to load the model  
* **disk-space** \- the total size of stack components plus runtime additions  
* **components** \- list of snap components required by the stack  
* **configurations** \- default snap configurations  
  * engine \- name of engine snap component (mandatory)  
  * model \- one of:  
    * Name of model snap component  
    * Path to local directory or file containing the model  
    * nil \- indicates that a separate model is not needed, i.e., when one is already embedded in the engine (e.g. [llamafile](https://github.com/Mozilla-Ocho/llamafile))

Device entries may also have the following device-specific fields:  
For CPUs:  
**architectures \-** CPU architecture in Debian nomenclature (amd64, arm64). This field is  mandatory. The remaining fields are architecture specific.

* amd64:  
  * **manufacturer-id** \- reported by CPUID instruction  
  * **flags** \- list of required CPU flags  
  * *family-id* \- not used   
* arm64:  
  * **implementer-id** \- from Main ID register  
  * **part-number** \- from Main ID register  
  * *features* \- not used

**For PCI peripherals**

* **bus: pci**  
* **vendor-id** \- PCI vendor ID hex number  
* **device-id** \- PCI device ID hex number  
* *device-class \-* not used  
* *subvendor-id \-* not used  
* *subdevice-id \-* not used  
* *programming-interface \-* not used

**For GPUs:**

* **bus** \- the bus or protocol used by the device (supported values: pci, default: pci)  
* **vram** \- minimum required Video RAM  
* **compute-capability** \- applicable to NVIDIA (vendor-id 0x10de) only. A list of MAJOR.MINOR version strings prefixed with a comparison operator (==, \<=, \>=, \>, \<[^1]). MAJOR and MINOR must both be integers.

When bus is set to PCI, this object inherits all PCI peripheral fields.

**For NPUs:**

* **bus** \- the bus or protocol used by the device (supported values: pci, default: pci)

When bus is set to PCI, this object inherits all PCI peripheral fields.

**TPUs:** not used

* **bus** \- the bus or protocol used by the device (supported values: pci, usb, default: pci)

Additional rules:

* The serialization omits prefixes such as “required”, “minimum”. All listed devices, memory, disk space, and components are to be treated as requirements to use the stack.   
* If a stack is supported on a limited set of architectures, it should define that requirement via a CPU entry.  
* Disk/memory capacities should be given either in gigabytes or megabytes, with **G** or **M** as suffixes.  
* Device Bus is to uniquely identify the vendor, as Vendor IDs are only unique within their own subsystem. For example, the Nvidia vendor ID is 0x10DE for PCI and 0x0955 for USB devices. For graphics accelerators, we only expect pci or usb as the possible values. For the full list from Linux, refer [here](https://github.com/torvalds/linux/blob/158f238aa69d91ad74e535c73f552bd4b025109c/scripts/mod/file2alias.c#L1546-L1599).

Example YAML serialization of a stack manifest:

```YAML
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
      compute-capability: [==5.3, >=6.2]
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

[^1]:  The comparison operators used for Compute Capability are based on [PEP](https://peps.python.org/pep-0440/#version-specifiers) and [Snapcraft’s Python plugin](https://snapcraft.io/docs/python-plugin). When compared to [npm’s semantic versioning](https://docs.npmjs.com/about-semantic-versioning) and [Debian package’s Depend syntax](https://www.debian.org/doc/debian-policy/ch-relationships.html), all operators except the equal operator match.

[image1]: <data:image/png;base64,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M48CghYAAGUQtBpG0KqMoAUAABAJQQsAACASghYAAEAkBC0AAIBICFoAAACRELQAAAAiIWgBAABEQtACAACIhKAFAAAQCUELAAAgEoIWAABAJAQtAACASAhaAAAAkRC0AAAAIiFoAQAARELQAgAAiISgBQAAEAlBCwAAIBKCFgAAZTwz6uFE1fXrZ2fW8phwYQ93ntGn3ZxZy6Pndy5y55k3aGpmLY/WHz/enWfdCysza3mpdmzdnplHAUELAIAy3nvjHRckxpx5S2Ytj8XjZrvzzL9zemYtj3sv6+3Os2H5msxaHt2+ea47z749ezNreVmF8yggaAEAAERC0AIAAIiEoAUAABAJQQsAACASghYAAEAkBC0AAIBICFoAAACRELQAAAAiIWgBAABEQtACAACIhKAFAAAQCUELAAAgEoIWAABAJAQtAACASAhaAAAAkRC0AAAAIiFoAQAARELQAgAAiISgBQBAGZNbDkhU7T71u8xaHgOObunO0+fI5pm1PDp89mR3nnsu7plZy0v19MgHM/N57duzN1n3wsrMPAoIWgAAlLFr2w4XSKZcOSCzlsfyWc+68yyZMjezlsfMdsPdebZt2pJZy6Pfjy5351GFa3m0OuS4A3KepoygBQBABQN/dnVmrhaDf3FdZq4Wg469JjNXi75HtcjM1aLrYeckbT5xQmYeBQQtAACASAhaAAAAkRC0AAAAIiFoAQAARELQAgAAiISgBQAAEAlBCwAAIBKCFgAAQCQELQAAgEgIWgAAAJEQtAAAACIhaAEAAERC0AIAAIiEoAUAABAJQQsAACASghYAAEAkBC0AAIBICFoAAACRELQAAAAiIWgBAABEQtACAACIhKAFAAAQCUELAAAgEoIWAABAJAQtAACASAhaAAAAkRC0AAAAIiFoAQAARELQAgAAiISgBQAAEAlBCwAAIBKCFgAAQCQELQAAgEgIWgAAAJEQtAAAACIhaAEAAERC0AIAAIiEoAUAABAJQQsAACASghYAAEAkBC0AAIBICFoAAACRELQAAAAiIWgBAABEQtACAACIhKAFAAAQCUELAAAgEoIWAABAJAQtAACASAhaAAAAkRC0AAAAIiFoAQAARELQAgAAiISgBQAAEAlBCwAAIBKCFgAAQCQELQAAgEgIWgAAAJH8D71p5l9jqPZ3AAAAAElFTkSuQmCC>