///<reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
System.register(['lodash', 'jquery'], function(exports_1) {
    var lodash_1, jquery_1;
    var ClusterInfoCtrl;
    function slugify(str) {
        var slug = str.replace("@", "at").replace("&", "and").replace(/[.]/g, "_").replace("/\W+/", "");
        return slug;
    }
    function getComponentHealth(component) {
        var health = "unhealthy";
        var message = '';
        lodash_1.default.forEach(component.conditions, function (condition) {
            if (condition.type === "Healthy" &&
                condition.status === "True") {
                health = "ok";
            }
            else {
                message = condition.message;
            }
        });
        return getHealthState(health, message);
    }
    function getNodeHealth(node) {
        var health = "unhealthy";
        var message = '';
        lodash_1.default.forEach(node.status.conditions, function (condition) {
            if (condition.type === "Ready" &&
                condition.status === "True") {
                health = "ok";
            }
            else {
                message = condition.message;
            }
        });
        return getHealthState(health, message);
    }
    function getHealthState(health, message) {
        switch (health) {
            case 'ok': {
                return {
                    text: 'OK',
                    iconClass: 'icon-gf icon-gf-online',
                    stateClass: 'alert-state-ok',
                    message: ''
                };
            }
            case 'unhealthy': {
                return {
                    text: 'UNHEALTHY',
                    iconClass: 'icon-gf icon-gf-critical',
                    stateClass: 'alert-state-critical',
                    message: message || ''
                };
            }
            case 'warning': {
                return {
                    text: 'warning',
                    iconClass: "icon-gf icon-gf-critical",
                    stateClass: 'alert-state-warning',
                    message: message || ''
                };
            }
        }
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
            ClusterInfoCtrl = (function () {
                /** @ngInject */
                function ClusterInfoCtrl($scope, $injector, backendSrv, datasourceSrv, $q, $location, alertSrv) {
                    var _this = this;
                    this.backendSrv = backendSrv;
                    this.datasourceSrv = datasourceSrv;
                    this.$q = $q;
                    this.$location = $location;
                    this.alertSrv = alertSrv;
                    this.$q = $q;
                    document.title = 'Grafana Kubernetes App';
                    this.pageReady = false;
                    this.cluster = {};
                    this.componentStatuses = [];
                    this.namespaces = [];
                    this.namespace = "";
                    this.nodes = [];
                    if (!("cluster" in $location.search())) {
                        alertSrv.set("no cluster specified.", "no cluster specified in url", 'error');
                        return;
                    }
                    this.getCluster($location.search().cluster)
                        .then(function (clusterDS) {
                        _this.clusterDS = clusterDS;
                        _this.pageReady = true;
                        _this.getClusterInfo();
                    });
                }
                ClusterInfoCtrl.prototype.getCluster = function (id) {
                    var _this = this;
                    return this.backendSrv.get('api/datasources/' + id).then(function (ds) {
                        _this.cluster = ds;
                        return _this.datasourceSrv.get(ds.name);
                    });
                };
                ClusterInfoCtrl.prototype.getClusterInfo = function () {
                    var _this = this;
                    this.clusterDS.getComponentStatuses().then(function (stats) {
                        _this.componentStatuses = lodash_1.default.map(stats, function (stat) {
                            stat.healthState = getComponentHealth(stat);
                            return stat;
                        });
                    });
                    this.clusterDS.getNamespaces().then(function (namespaces) {
                        _this.namespaces = namespaces;
                    });
                    this.clusterDS.getNodes().then(function (nodes) {
                        _this.nodes = lodash_1.default.map(nodes, function (node) {
                            node.healthState = getNodeHealth(node);
                            return node;
                        });
                    });
                };
                ClusterInfoCtrl.prototype.goToClusterDashboard = function () {
                    this.$location.path("dashboard/db/k8s-cluster")
                        .search({
                        "var-datasource": this.cluster.jsonData.ds,
                        "var-cluster": this.cluster.name
                    });
                };
                ClusterInfoCtrl.prototype.goToPodDashboard = function () {
                    this.$location.path("dashboard/db/k8s-container")
                        .search({
                        "var-datasource": this.cluster.jsonData.ds,
                        "var-cluster": this.cluster.name,
                        "var-node": 'All',
                        "var-namespace": 'All',
                        "var-pod": 'All'
                    });
                };
                ClusterInfoCtrl.prototype.goToNodeDashboard = function (node, evt) {
                    var clickTargetIsLinkOrHasLinkParents = jquery_1.default(evt.target).closest('a').length > 0;
                    if (clickTargetIsLinkOrHasLinkParents === false) {
                        this.$location.path("dashboard/db/k8s-node")
                            .search({
                            "var-datasource": this.cluster.jsonData.ds,
                            "var-cluster": this.cluster.name,
                            "var-node": node === 'All' ? 'All' : slugify(node.metadata.name)
                        });
                    }
                };
                ClusterInfoCtrl.prototype.goToWorkloads = function (ns, evt) {
                    var clickTargetIsLinkOrHasLinkParents = jquery_1.default(evt.target).closest('a').length > 0;
                    if (clickTargetIsLinkOrHasLinkParents === false) {
                        this.$location.path("plugins/grafana-kubernetes-app/page/cluster-workloads")
                            .search({
                            "cluster": this.cluster.id,
                            "namespace": slugify(ns.metadata.name)
                        });
                    }
                };
                ClusterInfoCtrl.prototype.goToNodeInfo = function (node, evt) {
                    var clickTargetIsLinkOrHasLinkParents = jquery_1.default(evt.target).closest('a').length > 0;
                    var closestElm = lodash_1.default.head(jquery_1.default(evt.target).closest('div'));
                    var clickTargetClickAttr = lodash_1.default.find(closestElm.attributes, { name: "ng-click" });
                    var clickTargetIsNodeDashboard = clickTargetClickAttr ? clickTargetClickAttr.value === "ctrl.goToNodeDashboard(node, $event)" : false;
                    if (clickTargetIsLinkOrHasLinkParents === false &&
                        clickTargetIsNodeDashboard === false) {
                        this.$location.path("plugins/grafana-kubernetes-app/page/node-info")
                            .search({
                            "cluster": this.cluster.id,
                            "node": node.metadata.name
                        });
                    }
                };
                ClusterInfoCtrl.templateUrl = 'components/clusters/partials/cluster_info.html';
                return ClusterInfoCtrl;
            })();
            exports_1("ClusterInfoCtrl", ClusterInfoCtrl);
        }
    }
});
//# sourceMappingURL=clusterInfo.js.map