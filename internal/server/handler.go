package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"personal/slowly/pkg"
)

const apiSlowPath = "/api/slow"

var timeout = time.Second * 5

func serverHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != apiSlowPath || r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}

		req := &pkg.SlowRequest{}
		d := json.NewDecoder(r.Body)
		err := d.Decode(req)
		if err != nil {
			err = errors.New(fmt.Sprintf("can't parse request: %s", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		select {
		case <-time.After(time.Millisecond * time.Duration(req.Timeout)):
			resp := &pkg.SlowResponse{Status: "ok"}
			bytes, err := json.Marshal(resp)
			if err != nil {
				log.Printf("can't unmarshal SlowResponse: %s", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			_, err = w.Write(bytes)
			if err != nil {
				log.Printf("can't write data to SlowResponse: %s", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		}
	})
}
