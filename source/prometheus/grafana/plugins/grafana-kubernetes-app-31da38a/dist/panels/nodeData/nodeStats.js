///<reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
System.register(['app/core/utils/kbn', 'lodash', 'moment'], function(exports_1) {
    var kbn_1, lodash_1, moment_1;
    var NodeStatsDatasource;
    function slugify(str) {
        var slug = str.replace("@", "at").replace("&", "and").replace(/[.]/g, "_").replace("/\W+/", "");
        return slug;
    }
    return {
        setters:[
            function (kbn_1_1) {
                kbn_1 = kbn_1_1;
            },
            function (lodash_1_1) {
                lodash_1 = lodash_1_1;
            },
            function (moment_1_1) {
                moment_1 = moment_1_1;
            }],
        execute: function() {
            NodeStatsDatasource = (function () {
                function NodeStatsDatasource(datasourceSrv, timeSrv) {
                    this.datasourceSrv = datasourceSrv;
                    this.timeSrv = timeSrv;
                }
                NodeStatsDatasource.prototype.issuePrometheusQuery = function (prometheusDS, query) {
                    return this.datasourceSrv.get(prometheusDS)
                        .then(function (datasource) {
                        var metricsQuery = {
                            range: { from: moment_1.default().subtract(5, 'minute'), to: moment_1.default() },
                            targets: [{ expr: query.expr, format: 'time_series' }],
                            legendFormat: query.legend,
                            interval: '60s',
                        };
                        return datasource.query(metricsQuery);
                    }).then(function (result) {
                        if (result && result.data) {
                            return result.data;
                        }
                        return {};
                    });
                };
                NodeStatsDatasource.prototype.getNodeStats = function (cluster_id, prometheusDS) {
                    var _this = this;
                    var podsPerNode, cpuPerNode, memoryPerNode;
                    var podQuery = {
                        expr: 'sum(label_join(kubelet_running_pod_count, "node",  "", "kubernetes_io_hostname")) by (node)',
                        legend: "{{node}}",
                    };
                    var cpuQuery = {
                        expr: 'sum(kube_pod_container_resource_requests_cpu_cores) by (node)',
                        legend: "{{node}}",
                    };
                    var memoryQuery = {
                        expr: 'sum(kube_pod_container_resource_requests_memory_bytes) by (node)',
                        legend: "{{node}}",
                    };
                    return this.issuePrometheusQuery(prometheusDS, podQuery)
                        .then(function (data) {
                        podsPerNode = data;
                        return;
                    }).then(function () {
                        return _this.issuePrometheusQuery(prometheusDS, cpuQuery);
                    })
                        .then(function (data) {
                        cpuPerNode = data;
                        return;
                    }).then(function () {
                        return _this.issuePrometheusQuery(prometheusDS, memoryQuery);
                    })
                        .then(function (data) {
                        memoryPerNode = data;
                        return { podsPerNode: podsPerNode, cpuPerNode: cpuPerNode, memoryPerNode: memoryPerNode };
                    });
                };
                NodeStatsDatasource.prototype.updateNodeWithStats = function (node, nodeStats) {
                    var formatFunc = kbn_1.default.valueFormats['percentunit'];
                    var nodeName = slugify(node.metadata.name);
                    var findFunction = function (o) { return o.target.substring(7, o.target.length - 2) === nodeName; };
                    var podsUsedData = lodash_1.default.find(nodeStats.podsPerNode, findFunction);
                    if (podsUsedData) {
                        node.podsUsed = lodash_1.default.last(podsUsedData.datapoints)[0];
                        node.podsUsedPerc = formatFunc(node.podsUsed / node.status.capacity.pods, 2, 5);
                    }
                    var cpuData = lodash_1.default.find(nodeStats.cpuPerNode, findFunction);
                    if (cpuData) {
                        node.cpuUsage = lodash_1.default.last(cpuData.datapoints)[0];
                        node.cpuUsageFormatted = kbn_1.default.valueFormats['none'](node.cpuUsage, 2, null);
                        node.cpuUsagePerc = formatFunc(node.cpuUsage / node.status.capacity.cpu, 2, 5);
                    }
                    var memData = lodash_1.default.find(nodeStats.memoryPerNode, findFunction);
                    if (memData) {
                        node.memoryUsage = lodash_1.default.last(memData.datapoints)[0];
                        var memCapacity = node.status.capacity.memory.substring(0, node.status.capacity.memory.length - 2) * 1000;
                        node.memUsageFormatted = kbn_1.default.valueFormats['bytes'](node.memoryUsage, 2, null);
                        node.memCapacityFormatted = kbn_1.default.valueFormats['bytes'](memCapacity, 2, null);
                        node.memoryUsagePerc = formatFunc((node.memoryUsage / memCapacity), 2, 5);
                    }
                    return node;
                };
                return NodeStatsDatasource;
            })();
            exports_1("default", NodeStatsDatasource);
        }
    }
});
//# sourceMappingURL=nodeStats.js.map