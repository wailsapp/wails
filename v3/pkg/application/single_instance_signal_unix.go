//go:build darwin || linux

package application

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
)

func (m *singleInstanceManager) setupSignalHandler() {
	if m == nil || m.options == nil || m.options.OnSecondInstanceLaunch == nil {
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGUSR1)

	go func() {
		for range sigChan {
			// Read data from the temporary file
			dataFile := getLockPath(m.options.UniqueID) + ".data"
			data, err := os.ReadFile(dataFile)
			if err != nil {
				continue
			}

			// Clean up the data file
			os.Remove(dataFile)

			// Parse the data
			var secondInstanceData SecondInstanceData
			if err := json.Unmarshal(data, &secondInstanceData); err != nil {
				continue
			}

			// Call the callback
			m.options.OnSecondInstanceLaunch(secondInstanceData)
		}
	}()
}
