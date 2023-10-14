package conversion

import (
	"time"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetRate(country string, date time.Time) (string, error) {
	return "", nil
}
