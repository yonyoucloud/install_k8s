安装最新稳定版：
===============
#### 版本说明

| service              | version |
|----------------------|---------|
| `kubernetes`         | v1.32.3 |
| `etcd`               | 3.5.21  |
| `istio`              | 1.25.1  |
| `containerd`         | v1.7.27 |
| `calico`             | v3.29.3 |
| `coredns`            | v1.12.0 |
| `dashboard`          | v2.7.0  |
| `metrics-server`     | v0.7.2  |
| `prometheus`         | v3.2.1  |
| `alertmanager`       | v0.28.1 |
| `grafana`            | 11.6.0  |
| `kube-state-metrics` | v2.15.0 |
| `node-exporter`      | v1.9.0  |
| `helm`               | v3.17.2 |
| `cfssl`              | v1.6.5  |
| `runc`               | v1.2.6  |
| `crictl`             | v1.32.0 |
| `nerdctl`            | v2.0.4  |
| `cni`                | v1.6.2  |

#### 安装说明

##### 1、下载并解压缩二进制安装文件(installk8s-v1.32.3-20250401.gz):
[v1.32.3版本静态安装包下载](https://pan.baidu.com/s/1OuwYBsirXzZL-kJJdBrocQ?pwd=dznm)

[v1.31.4版本静态安装包下载](https://pan.baidu.com/s/1TOqKINTELq2qpxHfOSJuzw?pwd=ap6v)

[v1.26.4版本静态安装包下载](https://pan.baidu.com/s/1Q5XaSDyCKzkT_mtJcOV5dA?pwd=mv7n)

[v1.23.6版本静态安装包下载](https://pan.baidu.com/s/1kJ4vc9yMrskW-UyXXZ2Hng?pwd=hc3a)
```
cd /data/ && tar zxvf installk8s-v1.32.3-20250401.gz -C /data

cd /data/installk8s/sysbase
/data/installk8s/sysbase
├── bin # 根据系统执行可执行文件，启动安装服务
│   ├── sysbase-v1.0.0-darwin-arm64
│   └── sysbase-v1.0.0-linux-amd64
├── etc
│   ├── config-demo.yaml
│   ├── config.js
│   └── config.yaml # 修改此配置文件，注意数据库配置及二进制安装文件路径（先创建 sysbase 数据库，程序启动会自动创建表）
└── static
    ├── config.js # 需要修改一下 apiHost，写成运行此程序的机器 ip，端口号保持和此服务运行端口一致
    ├── css
    │   ├── chunk-vendors.537be47b.css
    │   └── index.d169f2ef.css
    ├── favicon.ico
    ├── fonts
    │   ├── element-icons.535877f5.woff
    │   └── element-icons.732389de.ttf
    ├── index.html
    └── js
        ├── about.3c68e217.js
        ├── about.3c68e217.js.map
        ├── chunk-vendors.ccefee03.js
        ├── chunk-vendors.ccefee03.js.map
        ├── index.0ab2cb87.js
        └── index.0ab2cb87.js.map
```

##### 2、运行安装服务：
```
cd /data/installk8s/sysbase
./bin/sysbase-v1.0.0-linux-amd64 (这里不同平台选择不同可执行文件) 
访问安装服务: http://192.168.58.2:8081/static/ (这里的 192.168.58.2 根据实际情况，是运行安装服务的 IP 地址)
```

##### 3、添加资源:
```
资源类型选择 vps，特定描述这些都是必须的，其中 etcd、master、node 可以配置多台机器，其他几个确保唯一， 另外，
支持这些全配置一台机器，即单机也可以运行，建议单台最低配置 8 核 16G，其实 4 核 8G 也可以运行起来。
```
![image](/images/install_resource.jpeg)


##### 4、创建 K8sCluster:
```
一条记录代表一个 k8s 集群，一定要选择前面添加的资源列表
```
![image](/images/install_k8scluster.jpeg)

##### 5、执行安装:
```
点击一键安装前，可以先点击内核升级，因为内核升级会重启机器，一键安装逻辑也会判断内核是否已升级，如果 未升级，也会触发升级、重启。
```
![image](/images/install_install.jpeg)

##### 6、安装后重要目录说明:
```
/data/installk8s/addons/certs 安装过程中生成的 TLS 证书文件，可以将 k8s.com.crt 导入到系统并信任。
/data/installk8s/addons/gateways Istio 网关及虚拟服务设置，这里是站点入口配置处。

绑定Hosts，访问以下站点：
192.168.58.2 dashboard.k8s.com grafana.k8s.com prometheus.k8s.com kiali.k8s.com test.k8s.com

安装程序开源目录地址：
./sysbase
采用golang+vue编写
```

##### 7、安装后效果图:
![image](/images/instsall_example1.jpeg)
![image](/images/instsall_example2.jpeg)
![image](/images/instsall_example3.jpeg)
![image](/images/instsall_example4.jpeg)
![image](/images/instsall_example5.jpeg)
![image](/images/instsall_example6.jpeg)
![image](/images/instsall_example7.jpeg)
