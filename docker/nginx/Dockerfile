FROM ubuntu:14.04

RUN apt-get -y update && apt-get -y install \
  curl \
  wget

# https://www.nginx.com/resources/admin-guide/installing-nginx-open-source/
RUN apt-get -y install build-essential
RUN apt-get -y install libpcre3 libpcre3-dev zlib1g-dev libssl-dev

ADD headers-more-nginx-module /headers-more-nginx-module

RUN wget http://nginx.org/download/nginx-1.10.2.tar.gz
RUN tar zxf nginx-1.10.2.tar.gz && cd nginx-1.10.2 && ./configure --prefix=/etc/nginx --sbin-path=/usr/sbin/nginx --conf-path=/etc/nginx/nginx.conf --error-log-path=/var/log/nginx/error.log --http-log-path=/var/log/nginx/access.log --pid-path=/var/run/nginx.pid --lock-path=/var/run/nginx.lock --add-module=/headers-more-nginx-module && make && sudo make install

EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]
