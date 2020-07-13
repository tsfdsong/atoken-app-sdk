package vex

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/tsfdsong/eos-go"
)

func TestVexTransfer(t *testing.T) {

	input := &TransferInfo{
		From:     "atokentry123",
		To:       "atokenmai123",
		Quantity: "1.0000 VEX",
		Memo:     "test",
	}

	fmt.Printf("input: %v\n", input)

	api := eos.New("https://explorer.vexanium.com:6960/")
	info, err := api.GetInfo(context.Background())
	if err != nil {
		t.Errorf("GetInfo: %v\n", err)
		return
	}

	tx, err := transferAmount(info, input, "5JknoozotmRa18kRNcfYdNVHXPHTTf9pwzUaiVtDCBPpM2J73hE", "active")
	if err != nil {
		t.Errorf("transferAmount: %v\n", err)
		return
	}

	var packedTx eos.PackedTransaction
	json.Unmarshal([]byte(tx), &packedTx)

	//brodcast
	response, err := api.PushTransaction(context.Background(), &packedTx)
	if err != nil {
		t.Errorf("PushTransaction: %v\n", err)
		return
	}

	fmt.Printf("Transaction [%s] submitted to the network succesfully.\n", hex.EncodeToString(response.Processed.ID))

	fmt.Printf("tx : %v\n", tx)
}

func TestVexNewAccount(t *testing.T) {
	ownerKeyWeight := KeyWeightInfo{
		PublicKey: "VEX5awQ9H9Tj2qorH1ETLDDynvdQgxxFZepfr1iS5EU3WtWvzyfe1",
		Weight:    1,
	}
	ownerAuthor := AuthorityInfo{
		Threshold: 1,
		Keys:      []KeyWeightInfo{ownerKeyWeight},
	}

	activeKeyWeight := KeyWeightInfo{
		PublicKey: "VEX5awQ9H9Tj2qorH1ETLDDynvdQgxxFZepfr1iS5EU3WtWvzyfe1",
		Weight:    1,
	}
	activeAuthor := AuthorityInfo{
		Threshold: 1,
		Keys:      []KeyWeightInfo{activeKeyWeight},
	}

	newAccount := NewAccountInfo{
		Creator: "atokentry123",
		Name:    "atokenmax345",
		Owner:   ownerAuthor,
		Active:  activeAuthor,
	}

	buyRAMBytes := BuyRAMBytes{
		Payer:    "atokentry123",
		Receiver: "atokenmax345",
		Bytes:    4096,
	}

	deleBW := DelegateBWInfo{
		From:     "atokentry123",
		Receiver: "atokenmax345",
		StakeNet: "0.2500 VEX",
		StakeCPU: "0.7500 VEX",
		Transfer: 1,
	}

	input := &CreateAccountInfo{
		Account: newAccount,
		RAM:     buyRAMBytes,
		BW:      deleBW,
	}

	api := eos.New("https://explorer.vexanium.com:6960/")
	info, err := api.GetInfo(context.Background())
	if err != nil {
		t.Errorf("GetInfo: %v\n", err)
		return
	}

	tx, err := createAccount(info, input, "5JknoozotmRa18kRNcfYdNVHXPHTTf9pwzUaiVtDCBPpM2J73hE", "active")
	if err != nil {
		t.Errorf("createAccount: %v\n", err)
		return
	}

	fmt.Printf("tx: %v\n", tx)

	var packedTx eos.PackedTransaction
	err = json.Unmarshal([]byte(tx), &packedTx)
	if err != nil {
		t.Errorf("Unmarshal: %v\n", err)
		return
	}

	//brodcast
	response, err := api.PushTransaction(context.Background(), &packedTx)
	if err != nil {
		t.Errorf("PushTransaction: %v\n", err)
		return
	}

	fmt.Printf("Transaction [%s] submitted to the network succesfully.\n", hex.EncodeToString(response.Processed.ID))

	fmt.Printf("tx : %v\n", tx)
}

func TestVexSellRAM(t *testing.T) {

	input := &SellRAMInfo{
		Account: "atokentry123",
		Bytes:   "3000000",
	}

	fmt.Printf("input: %v\n", input)

	api := eos.New("https://explorer.vexanium.com:6960/")
	info, err := api.GetInfo(context.Background())
	if err != nil {
		t.Errorf("GetInfo: %v\n", err)
		return
	}

	byda, _ := json.Marshal(&info)
	fmt.Printf("%v\n", string(byda))

	tx, err := sellRAM(info, input, "5JknoozotmRa18kRNcfYdNVHXPHTTf9pwzUaiVtDCBPpM2J73hE", "active")
	if err != nil {
		t.Errorf("sellRAM: %v\n", err)
		return
	}

	var packedTx eos.PackedTransaction
	json.Unmarshal([]byte(tx), &packedTx)

	//brodcast
	response, err := api.PushTransaction(context.Background(), &packedTx)
	if err != nil {
		t.Errorf("sellRAM PushTransaction: %v\n", err)
		return
	}

	fmt.Printf("sellRAM tx : %v\n", tx)

	fmt.Printf("sellRAM Transaction [%s] submitted to the network succesfully.\n", hex.EncodeToString(response.Processed.ID))
}

func TestVexBuyRAM(t *testing.T) {

	input := &BuyRAMInfo{
		Payer:    "atokentry123",
		Receiver: "atokentry123",
		Quant:    "0.5000 VEX",
	}

	fmt.Printf("input: %v\n", input)

	api := eos.New("https://explorer.vexanium.com:6960/")
	info, err := api.GetInfo(context.Background())
	if err != nil {
		t.Errorf("GetInfo: %v\n", err)
		return
	}

	byda, _ := json.Marshal(&info)
	fmt.Printf("%v\n", string(byda))

	tx, err := buyRAM(info, input, "5JknoozotmRa18kRNcfYdNVHXPHTTf9pwzUaiVtDCBPpM2J73hE", "active")
	if err != nil {
		t.Errorf("sellRAM: %v\n", err)
		return
	}

	var packedTx eos.PackedTransaction
	json.Unmarshal([]byte(tx), &packedTx)

	//brodcast
	response, err := api.PushTransaction(context.Background(), &packedTx)
	if err != nil {
		t.Errorf("sellRAM PushTransaction: %v\n", err)
		return
	}

	fmt.Printf("sellRAM tx : %v\n", tx)

	fmt.Printf("sellRAM Transaction [%s] submitted to the network succesfully.\n", hex.EncodeToString(response.Processed.ID))
}

func TestVexDelegateBWNoTransfer(t *testing.T) {

	input := &DelegateBWInfo{
		From:     "atokentry123",
		Receiver: "atokenmai123",
		StakeNet: "1.0000 VEX",
		StakeCPU: "1.0000 VEX",
		Transfer: 0,
	}

	fmt.Printf("input: %v\n", input)

	api := eos.New("https://explorer.vexanium.com:6960/")
	info, err := api.GetInfo(context.Background())
	if err != nil {
		t.Errorf("GetInfo: %v\n", err)
		return
	}

	byda, _ := json.Marshal(&info)
	fmt.Printf("%v\n", string(byda))

	tx, err := delegateBW(info, input, "5JknoozotmRa18kRNcfYdNVHXPHTTf9pwzUaiVtDCBPpM2J73hE", "active")
	if err != nil {
		t.Errorf("delegateBW: %v\n", err)
		return
	}

	var packedTx eos.PackedTransaction
	json.Unmarshal([]byte(tx), &packedTx)

	//brodcast
	response, err := api.PushTransaction(context.Background(), &packedTx)
	if err != nil {
		t.Errorf("delegateBW PushTransaction: %v\n", err)
		return
	}

	fmt.Printf("delegateBW tx : %v\n", tx)

	fmt.Printf("delegateBW Transaction [%s] submitted to the network succesfully.\n", hex.EncodeToString(response.Processed.ID))
}

func TestVexUnDelegateBWNoTransfer(t *testing.T) {
	oneBW := UnDelegateBWInfo{
		From:       "atokentry123",
		Receiver:   "atokenmai123",
		UnstakeNet: "0.1500 VEX",
		UnstakeCPU: "0.1500 VEX",
	}

	input := []UnDelegateBWInfo{oneBW}

	fmt.Printf("input: %v\n", input)

	api := eos.New("https://explorer.vexanium.com:6960/")
	info, err := api.GetInfo(context.Background())
	if err != nil {
		t.Errorf("GetInfo: %v\n", err)
		return
	}

	byda, _ := json.Marshal(&info)
	fmt.Printf("%v\n", string(byda))

	tx, err := unDelegateBW(info, input, "5JknoozotmRa18kRNcfYdNVHXPHTTf9pwzUaiVtDCBPpM2J73hE", "active")
	if err != nil {
		t.Errorf("unDelegateBW: %v\n", err)
		return
	}

	var packedTx eos.PackedTransaction
	json.Unmarshal([]byte(tx), &packedTx)

	//brodcast
	response, err := api.PushTransaction(context.Background(), &packedTx)
	if err != nil {
		t.Errorf("unDelegateBW PushTransaction: %v\n", err)
		return
	}

	fmt.Printf("unDelegateBW tx : %v\n", tx)

	fmt.Printf("unDelegateBW Transaction [%s] submitted to the network succesfully.\n", hex.EncodeToString(response.Processed.ID))
}

func TestVexDelegateBWTransfer(t *testing.T) {

	input := &DelegateBWInfo{
		From:     "atokentry123",
		Receiver: "atokenmai345",
		StakeNet: "0.1000 VEX",
		StakeCPU: "0.0000 VEX",
		Transfer: 1,
	}

	fmt.Printf("input: %v\n", input)

	api := eos.New("https://explorer.vexanium.com:6960/")
	info, err := api.GetInfo(context.Background())
	if err != nil {
		t.Errorf("GetInfo: %v\n", err)
		return
	}

	byda, _ := json.Marshal(&info)
	fmt.Printf("%v\n", string(byda))

	tx, err := delegateBW(info, input, "5JknoozotmRa18kRNcfYdNVHXPHTTf9pwzUaiVtDCBPpM2J73hE", "active")
	if err != nil {
		t.Errorf("delegateBW: %v\n", err)
		return
	}

	var packedTx eos.PackedTransaction
	json.Unmarshal([]byte(tx), &packedTx)

	//brodcast
	response, err := api.PushTransaction(context.Background(), &packedTx)
	if err != nil {
		t.Errorf("delegateBW PushTransaction: %v\n", err)
		return
	}

	fmt.Printf("delegateBW tx : %v\n", tx)

	fmt.Printf("delegateBW Transaction [%s] submitted to the network succesfully.\n", hex.EncodeToString(response.Processed.ID))
}

func TestVexUnDelegateBWTransfer(t *testing.T) {
	oneBW := UnDelegateBWInfo{
		From:       "atokenmai345",
		Receiver:   "atokenmai345",
		UnstakeNet: "0.0800 VEX",
		UnstakeCPU: "0.4800 VEX",
	}

	input := []UnDelegateBWInfo{oneBW}

	fmt.Printf("input: %v\n", input)

	api := eos.New("https://explorer.vexanium.com:6960/")
	info, err := api.GetInfo(context.Background())
	if err != nil {
		t.Errorf("GetInfo: %v\n", err)
		return
	}

	byda, _ := json.Marshal(&info)
	fmt.Printf("%v\n", string(byda))

	tx, err := unDelegateBW(info, input, "5JLYhubP9whvSjUsPMGtm9bkHCPEpFa9QNFKfSdxU9kpewApixj", "active")
	if err != nil {
		t.Errorf("unDelegateBW: %v\n", err)
		return
	}

	var packedTx eos.PackedTransaction
	json.Unmarshal([]byte(tx), &packedTx)

	//brodcast
	response, err := api.PushTransaction(context.Background(), &packedTx)
	if err != nil {
		t.Errorf("unDelegateBW PushTransaction: %v\n", err)
		return
	}

	fmt.Printf("unDelegateBW tx : %v\n", tx)

	fmt.Printf("unDelegateBW Transaction [%s] submitted to the network succesfully.\n", hex.EncodeToString(response.Processed.ID))
}
