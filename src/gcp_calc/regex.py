import os
import re
from google.cloud import storage

os.environ["GCLOUD_PROJECT"] = "ece461-module-registry"
bucket_name = "tmr-bucket"

def find_regex_matches(regex_string):

    if (len(regex_string) == 0):
        return ""

    storage_client = storage.Client()
    blobs = storage_client.list_blobs(bucket_name)

    regex_pattern = re.compile(regex_string)

    results = []

    for blob in blobs:

        blobWritten = 0
        result_append = []
        
        if regex_pattern.match(blob.name) != None:
            moduleName = blob.name[:-4]
            result_append.append(moduleName)

            moduleVersion = (blob.metadata)['version']
            result_append.append(moduleVersion)

            blobWritten = 1
        
        if regex_pattern.match((blob.metadata)['README']) != None and blobWritten == 0:
            moduleName = blob.name[:-4]
            result_append.append(moduleName)

            moduleVersion = (blob.metadata)['version']
            result_append.append(moduleVersion)

            blobWritten = 1

        if blobWritten == 1:
            results.append(result_append)

    return results

print(find_regex_matches("require('lodash/fp/curryN')"))