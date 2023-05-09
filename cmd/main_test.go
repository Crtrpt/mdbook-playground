package main

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/crtrpt/mdbook-playground/internal"
	"github.com/stretchr/testify/assert"
)

const urlEnd = "http://127.0.0.1:9080"

func TestMain(t *testing.T) {

	resp, _ := http.PostForm(urlEnd, url.Values{"image": {"centos:latest"}, "code": {"aaa"}})
	assert.Equal(t, resp.StatusCode, 200)
}

func TestStartC1(t *testing.T) {
	InitDocker()

	ctx := context.WithValue(context.Background(), "client", Cli)
	result, err := internal.StartC1(ctx, "centos:latest", "echo aaa")
	assert.Equal(t, err, nil)
	assert.Contains(t, result, "aaa")
}
