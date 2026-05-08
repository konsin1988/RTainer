package health

import "context"

type Repository interface {
  Ping(ctx context.Context) error
}

type Service struct {
  repo Repository
}

func NewService(repo Repository) *Service {
  return &Service{repo: repo}
}

func (s *Service) Check(ctx context.Context) error {
  return s.repo.Ping(ctx)
}
