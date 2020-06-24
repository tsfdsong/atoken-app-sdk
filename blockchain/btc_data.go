package blockchain

//TransactionBTC btc transaction object
type TransactionBTC struct {
	TxID  string
	HexTx string
}

//BTCTxInput input of building BTC
type BTCTxInput struct {
	CoinType       string `json:"cointype"`
	Utxos          []Utxo `json:"utxos"`
	To             []WlTo `json:"to"`
	ChangeAddress  string `json:"changeaddress"`
	Fee            int64  `json:"fee"`
	BlockHash      string `json:"blockhash"`
	Dust           int64  `json:"dust"`
	OmniAddress    string `json:"omniAddress"`
	OmniCurrencyID int64  `json:"omniCurrencyID"`
	OmniAmount     int64  `json:"omniAmount"`
	NeedOmniOut    int    `json:"needOmniOut"`
}

//Utxo btc input
type Utxo struct {
	Address     string `json:"address"`
	TxID        string `json:"txid"`
	OutputIndex int    `json:"outputindex"`
	PkScript    string `json:"pkscript"` //last publickeyscript of utxo.vout
	Satoshis    int64  `json:"satoshis"`
	Public      string `json:"public"`
	Private     string `json:"private"`
}

//WlTo btc output
type WlTo struct {
	To       string `json:"to"`
	Satoshis int64  `json:"satoshis"`
}

func (input BTCTxInput) getChangeAmount() int64 {
	toAmount := int64(0)
	for _, txout := range input.To {
		toAmount += txout.Satoshis
	}

	fromAmount := int64(0)
	for _, txin := range input.Utxos {
		fromAmount += txin.Satoshis
	}
	//change amount

	valueNeed := toAmount + input.Fee

	if input.OmniCurrencyID != 0 {
		valueNeed -= input.Dust
	}

	changeAmount := fromAmount - valueNeed

	return changeAmount
}
