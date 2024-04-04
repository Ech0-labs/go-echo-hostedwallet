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
	Tx        string `json:"tx"`
	Addr      string `json:"addr"`
	Timestamp int64  `json:"timestamp"`
	Data      string `json:"data"`
}

func getTxsFromAddr(client *rpcclient.Client, addr string) ([]string, error) {
	addrs, err := client.ListReceivedByAddressIncludeEmpty(0, true)
	if err != nil {
		return nil, err
	}

	for i := range addrs {
		if addrs[i].Address == addr {
			return addrs[i].TxIDs, nil
		}
	}

	return nil, errors.New("addr not found")
}

func getRawTx(client *rpcclient.Client, strHash string) (*btcjson.TxRawResult, error) {
	hash, err := chainhash.NewHashFromStr(strHash)
	if err != nil {
		return nil, err
	}

	tx, err := client.GetTransaction(hash)
	if err != nil {
		return nil, err
	}

	byteSlice, err := hex.DecodeString(tx.Hex)
	if err != nil {
		return nil, err
	}

	rawTx, err := client.DecodeRawTransaction(byteSlice)
	return rawTx, err
}

func getSender(client *rpcclient.Client, tx *btcjson.TxRawResult) (string, error) {
	hash, err := chainhash.NewHashFromStr(tx.Vin[0].Txid)
	if err != nil {
		return "", err
	}

	txIn, err := client.GetTransaction(hash)
	if err != nil {
		return "", err
	}

	for _, out := range txIn.Details {
		if out.Vout == tx.Vin[0].Vout {
			return out.Address, nil
		}
	}

	return "", errors.New("vout not found")
}

func getTime(client *rpcclient.Client, txHash string) (int64, error) {
	hash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return 0, err
	}

	tx, err := client.GetTransaction(hash)
	if err != nil {
		return 0, err
	}

	blockHash, err := chainhash.NewHashFromStr(tx.BlockHash)
	if err != nil {
		return 0, err
	}

	header, err := client.GetBlockHeader(blockHash)
	if err != nil {
		return 0, err
	}

	return header.Timestamp.Unix(), nil
}

func getOpReturn(outs []btcjson.Vout) string {
	for _, out := range outs {
		if out.ScriptPubKey.Type == "nulldata" {
			data := strings.TrimPrefix(out.ScriptPubKey.Asm, "OP_RETURN ")
			msg, err := hex.DecodeString(data)
			if err != nil {
				continue
			}
			return string(msg)
		}
	}

	return ""
}

func ListMessages(client *rpcclient.Client) ([]Message, error) {
	txIds, err := getTxsFromAddr(client, echoAddr)
	if err != nil {
		return nil, err
	}

	var messages []Message
	for _, txHash := range txIds {
		rawTx, err := getRawTx(client, txHash)
		if err != nil {
			continue
		}

		addr, err := getSender(client, rawTx)
		if err != nil {
			continue
		}

		msg := getOpReturn(rawTx.Vout)

		time, err := getTime(client, rawTx.Txid)
		if err != nil {
			continue
		}

		messages = append(messages, Message{txHash, addr, time, msg})
	}

	return messages, nil
}
