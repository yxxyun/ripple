package main

import (
	"context"
	"fmt"

	"github.com/yxxyun/ripple/crypto"
	"github.com/yxxyun/ripple/data"
	"github.com/yxxyun/ripple/rpc"
	"google.golang.org/grpc"
)

func main() {
	seed, _ := crypto.NewRippleHashCheck("sapqGRrejEA8Z3mbAGqiuBNak4HHs", crypto.RIPPLE_FAMILY_SEED)
	key, _ := crypto.NewECDSAKey(seed.Payload())
	zero := uint32(0)
	src, _ := data.NewAccountFromAddress("rpvsmj5w2iuFv78SyN44NcJt9nrskfRYCA")
	dest, _ := data.NewAccountFromAddress("r3mQPV1CKHo67VdpUmu4iWVfyTxNTiDujn")
	amount, _ := data.NewAmount("100/XRP")
	fee, _ := data.NewNativeValue(int64(12))
	txflag := new(data.TransactionFlag)
	*txflag = *txflag | data.TxCanonicalSignature
	conn, err := grpc.Dial("test.xrp.xpring.io:50051", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close()

	client := rpc.NewXRPLedgerAPIServiceClient(conn)

	addr := rpc.AccountAddress{
		Address: "rpvsmj5w2iuFv78SyN44NcJt9nrskfRYCA",
	}
	acctinforeq := rpc.GetAccountInfoRequest{
		Account: &addr,
	}
	acctinforesp, _ := client.GetAccountInfo(context.Background(), &acctinforeq)
	ledgerSequence := acctinforesp.GetLedgerIndex() + 3
	accountSequence := acctinforesp.GetAccountData().GetSequence().GetValue()
	tx := &data.Payment{
		Destination: *dest,
		Amount:      *amount,
	}
	tx.TxBase = data.TxBase{
		Account:            *src,
		Fee:                *fee,
		TransactionType:    data.PAYMENT,
		Flags:              txflag,
		LastLedgerSequence: &ledgerSequence,
		Sequence:           accountSequence,
	}
	data.Sign(tx, key, &zero)
	_, raw, _ := data.Raw(tx)
	submittx := rpc.SubmitTransactionRequest{
		SignedTransaction: raw,
	}
	resp, _ := client.SubmitTransaction(context.Background(), &submittx)
	fmt.Println(resp)
}
