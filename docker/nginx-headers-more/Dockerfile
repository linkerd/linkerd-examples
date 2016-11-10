FROM ubuntu:14.04

RUN apt-get -y update
RUN apt-get install -y curl

# install nginx and third party nginx modules (we need the Headers More module)
RUN apt-get -y install \
  nginx \
  nginx-extras

# forward request and error logs to docker log collector
RUN ln -sf /dev/stdout /var/log/nginx/access.log
RUN ln -sf /dev/stderr /var/log/nginx/error.log

ADD ../nginx/nginx.conf /etc/nginx/

EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]