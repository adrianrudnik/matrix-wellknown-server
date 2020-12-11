#!/bin/bash -ex

docker build -t matrix-schema-server .
docker run --rm -it -p 8080:8080 -v "${PWD}/example:/var/schema" matrix-schema-server
