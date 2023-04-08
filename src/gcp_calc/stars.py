from google.cloud import storage

def increment_download(bucket_name, blob):
    """Increments download metadata"""

    # Connects to storage client
    storage_client = storage.Client()

    # Selects bucket
    bucket = storage_client.get_bucket(bucket_name)

    # Selects blob
    blob = bucket.get_blob(blob)

    # Race condition handler
    metageneration_match_precondition = blob.metageneration
    metadata = blob.metadata

    # Increment downloads by 1
    metadata["downloads"] = metadata["downloads"] + 1
    blob.metadata = metadata
    blob.patch(if_metageneration_match=metageneration_match_precondition)

    return 1

def setup_blob_metadata(bucket_name, blob):
    """Initializes blob metadata (downloads, stars)"""

    # Connects to storage client
    storage_client = storage.Client()

    # Selects bucket
    bucket = storage_client.get_bucket(bucket_name)

    # Selects blob
    blob = bucket.get_blob(blob)

    # Race condition handler
    metageneration_match_precondition = blob.metageneration
    metadata = blob.metadata

    # Initializes all metadata to default values
    metadata["downloads"] = 0
    metadata["stars"] = 0

    blob.metadata = metadata
    blob.patch(if_metageneration_match=metageneration_match_precondition)

def increment_stars(bucket_name, blob):
    """Initializes blob metadata (downloads, stars, popularity)"""

    # Connects to storage client
    storage_client = storage.Client()

    # Selects bucket
    bucket = storage_client.get_bucket(bucket_name)

    # Selects blob
    blob = bucket.get_blob(blob)

    # Race condition handler
    metageneration_match_precondition = blob.metageneration
    metadata = blob.metadata

    # Increment Stars by 1
    metadata["stars"] = metadata["stars"] + 1
    
    blob.metadata = metadata
    blob.patch(if_metageneration_match=metageneration_match_precondition)

def get_popularity(bucket_name, blob):
    """Returns the popularity values: (stars, downloads)"""

    # Connects to storage client
    storage_client = storage.Client()

    # Selects bucket
    bucket = storage_client.get_bucket(bucket_name)

    # Selects blob
    blob = bucket.get_blob(blob)

    # Race condition handler
    metageneration_match_precondition = blob.metageneration

    stars = blob.metadata["stars"]
    downloads = blob.metadata["downloads"]

    return stars, downloads

def set_version(bucket_name, blob, version):
    """Returns the popularity values: (stars, downloads)"""

    # Connects to storage client
    storage_client = storage.Client()

    # Selects bucket
    bucket = storage_client.get_bucket(bucket_name)

    # Selects blob
    blob = bucket.get_blob(blob)

    # Race condition handler
    metageneration_match_precondition = blob.metageneration

    metadata = blob.metadata

    metadata["version"] = version

    blob.metadata = metadata
    blob.patch(if_metageneration_match=metageneration_match_precondition)



set_version("tmr-bucket", "lodash.txt", "1.2.4")
