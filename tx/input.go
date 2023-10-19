package tx

import "github.com/ltcsuite/ltcd/btcjson"

func CreateInputs(utxo ...btcjson.ListUnspentResult) ([]btcjson.TransactionInput, float64) {
	txInputs := make([]btcjson.TransactionInput, len(utxo))

	amount := 0.0
	for i, input := range utxo {
		amount += input.Amount
		txInputs[i] = btcjson.TransactionInput{Txid: input.TxID, Vout: input.Vout}
	}

	return txInputs, amount
}
