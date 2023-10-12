package server

import (
	"fmt"
	"os"
	"strconv"

	jsonpatch "github.com/evanphx/json-patch/v5"
)

func makeLog(version int) {
	for {
		inputTransaction := <-transactionManagerGlobal
		if inputTransaction.Source == source && inputTransaction.Id >= IDTransaction {
			IDTransaction = inputTransaction.Id - IDTransaction
		}

		_, ok := clock[inputTransaction.Source]
		if !ok {
			clock[inputTransaction.Source] = -1
		}

		if clock[inputTransaction.Source] > inputTransaction.Id && inputTransaction.Source != source {
			continue
		}

		clock[inputTransaction.Source] = inputTransaction.Id
		fmt.Println("input", inputTransaction.Payload)
		patch, err := jsonpatch.DecodePatch([]byte("[" + inputTransaction.Payload + "]"))
		if err != nil {
			fmt.Printf("error in decoding patch: %v", err)
		}

		modified, err := patch.Apply([]byte(snap))
		if err != nil {
			fmt.Errorf("error in making new patch: %v", err)
		}

		snap = string(modified)
		fmt.Println("actual snap", string(snap))
		wal = append(wal, inputTransaction)

		if inputTransaction.Source == source {
			f, err := os.OpenFile("internal/logs/transaction_log_"+strconv.Itoa(version)+".txt",
				os.O_APPEND|os.O_WRONLY, 0777)
			defer f.Close()

			if err != nil {
				f, err = os.Create("internal/logs/transaction_log_" + strconv.Itoa(version) + ".txt")
				defer f.Close()
				if err != nil {
					fmt.Println(err)
				}

				_, err = f.WriteString(inputTransaction.Payload)
				if err != nil {
					fmt.Println(err)
				}

				continue
			}

			_, err = f.WriteString("\n" + inputTransaction.Payload)
			if err != nil {
				fmt.Println(err)
			}

			select {
			case <-timer.C:
				version++
			default:
				continue
			}
		}

	}
}
