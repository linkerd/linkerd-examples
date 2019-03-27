FROM grafana/grafana:5.4.3

USER root

RUN apt-get update && \
    apt-get -y --no-install-recommends install curl

RUN mkdir -p /var/lib/grafana/dashboards
COPY ./grafana.json /usr/share/grafana/public/dashboards/home.json

COPY ./bootstrap.sh /
ENTRYPOINT ["/bootstrap.sh"]
