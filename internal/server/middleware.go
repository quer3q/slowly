package server

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"personal/slowly/pkg"
)

func timeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cCtx, cCancel := context.WithCancel(r.Context())
		tCtx, tCancel := context.WithTimeout(r.Context(), timeout)
		defer tCancel()

		tw := &timeoutWriter{rw: w, buf: &bytes.Buffer{}}

		go func() {
			next.ServeHTTP(tw, r.WithContext(cCtx))
			cCancel()
		}()

		select {
		case <-cCtx.Done():
			w.WriteHeader(tw.status)
			_, err := w.Write(tw.buf.Bytes())
			if err != nil {
				log.Printf("can't write data: %s", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		case <-tCtx.Done():
			cCancel()

			resp := &pkg.SlowResponse{Error: "timeout too long"}
			b, err := json.Marshal(resp)
			if err != nil {
				log.Printf("can't unmarshal SlowResponse: %s", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write(b)
			if err != nil {
				log.Printf("can't write data to SlowResponse: %s", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		}
	})
}

type timeoutWriter struct {
	rw http.ResponseWriter

	status int
	buf    *bytes.Buffer
}

func (tw timeoutWriter) Header() http.Header {
	return tw.rw.Header()
}

func (tw *timeoutWriter) WriteHeader(status int) {
	tw.status = status
}

func (tw *timeoutWriter) Write(b []byte) (int, error) {
	if tw.status == 0 {
		tw.status = http.StatusOK
	}

	return tw.buf.Write(b)
}
