package generator

import (
	"context"
	"github/sta-zot/pcs/internal/domain"
	"io"
)

type generationServise struct {
}

func (gs *generationServise) Generate(ctx context.Context, settings *domain.PhoneSettings) (io.ByteReader, error) {
	generator, err := generatorFactory(settings.Vendor(), settings.Model())
	if err != nil {
		return nil, err
	}
	return generator.Generate(ctx, settings)
}
