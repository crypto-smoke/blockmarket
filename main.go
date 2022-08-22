package main

import (
	"context"
	"fmt"
	"github.com/crypto-smoke/blockchain"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math/big"
)

func main() {
	log.Logger = zerolog.New(zerolog.NewConsoleWriter())

	client, err := ethclient.Dial("https://rpc.ankr.com/eth")
	if err != nil {
		log.Fatal().Err(err).Msg("failed dialing node")
	}
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get block number")
	}
	log.Info().
		Uint64("block", blockNumber).
		Msg("got current block")

	swap, err := blockchain.NewSwap(common.HexToAddress("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"), client)
	if err != nil {
		log.Fatal().Err(err).Msg("failed creating swap")
	}

	lastLength := new(big.Int)

	lastBlock := big.NewInt(10207858)
	for i := 0; i < 1000; i++ {

		opts := bind.CallOpts{
			Pending:     false,
			From:        common.Address{},
			BlockNumber: lastBlock,
			Context:     nil,
		}
		length, err := swap.AllPairsLength(&opts)
		if err != nil {
			log.Fatal().Err(err).Msg("failed getting swap pair length")
		}
		if length.Uint64() != lastLength.Uint64() {
			log.Info().Int64("block", lastBlock.Int64()).Uint64("length", length.Uint64()).Msg("new pairs")
			lastLength = length
		}
		fmt.Println(length.String())
		lastBlock = new(big.Int).Add(lastBlock, big.NewInt(1))
	}

}
