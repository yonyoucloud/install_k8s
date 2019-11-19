///<reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />

import moment from 'moment';
import {PanelCtrl} from 'app/plugins/sdk';
import _ from 'lodash';
import NodeStatsDatasource from './nodeStats';

const panelDefaults = {
};

export class NodeDataCtrl extends PanelCtrl {
  templateVariables: any;
  nodeStatsDatasource: NodeStatsDatasource;
  pageReady: boolean;
  cluster: any;
  clusterDS: any;
  node: any;
  isInListMode: boolean;
  nodes: any[];


  static templateUrl = 'panels/nodeData/partials/node_info.html';
  static scrollable = true;

  /** @ngInject */
  constructor($scope, $injector, private backendSrv, private datasourceSrv, private $location, private alertSrv, private timeSrv, private variableSrv) {
    super($scope, $injector);
    _.defaults(this.panel, panelDefaults);

    this.templateVariables = this.variableSrv.variables;
    this.nodeStatsDatasource = new NodeStatsDatasource(datasourceSrv, timeSrv);
    document.title = 'Grafana Kubernetes App';

    this.pageReady = false;
    this.cluster = {};
    this.clusterDS = {};
    this.node = {};

    this.isInListMode = false;
    this.nodes = [];

    this.loadCluster();
  }

  loadCluster() {
    const cluster = _.find(this.templateVariables, {'name': 'cluster'});
    if (!cluster) {
      this.alertSrv.set("no cluster specified.", "no cluster specified in url", 'error');
      return;
    } else {
      const cluster_id = cluster.current.value;
      const nodeVar = _.find(this.templateVariables, {'name': 'node'});
      const node_name  = nodeVar.current.value !== '$__all' ? nodeVar.current.value : 'All';
      const prometheusDS  = _.find(this.templateVariables, {'name': 'datasource'}).current.value;

      this.loadDatasource(cluster_id).then(() => {
        return this.nodeStatsDatasource.getNodeStats(cluster_id, prometheusDS);
      }).then(nodeStats => {
        if (node_name === 'All') {
          this.isInListMode = true;
          this.clusterDS.getNodes().then(nodes => {
            this.nodes = _.map(nodes, node => {
              node.healthState = this.getNodeHealth(node);
              this.nodeStatsDatasource.updateNodeWithStats(node, nodeStats);

              return node;
            });
          });
        } else {
          this.isInListMode = false;
          this.clusterDS.getNode(unslugify(node_name)).then(node => {
            this.node = node;
            this.pageReady = true;
          });
        }
      });
    }
  }

  getNodeHealth(node) {
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
    return this.getHealthState(health, message);
  }

  getHealthState(health, message) {
    switch (health) {
      case 'ok': {
        return {
          text: 'OK',
          iconClass: 'icon-gf icon-gf-online',
          stateClass: 'alert-state-ok',
          message: '',
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

  refresh() {
    this.loadCluster();
  }

  loadDatasource(id) {
    return this.backendSrv.get('api/datasources')
      .then(result => {
        return _.filter(result, {"type": "grafana-kubernetes-datasource", "name": id})[0];
      })
      .then(ds => {
        if (!ds) {
          this.alertSrv.set("Failed to connect", "Could not connect to the specified cluster.", 'error');
          throw "Failed to connect to " + id;
        }
        this.cluster = ds;
        return this.datasourceSrv.get(ds.name);
      }).then(clusterDS => {
        this.clusterDS = clusterDS;
        return clusterDS;
      });
  }

  goToNodeDashboard(node) {
    const variable = _.find(this.templateVariables, {'name': 'node'});
    variable.current.text = node === 'All' ? 'All': slugify(node.metadata.name);
    variable.current.value = node === 'All' ? '$__all': slugify(node.metadata.name);

    this.variableSrv.variableUpdated(variable).then(() => {
      this.$scope.$emit('template-variable-value-updated');
      this.$scope.$root.$broadcast('refresh');
    });
  }

  conditionStatus(condition) {
    var status;
    if (condition.type === "Ready") {
      status = condition.status === "True";
    } else {
      status = condition.status === "False";
    }

    return {
      value: status,
      text: status ? "Ok" : "Error"
    };
  }

  isConditionOk(condition) {
    return this.conditionStatus(condition).value;
  }

  conditionLastTransitionTime(condition) {
    return moment(condition.lastTransitionTime).format('YYYY-MM-DD HH:mm:ss');
  }
}

function slugify(str) {
  var slug = str.replace("@", "at").replace("&", "and").replace(/[.]/g, "_").replace("/\W+/", "");
  return slug;
}

function unslugify(str) {
  var slug = str.replace(/[_]/g, ".");
  return slug;
}
