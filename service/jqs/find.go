package jqs

import (
	"context"

	"github.com/colbee1/jobq"
)

func (s *Service) FindByStatus(ctx context.Context, status jobq.Status, offset int, limit int) ([]jobq.ID, error) {
	tx, err := s.jobRepo.NewTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Close()

	return tx.FindByStatus(ctx, status, offset, limit)
}
