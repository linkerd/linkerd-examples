FROM grafana/grafana:5.2.2

COPY grafana.ini              $GF_PATHS_CONFIG
COPY datasources.yaml         $GF_PATHS_PROVISIONING/datasources/
COPY dashboards.yaml          $GF_PATHS_PROVISIONING/dashboards/
COPY dashboards/*             $GF_PATHS_PROVISIONING/dashboards/
COPY dashboards/l5d-perf.json $GF_PATHS_HOME/public/dashboards/home.json
