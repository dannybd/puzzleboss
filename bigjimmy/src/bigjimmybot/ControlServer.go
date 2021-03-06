package bigjimmybot

import (
	"fmt"
	"net/http"
	"html"
	"encoding/json"
	"strconv"
	_ "github.com/go-sql-driver/mysql"
	l4g "code.google.com/p/log4go"
)

type BigJimmyControlData struct {
	Version string
}

var httpControlPort, httpControlPath string

func ControlServer(controlPort string, controlPath string) {
	httpControlPort = controlPort
	httpControlPath = controlPath
	startHttpControl()
}


func startHttpControl() {
	// Start an http server for control
	http.HandleFunc("/"+httpControlPath+"/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			log.Logf(l4g.INFO, "Got POST on control port to %v", html.EscapeString(r.URL.Path))
			
			switch r.URL.Path {
			case "/"+httpControlPath+"/version": {
					defer r.Body.Close()
					var data BigJimmyControlData
					err := json.NewDecoder(r.Body).Decode(&data)
					if err != nil {
						log.Logf(l4g.ERROR, "Error decoding JSON from body of POST: %v", err)
						w.WriteHeader(http.StatusBadRequest)
						fmt.Fprintf(w, "Error decoding JSON from body of POST: %v", err)
					}
					log.Logf(l4g.INFO, "Processed data from version POST: %+v", data)
					
					newVersion, err := strconv.ParseInt(data.Version, 10, 0)
					if err != nil {
						log.Logf(l4g.ERROR, "Could not parse version as int: %v", err)
						w.WriteHeader(http.StatusBadRequest)
						fmt.Fprintf(w, "Could not parse version as int: %v", err)

					}
					
					if newVersion > Version {
						// get and process version diff 
						//go PbGetVersionDiff(Version)
						// send the version down the channel for processing
						versionChan <- newVersion
					}
					w.WriteHeader(http.StatusOK)
					fmt.Fprintf(w, "Processed version POST: %+v", data)

				}
			default: {
					w.WriteHeader(http.StatusNotFound)
					
					fmt.Fprintf(w, "Got POST on control port to %v", html.EscapeString(r.URL.Path))
				}
			}
		} else {
			log.Logf(l4g.ERROR, "method %v not supported on control port (requested %v)", r.Method, html.EscapeString(r.URL.Path))
			w.WriteHeader(http.StatusNotImplemented)
			fmt.Fprintf(w, "method %v not supported on control port (requested %v)", r.Method, html.EscapeString(r.URL.Path))
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		log.Logf(l4g.ERROR, "returning StatusNotFound (requested %v)", html.EscapeString(r.URL.Path))
	})

	

	log.Logf(l4g.INFO, "starting http server on port %v", httpControlPort)
	l4g.Crashf("http server died: %v", http.ListenAndServe(":"+httpControlPort, nil))
}
