package config

import (
	"captcha-service/internal/domain/entity"
	"github.com/caarlos0/env/v11"
)

func LoadDemoConfig() (*entity.DemoConfig, error) {
	config := &entity.DemoConfig{}

	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}
