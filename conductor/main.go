package main

import (
	"net/http"

	"github.com/erkrnt/symphony/services"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	port = ":36837"
)

func main() {
	flags := GetFlags()
	if flags.Verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	db, err := GetDatabase(flags)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := mux.NewRouter()

	r.Path("/service").Queries("hostname", "{hostname}").Handler(services.RegisterHandler(GetServiceByHostnameHandler(db))).Methods("GET")
	r.Path("/service").Handler(services.RegisterHandler(PostServiceHandler(db))).Methods("POST")
	r.Path("/service/{id}").Handler(services.RegisterHandler(GetServiceByIDHandler(db))).Methods("GET")

	r.Path("/servicetype").Queries("name", "{name}").Handler(services.RegisterHandler(GetServiceTypeByNameHandler(db))).Methods("GET")

	r.Path("/physicalvolume").Queries("device", "{device}").Queries("service_id", "{service_id}").Handler(services.RegisterHandler(GetPhysicalVolumeByDeviceHandler(db))).Methods("GET")
	r.Path("/physicalvolume").Handler(services.RegisterHandler(PostPhysicalVolumeHandler(db))).Methods("POST")
	r.Path("/physicalvolume/{id}").Handler(services.RegisterHandler(DeletePhysicalVolumeHandler(db))).Methods("DELETE")

	r.Path("/volumegroup").Queries("physical_volume_id", "{physical_volume_id}").Queries("service_id", "{service_id}").Handler(services.RegisterHandler(GetVolumeGroupByPhysicalVolumeIDAndServiceIDHandler(db))).Methods("GET")
	r.Path("/volumegroup").Handler(services.RegisterHandler(PostVolumeGroupHandler(db))).Methods("POST")
	r.Path("/volumegroup/{id}").Handler(services.RegisterHandler(DeleteVolumeGroupHandler(db))).Methods("DELETE")

	logrus.WithFields(logrus.Fields{"port": port}).Info("Started conductor service.")

	logrus.Fatal(http.ListenAndServe(port, r))
}
