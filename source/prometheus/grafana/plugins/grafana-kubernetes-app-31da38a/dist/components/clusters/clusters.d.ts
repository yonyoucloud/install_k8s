/// <reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
export declare class ClustersCtrl {
    private backendSrv;
    private contextSrv;
    private $location;
    cluster: any;
    pageReady: boolean;
    datasources: [any];
    clusters: {};
    isOrgEditor: boolean;
    static templateUrl: string;
    /** @ngInject */
    constructor($scope: any, $injector: any, backendSrv: any, contextSrv: any, $location: any);
    getClusters(): any;
    confirmDelete(id: any): void;
    deleteCluster(cluster: any): void;
    clusterInfo(cluster: any): void;
}
