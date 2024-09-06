package client

import (
	"context"
	"contract-indexer/config"
	"contract-indexer/metrics"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type IndexerClient struct {
	ethClient     *ethclient.Client
	metricService *metrics.MetricService
	logsCh        chan types.Log
	cfg           *config.Config
}

func NewIndexerClient(cfg *config.Config, metricService *metrics.MetricService) *IndexerClient {
	// Initialize EthClient
	ethClient, err := ethclient.Dial(cfg.NodeConfig.EthereumNodeURL)
	if err != nil {
		panic(fmt.Errorf("indexer failed to connect to Ethereum client: %v", err))
	}
	println("Initialized connection to eth node url: ", cfg.NodeConfig.EthereumNodeURL)
	println("USDT Contract Address: ", cfg.ContractConfig.USDTContractAddress)
	// Events from the Listener will be passed into this for processing
	logsCh := make(chan types.Log)

	return &IndexerClient{
		ethClient:     ethClient,
		metricService: metricService,
		cfg:           cfg,
		logsCh:        logsCh,
	}
}

// ListenForERC20Transfers subscribes to ERC20 Transfer events
// Sends them to a log channel to process them
func (indexerClient *IndexerClient) ListenForERC20Transfers() {
	// Parse the ERC20 ABI
	contractABI, err := abi.JSON(strings.NewReader(indexerClient.cfg.ContractConfig.ERC20ABI))
	if err != nil {
		panic(fmt.Errorf("indexer failed to parse ERC20 ABI, err=%v", err))
	}

	// Set query filter
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(indexerClient.cfg.ContractConfig.USDTContractAddress)},
	}

	// Subscribe to filter logs
	sub, err := indexerClient.ethClient.SubscribeFilterLogs(context.Background(), query, indexerClient.logsCh)
	if err != nil {
		panic(fmt.Errorf("indexer failed to subscribe to logs, err=%v", err))
	}

	// Goroutine channel to handle logs
	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Fatalf("event listener error, err=%v", err)
			case eventLog := <-indexerClient.logsCh:
				indexerClient.processTransferEvent(contractABI, eventLog)
			}
		}
	}()
}

// processTransferEvent unpacks and processes the Transfer event.
func (indexerClient *IndexerClient) processTransferEvent(contractABI abi.ABI, eventLog types.Log) {
	// Event signature for the Transfer event
	transferEventSignature := []byte("Transfer(address,address,uint256)")
	transferEventHash := common.BytesToHash(crypto.Keccak256(transferEventSignature))

	// Skip events that are not Transfer events
	if eventLog.Topics[0] != transferEventHash {
		return
	}

	// Extract the indexed `from` and `to` addresses from the topics
	from := common.HexToAddress(eventLog.Topics[1].Hex())
	to := common.HexToAddress(eventLog.Topics[2].Hex())

	// Extract transfer value from event
	var transferData struct {
		Value *big.Int
	}
	err := contractABI.UnpackIntoInterface(&transferData, "Transfer", eventLog.Data)
	if err != nil {
		log.Printf("Failed to unpack log data: %v", err)
		return
	}

	convertedValueInUSD := float64(transferData.Value.Int64()) / math.Pow10(6)

	// Update metrics
	indexerClient.metricService.IncUSDTTxCount()
	indexerClient.metricService.IncUSDTTransferCount()
	indexerClient.metricService.IncUSDTTransferredAmountTotal(convertedValueInUSD)

	// Log the transfer event
	fmt.Printf("Transfer of %.2f tokens from %s to %s \n\n", convertedValueInUSD, from.Hex(), to.Hex())
}
