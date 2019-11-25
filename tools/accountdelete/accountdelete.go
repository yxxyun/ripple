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
	host = flag.String("host", "wss://s.devnet.rippletest.net:51233", "websockets host to connect to")
)

func main() {
	flag.Parse()
	fmt.Println("")
	fmt.Println("瑞波账号删除，5.0XRP做为手续费销毁，其余XRP发送到接收账号")
	fmt.Println("如有非XRP资产或信任线存在，将自动归还资产到发行者地址并删除信任线")
	var sec, destacct string
	var tag uint32
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
	minfee, _ := data.NewNativeValue(int64(12))
	txflag := new(data.TransactionFlag)
	*txflag = *txflag | data.TxCanonicalSignature
	trustflag := new(data.TransactionFlag)
	*trustflag = *trustflag | data.TxSetNoRipple
	*trustflag = *trustflag | data.TxCanonicalSignature
	r, err := websockets.NewRemote(*host)
	checkErr(err, true)
	var account data.Account
	copy(account[:], key.Id(&zero))

	airesult, err := r.AccountInfo(account)
	checkErr(err, true)
	accountSequence := *airesult.AccountData.Sequence
	ledgerSequence := airesult.LedgerSequence + 6
	if airesult.LedgerSequence-accountSequence < 256 {
		fmt.Println("稍后再试，The current ledger index must be at least 256 higher than the account's sequence number.")
		os.Exit(1)
	}
	alresult, err := r.AccountLines(account, "current")
	checkErr(err, true)
	if alresult.Lines.Len() > 0 {
		fmt.Println("存在信任线，将自动归还非XRP资产到发行者地址，并删除信任线")
		for _, line := range alresult.Lines {
			if !line.Balance.IsZero() && !line.Balance.IsNegative() {
				sendAmount, _ := data.NewAmount(line.Balance.String() + "/" + line.Asset().String())
				payment := &data.Payment{
					Destination:    line.Account,
					DestinationTag: &zero,
					Amount:         *sendAmount,
				}
				payment.TxBase = data.TxBase{
					Account:            account,
					Fee:                *minfee,
					TransactionType:    data.PAYMENT,
					Flags:              txflag,
					LastLedgerSequence: &ledgerSequence,
					Sequence:           accountSequence,
				}
				data.Sign(payment, key, &zero)
				r.Submit(payment)
				accountSequence++
				ledgerSequence++
			}
			if !line.Balance.IsNegative() {
				amount, _ := data.NewAmount("0/" + line.Asset().String())
				trust := &data.TrustSet{
					LimitAmount: *amount,
					QualityIn:   &zero,
					QualityOut:  &zero,
				}
				trust.TxBase = data.TxBase{
					Account:            account,
					TransactionType:    data.TRUST_SET,
					Flags:              trustflag,
					Fee:                *minfee,
					LastLedgerSequence: &ledgerSequence,
					Sequence:           accountSequence,
				}
				data.Sign(trust, key, &zero)
				r.Submit(trust)
				accountSequence++
				ledgerSequence++
			}
		}
	}
	tx := &data.AccountDelete{
		Destination:    *dest,
		DestinationTag: &tag,
	}
	tx.TxBase = data.TxBase{
		Account:            account,
		Fee:                *acctdelfee,
		TransactionType:    data.ACCOUNT_DELETE,
		Flags:              txflag,
		LastLedgerSequence: &ledgerSequence,
		Sequence:           accountSequence,
	}
	checkErr(data.Sign(tx, key, &zero), true)
	ret, err := r.Submit(tx)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(ret)
	}
}
