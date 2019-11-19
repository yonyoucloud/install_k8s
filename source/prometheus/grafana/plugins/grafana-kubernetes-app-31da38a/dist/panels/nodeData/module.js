///<reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
System.register(['./nodeData', 'app/plugins/sdk'], function(exports_1) {
    var nodeData_1, sdk_1;
    return {
        setters:[
            function (nodeData_1_1) {
                nodeData_1 = nodeData_1_1;
            },
            function (sdk_1_1) {
                sdk_1 = sdk_1_1;
            }],
        execute: function() {
            sdk_1.loadPluginCss({
                dark: 'plugins/grafana-kubernetes-app/css/dark.css',
                light: 'plugins/grafana-kubernetes-app/css/light.css'
            });
            exports_1("PanelCtrl", nodeData_1.NodeDataCtrl);
        }
    }
});
//# sourceMappingURL=module.js.map