package main

import (
	"context"
	"fmt"
	"github.com/crypto-smoke/blockchain"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math/big"
	"sync"
)

type pairData struct {
	block *big.Int
	count *big.Int
}
type reserves struct {
	Block              *big.Int
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}

func getLPfromPID(pid *big.Int, s *blockchain.Swap, c *ethclient.Client) (*blockchain.LiquidityPool, error) {
	pool, err := s.AllPairs(nil, pid)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get pool from pid")
	}
	lp, err := blockchain.NewLP(pool, c)
	if err != nil {
		return nil, errors.Wrap(err, "unable to getlp")
	}
	return lp, nil
}

func getPoolReservesAtBlock(lp *blockchain.LiquidityPool, block *big.Int) {

}
func priceOf(baseDecimals, quoteDecimals, reserveQuote, reserveBase *big.Int) *big.Int {
	baseExponent := new(big.Int).Exp(big.NewInt(10), baseDecimals, nil)
	input := new(big.Int).Mul(big.NewInt(1), baseExponent)
	output := amountOut(input, reserveBase, reserveQuote)
	//	fmt.Println("input", input)
	//	fmt.Println("output", output)
	quoteExponent := new(big.Int).Exp(big.NewInt(10), quoteDecimals, nil)
	x := new(big.Int).Div(output, quoteExponent)
	//	fmt.Println("x", x)
	return x
}
func main() {
	/*
		reserveA, success := new(big.Int).SetString("36054785717101931180384", 10)
		if !success {
			panic("no bueno")
		}
		reserveB, success := new(big.Int).SetString("56826509815735", 10)
		if !success {
			panic("no bueno")
		}

		priceOf(big.NewInt(18), big.NewInt(6), reserveB, reserveA)
		return
		//amountIn := big.NewInt(16000000)
		amountIn := big.NewInt(1000000000000000000)

		output := amountOut(amountIn, reserveA, reserveB)
		fmt.Println(output)
		fmt.Println(new(big.Int).Div(amountIn, output))
		return
		fmt.Println(calculatePrice(reserveA, reserveB))
		fmt.Println(calculatePrice(reserveB, reserveA))

		return


	*/
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

	// uniswap
	swap, err := blockchain.NewSwap(common.HexToAddress("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"), client)
	if err != nil {
		log.Fatal().Err(err).Msg("failed creating swap")
	}

	// weth / usdc
	lp, err := getLPfromPID(new(big.Int), swap, client)
	if err != nil {
		log.Fatal().Err(err).Msg("failed getting LP")
	}

	// 7 days worth of blocks
	var start uint64 = 60 * 60 * 24 * 7 / 13
	howManyWeWant := start

	maxGoroutines := 1000
	guard := make(chan struct{}, maxGoroutines)
	ch := make(chan reserves)
	var stuff = make(map[uint64]string)
	var wg sync.WaitGroup

	// read data
	go func() {
		wg.Add(1)
		defer wg.Done()
		var count uint64
		for {
			data := <-ch
			count++

			if data.Reserve0.Cmp(big.NewInt(0)) < 1 || data.Reserve1.Cmp(big.NewInt(0)) < 1 {
				continue
			}
			price := priceOf(big.NewInt(18), big.NewInt(6), data.Reserve0, data.Reserve1)
			//fmt.Println("price", price)
			//price := new(big.Int).Div(data.Reserve0, data.Reserve1)
			priceString := price.String()
			stuff[data.Block.Uint64()] = priceString
			//fmt.Println(priceString, new(big.Float).Quo(r1, r0).Text('f', -1))

			if count == howManyWeWant {
				fmt.Println("done reading")
				break
			}
		}
	}()

	//lastBlock := big.NewInt(10207858) // startblock of the uniswapv2 router2 contract

	// lastBlock := big.NewInt(10008355) // weth/usdc start
	//lastBlock := big.NewInt(10019997) // first LP approve
	//fmt.Println(blockNumber, start)
	//return
	lastBlock := big.NewInt(int64(blockNumber - start))
	for i := 0; i < int(howManyWeWant); i++ {
		guard <- struct{}{}
		wg.Add(1)
		go func(lastBlock *big.Int) {
			defer wg.Done()
			opts := bind.CallOpts{
				BlockNumber: lastBlock,
			}
			r, err := lp.GetReserves(&opts)
			if err != nil {
				log.Fatal().Err(err).Msg("failed getting swap pair length")
			}
			//fmt.Println(r.BlockTimestampLast, r.Reserve1, r.Reserve0)
			ch <- reserves{
				Block:              lastBlock,
				Reserve0:           r.Reserve0,
				Reserve1:           r.Reserve1,
				BlockTimestampLast: r.BlockTimestampLast,
			}
			<-guard
		}(lastBlock)

		lastBlock = new(big.Int).Add(lastBlock, big.NewInt(1))
	}

	fmt.Println("waiting")
	wg.Wait()
	fmt.Println("done")
	//for k, v := range stuff {
	//fmt.Println(k, v)
	//}
}
