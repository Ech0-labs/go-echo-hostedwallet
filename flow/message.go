package flow

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/ltcsuite/ltcd/btcjson"
	"github.com/ltcsuite/ltcd/chaincfg/chainhash"
	"github.com/ltcsuite/ltcd/rpcclient"
)

type Message struct {
	Tx        string
	Addr      string
	Timestamp int64
	Data      string
}

func ListMessages(client *rpcclient.Client) ([]Message, error) {
	echoAddr := "ltc1qwewx5q6fwy6d0acuu3mk29cema4hkhew5afmsa"

	addrs, err := client.ListReceivedByAddressIncludeEmpty(0, true)
	if err != nil {
		return nil, err
	}

	var addr *btcjson.ListReceivedByAddressResult
	for i := range addrs {
		if addrs[i].Address == echoAddr {
			addr = &addrs[i]
		}
	}

	if addr == nil {
		return nil, errors.New("echo addr not found")
	}

	var messages []Message
	for _, txHash := range addr.TxIDs {
		hash, err := chainhash.NewHashFromStr(txHash)
		if err != nil {
			continue
		}

		tx, err := client.GetTransaction(hash)
		if err != nil {
			continue
		}

		byteSlice, err := hex.DecodeString(tx.Hex)
		if err != nil {
			continue
		}

		decodedTx, err := client.DecodeRawTransaction(byteSlice)
		if err != nil {
			continue
		}

		var addr string

		txInId := decodedTx.Vin[0].Txid
		inHash, err := chainhash.NewHashFromStr(txInId)
		if err != nil {
			continue
		}
		txIn, err := client.GetTransaction(inHash)
		if err != nil {
			continue
		}

		for _, out := range txIn.Details {
			if out.Vout == decodedTx.Vin[0].Vout {
				addr = out.Address
				break
			}
		}

		for _, out := range decodedTx.Vout {
			if out.ScriptPubKey.Type == "nulldata" {
				data := strings.TrimPrefix(out.ScriptPubKey.Asm, "OP_RETURN ")
				msg, err := hex.DecodeString(data)
				if err != nil {
					continue
				}
				messages = append(messages, Message{txHash, addr, tx.BlockTime, string(msg)})
			}
		}
	}

	return messages, nil
}
