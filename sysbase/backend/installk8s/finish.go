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
	ik.Stdout <- "启动Publish"

	ik.Params.DoWhat = "start"
	ik.Params.DoDocker = true
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

	ik.Stdout <- "启动Dns"
	ik.Params.DoWhat = "start"
	ik.ServiceDns()
	ik.wait(5, "等待...")

	ik.Stdout <- "[结束]启动所有服务"

	ik.Stdout <- "[开始]验证k8s集群"
	publishRole.WaitOutput = true
	ik.er.SetRole(publishRole)
	for i := 1; i <= 180; i++ {
		ik.Stdout <- fmt.Sprintf("等待倒计时(%d)s...", i)
		res := ik.er.Run("kubectl get nodes -o wide | grep NotReady > /dev/null ; echo $?")
		if strings.Contains(strings.Join(res, ""), "-> 1") {
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
		res := ik.er.Run("kubectl -n kube-system get pods -o wide | grep kubernetes-dashboard | grep Running > /dev/null ; echo $?")
		if strings.Contains(strings.Join(res, ""), "-> 0") {
			publishRole.WaitOutput = false
			ik.er.SetRole(publishRole)
			ik.er.Run(`kubectl -n istio-system get pod -o wide | grep istio-ingressgateway | grep Running | awk '{print "设置Hosts: "$7" dashboard.k8s.com 然后您可以访问kubernetes-dashboard: https://dashboard.k8s.com:10443"}'`)
			ik.Stdout <- "用下面输出的token登录kubernetes-dashboard"
			ik.er.Run(`kubectl describe secret $(kubectl get secret -n kube-system | grep admin-token | awk '{print $1}') -n kube-system | grep token: | awk '{print $2}'`)
			break
		}
		time.Sleep(1 * time.Second)
	}

	publishRole.WaitOutput = true
	ik.er.SetRole(publishRole)
	for i := 1; i <= 180; i++ {
		ik.Stdout <- fmt.Sprintf("等待monitoring-grafana running(%d)s...", i)
		res := ik.er.Run("kubectl -n kube-system get pods -o wide | grep monitoring-grafana | grep Running > /dev/null ; echo $?")
		if strings.Contains(strings.Join(res, ""), "-> 0") {
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
		res := ik.er.Run("kubectl -n kube-system get pods -o wide | grep prometheus | grep Running > /dev/null ; echo $?")
		if strings.Contains(strings.Join(res, ""), "-> 0") {
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
		res := ik.er.Run("kubectl -n istio-system get pods -o wide | grep kiali | grep Running > /dev/null ; echo $?")
		if strings.Contains(strings.Join(res, ""), "-> 0") {
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
		res := ik.er.Run("kubectl -n test-system get pods -o wide | grep web-test | grep Running > /dev/null ; echo $?")
		if strings.Contains(strings.Join(res, ""), "-> 0") {
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
	res := ik.er.Run("kubectl -n test-system get pods | grep web-test | awk '{print $1}'")
	resArr := strings.Split(strings.Join(res, ""), "-> ")
	pod := strings.Trim(resArr[1], "\n")

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

	pridockerRole, ok := ik.resources["pridocker"]
	if !ok {
		ik.Stdout <- "没有pridocker资源"
		return
	}

	priDockerHost := strings.Split(pridockerRole.Hosts[0], ":")[0]

	ik.er.SetRole(publishRole)
	cmds := []string{
		fmt.Sprintf("docker images | grep 'test-base' || (cd %s/images && sha256=`docker load -i test-containers~test-base:1.0.tar | grep Loaded | cut -f 3,4 -d ' ' | cut -f 2 -d ' ' | sed 's/sha256://g'` && docker tag $sha256 %s:5000/test-containers/test-base:1.0)", ik.SourceDir, priDockerHost),

		fmt.Sprintf("docker images | grep 'pause-amd64' || (cd %s/images && sha256=`docker load -i gcr.io~google_containers~pause-amd64:3.2.tar | grep Loaded | cut -f 3,4 -d ' ' | cut -f 2 -d ' ' | sed 's/sha256://g'` && docker tag $sha256 %s:5000/gcr.io/google_containers/pause-amd64:3.2 && docker push %s:5000/gcr.io/google_containers/pause-amd64:3.2)", ik.SourceDir, priDockerHost, priDockerHost),

		fmt.Sprintf("docker images | grep 'kubernetesui/dashboard' || (cd %s/images && sha256=`docker load -i kubernetesui~dashboard:v2.0.1.tar | grep Loaded | cut -f 3,4 -d ' ' | cut -f 2 -d ' ' | sed 's/sha256://g'` && docker tag $sha256 %s:5000/kubernetesui/dashboard:v2.0.1 && docker push %s:5000/kubernetesui/dashboard:v2.0.1)", ik.SourceDir, priDockerHost, priDockerHost),

		fmt.Sprintf("docker images | grep 'kubernetesui/metrics-scraper' || (cd %s/images && sha256=`docker load -i kubernetesui~metrics-scraper:v1.0.4.tar | grep Loaded | cut -f 3,4 -d ' ' | cut -f 2 -d ' ' | sed 's/sha256://g'` && docker tag $sha256 %s:5000/kubernetesui/metrics-scraper:v1.0.4 && docker push %s:5000/kubernetesui/metrics-scraper:v1.0.4)", ik.SourceDir, priDockerHost, priDockerHost),

		fmt.Sprintf("docker images | grep 'busybox' || (cd %s/images && sha256=`docker load -i busybox:latest.tar | grep Loaded | cut -f 3,4 -d ' ' | cut -f 2 -d ' ' | sed 's/sha256://g'` && docker tag $sha256 %s:5000/busybox:latest && docker push %s:5000/busybox:latest)", ik.SourceDir, priDockerHost, priDockerHost),

		fmt.Sprintf("docker images | grep 'grafana' || (cd %s/images && sha256=`docker load -i grafana~grafana:7.4.3.tar | grep Loaded | cut -f 3,4 -d ' ' | cut -f 2 -d ' ' | sed 's/sha256://g'` && docker tag $sha256 %s:5000/grafana/grafana:7.4.3 && docker push %s:5000/grafana/grafana:7.4.3)", ik.SourceDir, priDockerHost, priDockerHost),

		fmt.Sprintf("docker images | grep 'jimmidyson/configmap-reload' || (cd %s/images && sha256=`docker load -i jimmidyson~configmap-reload:v0.5.0.tar | grep Loaded | cut -f 3,4 -d ' ' | cut -f 2 -d ' ' | sed 's/sha256://g'` && docker tag $sha256 %s:5000/jimmidyson/configmap-reload:v0.5.0 && docker push %s:5000/jimmidyson/configmap-reload:v0.5.0)", ik.SourceDir, priDockerHost, priDockerHost),

		fmt.Sprintf("docker images | grep 'addon-resizer' || (cd %s/images && sha256=`docker load -i k8s.gcr.io~addon-resizer:1.8.11.tar | grep Loaded | cut -f 3,4 -d ' ' | cut -f 2 -d ' ' | sed 's/sha256://g'` && docker tag $sha256 %s:5000/k8s.gcr.io/addon-resizer:1.8.11 && docker push %s:5000/k8s.gcr.io/addon-resizer:1.8.11)", ik.SourceDir, priDockerHost, priDockerHost),

		fmt.Sprintf("docker images | grep 'prom/alertmanager' || (cd %s/images && sha256=`docker load -i prom~alertmanager:v0.21.0.tar | grep Loaded | cut -f 3,4 -d ' ' | cut -f 2 -d ' ' | sed 's/sha256://g'` && docker tag $sha256 %s:5000/prom/alertmanager:v0.21.0 && docker push %s:5000/prom/alertmanager:v0.21.0)", ik.SourceDir, priDockerHost, priDockerHost),

		fmt.Sprintf("docker images | grep 'prom/node-exporter' || (cd %s/images && sha256=`docker load -i prom~node-exporter:v1.1.2.tar | grep Loaded | cut -f 3,4 -d ' ' | cut -f 2 -d ' ' | sed 's/sha256://g'` && docker tag $sha256 %s:5000/prom/node-exporter:v1.1.2 && docker push %s:5000/prom/node-exporter:v1.1.2)", ik.SourceDir, priDockerHost, priDockerHost),

		fmt.Sprintf("docker images | grep 'prom/prometheus' || (cd %s/images && sha256=`docker load -i prom~prometheus:v2.25.0.tar | grep Loaded | cut -f 3,4 -d ' ' | cut -f 2 -d ' ' | sed 's/sha256://g'` && docker tag $sha256 %s:5000/prom/prometheus:v2.25.0 && docker push %s:5000/prom/prometheus:v2.25.0)", ik.SourceDir, priDockerHost, priDockerHost),

		fmt.Sprintf("docker images | grep 'kube-state-metrics' || (cd %s/images && sha256=`docker load -i quay.io~coreos~kube-state-metrics:v1.9.8.tar | grep Loaded | cut -f 3,4 -d ' ' | cut -f 2 -d ' ' | sed 's/sha256://g'` && docker tag $sha256 %s:5000/quay.io/coreos/kube-state-metrics:v1.9.8 && docker push %s:5000/quay.io/coreos/kube-state-metrics:v1.9.8)", ik.SourceDir, priDockerHost, priDockerHost),

		fmt.Sprintf("docker images | grep 'metrics-server-amd64' || (cd %s/images && sha256=`docker load -i k8s.gcr.io~metrics-server-amd64:v0.3.6.tar | grep Loaded | cut -f 3,4 -d ' ' | cut -f 2 -d ' ' | sed 's/sha256://g'` && docker tag $sha256 %s:5000/k8s.gcr.io/metrics-server-amd64:v0.3.6 && docker push %s:5000/k8s.gcr.io/metrics-server-amd64:v0.3.6)", ik.SourceDir, priDockerHost, priDockerHost),

		fmt.Sprintf("docker images | grep 'calico/cni' || (cd %s/images && sha256=`docker load -i calico~cni:v3.18.0.tar | grep Loaded | cut -f 3,4 -d ' ' | cut -f 2 -d ' ' | sed 's/sha256://g'` && docker tag $sha256 %s:5000/calico/cni:v3.18.0 && docker push %s:5000/calico/cni:v3.18.0)", ik.SourceDir, priDockerHost, priDockerHost),

		fmt.Sprintf("docker images | grep 'calico/kube-controllers' || (cd %s/images && sha256=`docker load -i calico~kube-controllers:v3.18.0.tar | grep Loaded | cut -f 3,4 -d ' ' | cut -f 2 -d ' ' | sed 's/sha256://g'` && docker tag $sha256 %s:5000/calico/kube-controllers:v3.18.0 && docker push %s:5000/calico/kube-controllers:v3.18.0)", ik.SourceDir, priDockerHost, priDockerHost),

		fmt.Sprintf("docker images | grep 'calico/node' || (cd %s/images && sha256=`docker load -i calico~node:v3.18.0.tar | grep Loaded | cut -f 3,4 -d ' ' | cut -f 2 -d ' ' | sed 's/sha256://g'` && docker tag $sha256 %s:5000/calico/node:v3.18.0 && docker push %s:5000/calico/node:v3.18.0)", ik.SourceDir, priDockerHost, priDockerHost),

		fmt.Sprintf("docker images | grep 'calico/pod2daemon-flexvol' || (cd %s/images && sha256=`docker load -i calico~pod2daemon-flexvol:v3.18.0.tar | grep Loaded | cut -f 3,4 -d ' ' | cut -f 2 -d ' ' | sed 's/sha256://g'` && docker tag $sha256 %s:5000/calico/pod2daemon-flexvol:v3.18.0 && docker push %s:5000/calico/pod2daemon-flexvol:v3.18.0)", ik.SourceDir, priDockerHost, priDockerHost),

		fmt.Sprintf("docker images | grep 'coredns' || (cd %s/images && sha256=`docker load -i k8s.gcr.io~coredns:1.7.0.tar | grep Loaded | cut -f 3,4 -d ' ' | cut -f 2 -d ' ' | sed 's/sha256://g'` && docker tag $sha256 %s:5000/k8s.gcr.io/coredns:1.7.0 && docker push %s:5000/k8s.gcr.io/coredns:1.7.0)", ik.SourceDir, priDockerHost, priDockerHost),

		fmt.Sprintf("docker images | grep 'istio/pilot' || (cd %s/images && sha256=`docker load -i istio~pilot:1.9.1.tar | grep Loaded | cut -f 3,4 -d ' ' | cut -f 2 -d ' ' | sed 's/sha256://g'` && docker tag $sha256 %s:5000/istio/pilot:1.9.1 && docker push %s:5000/istio/pilot:1.9.1)", ik.SourceDir, priDockerHost, priDockerHost),

		fmt.Sprintf("docker images | grep 'istio/proxyv2' || (cd %s/images && sha256=`docker load -i istio~proxyv2:1.9.1.tar | grep Loaded | cut -f 3,4 -d ' ' | cut -f 2 -d ' ' | sed 's/sha256://g'` && docker tag $sha256 %s:5000/istio/proxyv2:1.9.1 && docker push %s:5000/istio/proxyv2:1.9.1)", ik.SourceDir, priDockerHost, priDockerHost),

		fmt.Sprintf("docker images | grep 'quay.io/kiali/kiali' || (cd %s/images && sha256=`docker load -i quay.io~kiali~kiali:v1.29.tar | grep Loaded | cut -f 3,4 -d ' ' | cut -f 2 -d ' ' | sed 's/sha256://g'` && docker tag $sha256 %s:5000/quay.io/kiali/kiali:v1.29 && docker push %s:5000/quay.io/kiali/kiali:v1.29)", ik.SourceDir, priDockerHost, priDockerHost),

		fmt.Sprintf("docker images | grep 'jaegertracing/all-in-one' || (cd %s/images && sha256=`docker load -i jaegertracing~all-in-one:1.20.tar | grep Loaded | cut -f 3,4 -d ' ' | cut -f 2 -d ' ' | sed 's/sha256://g'` && docker tag $sha256 %s:5000/jaegertracing/all-in-one:1.20 && docker push %s:5000/jaegertracing/all-in-one:1.20)", ik.SourceDir, priDockerHost, priDockerHost),

		fmt.Sprintf(`docker images | grep %s:5000 | awk '{print "docker push "$1":"$2}' | sh`, priDockerHost),
	}
	ik.er.Run(cmds...)
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

	pridockerRole, ok := ik.resources["pridocker"]
	if !ok {
		ik.Stdout <- "没有pridocker资源"
		return
	}

	nodeRole, ok := ik.resources["node"]
	if !ok {
		ik.Stdout <- "没有node资源"
		return
	}

	etcdLbHost := strings.Split(etcdLbRole.Hosts[0], ":")[0]
	priDockerHost := strings.Split(pridockerRole.Hosts[0], ":")[0]

	ik.er.SetRole(publishRole)
	cmds := []string{
		fmt.Sprintf("kubectl delete -f %s/calico", ik.SourceDir),
		fmt.Sprintf(`sed "s#PRI_DOCKER_HOST#%s#g" %s/calico/calico.yaml.tpl > %s/calico/calico.yaml`, priDockerHost, ik.SourceDir, ik.SourceDir),
		fmt.Sprintf(`sed -i "s#ETCD_LVS_HOST#%s#g" %s/calico/calico.yaml`, etcdLbHost, ik.SourceDir),
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
		res := ik.er.Run("kubectl get pods -o wide -n kube-system | grep calico | grep Running | wc -l")
		resArr := strings.Split(strings.Join(res, ""), "-> ")
		numStr := strings.Trim(resArr[1], "\n")
		num, _ := strconv.Atoi(numStr)
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

	pridockerRole, ok := ik.resources["pridocker"]
	if !ok {
		ik.Stdout <- "没有pridocker资源"
		return
	}

	pridnsRole, ok := ik.resources["pridns"]
	if !ok {
		ik.Stdout <- "没有pridns资源"
		return
	}

	priDockerHost := strings.Split(pridockerRole.Hosts[0], ":")[0]
	pridnsHost := strings.Split(pridnsRole.Hosts[0], ":")[0]

	ik.er.SetRole(publishRole)
	cmds := []string{
		fmt.Sprintf(`sed "s#PRI_DOCKER_HOST#%s#g" %s/dashboard/dashboard.yaml.tpl > %s/dashboard/dashboard.yaml`, priDockerHost, ik.SourceDir, ik.SourceDir),
		fmt.Sprintf(`sed "s#PRI_DOCKER_HOST#%s#g" %s/dns/coredns.yaml.tpl > %s/dns/coredns.yaml`, priDockerHost, ik.SourceDir, ik.SourceDir),
		fmt.Sprintf(`sed "s#PRI_DOCKER_HOST#%s#g" %s/metrics-server/metrics-server-deployment.yaml.tpl > %s/metrics-server/metrics-server-deployment.yaml`, priDockerHost, ik.SourceDir, ik.SourceDir),
		fmt.Sprintf(`sed "s#PRI_DOCKER_HOST#%s#g" %s/prometheus/alertmanager-deployment.yaml.tpl > %s/prometheus/alertmanager-deployment.yaml`, priDockerHost, ik.SourceDir, ik.SourceDir),
		fmt.Sprintf(`sed "s#PRI_DOCKER_HOST#%s#g" %s/prometheus/grafana.yaml.tpl > %s/prometheus/grafana.yaml`, priDockerHost, ik.SourceDir, ik.SourceDir),
		fmt.Sprintf(`sed "s#PRI_DOCKER_HOST#%s#g" %s/prometheus/kube-state-metrics-deployment.yaml.tpl > %s/prometheus/kube-state-metrics-deployment.yaml`, priDockerHost, ik.SourceDir, ik.SourceDir),
		fmt.Sprintf(`sed "s#PRI_DOCKER_HOST#%s#g" %s/prometheus/node-exporter-ds.yaml.tpl > %s/prometheus/node-exporter-ds.yaml`, priDockerHost, ik.SourceDir, ik.SourceDir),
		fmt.Sprintf(`sed "s#PRI_DOCKER_HOST#%s#g" %s/prometheus/prometheus-statefulset.yaml.tpl > %s/prometheus/prometheus-statefulset.yaml`, priDockerHost, ik.SourceDir, ik.SourceDir),
		fmt.Sprintf(`sed -i "s#HOST#%s#g" %s/dns/coredns.yaml`, pridnsHost, ik.SourceDir),
		fmt.Sprintf(`rm -rf %s/dashboard/certs ; mkdir -p %s/dashboard/certs`, ik.SourceDir, ik.SourceDir),
		fmt.Sprintf(`openssl req -out %s/dashboard/certs/dashboard.csr -newkey rsa:2048 -nodes -keyout %s/dashboard/certs/dashboard.key -subj "/CN=dashboard.k8s.com/O=dashboard organization"`, ik.SourceDir, ik.SourceDir),
		fmt.Sprintf(`openssl x509 -req -days 36500 -CA %s/master/etc/kubernetes/pki/ca.pem -CAkey %s/master/etc/kubernetes/pki/ca-key.pem -CAcreateserial -in %s/dashboard/certs/dashboard.csr -out %s/dashboard/certs/dashboard.crt`, ik.SourceDir, ik.SourceDir, ik.SourceDir, ik.SourceDir),
		`kubectl -n kube-system delete secret kubernetes-dashboard-certs`,
		fmt.Sprintf(`kubectl -n kube-system create secret generic kubernetes-dashboard-certs --from-file=%s/dashboard/certs`, ik.SourceDir),
		fmt.Sprintf(`kubectl apply -f %s/dashboard`, ik.SourceDir),
		fmt.Sprintf(`kubectl apply -f %s/dns`, ik.SourceDir),
		fmt.Sprintf(`kubectl apply -f %s/metrics-server`, ik.SourceDir),
		fmt.Sprintf(`kubectl apply -f %s/prometheus`, ik.SourceDir),
	}
	ik.er.Run(cmds...)
}

func (ik *InstallK8s) installIstio() {
	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return
	}

	pridockerRole, ok := ik.resources["pridocker"]
	if !ok {
		ik.Stdout <- "没有pridocker资源"
		return
	}

	priDockerHost := strings.Split(pridockerRole.Hosts[0], ":")[0]

	ik.er.SetRole(publishRole)
	cmds := []string{
		fmt.Sprintf(`cd %s/istio/manifests/profiles && sed "s#PRI_DOCKER_HOST#%s#g" default.yaml.tpl > default.yaml`, ik.SourceDir, priDockerHost),
		fmt.Sprintf(`cd %s/istio/samples/addons && sed "s#PRI_DOCKER_HOST#%s#g" kiali.yaml.tpl > kiali.yaml`, ik.SourceDir, priDockerHost),
		fmt.Sprintf(`cd %s/istio/samples/addons && sed "s#PRI_DOCKER_HOST#%s#g" jaeger.yaml.tpl > jaeger.yaml`, ik.SourceDir, priDockerHost),
		fmt.Sprintf(`istioctl manifest generate > %s/istio/generated-manifest.yaml`, ik.SourceDir),
		fmt.Sprintf(`istioctl install --manifests=%s/istio/manifests -y`, ik.SourceDir),
		`kubectl -n istio-system get deployment`,
		fmt.Sprintf(`istioctl verify-install -f %s/istio/generated-manifest.yaml`, ik.SourceDir),
		fmt.Sprintf(`kubectl apply -f %s/istio/samples/addons/kiali.yaml`, ik.SourceDir),
		fmt.Sprintf(`kubectl apply -f %s/istio/samples/addons/jaeger.yaml`, ik.SourceDir),
		fmt.Sprintf(`rm -rf %s/istio/samples/addons/gateways/certs ; mkdir -p %s/istio/samples/addons/gateways/certs`, ik.SourceDir, ik.SourceDir),
		fmt.Sprintf(`openssl req -x509 -sha256 -nodes -days 36500 -newkey rsa:2048 -subj '/O=k8s Inc./CN=k8s.com' -keyout %s/istio/samples/addons/gateways/certs/k8s.com.key -out %s/istio/samples/addons/gateways/certs/k8s.com.crt`, ik.SourceDir, ik.SourceDir),
		`kubectl -n istio-system delete secret k8s-com-certs`,
		// # 这个命名空间必须和ingressgateway容器服务在一起，否则加载不到证书，https站点无法访问，没有报错，被坑了很久
		fmt.Sprintf(`kubectl -n istio-system create secret tls k8s-com-certs --key=%s/istio/samples/addons/gateways/certs/k8s.com.key --cert=%s/istio/samples/addons/gateways/certs/k8s.com.crt`, ik.SourceDir, ik.SourceDir),
		fmt.Sprintf(`kubectl apply -f %s/istio/samples/addons/gateways`, ik.SourceDir),
	}
	ik.er.Run(cmds...)
}

func (ik *InstallK8s) initWebTest() {
	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return
	}

	pridockerRole, ok := ik.resources["pridocker"]
	if !ok {
		ik.Stdout <- "没有pridocker资源"
		return
	}

	priDockerHost := strings.Split(pridockerRole.Hosts[0], ":")[0]

	ik.er.SetRole(publishRole)
	cmds := []string{
		`kubectl create namespace test-system ; kubectl label namespace test-system istio-injection=enabled`,
		fmt.Sprintf(`cd %s/web-test && sed "s#PRI_DOCKER_HOST#%s#g" Dockerfile.tpl > Dockerfile`, ik.SourceDir, priDockerHost),
		fmt.Sprintf(`cd %s/web-test && sed "s#PRI_DOCKER_HOST#%s#g" create.sh.tpl > create.sh && chmod 750 create.sh`, ik.SourceDir, priDockerHost),
		fmt.Sprintf(`cd %s/web-test && ./create.sh`, ik.SourceDir),
	}
	ik.er.Run(cmds...)
}

func (ik *InstallK8s) wait(second int, desc string) {
	for i := 1; i <= second; i++ {
		time.Sleep(1 * time.Second)
		ik.Stdout <- fmt.Sprintf("%s(%d秒)", desc, i)
	}
}
