package grpcserver

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/goplateframework/internal/worker/pb"
)

// https://storage.googleapis.com/goplateframework.appspot.com/menu_topings_images/d83a8d9e-2d9c-401a-a891-1c5527d5810b.webp

func (s *server) DeleteImage(ctx context.Context, req *pb.DeleteImageRequest) (*pb.Empty, error) {
	now := time.Now()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer func() {
		cancel()
		s.log.Infof("DeleteImage took %s", time.Since(now))
	}()

	trimmedUrl := strings.TrimPrefix(req.ImageUrl, FirebaseStorageDomain)

	fullpaths := strings.Split(trimmedUrl, "/")
	folderName := fullpaths[len(fullpaths)-2]
	fileName := fullpaths[len(fullpaths)-1]

	bucket := s.storage.Bucket(s.conf.GoogleStorage.BucketName)
	object := bucket.Object(fmt.Sprintf("%s/%s", folderName, fileName))

	_, err := object.Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			s.log.Errorf("Image %s does not exist", req.ImageUrl)
		}

		s.log.Errorf("Error getting image %s: %v", req.ImageUrl, err)
		return &pb.Empty{}, nil
	}

	if err := object.Delete(ctx); err != nil {
		s.log.Errorf("Error deleting image %s: %v", req.ImageUrl, err)
		return &pb.Empty{}, nil
	}

	return &pb.Empty{}, nil
}
