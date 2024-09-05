package config

import (
	"encoding/json"
	"os"
)

// Config holds the application configuration
type Config struct {
	NodeConfig     NodeConfig     `json:"node_config"`
	ContractConfig ContractConfig `json:"contract_config"`
	MetricsConfig  MetricsConfig  `json:"metrics_config"`
}

type NodeConfig struct {
	EthereumNodeURL string `json:"ethereum_node_url"`
}

type ContractConfig struct {
	USDTContractAddress string `json:"usdt_contract_address"`
	ERC20ABI            string `json:"erc20_abi"`
}

type MetricsConfig struct {
	Port uint16 `json:"port"`
}

// For local config files
func ParseConfigFromFile(configPath string) (*Config, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		panic(err)
	}
	return &config, nil
}
