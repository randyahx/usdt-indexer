# USDT Indexer

This is an indexer written in Golang to subscribe to USDT Transfer events on the Ethereum Mainnet. The data is displayed on Grafana according to Openmetrics standards.  

## Running the program locally 
Add the required fields to the **config file** before running.  

### Without Docker
```
go mod tidy
make build
./usdt-indexer --config-file ./config.json
```

### With Docker  
Ensure that docker is running and port 2112 is available.  
```
make build_docker
make run_docker
```

Note: This is an Interview assignment. Not ready for production.  
