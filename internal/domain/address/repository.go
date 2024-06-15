package address

import "context"

type DBRepository interface {
	Update(ctx context.Context, a *NewAddressDTO, id string) (*Schema, error)
}
