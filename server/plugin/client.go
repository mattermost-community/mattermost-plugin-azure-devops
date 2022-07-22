package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/serializers"
	"github.com/pkg/errors"
)

type Client interface {
	TestApi() (string, error) // TODO: remove later
	GenerateOAuthToken(encodedFormValues string) (*serializers.OAuthSuccessResponse, error)
}

type client struct {
	plugin     *Plugin
	httpClient *http.Client
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// TODO: remove later
func (c *client) TestApi() (string, error) {
	return "hello world", nil
}

func (c *client) GenerateOAuthToken(encodedFormValues string) (*serializers.OAuthSuccessResponse, error) {
	var oAuthSuccessResponse *serializers.OAuthSuccessResponse

	_, err := c.callFormURLEncoded(constants.BaseOauthURL, constants.PathToken, http.MethodPost, nil, &oAuthSuccessResponse, encodedFormValues)
	if err != nil {
		return nil, err
	}

	return oAuthSuccessResponse, nil
}

// Wrapper to make REST API requests with "application/json" type content
func (c *client) callJSON(url, path, method string, in, out interface{}) (responseData []byte, err error) {
	contentType := "application/json"
	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(in)
	if err != nil {
		return nil, err
	}
	return c.call(url, method, path, contentType, buf, out, "")
}

// Wrapper to make REST API requests with "application/x-www-form-urlencoded" type content
func (c *client) callFormURLEncoded(url, path, method string, in, out interface{}, formValues string) (responseData []byte, err error) {
	contentType := "application/x-www-form-urlencoded"
	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(in)
	if err != nil {
		return nil, err
	}
	return c.call(url, method, path, contentType, buf, out, formValues)
}

// Makes HTTP request to REST APIs
func (c *client) call(basePath, method, path, contentType string, inBody io.Reader, out interface{}, formValues string) (responseData []byte, err error) {
	errContext := fmt.Sprintf("Azure Devops: Call failed: method:%s, path:%s", method, path)
	pathURL, err := url.Parse(path)
	if err != nil {
		return nil, errors.WithMessage(err, errContext)
	}

	if pathURL.Scheme == "" || pathURL.Host == "" {
		var baseURL *url.URL
		baseURL, err = url.Parse(basePath)
		if err != nil {
			return nil, errors.WithMessage(err, errContext)
		}
		if path[0] != '/' {
			path = "/" + path
		}
		path = baseURL.String() + path
	}

	var req *http.Request
	if formValues != "" {
		req, err = http.NewRequest(method, path, strings.NewReader(formValues))
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequest(method, path, inBody)
		if err != nil {
			return nil, err
		}
	}

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.Body == nil {
		return nil, nil
	}
	defer resp.Body.Close()

	responseData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		if out != nil {
			err = json.Unmarshal(responseData, out)
			if err != nil {
				return responseData, err
			}
		}
		return responseData, nil

	case http.StatusNoContent:
		return nil, nil

	case http.StatusNotFound:
		return nil, ErrNotFound
	}

	errResp := ErrorResponse{}
	err = json.Unmarshal(responseData, &errResp)
	if err != nil {
		return responseData, errors.WithMessagef(err, "status: %s", resp.Status)
	}
	return responseData, fmt.Errorf("errorMessage %s", errResp.Message)
}

func InitClient(p *Plugin) Client {
	return &client{
		plugin:     p,
		httpClient: &http.Client{},
	}
}
