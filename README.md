# Inference Snaps

![Status: Beta](https://img.shields.io/badge/status-beta-orange)

> Install, run, and serve AI models as native Snaps — optimized for your hardware.

Inference Snaps turns complex generative model deployment into a single command.
It is built on the idea that AI runtimes should be portable, composable, and local-first.
Whether you’re experimenting with Qwen-VL, Whisper, or Llama 3, you can instantly install a pre-optimized runtime and spin up an API endpoint locally. Each model snap includes all dependencies and is tuned for your silicon — CPU, GPU, or NPU.


Every device should be able to run intelligent software without dependency hell or optimization nightmares.


## Quick start

```bash
# Install a model snap
sudo snap install qwen-vl --channel 2.5/beta

# Check status
qwen-vl status

# Chat with the LLM
qwen-vl chat
```

## What it does

- **Plug-and-play AI models**: Each snap bundles the model weights, dependencies, and optimized runtimes.
- **Hardware-aware**: Automatically selects a silicon-optimized runtime for Intel, Ampere, or NVIDIA.
- **API-first design:** Exposes a local OpenAI-compatible REST API for integration.
- **Offline-ready:** Once installed, the model runs entirely on your machine.
- **Reproducible & isolated:** Snaps guarantee consistent environments across supported Linux distros.

## Available Inference Snaps

| Model | Description | Install | Source |
|---|---|---|---|
| qwen-vl | Vision-language model | `snap install qwen-vl` | [repo](https://github.com/canonical/qwen-vl-snap)
| deepseek-r1 | Thinking model | `snap install deepseek-r1` | [repo](https://github.com/canonical/deepseek-r1-snap)

More open source models are coming soon.
[Share your feedback](https://github.com/canonical/inference-snaps/discussions) and vote for what open-source models to support next.


## Documentation

Head to https://documentation.ubuntu.com/inference-snaps

The documentation is maintained inside the [docs](./docs) directory. 


## Contributing

We welcome pull requests!
Please see [contribution guide](docs/CONTRIBUTING.md) for tooling, build instructions, and model submission guidelines.

Join the discussion in our [GitHub Discussions](https://github.com/canonical/inference-snaps/discussions), or [report issues](https://github.com/canonical/inference-snaps/issues/new/choose).


## License

See [License and Copyright](https://canonical-inference-snaps-docs.readthedocs-hosted.com/contributing/#license-and-copyright) section in our documentation.
