(sideload-model)=
# Sideload a model
Sideloading a model allows you to test a new LLM without modifying the snap's configuration files or rebuilding the snap.

To get started, inspect the current environment variables in the snap:

```shell
sudo <inference-snap> run -- env
```

This command prints all environment variables available in the snap.
Now suppose you want to sideload a model that requires a few environment variables:

```shell
sudo <inference-snap> set env.MODEL_FILE=/path/to/my_model
sudo <inference-snap> set env.MODEL_NAME=my_model
sudo <inference-snap> set env.MMPROJ_FILE=/path/to/my_model.mmproj # if the model also requires an mmproj file
```

Restart the snap to apply the changes:

```shell
sudo snap restart <inference-snap>
```