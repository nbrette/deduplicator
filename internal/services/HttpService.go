package services

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

type HTTPService struct {
	client *http.Client
}

func NewHTTPService() *HTTPService {
	return &HTTPService{client: http.DefaultClient}
}

//go:generate mockgen -destination=$PROJECT_PATH/internals/tests/mocks/mock_IHTTPService.go -package=mocks . IHTTPService
type IHTTPService interface {
	Post(url string, body []byte, headers map[string]string) error
}

func (httpService HTTPService) Post(url string, body []byte, headers map[string]string) error {
	request, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	for hKey, hVal := range headers {
		request.Header.Add(hKey, hVal)
	}
	res, err := httpService.client.Do(request)
	if err == nil {
		defer res.Body.Close()
		if res.StatusCode < 200 || res.StatusCode > 206 {
			return fmt.Errorf("destination returned a code other than 2XX. Response: %v", res)
		}
	}

	return err
}
