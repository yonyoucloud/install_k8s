# Grafana App for Kubernetes

[Kubernetes](http://kubernetes.io/) is an open-source system for automating deployment, scaling, and management of containerized applications.

The Grafana Kubernetes App allows you to monitor your Kubernetes cluster's performance. It includes 4 dashboards, Cluster, Node, Pod/Container and Deployment. It allows for the automatic deployment of the required Prometheus exporters and a default scrape config to use with your in cluster Prometheus deployment. The metrics collected are high-level cluster and node stats as well as lower level pod and container stats. Use the high-level metrics to alert on and the low-level metrics to troubleshoot.

![Container Dashboard](https://github.com/grafana/kubernetes-app/blob/master/src/img/cluster-dashboard-screenshot.png?raw=true)

![Container Dashboard](https://github.com/grafana/kubernetes-app/blob/master/src/img/container-dashboard-screenshot.png?raw=true)

![Node Dashboard](https://github.com/grafana/kubernetes-app/blob/master/src/img/node-dashboard-screenshot.png?raw=true)

### Requirements

1. Currently only has support for [**Prometheus**](https://prometheus.io/docs/prometheus/latest/querying/basics/)
2. For automatic deployment of the exporters, then Kubernetes 1.6 or higher is required.
3. Grafana 5.0.0+

### Features

- The app uses Kubernetes tags to allow you to filter pod metrics. Kubernetes clusters tend to have a lot of pods and a lot of pod metrics. The Pod/Container dashboard leverages the pod tags so you can easily find the relevant pod or pods.

- Easy installation of exporters, either a one click deploy from Grafana or detailed instructions to deploy them manually them with kubectl (also quite easy!)

- Cluster level metrics that are not available in Heapster, like CPU Capacity vs CPU Usage.

### Cluster Metrics

- Pod Capacity/Usage
- Memory Capacity/Usage
- CPU Capacity/Usage
- Disk Capacity/Usage
- Overview of Nodes, Pods and Containers

### Node Metrics

- CPU
- Memory Available
- Load per CPU
- Read IOPS
- Write IOPS
- %Util
- Network Traffic/second
- Network Packets/second
- Network Errors/second

### Pod/Container Metrics

- Memory Usage
- Network Traffic
- CPU Usage
- Read IOPS
- Write IOPS

### Documentation

#### Installation

1. Use the grafana-cli tool to install kubernetes from the commandline:

```
grafana-cli plugins install kubernetes-app
```

2. Restart your Grafana server.

3. Log into your Grafana instance. Navigate to the Plugins section, found in the Grafana main menu. Click the Apps tabs in the Plugins section and select the newly installed Kubernetes app. To enable the app, click the Config tab and click on the Enable button.

#### Connecting to your Cluster

1. Go to the Cluster List page via the Kubernetes app menu.

   ![Cluster List in main menu](https://github.com/grafana/kubernetes-app/blob/master/src/img/app-menu-screenshot.png?raw=true)

2. Click the `New Cluster` button.

3. Fill in the Auth details for your cluster.

4. Choose the Prometheus datasource that will be used for reading data in the dashboards.

6. Click `Deploy`. This will deploy a Node Exporter DaemonSet, to collect health metrics for every node, and a Deployment that collects cluster metrics.

### Feedback and Questions

Please submit any issues with the app on [Github](https://github.com/grafana/kubernetes-app/issues).
