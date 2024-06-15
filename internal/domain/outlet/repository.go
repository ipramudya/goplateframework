package outlet

import "context"

type DBRepository interface {
	GetOneByID(ctx context.Context, id string) (*SchemaWithAddress, error)
	AddOne(ctx context.Context, s *NewOutletDTO) (*SchemaWithAddress, error)
	// Update(ctx context.Context, s *NewOutletDTO) error
}
