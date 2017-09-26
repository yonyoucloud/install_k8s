#!/bin/sh

mkdir -p /data/docker/private-registry/{storage,registry}
sha256=`docker load -i /tmp/registry.tar | grep Loaded | awk '{print $4}' | awk -F ':' '{print $2}'`
docker tag $sha256 registry:latest
docker run \
    -d \
    --restart=always \
    --name registry \
    -e STORAGE_PATH=/registry-storage \
    -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/HOST.crt \
    -e REGISTRY_HTTP_TLS_KEY=/certs/HOST.key \
    -v /etc/certs:/certs \
    -v /data/docker/private-registry/storage:/registry-storage \
    -v /data/docker/private-registry/registry:/var/lib/registry \
    -u root \
    -p 5000:5000 \
    registry:latest
