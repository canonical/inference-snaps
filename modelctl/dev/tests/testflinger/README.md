# Test

This snap can be tested using Testflinger.
The snap and components are downloaded from a specified channel in the Snap Store.

Modify the exported environment variables in `run.sh` to match your test.

* SNAP_NAME: The name of the snap to be tested.
* SNAP_CHANNEL: The channel from where to install the snap.
* JOB_QUEUE: Set to the testflinger job queue that should be used.
* SELECT_ENGINE: Optional. Set to the name of the engine that will be manually selected after installation. The auto selected engine will be used if this is unset.
* EXPECTED_ENGINE: The selected engine will be checked before running the server to see if it matches this value. The test will fail if the names do not match.
* EXPECTED_TPS: A minimum number of tokens per second this engine should perform at on this specific job queue. The test will fail if the measured TPS value is below this threshold.
* INSTALL_NVIDIA_DRIVER_VERSION: Optional. Set it to install the nvidia driver, format e.g. 550.  

After setting the environment variables, you can execute `run.sh`.
Example:
```
export SNAP_NAME=deepseek-r1
export SNAP_CHANNEL=latest/edge
export JOB_QUEUE=201909-27366
export EXPECTED_ENGINE=cpu-avx512
export EXPECTED_TPS=3
./run.sh
```
