package blockchain

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/tsfdsong/hdwallet"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

const (
	// CurrentTxInSequenceNum is MaxTxInSequenceNum -2
	CurrentTxInSequenceNum uint32 = 0xfffffffd

	MinDustOutput int64 = 546
)

//GetAddressFromPrivKey get address from privatekey
func getAddressFromPrivKey(prikey *btcec.PrivateKey, isWitSeg bool) string {
	priKey := btcec.PrivateKey(*prikey)
	esdsaPubKey := priKey.PubKey()
	pubkeyBytes := esdsaPubKey.SerializeCompressed()

	var address string
	if isWitSeg {
		address = hdwallet.ToBTC(pubkeyBytes, true)
	} else {
		address = hdwallet.ToBTC(pubkeyBytes, false)
	}

	return address
}

//GetPayToAddrScript add script
func getPayToAddrScript(address string) []byte {
	rcvAddress, _ := btcutil.DecodeAddress(address, &chaincfg.MainNetParams)
	rcvScript, _ := txscript.PayToAddrScript(rcvAddress)
	return rcvScript
}

func getTxOut(address string, amount int64) *wire.TxOut {
	// create TxOut
	rcvscript := getPayToAddrScript(address)

	txOut := wire.NewTxOut(amount, rcvscript)
	return txOut
}

func createOmniData(currencyID, amount int64) string {
	omniPrefix := "6f6d6e69" //omni ascii
	txVersion := "0000"      //transaction version
	txType := "0000"         //transaction type: simple send

	curID := fmt.Sprintf("%08x", currencyID) //currency id
	amt := fmt.Sprintf("%016x", amount)      //amount

	result := omniPrefix + txVersion + txType + curID + amt

	return result
}

func getOmniTxOut(currencyID, amount int64) *wire.TxOut {
	omniData := createOmniData(currencyID, amount)

	data, _ := hex.DecodeString(omniData)

	b := txscript.NewScriptBuilder()
	b.AddOp(txscript.OP_RETURN)
	b.AddData([]byte(data))

	sigscript, _ := b.Script()

	// fmt.Printf("opreturn: %v\n", hex.EncodeToString(sigscript))

	txOut := wire.NewTxOut(0, sigscript)
	return txOut
}

func txToHex(tx *wire.MsgTx) string {
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	tx.Serialize(buf)
	return hex.EncodeToString(buf.Bytes())
}

//IsWitSehAddress check address type
func isWitSehAddress(addr string) bool {
	rcvAddress, _ := btcutil.DecodeAddress(addr, &chaincfg.MainNetParams)

	switch rcvAddress.(type) {
	case *btcutil.AddressWitnessPubKeyHash:
		return true
	case *btcutil.AddressWitnessScriptHash:
		return true
	case *btcutil.AddressScriptHash:
		return true
	default:
		return false
	}
}

//buildBTCTx construct btc transaction
func buildBTCTx(input BTCTxInput) (*wire.MsgTx, error) {

	//0. create new empty transaction
	redemTx := wire.NewMsgTx(wire.TxVersion)

	//1. calculate change of btc
	changeAmount := input.getChangeAmount()
	if changeAmount > 0 {
		if changeAmount > MinDustOutput {
			//change
			redemTx.AddTxOut(getTxOut(input.ChangeAddress, changeAmount))
		} else {
			//skip change ,make it miner fee
		}

	}

	//2.construct vout
	if input.OmniCurrencyID != 0 {
		if input.NeedOmniOut == 1 {
			redemTx.AddTxOut(getTxOut(input.ChangeAddress, input.Dust))
		}

		//add omni txout
		redemTx.AddTxOut(getOmniTxOut(input.OmniCurrencyID, input.OmniAmount))

		for _, v := range input.To {
			redemTx.AddTxOut(getTxOut(v.To, input.Dust))
		}
	} else {
		for _, v := range input.To {
			redemTx.AddTxOut(getTxOut(v.To, v.Satoshis))
		}
	}

	//3. vins
	for _, txin := range input.Utxos {
		hash, err := chainhash.NewHashFromStr(txin.TxID)
		if err != nil {
			return nil, fmt.Errorf("could not get hash from transaction ID: %v", err)
		}

		// create TxIn
		outPoint := wire.NewOutPoint(hash, uint32(txin.OutputIndex))
		txIn := wire.NewTxIn(outPoint, nil, nil)
		txIn.Sequence = CurrentTxInSequenceNum
		redemTx.AddTxIn(txIn)
	}

	//filled tx.vin.scriptsig
	for i := range input.Utxos {
		// sign transaction
		pkScript, err := hex.DecodeString(input.Utxos[i].PkScript)
		if err != nil {
			return nil, fmt.Errorf("could not get pkscript: %v", err)
		}

		myPrivateKey, err := hdwallet.HexToECDSAPrivateKey(input.Utxos[i].Private)
		if err != nil {
			return nil, err
		}

		if isWitSehAddress(input.Utxos[i].Address) {
			txSigHashes := txscript.NewTxSigHashes(redemTx)

			witnessTx, err := txscript.WitnessSignature(
				redemTx, // The tx to be signed.
				txSigHashes,
				i, // The index of the txin the signature is for.
				input.Utxos[i].Satoshis,
				pkScript,            // The other half of the script from the PubKeyHash.
				txscript.SigHashAll, // The signature flags that indicate what the sig covers.
				myPrivateKey,        // The key to generate the signature with.
				true)                // The compress sig flag. This saves space on the blockchain.

			if err != nil {
				return nil, fmt.Errorf("could not generate signature: %v", err)
			}

			redemTx.TxIn[i].Witness = witnessTx

			for _, wit := range redemTx.TxIn[i].Witness {
				fmt.Printf("wit: %v\n", hex.EncodeToString(wit))
			}

			//scriptSig
			pk := (*btcec.PublicKey)(&myPrivateKey.PublicKey)

			pkData := pk.SerializeCompressed()

			address, err := btcutil.NewAddressWitnessPubKeyHash(
				btcutil.Hash160(pkData), &chaincfg.MainNetParams)
			if err != nil {
				return nil, err
			}

			scriptsig, err := txscript.PayToAddrScript(address)
			if err != nil {
				return nil, err
			}

			buf := bytes.NewBuffer(make([]byte, 0, len(scriptsig)+2))
			buf.WriteByte(byte(len(scriptsig)))
			buf.Write(scriptsig)

			// sigStr := hex.EncodeToString(buf.Bytes())
			// fmt.Printf("i => sig:   %v %v\n", i, sigStr)

			redemTx.TxIn[i].SignatureScript = buf.Bytes()

		} else {
			scriptsig, err := txscript.SignatureScript(
				redemTx,             // The tx to be signed.
				i,                   // The index of the txin the signature is for.
				pkScript,            // The other half of the script from the PubKeyHash.
				txscript.SigHashAll, // The signature flags that indicate what the sig covers.
				myPrivateKey,        // The key to generate the signature with.
				true)                // The compress sig flag. This saves space on the blockchain.

			if err != nil {
				return nil, fmt.Errorf("could not generate signature: %v", err)
			}

			// sigStr := hex.EncodeToString(scriptsig)
			// fmt.Printf("i => sig:   %v %v\n", i, sigStr)

			redemTx.TxIn[i].SignatureScript = scriptsig

			//Validate signature
			vm, err := txscript.NewEngine(pkScript, redemTx, i, txscript.StandardVerifyFlags, nil, nil, input.Utxos[i].Satoshis)
			if err != nil {
				return nil, fmt.Errorf("validate signature: %v", err)
			}

			if err := vm.Execute(); err != nil {
				return nil, fmt.Errorf("vm.Execute: %v", err)
			}
		}

	}

	return redemTx, nil
}

//TransferBTC make btc transaction
func TransferBTC(btc string) (*TransactionBTC, error) {
	var input BTCTxInput
	err := json.Unmarshal([]byte(btc), &input)
	if err != nil {
		return nil, err
	}

	tx, err := buildBTCTx(input)
	if err != nil {
		return nil, err
	}

	hexTx := txToHex(tx)
	txid := tx.TxHash().String()

	return &TransactionBTC{
		HexTx: hexTx,
		TxID:  txid,
	}, nil
}
