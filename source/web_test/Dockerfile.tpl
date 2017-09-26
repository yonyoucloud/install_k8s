# docker build -t="PRI_DOCKER_HOST:5000/esn-containers/hello:1.0" .
# CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o server ./hello.go

FROM PRI_DOCKER_HOST:5000/esn-containers/esn_base:1.0

MAINTAINER shenghua bi <net.bsh@gmail.com>

COPY server /server
