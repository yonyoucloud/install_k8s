///<reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />

import kbn from 'app/core/utils/kbn';
import _ from 'lodash';
import moment from 'moment';

export default class NodeStatsDatasource {
  constructor(private datasourceSrv, private timeSrv) {}

  issuePrometheusQuery(prometheusDS, query) {
    return this.datasourceSrv.get(prometheusDS)
      .then((datasource) => {
        var metricsQuery = {
          range: { from: moment().subtract(5, 'minute'), to: moment() },
          targets: [{ expr: query.expr, format: 'time_series' }],
          legendFormat: query.legend,
          interval: '60s',
        };
        return datasource.query(metricsQuery);
      }).then((result) => {
        if (result && result.data) {
          return result.data;
        }
        return {};
      });
  }

  getNodeStats(cluster_id, prometheusDS) {
    let podsPerNode, cpuPerNode, memoryPerNode;

    const podQuery = {
      expr: 'sum(label_join(kubelet_running_pod_count, "node",  "", "kubernetes_io_hostname")) by (node)',
      legend: "{{node}}",
    };
    const cpuQuery = {
      expr: 'sum(kube_pod_container_resource_requests_cpu_cores) by (node)',
      legend: "{{node}}",
    };
    const memoryQuery = {
      expr: 'sum(kube_pod_container_resource_requests_memory_bytes) by (node)',
      legend: "{{node}}",
    };

    return this.issuePrometheusQuery(prometheusDS, podQuery)
      .then(data => {
        podsPerNode = data;
        return;
      }).then(() => {
        return this.issuePrometheusQuery(prometheusDS, cpuQuery);
      })
      .then(data => {
        cpuPerNode = data;
        return;
      }).then(() => {
        return this.issuePrometheusQuery(prometheusDS, memoryQuery);
      })
      .then(data => {
        memoryPerNode = data;
        return {podsPerNode, cpuPerNode, memoryPerNode};
      });
  }

  updateNodeWithStats(node, nodeStats) {
    var formatFunc = kbn.valueFormats['percentunit'];
    const nodeName = slugify(node.metadata.name);
    const findFunction = function(o) {return o.target.substring(7, o.target.length - 2) === nodeName;};
    const podsUsedData = _.find(nodeStats.podsPerNode, findFunction);
    if (podsUsedData) {
      node.podsUsed = _.last(podsUsedData.datapoints)[0];
      node.podsUsedPerc = formatFunc(node.podsUsed / node.status.capacity.pods, 2, 5);
    }

    const cpuData = _.find(nodeStats.cpuPerNode, findFunction);
    if (cpuData) {
      node.cpuUsage = _.last(cpuData.datapoints)[0];
      node.cpuUsageFormatted = kbn.valueFormats['none'](node.cpuUsage, 2, null);
      node.cpuUsagePerc = formatFunc(node.cpuUsage / node.status.capacity.cpu, 2, 5);
    }

    const memData = _.find(nodeStats.memoryPerNode, findFunction);
    if (memData) {
      node.memoryUsage = _.last(memData.datapoints)[0];
      const memCapacity = node.status.capacity.memory.substring(0, node.status.capacity.memory.length - 2)  * 1000;
      node.memUsageFormatted = kbn.valueFormats['bytes'](node.memoryUsage, 2, null);
      node.memCapacityFormatted = kbn.valueFormats['bytes'](memCapacity, 2, null);
      node.memoryUsagePerc = formatFunc((node.memoryUsage / memCapacity), 2, 5);
    }

    return node;
  }
}

function slugify(str) {
  var slug = str.replace("@", "at").replace("&", "and").replace(/[.]/g, "_").replace("/\W+/", "");
  return slug;
}
