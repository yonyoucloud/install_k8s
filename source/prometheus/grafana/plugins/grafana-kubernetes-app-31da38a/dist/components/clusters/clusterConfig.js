///<reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
System.register(['lodash', 'app/core/app_events', 'angular'], function(exports_1) {
    var lodash_1, app_events_1, angular_1;
    var nodeExporterImage, kubestateImage, kubestateDeployment, nodeExporterDaemonSet, ClusterConfigCtrl;
    return {
        setters:[
            function (lodash_1_1) {
                lodash_1 = lodash_1_1;
            },
            function (app_events_1_1) {
                app_events_1 = app_events_1_1;
            },
            function (angular_1_1) {
                angular_1 = angular_1_1;
            }],
        execute: function() {
            nodeExporterImage = 'quay.io/prometheus/node-exporter:v0.15.0';
            kubestateImage = 'quay.io/coreos/kube-state-metrics:v1.1.0';
            kubestateDeployment = {
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
            nodeExporterDaemonSet = {
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
            ClusterConfigCtrl = (function () {
                /** @ngInject */
                function ClusterConfigCtrl($scope, $injector, backendSrv, $q, contextSrv, $location, $window, alertSrv) {
                    this.backendSrv = backendSrv;
                    this.$q = $q;
                    this.contextSrv = contextSrv;
                    this.$location = $location;
                    this.$window = $window;
                    this.alertSrv = alertSrv;
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
                    this.getDatasources().then(function () {
                        self.pageReady = true;
                    });
                }
                ClusterConfigCtrl.prototype.toggleHelp = function () {
                    this.showHelp = !this.showHelp;
                };
                ClusterConfigCtrl.prototype.togglePrometheusExample = function () {
                    this.showPrometheusExample = !this.showPrometheusExample;
                };
                ClusterConfigCtrl.prototype.getDatasources = function () {
                    var self = this;
                    var promises = [];
                    if ("cluster" in self.$location.search()) {
                        promises.push(self.getCluster(this.$location.search().cluster).then(function () {
                            return self.getDeployments().then(function (ds) {
                                lodash_1.default.forEach(ds.items, function (deployment) {
                                    if (deployment.metadata.name === "prometheus-deployment") {
                                        self.prometheusDeployed = true;
                                    }
                                });
                            });
                        }));
                    }
                    promises.push(self.getPrometheusDatasources());
                    return this.$q.all(promises);
                };
                ClusterConfigCtrl.prototype.getCluster = function (id) {
                    var self = this;
                    return this.backendSrv.get('/api/datasources/' + id)
                        .then(function (ds) {
                        if (!(ds.jsonData.ds)) {
                            ds.jsonData.ds = "";
                        }
                        self.cluster = ds;
                    });
                };
                ClusterConfigCtrl.prototype.getPrometheusDatasources = function () {
                    var self = this;
                    return this.backendSrv.get('/api/datasources')
                        .then(function (result) {
                        // self.hostedMetricsDS = _.filter(result, obj =>
                        //   /grafana.net\/(graphite|prometheus)$/.test(obj.url)
                        // );
                        self.datasources = lodash_1.default.filter(result, {
                            "type": "prometheus"
                        });
                    });
                };
                ClusterConfigCtrl.prototype.getDeployments = function () {
                    var self = this;
                    return this.backendSrv.request({
                        url: 'api/datasources/proxy/' + self.cluster.id + '/apis/apps/v1beta1/namespaces/kube-system/deployments',
                        method: 'GET',
                        headers: {
                            'Content-Type': 'application/json'
                        }
                    });
                };
                ClusterConfigCtrl.prototype.save = function () {
                    var _this = this;
                    return this.saveDatasource()
                        .then(function () {
                        return _this.getDatasources();
                    })
                        .then(function () {
                        _this.alertSrv.set("Saved", "Saved and successfully connected to " + _this.cluster.name, 'success', 3000);
                    })
                        .catch(function (err) {
                        _this.alertSrv.set("Saved", "Saved but failed to connect to " + _this.cluster.name + '. Error: ' + err, 'error', 5000);
                    });
                };
                ClusterConfigCtrl.prototype.savePrometheusConfigToFile = function () {
                    var blob = new Blob([this.generatePrometheusConfig()], {
                        type: "application/yaml"
                    });
                    this.saveToFile('prometheus.yml', blob);
                };
                ClusterConfigCtrl.prototype.saveNodeExporterDSToFile = function () {
                    var blob = new Blob([angular_1.default.toJson(nodeExporterDaemonSet, true)], {
                        type: "application/json"
                    });
                    this.saveToFile('grafanak8s-node-exporter-ds.json', blob);
                };
                ClusterConfigCtrl.prototype.saveKubeStateDeployToFile = function () {
                    var blob = new Blob([angular_1.default.toJson(kubestateDeployment, true)], {
                        type: "application/json"
                    });
                    this.saveToFile('grafanak8s-kubestate-deploy.json', blob);
                };
                ClusterConfigCtrl.prototype.saveToFile = function (filename, blob) {
                    var blobUrl = window.URL.createObjectURL(blob);
                    var element = document.createElement('a');
                    element.setAttribute('href', blobUrl);
                    element.setAttribute('download', filename);
                    element.style.display = 'none';
                    document.body.appendChild(element);
                    element.click();
                    document.body.removeChild(element);
                };
                ClusterConfigCtrl.prototype.deploy = function () {
                    var _this = this;
                    var question = !this.prometheusDeployed ?
                        'This action will deploy Prometheus exporters to your Kubernetes cluster.' +
                            'Are you sure you want to deploy?' :
                        'This action will update the Prometheus exporters on your Kubernetes cluster. ' +
                            'Are you sure you want to deploy?';
                    app_events_1.default.emit('confirm-modal', {
                        title: 'Deploy to Kubernetes Cluster',
                        text: question,
                        yesText: "Deploy",
                        icon: "fa-question",
                        onConfirm: function () {
                            _this.saveAndDeploy();
                        }
                    });
                };
                ClusterConfigCtrl.prototype.undeploy = function () {
                    var _this = this;
                    var question = 'This action will remove the DaemonSet on your Kubernetes cluster that collects health metrics. ' +
                        'Are you sure you want to remove it?';
                    app_events_1.default.emit('confirm-modal', {
                        title: 'Remove Daemonset Collector',
                        text: question,
                        yesText: "Remove",
                        icon: "fa-question",
                        onConfirm: function () {
                            _this.undeployPrometheus();
                        }
                    });
                };
                ClusterConfigCtrl.prototype.saveDatasource = function () {
                    if (this.cluster.id) {
                        return this.backendSrv.put('/api/datasources/' + this.cluster.id, this.cluster);
                    }
                    else {
                        return this.backendSrv.post('/api/datasources', this.cluster);
                    }
                };
                ClusterConfigCtrl.prototype.saveAndDeploy = function () {
                    var _this = this;
                    return this.saveDatasource()
                        .then(function () {
                        return _this.deployPrometheus();
                    });
                };
                ClusterConfigCtrl.prototype.checkApiVersion = function (clusterId) {
                    return this.backendSrv.request({
                        url: 'api/datasources/proxy/' + clusterId + '/apis/extensions/v1beta1',
                        method: 'GET',
                        headers: {
                            'Content-Type': 'application/json'
                        }
                    }).then(function (result) {
                        if (!result.resources || result.resources.length === 0) {
                            throw "This Kubernetes cluster does not support v1beta1 of the API which is needed to deploy automatically. " +
                                "You can install manually using the instructions at the bottom of the page.";
                        }
                    });
                };
                ClusterConfigCtrl.prototype.createConfigMap = function (clusterId, cm) {
                    return this.backendSrv.request({
                        url: 'api/datasources/proxy/' + clusterId + '/api/v1/namespaces/kube-system/configmaps',
                        method: 'POST',
                        data: cm,
                        headers: {
                            'Content-Type': 'application/json'
                        }
                    });
                };
                ClusterConfigCtrl.prototype.createDaemonSet = function (clusterId, daemonSet) {
                    return this.backendSrv.request({
                        url: 'api/datasources/proxy/' + clusterId + '/apis/extensions/v1beta1/namespaces/kube-system/daemonsets',
                        method: 'POST',
                        data: daemonSet,
                        headers: {
                            'Content-Type': "application/json"
                        }
                    });
                };
                ClusterConfigCtrl.prototype.deleteDaemonSet = function (clusterId) {
                    return this.backendSrv.request({
                        url: 'api/datasources/proxy/' + clusterId + '/apis/extensions/v1beta1/namespaces/kube-system/daemonsets/node-exporter',
                        method: 'DELETE',
                    });
                };
                ClusterConfigCtrl.prototype.createDeployment = function (clusterId, deployment) {
                    return this.backendSrv.request({
                        url: 'api/datasources/proxy/' + clusterId + '/apis/apps/v1beta1/namespaces/kube-system/deployments',
                        method: 'POST',
                        data: deployment,
                        headers: {
                            'Content-Type': "application/json"
                        }
                    });
                };
                ClusterConfigCtrl.prototype.deleteDeployment = function (clusterId, deploymentName) {
                    var _this = this;
                    return this.backendSrv.request({
                        url: 'api/datasources/proxy/' + clusterId + '/apis/apps/v1beta1/namespaces/kube-system/deployments/' + deploymentName,
                        method: 'DELETE'
                    }).then(function () {
                        return _this.backendSrv.request({
                            url: 'api/datasources/proxy/' + clusterId +
                                '/apis/extensions/v1beta1/namespaces/kube-system/replicasets?labelSelector=grafanak8sapp%3Dtrue',
                            method: 'DELETE'
                        });
                    });
                };
                ClusterConfigCtrl.prototype.deleteConfigMap = function (clusterId, cmName) {
                    return this.backendSrv.request({
                        url: 'api/datasources/proxy/' + clusterId + '/api/v1/namespaces/kube-system/configmaps/' + cmName,
                        method: 'DELETE'
                    });
                };
                ClusterConfigCtrl.prototype.deletePods = function () {
                    var _this = this;
                    var self = this;
                    return this.backendSrv.request({
                        url: 'api/datasources/proxy/' + self.cluster.id +
                            '/api/v1/namespaces/kube-system/pods?labelSelector=grafanak8sapp%3Dtrue',
                        method: 'GET',
                        headers: {
                            'Content-Type': 'application/json'
                        }
                    }).then(function (pods) {
                        if (!pods || pods.items.length === 0) {
                            throw "No pods found to update.";
                        }
                        var promises = [];
                        lodash_1.default.forEach(pods.items, function (pod) {
                            promises.push(_this.backendSrv.request({
                                url: 'api/datasources/proxy/' + self.cluster.id + '/api/v1/namespaces/kube-system/pods/' + pod.metadata.name,
                                method: 'DELETE',
                            }));
                        });
                        return _this.$q.all(promises);
                    });
                };
                ClusterConfigCtrl.prototype.cancel = function () {
                    this.$window.history.back();
                };
                ClusterConfigCtrl.prototype.deployPrometheus = function () {
                    var _this = this;
                    var self = this;
                    if (!this.cluster || !this.cluster.id) {
                        this.alertSrv.set("Error", "Could not connect to cluster.", 'error');
                        return;
                    }
                    return this.checkApiVersion(self.cluster.id)
                        .then(function () {
                        return _this.createDeployment(self.cluster.id, kubestateDeployment);
                    })
                        .catch(function (err) {
                        _this.alertSrv.set("Error", err, 'error');
                    })
                        .then(function () {
                        return _this.createDaemonSet(self.cluster.id, nodeExporterDaemonSet);
                    })
                        .catch(function (err) {
                        _this.alertSrv.set("Error", err, 'error');
                    })
                        .then(function () {
                        _this.prometheusDeployed = true;
                        _this.alertSrv.set("Deployed", "Prometheus and exporters have been deployed to " + self.cluster.name, 'success', 5000);
                    });
                };
                ClusterConfigCtrl.prototype.undeployPrometheus = function () {
                    var _this = this;
                    var self = this;
                    return this.checkApiVersion(self.cluster.id)
                        .then(function () {
                        return _this.deleteDeployment(self.cluster.id, 'kube-state-metrics');
                    })
                        .catch(function (err) {
                        _this.alertSrv.set("Error", err, 'error');
                    })
                        .then(function () {
                        return _this.deleteDaemonSet(self.cluster.id);
                    })
                        .catch(function (err) {
                        _this.alertSrv.set("Error", err, 'error');
                    })
                        .then(function () {
                        return _this.deletePods();
                    })
                        .catch(function (err) {
                        _this.alertSrv.set("Error", err, 'error');
                    })
                        .then(function () {
                        _this.prometheusDeployed = false;
                        _this.alertSrv.set("Grafana K8s removed", "Prometheus and exporters removed from " + self.cluster.name, 'success', 5000);
                    });
                };
                ClusterConfigCtrl.prototype.generatePrometheusConfig = function () {
                    return "scrape_configs:\n- job_name: 'kubernetes-kubelet'\n  scheme: https\n  tls_config:\n    ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt\n    insecure_skip_verify: true\n  bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token\n  kubernetes_sd_configs:\n  - role: node\n  relabel_configs:\n  - action: labelmap\n    regex: __meta_kubernetes_node_label_(.+)\n  - target_label: __address__\n    replacement: kubernetes.default.svc:443\n  - source_labels: [__meta_kubernetes_node_name]\n    regex: (.+)\n    target_label: __metrics_path__\n    replacement: /api/v1/nodes/${1}/proxy/metrics\n- job_name: 'kubernetes-cadvisor'\n  scheme: https\n  tls_config:\n    ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt\n    insecure_skip_verify: true\n  bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token\n  kubernetes_sd_configs:\n  - role: node\n  relabel_configs:\n  - action: labelmap\n    regex: __meta_kubernetes_node_label_(.+)\n  - target_label: __address__\n    replacement: kubernetes.default.svc:443\n  - source_labels: [__meta_kubernetes_node_name]\n    regex: (.+)\n    target_label: __metrics_path__\n    replacement: /api/v1/nodes/${1}/proxy/metrics/cadvisor\n- job_name: 'kubernetes-kube-state'\n  kubernetes_sd_configs:\n  - role: pod\n  relabel_configs:\n  - action: labelmap\n    regex: __meta_kubernetes_pod_label_(.+)\n  - source_labels: [__meta_kubernetes_namespace]\n    action: replace\n    target_label: kubernetes_namespace\n  - source_labels: [__meta_kubernetes_pod_name]\n    action: replace\n    target_label: kubernetes_pod_name\n  - source_labels: [__meta_kubernetes_pod_label_grafanak8sapp]\n    regex: .*true.*\n    action: keep\n  - source_labels: ['__meta_kubernetes_pod_label_daemon', '__meta_kubernetes_pod_node_name']\n    regex: 'node-exporter;(.*)'\n    action: replace\n    target_label: nodename";
                };
                ClusterConfigCtrl.prototype.generatePrometheusConfigMap = function () {
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
                };
                ClusterConfigCtrl.templateUrl = 'components/clusters/partials/cluster_config.html';
                return ClusterConfigCtrl;
            })();
            exports_1("ClusterConfigCtrl", ClusterConfigCtrl);
        }
    }
});
//# sourceMappingURL=clusterConfig.js.map