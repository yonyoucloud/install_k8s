///<reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />

import _ from 'lodash';
import appEvents from 'app/core/app_events';
import angular from 'angular';

const nodeExporterImage='quay.io/prometheus/node-exporter:v0.15.0';
const kubestateImage = 'quay.io/coreos/kube-state-metrics:v1.1.0';

let kubestateDeployment = {
  "apiVersion": "apps/v1beta1",
  "kind": "Deployment",
  "metadata": {
    "name": "kube-state-metrics",
    "namespace": "kube-system"
  },
  "spec": {
    "selector": {
      "matchLabels": {
        "k8s-app": "kube-state-metrics",
        "grafanak8sapp": "true"
      }
    },
    "replicas": 1,
    "template": {
      "metadata": {
        "labels": {
          "k8s-app": "kube-state-metrics",
          "grafanak8sapp": "true"
        }
      },
      "spec": {
        "containers": [{
          "name": "kube-state-metrics",
          "image": kubestateImage,
          "ports": [{
            "name": "http-metrics",
            "containerPort": 8080
          }],
          "readinessProbe": {
            "httpGet": {
              "path": "/healthz",
              "port": 8080
            },
            "initialDelaySeconds": 5,
            "timeoutSeconds": 5
          }
        }]
      }
    }
  }
};

const nodeExporterDaemonSet = {
  "kind": "DaemonSet",
  "apiVersion": "extensions/v1beta1",
  "metadata": {
    "name": "node-exporter",
    "namespace": "kube-system"
  },
  "spec": {
    "selector": {
      "matchLabels": {
        "daemon": "node-exporter",
        "grafanak8sapp": "true"
      }
    },
    "template": {
      "metadata": {
        "name": "node-exporter",
        "labels": {
          "daemon": "node-exporter",
          "grafanak8sapp": "true"
        }
      },
      "spec": {
        "volumes": [
          {
            "name": "proc",
            "hostPath": {
              "path": "/proc"
            }
          },
          {
            "name": "sys",
            "hostPath": {
              "path": "/sys"
            }
          }
        ],
        "containers": [{
          "name": "node-exporter",
          "image": nodeExporterImage,
          "args": [
            "--path.procfs=/proc_host",
            "--path.sysfs=/host_sys"
          ],
          "ports": [{
            "name": "node-exporter",
            "hostPort": 9100,
            "containerPort": 9100
          }],
          "volumeMounts": [{
              "name": "sys",
              "readOnly": true,
              "mountPath": "/host_sys"
            },
            {
              "name": "proc",
              "readOnly": true,
              "mountPath": "/proc_host"
            }
          ],
          "imagePullPolicy": "IfNotPresent"
        }],
        "restartPolicy": "Always",
        "hostNetwork": true,
        "hostPID": true
      }
    }
  }
};

export class ClusterConfigCtrl {
  cluster: any;
  isOrgEditor: boolean;
  pageReady: boolean;
  prometheusDeployed: boolean;
  showHelp: boolean;
  showPrometheusExample: boolean;
  datasources: [any];
  
  static templateUrl = 'components/clusters/partials/cluster_config.html';

  /** @ngInject */
  constructor($scope, $injector, private backendSrv, private $q, private contextSrv, private $location, private $window, private alertSrv) {
    var self = this;
    this.isOrgEditor = contextSrv.hasRole('Editor') || contextSrv.hasRole('Admin');
    this.cluster = {
      type: 'grafana-kubernetes-datasource'
    };
    this.pageReady = false;
    this.prometheusDeployed = false;
    this.showHelp = false;
    this.showPrometheusExample = false;
    document.title = 'Grafana Kubernetes App';

    this.getDatasources().then(() => {
      self.pageReady = true;
    });
  }

  toggleHelp() {
    this.showHelp = !this.showHelp;
  }

  togglePrometheusExample() {
    this.showPrometheusExample = !this.showPrometheusExample;
  }

  getDatasources() {
    var self = this;
    var promises = [];
    if ("cluster" in self.$location.search()) {
      promises.push(self.getCluster(this.$location.search().cluster).then(() => {
        return self.getDeployments().then(ds => {
          _.forEach(ds.items, function (deployment) {
            if (deployment.metadata.name === "prometheus-deployment") {
              self.prometheusDeployed = true;
            }
          });
        });
      }));
    }

    promises.push(self.getPrometheusDatasources());

    return this.$q.all(promises);
  }

  getCluster(id) {
    var self = this;
    return this.backendSrv.get('/api/datasources/' + id)
      .then((ds) => {
        if (!(ds.jsonData.ds)) {
          ds.jsonData.ds = "";
        }
        self.cluster = ds;
      });
  }

  getPrometheusDatasources() {
    var self = this;
    return this.backendSrv.get('/api/datasources')
    .then((result) => {
      // self.hostedMetricsDS = _.filter(result, obj =>
      //   /grafana.net\/(graphite|prometheus)$/.test(obj.url)
      // );
      self.datasources = _.filter(result, {
        "type": "prometheus"
      });
    });
  }

  getDeployments() {
    var self = this;
    return this.backendSrv.request({
      url: 'api/datasources/proxy/' + self.cluster.id + '/apis/apps/v1beta1/namespaces/kube-system/deployments',
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    });
  }

  save() {
    return this.saveDatasource()
      .then(() => {
        return this.getDatasources();
      })
      .then(() => {
        this.alertSrv.set("Saved", "Saved and successfully connected to " + this.cluster.name, 'success', 3000);
      })
      .catch(err => {
        this.alertSrv.set("Saved", "Saved but failed to connect to " + this.cluster.name + '. Error: ' + err, 'error', 5000);
      });
  }

  savePrometheusConfigToFile() {
    let blob = new Blob([this.generatePrometheusConfig()], {
      type: "application/yaml"
    });
    this.saveToFile('prometheus.yml', blob);
  }

  saveNodeExporterDSToFile() {
    let blob = new Blob([angular.toJson(nodeExporterDaemonSet, true)], {
      type: "application/json"
    });
    this.saveToFile('grafanak8s-node-exporter-ds.json', blob);
  }

  saveKubeStateDeployToFile() {
    let blob = new Blob([angular.toJson(kubestateDeployment, true)], {
      type: "application/json"
    });
    this.saveToFile('grafanak8s-kubestate-deploy.json', blob);
  }

  saveToFile(filename, blob) {
    let blobUrl = window.URL.createObjectURL(blob);

    let element = document.createElement('a');
    element.setAttribute('href', blobUrl);
    element.setAttribute('download', filename);
    element.style.display = 'none';
    document.body.appendChild(element);
    element.click();
    document.body.removeChild(element);
  }

  deploy() {
    var question = !this.prometheusDeployed ?
      'This action will deploy Prometheus exporters to your Kubernetes cluster.' +
      'Are you sure you want to deploy?' :
      'This action will update the Prometheus exporters on your Kubernetes cluster. ' +
      'Are you sure you want to deploy?';
    appEvents.emit('confirm-modal', {
      title: 'Deploy to Kubernetes Cluster',
      text: question,
      yesText: "Deploy",
      icon: "fa-question",
      onConfirm: () => {
        this.saveAndDeploy();
      }
    });
  }

  undeploy() {
    var question = 'This action will remove the DaemonSet on your Kubernetes cluster that collects health metrics. ' +
      'Are you sure you want to remove it?';

    appEvents.emit('confirm-modal', {
      title: 'Remove Daemonset Collector',
      text: question,
      yesText: "Remove",
      icon: "fa-question",
      onConfirm: () => {
        this.undeployPrometheus();
      }
    });
  }

  saveDatasource() {
    if (this.cluster.id) {
      return this.backendSrv.put('/api/datasources/' + this.cluster.id, this.cluster);
    } else {
      return this.backendSrv.post('/api/datasources', this.cluster);
    }
  }

  saveAndDeploy() {
    return this.saveDatasource()
      .then(() => {
        return this.deployPrometheus();
      });
  }

  checkApiVersion(clusterId) {
    return this.backendSrv.request({
      url: 'api/datasources/proxy/' + clusterId + '/apis/extensions/v1beta1',
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    }).then(result => {
      if (!result.resources || result.resources.length === 0) {
        throw "This Kubernetes cluster does not support v1beta1 of the API which is needed to deploy automatically. " +
          "You can install manually using the instructions at the bottom of the page.";
      }
    });
  }

  createConfigMap(clusterId, cm) {
    return this.backendSrv.request({
      url: 'api/datasources/proxy/' + clusterId + '/api/v1/namespaces/kube-system/configmaps',
      method: 'POST',
      data: cm,
      headers: {
        'Content-Type': 'application/json'
      }
    });
  }

  createDaemonSet(clusterId, daemonSet) {
    return this.backendSrv.request({
      url: 'api/datasources/proxy/' + clusterId + '/apis/extensions/v1beta1/namespaces/kube-system/daemonsets',
      method: 'POST',
      data: daemonSet,
      headers: {
        'Content-Type': "application/json"
      }
    });
  }

  deleteDaemonSet(clusterId) {
    return this.backendSrv.request({
      url: 'api/datasources/proxy/' + clusterId + '/apis/extensions/v1beta1/namespaces/kube-system/daemonsets/node-exporter',
      method: 'DELETE',
    });
  }

  createDeployment(clusterId, deployment) {
    return this.backendSrv.request({
      url: 'api/datasources/proxy/' + clusterId + '/apis/apps/v1beta1/namespaces/kube-system/deployments',
      method: 'POST',
      data: deployment,
      headers: {
        'Content-Type': "application/json"
      }
    });
  }

  deleteDeployment(clusterId, deploymentName) {
    return this.backendSrv.request({
      url: 'api/datasources/proxy/' + clusterId + '/apis/apps/v1beta1/namespaces/kube-system/deployments/' + deploymentName,
      method: 'DELETE'
    }).then(() => {
      return this.backendSrv.request({
        url: 'api/datasources/proxy/' + clusterId +
          '/apis/extensions/v1beta1/namespaces/kube-system/replicasets?labelSelector=grafanak8sapp%3Dtrue',
        method: 'DELETE'
      });
    });
  }

  deleteConfigMap(clusterId, cmName) {
    return this.backendSrv.request({
      url: 'api/datasources/proxy/' + clusterId + '/api/v1/namespaces/kube-system/configmaps/' + cmName,
      method: 'DELETE'
    });
  }

  deletePods() {
    var self = this;
    return this.backendSrv.request({
      url: 'api/datasources/proxy/' + self.cluster.id +
        '/api/v1/namespaces/kube-system/pods?labelSelector=grafanak8sapp%3Dtrue',
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    }).then(pods => {
      if (!pods || pods.items.length === 0) {
        throw "No pods found to update.";
      }

      var promises = [];

      _.forEach(pods.items, pod => {
        promises.push(this.backendSrv.request({
          url: 'api/datasources/proxy/' + self.cluster.id + '/api/v1/namespaces/kube-system/pods/' + pod.metadata.name,
          method: 'DELETE',
        }));
      });

      return this.$q.all(promises);
    });
  }

  cancel() {
    this.$window.history.back();
  }

  deployPrometheus() {
    let self = this;
    if (!this.cluster || !this.cluster.id) {
      this.alertSrv.set("Error", "Could not connect to cluster.", 'error');
      return;
    }
    return this.checkApiVersion(self.cluster.id)
      .then(() => {
        return this.createDeployment(self.cluster.id, kubestateDeployment);
      })
      .catch(err => {
        this.alertSrv.set("Error", err, 'error');
      })
      .then(() => {
        return this.createDaemonSet(self.cluster.id, nodeExporterDaemonSet);
      })
      .catch(err => {
        this.alertSrv.set("Error", err, 'error');
      })
      .then(() => {
        this.prometheusDeployed = true;
        this.alertSrv.set("Deployed", "Prometheus and exporters have been deployed to " + self.cluster.name, 'success', 5000);
      });
  }

  undeployPrometheus() {
    var self = this;
    return this.checkApiVersion(self.cluster.id)
      .then(() => {
        return this.deleteDeployment(self.cluster.id, 'kube-state-metrics');
      })
      .catch(err => {
        this.alertSrv.set("Error", err, 'error');
      })
      .then(() => {
        return this.deleteDaemonSet(self.cluster.id);
      })
      .catch(err => {
        this.alertSrv.set("Error", err, 'error');
      })
      .then(() => {
        return this.deletePods();
      })
      .catch(err => {
        this.alertSrv.set("Error", err, 'error');
      })
      .then(() => {
        this.prometheusDeployed = false;
        this.alertSrv.set("Grafana K8s removed", "Prometheus and exporters removed from " + self.cluster.name, 'success', 5000);
      });
  }

  generatePrometheusConfig() {
    return `scrape_configs:
- job_name: \'kubernetes-kubelet\'
  scheme: https
  tls_config:
    ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
    insecure_skip_verify: true
  bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
  kubernetes_sd_configs:
  - role: node
  relabel_configs:
  - action: labelmap
    regex: __meta_kubernetes_node_label_(.+)
  - target_label: __address__
    replacement: kubernetes.default.svc:443
  - source_labels: [__meta_kubernetes_node_name]
    regex: (.+)
    target_label: __metrics_path__
    replacement: /api/v1/nodes/\${1}/proxy/metrics
- job_name: \'kubernetes-cadvisor\'
  scheme: https
  tls_config:
    ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
    insecure_skip_verify: true
  bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
  kubernetes_sd_configs:
  - role: node
  relabel_configs:
  - action: labelmap
    regex: __meta_kubernetes_node_label_(.+)
  - target_label: __address__
    replacement: kubernetes.default.svc:443
  - source_labels: [__meta_kubernetes_node_name]
    regex: (.+)
    target_label: __metrics_path__
    replacement: /api/v1/nodes/\${1}/proxy/metrics/cadvisor
- job_name: \'kubernetes-kube-state\'
  kubernetes_sd_configs:
  - role: pod
  relabel_configs:
  - action: labelmap
    regex: __meta_kubernetes_pod_label_(.+)
  - source_labels: [__meta_kubernetes_namespace]
    action: replace
    target_label: kubernetes_namespace
  - source_labels: [__meta_kubernetes_pod_name]
    action: replace
    target_label: kubernetes_pod_name
  - source_labels: [__meta_kubernetes_pod_label_grafanak8sapp]
    regex: .*true.*
    action: keep
  - source_labels: ['__meta_kubernetes_pod_label_daemon', '__meta_kubernetes_pod_node_name']
    regex: 'node-exporter;(.*)'
    action: replace
    target_label: nodename`;
  }

  generatePrometheusConfigMap() {
    return {
      "apiVersion": "v1",
      "kind": "ConfigMap",
      "metadata": {
        "name": "prometheus-configmap"
      },
      "data": {
        "prometheus.yml": this.generatePrometheusConfig()
      }
    };
  }
}
