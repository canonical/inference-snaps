(network-ports)=

# Network ports registry

Inference snaps bind on one or multiple TCP ports to expose their APIs. If a port is already in use by another service on your system, the inference snap service will fail to start. In such cases, you can {ref}`configure the snap <configure-inference-snaps>` to listen on a different port.

The following table lists the default port assignments for each inference snap. The port range starts from the unassigned IANA registration subset 8323-8350 and continues into higher numbers. Each snap is assigned an even-numbered port, with the following odd number reserved for future use.

| Snap | Reserved ports |
|------|-------|
| [DeepSeek R1](https://snapcraft.io/deepseek-r1) | 8324, 8325 |
| [Qwen 2.5 VL](https://snapcraft.io/qwen-vl) | 8326, 8327 |
| [Gemma3](https://snapcraft.io/gemma3) | 8328, 8329 |
| [Nemotron 3 Nano](https://snapcraft.io/nemotron-3-nano) | 8330, 8331 |
| [Nemotron 3 Super](https://snapcraft.io/nemotron-3-super) | 8332, 8333 |
| [Nemotron 3 Nano Omni](https://snapcraft.io/nemotron-3-nano-omni) | 8334, 8335 |
| [Gemma4](https://snapcraft.io/gemma4) | 8336, 8337 |
| [Qwen 3](https://snapcraft.io/qwen3) | 8338, 8339 |
| [Qwen 3 Coder](https://snapcraft.io/qwen3-coder) | 8340, 8341 |
| [Qwen 3.6](https://snapcraft.io/qwen3-6) | 8342, 8343 |
| [SmolLM2](https://snapcraft.io/smollm2) | 8344, 8345 |
| Qwen 3.5 | 8352, 8353 |
| [GLM 4.7 Flash](https://snapcraft.io/glm-4-7-flash) | 8354, 8355 |
