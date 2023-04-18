package ratom

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

func CreateBucket(bucketName string) error {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := "ece461-module-registry"

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)

	// Creates the new bucket.
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := bucket.Create(ctx, projectID, nil); err != nil {
		return fmt.Errorf("Failed to create bucket: %v", err)
	}

	fmt.Printf("Bucket %v created.\n", bucketName)
	return nil
}

func uploadModule(module string, bucketName string) error {
	ctx := context.Background()

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Sets the name for the new bucket.
	//_, object, _ := strings.Cut(module, "/")
	object := module

	// Open local file.
	f, err := os.Open("temp/" + module + ".zip")
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	o := client.Bucket(bucketName).Object(object)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	// For an object that does not yet exist, set the DoesNotExist precondition.
	o = o.If(storage.Conditions{DoesNotExist: true})

	// If the live object already exists in your bucket, set instead a
	// generation-match precondition using the live object's generation number.
	// attrs, err := o.Attrs(ctx)
	// if err != nil {
	//      return fmt.Errorf("object.Attrs: %v", err)
	// }
	// o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	// Upload an object with storage.Writer.
	wc := o.NewWriter(ctx)
	var w bytes.Buffer
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	fmt.Fprintf(&w, "Blob %v uploaded.\n", object)

	return nil
}

// setMetadata sets an object's metadata.
func setMetadata(w io.Writer, bucket, object string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := client.Bucket(bucket).Object(object)

	// Optional: set a metageneration-match precondition to avoid potential race
	// conditions and data corruptions. The request to update is aborted if the
	// object's metageneration does not match your precondition.
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("object.Attrs: %v", err)
	}
	o = o.If(storage.Conditions{MetagenerationMatch: attrs.Metageneration})

	// Update the object to set the metadata.
	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		Metadata: map[string]string{
			"keyToAddOrUpdate": "value",
		},
	}
	if _, err := o.Update(ctx, objectAttrsToUpdate); err != nil {
		return fmt.Errorf("ObjectHandle(%q).Update: %v", object, err)
	}
	fmt.Fprintf(w, "Updated custom metadata for object %v in bucket %v.\n", object, bucket)
	return nil
}

func deleteBucket(bucketName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	bucket := client.Bucket(bucketName)
	if err := bucket.Delete(ctx); err != nil {
		return fmt.Errorf("Bucket(%q).Delete: %v", bucketName, err)
	}

	fmt.Printf("Bucket %v deleted\n", bucketName)
	return nil
}
