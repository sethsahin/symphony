package cli

import (
	"github.com/erkrnt/symphony/api"
	"github.com/erkrnt/symphony/internal/service"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// ServiceNewOptions : options for initializing a service
type ServiceNewOptions struct {
	APIServerAddr *string
	ClusterID     *string
	ServiceType   *string
	SocketPath    *string
}

// ServiceNew : initializes a service for use
func ServiceNew(opts ServiceNewOptions) {
	if opts.APIServerAddr == nil {
		logrus.Fatal("Missing --apiserver-addr option. Check --help for more.")
	}

	if opts.SocketPath == nil {
		logrus.Fatal("Missing --socket-path option. Check --help for more.")
	}

	conn := service.NewClientConnSocket(*opts.SocketPath)

	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), service.ContextTimeout)

	defer cancel()

	c := api.NewControlClient(conn)

	if *opts.ServiceType != "apiserver" && *opts.ServiceType != "block" {
		logrus.Fatal("invalid_service_type")
	}

	var st api.ServiceType

	switch *opts.ServiceType {
	case "apiserver":
		st = api.ServiceType_APISERVER
	case "block":
		st = api.ServiceType_BLOCK
	}

	options := &api.RequestServiceNew{
		APIServerAddr: *opts.APIServerAddr,
		ClusterID:     *opts.ClusterID,
		ServiceType:   st,
	}

	cluster, err := c.ServiceNew(ctx, options)

	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info(cluster)
}
