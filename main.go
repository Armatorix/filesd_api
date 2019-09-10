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

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	configsPath = "/etc/prometheus/filesd"
	nameScheme  = "filesd_api_%s.json"
)

type FileSDConfig struct {
	Targets []string          `json:targets"`
	Labels  map[string]string `json:labels"`
}

type FileSDConfigs []FileSDConfig

func (f *FileSDConfigs) Hash() (string, error) {
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
func (f *FileSDConfigs) ToPrometheusString() ([]byte, error) {
	return json.Marshal(f)

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
func CreateFilesd(w http.ResponseWriter, r *http.Request) {
	le := log.WithField("endpoint", "CreateFileSDConfig")
	le.Debug("Start")
	var c FileSDConfigs
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		le.WithError(err).Error("body decode")
		return
	}
	log.Info(c)
	confBytes, err := c.ToPrometheusString()
	if err != nil {
		le.WithError(err).Error("marshal")
		return
	}
	if err != nil {
		le.WithError(err).Error("files amount")
		return
	}
	id, err := c.Hash()
	if err != nil {
		le.WithError(err).Error("hash calc")
		return
	}
	err = ioutil.WriteFile(filePath(id), confBytes, 0644)
	if err != nil {
		le.WithError(err).Error("write to file")
		return
	}
	w.WriteHeader(http.StatusOK)
	le.Debug("Done")
}

// DeleteFilesd endpoint deletes file with specific ID in configs path
func DeleteFilesd(w http.ResponseWriter, r *http.Request) {
	le := log.WithField("endpoint", "DeleteFileSDConfig")
	vars := mux.Vars(r)
	fileID := vars["id"]

	path := filePath(fileID)
	le.Debug("removing:" + path)
	err := os.Remove(path)
	if err != nil {
		le.WithError(err).Debug("cannot remove")
		return
	}
	w.WriteHeader(http.StatusOK)
	le.Debug("Done")
}

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	const PORT = "2137"

	log.WithField("port", PORT).Info("API started")

	router := mux.NewRouter()
	router.HandleFunc("/filesd", CreateFilesd).Methods("POST")
	router.HandleFunc("/filesd/{id:[0-9a-fA-F]{40}}", DeleteFilesd).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
