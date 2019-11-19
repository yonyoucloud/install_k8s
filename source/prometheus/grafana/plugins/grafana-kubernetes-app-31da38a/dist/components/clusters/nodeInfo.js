///<reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
System.register(['moment'], function(exports_1) {
    var moment_1;
    var NodeInfoCtrl;
    function slugify(str) {
        var slug = str.replace("@", "at").replace("&", "and").replace(/[.]/g, "_").replace("/\W+/", "");
        return slug;
    }
    return {
        setters:[
            function (moment_1_1) {
                moment_1 = moment_1_1;
            }],
        execute: function() {
            NodeInfoCtrl = (function () {
                /** @ngInject */
                function NodeInfoCtrl($scope, $injector, backendSrv, datasourceSrv, $q, $location, alertSrv) {
                    var _this = this;
                    this.backendSrv = backendSrv;
                    this.datasourceSrv = datasourceSrv;
                    this.$q = $q;
                    this.$location = $location;
                    this.alertSrv = alertSrv;
                    document.title = 'Grafana Kubernetes App';
                    this.pageReady = false;
                    this.cluster = {};
                    this.clusterDS = {};
                    this.node = {};
                    if (!("cluster" in $location.search())) {
                        alertSrv.set("no cluster specified.", "no cluster specified in url", 'error');
                        return;
                    }
                    else {
                        var cluster_id = $location.search().cluster;
                        var node_name = $location.search().node;
                        this.loadDatasource(cluster_id).then(function () {
                            _this.clusterDS.getNode(node_name).then(function (node) {
                                _this.node = node;
                                _this.pageReady = true;
                            });
                        });
                    }
                }
                NodeInfoCtrl.prototype.loadDatasource = function (id) {
                    var _this = this;
                    return this.backendSrv.get('api/datasources/' + id)
                        .then(function (ds) {
                        _this.cluster = ds;
                        return _this.datasourceSrv.get(ds.name);
                    }).then(function (clusterDS) {
                        _this.clusterDS = clusterDS;
                        return clusterDS;
                    });
                };
                NodeInfoCtrl.prototype.goToNodeDashboard = function () {
                    this.$location.path("dashboard/db/k8s-node")
                        .search({
                        "var-datasource": this.cluster.jsonData.ds,
                        "var-cluster": this.cluster.name,
                        "var-node": slugify(this.node.metadata.name)
                    });
                };
                NodeInfoCtrl.prototype.conditionStatus = function (condition) {
                    var status;
                    if (condition.type === "Ready") {
                        status = condition.status === "True";
                    }
                    else {
                        status = condition.status === "False";
                    }
                    return {
                        value: status,
                        text: status ? "Ok" : "Error"
                    };
                };
                NodeInfoCtrl.prototype.isConditionOk = function (condition) {
                    return this.conditionStatus(condition).value;
                };
                NodeInfoCtrl.prototype.conditionLastTransitionTime = function (condition) {
                    return moment_1.default(condition.lastTransitionTime).format('YYYY-MM-DD HH:mm:ss');
                };
                NodeInfoCtrl.templateUrl = 'components/clusters/partials/node_info.html';
                return NodeInfoCtrl;
            })();
            exports_1("NodeInfoCtrl", NodeInfoCtrl);
        }
    }
});
//# sourceMappingURL=nodeInfo.js.map