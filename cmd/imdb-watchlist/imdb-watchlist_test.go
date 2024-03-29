package main

import (
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"testing"
	"time"
)

func Test_Main(t *testing.T) {
	err := Main(slog.Default())
	assert.Error(t, err)

	*listID = "1234,5678"
	*debug = true
	go main()

	assert.Eventually(t, func() bool {
		resp, err := http.Get("http://localhost:9090/metrics")
		return err == nil && resp.StatusCode == http.StatusOK
	}, time.Second, time.Millisecond)
}
