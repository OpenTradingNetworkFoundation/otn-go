package api

import (
	"fmt"

	"github.com/opentradingnetworkfoundation/otn-go/objects"
)

type FullAccountHandler func(acc *objects.FullAccount) error

func LoadAccounts(rpc DatabaseAPI, handler FullAccountHandler) error {
	count, err := rpc.GetAccountCount()
	if err != nil {
		return err
	}

	batchSize := uint64(20)

	for i := uint64(0); i < count; i += batchSize {
		currentBatch := batchSize
		if i+batchSize >= count {
			currentBatch = count - i
		}
		ids := make([]objects.GrapheneObject, currentBatch)
		for j := uint64(0); j < currentBatch; j++ {
			ids[j] = objects.NewGrapheneID(objects.ObjectID(fmt.Sprintf("1.2.%d", i+j)))
		}
		accs, err := rpc.GetFullAccounts(false, ids...)
		if err != nil {
			return err
		}

		for idx := range accs {
			if err := handler(&accs[idx].FullAccount); err != nil {
				return err
			}
		}
	}

	return nil
}
