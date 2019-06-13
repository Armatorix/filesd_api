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
	"os"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	configsPath = "/etc/prometheus/filesd"
	nameScheme  = "filesd_api_%s.json"
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
	hash := h.Sum(nil)
	return fmt.Sprintf("%x", hash), nil
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

func isValidHexadecimal(id string) bool {
	for _, l := range strings.ToLower(id) {
		if (l >= '0' && l <= '9') ||
			(l >= 'a' && l <= 'f') {
			continue
		}
		return false
	}
	return true
}

// CreateFilesd POST endpoint, which creates file in configs path
// naming is based on specified scheme and on hash of input data
func CreateFilesd(w http.ResponseWriter, r *http.Request) {
	const endpointName = "CreateFileSDConfig"
	log.WithField("endpoint", endpointName).Debug("Start")
	var c FileSDConfig
	err := json.NewDecoder(r.Body).Decode(&c)
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
	err = ioutil.WriteFile(filePath(id), confBytes, 0644)
	if err != nil {
		log.WithField("endpoint", endpointName).WithError(err).Error("write to file")
		return
	}
	w.WriteHeader(http.StatusOK)
	log.WithField("endpoint", endpointName).Debug("Done")
}

func DeleteFilesd(w http.ResponseWriter, r *http.Request) {
	const endpointName = "DeleteFileSDConfig"
	vars := mux.Vars(r)
	fileID := vars["id"]
	log.Info(fileID)
	if !isValidHexadecimal(fileID) {
		log.WithField("endpoint", endpointName).Error("invalid hexadecimal")
		return
	}
	path := filePath(fileID)
	log.WithField("endpoint", endpointName).Debug("removing:" + path)
	err := os.Remove(path)
	if err != nil {
		log.WithField("endpoint", endpointName).WithError(err).Debug("cannot remove")
		return
	}
	w.WriteHeader(http.StatusOK)
	log.WithField("endpoint", endpointName).Debug("Done")
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
