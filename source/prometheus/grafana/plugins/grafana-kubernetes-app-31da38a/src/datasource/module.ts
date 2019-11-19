///<reference path="../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />

import {K8sDatasource} from './datasource';
import {K8sQueryCtrl} from './query_ctrl';

class K8sConfigCtrl {
  static templateUrl = 'datasource/partials/config.html'; 
}

export {
  K8sDatasource as Datasource,
  K8sQueryCtrl as QueryCtrl,
  K8sConfigCtrl as ConfigCtrl
};