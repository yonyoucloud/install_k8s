package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/keyutil"

	"updatesecret/serviceaccount"
)

var (
	kubeConfig     = pflag.String("kubeconfig", "~/.kube/config", "config认证文件路径")
	privateKeyFile = pflag.String("private-key-file", "/etc/kubernetes/pki/ca-key.pem", "CA证书私钥文件路径")

	rsaPrivateKey string
)

func init() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	_ = flag.CommandLine.Parse(make([]string, 0)) // Init for glog calls in kubernetes packages
}

func main() {
	if home := homedir.HomeDir(); home != "" {
		*kubeConfig = filepath.Join(home, ".kube", "config")
	}

	if len(*kubeConfig) == 0 {
		fmt.Println("必须指定config认证文件路径")
		return
	}

	if !exists(*kubeConfig) {
		fmt.Printf("%s 文件路径不存在\n", *kubeConfig)
		return
	}

	if len(*privateKeyFile) == 0 {
		fmt.Println("必须指定CA证书私钥文件路径")
		return
	}

	if !exists(*privateKeyFile) {
		fmt.Printf("%s 文件路径不存在\n", *privateKeyFile)
		return
	}

	data, err := ioutil.ReadFile(*privateKeyFile)
	if err != nil {
		panic(err.Error())
	}
	rsaPrivateKey = string(data)

	if rsaPrivateKey == "" {
		fmt.Println("CA证书文件内容不能为空")
		return
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	secrets, err := clientset.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, secret := range secrets.Items {
		if secret.Type != "kubernetes.io/service-account-token" {
			continue
		}

		// if secret.Name != "calico-node-token-lv27t" {
		// 	continue
		// }

		serviceAccount, err := clientset.CoreV1().ServiceAccounts(secret.Namespace).Get(context.TODO(), secret.GetAnnotations()["kubernetes.io/service-account.name"], metav1.GetOptions{})
		if err != nil {
			fmt.Printf("获取 Secret[%s/%s]ServiceAccount 对象报错: %s\n", secret.Namespace, secret.Name, err.Error())
			continue
		}

		rsaToken, err := generateToken(serviceAccount, secret)
		if err != nil {
			fmt.Printf("生成 Secret[%s/%s]Token 报错: %s\n", secret.Namespace, secret.Name, err.Error())
			continue
		}

		secret.Data["token"] = []byte(rsaToken)

		_, err = clientset.CoreV1().Secrets(secret.Namespace).Update(context.TODO(), &secret, metav1.UpdateOptions{})
		if err != nil {
			fmt.Printf("更新 Secret[%s/%s]Token 报错: %s\n", secret.Namespace, secret.Name, err.Error())
			continue
		}

		fmt.Printf("更新 Secret[%s/%s]Token 更新成功\n", secret.Namespace, secret.Name)
	}
}

func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func getPrivateKey(data string) interface{} {
	key, err := keyutil.ParsePrivateKeyPEM([]byte(data))
	if err != nil {
		panic(fmt.Errorf("unexpected error parsing private key: %v", err))
	}
	return key
}

func generateToken(serviceAccount *v1.ServiceAccount, secret v1.Secret) (string, error) {
	rsaGenerator, err := serviceaccount.JWTTokenGenerator(serviceaccount.LegacyIssuer, getPrivateKey(rsaPrivateKey))
	if err != nil {
		return "", err
	}
	rsaToken, err := rsaGenerator.GenerateToken(serviceaccount.LegacyClaims(*serviceAccount, secret))
	if err != nil {
		return "", err
	}

	return rsaToken, nil
}
