///<reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
System.register(['./podNav', 'app/plugins/sdk'], function(exports_1) {
    var podNav_1, sdk_1;
    return {
        setters:[
            function (podNav_1_1) {
                podNav_1 = podNav_1_1;
            },
            function (sdk_1_1) {
                sdk_1 = sdk_1_1;
            }],
        execute: function() {
            sdk_1.loadPluginCss({
                dark: 'plugins/grafana-kubernetes-app/css/dark.css',
                light: 'plugins/grafana-kubernetes-app/css/light.css'
            });
            exports_1("PanelCtrl", podNav_1.PodNavCtrl);
        }
    }
});
//# sourceMappingURL=module.js.map