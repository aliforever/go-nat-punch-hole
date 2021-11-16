package args

import "flag"

var (
	Type                   string
	LocalAddress           string
	SignalingServerAddress string
)

func init() {
	flag.StringVar(&Type, "type", "", "--type=server/client")
	flag.StringVar(&LocalAddress, "local", "127.0.0.1:8181", "--local=:8181")
	flag.StringVar(&SignalingServerAddress, "server", "", "--server=142.256.14.21")
	flag.Parse()
}
