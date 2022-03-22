package httputils

import (
	"net/http"
	"time"

	"github.com/tenstad/jfrog-client-go/utils"
)

type HttpClientDetails struct {
	User        string
	Password    string
	ApiKey      string
	AccessToken string
	Headers     map[string]string
	Transport   *http.Transport
	HttpTimeout time.Duration
}

func (httpClientDetails HttpClientDetails) Clone() *HttpClientDetails {
	headers := make(map[string]string)
	utils.MergeMaps(httpClientDetails.Headers, headers)
	return &HttpClientDetails{
		User:        httpClientDetails.User,
		Password:    httpClientDetails.Password,
		ApiKey:      httpClientDetails.ApiKey,
		AccessToken: httpClientDetails.AccessToken,
		Headers:     headers}
}
