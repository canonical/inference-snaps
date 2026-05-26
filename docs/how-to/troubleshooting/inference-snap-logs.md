(inference-snap-logs)=
# Inference Snap Logs

The inference snap provides logs that can help you troubleshoot issues and understand the behavior of the snap.

## Viewing logs
To view the logs of the inference snap, use the following command:

```shell
sudo snap logs <inference-snap>
```

In order to follow the logs in real-time, you can use the `-f` flag:

```shell
sudo snap logs -f <inference-snap>
```

If the snap has not been just restarted the logs may be truncated. To view the full logs, you can use the `-n` flag:

```shell
sudo snap logs <inference-snap> -f -n N_LINES
``` 

Where `N_LINES` is the number of lines you want to see from the logs, `all` can be used to see all logs.