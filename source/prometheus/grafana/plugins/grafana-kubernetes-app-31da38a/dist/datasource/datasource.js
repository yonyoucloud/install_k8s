///<reference path="../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
System.register(['lodash'], function(exports_1) {
    var lodash_1;
    var K8sDatasource;
    function addNamespace(namespace) {
        return namespace ? 'namespaces/' + namespace + '/' : '';
    }
    function addLabels(labels) {
        var querystring = '';
        lodash_1.default.forEach(labels, function (value, label) {
            querystring += label + '%3D' + value + '%2C';
        });
        return lodash_1.default.trimEnd(querystring, '%2C');
    }
    return {
        setters:[
            function (lodash_1_1) {
                lodash_1 = lodash_1_1;
            }],
        execute: function() {
            K8sDatasource = (function () {
                function K8sDatasource(instanceSettings, backendSrv, templateSrv, $q) {
                    this.backendSrv = backendSrv;
                    this.templateSrv = templateSrv;
                    this.$q = $q;
                    this.type = instanceSettings.type;
                    this.url = instanceSettings.url;
                    this.name = instanceSettings.name;
                    this.id = instanceSettings.id;
                    this.ds = instanceSettings.jsonData.ds;
                    this.backendSrv = backendSrv;
                    this.$q = $q;
                }
                K8sDatasource.prototype.testDatasource = function () {
                    return this.backendSrv.datasourceRequest({
                        url: this.url + '/',
                        method: 'GET'
                    }).then(function (response) {
                        if (response.status === 200) {
                            return { status: "success", message: "Data source is working", title: "Success" };
                        }
                    });
                };
                K8sDatasource.prototype._get = function (apiResource) {
                    return this.backendSrv.datasourceRequest({
                        url: this.url + apiResource,
                        method: "GET",
                        headers: { 'Content-Type': 'application/json' }
                    }).then(function (response) {
                        return response.data;
                    }, function (error) {
                        return error;
                    });
                };
                K8sDatasource.prototype.getNodes = function () {
                    return this._get('/api/v1/nodes')
                        .then(function (result) {
                        return result.items;
                    });
                };
                K8sDatasource.prototype.getNode = function (name) {
                    return this._get('/api/v1/nodes/' + name);
                };
                K8sDatasource.prototype.getNamespaces = function () {
                    return this._get('/api/v1/namespaces')
                        .then(function (result) {
                        return result.items;
                    });
                };
                K8sDatasource.prototype.getComponentStatuses = function () {
                    return this._get('/api/v1/componentstatuses')
                        .then(function (result) {
                        return result.items;
                    });
                };
                K8sDatasource.prototype.getDaemonSets = function (namespace) {
                    return this._get('/apis/extensions/v1beta1/' + addNamespace(namespace) + 'daemonsets')
                        .then(function (result) {
                        return result.items;
                    });
                };
                K8sDatasource.prototype.getReplicationControllers = function (namespace) {
                    return this._get('/api/v1/' + addNamespace(namespace) + 'replicationcontrollers')
                        .then(function (result) {
                        return result.items;
                    });
                };
                K8sDatasource.prototype.getDeployments = function (namespace) {
                    return this._get('/apis/extensions/v1beta1/' + addNamespace(namespace) + 'deployments')
                        .then(function (result) {
                        return result.items;
                    });
                };
                K8sDatasource.prototype.getPods = function (namespace) {
                    return this._get('/api/v1/' + addNamespace(namespace) + 'pods')
                        .then(function (result) {
                        return result.items;
                    });
                };
                K8sDatasource.prototype.getPodsByLabel = function (namespace, labels) {
                    return this._get('/api/v1/' + addNamespace(namespace) + 'pods?labelSelector=' + addLabels(labels))
                        .then(function (result) {
                        return result.items;
                    });
                };
                K8sDatasource.prototype.getPod = function (name) {
                    return this._get('/api/v1/pods/?fieldSelector=metadata.name%3D' + name)
                        .then(function (result) {
                        if (result.items && result.items.length === 1) {
                            return result.items[0];
                        }
                        else {
                            return result.items;
                        }
                    });
                };
                K8sDatasource.prototype.getPodsByName = function (names) {
                    var _this = this;
                    var promises = [];
                    if (Array.isArray(names)) {
                        lodash_1.default.forEach(names, function (name) {
                            promises.push(_this.getPod(name));
                        });
                        return this.$q.all(promises);
                    }
                    else {
                        return this.getPod(names)
                            .then(function (pod) {
                            return [pod];
                        });
                    }
                };
                K8sDatasource.prototype.query = function (options) {
                    throw new Error("Query Support not implemented yet.");
                };
                K8sDatasource.prototype.annotationQuery = function (options) {
                    throw new Error("Annotation Support not implemented yet.");
                };
                K8sDatasource.prototype.metricFindQuery = function (query) {
                    var promises = [];
                    var namespaces;
                    if (!query) {
                        return Promise.resolve([]);
                    }
                    var interpolated = this.templateSrv.replace(query, {});
                    var query_list = interpolated.split(" ");
                    if (query_list.length > 1) {
                        namespaces = query_list[1].replace("{", "").replace("}", "").split(",");
                    }
                    else {
                        namespaces = [""]; //Gets all pods/deployments
                    }
                    switch (query_list[0]) {
                        case 'pod':
                            for (var _i = 0; _i < namespaces.length; _i++) {
                                var ns = namespaces[_i];
                                promises.push(this.getPods(ns));
                            }
                            return Promise.all(promises).then(function (res) {
                                var data = [];
                                var pods = lodash_1.default.flatten(res).filter(function (n) { return n; });
                                for (var _i = 0; _i < pods.length; _i++) {
                                    var pod = pods[_i];
                                    data.push({
                                        text: pod.metadata.name,
                                        value: pod.metadata.name,
                                    });
                                }
                                return data;
                            });
                        case 'deployment':
                            for (var _a = 0; _a < namespaces.length; _a++) {
                                var ns = namespaces[_a];
                                promises.push(this.getDeployments(ns));
                            }
                            return Promise.all(promises).then(function (res) {
                                var data = [];
                                var deployments = lodash_1.default.flatten(res).filter(function (n) { return n; });
                                for (var _i = 0; _i < deployments.length; _i++) {
                                    var deployment = deployments[_i];
                                    data.push({
                                        text: deployment.metadata.name,
                                        value: deployment.metadata.name,
                                    });
                                }
                                return data;
                            });
                        case 'namespace':
                            return this.getNamespaces().then(function (namespaces) {
                                var data = [];
                                for (var _i = 0; _i < namespaces.length; _i++) {
                                    var ns = namespaces[_i];
                                    data.push({
                                        text: ns.metadata.name,
                                        value: ns.metadata.name,
                                    });
                                }
                                ;
                                return data;
                            });
                        case 'node':
                            return this.getNodes().then(function (nodes) {
                                var data = [];
                                for (var _i = 0; _i < nodes.length; _i++) {
                                    var node = nodes[_i];
                                    data.push({
                                        text: node.metadata.name,
                                        value: node.metadata.name,
                                    });
                                }
                                ;
                                return data;
                            });
                        case 'datasource':
                            return Promise.resolve([{
                                    text: this.ds,
                                    value: this.ds,
                                }]);
                    }
                };
                K8sDatasource.baseApiUrl = '/api/v1/';
                return K8sDatasource;
            })();
            exports_1("K8sDatasource", K8sDatasource);
        }
    }
});
//# sourceMappingURL=datasource.js.map