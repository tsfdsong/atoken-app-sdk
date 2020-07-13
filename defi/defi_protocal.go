package defi

import (
	"encoding/json"
	"fmt"

	"github.com/tsfdsong/atoken-app-sdk/defi/aave"
)

const (
	tDefiApprove int = iota
	tDefiAaveDeposit
	tDefiAaveRedeem
)

//DefiProtocal cmdType: operator code;data: json string of input parameter
func DefiProtocal(cmdType int, data string) (string, error) {
	switch cmdType {
	case tDefiApprove:
		{
			var input aave.AaveApprover
			err := json.Unmarshal([]byte(data), &input)
			if err != nil {
				return "", fmt.Errorf("Approve unmarshal, %v", err)
			}

			return aave.Approve(&input)
		}
	case tDefiAaveDeposit:
		{
			var input aave.AaveDepositer
			err := json.Unmarshal([]byte(data), &input)
			if err != nil {
				return "", fmt.Errorf("Deposit unmarshal, %v", err)
			}

			return aave.Deposit(&input)
		}
	case tDefiAaveRedeem:
		{
			var input aave.AaveRedeemer
			err := json.Unmarshal([]byte(data), &input)
			if err != nil {
				return "", fmt.Errorf("Redeem unmarshal, %v", err)
			}

			return aave.Redeem(&input)
		}
	default:
		return "", fmt.Errorf("DefiProtocal not support %v operator", cmdType)
	}
}
