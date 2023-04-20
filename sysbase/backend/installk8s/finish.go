package installk8s

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (ik *InstallK8s) FinishInstall() {
	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return
	}

	nodeRole, ok := ik.resources["node"]
	if !ok {
		ik.Stdout <- "没有node资源"
		return
	}

	ik.Stdout <- "[开始]启动所有服务"

	ik.Stdout <- "启动Dns"
	ik.Params.DoWhat = "start"
	ik.ServiceDns()
	ik.wait(5, "等待...")

	ik.Stdout <- "启动Publish"
	ik.Params.DoWhat = "start"
	ik.Params.DoContainerd = true
	ik.ServicePublish()

	ik.Stdout <- "启动Etcd集群"
	ik.ServiceEtcd()
	ik.Params.DoWhat = "restart"
	ik.ServiceEtcd()
	ik.wait(10, "等待...")

	ik.Stdout <- "启动Master集群"
	ik.Params.DoWhat = "start"
	ik.ServiceMaster()
	ik.Params.DoWhat = "restart"
	ik.ServiceMaster()
	ik.wait(8, "等待...")

	ik.Stdout <- "启动Node集群"
	ik.Params.DoWhat = "start"
	ik.ServiceNode()
	ik.Params.DoWhat = "restart"
	ik.ServiceNode()
	ik.wait(8, "等待...")

	ik.Stdout <- "[结束]启动所有服务"

	ik.Stdout <- "[开始]验证k8s集群"
	publishRole.WaitOutput = true
	ik.er.SetRole(publishRole)
	for i := 1; i <= 180; i++ {
		ik.Stdout <- fmt.Sprintf("等待倒计时(%d)s...", i)
		ik.er.Run("kubectl get nodes -o wide | grep NotReady > /dev/null ; echo $?")
		if ik.er.GetCmdReturn()[0] == "1" {
			publishRole.WaitOutput = false
			ik.er.SetRole(publishRole)
			ik.er.Run("kubectl get nodes -o wide")
			break
		}
		time.Sleep(1 * time.Second)
	}
	ik.Stdout <- "[结束]验证k8s集群"

	ik.Stdout <- "[开始]初始化镜像"
	ik.initImages()
	ik.Stdout <- "[结束]初始化镜像"

	ik.Stdout <- "[开始]初始化calico"
	ik.initCalico()
	ik.kubeletcniNode(nodeRole)
	ik.Stdout <- "[结束]初始化calico"

	ik.Stdout <- "[开始]初始k8s系统镜像服务"
	ik.initK8sSystem()
	ik.Stdout <- "[结束]初始k8s系统镜像服务"

	ik.Stdout <- "[开始]安装Istio"
	ik.installIstio()
	ik.Stdout <- "[结束]安装Istio"

	ik.Stdout <- "[开始]初始化测试微服务"
	ik.initWebTest()
	ik.Stdout <- "[结束]初始化测试微服务"

	ik.Stdout <- "[开始]需要您验证测试以下说明"
	publishRole.WaitOutput = true
	ik.er.SetRole(publishRole)
	for i := 1; i <= 180; i++ {
		ik.Stdout <- fmt.Sprintf("等待kubernetes-dashboard running(%d)s...", i)
		ik.er.Run("kubectl -n kube-system get pods -o wide | grep kubernetes-dashboard | grep Running > /dev/null ; echo $?")
		if ik.er.GetCmdReturn()[0] == "0" {
			publishRole.WaitOutput = false
			ik.er.SetRole(publishRole)
			ik.er.Run(`kubectl -n istio-system get pod -o wide | grep istio-ingressgateway | grep Running | awk '{print "设置Hosts: "$7" dashboard.k8s.com 然后您可以访问kubernetes-dashboard: https://dashboard.k8s.com:10443"}'`)
			ik.Stdout <- "用下面输出的token登录kubernetes-dashboard"
			// ik.er.Run(`kubectl describe secret $(kubectl get secret -n kube-system | grep admin-token | awk '{print $1}') -n kube-system | grep token: | awk '{print $2}'`)
			ik.er.Run(`kubectl -n kube-system get secret admin-user-token -oyaml | grep token: | awk -F 'token: ' '{print $2}' | base64 -d && echo`)
			break
		}
		time.Sleep(1 * time.Second)
	}

	publishRole.WaitOutput = true
	ik.er.SetRole(publishRole)
	for i := 1; i <= 180; i++ {
		ik.Stdout <- fmt.Sprintf("等待grafana running(%d)s...", i)
		ik.er.Run("kubectl -n kube-system get pods -o wide | grep grafana | grep Running > /dev/null ; echo $?")
		if ik.er.GetCmdReturn()[0] == "0" {
			publishRole.WaitOutput = false
			ik.er.SetRole(publishRole)
			ik.er.Run(`kubectl -n istio-system get pod -o wide | grep istio-ingressgateway | grep Running | awk '{print "设置Hosts: "$7" grafana.k8s.com 然后您可以访问grafana: http://grafana.k8s.com:10080 或 https://grafana.k8s.com:10443"}'`)
			ik.Stdout <- "账号密码为：admin/123456"
			// ik.Stdout <- "注意：需要配置一下grafana的k8s插件中的URL地址及三个认证证书（base64解码~/.kube/config中的相应证书）"
			break
		}
		time.Sleep(1 * time.Second)
	}

	publishRole.WaitOutput = true
	ik.er.SetRole(publishRole)
	for i := 1; i <= 180; i++ {
		ik.Stdout <- fmt.Sprintf("等待prometheus running(%d)s...", i)
		ik.er.Run("kubectl -n kube-system get pods -o wide | grep prometheus | grep Running > /dev/null ; echo $?")
		if ik.er.GetCmdReturn()[0] == "0" {
			publishRole.WaitOutput = false
			ik.er.SetRole(publishRole)
			ik.er.Run(`kubectl -n istio-system get pod -o wide | grep istio-ingressgateway | grep Running | awk '{print "设置Hosts: "$7" prometheus.k8s.com 然后您可以访问prometheus: http://prometheus.k8s.com:10080 或 https://prometheus.k8s.com:10443"}'`)
			break
		}
		time.Sleep(1 * time.Second)
	}

	publishRole.WaitOutput = true
	ik.er.SetRole(publishRole)
	for i := 1; i <= 180; i++ {
		ik.Stdout <- fmt.Sprintf("等待kiali running(%d)s...", i)
		ik.er.Run("kubectl -n istio-system get pods -o wide | grep kiali | grep Running > /dev/null ; echo $?")
		if ik.er.GetCmdReturn()[0] == "0" {
			publishRole.WaitOutput = false
			ik.er.SetRole(publishRole)
			ik.er.Run(`kubectl -n istio-system get pod -o wide | grep istio-ingressgateway | grep Running | awk '{print "设置Hosts: "$7" kiali.k8s.com 然后您可以访问kiali: http://kiali.k8s.com:10080 或 https://kiali.k8s.com:10443"}'`)
			break
		}
		time.Sleep(1 * time.Second)
	}

	publishRole.WaitOutput = true
	ik.er.SetRole(publishRole)
	for i := 1; i <= 180; i++ {
		ik.Stdout <- fmt.Sprintf("等待web-test running(%d)s...", i)
		ik.er.Run("kubectl -n test-system get pods -o wide | grep web-test | grep Running > /dev/null ; echo $?")
		if ik.er.GetCmdReturn()[0] == "0" {
			publishRole.WaitOutput = false
			ik.er.SetRole(publishRole)
			cmds := []string{
				fmt.Sprintf(`chmod 600 %s/test-base/ssh/root/* %s/test-base/ssh/esn/*`, ik.SourceDir, ik.SourceDir),
				fmt.Sprintf(`kubectl -n istio-system get pod -o wide | grep istio-ingressgateway | grep Running | awk '{print "设置Hosts: "$7" test.k8s.com 然后您可以访问web-test: http://test.k8s.com:10080 或 https://test.k8s.com:10443";print "您可以执行: ssh -i %s/test-base/ssh/root/id_rsa root@"$6" 直接登录到容器中";print "您也可以执行: ssh -i %s/test-base/ssh/esn/id_rsa esn@"$6" 直接登录到容器中"}'`, ik.SourceDir, ik.SourceDir),
			}
			ik.er.Run(cmds...)
			break
		}
		time.Sleep(1 * time.Second)
	}

	ik.Stdout <- "您可以进入到容器中执行: ping t.test.com 看是否解析到10.10.10.10上, 或看下面测试输出"

	publishRole.WaitOutput = true
	ik.er.SetRole(publishRole)
	ik.er.Run("kubectl -n test-system get pods | grep web-test | awk '{print $1}'")
	pod := ik.er.GetCmdReturn()[0]

	publishRole.WaitOutput = false
	ik.er.SetRole(publishRole)
	ik.er.Run(fmt.Sprintf(`kubectl -n test-system exec %s -- ping -c 5 t.test.com`, pod))

	ik.Stdout <- "[结束]需要您验证测试以下说明"
	ik.Stdout <- "祝您好运，安全稳定的k8s集群安装完毕！"
}

func (ik *InstallK8s) initImages() {
	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return
	}

	ik.er.SetRole(publishRole)
	ik.er.Run(fmt.Sprintf(`cd %s/images && for file in $(find . -type f); do file=${file/\.\/}; saveIFS=$IFS IFS="/" fileArr=($file) IFS=$saveIFS; count=${#fileArr[*]} repository=$(IFS="/"; echo "${fileArr[*]: 0: $count-1}") tag=${fileArr[*]: -1}; nerdctl -n k8s.io load -i $file && nerdctl -n k8s.io push "${repository}:${tag}";done`, ik.SourceDir))
}

func (ik *InstallK8s) initCalico() {
	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return
	}

	etcdLbRole, ok := ik.resources["etcdlb"]
	if !ok {
		ik.Stdout <- "没有etcdlb资源"
		return
	}

	nodeRole, ok := ik.resources["node"]
	if !ok {
		ik.Stdout <- "没有node资源"
		return
	}

	etcdLbHost := strings.Split(etcdLbRole.Hosts[0], ":")[0]

	ik.er.SetRole(publishRole)
	cmds := []string{
		fmt.Sprintf("kubectl delete -f %s/calico", ik.SourceDir),
		fmt.Sprintf(`sed "s#ETCD_LVS_HOST#%s#g" %s/calico/calico.yaml.tpl > %s/calico/calico.yaml`, etcdLbHost, ik.SourceDir, ik.SourceDir),
		fmt.Sprintf(`TLS_ETCD_KEY=$(cat %s/etcd/etc/etcd/ssl/etcd-key.pem | base64 | tr -d "\n") && sed -i "s#TLS_ETCD_KEY#$TLS_ETCD_KEY#g" %s/calico/calico.yaml`, ik.SourceDir, ik.SourceDir),
		fmt.Sprintf(`TLS_ETCD_CERT=$(cat %s/etcd/etc/etcd/ssl/etcd.pem | base64 | tr -d "\n") && sed -i "s#TLS_ETCD_CERT#$TLS_ETCD_CERT#g" %s/calico/calico.yaml`, ik.SourceDir, ik.SourceDir),
		fmt.Sprintf(`TLS_ETCD_CA=$(cat %s/etcd/etc/etcd/ssl/ca.pem | base64 | tr -d "\n") && sed -i "s#TLS_ETCD_CA#$TLS_ETCD_CA#g" %s/calico/calico.yaml`, ik.SourceDir, ik.SourceDir),
		fmt.Sprintf("kubectl apply -f %s/calico", ik.SourceDir),
	}
	ik.er.Run(cmds...)

	total := len(nodeRole.Hosts) + 1
	i := 0
	for {
		i++
		publishRole.WaitOutput = true
		ik.er.SetRole(publishRole)
		ik.er.Run("kubectl get pods -o wide -n kube-system | grep calico | grep Running | wc -l")
		num, _ := strconv.Atoi(ik.er.GetCmdReturn()[0])
		ik.Stdout <- fmt.Sprintf("等待所有节点calico容器正常运行(%ds)(%d = %d)", i, total, num)
		if num == total {
			break
		}
		if i == 30 {
			ik.Params.DoWhat = "restart"
			ik.ServiceMaster()
		}
		time.Sleep(3 * time.Second)
	}
}

func (ik *InstallK8s) initK8sSystem() {
	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return
	}

	pridnsRole, ok := ik.resources["pridns"]
	if !ok {
		ik.Stdout <- "没有pridns资源"
		return
	}

	pridnsHost := strings.Split(pridnsRole.Hosts[0], ":")[0]

	ik.er.SetRole(publishRole)
	cmds := []string{
		fmt.Sprintf(`sed "s#HOST#%s#g" %s/dns/coredns.yaml.tpl > %s/dns/coredns.yaml`, pridnsHost, ik.SourceDir, ik.SourceDir),

		// 生成TLS证书
		fmt.Sprintf(`rm -rf %s/addons/certs ; mkdir -p %s/addons/certs`, ik.SourceDir, ik.SourceDir),
		fmt.Sprintf(`cat > %s/addons/certs/extfile.cnf <<-EOF
[ v3_ca ]
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1=k8s.com
DNS.2=*.k8s.com
EOF`, ik.SourceDir),
		fmt.Sprintf(`openssl req -out %s/addons/certs/k8s.com.csr -newkey rsa:2048 -nodes -keyout %s/addons/certs/k8s.com.key -subj "/CN=k8s.com/O=k8s organization"`, ik.SourceDir, ik.SourceDir),
		fmt.Sprintf(`openssl x509 -req -days 36500 -CA %s/master/etc/kubernetes/pki/ca.pem -CAkey %s/master/etc/kubernetes/pki/ca-key.pem -CAcreateserial -in %s/addons/certs/k8s.com.csr -out %s/addons/certs/k8s.com.crt -extfile %s/addons/certs/extfile.cnf -extensions v3_ca`, ik.SourceDir, ik.SourceDir, ik.SourceDir, ik.SourceDir, ik.SourceDir),

		`kubectl -n kube-system delete secret kubernetes-dashboard-certs`,
		fmt.Sprintf(`kubectl -n kube-system create secret tls kubernetes-dashboard-certs --key=%s/addons/certs/k8s.com.key --cert=%s/addons/certs/k8s.com.crt`, ik.SourceDir, ik.SourceDir),

		fmt.Sprintf(`kubectl apply -f %s/dns`, ik.SourceDir),
		fmt.Sprintf(`kubectl apply -f %s/addons/dashboard`, ik.SourceDir),
		fmt.Sprintf(`kubectl apply -f %s/addons/metrics-server`, ik.SourceDir),
		fmt.Sprintf(`kubectl apply -f %s/addons/kube-state-metrics`, ik.SourceDir),
		fmt.Sprintf(`kubectl apply -f %s/addons/prometheus`, ik.SourceDir),
	}
	ik.er.Run(cmds...)
}

func (ik *InstallK8s) installIstio() {
	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return
	}

	ik.er.SetRole(publishRole)
	cmds := []string{
		// install istio
		fmt.Sprintf(`istioctl install --manifests=%s/istio/manifests -y`, ik.SourceDir),
		fmt.Sprintf(`istioctl manifest generate > %s/addons/istio/generated-manifest.yaml`, ik.SourceDir),
		fmt.Sprintf(`istioctl verify-install -f %s/addons/istio/generated-manifest.yaml`, ik.SourceDir),

		fmt.Sprintf(`kubectl apply -f %s/addons/istio`, ik.SourceDir),
		`kubectl -n istio-system get deployment`,

		// init gateways
		`kubectl -n istio-system delete secret k8s-com-certs`,
		// # 这个命名空间必须和ingressgateway容器服务在一起，否则加载不到证书，https站点无法访问，没有报错，被坑了很久
		fmt.Sprintf(`kubectl -n istio-system create secret tls k8s-com-certs --key=%s/addons/certs/k8s.com.key --cert=%s/addons/certs/k8s.com.crt`, ik.SourceDir, ik.SourceDir),
		fmt.Sprintf(`kubectl apply -f %s/addons/gateways`, ik.SourceDir),
	}
	ik.er.Run(cmds...)
}

func (ik *InstallK8s) initWebTest() {
	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return
	}

	ik.er.SetRole(publishRole)
	cmds := []string{
		`kubectl create namespace test-system ; kubectl label namespace test-system istio-injection=enabled`,
		fmt.Sprintf(`cd %s/web-test && kubectl create namespace test-system`, ik.SourceDir),
		fmt.Sprintf(`cd %s/web-test && kubectl apply -f yaml`, ik.SourceDir),
	}
	ik.er.Run(cmds...)
}

func (ik *InstallK8s) wait(second int, desc string) {
	for i := 1; i <= second; i++ {
		time.Sleep(1 * time.Second)
		ik.Stdout <- fmt.Sprintf("%s(%d秒)", desc, i)
	}
}
