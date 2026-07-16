# Inference Snaps

> Install, run, and serve AI models as snaps — optimized for your hardware.

Inference Snaps turn complex generative model deployment into a single command.
They are built on the idea that AI runtimes should be portable, composable, and local-first.
Whether you’re experimenting with Gemma 4, Qwen 3 Coder, or DeepSeek R1, you can instantly install a pre-optimized runtime and spin up an API endpoint locally. Each model snap includes all dependencies and is tuned for your silicon — CPU, GPU, or NPU.

Every device should be able to run intelligent software without dependency hell or optimization nightmares.


## Quick start

Install the [accelerator driver](https://documentation.ubuntu.com/inference-snaps/how-to/install-drivers/).

Then:

```bash
# Install a model snap
sudo snap install gemma4

# Check status
gemma4 status

# Chat with the LLM
gemma4 chat
```

## What it does

- **Plug-and-play AI models**: Each snap bundles a range of optimized model weights and runtimes. 
- **Hardware-aware**: Automatically selects and deploys silicon-optimized runtime and model weights for your machine.
- **API-first design:** Exposes local standard APIs for integration.
- **Offline-ready:** Once installed, the model runs entirely on your machine.
- **Reproducible & isolated:** Snaps guarantee consistent environments across supported Linux distros.

## Available inference snaps

| Model | Capabilities | Install | Source |
|---|---|---|---|
| DeepSeek R1 | text, thinking | `snap install deepseek-r1` | [repo](https://github.com/canonical/deepseek-r1-snap)
| Gemma 3 | text, thinking, vision | `snap install gemma3` | [repo](https://github.com/canonical/gemma3-snap)
| Gemma 4 | text, thinking, vision, tools | `snap install gemma4` | [repo](https://github.com/canonical/gemma4-snap)
| GLM 4.7 Flash | text, thinking, tools | `snap install glm-4-7-flash` | [repo](https://github.com/canonical/glm-4.7-flash-snap)
| Nemotron 3 Nano | text, thinking, tools | `snap install nemotron-3-nano` | [repo](https://github.com/canonical/nemotron-3-nano-snap)
| Nemotron 3 Nano Omni | text, thinking, vision, tools | `snap install nemotron-3-nano-omni` | [repo](https://github.com/canonical/nemotron-3-nano-omni-snap)
| Qwen 2.5 VL | text, vision | `snap install qwen-vl` | [repo](https://github.com/canonical/qwen-vl-snap)
| Qwen 3 | text, thinking, tools | `snap install qwen3` | [repo](https://github.com/canonical/qwen3-snap)
| Qwen 3 Coder | text, tools | `snap install qwen3-coder` | [repo](https://github.com/canonical/qwen3-coder-snap)
| Qwen 3.6 | text, thinking, vision, tools | `snap install qwen3-6` | [repo](https://github.com/canonical/qwen3.6-snap)

More open source models are coming soon.
[Share your feedback](https://github.com/canonical/inference-snaps/discussions) and vote for what open-source models to support next.


## Documentation

Head to https://documentation.ubuntu.com/inference-snaps

The documentation is maintained inside the [docs](./docs) directory. 


## Contributing

We welcome pull requests!
Please see [contribution guide](https://documentation.ubuntu.com/inference-snaps/contributing/) for tooling, build instructions, and model submission guidelines.

Join the discussion in our [GitHub Discussions](https://github.com/canonical/inference-snaps/discussions), or [report issues](https://github.com/canonical/inference-snaps/issues/new/choose).


