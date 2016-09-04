#!/usr/bin/env bash

# consul agent -config-file consul_server.conf -bind=127.0.0.1 -dev
# consul agent -config-file consul01.conf -bind=127.0.0.1 -data-dir=/tmp
protoc -I structs/ structs/*.proto  --go_out=plugins=grpc:structs
protoc -I members/ members/*.proto  --go_out=plugins=grpc:members
