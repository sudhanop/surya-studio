package services

import (
	"context"

	"github.com/suryaphotography/backend/internal/models"
	"github.com/suryaphotography/backend/internal/repositories"
)

type DashboardService struct {
	funcRepo    *repositories.FunctionRepo
	inquiryRepo *repositories.InquiryRepo
	mediaRepo   *repositories.MediaRepo
}

func NewDashboardService(fr *repositories.FunctionRepo, ir *repositories.InquiryRepo, mr *repositories.MediaRepo) *DashboardService {
	return &DashboardService{funcRepo: fr, inquiryRepo: ir, mediaRepo: mr}
}

func (s *DashboardService) Stats(ctx context.Context) (*models.DashboardStats, error) {
	stats := &models.DashboardStats{}
	var err error
	if stats.UpcomingFunctions, err = s.funcRepo.CountUpcoming(ctx); err != nil {
		return nil, err
	}
	if stats.RecentInquiries, err = s.inquiryRepo.CountRecent(ctx, 7); err != nil {
		return nil, err
	}
	if stats.PendingAlbums, err = s.funcRepo.CountPendingAlbums(ctx); err != nil {
		return nil, err
	}
	if stats.PendingVideoEdits, err = s.funcRepo.CountPendingVideos(ctx); err != nil {
		return nil, err
	}
	if stats.PendingDeliveries, err = s.funcRepo.CountPendingDeliveries(ctx); err != nil {
		return nil, err
	}
	if stats.TotalUploads, err = s.mediaRepo.CountAll(ctx); err != nil {
		return nil, err
	}
	if stats.LatestPortfolioCount, err = s.mediaRepo.CountRecent(ctx, 7); err != nil {
		return nil, err
	}
	return stats, nil
}
