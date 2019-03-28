# Prometheus Exporter for Tibber

https://developer.tibber.com/explorer

## Get test token
https://developer.tibber.com/explorer

## Run
```bash
Usage:
  tibber-exporter [OPTIONS]

Application Options:
  -t, --token=    Authorization token [$TOKEN]
  -e, --endpoint= Endpoint (default: https://api.tibber.com/v1-beta/gql) [$ENDPOINT]

Help Options:
  -h, --help      Show this help message
```

All arguments can be set either by command line argument or environment

```bash
export TIBBER_TOKEN=d1007ead2dc84a2b82f0de19451c5fb22112f7ae11d19bf2bedb224a003ff74a

tibber-exporter
```


