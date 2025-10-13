# Use an LLM snap

This tutorial will guide you through the steps to install and use a Large Language Model packaged in one of our Inference Snap.
We will specifically be using Qwen 2.5 VL, which is a multi-modal LLM.
These models are also known as Vision Language Models (VLM), as they can ingest and interpret images, along with text.

## Installing the snap

We install Qwen VL from the `2.5` track, which indicates the version of the model.
The snap is currently released as `beta` stability, so we will be using that as the channel.

```
sudo snap install qwen-vl --channel 2.5/beta
```

During installation the snap detects information about your system like installed hardware and available memory.
Based on this an appropriate engine is chosen.
