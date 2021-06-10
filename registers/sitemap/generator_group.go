package sitemap

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-crawler/registers"
	"golang.org/x/sync/errgroup"
)

// GeneratorGroup ...
type GeneratorGroup []registers.LinkGenerator

// GenerateLinks ...
func (generators GeneratorGroup) GenerateLinks(
	ctx context.Context,
	baseLink string,
) (
	[]string,
	error,
) {
	waiter, ctx := errgroup.WithContext(ctx)

	linkGroups := make([][]string, len(generators))
	for index, generator := range generators {
		index, generator := index, generator

		waiter.Go(func() error {
			links, err := generator.GenerateLinks(ctx, baseLink)
			if err != nil {
				return errors.Wrapf(err, "error with generator #%d", index)
			}

			linkGroups[index] = links
			return nil
		})
	}

	if err := waiter.Wait(); err != nil {
		return nil, errors.Wrap(err, "unable to generate Sitemap links")
	}

	var totalLinks []string
	for _, linkGroup := range linkGroups {
		totalLinks = append(totalLinks, linkGroup...)
	}

	return totalLinks, nil
}
