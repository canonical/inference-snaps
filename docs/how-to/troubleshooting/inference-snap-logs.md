(inference-snap-logs)=
# Inference Snap Logs

The inference snap provides logs that can help you troubleshoot issues and understand the behavior of the snap.

# View inference snap logs  
To view the logs of the inference snap, use the following command:

```shell
sudo snap logs <inference-snap>
```

In order to follow the logs, set the `-f` flag:  

```shell
sudo snap logs -f <inference-snap>
```

To query additional logs, set the `-n` flag with the number of lines, for example:

```shell
sudo snap logs <inference-snap> -f -n N_LINES
``` 

Set `-n all` to query all logs.
