[Unit]
Description=Etcd Server
After=network.target
After=network-online.target
Wants=network-online.target

[Service]
Type=notify
WorkingDirectory=/data/etcd
User=etcd
# set GOMAXPROCS to number of processors
ExecStart=/usr/bin/etcd --name ETCD_NAME --data-dir /data/etcd --listen-client-urls https://0.0.0.0:2379 --listen-peer-urls https://ETCD_HOST:2380 --advertise-client-urls https://ETCD_HOST:2379 --initial-cluster-token etcd-cluster --initial-cluster ETCD_INITIAL_CLUSTER --initial-cluster-state ETCD_INITIAL_CLUSTER_STATE --initial-advertise-peer-urls https://ETCD_HOST:2380 --cert-file /etc/etcd/ssl/etcd.pem --key-file /etc/etcd/ssl/etcd-key.pem --peer-cert-file /etc/etcd/ssl/etcd.pem --peer-key-file /etc/etcd/ssl/etcd-key.pem --trusted-ca-file /etc/etcd/ssl/ca.pem --peer-trusted-ca-file /etc/etcd/ssl/ca.pem --logger=zap
Restart=on-failure
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target

