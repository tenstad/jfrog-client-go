package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tenstad/jfrog-client-go/artifactory"
	"github.com/tenstad/jfrog-client-go/config"
)

func TestGetArtifactoryVersionWithCustomHttpClient(t *testing.T) {
	initArtifactoryTest(t)
	rtDetails := GetRtDetails()

	client := http.DefaultClient

	serviceConfig, err := config.NewConfigBuilder().
		SetServiceDetails(rtDetails).
		SetDryRun(false).
		SetHttpClient(client).
		Build()
	if err != nil {
		t.Error(err)
	}

	rtManager, err := artifactory.New(serviceConfig)
	if err != nil {
		t.Error(err)
	}

	version, err := rtManager.GetVersion()
	assert.NoError(t, err, "Should not fail")
	if version == "" {
		t.Error("Expected a version, got empty string")
	}
}

func TestGetArtifactoryVersionWithProxyShouldFail(t *testing.T) {
	initArtifactoryTest(t)
	rtDetails := GetRtDetails()

	proxyUrl, err := url.Parse("http://invalidproxy:12345")
	assert.NoError(t, err)
	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
	}

	serviceConfig, err := config.NewConfigBuilder().
		SetServiceDetails(rtDetails).
		SetDryRun(false).
		SetHttpClient(client).
		Build()
	if err != nil {
		t.Error(err)
	}

	rtManager, err := artifactory.New(serviceConfig)
	if err != nil {
		t.Error(err)
	}

	_, err = rtManager.GetVersion()
	assert.Error(t, err, "Should fail with invalid proxy")
}
