# IntelAPI

### Context

The IntelAPI gathers known bad IP addresses from a set of DataSources. The data is made available via the **/blocklist** endpoint.
Current version contains the following DataSources:
* blocklist.de
* abuseipdb

The latter requires an API key. We can set this API key using the following environment variable:

```
export ABUSEIPDB_API_TOKEN="API_TOKEN"
```

### Run locally

We can use **run.sh** to launch a Docker container with the API installed. The container will expose port 8080 to the
local system.

```
chmod +x ./run.sh
./run.sh
```
