package address

import "context"

type DBRepository interface {
	Update(ctx context.Context, s *AddressDTO) error
}
