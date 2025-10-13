(use-ai-model-snap-in-ide)=
# Use an AI model snap in your favorite IDE

The AI model snaps provide an API that can be integrated with other software.
This guide explains how to integrate an AI model snap with your favorite IDE.
It assumes that the snap has already been [installed and configured](install-deepseek-r1-snap.md).

## Install Continue

[Continue][continue-docs] enables integration of locally running models with the IDE.

````{tabs}

```{group-tab} VS Code
Install the Continue extension by searching for **Continue - open-source AI code assistant** in the VS Code [Extensions view][vs-code-extensions].
```

```{group-tab} JetBrains IDEs
Install the Continue plugin by searching for **Continue** in the plugins menu on the welcome screen of the IDE, or by going to **File** > **Settings** > **Plugins** and searching for **Continue**.
```

````

## Configure Continue

Open the configuration file `config.yaml`.
The [config.yaml reference page][continue-ref] describes the possible locations of this file.
It’s usually located at `~/.continue/config.yaml`.


```{note}
If you find a config.json file instead of the YAML file, refer to this [YAML migration guide](https://docs.continue.dev/reference/yaml-migration). 
```

````{tabs}

```{group-tab} VS Code
You can also open it from the extension as follows:

1. In the [Activity Bar][activity-bar], select the Continue logo
2. Click on the **Select model** drop down menu
3. In the new window, click the configuration icon (⚙️) to open `config.yaml`
```

```{group-tab} JetBrains IDEs
You can also open it from the plugin as follows:
1. In the right [tool window bar][tool-window-bar], select the Continue logo
2. Click on the **Select model** drop down menu
3. In the new window, click the configuration icon (⚙️) to open `config.yaml`
```
````


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

The values above are examples based on the DeepSeek-R1 snap.
Update `name`, `model`, and `apiBase` to match your specific snap and its configuration.
To identify the correct `apiBase` and `model` name, check out this guide on {ref}`using the snap via its OpenAI API<use-openai-api>`.

For additional configuration options, visit the [Continue reference page][continue-ref].

## Use the AI snap with Continue

Once the model is configured, it can be selected from the Select Model drop down at the bottom of the Continue chat box.
Any requests in the chat box will be sent to the selected model.

```{tip}
Explore the [Continue documentation][continue-docs] to learn how to use it for coding, chat and more.
```

<!-- links -->
[continue-ref]: https://docs.continue.dev/reference
[continue-docs]: https://docs.continue.dev/
[vs-code-extensions]: https://code.visualstudio.com/docs/getstarted/extensions#_install-a-vs-code-extension
[activity-bar]: https://code.visualstudio.com/docs/getstarted/userinterface#_basic-layout
[tool-window-bar]: https://www.jetbrains.com/help/idea/tool-windows.html#bars_and_buttons
