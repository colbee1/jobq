package jqs

import (
	"context"

	"github.com/colbee1/jobq"
)

func (s *Service) ListByStatus(ctx context.Context, status jobq.JobStatus, offset int, limit int) ([]jobq.JobID, error) {
	tx, err := s.jobRepo.NewTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Close()

	return tx.ListByStatus(ctx, status, offset, limit)
}

func (s *Service) GetInfos(ctx context.Context, jids []jobq.JobID) ([]*jobq.JobInfo, error) {
	tx, err := s.jobRepo.NewTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Close()

	return tx.GetInfos(ctx, jids)
}
