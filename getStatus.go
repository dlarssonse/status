package status

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	gops "github.com/mitchellh/go-ps"
	"github.com/ricochet2200/go-disk-usage/du"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// GetStatus ...
func GetStatus(filename string) (*Status, error) {
	stat := Status{}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &stat)
	if err != nil {
		return nil, err
	}

	for i, service := range stat.Services {
		switch service.Type {

		// Process Check
		case "Process":
			if err = GetProcessStatus(&service); err != nil {
				stat.Services[i].Status = "Error"
				stat.Services[i].Error = err.Error()
			} else {
				stat.Services[i] = service
			}

		// Disk space
		case "Drive Space":
			if err = GetDriveSpaceStatus(&service); err != nil {
				stat.Services[i].Status = "Error"
				stat.Services[i].Error = err.Error()
			} else {
				stat.Services[i] = service
			}
		}

	}

	return &stat, nil
}

// GetDriveSpaceStatus ...
func GetDriveSpaceStatus(service *Service) error {

	var KB = uint64(1024)

	usage := du.NewDiskUsage(service.Filename)
	/*
		fmt.Println("Free:", usage.Free()/(KB*KB))
		fmt.Println("Available:", usage.Available()/(KB*KB))
		fmt.Println("Size:", usage.Size()/(KB*KB))
		fmt.Println("Used:", usage.Used()/(KB*KB))
		fmt.Println("Usage:", usage.Usage()*100, "%")
	*/

	if usage != nil {
		p := message.NewPrinter(language.English)
		service.Status = p.Sprintf("Free: %d MB (%0.1f%%)", usage.Free()/(KB*KB), 100-usage.Usage()*100)
		return nil
	}

	service.Status = "Not Working"
	service.Error = ""
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
			service.Error = ""
			return nil
		}
	}

	service.Status = "Not Running"
	service.Error = ""
	return nil
}
