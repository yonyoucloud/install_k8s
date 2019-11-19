/// <reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
export declare class ClusterConfigCtrl {
    private backendSrv;
    private $q;
    private contextSrv;
    private $location;
    private $window;
    private alertSrv;
    cluster: any;
    isOrgEditor: boolean;
    pageReady: boolean;
    prometheusDeployed: boolean;
    showHelp: boolean;
    showPrometheusExample: boolean;
    datasources: [any];
    static templateUrl: string;
    /** @ngInject */
    constructor($scope: any, $injector: any, backendSrv: any, $q: any, contextSrv: any, $location: any, $window: any, alertSrv: any);
    toggleHelp(): void;
    togglePrometheusExample(): void;
    getDatasources(): any;
    getCluster(id: any): any;
    getPrometheusDatasources(): any;
    getDeployments(): any;
    save(): any;
    savePrometheusConfigToFile(): void;
    saveNodeExporterDSToFile(): void;
    saveKubeStateDeployToFile(): void;
    saveToFile(filename: any, blob: any): void;
    deploy(): void;
    undeploy(): void;
    saveDatasource(): any;
    saveAndDeploy(): any;
    checkApiVersion(clusterId: any): any;
    createConfigMap(clusterId: any, cm: any): any;
    createDaemonSet(clusterId: any, daemonSet: any): any;
    deleteDaemonSet(clusterId: any): any;
    createDeployment(clusterId: any, deployment: any): any;
    deleteDeployment(clusterId: any, deploymentName: any): any;
    deleteConfigMap(clusterId: any, cmName: any): any;
    deletePods(): any;
    cancel(): void;
    deployPrometheus(): any;
    undeployPrometheus(): any;
    generatePrometheusConfig(): string;
    generatePrometheusConfigMap(): {
        "apiVersion": string;
        "kind": string;
        "metadata": {
            "name": string;
        };
        "data": {
            "prometheus.yml": string;
        };
    };
}
