package flow

import (
	"os"

	"github.com/Ech0-labs/go-echo-prototype/tx"
	"github.com/Ech0-labs/go-echo-prototype/utils"
	"github.com/joho/godotenv"
	"github.com/ltcsuite/ltcd/chaincfg/chainhash"
	"github.com/ltcsuite/ltcd/rpcclient"
)

func Send(client *rpcclient.Client, message string) (*chainhash.Hash, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	userAddr := os.Getenv("USER_ADDR")
	echoAddr := os.Getenv("ECHO_ADDR")

	utxo, err := client.ListUnspent()
	if err != nil {
		return nil, err
	}

	userUTXO := utils.FilterUserAddrUTXO(userAddr, utxo)
	inputs, amount := tx.CreateInputs(userUTXO...)

	echoOutput, err := tx.CreateOutput(echoAddr, dust)
	if err != nil {
		return nil, err
	}

	restOutput, err := tx.CreateOutput(userAddr, amount-dust-fees)
	if err != nil {
		return nil, err
	}

	outputs := tx.CreateOutputs(echoOutput, restOutput)
	data := []byte(message)

	minRelayFees, err := tx.MinRelayFees(client, inputs, outputs, data)
	if err != nil {
		return nil, err
	}

	restOutput, err = tx.CreateOutput(userAddr, amount-dust-minRelayFees)
	if err != nil {
		return nil, err
	}
	outputs = tx.CreateOutputs(echoOutput, restOutput)

	return tx.Send(client, inputs, outputs, data)
}
