module updatesecret

go 1.15

replace updatesecret => ./updatesecret

require (
	github.com/spf13/pflag v1.0.5
	gopkg.in/square/go-jose.v2 v2.5.1
	k8s.io/api v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/apiserver v0.20.4
	k8s.io/client-go v0.20.4
	k8s.io/component-base v0.20.4
	k8s.io/klog/v2 v2.4.0
// k8s.io/kubernetes v1.20.4
)
