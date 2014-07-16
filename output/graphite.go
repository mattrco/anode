package output

import (
	"fmt"
	"net"

	"github.com/golang/glog"

	"github.com/mattrco/anode.exp/data"
)

type Graphite struct {
	name    string
	conn    net.Conn
	outchan chan data.Datapoint
}

func (out *Graphite) Init(graphiteHost string, name string, outchan chan data.Datapoint) error {
	conn, err := net.Dial("tcp", graphiteHost)
	if err != nil {
		return err
	}
	out.conn = conn
	out.outchan = outchan
	out.name = name
	return nil
}

func (out *Graphite) Run() {
	for o := range out.outchan {
		// Graphite TCP format: <name> <value> <timestamp>\n
		outs := fmt.Sprintf("%s %f %d\n", out.name, o.Value, o.Timestamp)
		_, err := fmt.Fprintf(out.conn, outs)
		if err != nil {
			glog.Error(err)
		}
	}
}
