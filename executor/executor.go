package executor

import (
	"context"
)

type Service struct {
	cfg     *Config
	watcher *Watcher
}

func NewService(cfg *Config) *Service {
	watcher, err := NewWatcher()
	if err != nil {
		panic(err)
	}

	return &Service{
		cfg:     cfg,
		watcher: watcher,
	}
}

func (s *Service) Run(ctx context.Context) error {
	go func() {
		select {
		case <-ctx.Done():
			_ = s.watcher.Close()
		}
	}()

	for _, rule := range s.cfg.Rules {
		err := s.watcher.Add(rule)
		if err != nil {
			return err
		}
	}

	return s.watcher.Watch()
}
