(network-ports)=

# Network ports registry

Inference snaps bind on one or multiple TCP ports to expose their APIs. If a port is already in use by another service on your system, the inference snap service will fail to start. In such cases, you can {ref}`configure the snap <configure-inference-snaps>` to listen on a different port.

The following table lists the default port assignments for each inference snap. The port range starts from the unassigned IANA registration subset 8323-8350 and continues into higher numbers. Each snap is assigned an even-numbered port, with the following odd number reserved for future use.

| Snap | Reserved ports |
|------|-------|
| deepseek-r1 | 8324, 8325 |
| qwen-vl | 8326, 8327 |
| gemma3 | 8328, 8329 |
| nemotron-3-nano | 8330, 8331 |
| nemotron-3-super | 8332, 8333 |
| nemotron-3-nano-omni | 8334, 8335 |
| gemma4 | 8336, 8337 |
| qwen3 | 8338, 8339 |
| qwen3-coder | 8340, 8341 |
| qwen3.6 | 8342, 8343 |
| smollm2 | 8344, 8345 |
| qwen3.5 | 8352, 8353 |
| glm-4.7-flash | 8354, 8355 |
