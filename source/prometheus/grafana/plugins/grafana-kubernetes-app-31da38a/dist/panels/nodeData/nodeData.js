///<reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
System.register(['moment', 'app/plugins/sdk', 'lodash', './nodeStats'], function(exports_1) {
    var __extends = (this && this.__extends) || function (d, b) {
        for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p];
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
    var moment_1, sdk_1, lodash_1, nodeStats_1;
    var panelDefaults, NodeDataCtrl;
    function slugify(str) {
        var slug = str.replace("@", "at").replace("&", "and").replace(/[.]/g, "_").replace("/\W+/", "");
        return slug;
    }
    function unslugify(str) {
        var slug = str.replace(/[_]/g, ".");
        return slug;
    }
    return {
        setters:[
            function (moment_1_1) {
                moment_1 = moment_1_1;
            },
            function (sdk_1_1) {
                sdk_1 = sdk_1_1;
            },
            function (lodash_1_1) {
                lodash_1 = lodash_1_1;
            },
            function (nodeStats_1_1) {
                nodeStats_1 = nodeStats_1_1;
            }],
        execute: function() {
            panelDefaults = {};
            NodeDataCtrl = (function (_super) {
                __extends(NodeDataCtrl, _super);
                /** @ngInject */
                function NodeDataCtrl($scope, $injector, backendSrv, datasourceSrv, $location, alertSrv, timeSrv, variableSrv) {
                    _super.call(this, $scope, $injector);
                    this.backendSrv = backendSrv;
                    this.datasourceSrv = datasourceSrv;
                    this.$location = $location;
                    this.alertSrv = alertSrv;
                    this.timeSrv = timeSrv;
                    this.variableSrv = variableSrv;
                    lodash_1.default.defaults(this.panel, panelDefaults);
                    this.templateVariables = this.variableSrv.variables;
                    this.nodeStatsDatasource = new nodeStats_1.default(datasourceSrv, timeSrv);
                    document.title = 'Grafana Kubernetes App';
                    this.pageReady = false;
                    this.cluster = {};
                    this.clusterDS = {};
                    this.node = {};
                    this.isInListMode = false;
                    this.nodes = [];
                    this.loadCluster();
                }
                NodeDataCtrl.prototype.loadCluster = function () {
                    var _this = this;
                    var cluster = lodash_1.default.find(this.templateVariables, { 'name': 'cluster' });
                    if (!cluster) {
                        this.alertSrv.set("no cluster specified.", "no cluster specified in url", 'error');
                        return;
                    }
                    else {
                        var cluster_id = cluster.current.value;
                        var nodeVar = lodash_1.default.find(this.templateVariables, { 'name': 'node' });
                        var node_name = nodeVar.current.value !== '$__all' ? nodeVar.current.value : 'All';
                        var prometheusDS = lodash_1.default.find(this.templateVariables, { 'name': 'datasource' }).current.value;
                        this.loadDatasource(cluster_id).then(function () {
                            return _this.nodeStatsDatasource.getNodeStats(cluster_id, prometheusDS);
                        }).then(function (nodeStats) {
                            if (node_name === 'All') {
                                _this.isInListMode = true;
                                _this.clusterDS.getNodes().then(function (nodes) {
                                    _this.nodes = lodash_1.default.map(nodes, function (node) {
                                        node.healthState = _this.getNodeHealth(node);
                                        _this.nodeStatsDatasource.updateNodeWithStats(node, nodeStats);
                                        return node;
                                    });
                                });
                            }
                            else {
                                _this.isInListMode = false;
                                _this.clusterDS.getNode(unslugify(node_name)).then(function (node) {
                                    _this.node = node;
                                    _this.pageReady = true;
                                });
                            }
                        });
                    }
                };
                NodeDataCtrl.prototype.getNodeHealth = function (node) {
                    var health = "unhealthy";
                    var message = '';
                    lodash_1.default.forEach(node.status.conditions, function (condition) {
                        if (condition.type === "Ready" &&
                            condition.status === "True") {
                            health = "ok";
                        }
                        else {
                            message = condition.message;
                        }
                    });
                    return this.getHealthState(health, message);
                };
                NodeDataCtrl.prototype.getHealthState = function (health, message) {
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
                };
                NodeDataCtrl.prototype.refresh = function () {
                    this.loadCluster();
                };
                NodeDataCtrl.prototype.loadDatasource = function (id) {
                    var _this = this;
                    return this.backendSrv.get('api/datasources')
                        .then(function (result) {
                        return lodash_1.default.filter(result, { "type": "grafana-kubernetes-datasource", "name": id })[0];
                    })
                        .then(function (ds) {
                        if (!ds) {
                            _this.alertSrv.set("Failed to connect", "Could not connect to the specified cluster.", 'error');
                            throw "Failed to connect to " + id;
                        }
                        _this.cluster = ds;
                        return _this.datasourceSrv.get(ds.name);
                    }).then(function (clusterDS) {
                        _this.clusterDS = clusterDS;
                        return clusterDS;
                    });
                };
                NodeDataCtrl.prototype.goToNodeDashboard = function (node) {
                    var _this = this;
                    var variable = lodash_1.default.find(this.templateVariables, { 'name': 'node' });
                    variable.current.text = node === 'All' ? 'All' : slugify(node.metadata.name);
                    variable.current.value = node === 'All' ? '$__all' : slugify(node.metadata.name);
                    this.variableSrv.variableUpdated(variable).then(function () {
                        _this.$scope.$emit('template-variable-value-updated');
                        _this.$scope.$root.$broadcast('refresh');
                    });
                };
                NodeDataCtrl.prototype.conditionStatus = function (condition) {
                    var status;
                    if (condition.type === "Ready") {
                        status = condition.status === "True";
                    }
                    else {
                        status = condition.status === "False";
                    }
                    return {
                        value: status,
                        text: status ? "Ok" : "Error"
                    };
                };
                NodeDataCtrl.prototype.isConditionOk = function (condition) {
                    return this.conditionStatus(condition).value;
                };
                NodeDataCtrl.prototype.conditionLastTransitionTime = function (condition) {
                    return moment_1.default(condition.lastTransitionTime).format('YYYY-MM-DD HH:mm:ss');
                };
                NodeDataCtrl.templateUrl = 'panels/nodeData/partials/node_info.html';
                NodeDataCtrl.scrollable = true;
                return NodeDataCtrl;
            })(sdk_1.PanelCtrl);
            exports_1("NodeDataCtrl", NodeDataCtrl);
        }
    }
});
//# sourceMappingURL=nodeData.js.map