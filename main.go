package main

import (
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	Author  = "webdevops.io"
	Version = "0.1.0"

	PAGERDUTY_INCIDENT_LIMIT = 100
)

var (
	argparser       *flags.Parser
	args            []string
	Verbose         bool
	Logger          *DaemonLogger
	PagerDutyClient *pagerduty.Client
)

var opts struct {
	// general settings
	Verbose []bool `long:"verbose" short:"v"  env:"VERBOSE"  description:"Verbose mode"`

	// server settings
	ServerBind string        `long:"bind"         env:"SERVER_BIND"   description:"Server address" default:":8080"`
	ScrapeTime time.Duration `long:"scrape-time"  env:"SCRAPE_TIME"   description:"Scrape time (time.duration)" default:"5m"`

	// PagerDuty settings
	PagerDutyAuthToken string        `long:"pagerduty.authtoken"   env:"PAGERDUTY_AUTH_TOKEN"  description:"PagerDuty auth token" required:"true"`
	PagerDutySince     time.Duration `long:"pagerduty.date-range"  env:"PAGERDUTY_DATE_RANGE"  description:"PagerDuty date range" default:"168h"`

	// ElasticSearch settings
	ElasticsearchAddresses []string `long:"elasticsearch.address"  env:"ELASTICSEARCH_ADDRESS"  delim:" "  description:"ElasticSearch urls" required:"true"`
	ElasticsearchIndex     string   `long:"elasticsearch.index"    env:"ELASTICSEARCH_INDEX"               description:"ElasticSearch index name" default:"pagerduty"`
}

func main() {
	initArgparser()

	// set verbosity
	Verbose = len(opts.Verbose) >= 1

	// Init logger
	Logger = NewLogger(log.Lshortfile, Verbose)
	defer Logger.Close()

	Logger.Infof("Init Pagerduty2ElasticSearch exporter v%s (written by %v)", Version, Author)

	Logger.Infof("Init exporter")
	exporter := PagerdutyElasticsearchExporter{}
	exporter.Init()
	exporter.SetScrapeTime(opts.ScrapeTime)
	exporter.SetPagerdutyDateRange(opts.PagerDutySince)
	exporter.ConnectPagerduty(opts.PagerDutyAuthToken)

	cfg := elasticsearch.Config{
		Addresses: opts.ElasticsearchAddresses,
	}
	exporter.ConnectElasticsearch(cfg, opts.ElasticsearchIndex)
	exporter.Run()

	Logger.Infof("Starting http server on %s", opts.ServerBind)
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
	Logger.Fatal(http.ListenAndServe(opts.ServerBind, nil))
}
