///<reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
System.register([], function(exports_1) {
    var KubernetesConfigCtrl;
    return {
        setters:[],
        execute: function() {
            KubernetesConfigCtrl = (function () {
                /** @ngInject */
                function KubernetesConfigCtrl($scope, $injector, $q) {
                    this.$q = $q;
                    this.enabled = false;
                    this.appEditCtrl.setPostUpdateHook(this.postUpdate.bind(this));
                }
                KubernetesConfigCtrl.prototype.postUpdate = function () {
                    var _this = this;
                    if (!this.appModel.enabled) {
                        return this.$q.resolve();
                    }
                    return this.appEditCtrl.importDashboards().then(function () {
                        _this.enabled = true;
                        return {
                            url: "plugins/grafana-kubernetes-app/page/clusters",
                            message: "Kubernetes App enabled!"
                        };
                    });
                };
                KubernetesConfigCtrl.templateUrl = 'components/config/config.html';
                return KubernetesConfigCtrl;
            })();
            exports_1("KubernetesConfigCtrl", KubernetesConfigCtrl);
        }
    }
});
//# sourceMappingURL=config.js.map