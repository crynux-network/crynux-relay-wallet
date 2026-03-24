package config

import (
	"crypto/ecdsa"
	"errors"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"
)

var appConfig *AppConfig

// InitConfig Init is an exported method that takes the config from the config file
// and unmarshal it into AppConfig struct
func InitConfig(configPath string) error {
	v := viper.New()
	v.SetConfigType("yml")
	v.SetConfigName("config")

	if configPath != "" {
		v.AddConfigPath(configPath)
	} else {
		v.AddConfigPath("/app/config")
		v.AddConfigPath("config")
	}

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	appConfig = &AppConfig{}

	if err := v.Unmarshal(appConfig); err != nil {
		return err
	}

	if appConfig.Environment == EnvTest {
		privKey := GetTestPrivateKey()
		for network := range appConfig.Blockchains {
			blockchain := appConfig.Blockchains[network]
			blockchain.Account.PrivateKey = privKey
			appConfig.Blockchains[network] = blockchain
		}
		appConfig.Relay.Api.PrivateKey = privKey
	} else {
		// Load hard-coded private key
		for network := range appConfig.Blockchains {
			blockchain := appConfig.Blockchains[network]
			blockchain.Account.PrivateKey = GetPrivateKey(blockchain.Account.PrivateKeyFile)
			appConfig.Blockchains[network] = blockchain
		}
		appConfig.Relay.Api.PrivateKey = GetPrivateKey(appConfig.Relay.Api.PrivateKeyFile)
	}
	if err := checkBlockchainAccount(); err != nil {
		return err
	}
	if err := checkRelayConfig(); err != nil {
		return err
	}

	return nil
}

func checkBlockchainAccount() error {

	for _, blockchain := range appConfig.Blockchains {
		if blockchain.Account.PrivateKey == "" {
			return errors.New("blockchain account private key not set")
		}

		if blockchain.Account.Address == "" {
			return errors.New("blockchain account address not set")
		}

		// Check private key and address
		privateKey, err := crypto.HexToECDSA(blockchain.Account.PrivateKey)
		if err != nil {
			return err
		}

		publicKey := privateKey.Public()

		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			return errors.New("error casting public key to ECDSA")
		}

		address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

		if address != blockchain.Account.Address {
			return errors.New("account address and private key mismatch")
		}
	}

	return nil
}

func checkRelayConfig() error {
	depositAddress := appConfig.Relay.DepositAddress
	if depositAddress == "" {
		return errors.New("relay deposit address not set")
	}
	if !common.IsHexAddress(depositAddress) {
		return errors.New("relay deposit address is invalid")
	}
	return nil
}

func GetPrivateKey(file string) string {
	b, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return normalizePrivateKey(string(b))
}

func normalizePrivateKey(privateKey string) string {
	privateKey = strings.TrimSpace(privateKey)
	if strings.HasPrefix(privateKey, "0x") || strings.HasPrefix(privateKey, "0X") {
		return privateKey[2:]
	}
	return privateKey
}

func GetTestPrivateKey() string {
	return ""
}

func GetConfig() *AppConfig {
	return appConfig
}
