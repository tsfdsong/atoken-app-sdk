package aave

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

//GetMethodID get id of method
func GetMethodID(methodName string) []byte {
	var methodID []byte
	switch methodName {
	case "approve":
		{
			hash := sha3.NewLegacyKeccak256()
			hash.Write([]byte("approve(address,uint256)"))
			methodID = hash.Sum(nil)[:4]
		}
	case "deposit":
		{
			hash := sha3.NewLegacyKeccak256()
			hash.Write([]byte("deposit(address,uint256,uint16)"))
			methodID = hash.Sum(nil)[:4]
		}
	case "redeem":
		{
			hash := sha3.NewLegacyKeccak256()
			hash.Write([]byte("redeem(uint256)"))
			methodID = hash.Sum(nil)[:4]
		}
	default:
		{
			methodID = []byte{0}
		}
	}

	return methodID
}

//MakeTransaction construct eth transaction
func MakeTransaction(privateKeyString, to string, value, nonce, gasLimit, gasPrice, chainID uint64, data []byte) (*types.Transaction, error) {
	//private key
	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		return nil, fmt.Errorf("HexToECDSA: %v", err)
	}

	tx := types.NewTransaction(nonce, common.HexToAddress(to), big.NewInt(int64(value)), gasLimit, big.NewInt(int64(gasPrice)), data)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(int64(chainID))), privateKey)
	if err != nil {
		return nil, fmt.Errorf("SignTx error: %v", err)
	}

	return signedTx, nil
}

//Approve Approve ERC20
func Approve(input *AaveApprover) (string, error) {
	//make data
	//spender
	toAddr := common.HexToAddress(input.LendingPoolCoreAddress)
	paddedToAddress := common.LeftPadBytes(toAddr.Bytes(), 32)

	//value
	total := big.NewInt(input.Amount)
	paddedValue := common.LeftPadBytes(total.Bytes(), 32)

	//tx data
	var data []byte
	methodID := GetMethodID("approve")
	data = append(data, methodID...)
	data = append(data, paddedToAddress...)
	data = append(data, paddedValue...)

	signedTx, err := MakeTransaction(input.HexPrivateKey, input.ToAddress, input.Value, input.Nonce, input.GasLimit, input.GasPrice, input.ChainID, data)
	if err != nil {
		return "", fmt.Errorf("Approve %v", err)
	}

	txBytes, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return "", fmt.Errorf("Approve rlp encode, %v", err)
	}

	res := hexutil.Encode(txBytes)

	fmt.Printf("Approve success: %v \n", res)

	return string(txBytes), nil
}

//Deposit deposit on Aave protocal
func Deposit(input *AaveDepositer) (string, error) {
	//spender
	paddedReserveAddress := common.LeftPadBytes(common.HexToAddress(input.ReserveAddress).Bytes(), 32)

	//value
	amount := big.NewInt(input.Amount)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	referCode := big.NewInt(input.ReferCode)
	paddedReferCode := common.LeftPadBytes(referCode.Bytes(), 32)

	//tx data
	var data []byte
	methodID := GetMethodID("deposit")
	data = append(data, methodID...)
	data = append(data, paddedReserveAddress...)
	data = append(data, paddedAmount...)
	data = append(data, paddedReferCode...)

	signedTx, err := MakeTransaction(input.HexPrivateKey, input.ToAddress, input.Value, input.Nonce, input.GasLimit, input.GasPrice, input.ChainID, data)
	if err != nil {
		return "", fmt.Errorf("Deposit %v", err)
	}

	txBytes, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return "", fmt.Errorf("Deposit rlp encode, %v", err)
	}

	res := hexutil.Encode(txBytes)

	fmt.Printf("Deposit success: %v \n", res)

	return res, nil
}

//Redeem redeem balance
func Redeem(input *AaveRedeemer) (string, error) {

	//value
	paddedAmount := common.LeftPadBytes(big.NewInt(input.Amount).Bytes(), 32)

	//tx data
	var data []byte
	methodID := GetMethodID("redeem")
	data = append(data, methodID...)
	data = append(data, paddedAmount...)

	signedTx, err := MakeTransaction(input.HexPrivateKey, input.ToAddress, input.Value, input.Nonce, input.GasLimit, input.GasPrice, input.ChainID, data)
	if err != nil {
		return "", fmt.Errorf("Redeem %v", err)
	}

	txBytes, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return "", fmt.Errorf("Redeem rlp encode, %v", err)
	}

	res := hexutil.Encode(txBytes)

	fmt.Printf("Redeem success: %v \n", res)

	return string(txBytes), nil
}
