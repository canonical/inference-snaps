(switch-between-engines)=
# How to switch between engines

The engine that best matches your system is automatically selected when you install an AI model snap.
However, you can switch to another engine after installation.

## Check available engines

In the model's CLI, list all available engines:

```shell
<my-model> list-engines
```

The {spellexception}`COMPAT` column shows the compatibility of each engine with your system:

* `yes`: Compatible and stable
* `beta`: Compatible but not stable
* `no`: Not compatible

## Switch to another engine

To use a different engine with your current model:

```shell
sudo <my-model> use-engine <engine-name>
```

The selected engine will be downloaded and installed. 

## Verify selected engine
Once the installation is complete, use the status command to confirm the model is using the new engine:

```shell
<my-model> status
```

