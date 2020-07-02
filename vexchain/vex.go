package vex

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/tsfdsong/eos-go"
	"github.com/tsfdsong/eos-go/ecc"
	"github.com/tsfdsong/eos-go/system"
	"github.com/tsfdsong/eos-go/token"
)

func enc(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)

	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func getRawTxData(tx *eos.Transaction, chainID []byte, wifPriKey string) (string, error) {
	bytes, _ := json.Marshal(tx)
	fmt.Printf("get raw tx data: %v\n", string(bytes))

	//tx sig digest
	sigTx := eos.NewSignedTransaction(tx)

	txdata, cfd, err := sigTx.PackedTransactionAndCFD()
	if err != nil {
		return "", fmt.Errorf("packed transaction: %v", err)
	}

	sigDigest := eos.SigDigest(chainID, txdata, cfd)

	//
	privateKey, err := ecc.NewPrivateKey(wifPriKey)
	if err != nil {
		return "", fmt.Errorf("NewPrivateKey: %s", err)
	}
	str := hex.EncodeToString(sigDigest)
	fmt.Printf("hex: %v\n", str)

	sig, err := privateKey.Sign(sigDigest)
	if err != nil {
		return "", fmt.Errorf("signing through privatekey: %s", err)
	}
	sigTx.Signatures = append(sigTx.Signatures, sig)

	packedTx, err := sigTx.Pack(eos.CompressionNone)
	if err != nil {
		return "", fmt.Errorf("pack: %v", err)
	}

	txBytes, err := enc(packedTx)
	if err != nil {
		return "", fmt.Errorf("encode transaction: %s", err)
	}

	txres := string(txBytes)

	return txres, nil
}

func transferAmount(info *eos.InfoResp, in *TransferInfo, wifPriKey string) (string, error) {

	txOpts := &eos.TxOptions{
		ChainID:     info.ChainID,
		HeadBlockID: info.HeadBlockID,
	}

	quantity, err := eos.NewFixedSymbolAssetFromString(VEXSymbol, in.Quantity)
	if err != nil {
		return "", fmt.Errorf("transferAmount new asset %v", err)
	}

	from := eos.AccountName(in.From)
	to := eos.AccountName(in.To)
	memo := in.Memo

	act := &eos.Action{
		Account: token.AN("vex.token"),
		Name:    token.ActN("transfer"),
		Authorization: []eos.PermissionLevel{
			{Actor: from, Permission: token.PN("active")},
		},
		ActionData: eos.NewActionData(token.Transfer{
			From:     from,
			To:       to,
			Quantity: quantity,
			Memo:     memo,
		}),
	}

	tx := eos.NewTransaction([]*eos.Action{act}, txOpts)

	data, err := getRawTxData(tx, info.ChainID, wifPriKey)
	if err != nil {
		return "", fmt.Errorf("transferAmount %v", err)
	}

	return data, nil
}

func newVexAccount(in *NewAccountInfo) (*eos.Action, error) {
	ownerKeys := make([]eos.KeyWeight, 0)
	for _, onk := range in.Owner.Keys {
		pubKey, err := ecc.NewPublicKey(onk.PublicKey)
		if err != nil {
			fmt.Printf("NewPublicKey: %v\n", err)
			return nil, fmt.Errorf("new public key: %v", err)
		}

		tmp := eos.KeyWeight{
			PublicKey: pubKey,
			Weight:    onk.Weight,
		}

		ownerKeys = append(ownerKeys, tmp)
	}

	activeKeys := make([]eos.KeyWeight, 0)
	for _, onk := range in.Active.Keys {
		pubKey, err := ecc.NewPublicKey(onk.PublicKey)
		if err != nil {
			fmt.Printf("NewPublicKey: %v\n", err)
			return nil, fmt.Errorf("new public key: %v", err)
		}

		tmp := eos.KeyWeight{
			PublicKey: pubKey,
			Weight:    onk.Weight,
		}

		activeKeys = append(activeKeys, tmp)
	}

	ownerAccounts := make([]eos.PermissionLevelWeight, 0)
	for _, onp := range in.Owner.Accounts {
		tmp := eos.PermissionLevelWeight{
			Permission: eos.PermissionLevel{
				Actor:      eos.AccountName(onp.Permission.Actor),
				Permission: eos.PermissionName(onp.Permission.Permission),
			},
			Weight: onp.Weight,
		}

		ownerAccounts = append(ownerAccounts, tmp)
	}

	activeAccounts := make([]eos.PermissionLevelWeight, 0)
	for _, onp := range in.Active.Accounts {
		tmp := eos.PermissionLevelWeight{
			Permission: eos.PermissionLevel{
				Actor:      eos.AccountName(onp.Permission.Actor),
				Permission: eos.PermissionName(onp.Permission.Permission),
			},
			Weight: onp.Weight,
		}

		activeAccounts = append(activeAccounts, tmp)
	}

	return &eos.Action{
		Account: eos.AN("vexcore"),
		Name:    eos.ActN("newaccount"),
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AccountName(in.Creator), Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(system.NewAccount{
			Creator: eos.AccountName(in.Creator),
			Name:    eos.AccountName(in.Name),
			Owner: eos.Authority{
				Threshold: in.Owner.Threshold,
				Keys:      ownerKeys,
				Accounts:  ownerAccounts,
			},
			Active: eos.Authority{
				Threshold: in.Active.Threshold,
				Keys:      activeKeys,
				Accounts:  activeAccounts,
			},
		}),
	}, nil
}

func createAccount(info *eos.InfoResp, account *CreateAccountInfo, wifPriKey string) (string, error) {
	txOpts := &eos.TxOptions{
		ChainID:     info.ChainID,
		HeadBlockID: info.HeadBlockID,
	}
	actAccount, err := newVexAccount(&account.Account)
	if err != nil {
		return "", fmt.Errorf("newVexAccount: %v", err)
	}

	actBuyRAM := &eos.Action{
		Account: eos.AN("vexcore"),
		Name:    eos.ActN("buyrambytes"),
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AccountName(account.RAM.Payer), Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(system.BuyRAMBytes{
			Payer:    eos.AccountName(account.RAM.Payer),
			Receiver: eos.AccountName(account.RAM.Receiver),
			Bytes:    uint32(account.RAM.Bytes),
		}),
	}

	netAsset, err := eos.NewAssetFromString(account.BW.StakeNet)
	if err != nil {
		return "", fmt.Errorf("delegateBW net asset %v", err)
	}

	cpuAsset, err := eos.NewAssetFromString(account.BW.StakeCPU)
	if err != nil {
		return "", fmt.Errorf("delegateBW cpu asset %v", err)
	}

	isTransfer := false
	if account.BW.Transfer > 0 {
		isTransfer = true
	}

	actBW := &eos.Action{
		Account: eos.AN("vexcore"),
		Name:    eos.ActN("delegatebw"),
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AccountName(account.BW.From), Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(system.DelegateBW{
			From:     eos.AccountName(account.BW.From),
			Receiver: eos.AccountName(account.BW.Receiver),
			StakeNet: netAsset,
			StakeCPU: cpuAsset,
			Transfer: eos.Bool(isTransfer),
		}),
	}

	tx := eos.NewTransaction([]*eos.Action{actAccount, actBuyRAM, actBW}, txOpts)

	data, err := getRawTxData(tx, info.ChainID, wifPriKey)
	if err != nil {
		return "", fmt.Errorf("createAccount %v", err)
	}

	return data, nil
}

func sellRAM(info *eos.InfoResp, in *SellRAMInfo, wifPriKey string) (string, error) {
	txOpts := &eos.TxOptions{
		ChainID:     info.ChainID,
		HeadBlockID: info.HeadBlockID,
	}

	account := eos.AccountName(in.Account)
	bytes, err := strconv.ParseUint(in.Bytes, 10, 64)
	if err != nil {
		return "", fmt.Errorf("sellRAM parse bytes %v", err)
	}
	actSellRAM := &eos.Action{
		Account: eos.AN("vexcore"),
		Name:    eos.ActN("sellram"),
		Authorization: []eos.PermissionLevel{
			{Actor: account, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(system.SellRAM{
			Account: account,
			Bytes:   bytes,
		}),
	}

	tx := eos.NewTransaction([]*eos.Action{actSellRAM}, txOpts)

	data, err := getRawTxData(tx, info.ChainID, wifPriKey)
	if err != nil {
		return "", fmt.Errorf("sellRAM %v", err)
	}

	return data, nil
}

func buyRAM(info *eos.InfoResp, in *BuyRAMInfo, wifPriKey string) (string, error) {
	txOpts := &eos.TxOptions{
		ChainID:     info.ChainID,
		HeadBlockID: info.HeadBlockID,
	}

	ramAsset, err := eos.NewAssetFromString(in.Quant)
	if err != nil {
		return "", fmt.Errorf("buyRAM new asset %v", err)
	}

	actBuyRAM := &eos.Action{
		Account: eos.AN("vexcore"),
		Name:    eos.ActN("buyram"),
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AccountName(in.Payer), Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(system.BuyRAM{
			Payer:    eos.AccountName(in.Payer),
			Receiver: eos.AccountName(in.Receiver),
			Quantity: ramAsset,
		}),
	}

	tx := eos.NewTransaction([]*eos.Action{actBuyRAM}, txOpts)

	data, err := getRawTxData(tx, info.ChainID, wifPriKey)
	if err != nil {
		return "", fmt.Errorf("buyRAM %v", err)
	}

	return data, nil
}

func delegateBW(info *eos.InfoResp, bw *DelegateBWInfo, wifPriKey string) (string, error) {
	txOpts := &eos.TxOptions{
		ChainID:     info.ChainID,
		HeadBlockID: info.HeadBlockID,
	}

	netAsset, err := eos.NewAssetFromString(bw.StakeNet)
	if err != nil {
		return "", fmt.Errorf("delegateBW net asset %v", err)
	}

	cpuAsset, err := eos.NewAssetFromString(bw.StakeCPU)
	if err != nil {
		return "", fmt.Errorf("delegateBW cpu asset %v", err)
	}

	isTransfer := false
	if bw.Transfer > 0 {
		isTransfer = true
	}

	actBW := &eos.Action{
		Account: eos.AN("vexcore"),
		Name:    eos.ActN("delegatebw"),
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AccountName(bw.From), Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(system.DelegateBW{
			From:     eos.AccountName(bw.From),
			Receiver: eos.AccountName(bw.Receiver),
			StakeNet: netAsset,
			StakeCPU: cpuAsset,
			Transfer: eos.Bool(isTransfer),
		}),
	}

	tx := eos.NewTransaction([]*eos.Action{actBW}, txOpts)

	data, err := getRawTxData(tx, info.ChainID, wifPriKey)
	if err != nil {
		return "", fmt.Errorf("delegateBW %v", err)
	}

	return data, nil
}

func unDelegateBW(info *eos.InfoResp, bws []UnDelegateBWInfo, wifPriKey string) (string, error) {
	txOpts := &eos.TxOptions{
		ChainID:     info.ChainID,
		HeadBlockID: info.HeadBlockID,
	}

	actList := make([]*eos.Action, 0)
	for _, bw := range bws {
		netAsset, err := eos.NewAssetFromString(bw.UnstakeNet)
		if err != nil {
			return "", fmt.Errorf("unDelegateBW net asset %v", err)
		}

		cpuAsset, err := eos.NewAssetFromString(bw.UnstakeCPU)
		if err != nil {
			return "", fmt.Errorf("unDelegateBW cpu asset %v", err)
		}

		actBW := &eos.Action{
			Account: eos.AN("vexcore"),
			Name:    eos.ActN("undelegatebw"),
			Authorization: []eos.PermissionLevel{
				{Actor: eos.AccountName(bw.From), Permission: eos.PN("active")},
			},
			ActionData: eos.NewActionData(system.UndelegateBW{
				From:       eos.AccountName(bw.From),
				Receiver:   eos.AccountName(bw.Receiver),
				UnstakeNet: netAsset,
				UnstakeCPU: cpuAsset,
			}),
		}

		actList = append(actList, actBW)
	}

	tx := eos.NewTransaction(actList, txOpts)

	data, err := getRawTxData(tx, info.ChainID, wifPriKey)
	if err != nil {
		return "", fmt.Errorf("unDelegateBW %v", err)
	}

	return data, nil
}

func buyRAMBytes(info *eos.InfoResp, in *BuyRAMBytes, wifPriKey string) (string, error) {
	txOpts := &eos.TxOptions{
		ChainID:     info.ChainID,
		HeadBlockID: info.HeadBlockID,
	}

	actBuyRAM := &eos.Action{
		Account: eos.AN("vexcore"),
		Name:    eos.ActN("buyrambytes"),
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AccountName(in.Payer), Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(system.BuyRAMBytes{
			Payer:    eos.AccountName(in.Payer),
			Receiver: eos.AccountName(in.Receiver),
			Bytes:    uint32(in.Bytes),
		}),
	}

	tx := eos.NewTransaction([]*eos.Action{actBuyRAM}, txOpts)

	data, err := getRawTxData(tx, info.ChainID, wifPriKey)
	if err != nil {
		return "", fmt.Errorf("buyRAMBytes %v", err)
	}

	return data, nil
}
