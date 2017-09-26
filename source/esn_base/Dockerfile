# docker build -t="esn-containers/esn_base:1.0" .

FROM esn-containers/alpine:latest

LABEL maintainer="shenghua bi <net.bsh@gmail.com>"

ENV TZ 'Asia/Shanghai'

RUN apk upgrade --no-cache && \
    apk add --no-cache --update bash tzdata git openssh shadow && \
    mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    rm -rf /tmp/* /var/cache/apk/* && \
    rm -rf /etc/ssh/ssh_host_rsa_key /etc/ssh/ssh_host_dsa_key && \
    sed -i "s#umask 022#umask 027#g" /etc/profile && source /etc/profile && \
    mkdir -p /data && chmod 755 /data && useradd -m -d /data/www esn && useradd -m -d /data/www www && chown -R esn:esn /data/www && \
    usermod -G esn www && chmod 750 /data/www && \
    mkdir -p /data/log /data/sess /data/yy_log /data/tmp && \
    mkdir -p /data/log/nginx /data/log/php && \
    chown -R www:www /data/log /data/sess /data/yy_log /data/tmp && \
    chmod 750 /data/www /data//log /data/sess /data/yy_log /data/tmp

COPY ssh/root /root/.ssh
COPY ssh/esn /data/www/.ssh
COPY docker-entrypoint.sh /usr/local/bin

RUN chown -R esn:esn /data/www/.ssh && \
    chmod 750 /data/www/.ssh

EXPOSE 22
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["/usr/sbin/sshd", "-D"]
