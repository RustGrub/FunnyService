package goods

import (
	"context"
	"github.com/RustGrub/FunnyGoService/services/FunnyService/models"
)

type FsRepository interface {
	Create(ctx context.Context, req models.CreateGoodRequest) (res *models.Good, err error)
	Get(ctx context.Context, goodID, projectID int) (res *models.Good, err error)
	UpdateNameAndDescription(ctx context.Context, req models.UpdateGoodRequest) (res *models.Good, err error)
	Remove(ctx context.Context, req models.RemoveGoodRequest) (res *models.Good, err error)
	GetListByLimitAndOffset(ctx context.Context, limit, offset int) (res []models.Good, removed int, err error)
	Reprioritize(ctx context.Context, newP, goodP int, up bool) (res []models.Good, err error)
	UpdatePriority(ctx context.Context, goodID, projectID, priority int) (res *models.Good, err error)
}

type Cache interface {
	Set(ctx context.Context, good models.Good) error
	Get(ctx context.Context, goodID, projectID int) *models.Good
}
