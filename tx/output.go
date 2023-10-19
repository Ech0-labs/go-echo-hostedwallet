package tx

import (
	"github.com/ltcsuite/ltcd/chaincfg"
	"github.com/ltcsuite/ltcd/ltcutil"
	"github.com/ltcsuite/ltcd/txscript"
	"github.com/ltcsuite/ltcd/wire"
)

type TxOutput struct {
	adress ltcutil.Address
	ltcutil.Amount
}

func CreateOutput(addr string, amount float64) (*TxOutput, error) {
	ltcAddr, err := ltcutil.DecodeAddress(addr, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	return &TxOutput{ltcAddr, ltcutil.Amount(ltcutil.SatoshiPerBitcoin * amount)}, nil
}

func CreateOutputs(txOutputs ...*TxOutput) map[ltcutil.Address]ltcutil.Amount {
	outputs := make(map[ltcutil.Address]ltcutil.Amount, len(txOutputs))

	for _, txOut := range txOutputs {
		outputs[txOut.adress] = txOut.Amount
	}

	return outputs
}

func CreateOPReturn(data []byte) (*wire.TxOut, error) {
	nullData, err := txscript.NullDataScript(data)
	if err != nil {
		return nil, err
	}

	return wire.NewTxOut(0, nullData), nil
}
