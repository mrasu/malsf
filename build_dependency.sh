# consul agent -config-file consul01.conf -bind=127.0.0.1 -data-dir=/tmp
protoc -I members/ members/*.proto  --go_out=plugins=grpc:members