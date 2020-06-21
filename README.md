PagerDuty2Elasticsearch exporter
================================

[![license](https://img.shields.io/github/license/webdevops/pagerduty2elasticsearch-exporter.svg)](https://github.com/webdevops/pagerduty2elasticsearch-exporter/blob/master/LICENSE)
[![Docker](https://img.shields.io/badge/docker-webdevops%2Fpagerduty--exporter-blue.svg?longCache=true&style=flat&logo=docker)](https://hub.docker.com/r/webdevops/pagerduty2elasticsearch-exporter/)
[![Docker Build Status](https://img.shields.io/docker/build/webdevops/pagerduty2elasticsearch-exporter.svg)](https://hub.docker.com/r/webdevops/pagerduty2elasticsearch-exporter/)

Exporter for incidents and logentries from PagerDuty to ElasticSearch

Configuration
-------------

```
Usage:
  pagerduty2elasticsearch-exporter [OPTIONS]

Application Options:
  -v, --verbose                    verbose mode [$VERBOSE]
      --bind=                      Server address (default: :8080) [$SERVER_BIND]
      --scrape-time=               Scrape time (time.duration) (default: 5m) [$SCRAPE_TIME]
      --pagerduty.authtoken=       PagerDuty auth token [$PAGERDUTY_AUTH_TOKEN]
      --pagerduty.date-range=      PagerDuty date range (default: 168h) [$PAGERDUTY_DATE_RANGE]
      --pagerduty.max-connections= Maximum numbers of TCP connections to PagerDuty API (concurrency) (default: 4) [$PAGERDUTY_MAX_CONNECTIONS]
      --elasticsearch.address=     ElasticSearch urls [$ELASTICSEARCH_ADDRESS]
      --elasticsearch.index=       ElasticSearch index name (placeholders: %y for year, %m for month and %d for day) (default: pagerduty) [$ELASTICSEARCH_INDEX]
      --elasticsearch.batch-count= Number of documents which should be indexed in one request (default: 50) [$ELASTICSEARCH_BATCH_COUNT]
      --elasticsearch.retry-count= ElasticSearch request retry count (default: 5) [$ELASTICSEARCH_RETRY_COUNT]
      --elasticsearch.retry-delay= ElasticSearch request delay for reach retry (default: 5s) [$ELASTICSEARCH_RETRY_DELAY]

Help Options:
  -h, --help                       Show this help message
```

Metrics
-------

| Metric                                       | Description                                                        |
|----------------------------------------------|--------------------------------------------------------------------|
| `pagerduty2es_incident_total`                | Total number of processed incidents                                |
| `pagerduty2es_incident_logentry_total`       | Total number of processed logentries                               |
| `pagerduty2es_duration`                      | Scrape process duration                                            |
| `pagerduty2es_elasticsearch_requet_total`    | Number of total requests to ElasticSearch cluster                  |
| `pagerduty2es_elasticsearch_request_retries` | Number of retried requests to ElasticSearch cluster                |
