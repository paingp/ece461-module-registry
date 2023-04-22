import os
import re
from google.cloud import storage

os.environ["GCLOUD_PROJECT"] = "ece461-module-registry"
bucket_name = "tmr-bucket"

def rating_init():
    # Connects to storage client
    storage_client = storage.Client()

    blobs = storage_client.list_blobs(bucket_name)

    for blob in blobs:

        metageneration_match_precondition = blob.metageneration

        metadata = blob.metadata

        metadata['NET_SCORE'] = 0.0
        metadata['RESPONSIVE_MAINTAINER'] = 0.0
        metadata['RAMP_UP'] = 0.0
        metadata['BUS_FACTOR'] = 0.0
        metadata['CORRECTNESS'] = 0.0
        metadata['LICENSE'] = 0.0
        metadata['DEPENDENCIES'] = 0.0
        metadata['PULL_REQ_LOC'] = 0.0

        blob.metadata = metadata
        blob.patch(if_metageneration_match=metageneration_match_precondition)

rating_init()