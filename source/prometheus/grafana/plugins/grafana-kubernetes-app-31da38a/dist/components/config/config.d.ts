/// <reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
export declare class KubernetesConfigCtrl {
    private $q;
    static templateUrl: string;
    enabled: boolean;
    appEditCtrl: any;
    appModel: any;
    /** @ngInject */
    constructor($scope: any, $injector: any, $q: any);
    postUpdate(): any;
}
