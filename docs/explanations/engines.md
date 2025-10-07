# Engines

The term "Engine" in a Generative AI system is often used to describe the core that drives the primary functions of the system: creating new content. It is sometimes used interchangeably with the AI model itself, but can also refer to the broader system architecture that enables the model's functionality.

In our context, we use Engine to describe the full, end-to-end system architecture that manages the entire process of input, generation, and output. This system includes the AI model plus all the necessary supporting components to make it serve real-world applications.

This end-to-end system is what a user or developer interacts with and typically includes:
- **AI Model:** The pre-trained model (like an LLM)
- **Runtime:** Components and infrastructure necessary to efficiently perform inference

There may be other subsystems depending on the type of AI model. For example:
- **Vector embedding:** Mechanisms for LLMs to transform input (prompts) into a numerical format that the model can understand
- **Vision/audio encoder:** Transforming vision or audio into numerical representations
- **Projection:** Used for multimodal models to align features from non-text modalities with the embedded space of the model

The Engine on its own isn't a user-facing program. It can integrate with external components, such as frontend UIs, orchestration frameworks, and data ingestion pipelines to meet use case requirements.

In summary, an Engine refers to the complete system that enables AI models to be used effectively in real-world applications.
