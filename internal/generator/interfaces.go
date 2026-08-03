package generator

import (
	"context"
	"github/sta-zot/pcs/internal/domain"
	"io"
)

// Generator factory interface for creating Generators based on model and vendor
type GeneratorFactory interface {
	Get(model, vendor string) (Generator, error)
}

// Generator interface for generating device configuration files
type Generator interface {
	Generate(ctx context.Context, settings *domain.PhoneSettings) (io.ByteReader, error)
}

// Logger interface for logging
type logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
}
