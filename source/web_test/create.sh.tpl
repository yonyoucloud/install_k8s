#!/bin/sh

proName=web_test

source /etc/profile

date=$(date "+%Y%m%d.%H%M%S")
imageName=PRI_DOCKER_HOST:5000/esn-containers/$proName:$date

cp -rp yaml/pod.yaml.tpl yaml/pod.yaml
sed -i "s#PRO_IMAGE#$imageName#g" yaml/pod.yaml

docker build -t="$imageName" .
docker push $imageName

kubectl create namespace esn-system
kubectl apply -f yaml/pod.yaml
kubectl apply -f yaml/svc.yaml

echo $imageName >> ./image_version
