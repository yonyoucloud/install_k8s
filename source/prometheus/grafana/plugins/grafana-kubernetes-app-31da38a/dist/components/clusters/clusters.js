///<reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
System.register(['lodash', 'app/core/app_events'], function(exports_1) {
    var lodash_1, app_events_1;
    var ClustersCtrl;
    return {
        setters:[
            function (lodash_1_1) {
                lodash_1 = lodash_1_1;
            },
            function (app_events_1_1) {
                app_events_1 = app_events_1_1;
            }],
        execute: function() {
            ClustersCtrl = (function () {
                /** @ngInject */
                function ClustersCtrl($scope, $injector, backendSrv, contextSrv, $location) {
                    this.backendSrv = backendSrv;
                    this.contextSrv = contextSrv;
                    this.$location = $location;
                    var self = this;
                    this.isOrgEditor = contextSrv.hasRole('Editor') || contextSrv.hasRole('Admin');
                    document.title = 'Grafana Kubernetes App';
                    this.clusters = {};
                    this.pageReady = false;
                    this.getClusters().then(function () {
                        self.pageReady = true;
                    });
                }
                ClustersCtrl.prototype.getClusters = function () {
                    var self = this;
                    return this.backendSrv.get('/api/datasources')
                        .then(function (result) {
                        self.clusters = lodash_1.default.filter(result, { "type": "grafana-kubernetes-datasource" });
                    });
                };
                ClustersCtrl.prototype.confirmDelete = function (id) {
                    var _this = this;
                    this.backendSrv.delete('/api/datasources/' + id).then(function () {
                        _this.getClusters();
                    });
                };
                ClustersCtrl.prototype.deleteCluster = function (cluster) {
                    var _this = this;
                    app_events_1.default.emit('confirm-modal', {
                        title: 'Delete',
                        text: 'Are you sure you want to delete this data source? ' +
                            'If you need to undeploy the collectors, then do that before deleting the data source.',
                        yesText: "Delete",
                        icon: "fa-trash",
                        onConfirm: function () {
                            _this.confirmDelete(cluster.id);
                        }
                    });
                };
                ClustersCtrl.prototype.clusterInfo = function (cluster) {
                    this.$location.path("plugins/grafana-kubernetes-app/page/cluster-info").search({ "cluster": cluster.id });
                };
                ClustersCtrl.templateUrl = 'components/clusters/partials/clusters.html';
                return ClustersCtrl;
            })();
            exports_1("ClustersCtrl", ClustersCtrl);
        }
    }
});
//# sourceMappingURL=clusters.js.map