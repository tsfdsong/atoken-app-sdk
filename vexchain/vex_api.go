package vex

import (
	"encoding/json"
	"fmt"

	"github.com/tsfdsong/eos-go"
)

const (
	//赎回
	tVEXTransferTypeUndelegatebw int = iota
	//抵押
	tVEXTransferTypeDelegatebw
	//买内存
	tVEXTransferTypeBuyRam
	//卖内存
	tVEXTransferTypeSellRam
	//创建账号
	tVEXTransferTypeCreateAccount
	//转账交易
	tVEXTransferTypeTransferAmount
	// 购买内存bytes
	tVEXTransferTypeBuyRamBytes
)

//VexAPI common api
func VexAPI(cmdType int, infoStr, data, wifPriKey, perm string) (string, error) {
	var info eos.InfoResp
	err := json.Unmarshal([]byte(infoStr), &info)
	if err != nil {
		return "", fmt.Errorf("unmarshal info: %v", err)
	}

	switch cmdType {
	case tVEXTransferTypeUndelegatebw:
		{
			bw := make([]UnDelegateBWInfo, 0)
			err := json.Unmarshal([]byte(data), &bw)
			if err != nil {
				return "", fmt.Errorf("unmarshal UnDelegateBWInfo: %v", err)
			}

			return unDelegateBW(&info, bw, wifPriKey, perm)
		}
	case tVEXTransferTypeDelegatebw:
		{
			var bw DelegateBWInfo
			err := json.Unmarshal([]byte(data), &bw)
			if err != nil {
				return "", fmt.Errorf("unmarshal DelegateBWInfo: %v", err)
			}

			return delegateBW(&info, &bw, wifPriKey, perm)
		}
	case tVEXTransferTypeBuyRam:
		{
			var ram BuyRAMInfo
			err := json.Unmarshal([]byte(data), &ram)
			if err != nil {
				return "", fmt.Errorf("unmarshal BuyRAMInfo: %v", err)
			}

			return buyRAM(&info, &ram, wifPriKey, perm)
		}
	case tVEXTransferTypeSellRam:
		{
			var ram SellRAMInfo
			err := json.Unmarshal([]byte(data), &ram)
			if err != nil {
				return "", fmt.Errorf("unmarshal SellRAMInfo: %v", err)
			}

			return sellRAM(&info, &ram, wifPriKey, perm)
		}
	case tVEXTransferTypeCreateAccount:
		{
			var acct CreateAccountInfo
			err := json.Unmarshal([]byte(data), &acct)
			if err != nil {
				return "", fmt.Errorf("unmarshal CreateAccountInfo: %v", err)
			}

			return createAccount(&info, &acct, wifPriKey, perm)
		}
	case tVEXTransferTypeTransferAmount:
		{
			var trans TransferInfo
			err := json.Unmarshal([]byte(data), &trans)
			if err != nil {
				return "", fmt.Errorf("unmarshal TransferInfo: %v", err)
			}

			return transferAmount(&info, &trans, wifPriKey, perm)
		}
	case tVEXTransferTypeBuyRamBytes:
		{
			var ram BuyRAMBytes
			err := json.Unmarshal([]byte(data), &ram)
			if err != nil {
				return "", fmt.Errorf("unmarshal BuyRAMBytes: %v", err)
			}

			return buyRAMBytes(&info, &ram, wifPriKey, perm)
		}
	}

	return "", fmt.Errorf("unsupport operate type: %v", cmdType)
}
