package server

import (
	"fmt"
	"os"
	"strconv"
)

func makeLog(version int) {
	for {
		f, err := os.OpenFile("internal/logs/transaction_log_"+strconv.Itoa(version)+".txt",
			os.O_APPEND|os.O_WRONLY, 0777)
		defer f.Close()

		if err != nil {
			f, err = os.Create("internal/logs/transaction_log_" + strconv.Itoa(version) + ".txt")
			defer f.Close()
			if err != nil {
				fmt.Println(err)
			}

			_, err = f.WriteString(<-transactionManager)
			if err != nil {
				fmt.Println(err)
			}

			continue
		}

		_, err = f.WriteString("\n" + <-transactionManager)
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
