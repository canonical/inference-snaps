(manage-snap)=
# Manage an Inference Snap

An Inference Snap generally exposes a couple of user configuration.
It also includes a service that runs the inference server.
These can be managed by the user.

## Start and stop the server

All engines contain a server application.
It is started automatically after installation.
To conserve computing resources, it can be stopped and started on demand by the user.

To stop and disable the server:

```shell
sudo snap stop --disable <inference-snap>
```

To start the server:

```shell
sudo snap start <inference-snap>
```


## Checking the logs

To debug possible issues, you can query the server logs:

```shell
sudo snap logs <inference-snap>
```
To query more lines and follow the logs, use `-n -100 -f`.

## Configure the snap

Configurations can be explored and changed using the `get` and `set` commands.

All configurations can be seen by running:

```shell
<inference-snap> get
```

Not all configurations can be changed, as they are defined by the selected engine.

### HTTP server configurations

The HTTP server's bind host and port are user changeable and can be seen under the `http` key:

```shell 
<inference-snap> get http
```

To change the HTTP port to, for example, 8999:

```shell
sudo <inference-snap> set http.port=8999
```

If you need other machines on the network to use this server, you need to change the bind host to 0.0.0.0:

```shell
sudo <inference-snap> set http.host=0.0.0.0
```

Once changed, restart the server:

```shell
sudo snap restart <inference-snap>
```

### Model layers on GPU

In some cases CUDA-based engines allow you to limit the number of model layers that are loaded onto the GPU by using the `gpu-layers` snap option.
This is useful if the GPU does not have enough vRAM.
The remaining layers will run on the CPU.
```shell
sudo <inference-snap> set gpu-layers=20
```

## Uninstall the snap

To remove an Inference Snap:

```shell
sudo snap remove <inference-snap>
```
