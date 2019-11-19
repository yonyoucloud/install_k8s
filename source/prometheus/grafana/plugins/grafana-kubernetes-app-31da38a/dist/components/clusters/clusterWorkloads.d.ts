/// <reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
export declare class ClusterWorkloadsCtrl {
    private backendSrv;
    private datasourceSrv;
    private $q;
    private $location;
    private alertSrv;
    pageReady: boolean;
    cluster: any;
    namespaces: string[];
    namespace: string;
    daemonSets: any[];
    replicationControllers: any[];
    deployments: any[];
    pods: any[];
    clusterDS: any;
    static templateUrl: string;
    /** @ngInject */
    constructor($scope: any, $injector: any, backendSrv: any, datasourceSrv: any, $q: any, $location: any, alertSrv: any);
    getCluster(id: any): any;
    getWorkloads(): void;
    componentHealth(component: any): string;
    isComponentHealthy(component: any): boolean;
    goToPodDashboard(pod: any): void;
    goToDeploymentDashboard(deploy: any): void;
    goToPodInfo(pod: any, evt: any): void;
}
