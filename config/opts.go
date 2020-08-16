package config

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"time"
)

type (
	Opts struct {
		// logger
		Logger struct {
			Debug   bool `           long:"debug"        env:"DEBUG"    description:"debug mode"`
			Verbose bool `short:"v"  long:"verbose"      env:"VERBOSE"  description:"verbose mode"`
			LogJson bool `           long:"log.json"     env:"LOG_JSON" description:"Switch log output to json format"`
		}

		// PagerDuty settings
		PagerDuty struct {
			AuthToken      string        `long:"pagerduty.authtoken"        env:"PAGERDUTY_AUTH_TOKEN"        description:"PagerDuty auth token" required:"true" json:"-"`
			Since          time.Duration `long:"pagerduty.date-range"       env:"PAGERDUTY_DATE_RANGE"        description:"PagerDuty date range" default:"168h"`
			MaxConnections int           `long:"pagerduty.max-connections"  env:"PAGERDUTY_MAX_CONNECTIONS"   description:"Maximum numbers of TCP connections to PagerDuty API (concurrency)" default:"4"`
		}

		// ElasticSearch settings
		Elasticsearch struct {
			Addresses  []string      `long:"elasticsearch.address"      env:"ELASTICSEARCH_ADDRESS"  delim:" "  description:"ElasticSearch urls" required:"true"`
			Username   string        `long:"elasticsearch.username"     env:"ELASTICSEARCH_USERNAME"            description:"ElasticSearch username for HTTP Basic Authentication"`
			Password   string        `long:"elasticsearch.password"     env:"ELASTICSEARCH_PASSWORD"            description:"ElasticSearch password for HTTP Basic Authentication" json:"-"`
			ApiKey     string        `long:"elasticsearch.apikey"       env:"ELASTICSEARCH_APIKEY"              description:"ElasticSearch base64-encoded token for authorization; if set, overrides username and password" json:"-"`
			Index      string        `long:"elasticsearch.index"        env:"ELASTICSEARCH_INDEX"               description:"ElasticSearch index name (placeholders: %y for year, %m for month and %d for day)" default:"pagerduty"`
			BatchCount int           `long:"elasticsearch.batch-count"  env:"ELASTICSEARCH_BATCH_COUNT"         description:"Number of documents which should be indexed in one request"  default:"50"`
			RetryCount int           `long:"elasticsearch.retry-count"  env:"ELASTICSEARCH_RETRY_COUNT"         description:"ElasticSearch request retry count"                           default:"5"`
			RetryDelay time.Duration `long:"elasticsearch.retry-delay"  env:"ELASTICSEARCH_RETRY_DELAY"         description:"ElasticSearch request delay for reach retry"                 default:"5s"`
		}

		// general options
		ServerBind string        `long:"bind"     env:"SERVER_BIND"   description:"Server address"     default:":8080"`
		ScrapeTime time.Duration `long:"scrape-time"  env:"SCRAPE_TIME"   description:"Scrape time (time.duration)" default:"5m"`
	}
)

func (o *Opts) GetJson() []byte {
	jsonBytes, err := json.Marshal(o)
	if err != nil {
		log.Panic(err)
	}
	return jsonBytes
}
