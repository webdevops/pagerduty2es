PagerDuty2Elasticsearch exporter
================================

[![license](https://img.shields.io/github/license/webdevops/pagerduty2elasticsearch-exporter.svg)](https://github.com/webdevops/pagerduty2elasticsearch-exporter/blob/master/LICENSE)
[![Docker](https://img.shields.io/badge/docker-webdevops%2Fpagerduty--exporter-blue.svg?longCache=true&style=flat&logo=docker)](https://hub.docker.com/r/webdevops/pagerduty2elasticsearch-exporter/)
[![Docker Build Status](https://img.shields.io/docker/build/webdevops/pagerduty2elasticsearch-exporter.svg)](https://hub.docker.com/r/webdevops/pagerduty2elasticsearch-exporter/)

Exporter for incidents and logentries from PagerDuty to ElasticSearch

Configuration
-------------

| Environment variable                    | DefaultValue                | Description                                                              |
|-----------------------------------------|-----------------------------|--------------------------------------------------------------------------|
| `SCRAPE_TIME`                           | `5m`                        | Time (time.Duration) for general informations                            |
| `SERVER_BIND`                           | `:8080`                     | IP/Port binding                                                          |
| `PAGERDUTY_AUTH_TOKEN`                  | none                        | PagerDuty auth token                                                     |
| `PAGERDUTY_DATE_RANGE`                  | `168h`                      | Date range for importing historical data                                 |
| `PAGERDUTY_MAX_CONNECTIONS`             | `4`                         | Maximum numbers of HTTP connections to PagerDuty API                     |
| `ELASTICSEARCH_ADDRESS`                 | none, required              | ElasticSearch cluster addresses (multiple)                               |
| `ELASTICSEARCH_INDEX`                   | `pagerduty`                 | Name of ElasticSearch index                                              |
| `ELASTICSEARCH_BATCH_COUNT`             | `50`                        | Number of documents which should be indexed in one request                |
| `ELASTICSEARCH_RETRY_COUNT`             | `5`                         | ElasticSearch request retry count                                        |
| `ELASTICSEARCH_RETRY_DELAY`             | `5s`                        | ElasticSearch request delay for reach retry                              |

Metrics
-------

| Metric                                       | Description                                                        |
|----------------------------------------------|--------------------------------------------------------------------|
| `pagerduty2es_incident_total`                | Total number of processed incidents                                |
| `pagerduty2es_incident_logentry_total`       | Total number of processed logentries                               |
| `pagerduty2es_duration`                      | Scrape process duration                                            |
| `pagerduty2es_elasticsearch_requet_total`    | Number of total requests to ElasticSearch cluster                  |
| `pagerduty2es_elasticsearch_request_retries` | Number of retried requests to ElasticSearch cluster                |
