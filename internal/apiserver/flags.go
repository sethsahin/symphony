package apiserver

import (
	"errors"
	"net"

	"github.com/erkrnt/symphony/internal/service"
	"github.com/hashicorp/go-sockaddr"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Flags : command line flags
type Flags struct {
	ConfigDir  string
	ConsulAddr string
	ListenAddr *net.TCPAddr
	Verbose    bool
}

var (
	bindInterface = kingpin.Flag("bind-interface", "Sets the bind interface for listening services.").Required().String()
	configDir     = kingpin.Flag("config-dir", "Sets configuration directory for manager.").Default(".").String()
	consulAddr    = kingpin.Flag("consul-addr", "Sets the etcd endpoints list.").Required().String()
	verbose       = kingpin.Flag("verbose", "Sets the lowest level of service output.").Bool()
)

func getFlags() (*Flags, error) {
	kingpin.Parse()

	if bindInterface == nil {
		return nil, errors.New("invalid_bind_interface")
	}

	configDirPath, err := service.GetDirPath(configDir)

	if err != nil {
		return nil, err
	}

	ipAddr, err := sockaddr.GetInterfaceIP(*bindInterface)

	if err != nil {
		return nil, err
	}

	listenAddr, err := service.GetListenAddr(ipAddr, nil)

	if err != nil {
		return nil, err
	}

	flags := &Flags{
		ConfigDir:  *configDirPath,
		ConsulAddr: *consulAddr,
		ListenAddr: listenAddr,
		Verbose:    *verbose,
	}

	fields := logrus.Fields{
		"ConfigDir":  flags.ConfigDir,
		"ConsulAddr": flags.ConsulAddr,
		"ListenAddr": flags.ListenAddr.String(),
		"Verbose":    flags.Verbose,
	}

	logrus.WithFields(fields).Info("Service command-line flags loaded.")

	return flags, nil
}