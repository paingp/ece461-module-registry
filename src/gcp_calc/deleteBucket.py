#!/usr/bin/env python3
import sys
from google.cloud import storage
import os

os.environ["GCLOUD_PROJECT"] = "ece461-module-registry"

def list_blobs(bucket_name):
    storage_client = storage.Client()

    # Note: Client.list_blobs requires at least package version 1.17.0.
    blobs = storage_client.list_blobs(bucket_name)

    # Note: The call returns a response only when the iterator is consumed.
    #for blob in blobs:
    #    print(blob.name)

    return blobs

def delete_objects(blobs):
    storage_client = storage.Client()

    #with storage_client.batch():
    for blob in blobs:
        generation_match_precondition = None

        # Optional: set a generation-match precondition to avoid potential race conditions
        # and data corruptions. The request to delete is aborted if the object's
        # generation number does not match your precondition.
        blob.reload()  # Fetch blob metadata to use in generation_match_precondition.
        generation_match_precondition = blob.generation

        blob.delete(if_generation_match=generation_match_precondition)

        print(f"Blob {blob.name} deleted.")


def delete_bucket(bucket_name):
    """Deletes a bucket. The bucket must be empty."""
    # bucket_name = "your-bucket-name"

    storage_client = storage.Client()

    bucket = storage_client.get_bucket(bucket_name)
    bucket.delete()

    print(f"Bucket {bucket.name} deleted")

if __name__ == "__main__":

    print("ive been here")

    bucket_name = sys.argv[1]
    blobs = list_blobs(bucket_name)
    delete_objects(blobs)
    # delete_bucket(bucket_name)
    