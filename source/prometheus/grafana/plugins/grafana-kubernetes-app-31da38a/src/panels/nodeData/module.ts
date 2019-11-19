///<reference path="../../../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />

import {NodeDataCtrl} from './nodeData';
import {loadPluginCss} from 'app/plugins/sdk';

loadPluginCss({
  dark: 'plugins/grafana-kubernetes-app/css/dark.css',
  light: 'plugins/grafana-kubernetes-app/css/light.css'
});

export  {
  NodeDataCtrl as PanelCtrl
};
