///<reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
System.register(['lodash', 'jquery'], function(exports_1) {
    var lodash_1, jquery_1;
    var ClusterWorkloadsCtrl;
    function slugify(str) {
        var slug = str.replace("@", "at").replace("&", "and").replace(/[.]/g, "_").replace("/\W+/", "");
        return slug;
    }
    return {
        setters:[
            function (lodash_1_1) {
                lodash_1 = lodash_1_1;
            },
            function (jquery_1_1) {
                jquery_1 = jquery_1_1;
            }],
        execute: function() {
            ClusterWorkloadsCtrl = (function () {
                /** @ngInject */
                function ClusterWorkloadsCtrl($scope, $injector, backendSrv, datasourceSrv, $q, $location, alertSrv) {
                    var _this = this;
                    this.backendSrv = backendSrv;
                    this.datasourceSrv = datasourceSrv;
                    this.$q = $q;
                    this.$location = $location;
                    this.alertSrv = alertSrv;
                    document.title = 'Grafana Kubernetes App';
                    this.pageReady = false;
                    this.cluster = {};
                    this.namespaces = [];
                    this.namespace = "";
                    this.daemonSets = [];
                    this.replicationControllers = [];
                    this.deployments = [];
                    this.pods = [];
                    if (!("cluster" in $location.search())) {
                        alertSrv.set("no cluster specified.", "no cluster specified in url", 'error');
                        return;
                    }
                    if ("namespace" in $location.search()) {
                        this.namespace = $location.search().namespace;
                    }
                    this.getCluster($location.search().cluster)
                        .then(function (clusterDS) {
                        _this.clusterDS = clusterDS;
                        _this.pageReady = true;
                        _this.getWorkloads();
                    });
                }
                ClusterWorkloadsCtrl.prototype.getCluster = function (id) {
                    var _this = this;
                    return this.backendSrv.get('api/datasources/' + id).then(function (ds) {
                        _this.cluster = ds;
                        return _this.datasourceSrv.get(ds.name);
                    });
                };
                ClusterWorkloadsCtrl.prototype.getWorkloads = function () {
                    var _this = this;
                    var namespace = this.namespace;
                    this.clusterDS.getNamespaces().then(function (namespaces) {
                        _this.namespaces = namespaces;
                    });
                    this.clusterDS.getDaemonSets(namespace).then(function (daemonSets) {
                        _this.daemonSets = daemonSets;
                    });
                    this.clusterDS.getReplicationControllers(namespace).then(function (rc) {
                        _this.replicationControllers = rc;
                    });
                    this.clusterDS.getDeployments(namespace).then(function (deployments) {
                        _this.deployments = deployments;
                    });
                    this.clusterDS.getPods(namespace).then(function (pods) {
                        _this.pods = pods;
                    });
                };
                ClusterWorkloadsCtrl.prototype.componentHealth = function (component) {
                    var health = "unhealthy";
                    lodash_1.default.forEach(component.conditions, function (condition) {
                        if ((condition.type === "Healthy") && (condition.status === "True")) {
                            health = "healthy";
                        }
                    });
                    return health;
                };
                ClusterWorkloadsCtrl.prototype.isComponentHealthy = function (component) {
                    return this.componentHealth(component) === "healthy";
                };
                ClusterWorkloadsCtrl.prototype.goToPodDashboard = function (pod) {
                    this.$location.path("dashboard/db/k8s-container")
                        .search({
                        "var-datasource": this.cluster.jsonData.ds,
                        "var-cluster": this.cluster.name,
                        "var-node": slugify(pod.spec.nodeName),
                        "var-namespace": pod.metadata.namespace,
                        "var-pod": pod.metadata.name
                    });
                };
                ClusterWorkloadsCtrl.prototype.goToDeploymentDashboard = function (deploy) {
                    this.$location.path("dashboard/db/k8s-deployments")
                        .search({
                        "var-datasource": this.cluster.jsonData.ds,
                        "var-cluster": this.cluster.name,
                        "var-namespace": deploy.metadata.namespace,
                        "var-deployment": deploy.metadata.name
                    });
                };
                ClusterWorkloadsCtrl.prototype.goToPodInfo = function (pod, evt) {
                    var clickTargetIsLinkOrHasLinkParents = jquery_1.default(evt.target).closest('a').length > 0;
                    var closestElm = lodash_1.default.head(jquery_1.default(evt.target).closest('div'));
                    var clickTargetClickAttr = lodash_1.default.find(closestElm.attributes, { name: "ng-click" });
                    var clickTargetIsNodeDashboard = clickTargetClickAttr ? clickTargetClickAttr.value === "ctrl.goToPodDashboard(pod, $event)" : false;
                    if (clickTargetIsLinkOrHasLinkParents === false &&
                        clickTargetIsNodeDashboard === false) {
                        this.$location.path("plugins/grafana-kubernetes-app/page/pod-info")
                            .search({
                            "cluster": this.cluster.id,
                            "namespace": slugify(pod.metadata.namespace),
                            "pod": pod.metadata.name
                        });
                    }
                };
                ClusterWorkloadsCtrl.templateUrl = 'components/clusters/partials/cluster_workloads.html';
                return ClusterWorkloadsCtrl;
            })();
            exports_1("ClusterWorkloadsCtrl", ClusterWorkloadsCtrl);
        }
    }
});
//# sourceMappingURL=clusterWorkloads.js.map