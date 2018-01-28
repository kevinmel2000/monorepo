#!/bin/bash

export EXMPLENV=$2
/usr/local/bin/envoy -c /envoy/envoy_config.json &
/usr/local/bin/$1 -config_dir=/etc/$1 -log_level=debug