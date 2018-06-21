此项目包含大文件，故需安装git lfs，参考：https://git-lfs.github.com/<br>
git clone git@github.com:yonyoucloud/install_k8s.git<br>
cd install_k8s<br>
git lfs pull # 下载大文件<br>
<br>
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
第一步, 说明:<br>
└── install_k8s        安装包<br>
    ├── fabfile.py     基于fabric实现自动化安装k8s集群脚本<br>
    ├── install.sh     安装shell脚本，里面会调用fabfile.py中函数<br>
    ├── README         说明文件<br>
    ├── source         源文件目录，也包括配置<br>
    └── ssh            直连容器需要的秘钥<br>
    └── uninstall.sh   卸载脚本<br>
    └── start.sh       如果采用LVS代理Etcd、Master集群，LVS机器发生重启时，需要在其上执行此脚本<br>
    └── add_node.sh    修改fabric.py中newnode配置，执行此脚本可以添加node节点，支持一次添加多个<br>
<br>
第二步, 修改fabfile.py文件中主机登录密码及安装目的主机地址:<br>
注意: 整个集群支持安装到一台主机上面, 需要注意vip要在同一网段, 且etcd和master的vip必须不同<br>
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
env.roledefs = {<br>
    # 次发布脚本运行主机,需要把地址加到master证书,否则后面执行kubectl认证不通过,此机可做发布机用<br>
    'publish': {<br>
        'hosts': [<br>
            '10.211.55.25:22',<br>
        ],<br>
    },<br>
    # etcd节点安装主机(支持集群)<br>
    'etcd': {<br>
        'hosts': [<br>
            '10.211.55.25:22',<br>
        ],<br>
        # 负载均衡etcd入口ip(虚ip)<br>
        'vip': '10.211.55.201'<br>
    },<br>
    # master节点安装主机(支持集群)<br>
    'master': {<br>
        'hosts': [<br>
            '10.211.55.25:22',<br>
        ],<br>
        # 负载均衡master入口ip(虚ip)<br>
        'vip': '10.211.55.202'<br>
    },<br>
    # node节点安装主机(支持集群)<br>
    'node': {<br>
        'hosts': [<br>
            '10.211.55.25:22',<br>
        ]<br>
    },<br>
    # lvs负载均衡安装主机(暂不支持集群)<br>
    'lvs': {<br>
        'hosts': [<br>
            '10.211.55.25:22',<br>
        ]<br>
    },<br>
    # 私有docker镜像仓库安装主机(暂不支持集群)<br>
    'pridocker': {<br>
        'hosts': [<br>
            '10.211.55.25:22',<br>
        ]<br>
    },<br>
    # 私有dns服务器安装主机(暂不支持集群)<br>
    'pridns': {<br>
        'hosts': [<br>
            '10.211.55.25:22',<br>
        ]<br>
    },<br>
    # 新加Node节点(支持集群)<br>
    'newnode': {<br>
        'hosts': [<br>
            '10.211.55.26:22',<br>
        ]<br>
    },<br>
<br>
第三步:<br>
只需执行install.sh文件<br>
cd install_k8s<br>
./install.sh<br>
