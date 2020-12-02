package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yxxyun/ripple/data"
	"github.com/yxxyun/ripple/terminal"
	"github.com/yxxyun/ripple/websockets"
)

func checkErr(err error, quit bool) {
	if err != nil {
		terminal.Println(err.Error(), terminal.Default)
		if quit {
			os.Exit(1)
		}
	}
}

var (
	host     = flag.String("host", "wss://s1.ripple.com:443", "websockets host to connect to")
	proposed = flag.Bool("proposed", false, "include proposed transacions")
)

func main() {
	flag.Parse()
	r, err := websockets.NewRemote(*host, true)
	checkErr(err, true)

	for {
		if r.Conn {
			confirmation, err := r.Subscribe(true, !*proposed, *proposed, true)
			checkErr(err, false)
			terminal.Println(fmt.Sprint("Subscribed at: ", confirmation.LedgerSequence), terminal.Default)
			// go func() {

			// Consume messages as they arrive
			for {
				if !r.Conn {
					break
				}
				msg, ok := <-r.Incoming
				if !ok {
					break
				}

				switch msg := msg.(type) {
				case *websockets.LedgerStreamMsg:
					terminal.Println(msg, terminal.Default)
				case *websockets.TransactionStreamMsg:
					terminal.Println(&msg.Transaction, terminal.Indent)
					for _, path := range msg.Transaction.PathSet() {
						terminal.Println(path, terminal.DoubleIndent)
					}
					trades, err := data.NewTradeSlice(&msg.Transaction)
					checkErr(err, false)
					for _, trade := range trades {
						terminal.Println(trade, terminal.DoubleIndent)
					}
					balances, err := msg.Transaction.Balances()
					checkErr(err, false)
					for _, balance := range balances {
						terminal.Println(balance, terminal.DoubleIndent)
					}
				case *websockets.ServerStreamMsg:
					terminal.Println(msg, terminal.Default)
				}
			}
		}
		//}()

	}

}
