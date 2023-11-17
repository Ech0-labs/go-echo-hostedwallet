package tx

import (
	"errors"
	"strconv"
	"strings"

	"github.com/ltcsuite/ltcd/btcjson"
	"github.com/ltcsuite/ltcd/chaincfg/chainhash"
	"github.com/ltcsuite/ltcd/ltcutil"
	"github.com/ltcsuite/ltcd/rpcclient"
	"github.com/ltcsuite/ltcd/wire"
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

func send(client *rpcclient.Client, tx *wire.MsgTx) (*chainhash.Hash, error) {
	signedTx, ok, err := client.SignRawTransactionWithWallet(tx)
	if err != nil {
		return &chainhash.Hash{}, err
	}

	if !ok {
		return &chainhash.Hash{}, errors.New("The tx cannot be fully signed with this wallet")
	}

	return client.SendRawTransaction(signedTx, false)
}

func MinRelayFees(client *rpcclient.Client, inputs Inputs, outputs Outputs, data []byte) (float64, error) {
	rawTx, err := create(client, inputs, outputs, data)

	if err != nil {
		return 0.0, err
	}

	_, err = send(client, rawTx)

	parts := strings.Split(err.Error(), " ")
	minRelayFees, err := strconv.ParseFloat(parts[len(parts)-1], 64)
	if err != nil {
		return 0.0, err
	}

	return (minRelayFees + 1) / ltcutil.SatoshiPerBitcoin, nil
}

func Send(client *rpcclient.Client, inputs Inputs, outputs Outputs, data []byte) (*chainhash.Hash, error) {
	tx, err := create(client, inputs, outputs, data)
	if err != nil {
		return nil, err
	}

	return send(client, tx)
}
