package app

import (
	"contract-indexer/client"
	"contract-indexer/config"
	"contract-indexer/metrics"
)

type App struct {
	indexerClient *client.IndexerClient
	metricService *metrics.MetricService
}

func NewApp(cfg *config.Config) *App {
	metricService := metrics.NewMetricService(cfg)
	indexerClient := client.NewIndexerClient(cfg, metricService)

	return &App{
		metricService: metricService,
		indexerClient: indexerClient,
	}
}

func (app *App) Start() {
	go app.metricService.Start()
	app.indexerClient.ListenForERC20Transfers()
}
