package googlestorage

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/storage"
	"github.com/goplateframework/config"
	"google.golang.org/api/option"
)

func Init(ctx context.Context, conf *config.Config) (*storage.Client, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get current working directory error, %v", err)
	}

	storagePath := cwd + conf.GoogleStorage.Path
	clientOpt := option.WithCredentialsFile(storagePath)

	return storage.NewClient(ctx, clientOpt)
}
