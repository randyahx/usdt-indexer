# USDT Indexer

This is an indexer written in Golang to subscribe to USDT Transfer events on the Ethereum Mainnet. The data is displayed on Grafana according to Openmetrics standards.  

## Running the program  
Add the required fields to the **config file** before running.  
```
go mod tidy
go build
./usdt-indexer --config-file ./config.json
```  

Note: This is an Interview assignment. Not ready for production.  
