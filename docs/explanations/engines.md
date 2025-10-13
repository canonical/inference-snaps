(engines)=
# Engines

As a generative AI system, Inference Snaps use the term “Engine” to describe the core that drives the primary function of the snap: generating new content.
The Engine describes the full, end-to-end system architecture that manages the entire process of input, generation, and output.
This system includes the AI model weights plus all the necessary supporting components to make it serve real-world applications.

This end-to-end system is what a user or developer interacts with and typically includes:

* **AI Model Weights:** The pre-trained model (e.g. the weights of DeepSeek R1)
* **Runtime:** Execution logic and the necessary optimizations to efficiently perform matrix multiplications of the model architecture.
* **Supporting subsystems:**
  * **Language embedding:** Converts input text into numerical representations the engine can operate on.
  * **Vision or audio encoders:** Convert images or analog sound signals into numerical representations.
  * **Matrix Transformations:** Used for multimodal models to align features from non-text modalities with the embedded space of the model.

The Engine itself is not a user-facing application. Instead, it operates inside the snap, integrating with external systems, orchestration tools, and data pipelines to deliver AI capabilities to end users.
