(use-ai-model-snap-in-vs-code)=
# Use an AI model snap in VS Code

The AI model snaps provide an API that can be integrated with other software.
This guide explains how to integrate an AI model snap with the Visual Studio Code (VS Code) IDE.
It assumes that the snap has already been [installed and configured](install-deepseek-r1-snap.md).

## Install the Continue extension

The [Continue extension][continue-ext] enables integration of locally running models with the IDE.
Install the Continue extension by searching for **Continue - open-source AI code assistant** in the VS Code [Extensions view][vs-code-extensions].

## Configure Continue

Open the configuration file `config.yaml`.
The [config.yaml reference][continue-ref] describes the possible locations for this file.
It’s usually located at `~/.continue/config.yaml`.

You can also open it from the extension as follows:

- In the [Activity Bar][activity-bar], select the Continue logo
- Click on the **Select model** drop down menu
- In the new window, click the configuration icon (⚙️) to open `config.yaml`

Find the `models` list in the YAML file and add: 

```yaml
  - name: Deepseek-R1
    provider: openai
    model: deepseek-r1
    apiBase: http://localhost:8080/v1
    roles:
      - chat
      - edit
```

The values above are examples based on the Deepseek-R1 snap.
Update `name`, `model`, and `apiBase` to match your specific snap and its configuration.
To identify the correct `apiBase` and `model` name, refer to [this document](use-openai-api).

```{tip}
See the [Continue extension reference][continue-ref] for additional configuration options.
```

## Use the AI snap with Continue

Once the model is configured, it can be selected from the Select Model drop down at the bottom of the Continue extension's chat box.
Any requests in the chat box will be sent to the selected model.

```{tip}
Explore the [Continue documentation][continue-docs] to learn more.
```

<!-- links -->
[continue-ext]: https://marketplace.visualstudio.com/items?itemName=Continue.continue
[continue-ref]: https://docs.continue.dev/reference
[continue-docs]: https://docs.continue.dev/
[vs-code-extensions]: https://code.visualstudio.com/docs/getstarted/extensions#_install-a-vs-code-extension
[activity-bar]: https://code.visualstudio.com/docs/getstarted/userinterface#_basic-layout
