///<reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
System.register(['moment'], function(exports_1) {
    var moment_1;
    var PodInfoCtrl;
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
            PodInfoCtrl = (function () {
                /** @ngInject */
                function PodInfoCtrl($scope, $injector, backendSrv, datasourceSrv, $q, $location, alertSrv) {
                    var _this = this;
                    this.backendSrv = backendSrv;
                    this.datasourceSrv = datasourceSrv;
                    this.$q = $q;
                    this.$location = $location;
                    this.alertSrv = alertSrv;
                    document.title = 'Grafana Kubernetes App';
                    this.pageReady = false;
                    this.pod = {};
                    if (!("cluster" in $location.search())) {
                        alertSrv.set("no cluster specified.", "no cluster specified in url", 'error');
                        return;
                    }
                    else {
                        this.cluster_id = $location.search().cluster;
                        var pod_name = $location.search().pod;
                        this.loadDatasource(this.cluster_id).then(function () {
                            _this.clusterDS.getPod(pod_name).then(function (pod) {
                                _this.pod = pod;
                                _this.pageReady = true;
                            });
                        });
                    }
                }
                PodInfoCtrl.prototype.loadDatasource = function (id) {
                    var _this = this;
                    return this.backendSrv.get('api/datasources/' + id)
                        .then(function (ds) {
                        _this.datasource = ds.jsonData.ds;
                        return _this.datasourceSrv.get(ds.name);
                    }).then(function (clusterDS) {
                        _this.clusterDS = clusterDS;
                        return clusterDS;
                    });
                };
                PodInfoCtrl.prototype.conditionStatus = function (condition) {
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
                PodInfoCtrl.prototype.goToPodDashboard = function (pod) {
                    this.$location.path("dashboard/db/k8s-container")
                        .search({
                        "var-datasource": this.datasource,
                        "var-cluster": this.clusterDS.name,
                        "var-node": slugify(pod.spec.nodeName),
                        "var-namespace": pod.metadata.namespace,
                        "var-pod": pod.metadata.name
                    });
                };
                PodInfoCtrl.prototype.isConditionOk = function (condition) {
                    return this.conditionStatus(condition).value;
                };
                PodInfoCtrl.prototype.formatTime = function (time) {
                    return moment_1.default(time).format('YYYY-MM-DD HH:mm:ss');
                };
                PodInfoCtrl.templateUrl = 'components/clusters/partials/pod_info.html';
                return PodInfoCtrl;
            })();
            exports_1("PodInfoCtrl", PodInfoCtrl);
        }
    }
});
//# sourceMappingURL=podInfo.js.map