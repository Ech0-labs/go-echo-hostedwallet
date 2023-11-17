package tx

import (
	"errors"

	"github.com/ltcsuite/ltcd/btcjson"
	"github.com/ltcsuite/ltcd/chaincfg/chainhash"
	"github.com/ltcsuite/ltcd/ltcutil"
	"github.com/ltcsuite/ltcd/rpcclient"
)

type Inputs []btcjson.TransactionInput
type Outputs map[ltcutil.Address]ltcutil.Amount

func create(client *rpcclient.Client, inputs Inputs, outputs Outputs, data []byte) (*wire.MsgTx, error) {
	tx, err := client.CreateRawTransaction(inputs, outputs, nil)
	if err != nil {
		return &wire.MsgTx{}, err
	}

	opReturn, err := CreateOPReturn(data)
	if err != nil {
		return &wire.MsgTx{}, err
	}
	tx.AddTxOut(opReturn)

	return tx, nil
}

	signedTx, ok, err := client.SignRawTransactionWithWallet(tx)
	if err != nil {
		return &chainhash.Hash{}, err
	}

	if !ok {
		return &chainhash.Hash{}, errors.New("The tx cannot be fully signed with this wallet")
	}

	return client.SendRawTransaction(signedTx, false)
}
