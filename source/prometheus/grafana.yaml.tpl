apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: monitoring-grafana
  namespace: kube-system
  labels:
    k8s-app: grafana
    kubernetes.io/cluster-service: "true"
    addonmanager.kubernetes.io/mode: Reconcile
    version: v2.2.1
spec:
  serviceName: "grafana"
  replicas: 1
  podManagementPolicy: "Parallel"
  updateStrategy:
   type: "RollingUpdate"
  selector:
    matchLabels:
      k8s-app: grafana
  template:
    metadata:
      labels:
        k8s-app: grafana
      annotations:
        scheduler.alpha.kubernetes.io/critical-pod: ''
    spec:
      priorityClassName: system-cluster-critical
      #serviceAccountName: grafana
      initContainers:
      - name: "init-chown-data"
        image: "PRI_DOCKER_HOST:5000/busybox:latest"
        imagePullPolicy: "IfNotPresent"
        command: ["chown", "-R", "472:472", "/var/lib/grafana"]
        volumeMounts:
        - name: grafana-data
          mountPath: /var/lib/grafana
          subPath: "grafana"
      containers:
      - name: grafana
        image: PRI_DOCKER_HOST:5000/grafana/grafana:6.4.4.1
        imagePullPolicy: "IfNotPresent"
        ports:
        - containerPort: 3000
          protocol: TCP
        volumeMounts:
        - mountPath: /etc/ssl/certs
          name: ca-certificates
          readOnly: true
        - name: grafana-data
          mountPath: "/var/lib/grafana"
          subPath: "grafana"
        env:
        - name: GF_SECURITY_ADMIN_USER
          value: admin
        - name: GF_SECURITY_ADMIN_PASSWORD
          value: sMkJslfjHDSF19AG
        - name: INFLUXDB_HOST
          value: monitoring-influxdb
        - name: GF_SERVER_HTTP_PORT
          value: "3000"
          # The following env variables are required to make Grafana accessible via
          # the kubernetes api-server proxy. On production clusters, we recommend
          # removing these env variables, setup auth for grafana, and expose the grafana
          # service using a LoadBalancer or a public IP.
        - name: GF_AUTH_BASIC_ENABLED
          value: "false"
        - name: GF_AUTH_ANONYMOUS_ENABLED
          value: "false"
        - name: GF_AUTH_ANONYMOUS_ORG_ROLE
          value: Admin
        - name: GF_SERVER_ROOT_URL
          # If you're only using the API Server proxy, set this value instead:
          # value: /api/v1/namespaces/kube-system/services/monitoring-grafana/proxy
          value: /
      volumes:
      - name: ca-certificates
        hostPath:
          path: /etc/ssl/certs
      - name: grafana-data
        persistentVolumeClaim:
          claimName: grafana
      #securityContext:
      #  fsGroup: 472
      #  runAsUser: 472
      #nodeSelector:
      #  monitor: prometheus
  volumeClaimTemplates:
  - metadata:
      name: grafana-data
    spec:
      storageClassName: grafana
      accessModes:
        - ReadWriteOnce
      resources:
        requests:
          storage: "2Gi"
---
apiVersion: v1
kind: Service
metadata:
  labels:
    # For use as a Cluster add-on (https://github.com/kubernetes/kubernetes/tree/master/cluster/addons)
    # If you are NOT using this as an addon, you should comment out this line.
    kubernetes.io/cluster-service: 'true'
    kubernetes.io/name: monitoring-grafana
  name: monitoring-grafana
  namespace: kube-system
spec:
  # In a production setup, we recommend accessing Grafana through an external Loadbalancer
  # or through a public IP.
  # type: LoadBalancer
  # You could also use NodePort to expose the service at a randomly-generated port
  type: NodePort
  ports:
  - port: 80
    targetPort: 3000
    nodePort: 30001
  selector:
    k8s-app: grafana
