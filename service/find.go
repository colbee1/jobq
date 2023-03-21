package service

import (
	"context"

	"github.com/colbee1/jobq"
)

func (s *Service) Find(ctx context.Context, jids []jobq.JobID) ([]*jobq.Job, error) {
	return s.jobRepo.Find(ctx, jids)
}

func (s *Service) FindByStatus(ctx context.Context, state jobq.JobStatus, offset int, limit int) ([]*jobq.Job, error) {
	return s.jobRepo.FindByStatus(ctx, state, offset, limit)
}
