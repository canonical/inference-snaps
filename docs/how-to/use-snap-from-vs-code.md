(use-deepseek-r1-snap-with-vs-code)=
# Use Deepseek-R1 snap with VS Code

The Deepseek-R1 snap provides an API that can be integrated with other software.
This guide explains how to integrate the Deepseek-R1 snap with the Visual Studio Code (VS Code) IDE.
It assumes that the snap has already been [installed and configured](install-deepseek-r1-snap.md).

## Install the Continue extension

The [Continue extension][continue-ext] enables integration of locally running models with the IDE.
First, [install][vs-code-extensions] the Continue extension by searching for: **Continue - open-source AI code assistant**.

## Configure Continue

Next, open the configuration file `config.yaml`.
The [config.yaml reference][continue-ref] describes the possible locations for this file.
It’s usually located at `~/.continue/config.yaml`.

You can also open it from the extension.
 - In the extensions menu: select the Continue logo
 - Click: **Select model**
 - In the new window, click the configuration icon (⚙️)

A YAML file will open. In the `models` list, add:

```yaml
  - name: Deepseek-R1
    provider: openai
    model: deepseek-r1
    apiBase: http://127.0.0.1:8080/v1
    roles:
      - chat
      - edit
```

Adjust the `apiBase` according to your Deepseek-R1 snap configuration.

```{tip}
See the [Continue extension reference][continue-ref] for additional configuration options.
```

## Use Deepseek-r1 snap with Continue

Once the models is configured, it can be selected on the Continue panel and be used.

```{tip}
Explore the [Continue documentation][continue-docs] to learn more.
```

<!-- links -->
[continue-ext]: https://marketplace.visualstudio.com/items?itemName=Continue.continue
[continue-ref]: https://docs.continue.dev/reference
[continue-docs]: https://docs.continue.dev/
[vs-code-extensions]: https://code.visualstudio.com/docs/getstarted/extensions
