///<reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />

import _ from 'lodash';
import $ from 'jquery';

function slugify(str) {
  var slug = str.replace("@", "at").replace("&", "and").replace(/[.]/g, "_").replace("/\W+/", "");
  return slug;
}

export class ClusterInfoCtrl {
  cluster: any;
  pageReady: boolean;
  componentStatuses: any;
  namespaces: string[];
  namespace: string;
  nodes: any[];
  datasources: any;
  clusterDS: any;

  static templateUrl = 'components/clusters/partials/cluster_info.html';
  
  /** @ngInject */
  constructor($scope, $injector, private backendSrv, private datasourceSrv, private $q, private $location, private alertSrv) {
    this.$q = $q;
    document.title = 'Grafana Kubernetes App';

    this.pageReady = false;
    this.cluster = {};
    this.componentStatuses = [];
    this.namespaces = [];
    this.namespace = "";
    this.nodes = [];

    if (!("cluster" in $location.search())) {
      alertSrv.set("no cluster specified.", "no cluster specified in url", 'error');
      return;
    }

    this.getCluster($location.search().cluster)
      .then(clusterDS => {
        this.clusterDS = clusterDS;
        this.pageReady = true;
        this.getClusterInfo();
      });
  }

  getCluster(id) {
    return this.backendSrv.get('api/datasources/' + id).then(ds => {
      this.cluster = ds;
      return this.datasourceSrv.get(ds.name);
    });
  }

  getClusterInfo() {
    this.clusterDS.getComponentStatuses().then(stats => {
      this.componentStatuses = _.map(stats, stat => {
        stat.healthState = getComponentHealth(stat);
        return stat;
      });
    });
    this.clusterDS.getNamespaces().then(namespaces => {
      this.namespaces = namespaces;
    });
    this.clusterDS.getNodes().then(nodes => {
      this.nodes = _.map(nodes, node => {
        node.healthState = getNodeHealth(node);
        return node;
      });
    });
  }

  goToClusterDashboard() {
    this.$location.path("dashboard/db/k8s-cluster")
      .search({
        "var-datasource": this.cluster.jsonData.ds,
        "var-cluster": this.cluster.name
      });
  }

  goToPodDashboard() {
    this.$location.path("dashboard/db/k8s-container")
    .search({
      "var-datasource": this.cluster.jsonData.ds,
      "var-cluster": this.cluster.name,
      "var-node": 'All',
      "var-namespace": 'All',
      "var-pod": 'All'
    });
  }

  goToNodeDashboard(node, evt) {
    var clickTargetIsLinkOrHasLinkParents = $(evt.target).closest('a').length > 0;
    if (clickTargetIsLinkOrHasLinkParents === false) {
      this.$location.path("dashboard/db/k8s-node")
      .search({
        "var-datasource": this.cluster.jsonData.ds,
        "var-cluster": this.cluster.name,
        "var-node": node === 'All' ? 'All': slugify(node.metadata.name)
      });
    }
  }

  goToWorkloads(ns, evt) {
    var clickTargetIsLinkOrHasLinkParents = $(evt.target).closest('a').length > 0;
    if (clickTargetIsLinkOrHasLinkParents === false) {
      this.$location.path("plugins/grafana-kubernetes-app/page/cluster-workloads")
      .search({
        "cluster": this.cluster.id,
        "namespace": slugify(ns.metadata.name)
      });
    }
  }

  goToNodeInfo(node, evt) {
    var clickTargetIsLinkOrHasLinkParents = $(evt.target).closest('a').length > 0;

    var closestElm = _.head($(evt.target).closest('div'));
    var clickTargetClickAttr = _.find(closestElm.attributes, {name: "ng-click"});
    var clickTargetIsNodeDashboard = clickTargetClickAttr ? clickTargetClickAttr.value === "ctrl.goToNodeDashboard(node, $event)" : false;
    if (clickTargetIsLinkOrHasLinkParents === false &&
        clickTargetIsNodeDashboard === false) {
      this.$location.path("plugins/grafana-kubernetes-app/page/node-info")
      .search({
        "cluster": this.cluster.id,
        "node": node.metadata.name
      });
    }
  }
}

function getComponentHealth(component) {
  let health = "unhealthy";
  let message = '';
  _.forEach(component.conditions, condition => {
    if (condition.type   === "Healthy" &&
        condition.status === "True") {
      health = "ok";
    } else {
      message = condition.message;
    }
  });
  return getHealthState(health, message);
}

function getNodeHealth(node) {
  let health = "unhealthy";
  let message = '';
  _.forEach(node.status.conditions, condition => {
    if (condition.type   === "Ready" &&
        condition.status === "True") {
      health = "ok";
    } else {
      message = condition.message;
    }
  });
  return getHealthState(health, message);
}

function getHealthState(health, message) {
  switch (health) {
    case 'ok': {
      return {
        text: 'OK',
        iconClass: 'icon-gf icon-gf-online',
        stateClass: 'alert-state-ok',
        message: ''
      };
    }
    case 'unhealthy': {
      return {
        text: 'UNHEALTHY',
        iconClass: 'icon-gf icon-gf-critical',
        stateClass: 'alert-state-critical',
        message: message || ''
      };
    }
    case 'warning': {
      return {
        text: 'warning',
        iconClass: "icon-gf icon-gf-critical",
        stateClass: 'alert-state-warning',
        message: message || ''
      };
    }
  }
}
