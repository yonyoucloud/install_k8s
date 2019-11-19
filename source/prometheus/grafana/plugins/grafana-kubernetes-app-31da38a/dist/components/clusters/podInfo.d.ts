/// <reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
export declare class PodInfoCtrl {
    private backendSrv;
    private datasourceSrv;
    private $q;
    private $location;
    private alertSrv;
    pageReady: boolean;
    pod: any;
    cluster_id: any;
    clusterDS: any;
    datasource: any;
    static templateUrl: string;
    /** @ngInject */
    constructor($scope: any, $injector: any, backendSrv: any, datasourceSrv: any, $q: any, $location: any, alertSrv: any);
    loadDatasource(id: any): any;
    conditionStatus(condition: any): {
        value: any;
        text: string;
    };
    goToPodDashboard(pod: any): void;
    isConditionOk(condition: any): any;
    formatTime(time: any): any;
}
