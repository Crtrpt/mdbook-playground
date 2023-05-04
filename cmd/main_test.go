package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	resp, _ := http.PostForm("http://127.0.0.1:8080", url.Values{"image": {"go"}, "code": {"aaa"}})
	assert.Equal(t, resp.StatusCode, 200)
}
