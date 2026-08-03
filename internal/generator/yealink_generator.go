package generator

import (
	"context"
	"github/sta-zot/pcs/internal/domain"
	"io"
)

type yealinkGenerator struct{}

func (yg *yealinkGenerator) Generate(ctx context.Context, settings *domain.PhoneSettings) (io.ByteReader, error) {

	return nil, nil
}

func NewYealinkGenerator() *yealinkGenerator {
	return &yealinkGenerator{}
}
