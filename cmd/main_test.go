package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const urlEnd = "http://127.0.0.1:8080"

func TestMain(t *testing.T) {
	resp, _ := http.PostForm(urlEnd, url.Values{"image": {"go"}, "code": {"aaa"}})
	assert.Equal(t, resp.StatusCode, 200)
}
