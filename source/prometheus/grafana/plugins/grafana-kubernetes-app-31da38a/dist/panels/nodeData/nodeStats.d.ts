/// <reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
export default class NodeStatsDatasource {
    private datasourceSrv;
    private timeSrv;
    constructor(datasourceSrv: any, timeSrv: any);
    issuePrometheusQuery(prometheusDS: any, query: any): any;
    getNodeStats(cluster_id: any, prometheusDS: any): any;
    updateNodeWithStats(node: any, nodeStats: any): any;
}
