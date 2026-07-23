package repository

import (
	"context"

	"github.com/ruveydagundogan/cvmatcher/backend/internal/domain/jobdescription/model"
)

type JobDescriptionRepository interface {
	Save(ctx context.Context, jd *model.JobDescription) error
	FindByID(ctx context.Context, id string) (*model.JobDescription, error)
	FindByUserID(ctx context.Context, userID string, offset, limit int) ([]*model.JobDescription, int, error)
	Update(ctx context.Context, jd *model.JobDescription) error
	Delete(ctx context.Context, id string) error
}
