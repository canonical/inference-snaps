# How to switch between engines

When installed, the engine that best matches your system is selected automatically.
You can switch to another engine at any time using the steps below.

## Check available engines

```shell
<my-model> list-engines
```

The table shows available engines.
On the **COMPAT** column is shown the compatibility of each engine with your system:

* **yes** → compatible and stable
* **beta** → compatible but not stable
* **no** → not compatible

## Switch to another engine

```shell
sudo <my-model> use-engine <engine-name>
```

This downloads and installs the selected engine.

