(switch-between-engines)=
# How to switch between engines

The engine that best matches your system is automatically selected when you install an Inference Snap.
However, you can switch to another engine after installation.

## Check available engines

In the Inference Snaps's CLI, list all available engines:

```shell
<inference-snap> list-engines
```

The {spellexception}`COMPAT` column shows the compatibility of each engine with your system:

* `yes`: Compatible and stable
* `beta`: Compatible but not stable
* `no`: Not compatible

## Switch to another engine

To use a different engine with your current Inference Snap:

```shell
sudo <inference-snap> use-engine <engine-name>
```

The selected engine will be downloaded and installed. 

## Verify selected engine
Once the installation is complete, use the status command to confirm the Inference Snap is using the new engine:

```shell
<inference-snap> status
```

