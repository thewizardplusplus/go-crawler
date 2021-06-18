package sitemap

import (
	"compress/gzip"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	httputils "github.com/thewizardplusplus/go-http-utils"
)

// Loader ...
type Loader struct {
	HTTPClient httputils.HTTPClient
}

// LoadLink ...
func (loader Loader) LoadLink(link string, options interface{}) (
	[]byte,
	error,
) {
	request, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create the request")
	}
	request = request.WithContext(options.(context.Context))

	response, err := loader.HTTPClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "unable to send the request")
	}
	defer response.Body.Close() // nolint: errcheck

	responseReader := response.Body
	if response.Header.Get("Content-Encoding") == "gzip" {
		var err error
		responseReader, err = gzip.NewReader(responseReader)
		if err != nil {
			return nil,
				errors.Wrap(err, "unable to create the gzip reader for the response")
		}
		defer responseReader.Close() // nolint: errcheck
	}

	responseData, err := ioutil.ReadAll(responseReader)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read the response")
	}

	return responseData, nil
}
