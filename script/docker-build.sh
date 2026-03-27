#!/bin/bash

set -e

v=$(cat VERSION)
docker build -t github.com/server-selfish/backend:$v .
