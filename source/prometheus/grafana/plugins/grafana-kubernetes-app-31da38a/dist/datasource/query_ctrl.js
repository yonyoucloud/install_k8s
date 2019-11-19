System.register(['lodash', 'app/plugins/sdk'], function(exports_1) {
    var __extends = (this && this.__extends) || function (d, b) {
        for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p];
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
    var lodash_1, sdk_1;
    var K8sQueryCtrl;
    return {
        setters:[
            function (lodash_1_1) {
                lodash_1 = lodash_1_1;
            },
            function (sdk_1_1) {
                sdk_1 = sdk_1_1;
            }],
        execute: function() {
            K8sQueryCtrl = (function (_super) {
                __extends(K8sQueryCtrl, _super);
                /** @ngInject **/
                function K8sQueryCtrl($scope, $injector, templateSrv) {
                    _super.call(this, $scope, $injector);
                    this.templateSrv = templateSrv;
                    this.defaults = {};
                    lodash_1.default.defaultsDeep(this.target, this.defaults);
                    this.target.target = this.target.target || '';
                    this.target.type = this.target.type || 'timeserie';
                }
                K8sQueryCtrl.prototype.getOptions = function (query) {
                    return this.datasource.metricFindQuery('');
                };
                K8sQueryCtrl.prototype.onChangeInternal = function () {
                    this.panelCtrl.refresh(); // Asks the panel to refresh data.
                };
                K8sQueryCtrl.templateUrl = 'datasource/partials/query.editor.html';
                return K8sQueryCtrl;
            })(sdk_1.QueryCtrl);
            exports_1("K8sQueryCtrl", K8sQueryCtrl);
        }
    }
});
//# sourceMappingURL=query_ctrl.js.map