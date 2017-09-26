#docker build -t="PRI_DOCKER_HOST:5000/google_containers/heapster-grafana-amd64:v4.2.1" .

FROM      PRI_DOCKER_HOST:5000/google_containers/heapster-grafana-amd64:v4.2.0
MAINTAINER shenghua bi <net.bsh@gmail.com>

COPY etc/grafana /etc/grafana
