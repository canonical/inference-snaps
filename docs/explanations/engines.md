(engines)=
# Engines

Inference snaps, as generative AI systems, use the term “engine” to describe the core that drives their primary function, i.e., generating new content.
It refers to the entire system architecture that manages input, generation, and output of information.

Engines consist of the AI model weights plus every other component they need to be able to serve real-world applications. An engine will typically have:

* **AI Model Weights:** The pre-trained model (e.g. the weights of Qwen2.5 VL {spellexception}`quantized` at Q4_K)
* **Runtime:** Execution logic and the necessary optimizations to efficiently perform matrix multiplications of the model architecture.
* **Supporting subsystems:**
  * **Language embedding:** Converts input text into numerical representations the engine can use.
  * **Vision and audio encoders:** Convert images and analog sound signals into numerical representations.
  * **Matrix Transformations:** Used for multimodal models to align features from non-text modalities with the embedded space of the model.

An engine is not a user-facing application.
It is integrated with external systems, orchestration tools, and data pipelines to deliver AI capabilities to end users.
