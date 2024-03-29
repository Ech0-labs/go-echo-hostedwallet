package flow

import (
	"github.com/ltcsuite/ltcd/rpcclient"
)

type Addr struct {
	Addr   string  `json:"addr"`
	Label  string  `json:"label"`
	Amount float64 `json:"amount"`
}

func ListAddr(client *rpcclient.Client) ([]Addr, error) {
	addrs, err := client.ListReceivedByAddressIncludeEmpty(0, true)
	if err != nil {
		return nil, err
	}

	var res []Addr

	for _, addr := range addrs {
		if addr.InvolvesWatchonly {
			continue
		}

		info, err := client.GetAddressInfo(addr.Address)
		if err != nil {
			return nil, err
		}

		label := ""
		if len(info.Labels) > 0 {
			label = info.Labels[0]
		}
		res = append(res, Addr{Addr: addr.Address, Label: label, Amount: addr.Amount})
	}

	return res, nil
}
