package ratom

import (
	"context"
	"fmt"

	// "io"
	"time"

	"cloud.google.com/go/storage"
)

func SetMetadata(bucket, object string, r *Repo) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := client.Bucket(bucket).Object(object)

	attrs, err := o.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("object.Attrs: %v", err)
	}
	o = o.If(storage.Conditions{MetagenerationMatch: attrs.Metageneration})

	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		Metadata: map[string]string{
			"NET_SCORE": fmt.Sprintf("%f", r.NetScore), 
			"RESPONSIVE_MAINTAINER": fmt.Sprintf("%f", r.Responsiveness),
			"RAMP_UP" : fmt.Sprintf("%f", r.RampUpTime),
			"BUS_FACTOR" : fmt.Sprintf("%f", r.BusFactor),
			"CORRECTNESS" : fmt.Sprintf("%f", r.Correctness),
			"LICENSE" : fmt.Sprintf("%f", r.LicenseCompatibility),
			"DEPENDENCIES" : fmt.Sprintf("%f", r.Dependency),
			"PULL_REQ_LOC" : fmt.Sprintf("%f", r.LocPRCR),
		},
	}

	if _, err := o.Update(ctx, objectAttrsToUpdate); err != nil {
		return fmt.Errorf("ObjectHandle(%q).Update: %v", object, err)
	}

	return nil
}
