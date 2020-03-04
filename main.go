package main

import (
	"fmt"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"
)

const (
	author  = "webdevops.io"
	version = "0.2.0"

	// Limit of pagerduty incidents per call
	PagerdutyIncidentLimit = 100
)

var (
	argparser    *flags.Parser
	verbose      bool
	daemonLogger *DaemonLogger
)

var opts struct {
	// general settings
	Verbose []bool `long:"verbose" short:"v"  env:"VERBOSE"  description:"verbose mode"`

	// server settings
	ServerBind string        `long:"bind"         env:"SERVER_BIND"   description:"Server address" default:":8080"`
	ScrapeTime time.Duration `long:"scrape-time"  env:"SCRAPE_TIME"   description:"Scrape time (time.duration)" default:"5m"`

	// PagerDuty settings
	PagerDutyAuthToken      string        `long:"pagerduty.authtoken"        env:"PAGERDUTY_AUTH_TOKEN"  description:"PagerDuty auth token" required:"true"`
	PagerDutySince          time.Duration `long:"pagerduty.date-range"       env:"PAGERDUTY_DATE_RANGE"  description:"PagerDuty date range" default:"168h"`
	PagerDutyMaxConnections int           `long:"pagerduty.max-connections"  env:"PAGERDUTY_MAX_CONNECTIONS"                    description:"Maximum numbers of TCP connections to PagerDuty API (concurrency)" default:"4"`

	// ElasticSearch settings
	ElasticsearchAddresses  []string      `long:"elasticsearch.address"      env:"ELASTICSEARCH_ADDRESS"  delim:" "  description:"ElasticSearch urls" required:"true"`
	ElasticsearchIndex      string        `long:"elasticsearch.index"        env:"ELASTICSEARCH_INDEX"               description:"ElasticSearch index name" default:"pagerduty"`
	ElasticsearchRetryCount int           `long:"elasticsearch.retry-count"  env:"ELASTICSEARCH_RETRY_COUNT"         description:"ElasticSearch request retry count" default:"5"`
	ElasticsearchRetryDelay time.Duration `long:"elasticsearch.retry-delay"  env:"ELASTICSEARCH_RETRY_DELAY"         description:"ElasticSearch request delay for reach retry" default:"5s"`
}

func main() {
	initArgparser()

	// set verbosity
	verbose = len(opts.Verbose) >= 1

	// Init logger
	daemonLogger = NewLogger(log.Lshortfile, verbose)
	defer daemonLogger.Close()

	daemonLogger.Infof("Init Pagerduty2ElasticSearch exporter v%s (written by %v)", version, author)

	daemonLogger.Infof("Init exporter")
	exporter := PagerdutyElasticsearchExporter{}
	exporter.Init()
	exporter.SetScrapeTime(opts.ScrapeTime)
	exporter.SetPagerdutyDateRange(opts.PagerDutySince)
	exporter.ConnectPagerduty(
		opts.PagerDutyAuthToken,
		&http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxConnsPerHost:       opts.PagerDutyMaxConnections,
				MaxIdleConns:          opts.PagerDutyMaxConnections,
				IdleConnTimeout:       60 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
			},
		},
	)

	cfg := elasticsearch.Config{
		Addresses: opts.ElasticsearchAddresses,
	}
	exporter.ConnectElasticsearch(cfg, opts.ElasticsearchIndex)
	exporter.SetElasticsearchRetry(opts.ElasticsearchRetryCount, opts.ElasticsearchRetryDelay)
	exporter.Run()

	daemonLogger.Infof("Starting http server on %s", opts.ServerBind)
	startHttpServer()
}

// init argparser and parse/validate arguments
func initArgparser() {
	argparser = flags.NewParser(&opts, flags.Default)
	_, err := argparser.Parse()

	// check if there is an parse error
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println()
			argparser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}
}

// start and handle prometheus handler
func startHttpServer() {
	http.Handle("/metrics", promhttp.Handler())
	daemonLogger.Fatal(http.ListenAndServe(opts.ServerBind, nil))
}
