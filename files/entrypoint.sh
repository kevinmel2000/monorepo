#!/bin/bash

/usr/local/bin/envoy -c /envoy/envoy_config.json &
/usr/local/bin/$1