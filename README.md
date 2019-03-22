先展示一下结果
===============
安装后最终输出下图样子即为成功:<br>
![image](https://github.com/yonyoucloud/install_k8s/blob/master/images/finish-install.jpeg)<br>
获取pods命令正常输出:<br>
![image](https://github.com/yonyoucloud/install_k8s/blob/master/images/getpods.jpeg)<br>
访问dashboard:<br>
![image](https://github.com/yonyoucloud/install_k8s/blob/master/images/dashboard.jpeg)<br>
访问grafana:<br>
![image](https://github.com/yonyoucloud/install_k8s/blob/master/images/grafana.jpeg)<br>
访问微服务例子:<br>
![image](https://github.com/yonyoucloud/install_k8s/blob/master/images/web-test.jpeg)<br>
<br>
下面是安装方法（安装过程有疑问可以加我微信：bsh888）<br>
=============
第一步，克隆安装程序，并进入目录：<br>
git clone https://github.com/yonyoucloud/install_k8s<br>
cd install_k8s<br>
<br>
第二步，下载大文件二进制包（里面是x86_64下编译好的k8s二进制文件）：<br>
注意：下面的包可选，但要对应到相应的安装源码版本。<br>
<br>
v1.13.3.1版本下载地址（ls-files-v1.13.3.1.gz，建议用此版本，对k8s和runc做了禁止kmem特性，否则可能容器中进程被内核频繁oomkill）：<br>
链接: https://pan.baidu.com/s/17qBtZ8zanKNjXX4hnx48Zw 提取码: 8krp<br>
<br>
v1.13.3版本下载地址（ls-files-v1.13.3.gz）：<br>
链接: https://pan.baidu.com/s/1nNrfjA8fFqlkFa442jW47w 提取码: y69k<br>
<br>
v1.11.3版本下载地址（ls-files-v1.11.3.gz）：<br>
链接: https://pan.baidu.com/s/1gCBY6YgG1McnDUen6egfMg 提取码: r7rs<br>
<br>
mv ls-files-v1.13.3.gz install_k8s<br>
cd install_k8s<br>
tar zxvf ls-files-v1.13.3.gz<br>
<br>
第三步，修改fabfile.py文件中主机登录密码及安装目的主机地址，并执行安装脚本：<br>
建议：开始先测试一下单机部署，只需要替换一下fabfile.py中的ip地址及登录密码即可。<br>
执行安装脚本：<br>
cd install_k8s<br>
./install.sh<br>
<br>
<br>
下面是更详细的一些说明:<br>
=============
本安装包，运行在centos7上，包含的服务有:<br>
etcd集群<br>
kubernets master集群<br>
kubernets node集群<br>
calico 网络<br>
docker 私有镜像仓库<br>
ipvsadm(lvs) 负载均衡<br>
bind 私有dns服务,方便内网域名拦截解析<br>
<br>
镜像服务包括:<br>
kube-dns<br>
kubernetes-dashboard k8s UI<br>
heapster 监控<br>
<br>
测试例子微服务(golang写的一个小的输出服务):<br>
web_test<br>
<br>
安装脚本目录说明:<br>
└── install_k8s        安装包<br>
    ├── fabfile.py     基于fabric实现自动化安装k8s集群脚本<br>
    ├── install.sh     安装shell脚本，里面会调用fabfile.py中函数<br>
    ├── README         说明文件<br>
    ├── source         源文件目录，也包括配置<br>
    └── ssh            直连容器需要的秘钥<br>
    └── uninstall.sh   卸载脚本<br>
    └── add_node.sh    修改fabric.py中newnode配置，执行此脚本可以添加node节点，支持一次添加多个，执行完把newnode合并到node配置中，便于集中控制<br>
    └── add_etcd.sh    修改fabric.py中newetcd配置，执行此脚本可以添加etcd节点，支持一次添加多个，执行完把newetcd合并到etcd配置中，便于集中控制<br>
    └── add_master.sh    修改fabric.py中newmaster配置，执行此脚本可以添加master节点，支持一次添加多个，执行完把newmaster合并到master配置中，便于集中控制<br>
<br>
fabfile.py说明:<br>
注意: <br>
1、整个集群支持安装到一台主机上面, 需要注意vip要在同一网段, 且etcd和master的vip必须不同<br>
2、如果采用LVS方式，机器重启时需要执行相应的虚ip挂载（这个不一定是必须执行）<br>
fab service_lvs_start #全部启动<br>
fab service_lvs_etcd #启动etcd<br>
fab service_lvs_master #启动master<br>
<br>
特别注意三项:<br>
1、修改脚本中的主机密码信息<br>
2、修改脚本中的主机地址信息<br>
3、确保系统是centos7，并且网卡名字是eth0<br>
<br>
编辑脚本配置信息, vim fabfile.py:<br>
env.user = 'root'<br>
env.password = '123456' # 注意这里需要修改服务器密码，集群密码要统一，也可以用下面秘钥文件的方式<br>
#env.key_filename = "~/.ssh/id_rsa"<br>
env.port = 22<br>
env.abort_on_prompts = True<br>
env.colors = True<br>
<br>
# 如果在阿里云、华为云部署等云IaaS部署，请设置为False，env.roledefs['lvs']['hosts']置为空，<br>
# 并且配置env.roledefs['etcd']['vip']及env.roledefs['master']['vip']分别为etcd、master<br>
# 负载均衡地址，并且事先将端口及虚机设置好<br>
env.use_lvs = True<br>
<br>
env.roledefs = {<br>
    # 发布机，后面通过在此机器上执行kubectl命令控制k8s集群及部署应用<br>
    'publish': {<br>
        'hosts': [<br>
            '10.211.55.53:22',<br>
        ],<br>
    },<br>
    # etcd节点安装主机(支持集群)<br>
    'etcd': {<br>
        'hosts': [<br>
            '10.211.55.54:22',<br>
            '10.211.55.55:22',<br>
        ],<br>
        # 负载均衡etcd入口ip(虚ip)<br>
        'vip': '10.211.55.201'<br>
    },<br>
    # master节点安装主机(支持集群)<br>
    'master': {<br>
        'hosts': [<br>
            '10.211.55.54:22',<br>
            '10.211.55.55:22',<br>
        ],<br>
        # 负载均衡master入口ip(虚ip)<br>
        'vip': '10.211.55.202'<br>
    },<br>
    # node节点安装主机(支持集群)<br>
    'node': {<br>
        'hosts': [<br>
            '10.211.55.54:22',<br>
            '10.211.55.55:22',<br>
        ]<br>
    },<br>
    # lvs负载均衡安装主机(暂不支持集群)<br>
    # 特别要注意，如果etcd及master是多机部署，lvs上不要放etcd及master服务，且不要和发布机在一起，否则网络会有问题，如果是阿里云、华为云一定要换成对应的slb（需要提前配置好节点及端口），其实最好lvs单独部署，因为在其上面是无法访问其负载均衡的节点的，为了节省资源，上面可以放私有镜像仓库、私有dns服务<br>
    'lvs': {<br>
        'hosts': [<br>
            '10.211.55.56:22',<br>
        ]<br>
    },<br>
    # 私有docker镜像仓库安装主机(暂不支持集群)<br>
    'pridocker': {<br>
        'hosts': [<br>
            '10.211.55.56:22',<br>
        ]<br>
    },<br>
    # 私有dns服务器安装主机(暂不支持集群)<br>
    'pridns': {<br>
        'hosts': [<br>
            '10.211.55.53:22',<br>
        ]<br>
    },<br>
    # 新加Node节点(支持集群)<br>
    'newnode': {<br>
        'hosts': [<br>
            #'10.211.55.57:22',<br>
        ]<br>
    },<br>
    # 新加etcd节点(支持集群)<br>
    'newetcd': {<br>
        'hosts': [<br>
            #'10.211.55.58:22',<br>
        ]<br>
    },<br>
    # 新加master节点(支持集群)<br>
    'newmaster': {<br>
        'hosts': [<br>
            #'10.211.55.59:22',<br>
        ]<br>
    },<br>
}<br>
