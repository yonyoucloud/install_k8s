/// <reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
export declare class ClusterInfoCtrl {
    private backendSrv;
    private datasourceSrv;
    private $q;
    private $location;
    private alertSrv;
    cluster: any;
    pageReady: boolean;
    componentStatuses: any;
    namespaces: string[];
    namespace: string;
    nodes: any[];
    datasources: any;
    clusterDS: any;
    static templateUrl: string;
    /** @ngInject */
    constructor($scope: any, $injector: any, backendSrv: any, datasourceSrv: any, $q: any, $location: any, alertSrv: any);
    getCluster(id: any): any;
    getClusterInfo(): void;
    goToClusterDashboard(): void;
    goToPodDashboard(): void;
    goToNodeDashboard(node: any, evt: any): void;
    goToWorkloads(ns: any, evt: any): void;
    goToNodeInfo(node: any, evt: any): void;
}
