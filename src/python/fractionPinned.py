import requests
from dateutil import parser
import datetime
import sys
import os

headers = {"Authorization":  "token " + os.getenv("GITHUB_TOKEN")}

def run_query(query): # A simple function to use requests.post to make the API call. Note the json= section.
    request = requests.post('https://api.github.com/graphql', json={'query': query}, headers=headers)
    if request.status_code == 200:
        return request.json()
    else:
        raise Exception("Query failed to run by returning code of {}. {}".format(request.status_code, query))

def graphQL(url):
    query = """
{
  viewer {
    login
  }
  rateLimit {
    limit
    cost
    remaining
    resetAt
  }
}
"""

    result = run_query(query) # Execute the query
    remaining_rate_limit = result["data"]["rateLimit"]["remaining"] # Drill down the dictionary
    print("Remaining rate limit - {}".format(remaining_rate_limit))
  
if __name__ == "__main__":
    graphQL(sys.argv[1])
  #print(score)

