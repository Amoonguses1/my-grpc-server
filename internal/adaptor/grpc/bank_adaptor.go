package grpc

import (
	"context"
	"time"

	"github.com/amoonguses1/grpc-proto-study/protogen/go/bank"
	"google.golang.org/genproto/googleapis/type/date"
)

func (a *GrpcAdaptor) GetCurrentBalance(ctx context.Context, req *bank.CurrentBalanceRequest) (*bank.CurrentBalanceResponse, error) {
	now := time.Now()
	balance := a.bankService.FindCurrentBalance(req.AccountNumber)

	return &bank.CurrentBalanceResponse{
		Amount: balance,
		CurrentDate: &date.Date{
			Year:  int32(now.Year()),
			Month: int32(now.Month()),
			Day:   int32(now.Day()),
		},
	}, nil
}
