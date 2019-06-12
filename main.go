package main

// TODO add web UI
// valudate input with some logging
// integration with loki

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	configsPath = "/etc/prometheus/filesd"
	nameScheme  = "filesd_api_%x.json"
)

type FileSDConfig struct {
	Targets []string `json:targets"`
	Labels  map[string]string
}

func (f *FileSDConfig) Hash() (string, error) {
	content, err := json.Marshal(f)
	if err != nil {
		return "", err
	}
	h := sha1.New()
	h.Write([]byte(content))
	xd := h.Sum(nil)
	return string(xd), nil
}

// ToPrometheusString returns config usable by prometheus
func (f *FileSDConfig) ToPrometheusString() ([]byte, error) {
	return json.Marshal([]FileSDConfig{*f})

}

func configsLen() (int, error) {
	files, err := ioutil.ReadDir(configsPath)
	return len(files), err
}

func filePath(id string) string {
	return fmt.Sprintf(fmt.Sprintf("%s/%s", configsPath, nameScheme), id)
}

// CreateFilesd POST endpoint, which creates file in configs path
// naming is based on specified scheme and on hash of input data

func CreateFilesd(w http.ResponseWriter, req *http.Request) {
	const endpointName = "CreateFileSDConfig"
	log.WithField("endpoint", endpointName).Debug("Start")
	var c FileSDConfig
	err := json.NewDecoder(req.Body).Decode(&c)
	if err != nil {
		log.WithField("endpoint", endpointName).WithError(err).Error("body decode")
		return
	}
	confBytes, err := c.ToPrometheusString()
	if err != nil {
		log.WithField("endpoint", endpointName).WithError(err).Error("marshal")
		return
	}
	if err != nil {
		log.WithField("endpoint", endpointName).WithError(err).Error("files amount")
		return
	}
	id, err := c.Hash()
	if err != nil {
		log.WithField("endpoint", endpointName).WithError(err).Error("hash calc")
		return
	}
	err = ioutil.WriteFile(filePath(id), confBytes, 0755)
	if err != nil {
		log.WithField("endpoint", endpointName).WithError(err).Error("write to file")
		return
	}
	log.WithField("endpoint", endpointName).Debug("Done")
}

func DeleteFilesd(w http.ResponseWriter, req *http.Request) {

}

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	const PORT = "2137"

	log.WithField("port", PORT).Info("API started")

	router := mux.NewRouter()
	router.HandleFunc("/filesd", CreateFilesd).Methods("POST")
	router.HandleFunc("/filesd/{id}", DeleteFilesd).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
