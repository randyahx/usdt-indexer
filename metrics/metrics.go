package metrics

import (
	"contract-indexer/config"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

const (
	MetricUSDTTransferCount          = "usdt_transfer_count"
	MetricUSDTTransferredAmountTotal = "usdt_transferred_amount_total"
	MetricUSDTTxCount                = "usdt_tx_per_second"
)

type MetricService struct {
	MetricsMap map[string]prometheus.Metric
	cfg        *config.Config
}

func NewMetricService(cfg *config.Config) *MetricService {
	metricsMap := make(map[string]prometheus.Metric)

	usdtTransfersCount := prometheus.NewCounter(prometheus.CounterOpts{
		Name: MetricUSDTTransferCount,
		Help: "Total number of USDT transfers",
	})
	metricsMap[MetricUSDTTransferCount] = usdtTransfersCount
	prometheus.MustRegister(usdtTransfersCount)

	usdtTransferredAmountTotal := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: MetricUSDTTransferredAmountTotal,
		Help: "Total amount of USDT tokens transferred",
	})
	metricsMap[MetricUSDTTransferredAmountTotal] = usdtTransferredAmountTotal
	prometheus.MustRegister(usdtTransferredAmountTotal)

	usdtTxCount := prometheus.NewCounter(prometheus.CounterOpts{
		Name: MetricUSDTTxCount,
		Help: "Number of USDT transactions per second",
	})
	metricsMap[MetricUSDTTxCount] = usdtTxCount
	prometheus.MustRegister(usdtTxCount)

	return &MetricService{
		MetricsMap: metricsMap,
		cfg:        cfg,
	}
}

func (m *MetricService) Start() {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(fmt.Sprintf(":%d", m.cfg.MetricsConfig.Port), nil)
	if err != nil {
		panic(err)
	}
}

func (metricService *MetricService) IncUSDTTransferCount() {
	metricService.MetricsMap[MetricUSDTTransferCount].(prometheus.Counter).Inc()
}

func (metricService *MetricService) IncUSDTTransferredAmountTotal(amount float64) {
	metricService.MetricsMap[MetricUSDTTransferredAmountTotal].(prometheus.Gauge).Add(amount)
}

func (metricService *MetricService) IncUSDTTxCount() {
	metricService.MetricsMap[MetricUSDTTxCount].(prometheus.Counter).Inc()
}
