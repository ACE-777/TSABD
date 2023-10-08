package server

import (
	"fmt"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"os"
	"strconv"
)

func makeLog(version int) {
	for {
		inputTransaction := <-transactionManagerGlobal
		if inputTransaction.Source == source {
			IDTransaction++
		}

		if clock[inputTransaction.Source] > inputTransaction.Id && inputTransaction.Source != source {
			continue
		}

		clock[inputTransaction.Source]++
		patch, err := jsonpatch.DecodePatch([]byte(inputTransaction.Payload))
		if err != nil {
			fmt.Printf("error in decoding patch: %v", err)
		}

		modified, err := patch.Apply([]byte(snap))
		if err != nil {
			fmt.Errorf("error in making new patch: %v", err)
		}

		snap = string(modified)
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

func applyJSONPatch(sourceJSON []byte, ops []jsonpatch.Operation) ([]byte, error) {
	document, err := jsonpatch.DecodePatch(sourceJSON)
	if err != nil {
		return nil, err
	}

	patchedJSON, err := document.Apply(sourceJSON)
	if err != nil {
		return nil, err
	}

	return patchedJSON, nil
}
