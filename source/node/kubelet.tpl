KUBELET_HOSTNAME="--hostname-override=NODE_HOST"
KUBELET_PAUSE_IMAGE="--pod-infra-container-image=PRI_DOCKER_HOST:5000/google-containers/pause-amd64:3.1"
KUBELET_ARGS="--kubeconfig=/etc/kubernetes/kubelet.conf --config=/etc/kubernetes/kubelet.yaml"
