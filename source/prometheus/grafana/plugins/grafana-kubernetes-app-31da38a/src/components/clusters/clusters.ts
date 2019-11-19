///<reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />

import _ from 'lodash';
import appEvents from 'app/core/app_events';

export class ClustersCtrl {
  cluster: any;
  pageReady: boolean;
  datasources: [any];
  clusters: {};
  isOrgEditor: boolean;

  static templateUrl = 'components/clusters/partials/clusters.html';

  /** @ngInject */
  constructor($scope, $injector, private backendSrv, private contextSrv, private $location) {
    var self = this;
    this.isOrgEditor = contextSrv.hasRole('Editor') || contextSrv.hasRole('Admin');
    document.title = 'Grafana Kubernetes App';
    this.clusters = {};
    this.pageReady = false;
    this.getClusters().then(() => {
      self.pageReady = true;
    });
  }

  getClusters() {
    var self = this;
    return this.backendSrv.get('/api/datasources')
    .then((result) => {
      self.clusters = _.filter(result, {"type": "grafana-kubernetes-datasource"});
    });
  }

  confirmDelete(id) {
    this.backendSrv.delete('/api/datasources/' + id).then(() => {
      this.getClusters();
    });
  }

  deleteCluster(cluster) {
    appEvents.emit('confirm-modal', {
      title: 'Delete',
      text: 'Are you sure you want to delete this data source? ' +
        'If you need to undeploy the collectors, then do that before deleting the data source.',
      yesText: "Delete",
      icon: "fa-trash",
      onConfirm: () => {
        this.confirmDelete(cluster.id);
      }
    });
  }

  clusterInfo(cluster) {
    this.$location.path("plugins/grafana-kubernetes-app/page/cluster-info").search({"cluster": cluster.id});
  }
}
