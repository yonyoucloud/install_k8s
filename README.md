本安装包，运行在centos7上，包含的服务有:
etcd集群
kubernets master集群
kubernets node集群
flannel 网络
docker 私有镜像仓库
ipvsadm(lvs) 负载均衡
bind 私有dns服务,方便内网域名拦截解析

镜像服务包括:
kube-dns
kubernetes-dashboard k8s UI
heapster 监控

第一步, 说明:
└── install_k8s        安装包
    ├── fabfile.py     基于fabric实现自动化安装k8s集群脚本
    ├── install.sh     安装shell脚本，里面会调用fabfile.py中函数
    ├── README         说明文件
    ├── source         源文件目录，也包括配置
    └── ssh            直连容器需要的秘钥
    └── uninstall.sh   卸载脚本

第二步, 修改fabfile.py文件中主机登录密码及安装目的主机地址:
注意: 整个集群支持安装到一台主机上面, 需要注意vip要在同一网段, 且etcd和master的vip必须不同

特别注意三项:
1、修改脚本中的主机密码信息
2、修改脚本中的主机地址信息
3、确保系统是centos7，并且网口名字是eth0

编辑脚本配置信息, vim fabfile.py:
env.user = 'root'
env.password = '123456' # 注意这里需要修改服务器密码，集群密码要统一，也可以用下面秘钥文件的方式
#env.key_filename = "~/.ssh/id_rsa"
env.port = 22
env.abort_on_prompts = True
env.colors = True

env.roledefs = {
    # 次发布脚本运行主机,需要把地址加到master证书,否则后面执行kubectl认证不通过,此机可做发布机用
    'publish': {
        'hosts': [
            '10.211.55.23:22',
        ],
    },
    # etcd节点安装主机(支持集群)
    'etcd': {
        'hosts': [
            '10.211.55.23:22',
        ],
        # 负载均衡etcd入口ip(虚ip)
        'vip': '10.211.55.201'
    },
    # master节点安装主机(支持集群)
    'master': {
        'hosts': [
            '10.211.55.23:22',
        ],
        # 负载均衡master入口ip(虚ip)
        'vip': '10.211.55.202'
    },
    # node节点安装主机(支持集群)
    'node': {
        'hosts': [
            '10.211.55.23:22',
        ]
    },
    # lvs负载均衡安装主机(暂不支持集群)
    'lvs': {
        'hosts': [
            '10.211.55.23:22',
        ]
    },
    # 私有docker镜像仓库安装主机(暂不支持集群)
    'pridocker': {
        'hosts': [
            '10.211.55.23:22',
        ]
    },
    # 私有dns服务器安装主机(暂不支持集群)
    'pridns': {
        'hosts': [
            '10.211.55.23:22',
        ]
    },

第三步:
只需执行install.sh文件
cd install_k8s
./install.sh
