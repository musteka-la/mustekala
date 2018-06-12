FROM jwilder/nginx-proxy

RUN { \
      echo 'proxy_connect_timeout 86400;'; \
      echo 'proxy_send_timeout 86400;'; \
      echo 'proxy_read_timeout 86400;'; \
      echo 'send_timeout 86400;'; \
    } > /etc/nginx/conf.d/extended_timeout.conf
