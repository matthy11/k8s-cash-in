#! /bin/bash
yum update -y
yum install -y nginx
systemctl enable nginx
systemctl start nginx
rm /etc/nginx/conf.d/default.conf
cat <<EOF > /etc/nginx/nginx.conf
user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log;
pid /run/nginx.pid;

# Load dynamic modules. See /usr/share/nginx/README.dynamic.
include /usr/share/nginx/modules/*.conf;

events {
    worker_connections 1024;
}

http {
    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;

    sendfile            on;
    tcp_nopush          on;
    tcp_nodelay         on;
    keepalive_timeout   65;
    types_hash_max_size 2048;

    include             /etc/nginx/mime.types;
    default_type        application/octet-stream;

    include /etc/nginx/conf.d/*.conf;

}
EOF
cat <<EOF > /etc/nginx/conf.d/reverse-proxy.conf
server {
    listen 80;
    listen [::]:80;

    location / {
        deny  all;
    }
    location /depositValidations {
        proxy_pass http://34.102.193.139/deposit-validations;
    }
    location /receivedTransfers {
        proxy_pass http://34.102.193.139/received-transfers;
    }
    location /reversedTransfers {
        proxy_pass http://34.102.193.139/reversed-transfers;
    }
    location /accessTokens {
        proxy_pass http://34.102.193.139/access-tokens;
    }
}
EOF
systemctl restart nginx