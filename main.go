package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Ech0-labs/go-echo-prototype/tx"
	"github.com/Ech0-labs/go-echo-prototype/utils"
	"github.com/joho/godotenv"
	"github.com/ltcsuite/ltcd/ltcutil"
	"github.com/ltcsuite/ltcd/rpcclient"
)

const Message = "Elon musk, get ready for ratio"
const userAddr = "ltc1qxwlk7kh2vx8xgvff0rxgrtzh2v00lzp7v9updr"
const echoAddr = "ltc1qwewx5q6fwy6d0acuu3mk29cema4hkhew5afmsa"
const dust = 2940 / ltcutil.SatoshiPerBitcoin
const fees = 182 / ltcutil.SatoshiPerBitcoin

func Handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	err := godotenv.Load()
	Handle(err)

	host := os.Getenv("HOST")
	user := os.Getenv("USER")
	pass := os.Getenv("PASS")

	var conf = &rpcclient.ConnConfig{
		Host:         host,
		User:         user,
		Pass:         pass,
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	client, err := rpcclient.New(conf, nil)
	Handle(err)
	defer client.Shutdown()

	blockCount, err := client.GetBlockCount()
	Handle(err)

	fmt.Println("Block height :", blockCount)

	if true {
		//return
	}

	utxo, err := client.ListUnspent()
	userUTXO := utils.FilterUserAddrUTXO(userAddr, utxo)

	inputs, amount := tx.CreateInputs(userUTXO...)

	echoOutput, err := tx.CreateOutput(echoAddr, dust)
	Handle(err)

	restOutput, err := tx.CreateOutput(userAddr, amount-dust-fees)
	Handle(err)

	outputs := tx.CreateOutputs(echoOutput, restOutput)

	message := []byte(Message)

	hash, err := tx.Send(*client, inputs, outputs, message)
	Handle(err)

	fmt.Println("hash of the tx :", hash)
}
