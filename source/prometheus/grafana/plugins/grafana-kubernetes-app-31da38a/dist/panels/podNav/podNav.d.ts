/// <reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
import { PanelCtrl } from 'app/plugins/sdk';
export declare class PodNavCtrl extends PanelCtrl {
    private backendSrv;
    private datasourceSrv;
    private $location;
    private alertSrv;
    private variableSrv;
    private $q;
    templateVariables: any;
    namespace: string;
    currentTags: any;
    currentPods: any[];
    selectedPods: any;
    chosenTags: any;
    clusterName: string;
    clusterDS: any;
    static templateUrl: string;
    constructor($scope: any, $injector: any, backendSrv: any, datasourceSrv: any, $location: any, alertSrv: any, variableSrv: any, $q: any);
    refresh(): void;
    needsRefresh(): boolean;
    loadTags(): void;
    setDefaults(): void;
    getPods(): any;
    parseTagsFromPods(pods: any): void;
    onTagSelect(): void;
    getPodsByLabel(): any;
    updateTemplateVariableWithPods(): void;
    removeEmptyTags(): void;
    getCluster(): any;
    removeTag(tag: any): void;
    selectPod(podName: any): void;
    removePodTag(podName: any): void;
}
