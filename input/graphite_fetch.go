// Package input contains input plugins for collecting metrics and passing them to
// analysis functions.
//
// TODO: replace sleep with the design in http://talks.golang.org/2013/advconc.slide#1,
// which would also provide a quit channel.
package input

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/golang/glog"

	"github.com/mattrco/anode/data"
)

type GraphiteFetch struct {
	host   string
	metric string
	// Initial data range to fetch, e.g. "-24hr" will fetch last 24 hours.
	initRange string
	receivers []chan data.Datapoint
}

/* Structure of graphite-webapp JSON response:
[
  {
    "target": "host.metric.a",
    "datapoints": [[<metric> <timestamp>]]
  },
  {...}
]
*/
type GraphiteResponse []GraphiteMetric

type GraphiteMetric struct {
	Target     string           `json:"target"`
	Datapoints [][]*json.Number `json:"datapoints"`
}

func (gf *GraphiteFetch) Init(host string, metric string, initRange string, receivers []chan data.Datapoint) error {
	gf.host = host
	gf.metric = metric
	gf.initRange = initRange
	gf.receivers = receivers
	if glog.V(2) {
		glog.Infof("Init graphite: %s: %s\n", host, metric)
	}
	return nil
}

func (gf *GraphiteFetch) Run() error {
	// If the metric can't be fetched now, assume we cannot proceed.
	err := gf.fetch(gf.initRange)
	if err != nil {
		glog.Fatal(err)
	}

	// Calculate the interval between fetches based on the last two metric timestamps.
	// If an interval cannot be calculated, default to 60s.
	intvl := 60
	for {
		// Sleep for the interval, then fetch new metrics.
		duration := time.Duration(intvl) * time.Second
		if glog.V(2) {
			glog.Infof("Sleeping for %v\n", duration)
		}
		time.Sleep(duration)
		// Time is supplied relative, e.g. from=-60sec to fetch metrics
		// received in the last 60 seconds.
		err := gf.fetch(fmt.Sprintf("-%dsec", intvl))
		if err != nil {
			glog.Errorf("Error fetching graphite metric: %s\n", err)
		}
	}
}

func (gf *GraphiteFetch) fetch(from string) error {
	mURL := url.URL{
		Scheme: "http",
		Host:   gf.host,
		Path:   "render",
	}
	query := url.Values{}
	query.Set("target", gf.metric)
	query.Set("from", from)
	query.Set("format", "json")
	mURL.RawQuery = query.Encode()

	resp, err := http.Get(mURL.String())
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	seriesList := GraphiteResponse{}
	err = json.Unmarshal(body, &seriesList)
	if err != nil {
		return err
	}

	for _, v := range seriesList[0].Datapoints {
		// Graphite returns a datapoint for each timestamp even if there is
		// no metric value, so skip any for which the metric value is nil.
		if v[0] != nil {
			if glog.V(3) {
				glog.Infof("Metric %s: %v at %v\n", gf.metric, v[0], v[1])
			}
			ts, err := v[1].Int64()
			if err != nil {
				return err
			}
			fltVal, err := v[0].Float64()
			if err != nil {
				return err
			}

			d := data.Datapoint{
				Timestamp: ts,
				Value:     fltVal,
			}
			// Send new datapoint to all receivers.
			for _, r := range gf.receivers {
				r <- d
			}
		}
	}
	return nil
}
