package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yxxyun/ripple/crypto"
	"github.com/yxxyun/ripple/data"
	"github.com/yxxyun/ripple/websockets"
)

func checkErr(err error, quit bool) {
	if err != nil {
		fmt.Println(err.Error())
		if quit {
			os.Exit(1)
		}
	}
}

var (
	host = flag.String("host", "wss://xrpl.ws:443", "websockets host to connect to")
)

func main() {
	flag.Parse()
	fmt.Println("")
	fmt.Println("瑞波账号删除，5.0XRP做为手续费销毁，其余XRP发送到接收账号")
	var srcacct, sec, destacct string
	var tag uint32
	fmt.Println("输入要删除的账号地址：")
	fmt.Scanln(&srcacct)
	src, err := data.NewAccountFromAddress(srcacct)
	if err != nil {
		fmt.Println("瑞波账号错误")
		os.Exit(1)
	}

	fmt.Println("输入密钥：")
	fmt.Scanln(&sec)

	seed, err := crypto.NewRippleHashCheck(sec, crypto.RIPPLE_FAMILY_SEED)
	if err != nil {
		fmt.Println("密钥错误")
		os.Exit(1)
	}
	key, _ := crypto.NewECDSAKey(seed.Payload())
	zero := uint32(0)

	fmt.Println("输入接收账号：")
	fmt.Scanln(&destacct)
	dest, err := data.NewAccountFromAddress(destacct)
	if err != nil {
		fmt.Println("瑞波账号错误")
		os.Exit(1)
	}
	fmt.Println("输入tag，如没有可随意填写：")
	fmt.Scanln(&tag)
	acctdelfee, _ := data.NewNativeValue(int64(5000000))
	txflag := new(data.TransactionFlag)
	*txflag = *txflag | data.TxCanonicalSignature
	trustflag := new(data.TransactionFlag)
	*trustflag = *trustflag | data.TxSetNoRipple
	*trustflag = *trustflag | data.TxCanonicalSignature
	r, err := websockets.NewRemote(*host, false)
	checkErr(err, true)

	airesult, err := r.AccountInfo(*src)
	checkErr(err, true)
	accountSequence := *airesult.AccountData.Sequence
	ledgerSequence := airesult.LedgerSequence + 6
	tx := &data.AccountDelete{
		Destination:    *dest,
		DestinationTag: &tag,
	}
	tx.TxBase = data.TxBase{
		Account:            *src,
		Fee:                *acctdelfee,
		TransactionType:    data.ACCOUNT_DELETE,
		Flags:              txflag,
		LastLedgerSequence: &ledgerSequence,
		Sequence:           accountSequence,
	}
	checkErr(data.Sign(tx, key, &zero), true)
	ret, err := r.Submit(tx, true)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(ret)
	}
}
