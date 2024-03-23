package utils

import (
	"errors"
	"fmt"

	"github.com/ltcsuite/ltcd/btcjson"
	"github.com/ltcsuite/ltcd/chaincfg"
	"github.com/ltcsuite/ltcd/ltcutil"
	"github.com/ltcsuite/ltcd/rpcclient"
)

func Consolidate(client *rpcclient.Client, addr string) error {
	utxo, err := client.ListUnspent()
	if err != nil {
		return err
	}

	if len(utxo) < 2 {
		return nil
	}

	amount := 0.0
	inputs := make([]btcjson.TransactionInput, len(utxo))
	for i, tx := range utxo {
		inputs[i] = btcjson.TransactionInput{Txid: tx.TxID, Vout: tx.Vout}
		amount += tx.Amount
	}

	toAddr, err := ltcutil.DecodeAddress(addr, &chaincfg.MainNetParams)
	if err != nil {
		return err
	}

	outputs := map[ltcutil.Address]ltcutil.Amount{
		toAddr: ltcutil.Amount(ltcutil.SatoshiPerBitcoin*amount - 245),
	}

	tx, err := client.CreateRawTransaction(inputs, outputs, nil)
	if err != nil {
		return err
	}

	signedTx, ok, err := client.SignRawTransactionWithWallet(tx)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("error when signing tx")
	}

	hash, err := client.SendRawTransaction(signedTx, false)
	fmt.Println("Hash of the grouping tx :", hash)
	return err
}

func FilterUserAddrUTXO(addr string, utxo []btcjson.ListUnspentResult) []btcjson.ListUnspentResult {
	userUTXO := []btcjson.ListUnspentResult{}
	for _, tx := range utxo {
		if tx.Address == addr {
			userUTXO = append(userUTXO, tx)
		}
	}

	return userUTXO
}
