///<reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
System.register(['app/plugins/sdk', 'lodash'], function(exports_1) {
    var __extends = (this && this.__extends) || function (d, b) {
        for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p];
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
    var sdk_1, lodash_1;
    var panelDefaults, PodNavCtrl;
    return {
        setters:[
            function (sdk_1_1) {
                sdk_1 = sdk_1_1;
            },
            function (lodash_1_1) {
                lodash_1 = lodash_1_1;
            }],
        execute: function() {
            panelDefaults = {};
            PodNavCtrl = (function (_super) {
                __extends(PodNavCtrl, _super);
                function PodNavCtrl($scope, $injector, backendSrv, datasourceSrv, $location, alertSrv, variableSrv, $q) {
                    _super.call(this, $scope, $injector);
                    this.backendSrv = backendSrv;
                    this.datasourceSrv = datasourceSrv;
                    this.$location = $location;
                    this.alertSrv = alertSrv;
                    this.variableSrv = variableSrv;
                    this.$q = $q;
                    lodash_1.default.defaults(this.panel, panelDefaults);
                    this.templateVariables = this.variableSrv.variables;
                    this.namespace = "All";
                    this.currentTags = {};
                    this.currentPods = [];
                    this.selectedPods = [];
                    this.setDefaults();
                    this.loadTags();
                    this.chosenTags = {};
                }
                PodNavCtrl.prototype.refresh = function () {
                    if (this.needsRefresh()) {
                        this.currentTags = {};
                        this.currentPods = [];
                        this.chosenTags = {};
                        this.selectedPods = [];
                        this.setDefaults();
                        this.loadTags();
                    }
                };
                PodNavCtrl.prototype.needsRefresh = function () {
                    var cluster = lodash_1.default.find(this.templateVariables, { 'name': 'cluster' });
                    var ns = lodash_1.default.find(this.templateVariables, { 'name': 'namespace' });
                    if (this.clusterName !== cluster.current.value) {
                        return true;
                    }
                    if ((ns.current.value === '$__all' || ns.current.value[0] === '$__all')
                        && (this.namespace === ns.current.value || this.namespace === '')) {
                        return false;
                    }
                    if (ns.current.value !== this.namespace) {
                        return true;
                    }
                    return false;
                };
                PodNavCtrl.prototype.loadTags = function () {
                    var _this = this;
                    this.getCluster().then(function () {
                        return _this.getPods().then(function (pods) {
                            _this.parseTagsFromPods(pods);
                            _this.currentPods = lodash_1.default.uniq(lodash_1.default.map(pods, function (p) { return p.metadata.name; }));
                        });
                    });
                };
                PodNavCtrl.prototype.setDefaults = function () {
                    var cluster = lodash_1.default.find(this.templateVariables, { 'name': 'cluster' });
                    if (!cluster) {
                        this.alertSrv.set("no cluster specified.", "no cluster specified in url", 'error');
                        return;
                    }
                    var ns = lodash_1.default.find(this.templateVariables, { 'name': 'namespace' });
                    this.namespace = ns.current.value !== '$__all' && ns.current.value[0] !== '$__all' ? ns.current.value : '';
                    var podVariable = lodash_1.default.find(this.templateVariables, { 'name': 'pod' });
                    if (podVariable.current.value !== '$__all') {
                        this.selectedPods = lodash_1.default.isArray(podVariable.current.value) ? podVariable.current.value : [podVariable.current.value];
                    }
                };
                PodNavCtrl.prototype.getPods = function () {
                    var _this = this;
                    if (this.currentPods.length === 0) {
                        if (lodash_1.default.isArray(this.namespace)) {
                            var promises = [];
                            lodash_1.default.forEach(this.namespace, function (ns) {
                                promises.push(_this.clusterDS.getPods(ns));
                            });
                            return this.$q.all(promises)
                                .then(function (pods) {
                                return lodash_1.default.flatten(pods).filter(function (n) { return n; });
                            });
                        }
                        else {
                            return this.clusterDS.getPods(this.namespace);
                        }
                    }
                    else {
                        return this.clusterDS.getPodsByName(this.currentPods);
                    }
                };
                PodNavCtrl.prototype.parseTagsFromPods = function (pods) {
                    var _this = this;
                    this.currentTags = {};
                    lodash_1.default.forEach(pods, function (pod) {
                        lodash_1.default.forEach(pod.metadata.labels, function (value, label) {
                            if (!_this.currentTags[label]) {
                                _this.currentTags[label] = [];
                            }
                            if (!_this.currentTags[label].includes(value)) {
                                _this.currentTags[label].push(value);
                            }
                        });
                    });
                };
                PodNavCtrl.prototype.onTagSelect = function () {
                    var _this = this;
                    this.removeEmptyTags();
                    this.selectedPods = [];
                    this.getPodsByLabel()
                        .then(function (pods) {
                        _this.currentPods = lodash_1.default.uniq(lodash_1.default.map(pods, function (p) { return p.metadata.name; }));
                        _this.parseTagsFromPods(pods);
                        _this.updateTemplateVariableWithPods();
                    });
                };
                PodNavCtrl.prototype.getPodsByLabel = function () {
                    var _this = this;
                    if (lodash_1.default.isArray(this.namespace)) {
                        var promises = [];
                        lodash_1.default.forEach(this.namespace, function (ns) {
                            promises.push(_this.clusterDS.getPodsByLabel(ns, _this.chosenTags));
                        });
                        return this.$q.all(promises)
                            .then(function (pods) {
                            return lodash_1.default.flatten(pods).filter(function (n) { return n; });
                        });
                    }
                    else {
                        return this.clusterDS.getPodsByLabel(this.namespace, this.chosenTags);
                    }
                };
                PodNavCtrl.prototype.updateTemplateVariableWithPods = function () {
                    var _this = this;
                    var variable = lodash_1.default.find(this.templateVariables, { 'name': 'pod' });
                    if (this.selectedPods.length > 0) {
                        variable.current.text = this.selectedPods.join(' + ');
                        variable.current.value = this.selectedPods;
                    }
                    else {
                        variable.current.text = lodash_1.default.isEmpty(this.chosenTags) ? 'All' : this.currentPods.join(' + ');
                        variable.current.value = lodash_1.default.isEmpty(this.chosenTags) ? ['.+'] : this.currentPods;
                    }
                    this.variableSrv.updateOptions(variable).then(function () {
                        _this.variableSrv.variableUpdated(variable).then(function () {
                            _this.$scope.$emit('template-variable-value-updated');
                            _this.$scope.$root.$broadcast('refresh');
                        });
                    });
                };
                PodNavCtrl.prototype.removeEmptyTags = function () {
                    this.chosenTags = lodash_1.default.omitBy(this.chosenTags, function (val) { return !val; });
                };
                PodNavCtrl.prototype.getCluster = function () {
                    var _this = this;
                    var clusterName = lodash_1.default.find(this.templateVariables, { 'name': 'cluster' }).current.value;
                    this.clusterName = clusterName;
                    return this.backendSrv.get('/api/datasources')
                        .then(function (result) {
                        return lodash_1.default.filter(result, { "name": clusterName })[0];
                    })
                        .then(function (ds) {
                        if (!ds) {
                            _this.alertSrv.set("Failed to connect", "Could not connect to the specified cluster.", 'error');
                            throw "Failed to connect to " + clusterName;
                        }
                        if (!(ds.jsonData.ds)) {
                            ds.jsonData.ds = "";
                        }
                        return _this.datasourceSrv.get(ds.name);
                    }).then(function (clusterDS) {
                        _this.clusterDS = clusterDS;
                    });
                };
                PodNavCtrl.prototype.removeTag = function (tag) {
                    var _this = this;
                    delete this.chosenTags[tag];
                    this.getPodsByLabel()
                        .then(function (pods) {
                        _this.currentPods = lodash_1.default.uniq(lodash_1.default.map(pods, function (p) { return p.metadata.name; }));
                        _this.parseTagsFromPods(pods);
                        _this.updateTemplateVariableWithPods();
                    });
                };
                PodNavCtrl.prototype.selectPod = function (podName) {
                    this.chosenTags = {};
                    if (!this.selectedPods.includes(podName)) {
                        this.selectedPods.push(podName);
                    }
                    this.updateTemplateVariableWithPods();
                };
                PodNavCtrl.prototype.removePodTag = function (podName) {
                    lodash_1.default.remove(this.selectedPods, function (p) { return p === podName; });
                    this.updateTemplateVariableWithPods();
                    if (this.selectedPods.length === 0) {
                        this.currentPods = [];
                        this.loadTags();
                    }
                };
                PodNavCtrl.templateUrl = 'panels/podNav/partials/pod_nav.html';
                return PodNavCtrl;
            })(sdk_1.PanelCtrl);
            exports_1("PodNavCtrl", PodNavCtrl);
        }
    }
});
//# sourceMappingURL=podNav.js.map