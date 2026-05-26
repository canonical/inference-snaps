(inference-snap-logs)=
# View inference snap logs  
To view the logs of the inference snap, use the following command:

```shell
sudo snap logs <inference-snap>
```

In order to follow the logs, set the `-f` flag:  

```shell
sudo snap logs <inference-snap> -f
```

To query additional logs, set the `-n` flag with the number of lines, for example:

```shell
sudo snap logs <inference-snap> -f -n N_LINES
``` 

Set `-n all` to query all logs.