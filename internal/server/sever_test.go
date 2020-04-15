package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"personal/slowly/pkg"

	"github.com/stretchr/testify/require"
)

func Test_Server(t *testing.T) {
	const url = "http://" + listenAddr + apiSlowPath
	const contentType = "application/json"
	go Start()

	t.Run("ok response for correct timeout", func(t *testing.T) {
		t.Parallel()
		var correctTimeout = timeout.Milliseconds() - 1000
		slowReq := &pkg.SlowRequest{Timeout: correctTimeout}
		req, err := json.Marshal(slowReq)
		require.NoError(t, err)

		resp, err := http.Post(url, contentType, bytes.NewBuffer(req))
		require.NoError(t, err)
		d := json.NewDecoder(resp.Body)
		slowResp := &pkg.SlowResponse{}
		err = d.Decode(slowResp)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "ok", slowResp.Status)
	})

	t.Run("timeout response for exceed timeout value", func(t *testing.T) {
		t.Parallel()
		var exceedTimeout = timeout.Milliseconds() + 1000
		slowReq := &pkg.SlowRequest{Timeout: exceedTimeout}
		req, err := json.Marshal(slowReq)
		require.NoError(t, err)

		resp, err := http.Post(url, contentType, bytes.NewBuffer(req))
		require.NoError(t, err)
		d := json.NewDecoder(resp.Body)
		slowResp := &pkg.SlowResponse{}
		err = d.Decode(slowResp)
		require.NoError(t, err)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		require.Equal(t, "timeout too long", slowResp.Error)
	})
}
