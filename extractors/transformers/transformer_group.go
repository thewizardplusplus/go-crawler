package transformers

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-crawler/models"
)

// TransformerGroup ...
type TransformerGroup []models.LinkTransformer

// TransformLinks ...
func (transformers TransformerGroup) TransformLinks(
	links []string,
	response *http.Response,
	responseContent []byte,
) ([]string, error) {
	// processing of the transformers should be sequential
	// to one transformer can influence to another one
	transformedLinks := links
	for index, transformer := range transformers {
		var err error
		transformedLinks, err =
			transformer.TransformLinks(transformedLinks, response, responseContent)
		if err != nil {
			return nil,
				errors.Wrapf(err, "unable to transform links via transformer #%d", index)
		}
	}

	return transformedLinks, nil
}
