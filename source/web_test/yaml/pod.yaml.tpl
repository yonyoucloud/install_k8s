kind: Deployment
apiVersion: apps/v1
metadata:
  labels:
    app: web-test
  name: web-test
  namespace: esn-system
spec:
  replicas: 1
  minReadySeconds: 10     #滚动升级时10s后认为该pod就绪
  strategy:
    rollingUpdate:  ##由于replicas为3,则整个升级,pod个数在2-4个之间
      maxSurge: 1     #滚动升级时会先启动1个pod
      maxUnavailable: 1 #滚动升级时允许的最大Unavailable的pod个数
  selector:
    matchLabels:
      app: web-test
  template:
    metadata:
      labels:
        app: web-test
    spec:
      volumes:
        - name: accesslog
          hostPath:
            path: /data/log/nginx
        - name: errorlog
          hostPath:
            path: /data/log/php
      containers:
        - name: web-test
          #image: PRI_DOCKER_HOST:5000/esn-containers/web_test:1.0
          image: PRO_IMAGE
          resources:
            limits:
              memory: 100Mi
            requests:
              memory: 80Mi
          lifecycle:
            postStart:
              exec:
                command: ["/server/start_cmd.sh"]
            preStop:
              exec:
                #SIGTERM triggers a quick exit; gracefully terminate instead (SIGTERM触发快速退出；而优雅地终止) 将会在发送SIGTERM之前执行
                command: ["/server/stop_cmd.sh"]
          imagePullPolicy: Always
          ports:
          - containerPort: 80
            protocol: TCP
          livenessProbe: #(活性探针)
            tcpSocket:
              port: 80
            initialDelaySeconds: 6
            periodSeconds: 3
          readinessProbe:
            tcpSocket:
              port: 80
            initialDelaySeconds: 4
            periodSeconds: 2
          volumeMounts:
            - mountPath: /data/log/nginx
              name: accesslog
            - mountPath: /data/log/php
              name: errorlog
      #nodeSelector:
        #othersvc: othersys
