# USDT Indexer

This is an indexer written in Golang to subscribe to USDT Transfer events on the Ethereum Mainnet. The data is displayed on Grafana according to Openmetrics standards.  

## Running the program locally 
Add the required fields to the **config file** before running.  

### Without Docker
```
go mod tidy
go build
./usdt-indexer --config-file ./config.json
```

### With Docker  
Ensure that docker is running.  
```
docker build -t contract-indexer .
docker run -d --name contract-indexer -p 2112:2112 contract-indexer
```

Note: This is an Interview assignment. Not ready for production.  
