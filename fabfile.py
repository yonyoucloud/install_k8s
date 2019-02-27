#!/usr/bin/python
#coding:utf-8
# -------------------------------------------------------------------------------
# Filename:    fabfile.py
# Revision:    2.0
# Date:        2018/06/21
# Author:      bishenghua
# Email:       net.bsh@gmail.com
# Description: Script to install the kubernets system
# -------------------------------------------------------------------------------
# Copyright:   2018 (c) Bishenghua
# License:     GPL
#
# This program is free software; you can redistribute it and/or
# modify it under the terms of the GNU General Public License
# as published by the Free Software Foundation; either version 2
# of the License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty
# of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU General Public License for more details.
#
# you should have received a copy of the GNU General Public License
# along with this program (or with Nagios);
#
# Credits go to Ethan Galstad for coding Nagios
# If any changes are made to this script, please mail me a copy of the changes
# -------------------------------------------------------------------------------

from __future__ import with_statement

import sys
import subprocess
import time
import os

from fabric.api import env
from fabric.api import run
from fabric.api import parallel
from fabric.api import roles
from fabric.api import execute
from fabric.api import local
from fabric.api import get
from fabric.api import cd
from fabric.api import put

from fabric.api import hide
from fabric.api import settings

from fabric.api import task

env.user = 'root'
env.password = '123456'
#env.key_filename = "~/.ssh/id_rsa"
env.port = 22
env.abort_on_prompts = True
env.colors = True

# 如果在阿里云、华为云部署等云IaaS部署，请设置为False，env.roledefs['lvs']['hosts']置为空，
# 并且配置env.roledefs['etcd']['vip']及env.roledefs['master']['vip']分别为etcd、master
# 负载均衡地址，并且事先将端口及虚机设置好
env.use_lvs = False

env.roledefs = {
    # 发布机，后面通过在此机器上执行kubectl命令控制k8s集群及部署应用
    'publish': {
        'hosts': [
            '10.211.55.12:22',
        ],
    },
    # etcd节点安装主机(支持集群)
    'etcd': {
        'hosts': [
            '10.211.55.12:22',
        ],
        # 负载均衡etcd入口ip(虚ip)
        'vip': '10.211.55.12'
    },
    # master节点安装主机(支持集群)
    'master': {
        'hosts': [
            '10.211.55.12:22',
        ],
        # 负载均衡master入口ip(虚ip)
        'vip': '10.211.55.12'
    },
    # node节点安装主机(支持集群)
    'node': {
        'hosts': [
            '10.211.55.12:22',
        ]
    },
    # lvs负载均衡安装主机(暂不支持集群)
    # 特别要注意，如果etcd及master是多机部署，lvs上不要放etcd及master服务，且不要和发布机在一起，否则网络会有问题，如果是阿里云、华为云一定要换成对应的slb（需要提前配置好节点及端口），其实最好lvs单独部署，因为在其上面是无法访问其负载均衡的节点的，为了节省资源，上面可以放私有镜像仓库、私有dns服务
    'lvs': {
        'hosts': [
            '10.211.55.12:22',
        ]
    },
    # 私有docker镜像仓库安装主机(暂不支持集群)
    'pridocker': {
        'hosts': [
            '10.211.55.12:22',
        ]
    },
    # 私有dns服务器安装主机(暂不支持集群)
    'pridns': {
        'hosts': [
            '10.211.55.12:22',
        ]
    },
    # 新加Node节点(支持集群)
    'newnode': {
        'hosts': [
            #'10.211.55.57:22',
        ]
    },
    # 新加etcd节点(支持集群)
    'newetcd': {
        'hosts': [
            #'10.211.55.58:22',
        ]
    },
    # 新加master节点(支持集群)
    'newmaster': {
        'hosts': [
            #'10.211.55.59:22',
        ]
    },
}

def exec_cmd(cmd):
    p = subprocess.Popen(cmd, shell = True, stdout = subprocess.PIPE, stderr = subprocess.STDOUT)
    for line in p.stdout.readlines():
        print line.strip()
    retval = p.wait()


##########################[启动服务]############################
def service(dowhat = 'start'):
    execute(service_etcd, dowhat)
    execute(service_master, dowhat)
    execute(service_node, dowhat)
    execute(service_dns, dowhat)
    execute(service_publish, dowhat)
##########################[启动服务]############################


##########################[publish控制]############################
@parallel
@roles('publish')
def service_publish(dowhat = 'start'):
    run('systemctl ' + dowhat + ' docker')
    pass
##########################[publish控制]############################


##########################[etcd控制]############################
@parallel
@roles('etcd')
#fab service_etcd:status
def service_etcd(dowhat = 'start'):
    execute(_service_etcd, dowhat)
    pass

def _service_etcd(dowhat = 'start'):
    etcdlvs = env.roledefs['etcd']['vip']
    run('systemctl ' + dowhat + ' etcd')
    if dowhat == 'start' or dowhat == 'restart':
        run('iptables -P FORWARD ACCEPT ; echo "" > /dev/null')
        #local('etcdctl --ca-file=source/etcd/etc/etcd/ssl/ca.pem --cert-file=source/etcd/etc/etcd/ssl/etcd.pem --key-file=source/etcd/etc/etcd/ssl/etcd-key.pem --endpoints=https://' + etcdlvs + ':2379 set /esn.com/network/config \'{"Network":"172.30.0.0/16","SubnetLen":25,"Backend":{"Type":"vxlan"}}\'')
    pass

@parallel
@roles('newetcd')
def newetcd_service_etcd_start():
    if not len(env.roledefs['newetcd']['hosts']):
        return
    execute(_service_etcd, 'start')
    pass
##########################[etcd控制]############################


##########################[master控制]############################
@parallel
@roles('master')
def service_master(dowhat = 'start'):
    execute(_service_master, dowhat)
    pass

def _service_master(dowhat = 'start'):
    run('systemctl ' + dowhat + ' kube-apiserver')
    run('systemctl ' + dowhat + ' kube-controller-manager')
    run('systemctl ' + dowhat + ' kube-scheduler')
    if dowhat == 'start' or dowhat == 'restart':
        run('iptables -P FORWARD ACCEPT ; echo "" > /dev/null')
    if dowhat == 'stop':
        run("ps aux | grep kube-apiserver | grep -v grep | awk '{if($2 != \"\"){system(\"kill -9 \"$2)}}'")
        run("ps aux | grep kube-controller-manager | grep -v grep | awk '{if($2 != \"\"){system(\"kill -9 \"$2)}}'")
        run("ps aux | grep kube-scheduler | grep -v grep | awk '{if($2 != \"\"){system(\"kill -9 \"$2)}}'")
    pass

@parallel
@roles('newmaster')
def newmaster_service_master_start():
    if not len(env.roledefs['newmaster']['hosts']):
        return
    execute(_service_master, 'start')
    pass
##########################[master控制]############################


##########################[node控制]############################
@parallel
@roles('node')
def service_node(dowhat = 'start', docker = True):
    docker = docker == str(True)
    execute(_service_node, dowhat, docker)
    pass

@parallel
@roles('node', 'newnode')
def newmaster_service_node_restart():
    execute(_service_node, 'restart', False)
    pass

def newnode_service_node_start():
    execute(_newnode_service_node_start)

    i = 0
    while True:
        i = i + 1
        hosts = ''
        split = ''
        for host in env.roledefs['newnode']['hosts']:
            hosts += split + host.split(':')[0]
            split = '|'
        num = local('kubectl get nodes | grep -E "' + hosts + '" | grep Ready | wc -l', capture = True)
        total = len(env.roledefs['newnode']['hosts'])
        print '等待所有节点运行状态变为Ready(%ds)(%d = %s)' % (i, total, num)
        if int(num) == total:
            break
        time.sleep(3)

    i = 0
    while True:
        i = i + 1
        num = local('kubectl get pods -o wide -n kube-system | grep -E "' + hosts + '" | grep calico-node | grep Running | wc -l', capture = True)
        total = len(env.roledefs['newnode']['hosts'])
        print '等待所有节点calico-node容器正常运行(%ds)(%d = %s)' % (i, total, num)
        if int(num) == total:
            break
        time.sleep(3)
    pass

@parallel
@roles('newnode')
def _newnode_service_node_start():
    if not len(env.roledefs['newnode']['hosts']):
        return
    execute(_service_node, 'start')
    pass

def _service_node(dowhat = 'start', docker = True):
    run('systemctl ' + dowhat + ' kubelet')
    run('systemctl ' + dowhat + ' kube-proxy')
    if docker:
        run('systemctl ' + dowhat + ' docker')
    if dowhat == 'start' or dowhat == 'restart':
        run('iptables -P FORWARD ACCEPT ; echo "" > /dev/null')
    if dowhat == 'stop':
        run("ps aux | grep kubelet | grep -v grep | awk '{if($2 != \"\"){system(\"kill -9 \"$2)}}'")
        run("ps aux | grep kube-proxy | grep -v grep | awk '{if($2 != \"\"){system(\"kill -9 \"$2)}}'")
        if docker:
            run("ps aux | grep docker | grep -v grep | awk '{if($2 != \"\"){system(\"kill -9 \"$2)}}'")
    pass
##########################[node控制]############################


##########################[lvs控制]############################
def service_lvs_start(dowhat = 'start'):
    if not env.use_lvs:
        return
    execute('service_lvs')
    execute('service_lvs_etcd')
    execute('service_lvs_master')
    pass

@parallel
@roles('lvs')
def service_lvs():
    if not env.use_lvs:
        return
    run('/etc/rc.d/init.d/vip_route_lvs.sh')
    pass

@parallel
@roles('etcd', 'newetcd')
def service_lvs_etcd():
    if not env.use_lvs:
        return
    run('/etc/rc.d/init.d/vip_route_etcd.sh')
    pass

@parallel
@roles('master', 'newmaster')
def service_lvs_master():
    if not env.use_lvs:
        return
    run('/etc/rc.d/init.d/vip_route_master.sh')
    pass
##########################[lvs控制]############################


##########################[dns控制]############################
@parallel
@roles('pridns')
def service_dns(dowhat = 'start'):
    run('systemctl ' + dowhat + ' named-chroot')
    if dowhat == 'start' or dowhat == 'restart':
        run('iptables -P FORWARD ACCEPT ; echo "" > /dev/null')
    pass
##########################[dns控制]############################


##########################[基础安装]############################
@parallel
@roles('etcd', 'master', 'node', 'pridocker', 'lvs')
def install_base():
    execute(_install_base)
    pass

@parallel
@roles('publish')
def install_needbin():
    put('source/needbin.gz', '/tmp', mode=0640)
    run('tar zxvf /tmp/needbin.gz -C / && rm -rf /tmp/needbin.gz && mkdir -p /etc/calico')
    etcdlvs = env.roledefs['etcd']['vip']
    calicoctl_conf = '''apiVersion: projectcalico.org/v3
kind: CalicoAPIConfig
metadata:
spec:
  etcdEndpoints: "https://{0}:2379"
  etcdKeyFile: "/etc/cni/net.d/calico-tls/etcd-key"
  etcdCertFile: "/etc/cni/net.d/calico-tls/etcd-cert"
  etcdCACertFile: "/etc/cni/net.d/calico-tls/etcd-ca"
'''
    run("echo '" + calicoctl_conf.format(etcdlvs) + "' > /etc/calico/calicoctl.cfg")
    pass

@parallel
@roles('newnode')
def newnode_install_base():
    if not len(env.roledefs['newnode']['hosts']):
        return
    execute(_install_base)
    pass

@parallel
@roles('newetcd')
def newetcd_install_base():
    if not len(env.roledefs['newetcd']['hosts']):
        return
    execute(_install_base)
    pass

@parallel
@roles('newmaster')
def newmaster_install_base():
    if not len(env.roledefs['newmaster']['hosts']):
        return
    execute(_install_base)
    pass

def _install_base():
    run('yum install -y telnet net-tools openssl')
    run('mkdir /data > /dev/null 2>&1;if [ $? == 0 ];then useradd -d /data/www esn && useradd -d /data/www www && usermod -G esn www && chmod 750 /data/www && mkdir -p /data/log/php && mkdir -p /data/log/nginx && mkdir -p /data/yy_log && chown -R www:www /data/log /data/yy_log && chmod 750 /data/log /data/yy_log;fi')
    run('systemctl stop firewalld && systemctl disable firewalld ; echo "" > /dev/null')
    run('sed -i "s#SELINUX=enforcing#SELINUX=disabled#g" /etc/selinux/config && setenforce 0 ; echo "" > /dev/null')
    #run('sed -i "s#umask 022#umask 027#g" /etc/profile')
    run('cat /etc/sysctl.conf | grep net.ipv4.ip_forward > /dev/null 2>&1 ; if [ $? -ne 0 ];then echo "net.ipv4.ip_forward = 1" >> /etc/sysctl.conf && sysctl -p;fi')
    run('cat /etc/sysctl.conf | grep net.ipv4.conf.all.rp_filter > /dev/null 2>&1 ; if [ $? -ne 0 ];then echo "net.ipv4.conf.all.rp_filter = 1" >> /etc/sysctl.conf && sysctl -p;fi')
    pass
##########################[基础安装]############################


##########################[安装docker]############################
def install_docker():
    local('cd source/docker && tar zcvf conf_bin.gz etc usr')
    execute(_install_docker)
    pass

@parallel
@roles('publish', 'pridocker', 'master', 'node')
def _install_docker():
    execute(__install_docker)
    pass

@parallel
@roles('newnode')
def newnode_install_docker():
    if not len(env.roledefs['newnode']['hosts']):
        return
    execute(__install_docker)
    pass

def __install_docker():
    #put('source/docker/docker_engine_packages.gz', '/tmp', mode=0640)
    #run('cd /tmp && tar zxvf docker_engine_packages.gz && cd docker_engine_packages && yum -y localinstall * && rm -rf /tmp/docker_engine_packages.gz /tmp/docker_engine_packages')
    put('source/docker/conf_bin.gz', '/tmp', mode=0640)
    run('tar zxvf /tmp/conf_bin.gz -C / && rm -rf /tmp/conf_bin.gz && mkdir -p /data/docker && systemctl daemon-reload && systemctl enable docker')
    pass

@parallel
@roles('publish', 'pridocker', 'master', 'node')
def uninstall_docker():
    run('systemctl disable docker ; echo "" > /dev/null')
    #run('yum remove -y docker-engine')
    run('rm -rf /usr/bin/ctr /usr/bin/docker /usr/bin/dockerd /usr/bin/docker-init /usr/bin/docker-proxy /usr/bin/runc')
    run('rm -rf /data/docker /etc/docker')
    pass
##########################[安装docker]############################


##########################[安装lvs]############################
def install_lvs():
    if not env.use_lvs:
        return
    execute('remote_install_lvs')
    execute('install_lvsvip_etcd')
    execute('install_lvsvip_master')

@roles('lvs')
def remote_install_lvs():
    if not env.use_lvs:
        return
    run('yum install -y ipvsadm && systemctl enable ipvsadm')

    etcdvip = env.roledefs['etcd']['vip']
    mastervip = env.roledefs['master']['vip']

    while True:
        cmd = 'ifconfig eth0:lvs:0 ' + etcdvip + ' broadcast ' + etcdvip + ' netmask 255.255.255.255 up'
        run(cmd + ' && echo -e "#/bin/sh\\n# chkconfig:   2345 90 10\\n' + cmd + '" > /etc/rc.d/init.d/vip_route_lvs.sh')
        cmd = 'ifconfig eth0:lvs:1 ' + mastervip + ' broadcast ' + mastervip + ' netmask 255.255.255.255 up'
        run(cmd + ' && echo "' + cmd + '" >> /etc/rc.d/init.d/vip_route_lvs.sh')
        cmd = 'route add -host ' + etcdvip + ' dev eth0:lvs:0 ; echo "" > /dev/null'
        run(cmd + ' && echo "' + cmd + '" >> /etc/rc.d/init.d/vip_route_lvs.sh')
        cmd = 'route add -host ' + mastervip + ' dev eth0:lvs:1 ; echo "" > /dev/null'
        run(cmd + ' && echo "' + cmd + '" >> /etc/rc.d/init.d/vip_route_lvs.sh')
        run('chmod +x /etc/rc.d/init.d/vip_route_lvs.sh && chkconfig --add vip_route_lvs.sh && chkconfig vip_route_lvs.sh on')
        run('echo "1" > /proc/sys/net/ipv4/ip_forward')

        # etcd
        ipvsadm = '-A -t ' + etcdvip + ':2379 -s wrr\n'
        for host in env.roledefs['etcd']['hosts']:
            ipvsadm += '-a -t ' + etcdvip + ':2379 -r ' + host.split(':')[0] + ':2379 -g -w 1\n'

        # master
        ipvsadm += '-A -t ' + mastervip + ':6443 -s wrr\n'
        for host in env.roledefs['master']['hosts']:
            ipvsadm += '-a -t ' + mastervip + ':6443 -r ' + host.split(':')[0] + ':6443 -g -w 1\n'

        run('echo "' + ipvsadm + '" > /etc/sysconfig/ipvsadm')
        with settings(warn_only = True):
            result = run('systemctl restart ipvsadm && ipvsadm -Ln')
            if result.return_code == 0:
                break
    pass

@roles('lvs')
def remote_install_lvs_new():
    if not env.use_lvs:
        return
    etcdvip = env.roledefs['etcd']['vip']
    mastervip = env.roledefs['master']['vip']

    run('systemctl stop ipvsadm')

    # etcd
    ipvsadm = '-A -t ' + etcdvip + ':2379 -s wrr\n'
    etcd_hosts = env.roledefs['etcd']['hosts'] + env.roledefs['newetcd']['hosts']
    for index, host in enumerate(etcd_hosts):
        ipvsadm += '-a -t ' + etcdvip + ':2379 -r ' + host.split(':')[0] + ':2379 -g -w 1\n'

    # master
    ipvsadm += '-A -t ' + mastervip + ':6443 -s wrr\n'
    master_hosts = env.roledefs['master']['hosts'] + env.roledefs['newmaster']['hosts']
    for index, host in enumerate(master_hosts):
        ipvsadm += '-a -t ' + mastervip + ':6443 -r ' + host.split(':')[0] + ':6443 -g -w 1\n'

    run('echo "' + ipvsadm + '" > /etc/sysconfig/ipvsadm')
    with settings(warn_only = True):
        result = run('systemctl start ipvsadm && ipvsadm -Ln')
        if result.return_code == 0:
            print 'ipvsadm reload ok'
        else:
            print 'ipvsadm reload fail'
    pass

def uninstall_lvs():
    if not env.use_lvs:
        return
    execute('remote_uninstall_lvs')
    execute('uninstall_lvsvip_etcd')
    execute('uninstall_lvsvip_master')

@roles('lvs')
def remote_uninstall_lvs():
    if not env.use_lvs:
        return
    run('ifconfig eth0:lvs:0 down ; echo "" > /dev/null')
    run('ifconfig eth0:lvs:1 down ; echo "" > /dev/null')
    run('chkconfig vip_route_lvs.sh off && chkconfig --del vip_route_lvs.sh ; echo "" > /dev/null')
    run('yum remove -y ipvsadm && rm -rf /etc/sysconfig/ipvsadm /etc/rc.d/init.d/vip_route_lvs.sh')
    pass
##########################[安装lvs]############################


##########################[安装私有docker仓库]############################
@roles('pridocker')
def install_pridocker():
    curhost = env.host_string.split(':')[0]
    local('cd source/docker && sed "s#HOST#' + curhost + '#g" create_ssl.sh.tpl > create_ssl.sh && chmod 750 create_ssl.sh')
    local('cd source/docker && sed "s#HOST#' + curhost + '#g" start_registry.sh.tpl > start_registry.sh && chmod 750 start_registry.sh')
    local('cd source/docker && rm -rf ca.crt')

    put('source/docker/create_ssl.sh', '/tmp', mode=0750)
    run('/tmp/create_ssl.sh && rm -rf /tmp/create_ssl.sh')
    get('/etc/certs/' + curhost + '.crt', 'source/docker/ca.crt')
    run('systemctl restart docker')
    put('source/images/registry.tar', '/tmp', mode=0640)
    put('source/docker/start_registry.sh', '/tmp', mode=0750)
    run('/tmp/start_registry.sh ; echo "" > /dev/null && rm -rf /tmp/start_registry.sh')

    local('cd source/docker && chmod 640 ca.crt')
    pass

@roles('pridocker')
def uninstall_pridocker():
    run('systemctl stop docker && systemctl disable docker ; echo "" > /dev/null')
    run('yum remove -y docker-engine')
    run('rm -rf /data/docker /etc/docker')
    pass
##########################[安装私有docker仓库]############################


##########################[安装etcd]############################
etcd_index = 0
def install_etcd():
    # 证书要保证一样，所以只需要生成一次
    execute('create_ssl_etcd')
    execute('remote_install_etcd')
    pass

def newetcd_install_etcd():
    execute('create_ssl_etcd', True)
    execute('newetcd_remote_install_etcd')
    execute('newetcd_remote_modify_etcd_conf')
    execute('remote_update_master_etcd_ssl')
    pass

@roles('etcd')
def remote_install_etcd():
    execute(_remote_install_etcd)
    pass

@roles('etcd', 'newetcd')
def newetcd_remote_install_etcd():
    execute(_remote_install_etcd, True)
    pass

@roles('newetcd')
def newetcd_remote_modify_etcd_conf():
    if not len(env.roledefs['newetcd']['hosts']):
        return
    run('sed -i \'s#ETCD_INITIAL_CLUSTER_STATE="new"#ETCD_INITIAL_CLUSTER_STATE="existing"#g\' /etc/etcd/etcd.conf')
    pass

def newetcd_cluster_addnew():
    global etcd_index
    etcd_index = len(env.roledefs['etcd']['hosts'])
    etcdlvs = env.roledefs['etcd']['vip']
    for index, host in enumerate(env.roledefs['newetcd']['hosts']):
        etcd_index += 1
        local('etcdctl --ca-file=source/etcd/etc/etcd/ssl/ca.pem --cert-file=source/etcd/etc/etcd/ssl/etcd.pem --key-file=source/etcd/etc/etcd/ssl/etcd-key.pem --endpoints=https://' + etcdlvs + ':2379 member add etcd' + str(etcd_index) + ' https://' + host.split(':')[0]  + ':2380')
    pass

def newetcd_cluster_addnew_check():
    etcdlvs = env.roledefs['etcd']['vip']
    local('etcdctl --ca-file=source/etcd/etc/etcd/ssl/ca.pem --cert-file=source/etcd/etc/etcd/ssl/etcd.pem --key-file=source/etcd/etc/etcd/ssl/etcd-key.pem --endpoints=https://' + etcdlvs + ':2379 member list')
    local('etcdctl --ca-file=source/etcd/etc/etcd/ssl/ca.pem --cert-file=source/etcd/etc/etcd/ssl/etcd.pem --key-file=source/etcd/etc/etcd/ssl/etcd-key.pem --endpoints=https://' + etcdlvs + ':2379 cluster-health')
    pass

def _remote_install_etcd(addnew = False):
    global etcd_index
    curhost = env.host_string.split(':')[0]
    #if env.roledefs['etcd'].has_key('lvs'):
    cluster_hosts = ''
    tmpstr = ''
    etcd_index += 1
    etcdname = 'etcd' + str(etcd_index)
    etcd_hosts = env.roledefs['etcd']['hosts']
    if addnew:
        etcd_hosts = env.roledefs['etcd']['hosts'] + env.roledefs['newetcd']['hosts']
    for index, host in enumerate(etcd_hosts):
        cluster_hosts += tmpstr + ('etcd' + str(index + 1)) + '=https://' + host.split(':')[0]  + ':2380'
        tmpstr = ','

    local('cd source/etcd && sed "s#CLUSTER_HOSTS#' + cluster_hosts + '#g" etcd.conf.tpl > etc/etcd/etcd.conf')
    local('cd source/etcd && sed -i "s#ETCD_HOST#' + curhost + '#g" etc/etcd/etcd.conf')
    local('cd source/etcd && sed -i "s#CONFIG_ETCD_NAME#' + etcdname + '#g" etc/etcd/etcd.conf')

    local('cd source/etcd && tar zcvf etcd.gz etc usr')

    # run('yum -y install etcd')
    put('source/etcd/etcd.gz', '/etcd.gz', mode=0640)
    run('id etcd >& /dev/null || useradd -c "etcd user" -s /sbin/nologin -d /var/lib/etcd -r etcd')
    run('tar zxvf /etcd.gz -C / && rm -rf /etcd.gz && chown -R etcd:etcd /etc/etcd/ssl && mkdir -p /data/etcd && chown -R etcd:etcd /data/etcd && chmod 750 /data/etcd')
    local('rm -rf source/etcd/etcd.gz')
    run('systemctl daemon-reload && systemctl enable etcd')
    pass

def create_ssl_etcd(addnew = False):
    lvs = env.roledefs['etcd']['vip']
    hosts = ''
    lines_sed = 'N;'
    etcd_hosts = env.roledefs['etcd']['hosts']
    if addnew:
        etcd_hosts = env.roledefs['etcd']['hosts'] + env.roledefs['newetcd']['hosts']
    for index, host in enumerate(etcd_hosts):
        hosts += ',\\n      \\"' + host.split(':')[0] + '\\"'
        lines_sed += 'N;'

    local('cd source/etcd && sed "s#LVS#' + lvs + '#g" etcd-csr.json.tpl > etcd-csr.json')
    local('cd source/etcd && sed -i "s#HOSTS#' + hosts + '#g" etcd-csr.json')

    local('cd source/etcd && ./create_ssl.sh && /usr/bin/cp -rpf *.pem etc/etcd/ssl')

    #local('cd source/etcd && sed -i ":label;' + lines_sed + 's#' + hosts + '#HOSTS#;b label" etcd-csr.json')
    pass

@parallel
@roles('etcd')
def uninstall_etcd():
    run('systemctl disable etcd ; echo "" > /dev/null')
    # run('yum -y remove etcd')
    run('rm -rf /usr/bin/etcd')
    run('rm -rf /data/etcd /etc/etcd')
    pass

@parallel
@roles('etcd')
def install_lvsvip_etcd():
    execute(_install_lvsvip_etcd)
    pass

@parallel
@roles('newetcd')
def newetcd_install_lvsvip_etcd():
    if not len(env.roledefs['newetcd']['hosts']):
        return
    execute(_install_lvsvip_etcd)
    pass

def _install_lvsvip_etcd():
    lvsvip = env.roledefs['etcd']['vip']

    cmd = 'ifconfig lo:etcd:0 ' + lvsvip  + ' broadcast ' + lvsvip  + ' netmask 255.255.255.255 up'
    run(cmd + ' && echo -e "#/bin/sh\\n# chkconfig:   2345 90 10\\n' + cmd + '" > /etc/rc.d/init.d/vip_route_etcd.sh')
    cmd = 'route add -host ' + lvsvip  + ' dev lo:etcd:0 ; echo "" > /dev/null'
    run(cmd + ' && echo "' + cmd + '" >> /etc/rc.d/init.d/vip_route_etcd.sh')
    run('chmod +x /etc/rc.d/init.d/vip_route_etcd.sh && chkconfig --add vip_route_etcd.sh && chkconfig vip_route_etcd.sh on')

    run('echo "1" > /proc/sys/net/ipv4/conf/lo/arp_ignore && echo "2" > /proc/sys/net/ipv4/conf/lo/arp_announce && echo "1" > /proc/sys/net/ipv4/conf/all/arp_ignore && echo "2" > /proc/sys/net/ipv4/conf/all/arp_announce')
    pass

@parallel
@roles('etcd')
def uninstall_lvsvip_etcd():
    run('echo "0" > /proc/sys/net/ipv4/conf/lo/arp_ignore && echo "0" > /proc/sys/net/ipv4/conf/lo/arp_announce && echo "0" > /proc/sys/net/ipv4/conf/all/arp_ignore && echo "0" > /proc/sys/net/ipv4/conf/all/arp_announce')
    run('ifconfig lo:etcd:0 down ; echo "" > /dev/null')
    run('chkconfig vip_route_etcd.sh off && chkconfig --del vip_route_etcd.sh ; echo "" > /dev/null')
    run('rm -rf /etc/rc.d/init.d/vip_route_etcd.sh')
    pass
##########################[安装etcd]############################



##########################[安装master]############################
def install_master():
    # 证书要保证一样，所以只需要生成一次
    execute('create_ssl_master')
    execute('remote_install_master')

def newmaster_install_master():
    execute('create_ssl_master', True)
    execute('newmaster_remote_install_master')
    execute('newmaster_remote_modify_node_conf')
    pass

def newmaster_remote_modify_node_conf():
    execute('create_ssl_node')
    execute(_newmaster_remote_modify_node_conf)
    execute('kubeletcni_node', True)
    execute('newnode_kubeletcni_node', True)
    pass

@roles('node', 'newnode')
def _newmaster_remote_modify_node_conf():
    execute(_remote_install_node, True)
    pass

@roles('master')
def remote_install_master():
    execute(_remote_install_master)
    pass

@roles('master', 'newmaster')
def newmaster_remote_install_master():
    execute(_remote_install_master, True)
    pass

def _remote_install_master(addnew = False):
    curhost = env.host_string.split(':')[0]
    etcdlvs = env.roledefs['etcd']['vip']

    local('cd source/master && sed "s#K8S_HOST#' + curhost + '#g" config.tpl > etc/kubernetes/config')
    local('cd source/master && sed "s#ETCD_LVS_HOST#' + etcdlvs + '#g" apiserver.tpl > etc/kubernetes/apiserver')
    local('cd source/master && mkdir -p etc/kubernetes/pki/etcd && chmod 750 etc/kubernetes/pki/etcd')
    local('/usr/bin/cp -rpf source/etcd/etc/etcd/ssl/{ca.pem,etcd.pem,etcd-key.pem} source/master/etc/kubernetes/pki/etcd')

    if addnew and env.host_string in env.roledefs['master']['hosts']:
        local('cd source/master && tar zcvf master.gz etc')
    else:
        local('cd source/master && tar zcvf master.gz etc usr')
    put('source/master/master.gz', '/master.gz', mode=0640)
    run('tar zxvf /master.gz -C / && rm -rf /master.gz')
    local('rm -rf source/master/master.gz')
    run('systemctl daemon-reload && systemctl enable kube-apiserver && systemctl enable kube-controller-manager && systemctl enable kube-scheduler')
    pass

@roles('master')
def remote_update_master_etcd_ssl():
    local('/usr/bin/cp -rpf source/etcd/etc/etcd/ssl/{ca.pem,etcd.pem,etcd-key.pem} source/master/etc/kubernetes/pki/etcd')
    local('cd source/master/etc/kubernetes/pki && tar zcvf etcd.gz etcd')
    put('source/master/etc/kubernetes/pki/etcd.gz', '/etc/kubernetes/pki', mode=0640)
    run('cd /etc/kubernetes/pki && tar zxvf etcd.gz && rm -rf etcd.gz')
    local('rm -rf source/master/etc/kubernetes/pki/etcd.gz')
    pass

def create_ssl_master(addnew = False):
    hosts = ',\\n      \\"' + env.roledefs['publish']['hosts'][0].split(':')[0] + '\\"'
    lines_sed = 'N;'
    master_hosts = env.roledefs['master']['hosts']
    if addnew:
        master_hosts = env.roledefs['master']['hosts'] + env.roledefs['newmaster']['hosts']
    for index, host in enumerate(master_hosts):
        hosts += ',\\n      \\"' + host.split(':')[0] + '\\"'
        lines_sed += 'N;'
    masterlvs = env.roledefs['master']['vip']

    local('cd source/master && sed "s#LVS#' + masterlvs + '#g" admin-csr.json.tpl > admin-csr.json')
    local('cd source/master && sed "s#LVS#' + masterlvs + '#g" apiserver-csr.json.tpl > apiserver-csr.json')
    local('cd source/master && sed "s#LVS#' + masterlvs + '#g" ca-config.json.tpl > ca-config.json')
    local('cd source/master && sed "s#LVS#' + masterlvs + '#g" controller-manager-csr.json.tpl > controller-manager-csr.json')
    local('cd source/master && sed "s#LVS#' + masterlvs + '#g" scheduler-csr.json.tpl > scheduler-csr.json')
    local('cd source/master && sed -i "s#HOSTS#' + hosts + '#g" admin-csr.json apiserver-csr.json ca-config.json controller-manager-csr.json scheduler-csr.json')

    local('cd source/master && ./create_ssl.sh && /usr/bin/cp -rpf *.pem etc/kubernetes/pki')

    # admin.conf
    local('cd source/master && rm -rf etc/kubernetes/admin.conf && kubectl config set-cluster kubernetes --certificate-authority=ca.pem --embed-certs=true --server=https://' + masterlvs + ':6443 --kubeconfig=etc/kubernetes/admin.conf')
    local('cd source/master && kubectl config set-credentials kubernetes-admin --client-certificate=admin.pem --embed-certs=true --client-key=admin-key.pem --kubeconfig=etc/kubernetes/admin.conf')
    local('cd source/master && kubectl config set-context kubernetes-admin@kubernetes --cluster=kubernetes --user=kubernetes-admin --kubeconfig=etc/kubernetes/admin.conf')
    local('cd source/master && kubectl config use-context kubernetes-admin@kubernetes --kubeconfig=etc/kubernetes/admin.conf')
    local('cd source/master && mkdir -p /root/.kube && chmod 750 /root/.kube && /usr/bin/cp -rpf etc/kubernetes/admin.conf /root/.kube/config')

    # scheduler.conf
    local('cd source/master && rm -rf etc/kubernetes/scheduler.conf && kubectl config set-cluster kubernetes --certificate-authority=ca.pem --embed-certs=true --server=https://' + masterlvs + ':6443 --kubeconfig=etc/kubernetes/scheduler.conf')
    local('cd source/master && kubectl config set-credentials system:kube-scheduler --client-certificate=scheduler.pem --embed-certs=true --client-key=scheduler-key.pem --kubeconfig=etc/kubernetes/scheduler.conf')
    local('cd source/master && kubectl config set-context system:kube-scheduler@kubernetes --cluster=kubernetes --user=system:kube-scheduler --kubeconfig=etc/kubernetes/scheduler.conf')
    local('cd source/master && kubectl config use-context system:kube-scheduler@kubernetes --kubeconfig=etc/kubernetes/scheduler.conf')

    # controller-manager.conf
    local('cd source/master && rm -rf etc/kubernetes/controller-manager.conf && kubectl config set-cluster kubernetes --certificate-authority=ca.pem --embed-certs=true --server=https://' + masterlvs + ':6443 --kubeconfig=etc/kubernetes/controller-manager.conf')
    local('cd source/master && kubectl config set-credentials system:kube-controller-manager --client-certificate=controller-manager.pem --embed-certs=true --client-key=controller-manager-key.pem --kubeconfig=etc/kubernetes/controller-manager.conf')
    local('cd source/master && kubectl config set-context system:kube-controller-manager@kubernetes --cluster=kubernetes --user=system:kube-controller-manager --kubeconfig=etc/kubernetes/controller-manager.conf')
    local('cd source/master && kubectl config use-context system:kube-controller-manager@kubernetes --kubeconfig=etc/kubernetes/controller-manager.conf')
    pass


@parallel
@roles('master')
def uninstall_master():
    run('systemctl disable kube-apiserver && systemctl disable kube-controller-manager && systemctl disable kube-scheduler ; echo "" > /dev/null')
    run('rm -rf /etc/kubernetes /usr/bin/{kube-apiserver,kube-controller-manager,kube-scheduler} /usr/lib/systemd/system/{kube-apiserver.service,kube-controller-manager.service,kube-scheduler.service}')
    pass

@parallel
@roles('newmaster')
def newmaster_install_lvsvip_master():
    if not len(env.roledefs['newmaster']['hosts']):
        return
    execute(_install_lvsvip_master)
    pass

@parallel
@roles('master')
def install_lvsvip_master():
    execute(_install_lvsvip_master)
    pass

def _install_lvsvip_master():
    lvsvip = env.roledefs['master']['vip']
    cmd = 'ifconfig lo:master:0 ' + lvsvip  + ' broadcast ' + lvsvip  + ' netmask 255.255.255.255 up'
    run(cmd + ' && echo -e "#/bin/sh\\n# chkconfig:   2345 90 10\\n' + cmd + '" > /etc/rc.d/init.d/vip_route_master.sh')
    cmd = 'route add -host ' + lvsvip  + ' dev lo:master:0 ; echo "" > /dev/null'
    run(cmd + ' && echo "' + cmd + '" >> /etc/rc.d/init.d/vip_route_master.sh')
    run('chmod +x /etc/rc.d/init.d/vip_route_master.sh && chkconfig --add vip_route_master.sh && chkconfig vip_route_master.sh on')
    run('echo "1" > /proc/sys/net/ipv4/conf/lo/arp_ignore && echo "2" > /proc/sys/net/ipv4/conf/lo/arp_announce && echo "1" > /proc/sys/net/ipv4/conf/all/arp_ignore && echo "2" > /proc/sys/net/ipv4/conf/all/arp_announce')
    pass

@parallel
@roles('master')
def uninstall_lvsvip_master():
    run('echo "0" > /proc/sys/net/ipv4/conf/lo/arp_ignore && echo "0" > /proc/sys/net/ipv4/conf/lo/arp_announce && echo "0" > /proc/sys/net/ipv4/conf/all/arp_ignore && echo "0" > /proc/sys/net/ipv4/conf/all/arp_announce')
    run('ifconfig lo:master:0 down ; echo "" > /dev/null')
    run('chkconfig vip_route_master.sh off && chkconfig --del vip_route_master.sh ; echo "" > /dev/null')
    run('rm -rf /etc/rc.d/init.d/vip_route_master.sh')
    pass
##########################[安装master]############################


##########################[安装docker证书]############################
@roles('publish', 'node')
def install_dockercrt():
    execute(_install_dockercrt)
    pass

@roles('newnode')
def newnode_install_dockercrt():
    if not len(env.roledefs['newnode']['hosts']):
        return
    execute(_install_dockercrt)
    pass

def _install_dockercrt():
    pridocker = env.roledefs['pridocker']['hosts'][0].split(':')[0]

    local('cd source/docker && rm -rf etc/docker/certs.d/* && chmod 640 ca.crt && mkdir etc/docker/certs.d/' + pridocker + ':5000 && chmod 750 etc/docker/certs.d/' + pridocker + ':5000 && /usr/bin/cp -rpf ca.crt etc/docker/certs.d/' + pridocker + ':5000')

    local('cd source/docker && tar zcvf docker.gz etc')
    put('source/docker/docker.gz', '/docker.gz', mode=0640)
    run('tar zxvf /docker.gz -C / && rm -rf /docker.gz')
    local('rm -rf source/docker/docker.gz')
    pass

@roles('publish', 'node')
def uninstall_dockercrt():
    pridocker = env.roledefs['pridocker']['hosts'][0].split(':')[0]
    run('rm -rf /etc/docker/certs.d/' + pridocker + ':5000')
    pass
##########################[安装docker证书]############################


##########################[安装node]############################
def install_node():
    execute('create_ssl_node')
    execute('remote_install_node')

@roles('node')
def remote_install_node():
    execute(_remote_install_node)
    pass

@roles('newnode')
def newnode_install_node():
    if not len(env.roledefs['newnode']['hosts']):
        return
    execute(_remote_install_node)
    pass

def _remote_install_node(onlyconf = False):
    curhost = env.host_string.split(':')[0]
    masterlvs = env.roledefs['master']['vip']
    pridocker = env.roledefs['pridocker']['hosts'][0].split(':')[0]

    local('cd source/node && sed "s#NODE_HOST#' + curhost + '#g" kubelet-csr.json.tpl > kubelet-csr.json')
    local('cd source/node && sed "s#NODE_HOST#' + curhost + '#g" kubelet.tpl > etc/kubernetes/kubelet')
    local('cd source/node && sed "s#NODE_HOST#' + curhost + '#g" kubelet.yaml.tpl > etc/kubernetes/kubelet.yaml')
    local('cd source/node && sed "s#K8S_MASTER_LVS#' + masterlvs + '#g" config.tpl > etc/kubernetes/config')
    local('cd source/node && sed -i "s#K8S_MASTER_LVS#' + masterlvs + '#g" etc/kubernetes/kubelet')
    local('cd source/node && sed -i "s#PRI_DOCKER_HOST#' + pridocker + '#g" etc/kubernetes/kubelet')

    local('cd source/node && cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=frognew kubelet-csr.json | cfssljson -bare kubelet')

    # kubelet.conf
    local('cd source/node && rm -rf etc/kubernetes/kubelet.conf && kubectl config set-cluster kubernetes --certificate-authority=ca.pem --embed-certs=true --server=https://' + masterlvs + ':6443 --kubeconfig=etc/kubernetes/kubelet.conf')
    local('cd source/node && kubectl config set-credentials system:node:' + curhost + ' --client-certificate=kubelet.pem --embed-certs=true --client-key=kubelet-key.pem --kubeconfig=etc/kubernetes/kubelet.conf')
    local('cd source/node && kubectl config set-context system:node:' + curhost + '@kubernetes --cluster=kubernetes --user=system:node:' + curhost + ' --kubeconfig=etc/kubernetes/kubelet.conf')
    local('cd source/node && kubectl config use-context system:node:' + curhost + '@kubernetes --kubeconfig=etc/kubernetes/kubelet.conf')

    local('cd source/node && /usr/bin/cp -rpf *.pem etc/kubernetes/pki')

    if onlyconf:
        local('cd source/node && tar zcvf node.gz etc')
    else:
        local('cd source/node && tar zcvf node.gz etc usr')
    put('source/node/node.gz', '/node.gz', mode=0640)
    run('tar zxvf /node.gz -C / && rm -rf /node.gz')
    local('rm -rf source/node/node.gz')
    run('systemctl daemon-reload && systemctl enable kube-proxy && systemctl enable kubelet && mkdir -p /data/kubelet && chmod 750 /data/kubelet')
    pass

def create_ssl_node():
    masterlvs = env.roledefs['master']['vip']

    local('cd source/node && /usr/bin/cp -rpf ../master/{ca.pem,ca-key.pem,ca-config.json} .')
    local('cd source/node && cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=frognew kube-proxy-csr.json | cfssljson -bare kube-proxy')

    # kube-proxy.conf
    local('cd source/node && rm -rf etc/kubernetes/kube-proxy.conf && kubectl config set-cluster kubernetes --certificate-authority=ca.pem --embed-certs=true --server=https://' + masterlvs + ':6443 --kubeconfig=etc/kubernetes/kube-proxy.conf')
    local('cd source/node && kubectl config set-credentials system:kube-proxy --client-certificate=kube-proxy.pem --embed-certs=true --client-key=kube-proxy-key.pem --kubeconfig=etc/kubernetes/kube-proxy.conf')
    local('cd source/node && kubectl config set-context system:kube-proxy@kubernetes --cluster=kubernetes --user=system:kube-proxy --kubeconfig=etc/kubernetes/kube-proxy.conf')
    local('cd source/node && kubectl config use-context system:kube-proxy@kubernetes --kubeconfig=etc/kubernetes/kube-proxy.conf')
    pass

@parallel
@roles('node')
def uninstall_node():
    run('systemctl disable kube-proxy && systemctl disable kubelet ; echo "" > /dev/null')
    run('rm -rf /data/kubelet /etc/kubernetes /usr/bin/{kubelet,kube-proxy} /usr/lib/systemd/system/{kubelet.service,kube-proxy.service}')
    pass
##########################[安装node]############################


##########################[安装dns]############################
@parallel
@roles('pridns')
def install_dns():
    local('cd source/bind && tar zcvf bind.gz var etc')
    run('yum install -y bind-chroot')
    put('source/bind/bind.gz', '/tmp', mode=0640)
    run('tar zxvf /tmp/bind.gz -C / && rm -rf /tmp/bind.gz && chown -R named:named /var/named/zones && chown root:named /var/named /etc/named.conf /etc/named.rfc1912.zones && systemctl enable named-chroot')
    local('rm -rf source/bind/bind.gz')
    pass

@parallel
@roles('pridns')
def uninstall_dns():
    run('yum remove -y bind-chroot')
    pass
##########################[安装dns]############################


##########################[初始化镜像]############################
def init_images():
    pridocker = env.roledefs['pridocker']['hosts'][0].split(':')[0]

    local('docker images | grep "alpine" || (cd source/images && sha256=`docker load -i esn-containers~alpine:latest.tar | grep Loaded | awk \'{print $4}\' | awk -F \':\' \'{print $2}\'` && docker tag $sha256 ' + pridocker + ':5000/esn-containers/alpine:latest)')

    local('docker images | grep "esn_base" || (cd source/images && sha256=`docker load -i esn-containers~esn_base:1.0.tar | grep Loaded | awk \'{print $4}\' | awk -F \':\' \'{print $2}\'` && docker tag $sha256 ' + pridocker + ':5000/esn-containers/esn_base:1.0)')

    local('docker images | grep "pause-amd64" || (cd source/images && sha256=`docker load -i HOST:PORT~google-containers~pause-amd64:3.1.tar | grep Loaded | awk \'{print $4}\' | awk -F \':\' \'{print $2}\'` && docker tag $sha256 ' + pridocker + ':5000/google-containers/pause-amd64:3.1 && docker push ' + pridocker + ':5000/google-containers/pause-amd64:3.1)')

    local('docker images | grep "kubernetes-dashboard-amd64" || (cd source/images && sha256=`docker load -i HOST:PORT~k8s.gcr.io~kubernetes-dashboard-amd64:v1.10.1.tar | grep Loaded | awk \'{print $4}\' | awk -F \':\' \'{print $2}\'` && docker tag $sha256 ' + pridocker + ':5000/k8s.gcr.io/kubernetes-dashboard-amd64:v1.10.1 && docker push ' + pridocker + ':5000/k8s.gcr.io/kubernetes-dashboard-amd64:v1.10.1)')

    local('docker images | grep "coredns" || (cd source/images && sha256=`docker load -i HOST:PORT~k8s.gcr.io~coredns:1.2.6.tar | grep Loaded | awk \'{print $4}\' | awk -F \':\' \'{print $2}\'` && docker tag $sha256 ' + pridocker + ':5000/k8s.gcr.io/coredns:1.2.6 && docker push ' + pridocker + ':5000/k8s.gcr.io/coredns:1.2.6)')

    local('docker images | grep "heapster-influxdb-amd64" || (cd source/images && sha256=`docker load -i HOST:PORT~google_containers~heapster-influxdb-amd64:v1.3.3.tar | grep Loaded | awk \'{print $4}\' | awk -F \':\' \'{print $2}\'` && docker tag $sha256 ' + pridocker + ':5000/google_containers/heapster-influxdb-amd64:v1.3.3 && docker push ' + pridocker + ':5000/google_containers/heapster-influxdb-amd64:v1.3.3)')

    local('docker images | grep "heapster-grafana-amd64" || (cd source/images && sha256=`docker load -i HOST:PORT~google_containers~heapster-grafana-amd64:v4.4.3.tar | grep Loaded | awk \'{print $4}\' | awk -F \':\' \'{print $2}\'` && docker tag $sha256 ' + pridocker + ':5000/google_containers/heapster-grafana-amd64:v4.4.3 && docker push ' + pridocker + ':5000/google_containers/heapster-grafana-amd64:v4.4.3)')

    local('docker images | grep "heapster-amd64" || (cd source/images && sha256=`docker load -i HOST:PORT~google_containers~heapster-amd64:v1.5.3.tar | grep Loaded | awk \'{print $4}\' | awk -F \':\' \'{print $2}\'` && docker tag $sha256 ' + pridocker + ':5000/google_containers/heapster-amd64:v1.5.3 && docker push ' + pridocker + ':5000/google_containers/heapster-amd64:v1.5.3)')

    local('docker images | grep "calico" | grep node || (cd source/images && sha256=`docker load -i HOST:PORT~quay.io~calico~node:v3.1.3.tar | grep Loaded | awk \'{print $4}\' | awk -F \':\' \'{print $2}\'` && docker tag $sha256 ' + pridocker + ':5000/quay.io/calico/node:v3.1.3 && docker push ' + pridocker + ':5000/quay.io/calico/node:v3.1.3)')

    local('docker images | grep "calico" | grep cni || (cd source/images && sha256=`docker load -i HOST:PORT~quay.io~calico~cni:v3.1.3.tar | grep Loaded | awk \'{print $4}\' | awk -F \':\' \'{print $2}\'` && docker tag $sha256 ' + pridocker + ':5000/quay.io/calico/cni:v3.1.3 && docker push ' + pridocker + ':5000/quay.io/calico/cni:v3.1.3)')

    local('docker images | grep "calico" | grep kube-controllers || (cd source/images && sha256=`docker load -i HOST:PORT~quay.io~calico~kube-controllers:v3.1.3.tar | grep Loaded | awk \'{print $4}\' | awk -F \':\' \'{print $2}\'` && docker tag $sha256 ' + pridocker + ':5000/quay.io/calico/kube-controllers:v3.1.3 && docker push ' + pridocker + ':5000/quay.io/calico/kube-controllers:v3.1.3)')
    pass
##########################[初始化镜像]############################


##########################[初始化calico]############################
def init_calico():
    etcdlvs = env.roledefs['etcd']['vip']
    pridocker = env.roledefs['pridocker']['hosts'][0].split(':')[0]

    local('kubectl delete -f source/calico ; echo "" > /dev/null')

    local('sed "s#PRI_DOCKER_HOST#' + pridocker + '#g" source/calico/calico.yaml.tpl > source/calico/calico.yaml')
    local('sed -i "s#ETCD_LVS_HOST#' + etcdlvs + '#g" source/calico/calico.yaml')
    local('TLS_ETCD_KEY=$(cat source/etcd/etc/etcd/ssl/etcd-key.pem | base64 | tr -d "\n") && sed -i "s#TLS_ETCD_KEY#$TLS_ETCD_KEY#g" source/calico/calico.yaml')
    local('TLS_ETCD_CERT=$(cat source/etcd/etc/etcd/ssl/etcd.pem | base64 | tr -d "\n") && sed -i "s#TLS_ETCD_CERT#$TLS_ETCD_CERT#g" source/calico/calico.yaml')
    local('TLS_ETCD_CA=$(cat source/etcd/etc/etcd/ssl/ca.pem | base64 | tr -d "\n") && sed -i "s#TLS_ETCD_CA#$TLS_ETCD_CA#g" source/calico/calico.yaml')

    local('kubectl apply -f source/calico')

    i = 0
    while True:
        i = i + 1
        num = local('kubectl get pods -o wide -n kube-system | grep calico | grep Running | wc -l', capture = True)
        total = len(env.roledefs['node']['hosts']) + 1
        print '等待所有节点calico容器正常运行(%ds)(%d = %s)' % (i, total, num)
        if int(num) == total:
            break
        # 15次都没成功，重启一下master服务
        if i == 30:
            execute(service_master, dowhat = 'restart')
        time.sleep(3)
    pass
##########################[初始化calico]############################


##########################[修改kubelet配置，加载cni网络插件（calico启动后才会生成）]############################
@parallel
@roles('node')
def kubeletcni_node(onlyconf = False):
    execute(_kubeletcni_node, onlyconf)
    pass

@parallel
@roles('newnode')
def newnode_kubeletcni_node(onlyconf = False):
    if not len(env.roledefs['newnode']['hosts']):
        return
    execute(_kubeletcni_node, onlyconf)
    pass

def _kubeletcni_node(onlyconf = False):
    run("sed -i 's#--config=/etc/kubernetes/kubelet.yaml\"#--config=/etc/kubernetes/kubelet.yaml --network-plugin=cni --cni-conf-dir=/etc/cni/net.d --cni-bin-dir=/opt/cni/bin\"#g' /etc/kubernetes/kubelet")
    if not onlyconf:
        run('systemctl restart kubelet')
    pass
##########################[修改kubelet配置，加载cni网络插件（calico启动后才会生成）]############################


##########################[初始化k8s系统]############################
def init_k8s_system():
    pridocker = env.roledefs['pridocker']['hosts'][0].split(':')[0]
    pridns = env.roledefs['pridns']['hosts'][0].split(':')[0]

    local('sed "s#PRI_DOCKER_HOST#' + pridocker + '#g" source/dashboard/dashboard-controller.yaml.tpl > source/dashboard/dashboard-controller.yaml')
    local('sed "s#PRI_DOCKER_HOST#' + pridocker + '#g" source/dns/coredns.yaml.tpl > source/dns/coredns.yaml')
    local('sed "s#PRI_DOCKER_HOST#' + pridocker + '#g" source/heapster/grafana.yaml.tpl > source/heapster/grafana.yaml')
    local('sed "s#PRI_DOCKER_HOST#' + pridocker + '#g" source/heapster/heapster.yaml.tpl > source/heapster/heapster.yaml')
    local('sed "s#PRI_DOCKER_HOST#' + pridocker + '#g" source/heapster/influxdb.yaml.tpl > source/heapster/influxdb.yaml')
    local('sed "s#PRI_DOCKER_HOST#' + pridocker + '#g" source/heapster/Dockerfile.tpl > source/heapster/Dockerfile')
    local('sed -i "s#HOST#' + pridns + '#g" source/dns/coredns.yaml')

    local('kubectl apply -f source/dashboard')
    local('kubectl apply -f source/dns')
    local('kubectl apply -f source/heapster')
    pass
##########################[初始化k8s系统]############################


##########################[初始化web_test]############################
def init_web_test():
    pridocker = env.roledefs['pridocker']['hosts'][0].split(':')[0]

    local('cd source/web_test && sed "s#PRI_DOCKER_HOST#' + pridocker + '#g" Dockerfile.tpl > Dockerfile')
    local('cd source/web_test && sed "s#PRI_DOCKER_HOST#' + pridocker + '#g" create.sh.tpl > create.sh && chmod 750 create.sh')

    local('cd source/web_test && ./create.sh')
    pass
##########################[初始化web_test]############################


##########################[k8s系统报错]############################
@parallel
@roles('publish', 'etcd', 'master', 'node', 'newetcd', 'newmaster', 'newnode', 'lvs', 'pridocker')
def error():
    run('tail -f /var/log/messages')
    pass
##########################[k8s系统报错]############################
