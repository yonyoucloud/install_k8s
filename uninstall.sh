#!/bin/sh
# -------------------------------------------------------------------------------
# Filename:    uninstall.sh
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

echo -e "\033[32m{`date`}[开始]关闭所有服务.............................\033[0m"
fab service:stop || exit 1
echo -e "\033[32m{`date`}[结束]关闭所有服务.............................\n\n\n\n\n\n\033[0m"

echo -e "\033[32m{`date`}[开始]卸载dns.............................\033[0m"
fab uninstall_dns || exit 1
echo -e "\033[32m{`date`}[结束]卸载dns.............................\n\n\n\n\n\n\033[0m"

echo -e "\033[32m{`date`}[开始]卸载负载均衡.............................\033[0m"
fab uninstall_lvs || exit 1
echo -e "\033[32m{`date`}[结束]卸载负载均衡.............................\n\n\n\n\n\n\033[0m"

echo -e "\033[32m{`date`}[开始]卸载docker证书.............................\033[0m"
fab uninstall_dockercrt || exit 1
echo -e "\033[32m{`date`}[结束]卸载docker证书.............................\n\n\n\n\n\n\033[0m"

echo -e "\033[32m{`date`}[开始]卸载node节点.............................\033[0m"
fab uninstall_node || exit 1
echo -e "\033[32m{`date`}[结束]卸载node节点.............................\n\n\n\n\n\n\033[0m"

echo -e "\033[32m{`date`}[开始]卸载master节点.............................\033[0m"
fab uninstall_master || exit 1
echo -e "\033[32m{`date`}[结束]卸载master节点.............................\n\n\n\n\n\n\033[0m"

echo -e "\033[32m{`date`}[开始]卸载etcd.............................\033[0m"
fab uninstall_etcd || exit 1
echo -e "\033[32m{`date`}[结束]卸载etcd.............................\n\n\n\n\n\n\033[0m"

echo -e "\033[32m{`date`}[开始]卸载docker私有仓库.............................\033[0m"
fab uninstall_pridocker || exit 1
echo -e "\033[32m{`date`}[结束]卸载docker私有仓库.............................\n\n\n\n\n\n\033[0m"

echo -e "\033[32m{`date`}[开始]卸载docker.............................\033[0m"
fab uninstall_docker || exit 1
echo -e "\033[32m{`date`}[结束]卸载docker.............................\n\n\n\n\n\n\033[0m"

echo -e "\033[32m{`date`}[开始]卸载工具包.............................\033[0m"
rm -rf /usr/local/bin/cfssl*
rm -rf /usr/local/bin/{etcdctl,kubectl,kubeadm,kubemark}
echo -e "\033[32m{`date`}[结束]卸载工具包.............................\n\n\n\n\n\n\033[0m"

echo -e "\033[31m{`date`}[k8s集群卸载完毕！].............................\n\n\n\n\n\n\033[0m"
