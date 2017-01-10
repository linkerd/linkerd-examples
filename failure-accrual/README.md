# Failure accrual

This directory contains a docker-compose environment that runs a demo that you
can use to test [linkerd's failure accrual settings](
https://linkerd.io/features/circuit-breaking/#failure-accrual).

## Startup

To start linkerd preconfigured with 3 routers, routing to 15 go backends, using
different failure accrual settings for each router, generating traffic via
[slow_cooker](https://github.com/BuoyantIO/slow_cooker), and exporting stats to
prometheus and grafana, run:

```bash
$ docker-compose build && docker-compose up -d
```

## Dashboards

### Grafana

Grafana is running on port 3000 in your docker-compose environment. To see a
dashboard comparing the different failure accrual settings, load the Grafana
dashboard by going to port 3000 on your docker host. It should look like this:

![grafana](screenshot-grafana.png)

### linkerd admin

The linkerd admin server is running on port 9990 in your docker-compose
environment. To see the admin dashboard, go to port 9990 on your docker host. It
should look like this:

![linkerd](screenshot-linkerd.png)

## Troubleshooting

If you have any issues getting the demo up and running, pop into [linkerd's
Slack]( https://slack.linkerd.io) and we'll help you get it sorted out.

Thanks! 👋
