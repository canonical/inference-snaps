(manage-snap-services)=
# How to manage inference snap services


All inference snaps contain a server application.
The server is started as a service automatically after installation.
To conserve computing resources, it can be stopped and started on demand by the user.

To stop and disable the service:

```shell
sudo snap stop --disable <inference-snap>
```

To start the service:

```shell
sudo snap start <inference-snap>
```

To restart the service, e.g. to apply user configurations:
```shell
sudo snap restart <inference-snap>
```

## Check the logs

To debug possible issues, you can query the server logs:

```shell
sudo snap logs <inference-snap>
```
To query more lines and follow the logs, use `-n -100 -f`.

