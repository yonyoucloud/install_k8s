#!/bin/sh
# -------------------------------------------------------------------------------
# Filename:    add_node.sh
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

echo -e "\033[32m{`date`}[开始]关闭master服务服务.............................\033[0m"
fab service_master:stop || exit 1
echo -e "\033[32m{`date`}[结束]关闭master服务服务.............................\n\n\n\n\n\n\033[0m"

echo -e "\033[32m{`date`}[开始]安装基础环境.............................\033[0m"
fab newmaster_install_base || exit 1
echo -e "\033[32m{`date`}[结束]安装基础环境.............................\n\n\n\n\n\n\033[0m"

echo -e "\033[32m{`date`}[开始]安装master节点.............................\033[0m"
fab newmaster_install_master || exit 1
echo -e "\033[32m{`date`}[结束]安装master节点.............................\n\n\n\n\n\n\033[0m"

echo -e "\033[32m{`date`}[开始]启动新master节点.............................\033[0m"
fab service_master:start || exit 1
fab newmaster_service_master_start || exit 1
echo -e "\033[32m{`date`}[结束]启动新master节点.............................\n\n\n\n\n\n\033[0m"

echo -e "\033[32m{`date`}[开始]修改lvs配置.............................\033[0m"
fab newmaster_install_lvsvip_master || exit 1
fab remote_install_lvs_new || exit 1
echo -e "\033[32m{`date`}[结束]修改lvs配置.............................\n\n\n\n\n\n\033[0m"

echo -e "\033[32m{`date`}[开始]重启node服务.............................\033[0m"
fab newmaster_service_node_restart || exit 1
echo -e "\033[32m{`date`}[结束]重启node服务.............................\n\n\n\n\n\n\033[0m"

echo -e "\033[32m{`date`}[结束]所有节点添加完毕\033[0m\033[31m[祝您好运！]\033[0m\033[32m.............................\n\n\n\n\n\n\033[0m"
