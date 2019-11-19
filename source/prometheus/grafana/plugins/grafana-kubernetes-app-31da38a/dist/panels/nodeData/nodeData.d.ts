/// <reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
import { PanelCtrl } from 'app/plugins/sdk';
import NodeStatsDatasource from './nodeStats';
export declare class NodeDataCtrl extends PanelCtrl {
    private backendSrv;
    private datasourceSrv;
    private $location;
    private alertSrv;
    private timeSrv;
    private variableSrv;
    templateVariables: any;
    nodeStatsDatasource: NodeStatsDatasource;
    pageReady: boolean;
    cluster: any;
    clusterDS: any;
    node: any;
    isInListMode: boolean;
    nodes: any[];
    static templateUrl: string;
    static scrollable: boolean;
    /** @ngInject */
    constructor($scope: any, $injector: any, backendSrv: any, datasourceSrv: any, $location: any, alertSrv: any, timeSrv: any, variableSrv: any);
    loadCluster(): void;
    getNodeHealth(node: any): {
        text: string;
        iconClass: string;
        stateClass: string;
        message: any;
    };
    getHealthState(health: any, message: any): {
        text: string;
        iconClass: string;
        stateClass: string;
        message: any;
    };
    refresh(): void;
    loadDatasource(id: any): any;
    goToNodeDashboard(node: any): void;
    conditionStatus(condition: any): {
        value: any;
        text: string;
    };
    isConditionOk(condition: any): any;
    conditionLastTransitionTime(condition: any): any;
}
