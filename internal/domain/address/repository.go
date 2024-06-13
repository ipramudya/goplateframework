package address

import "context"

type DBRepository interface {
	GetOneByID(ctx context.Context, id string) (*Schema, error)
	AddOne(ctx context.Context, s *NewAddressDTO) error
	Update(ctx context.Context, s *NewAddressDTO) error
}
