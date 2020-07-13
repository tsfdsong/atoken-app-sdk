package aave

//AaveApprover approve input
type AaveApprover struct {
	HexPrivateKey string `json:"HexPrivateKey"`
	ToAddress     string `json:"ToAddress"`
	Value         uint64 `json:"Value"`
	Nonce         uint64 `json:"Nonce"`
	GasPrice      uint64 `json:"GasPrice"`
	GasLimit      uint64 `json:"GasLimit"`
	ChainID       uint64 `json:"ChainID"`
	//
	LendingPoolCoreAddress string `json:"LendingPoolCoreAddress"`
	Amount                 int64  `json:"Amount"`
}

//AaveDepositer deposit input
type AaveDepositer struct {
	HexPrivateKey string `json:"HexPrivateKey"`
	ToAddress     string `json:"ToAddress"`
	Value         uint64 `json:"Value"`
	Nonce         uint64 `json:"Nonce"`
	GasPrice      uint64 `json:"GasPrice"`
	GasLimit      uint64 `json:"GasLimit"`
	ChainID       uint64 `json:"ChainID"`

	//
	ReserveAddress string `json:"ReserveAddress"`
	Amount         int64  `json:"Amount"`
	ReferCode      int64  `json:"ReferCode"`
}

//AaveRedeemer redeem input
type AaveRedeemer struct {
	HexPrivateKey string `json:"HexPrivateKey"`
	ToAddress     string `json:"ToAddress"`
	Value         uint64 `json:"Value"`
	Nonce         uint64 `json:"Nonce"`
	GasPrice      uint64 `json:"GasPrice"`
	GasLimit      uint64 `json:"GasLimit"`
	ChainID       uint64 `json:"ChainID"`

	//
	Amount int64 `json:"Amount"`
}
