System.register(['./components/config/config', './components/clusters/clusters', './components/clusters/clusterConfig', './components/clusters/clusterInfo', './components/clusters/clusterWorkloads', './components/clusters/nodeInfo', './components/clusters/podInfo', 'app/plugins/sdk'], function(exports_1) {
    var config_1, clusters_1, clusterConfig_1, clusterInfo_1, clusterWorkloads_1, nodeInfo_1, podInfo_1, sdk_1;
    return {
        setters:[
            function (config_1_1) {
                config_1 = config_1_1;
            },
            function (clusters_1_1) {
                clusters_1 = clusters_1_1;
            },
            function (clusterConfig_1_1) {
                clusterConfig_1 = clusterConfig_1_1;
            },
            function (clusterInfo_1_1) {
                clusterInfo_1 = clusterInfo_1_1;
            },
            function (clusterWorkloads_1_1) {
                clusterWorkloads_1 = clusterWorkloads_1_1;
            },
            function (nodeInfo_1_1) {
                nodeInfo_1 = nodeInfo_1_1;
            },
            function (podInfo_1_1) {
                podInfo_1 = podInfo_1_1;
            },
            function (sdk_1_1) {
                sdk_1 = sdk_1_1;
            }],
        execute: function() {
            sdk_1.loadPluginCss({
                dark: 'plugins/grafana-kubernetes-app/css/dark.css',
                light: 'plugins/grafana-kubernetes-app/css/light.css'
            });
            exports_1("ConfigCtrl", config_1.KubernetesConfigCtrl);
            exports_1("ClustersCtrl", clusters_1.ClustersCtrl);
            exports_1("ClusterConfigCtrl", clusterConfig_1.ClusterConfigCtrl);
            exports_1("ClusterInfoCtrl", clusterInfo_1.ClusterInfoCtrl);
            exports_1("ClusterWorkloadsCtrl", clusterWorkloads_1.ClusterWorkloadsCtrl);
            exports_1("NodeInfoCtrl", nodeInfo_1.NodeInfoCtrl);
            exports_1("PodInfoCtrl", podInfo_1.PodInfoCtrl);
        }
    }
});
//# sourceMappingURL=module.js.map