package vex

import (
	"github.com/tsfdsong/eos-go"
	"github.com/tsfdsong/eos-go/ecc"
)

// VEXSymbol represents the standard VEX symbol on the chain.  It's
// here just to speed up things.
var VEXSymbol = eos.Symbol{Precision: 4, Symbol: "VEX"}

//TransferInfo input parameter of transfer
type TransferInfo struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Quantity string `json:"quantity"`
	Memo     string `json:"memo"`
}

//CreateAccountInfo create account info
type CreateAccountInfo struct {
	Account NewAccountInfo `json:"newaccount"`
	RAM     BuyRAMBytes    `json:"buyrambytes"`
	BW      DelegateBWInfo `json:"delegatebw"`
}

//WaitWeightInfo info
type WaitWeightInfo struct {
	WaitSec uint32 `json:"wait_sec"`
	Weight  uint16 `json:"weight"` // weight_type
}

//KeyWeightInfo KeyWeight info
type KeyWeightInfo struct {
	PublicKey string `json:"key"`
	Weight    uint16 `json:"weight"` // weight_type
}

//KeyWeight key
type KeyWeight struct {
	PublicKey ecc.PublicKey `json:"key"`
	Weight    uint16        `json:"weight"` // weight_type
}

//AuthorityInfo authority info
type AuthorityInfo struct {
	Threshold uint32                      `json:"threshold"`
	Keys      []KeyWeightInfo             `json:"keys"`
	Accounts  []eos.PermissionLevelWeight `json:"accounts"`
	Waits     []eos.WaitWeight            `json:"waits"`
}

//NewAccountInfo new account
type NewAccountInfo struct {
	Creator string        `json:"creator"`
	Name    string        `json:"name"`
	Owner   AuthorityInfo `json:"owner"`
	Active  AuthorityInfo `json:"active"`
}

//BuyRAMInfo buy ram
type BuyRAMInfo struct {
	Payer    string `json:"payer"`
	Receiver string `json:"receiver"`
	Quant    string `json:"quant"`
}

//SellRAMInfo sell ram
type SellRAMInfo struct {
	Account string `json:"account"`
	Bytes   string `json:"bytes"`
}

//DelegateBWInfo delegate BW
type DelegateBWInfo struct {
	From     string `json:"from"`
	Receiver string `json:"receiver"`
	StakeNet string `json:"stake_net"`
	StakeCPU string `json:"stake_cpu"`
	Transfer int    `json:"transfer"`
}

//UnDelegateBWInfo undelegate BW
type UnDelegateBWInfo struct {
	From       string `json:"from"`
	Receiver   string `json:"receiver"`
	UnstakeNet string `json:"unstake_net_quantity"`
	UnstakeCPU string `json:"unstake_cpu_quantity"`
}

//BuyRAMBytes buy ram bytes
type BuyRAMBytes struct {
	Payer    string `json:"payer"`
	Receiver string `json:"receiver"`
	Bytes    uint64 `json:"bytes"`
}
