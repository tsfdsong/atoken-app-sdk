package blockchain

import (
	"encoding/hex"
	"encoding/json"

	"github.com/btcsuite/btcd/btcec"
	"github.com/mr-tron/base58"
	"github.com/tsfdsong/atoken-app-sdk/hdwallet"

	"github.com/tyler-smith/go-bip39"
)

//AddresType ...
type AddresType struct {
	PrivateKey   string `json:"privatekey"`
	PublicKey    string `json:"publickey"`
	Address      string `json:"address"`
	AddressIndex int    `json:"addressindex"`
}

//WalletObject ...
type WalletObject struct {
	WalletID    string       `json:"walletid"`
	AddressList []AddresType `json:"addresslist"`
	Entropy     string       `json:"entropy"`
	Seed        string       `json:"seed"`
}

//Wallets list of WalletObject
type Wallets struct {
	WalletTable []*WalletObject `json:"table"`
}

//CreateWallets ...
func CreateWallets(mnemonic, coinType string, count int, isWIF bool) (string, error) {
	//0. get mnemonic
	mnemonic, err := hdwallet.CreateMnemonic(mnemonic)
	if err != nil {
		return "", err
	}

	result := make([]*WalletObject, 0)
	for i := 0; i < count; i++ {
		//1. common address
		wobj, err := getKeyPair(mnemonic, coinType, i, false, isWIF)
		if err != nil {
			return "", err
		}
		result = append(result, wobj)

		if coinType == "BTC" {
			//2. segwit address
			segobj, err := getKeyPair(mnemonic, coinType, i, true, isWIF)
			if err != nil {
				return "", err
			}
			result = append(result, segobj)
		}
	}

	tables := &Wallets{
		WalletTable: result,
	}

	strTables, err := json.Marshal(tables)
	if err != nil {
		return "", err
	}
	res := string(strTables)
	return res, nil
}

//CreateWallet ...
func CreateWallet(mnemonic, coinType string, isSegwit, isWIF bool) (string, error) {
	mnemonic, err := hdwallet.CreateMnemonic(mnemonic)
	if err != nil {
		return "", err
	}

	obj, err := getKeyPair(mnemonic, coinType, 0, isSegwit, isWIF)
	if err != nil {
		return "", err
	}

	strObj, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	res := string(strObj)
	return res, nil
}

//HexToWIF convert hex to wif
func HexToWIF(hexkey string) string {
	hexBytes, err := hex.DecodeString(hexkey)
	if err != nil {
		return ""
	}

	//0.add version 0x80
	versionPayload := append([]byte{byte(0x80)}, hexBytes...)

	checksumBytes := hdwallet.CheckSum(versionPayload)

	fullPayload := append(versionPayload, checksumBytes...)

	result := base58.Encode(fullPayload)

	return result
}

//getKeyPair ...
func getKeyPair(mnemonic, coinType string, addressIndex int, isSegwit, isWIF bool) (*WalletObject, error) {
	wallet, err := hdwallet.NewWallet(mnemonic, coinType)
	if err != nil {
		return nil, err
	}

	//get publickey and address
	publicKey, address, err := wallet.GetKeyAndAddress(coinType, addressIndex, isSegwit)
	if err != nil {
		return nil, err
	}

	//get private key
	var strPrivateKey string
	if isWIF {
		strPrivateKey, err = wallet.GetWIFPrivateKey(coinType, addressIndex, isSegwit)
	} else {
		strPrivateKey, err = wallet.GetPrivateKey(coinType, addressIndex, isSegwit)
	}

	if err != nil {
		return nil, err
	}

	//get wallet id
	walletID, err := wallet.GetWalletID()
	if err != nil {
		return nil, err
	}

	addrTypr := AddresType{
		PrivateKey:   strPrivateKey,
		PublicKey:    publicKey,
		Address:      address,
		AddressIndex: addressIndex,
	}

	addList := make([]AddresType, 0)
	addList = append(addList, addrTypr)

	return &WalletObject{
		WalletID:    walletID,
		AddressList: addList,
		Entropy:     wallet.Entropy,
		Seed:        wallet.Seed,
	}, nil
}

//ImportPrivateKey ...
func ImportPrivateKey(coinType, privateKey string, isSegwit, isWIF bool) (string, error) {
	//1. Recover private key from string
	var ecdsaPubKey *btcec.PublicKey
	var err error
	if isWIF {
		ecdsaPubKey, err = hdwallet.WIFToECDSAPublicKey(privateKey)
	} else {
		ecdsaPubKey, err = hdwallet.HexToECDSAPublicKey(privateKey)
	}

	if err != nil {
		return "", err
	}

	//2. Generate public key from private key
	publicKey, address, err := hdwallet.PublicKeyToAddress(coinType, ecdsaPubKey, isSegwit)
	if err != nil {
		return "", err
	}

	addrTypr := AddresType{
		PrivateKey:   privateKey,
		PublicKey:    publicKey,
		Address:      address,
		AddressIndex: 0,
	}

	addList := make([]AddresType, 0)
	addList = append(addList, addrTypr)

	obj := &WalletObject{
		WalletID:    "invalid",
		AddressList: addList,
		Entropy:     "",
		Seed:        "",
	}

	strObj, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	res := string(strObj)
	return res, nil
}

//EntropyFromMnemonic get entropy from mnemonic
func EntropyFromMnemonic(mnemonic string) (string, error) {
	enbyte, err := bip39.EntropyFromMnemonic(mnemonic)
	if err != nil {
		return "", err
	}

	hexEntropy := hex.EncodeToString(enbyte)

	return hexEntropy, nil
}

//MnemonicFromEntropy get mnemonic from entropy
func MnemonicFromEntropy(entropy string) (string, error) {
	byteEntropy, err := hex.DecodeString(entropy)
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(byteEntropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}
