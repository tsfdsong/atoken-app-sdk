package defi

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/tsfdsong/atoken-app-sdk/defi/aave"
)

func TestDefiProtocal(t *testing.T) {
	//deposit
	depositData := aave.AaveDepositer{
		HexPrivateKey:  "6226d8f8c181622d82a84e2d36e4c66c49f07c36cbf001b9b78abeb8ba313d41",
		ToAddress:      "0x398ec7346dcd622edc5ae82352f02be94c62d119",
		Value:          0,
		Nonce:          7,
		GasPrice:       27000000000,
		GasLimit:       300000,
		ChainID:        1,
		ReserveAddress: "0xdAC17F958D2ee523a2206206994597C13D831ec7",
		Amount:         1000000,
		ReferCode:      0,
	}

	data, err := json.Marshal(depositData)
	if err != nil {
		fmt.Printf("marshal deposit data failed, %v", err)
		return
	}

	tx, err := DefiProtocal(1, string(data))
	if err != nil {
		fmt.Printf("DefiProtocal failed, %v", err)
		return
	}

	fmt.Printf("DefiProtocal tx data, %v\n", len(tx))

	txData, err := hexutil.Decode(tx)
	if err != nil {
		fmt.Printf("hexutil decode tranasction failed: %v", err)
		return
	}
	var txRaw types.Transaction
	err = rlp.DecodeBytes(txData, &txRaw)
	if err != nil {
		fmt.Printf("rlpdecode tranasction failed: %v", err)
		return
	}

	fmt.Printf("DefiProtocal tx, %v\n", txData)
}

func TestDefiProtocalOfRedeem(t *testing.T) {
	//deposit
	redeemData := aave.AaveRedeemer{
		HexPrivateKey: "6226d8f8c181622d82a84e2d36e4c66c49f07c36cbf001b9b78abeb8ba313d41",
		ToAddress:     "0x71fc860F7D3A592A4a98740e39dB31d25db65ae8",
		Value:         0,
		Nonce:         8,
		GasPrice:      38000000000,
		GasLimit:      300000,
		ChainID:       1,
		Amount:        1000000,
	}

	data, err := json.Marshal(redeemData)
	if err != nil {
		fmt.Printf("marshal deposit data failed, %v", err)
		return
	}

	tx, err := DefiProtocal(2, string(data))
	if err != nil {
		fmt.Printf("DefiProtocal redeem failed, %v", err)
		return
	}

	fmt.Printf("DefiProtocal tx data, %v\n", len(tx))

	// txData, err := hexutil.Decode(tx)
	// if err != nil {
	// 	fmt.Printf("hexutil decode tranasction failed: %v", err)
	// 	return
	// }
	// var txRaw types.Transaction
	// err = rlp.DecodeBytes(txData, &txRaw)
	// if err != nil {
	// 	fmt.Printf("rlpdecode tranasction failed: %v", err)
	// 	return
	// }

	// fmt.Printf("DefiProtocal tx, %v\n", txData)
}
