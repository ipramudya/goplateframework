package grpcserver

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"path"
	"time"

	"cloud.google.com/go/storage"
	"github.com/disintegration/imaging"
	"github.com/goplateframework/internal/worker/pb"
	webpEncoder "github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

const (
	MaxWidthSupport       int    = 600
	DefaultQuality        int    = 75
	FirebaseStorageDomain string = "https://storage.googleapis.com"
)

func (s *server) ProcessImage(ctx context.Context, req *pb.ProcessImageRequest) (*pb.Empty, error) {
	now := time.Now()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer func() {
		cancel()
		s.log.Infof("ProcessImage took %s", time.Since(now))
	}()

	optimized, err := s.optimize(req.ImageData)
	if err != nil {
		s.log.Errorf("error optimizing image: %v", err)
		return nil, err
	}

	imageUrl, err := s.storeImage(ctx, optimized, req.Table, req.Id)
	if err != nil {
		s.log.Errorf("error storing image: %v", err)
		return nil, err
	}

	if err := s.writeIntoDatabase(ctx, imageUrl, req.Table, req.Id); err != nil {
		s.log.Errorf("error writing into database: %v", err)
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) optimize(imageData []byte) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %v", err)
	}

	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()
	var newWidth, newHeight int

	if originalWidth > MaxWidthSupport {
		newWidth = MaxWidthSupport
		newHeight = int((float64(MaxWidthSupport) / float64(originalWidth)) * float64(originalHeight))
	} else {
		// If the image is already smaller than 600px wide, keep original size
		newWidth = originalWidth
		newHeight = originalHeight
	}

	img = imaging.Resize(img, newWidth, newHeight, imaging.Linear)

	// re-encode image with optimized settings
	var buf bytes.Buffer
	switch format {
	case "jpeg":
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: DefaultQuality})
	case "png":
		encoder := png.Encoder{CompressionLevel: png.BestCompression}
		err = encoder.Encode(&buf, img)
	default:
		buf = *bytes.NewBuffer(imageData)
	}

	if err != nil {
		return nil, fmt.Errorf("error encoding image: %v", err)
	}

	// force to convert image into webp
	var webpBuf bytes.Buffer
	opts, err := webpEncoder.NewLossyEncoderOptions(
		webpEncoder.PresetDefault,
		float32(DefaultQuality),
	)

	if err != nil {
		return nil, fmt.Errorf("error creating webp encoder options: %v", err)
	}

	if err = webp.Encode(&webpBuf, img, opts); err != nil {
		return nil, fmt.Errorf("error encoding image: %v", err)
	}

	return webpBuf.Bytes(), nil
}

func (s *server) storeImage(ctx context.Context, imageData []byte, tableName string, id string) (string, error) {
	folderName := fmt.Sprintf("%s_images", tableName)
	fileName := fmt.Sprintf("%s.webp", id)
	fullPath := path.Join(folderName, fileName)

	bucket := s.storage.Bucket(s.conf.GoogleStorage.BucketName)
	object := bucket.Object(fullPath)

	writer := object.NewWriter(ctx)
	writer.ContentType = "image/webp"

	if _, err := io.Copy(writer, bytes.NewReader(imageData)); err != nil {
		return "", fmt.Errorf("error copying image data: %v", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("error closing writer: %v", err)
	}

	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("error setting ACL: %v", err)
	}

	return fmt.Sprintf("%s/%s/%s", FirebaseStorageDomain, s.conf.GoogleStorage.BucketName, fullPath), nil
}

func (s *server) writeIntoDatabase(ctx context.Context, imageUrl string, tableName string, id string) error {
	query := fmt.Sprintf("UPDATE %s SET image_url = $1 WHERE id = $2", tableName)

	_, err := s.db.ExecContext(ctx, query, imageUrl, id)
	if err != nil {
		return err
	}

	return nil
}
