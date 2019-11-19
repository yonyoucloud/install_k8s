///<reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />

import {PanelCtrl} from 'app/plugins/sdk';
import _ from 'lodash';

const panelDefaults = {
};

export class PodNavCtrl extends PanelCtrl {
  templateVariables: any;
  namespace: string;
  currentTags: any;
  currentPods: any[];
  selectedPods: any;
  chosenTags: any;
  clusterName: string;
  clusterDS: any;


  static templateUrl = 'panels/podNav/partials/pod_nav.html';

  constructor($scope, $injector, private backendSrv, private datasourceSrv, private $location, private alertSrv, private variableSrv, private $q) {
    super($scope, $injector);
    _.defaults(this.panel, panelDefaults);

    this.templateVariables = this.variableSrv.variables;
    this.namespace = "All";
    this.currentTags = {};
    this.currentPods = [];
    this.selectedPods = [];

    this.setDefaults();
    this.loadTags();
    this.chosenTags = {};
  }

  refresh() {
    if (this.needsRefresh()) {
      this.currentTags = {};
      this.currentPods = [];
      this.chosenTags = {};
      this.selectedPods = [];

      this.setDefaults();
      this.loadTags();
    }
  }

  needsRefresh() {
    const cluster = _.find(this.templateVariables, {'name': 'cluster'});
    const ns = _.find(this.templateVariables, {'name': 'namespace'});

    if (this.clusterName !== cluster.current.value) { return true; }

    if ((ns.current.value === '$__all' || ns.current.value[0] === '$__all')
      && (this.namespace === ns.current.value || this.namespace === '')) {
      return false;
    }

    if (ns.current.value !== this.namespace) { return true; }

    return false;
  }

  loadTags() {
    this.getCluster().then(() => {
      return this.getPods().then(pods => {
        this.parseTagsFromPods(pods);
        this.currentPods = _.uniq(_.map(pods, p => { return p.metadata.name; }));
      });
    });
  }

  setDefaults() {
    const cluster = _.find(this.templateVariables, {'name': 'cluster'});
    if (!cluster) {
      this.alertSrv.set("no cluster specified.", "no cluster specified in url", 'error');
      return;
    }

    const ns = _.find(this.templateVariables, {'name': 'namespace'});
    this.namespace = ns.current.value !== '$__all' && ns.current.value[0] !== '$__all' ? ns.current.value : '';
    const podVariable = _.find(this.templateVariables, {'name': 'pod'});

    if (podVariable.current.value !== '$__all') {
      this.selectedPods = _.isArray(podVariable.current.value) ? podVariable.current.value : [podVariable.current.value];
    }
  }

  getPods() {
    if (this.currentPods.length === 0) {
      if (_.isArray(this.namespace)) {
        const promises = [];
        _.forEach(this.namespace, ns => {
          promises.push(this.clusterDS.getPods(ns));
        });
        return this.$q.all(promises)
        .then(pods => {
          return _.flatten(pods).filter(n => n);
        });
      } else {
        return this.clusterDS.getPods(this.namespace);
      }
    } else {
      return this.clusterDS.getPodsByName(this.currentPods);
    }
  }

  parseTagsFromPods(pods) {
    this.currentTags = {};

    _.forEach(pods, pod => {
      _.forEach(pod.metadata.labels, (value, label) => {
        if (!this.currentTags[label]) {
          this.currentTags[label] = [];
        }
        if (!this.currentTags[label].includes(value)) {
          this.currentTags[label].push(value);
        }
      });
    });
  }

  onTagSelect() {
    this.removeEmptyTags();
    this.selectedPods = [];

    this.getPodsByLabel()
    .then(pods => {
      this.currentPods = _.uniq(_.map(pods, p => { return p.metadata.name; }));
      this.parseTagsFromPods(pods);
      this.updateTemplateVariableWithPods();
    });
  }

  getPodsByLabel() {
    if (_.isArray(this.namespace)) {
      const promises = [];
      _.forEach(this.namespace, ns => {
        promises.push(this.clusterDS.getPodsByLabel(ns, this.chosenTags));
      });
      return this.$q.all(promises)
      .then(pods => {
        return _.flatten(pods).filter(n => n);
      });
    } else {
      return this.clusterDS.getPodsByLabel(this.namespace, this.chosenTags);
    }
  }

  updateTemplateVariableWithPods() {
    const variable = _.find(this.templateVariables, {'name': 'pod'});

    if (this.selectedPods.length > 0) {
      variable.current.text = this.selectedPods.join(' + ');
      variable.current.value = this.selectedPods;
    } else {
      variable.current.text = _.isEmpty(this.chosenTags) ? 'All': this.currentPods.join(' + ');
      variable.current.value = _.isEmpty(this.chosenTags) ? ['.+']: this.currentPods;
    }

    this.variableSrv.updateOptions(variable).then(() => {
      this.variableSrv.variableUpdated(variable).then(() => {
        this.$scope.$emit('template-variable-value-updated');
        this.$scope.$root.$broadcast('refresh');
      });
    });
  }

  removeEmptyTags() {
    this.chosenTags = _.omitBy(this.chosenTags, val => { return !val;});
  }

  getCluster() {
    const clusterName = _.find(this.templateVariables, {'name': 'cluster'}).current.value;
    this.clusterName = clusterName;

    return this.backendSrv.get('/api/datasources')
    .then(result => {
      return _.filter(result, {"name": clusterName})[0];
    })
    .then((ds) => {
      if (!ds) {
        this.alertSrv.set("Failed to connect", "Could not connect to the specified cluster.", 'error');
        throw "Failed to connect to " + clusterName;
      }

      if (!(ds.jsonData.ds)) {
        ds.jsonData.ds = "";
      }
      return this.datasourceSrv.get(ds.name);
    }).then(clusterDS => {
      this.clusterDS = clusterDS;
    });
  }

  removeTag(tag) {
    delete this.chosenTags[tag];
    this.getPodsByLabel()
    .then(pods => {
      this.currentPods = _.uniq(_.map(pods, p => { return p.metadata.name; }));
      this.parseTagsFromPods(pods);
      this.updateTemplateVariableWithPods();
    });
  }

  selectPod(podName) {
    this.chosenTags = {};

    if (!this.selectedPods.includes(podName)) {
      this.selectedPods.push(podName);
    }

    this.updateTemplateVariableWithPods();
  }

  removePodTag(podName) {
    _.remove(this.selectedPods, p => { return p === podName;});
    this.updateTemplateVariableWithPods();

    if (this.selectedPods.length === 0) {
      this.currentPods = [];
      this.loadTags();
    }
  }
}
