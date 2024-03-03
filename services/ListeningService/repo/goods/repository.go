package goods

import (
	"context"
	"github.com/RustGrub/FunnyGoService/services/FunnyService/models"
)

type Repository interface {
	Set(ctx context.Context, good models.GoodAsLog) error
}
