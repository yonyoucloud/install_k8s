///<reference path="../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
System.register(['./datasource', './query_ctrl'], function(exports_1) {
    var datasource_1, query_ctrl_1;
    var K8sConfigCtrl;
    return {
        setters:[
            function (datasource_1_1) {
                datasource_1 = datasource_1_1;
            },
            function (query_ctrl_1_1) {
                query_ctrl_1 = query_ctrl_1_1;
            }],
        execute: function() {
            K8sConfigCtrl = (function () {
                function K8sConfigCtrl() {
                }
                K8sConfigCtrl.templateUrl = 'datasource/partials/config.html';
                return K8sConfigCtrl;
            })();
            exports_1("Datasource", datasource_1.K8sDatasource);
            exports_1("QueryCtrl", query_ctrl_1.K8sQueryCtrl);
            exports_1("ConfigCtrl", K8sConfigCtrl);
        }
    }
});
//# sourceMappingURL=module.js.map