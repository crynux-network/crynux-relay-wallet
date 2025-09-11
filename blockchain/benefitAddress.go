package blockchain

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func GetBenefitAddress(ctx context.Context, nodeAddress common.Address, network string) (common.Address, error) {
	client, err := GetBlockchainClient(network)
	if err != nil {
		return common.Address{}, err
	}

	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	opts := &bind.CallOpts{
		Pending: false,
		Context: callCtx,
	}

	ba, err := client.BenefitAddressContractInstance.GetBenefitAddress(opts, nodeAddress)
	if err != nil {
		return common.Address{}, err
	}
	if ba.Big().Cmp(big.NewInt(0)) == 0 {
		return nodeAddress, nil
	} else {
		return ba, nil
	}
}
