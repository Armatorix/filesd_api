package main

// TODO add web UI
// valudate input with some logging
// change naming scheme
// fix self registration
// integration with loki

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	configsPath = "/etc/prometheus/filesd"
)

type FileSDConfig struct {
	Targets []string `json:targets"`
	Labels  map[string]string
}

func configsLen() (int, error) {
	files, err := ioutil.ReadDir(configsPath)
	return len(files), err
}

func Config(w http.ResponseWriter, req *http.Request) {
	const endpointName = "CreateFileSDConfig"
	log.WithField("endpoint", endpointName).Debug("Start")
	var c FileSDConfig
	err := json.NewDecoder(req.Body).Decode(&c)
	if err != nil {
		log.WithField("endpoint", endpointName).WithError(err).Error("body decode")
		return
	}
	confBytes, err := json.Marshal([]FileSDConfig{c})
	if err != nil {
		log.WithField("endpoint", endpointName).WithError(err).Error("marshal")
		return
	}
	nextConf, err := configsLen()
	if err != nil {
		log.WithField("endpoint", endpointName).WithError(err).Error("files amount")
		return
	}
	err = ioutil.WriteFile(fmt.Sprintf(configsPath+"/filesd_api_%d.json", nextConf), confBytes, 0755)
	if err != nil {
		log.WithField("endpoint", endpointName).WithError(err).Error("write to file")
		return
	}
	log.WithField("endpoint", endpointName).Debug("Done")
}

func Register(w http.ResponseWriter, req *http.Request) {
	// log.Error(req.RemoteAddr)
	const endpointName = "CreateFileSDConfig"
	log.WithField("endpoint", endpointName).Debug("Start")
	var c FileSDConfig
	err := json.NewDecoder(req.Body).Decode(&c)
	if err != nil {
		log.WithField("endpoint", endpointName).WithError(err).Error("body decode")
		return
	}
	confBytes, err := json.Marshal([]FileSDConfig{c})
	if err != nil {
		log.WithField("endpoint", endpointName).WithError(err).Error("marshal")
		return
	}
	nextConf, err := configsLen()
	if err != nil {
		log.WithField("endpoint", endpointName).WithError(err).Error("files amount")
		return
	}
	err = ioutil.WriteFile(fmt.Sprintf(configsPath+"/filesd_api_%d.json", nextConf), confBytes, 0755)
	if err != nil {
		log.WithField("endpoint", endpointName).WithError(err).Error("write to file")
		return
	}
	log.WithField("endpoint", endpointName).Debug("Done")
}

func init() {

	log.SetLevel(log.DebugLevel)

}
func main() {
	const PORT = "2137"

	log.WithField("port", PORT).Info("API started")

	router := mux.NewRouter()
	router.HandleFunc("/register/config", Config).Methods("POST")
	router.HandleFunc("/register/self", Register).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
