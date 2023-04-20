module updatesecret

go 1.15

replace updatesecret => ./updatesecret

require (
	github.com/spf13/pflag v1.0.5
	gopkg.in/square/go-jose.v2 v2.6.0
	k8s.io/api v0.27.1
	k8s.io/apimachinery v0.27.1
	k8s.io/apiserver v0.27.1
	k8s.io/client-go v0.27.1
	k8s.io/component-base v0.27.1
	k8s.io/klog/v2 v2.90.1
)
