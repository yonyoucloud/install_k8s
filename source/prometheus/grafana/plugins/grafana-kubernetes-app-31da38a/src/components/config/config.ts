///<reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />

export class KubernetesConfigCtrl {
  static templateUrl = 'components/config/config.html';
  enabled: boolean;
  appEditCtrl: any;
  appModel: any;

  /** @ngInject */
  constructor($scope, $injector, private $q) {
    this.enabled = false;
    this.appEditCtrl.setPostUpdateHook(this.postUpdate.bind(this));
  }

  postUpdate() {
    if (!this.appModel.enabled) {
      return this.$q.resolve();
    }
    return this.appEditCtrl.importDashboards().then(() => {
      this.enabled = true;
      return {
        url: "plugins/grafana-kubernetes-app/page/clusters",
        message: "Kubernetes App enabled!"
      };
    });
  }
}
