package provisioning

import (
	"context"
	"github/sta-zot/pcs/internal/domain"
)

type fileGenerator interface {
	GenerateFile(config domain.PhoneSettings) ([]byte, error)
}

// Сервис отвечает за генерацию конфигурационных файлов для телефонов на основе переданной структуры с настройками
type service struct{}

func NewService(gen fileGenerator) *service {
	return &service{}
}

// GetConfig retrieves the configuration for the given MAC address and generates a file.
// Метод GetConfig получает структуру с конфигурацией для заданного MAC-адреса и генерирует файл.
// returns the generated file as a byte slice and any error that occurred during generation.
func (s *service) GetConfig(ctx context.Context, settings domain.PhoneSettings) ([]byte, error) {
	return nil, nil
}
