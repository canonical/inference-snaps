(available-snaps)=
# Available snaps

This page contains a list of inference snaps and the available optimizations.

## DeepSeek R1
[![deepseek-r1 snap](https://snapcraft.io/deepseek-r1/badge.svg)](https://snapcraft.io/deepseek-r1) 
[![deepseek-r1 code][gh-badge]](https://github.com/canonical/deepseek-r1-snap)

DeepSeek R1 is a reasoning Large Language Model mainly meant for chat completions.
Input and output is text-based.

This inference snap is optimized for the following hardware:

| Arch | Optimization | Description |
|--------------|--------------|-------------|
| amd64 | Intel GPU | Optimized for Intel integrated and discrete graphics |
| amd64 | Intel NPU | Intel Neural Processing Unit acceleration |
| amd64 | Intel CPU | Intel-specific CPU optimizations |
| amd64 | NVIDIA GPU | CUDA-enabled GPU acceleration |
| arm64 | Ampere {spellexception}`Altra`/One CPUs | Optimized for Ampere processors |

{{explore_optimizations}}

## Qwen VL
[![qwen-vl snap](https://snapcraft.io/qwen-vl/badge.svg)](https://snapcraft.io/qwen-vl) 
[![qwen-vl source][gh-badge]](https://github.com/canonical/qwen-vl-snap)


Qwen VL is a Vision Language Model which has the ability to process both visual and textual data.
The input can be a combination of an image and text, with the output being text-based.

The inference snap for Qwen 2.5 VL has been optimized for the following:

| Arch | Optimization | Description |
|--------------|--------------|-------------|
| amd64 | Intel GPU | Optimized for Intel integrated and discrete graphics |
| amd64 | Intel NPU | Intel Neural Processing Unit acceleration |
| amd64 | Intel CPU | Intel-specific CPU optimizations |
| amd64 | NVIDIA GPU | CUDA-enabled GPU acceleration |
| arm64 | Ampere {spellexception}`Altra`/One CPUs | Optimized for Ampere processors |

{{explore_optimizations}}

<!-- 
## Deprecated snaps

The following inference snaps are no longer maintained:
- [mistral-7b-instruct](https://github.com/canonical/mistral-7b-instruct-snap) 
-->


[gh-badge]: https://img.shields.io/badge/GitHub-repo?logo=github&color=%23104C35
