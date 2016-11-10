FROM nginx

RUN apt-get update
RUN apt-get install -y curl

RUN rm /etc/nginx/conf.d/default.conf

ADD nginx.conf /etc/nginx/

EXPOSE 80:80