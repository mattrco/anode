package output

import (
	"fmt"
	"net"

	"github.com/golang/glog"

	"github.com/mattrco/anode.exp/data"
)

type Graphite struct {
	conn  net.Conn
	input chan data.Datapoint
}

func (g *Graphite) Init(graphiteHost string, input chan data.Datapoint) error {
	conn, err := net.Dial("tcp", graphiteHost)
	if err != nil {
		return err
	}
	g.conn = conn
	g.input = input
	return nil
}

func (g *Graphite) Run() {
	for data := range g.input {
		// Graphite TCP format: <name> <value> <timestamp>\n
		outs := fmt.Sprintf("%s %f %d\n", data.MetricName, data.Value, data.Timestamp)
		_, err := fmt.Fprintf(g.conn, outs)
		if err != nil {
			glog.Error(err)
		}
	}
}
