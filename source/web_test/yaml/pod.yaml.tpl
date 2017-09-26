kind: Deployment
apiVersion: extensions/v1beta1
metadata:
  labels:
    app: web-test
  name: web-test
  namespace: esn-system
spec:
  replicas: 1
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
          #livenessProbe: #(活性探针)
            #httpGet:
              #path: /
              #port: 80
            #initialDelaySeconds: 30
            #timeoutSeconds: 30
          volumeMounts:
            - mountPath: /data/log/nginx
              name: accesslog
            - mountPath: /data/log/php
              name: errorlog
      #nodeSelector:
        #othersvc: othersys
