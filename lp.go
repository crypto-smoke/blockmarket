package main

import "math/big"

/*
   // given an input amount of an asset and pair reserves, returns the maximum output amount of the other asset
   function getAmountOut(uint amountIn, uint reserveIn, uint reserveOut) internal pure returns (uint amountOut) {
       require(amountIn > 0, 'UniswapV2Library: INSUFFICIENT_INPUT_AMOUNT');
       require(reserveIn > 0 && reserveOut > 0, 'UniswapV2Library: INSUFFICIENT_LIQUIDITY');
       uint amountInWithFee = amountIn.mul(997);
       uint numerator = amountInWithFee.mul(reserveOut);
       uint denominator = reserveIn.mul(1000).add(amountInWithFee);
       amountOut = numerator / denominator;
   }

amountOut = amountIn.mul(997).mul(reserveOut) / reserveIn.mul(1000).add(amountIn.mul(997))
*/
func amountOut(amountIn, reserveIn, reserveOut *big.Int) *big.Int {
	// feeless price = (reserveOut * 1000) / (reserveIn * 1000 + reserveOut)
	amountInWithFee := new(big.Int).Mul(amountIn, big.NewInt(997))
	numerator := new(big.Int).Mul(amountInWithFee, reserveOut)
	denominator := new(big.Int).Mul(reserveIn, big.NewInt(1000))
	denominator = new(big.Int).Add(denominator, amountInWithFee)

	return new(big.Int).Div(numerator, denominator)
	/*
		return new(big.Float).Quo(
			new(big.Float).SetInt(numerator),
			new(big.Float).SetInt(denominator))

	*/
}

/*
amountInWithFee = amountIn * 1000
numerator = amountInWithFee * reserveOut
denominator = (reserveIn*1000) + amountInWithFee
amountOut = numerator/denominator

amountOut = amountIn * 1000 * reserveOut / (reserveIn*1000) + amountInWithFee
amountOut = amountIn * 997 * reserveOut / (reserveIn*1000) + amountIn * 997

price = amountOut/amountIn

price = (amountIn * 1000 * reserveOut / (reserveIn*1000) + amountInWithFee) / amountIn
*/

func calculatePrice(reserveA, reserveB *big.Int) *big.Float {
	// feeless price = (reserveOut * 1000) / (reserveIn * 1000 + reserveOut)
	numerator := new(big.Int).Mul(reserveB, big.NewInt(1000))
	denominator := new(big.Int).Mul(reserveA, big.NewInt(1000))
	denominator = new(big.Int).Add(denominator, reserveB)

	return new(big.Float).Quo(
		new(big.Float).SetInt(numerator),
		new(big.Float).SetInt(denominator))
}
