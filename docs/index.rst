Inference Snaps
==========================

.. In a single sentence, say what the product is.
.. In a paragraph of one to three short sentences, describe what the product does.
.. In a paragraph of one to three short sentences, describe what need the product meets.
.. Whom is the product useful for?

Inference snaps are generative AI models that are packaged and optimized for efficient performance on local hardware.
They automatically detect your host machine's hardware and install runtime and model weight optimizations that best match its capabilities. 

Inference snaps provide an automatic, unified detection and installation mechanism.
This results in model and runtime choices that make effective use of your CPU, GPU, or NPU.

With inference snaps, you get fast, local AI capabilities for applications such as IDEs, chat bots, and media editing tools.
They also provide application developers with a standardized and reliable local API to integrate AI features into their software without the need for complex hardware-specific tuning.


In this documentation
---------------------
=================================        =========================================================

=================================        =========================================================
**Installing and configuring:**          :ref:`Install Qwen VL snap<qwen-vl-tutorial>` • :ref:`Switch engines <switch-between-engines>` • :ref:`Use OpenAI API<use-openai-api>`
**Hardware drivers:**                    :ref:`Intel GPU <install-intel-gpu-drivers>` • :ref:`Intel NPU <install-intel-npu-drivers>` • :ref:`Nvidia GPU <install-nvidia-gpu-drivers>`
**Building blocks:**                     :ref:`Snap architecture <architecture>` • :ref:`Engines <engines>` • :ref:`The command line interface<models-cli>`
=================================        =========================================================

How this documentation is organized
-----------------------------------
- The :ref:`tutorial<tutorials-index>` is a hands-on introduction to inference snaps for new users.
- :ref:`How-to guides<how-to-guides-index>` cover common tasks and key operations. These assume you have basic familiarity with inference snaps. 
- The :ref:`reference<reference-index>` section provides specifications and architectures.
- :ref:`Explanations<explanations-index>` contain detailed discussions on key concepts in the development and use of inference snaps.
   
Project and community
---------------------

The Inference Snaps project is a member of the Ubuntu family.
It is open source and warmly welcomes community projects, contributions, suggestions, fixes and constructive feedback.

- `Code of conduct <https://ubuntu.com/community/docs/ethos/code-of-conduct>`__
- :ref:`Contribute to our code <contribute-to-our-code>` 
- :ref:`Contribute to our documentation <contribute-to-our-docs>`
- :ref:`Release notes<release-notes>`


.. toctree::
   :hidden:
   :maxdepth: 1

   Tutorials </tutorial/index>
   How-to guides </how-to/index>
   Reference </reference/index>
   Explanations </explanations/index>
   contributing
   release-notes
