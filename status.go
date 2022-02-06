package hyper

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/infiniteloopcloud/log"
)

const (
	StatusAsLibrary = "lib"
	StatusAsService = "svc"
)

type StatusFn func() string

type StatusOpts struct {
	EnvironmentData []string
	Statuses        map[string]StatusFn
}

func StatusHandler(opts StatusOpts) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var envData = make(map[string]string)
		for _, e := range opts.EnvironmentData {
			envData[e] = os.Getenv(e)
		}

		var status = make(map[string]string)
		for name, fn := range opts.Statuses {
			status[name] = fn()
		}

		result := struct {
			EnvironmentData map[string]string `json:"environment_data"`
			Statuses        map[string]string `json:"statuses"`
		}{
			EnvironmentData: envData,
			Statuses:        status,
		}

		resp, err := json.Marshal(result)
		if err != nil {
			log.Debugf(r.Context(), "configDebug handler %s", err.Error())
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(resp)
		if err != nil {
			log.Debugf(r.Context(), "configDebug handler %s", err.Error())
		}
	}
}
