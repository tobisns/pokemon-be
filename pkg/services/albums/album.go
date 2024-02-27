package albums

import (
	"context"
)

type Repo interface {
	Get(ctx context.Context, id string) (Album, error)
}

type Service interface {
	Get(ctx context.Context, id string) (Album, error)
}

type album struct {
	repo Repo
}

// New Service instance
func New(repo Repo) Service {
	return &album{repo}
}

// Get sends the request straight to the repo
func (s *album) Get(ctx context.Context, id string) (Album, error) {
	return s.repo.Get(ctx, id)
}
