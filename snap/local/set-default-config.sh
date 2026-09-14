#!/bin/bash -eu

# Wipe existing package configuration
snapctl unset config.package

# Set default package configurations
modelctl set --package verbose="false"
modelctl set --package http.port="8322"
modelctl set --package http.host="127.0.0.1"
modelctl set --package webui.http.port="8323"
modelctl set --package webui.http.host="127.0.0.1"
modelctl set --package logger.http.port="8322"
modelctl set --package logger.http.host="127.0.0.1"
modelctl set --package ws.unix-socket="dummy.sock"
