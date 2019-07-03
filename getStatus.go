package status

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	gops "github.com/mitchellh/go-ps"
	"github.com/ricochet2200/go-disk-usage/du"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// accessControlHeaders
func accessControlHeaders(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers")
}

// APIErrorResponse ...
func errorResponse(w http.ResponseWriter, msg string) {
	accessControlHeaders(w)
	w.WriteHeader(http.StatusBadRequest)
	encoder := json.NewEncoder(w)
	err := map[string]string{"status": "ERROR", "message": msg}
	encoder.Encode(err)
}

// AddMuxRoute ...
func AddMuxRoute(path string, config *Status, router *mux.Router) {
	router.Handle(path, Handler(config)).Name("STATUS")
}

// Handler handles all requests
func Handler(config *Status) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			accessControlHeaders(w)
			return
		}

		if err := GetStatus(config); err != nil {
			errorResponse(w, err.Error())
			return
		}

		data, err := json.Marshal(config)
		if err != nil {
			errorResponse(w, err.Error())
			return
		}

		accessControlHeaders(w)
		w.Header().Add("Content-Type", "application/json")
		w.Write(data)
		return
	})
}

// ReadFile ...
func ReadFile(filename string) (*Status, error) {
	status := Status{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

// GetStatus ...
func GetStatus(status *Status) error {
	var err error

	for i, service := range status.Services {
		switch service.Type {

		// Process Check
		case "Process":
			if err = GetProcessStatus(&service); err != nil {
				status.Services[i].Status = "Error"
				status.Services[i].Error = err.Error()
			} else {
				status.Services[i] = service
			}

		// Disk space
		case "Drive Space":
			if err = GetDriveSpaceStatus(&service); err != nil {
				status.Services[i].Status = "Error"
				status.Services[i].Error = err.Error()
			} else {
				status.Services[i] = service
			}
		}

	}

	return nil
}

// GetDriveSpaceStatus ...
func GetDriveSpaceStatus(service *Service) error {

	var KB = uint64(1024)

	usage := du.NewDiskUsage(service.Filename)
	if usage != nil {
		if usage.Size() == 0 {
			service.Status = "Path Not Found"
			return fmt.Errorf("Path '%s' Not Found", service.Filename)
		}
		p := message.NewPrinter(language.English)
		service.Status = p.Sprintf("Free: %d GB (%0.1f%%)", usage.Free()/(KB*KB), 100-usage.Usage()*100)
		return nil
	}

	service.Status = "Not Working"
	return nil
}

// GetProcessStatus ...
func GetProcessStatus(service *Service) error {
	service.Status = "OK"
	processes, err := gops.Processes()
	if err != nil {
		return err
	}

	if len(processes) <= 0 {
		return fmt.Errorf("no processes found")
	}

	for _, p := range processes {
		if p.Executable() == service.Filename {
			service.Status = fmt.Sprintf("Running (PID: %d)", p.Pid())
			return nil
		}
	}

	service.Status = "Not Running"
	return nil
}
