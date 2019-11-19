/// <reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
export declare class NodeInfoCtrl {
    private backendSrv;
    private datasourceSrv;
    private $q;
    private $location;
    private alertSrv;
    pageReady: boolean;
    cluster: any;
    clusterDS: any;
    node: any;
    static templateUrl: string;
    /** @ngInject */
    constructor($scope: any, $injector: any, backendSrv: any, datasourceSrv: any, $q: any, $location: any, alertSrv: any);
    loadDatasource(id: any): any;
    goToNodeDashboard(): void;
    conditionStatus(condition: any): {
        value: any;
        text: string;
    };
    isConditionOk(condition: any): any;
    conditionLastTransitionTime(condition: any): any;
}
