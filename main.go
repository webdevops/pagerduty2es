package main

import (
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
)

const (
	Author               = "webdevops.io"
	Version              = "0.1.0"
)

var (
	argparser            *flags.Parser
	args                 []string
	Verbose              bool
	Logger               *DaemonLogger
	PagerDutyClient      *pagerduty.Client
)

var opts struct {
	// general settings
	Verbose []bool `long:"verbose" short:"v"        env:"VERBOSE"                description:"Verbose mode"`

	// server settings
	ServerBind     string        `long:"bind"               env:"SERVER_BIND"            description:"Server address"                                     default:":8080"`

	// PagerDuty settings
	PagerDutyAuthToken                 string        `long:"pagerduty.authtoken"                                         env:"PAGERDUTY_AUTH_TOKEN"                         description:"PagerDuty auth token" required:"true"`
}

func main() {
	initArgparser()

	// set verbosity
	Verbose = len(opts.Verbose) >= 1

	// Init logger
	Logger = NewLogger(log.Lshortfile, Verbose)
	defer Logger.Close()

	Logger.Infof("Init Pagerduty to ElasticSearch exporter v%s (written by %v)", Version, Author)

	Logger.Infof("Init PagerDuty client")
	initPagerDuty()

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
			fmt.Println(err)
			fmt.Println()
			argparser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}
}

// Init and build PagerDuty client
func initPagerDuty() {
	PagerDutyClient = pagerduty.NewClient(opts.PagerDutyAuthToken)
}

// start and handle prometheus handler
func startHttpServer() {
	http.Handle("/metrics", promhttp.Handler())
	Logger.Fatal(http.ListenAndServe(opts.ServerBind, nil))
}
