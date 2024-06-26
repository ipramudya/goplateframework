package menu

import "context"

type DBRepository interface {
	AddOne(ctx context.Context, nm *NewMenuDTO) (*Schema, error)
	Update(ctx context.Context, nm *NewMenuDTO, id string) (*Schema, error)
	GetAllByOutletID(ctx context.Context, outletID string) (*[]Schema, error)
}
