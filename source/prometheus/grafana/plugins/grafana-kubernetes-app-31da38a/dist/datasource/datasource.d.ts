/// <reference path="../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
export declare class K8sDatasource {
    private backendSrv;
    private templateSrv;
    private $q;
    id: number;
    name: string;
    url: string;
    type: string;
    ds: string;
    static baseApiUrl: string;
    constructor(instanceSettings: any, backendSrv: any, templateSrv: any, $q: any);
    testDatasource(): any;
    _get(apiResource: any): any;
    getNodes(): any;
    getNode(name: any): any;
    getNamespaces(): any;
    getComponentStatuses(): any;
    getDaemonSets(namespace: any): any;
    getReplicationControllers(namespace: any): any;
    getDeployments(namespace: any): any;
    getPods(namespace: any): any;
    getPodsByLabel(namespace: any, labels: any): any;
    getPod(name: any): any;
    getPodsByName(names: any): any;
    query(options: any): void;
    annotationQuery(options: any): void;
    metricFindQuery(query: string): any;
}
