package main

import (
	"context"
	"fmt"
	"io"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)



// DetectText gets text from the Vision API for an image at the given file path.
func DetectText(w io.Writer, file string) ([]*pb.EntityAnnotation, error) {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return nil, err
	}
	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		return nil, err
	}

	if len(annotations) == 0 {

		return nil, fmt.Errorf("no text found")
	}

	return annotations, nil
}
