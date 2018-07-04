kind: KubeletConfiguration
apiVersion: kubelet.config.k8s.io/v1beta1
featureGates:
  RotateKubeletServerCertificate: true
address: "NODE_HOST"
staticPodPath: "/etc/kubernetes/manifests"
clusterDNS: ["192.168.0.2"]
clusterDomain: "cluster.local"
tlsCertFile: "/etc/kubernetes/pki/kubelet.pem"
tlsPrivateKeyFile: "/etc/kubernetes/pki/kubelet-key.pem"
authorization:
  mode: Webhook
  webhook:
    cacheAuthorizedTTL: "5m0s"
    cacheUnauthorizedTTL: "30s"
authentication:
  x509:
    clientCAFile: "/etc/kubernetes/pki/ca.pem"
  webhook:
    enabled: false
    cacheTTL: "0s"
  anonymous:
    enabled: false
cgroupDriver: "cgroupfs"
tlsMinVersion: "VersionTLS12"
tlsCipherSuites:
- "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256"
- "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384"
- "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"
- "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
readOnlyPort: 0
port: 10250
# containerLogMaxSize: "10Mi"
# containerLogMaxFiles: 5
# evictionHard:
#   imagefs.available: "15%"
#   memory.available: "100Mi"
#   nodefs.available: "10%"
#   nodefs.inodesFree: "5%"
# evictionMaxPodGracePeriod: 0
# evictionPressureTransitionPeriod: "5m0s"
# fileCheckFrequency: "20s"
# imageGCHighThresholdPercent: 85
# imageGCLowThresholdPercent: 80
# maxOpenFiles: 1000000
maxPods: 300
failSwapOn: false
# imageMinimumGCAge: "2m0s"
# nodeStatusUpdateFrequency: "10s"
# runtimeRequestTimeout: "2m0s"
# streamingConnectionIdleTimeout: "4h0m0s"
# syncFrequency: "1m0s"
# volumeStatsAggPeriod: "1m0s"
