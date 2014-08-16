package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/golang/glog"

	"github.com/mattrco/anode.exp/analyze"
	"github.com/mattrco/anode.exp/data"
	"github.com/mattrco/anode.exp/input"
	"github.com/mattrco/anode.exp/output"
)

func main() {

	// A metric name to analyze must be provided, all other params have a sensible default.
	// TODO: usage help; check validity of flags; make graphite TCP port configurable.
	var metric = flag.String("metric", "", "graphite metric to retrieve")
	var backfill = flag.String("backfill", "-24hr", "range of backfill data to retrieve")
	var host = flag.String("host", "localhost", "graphite host")
	flag.Parse()
	if *metric == "" {
		glog.Fatal("Metric name required, none given")
	}

	// Start a graphite input instance with one channel for receiving new values.
	rec := make(chan data.Datapoint)
	receivers := []chan data.Datapoint{rec}

	gf := input.GraphiteFetch{}
	err := gf.Init(*host, *metric, *backfill, receivers)
	if err != nil {
		glog.Fatal(err)
	}
	go gf.Run()

	// Start an analyzer, passing in the channel that receives new graphite values,
	// and an output channel to propagate values back to graphite.
	outchan := make(chan data.Datapoint)

	change := analyze.Change{}
	// TODO: magic number alert.
	err = change.Init(rec, outchan, 120, 30, 10)
	if err != nil {
		glog.Fatal(err)
	}
	go change.Run()

	// Output values sent to outchan to graphite.
	output := output.Graphite{}
	err = output.Init(fmt.Sprintf("%s:2003", *host), outchan)
	if err != nil {
		glog.Fatal(err)
	}
	go output.Run()

	// Exit on SIGINT/SIGKILL.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	s := <-signals
	glog.Infof("Caught %s, exiting...\n", s)
}
