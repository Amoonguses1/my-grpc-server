package main

import (
	"database/sql"
	"log"
	"math/rand"
	"time"

	dbmigration "github.com/amoonguses1/my-grpc-server/db"
	mydb "github.com/amoonguses1/my-grpc-server/internal/adaptor/database"
	mygrpc "github.com/amoonguses1/my-grpc-server/internal/adaptor/grpc"
	app "github.com/amoonguses1/my-grpc-server/internal/application"
	"github.com/amoonguses1/my-grpc-server/internal/application/domain/bank"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(logWriter{})

	sqlDB, err := sql.Open("pgx", "postgres://harumaki@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalln("Can not connect database :", err)
	}

	dbmigration.Migrate(sqlDB)

	databaseAdaptor, err := mydb.NewDatabaseAdaptor(sqlDB)
	if err != nil {
		log.Fatalln("Cannot create database adaptor :", err)
	}

	hs := &app.HelloService{}
	bs := app.NewBankService(databaseAdaptor)

	go generateExchangeRates(bs, "USD", "JPN", 5*time.Second)

	grpcAdaptor := mygrpc.NewGrpcAdaptor(hs, bs, 9090)
	grpcAdaptor.Run()
}

func runDummyOrm(da *mydb.DatabaseAdaptor) {
	now := time.Now()

	uuid, _ := da.Save(
		&mydb.DummyOrm{
			UserID:    uuid.New(),
			UserName:  "hoge" + time.Now().Format("15:04:05"),
			CreatedAt: now,
			UpdatedAt: now,
		},
	)

	res, _ := da.GetByUuid(&uuid)

	log.Println("res :", res)
}

func generateExchangeRates(bs *app.BankService, fromCurrency, ToCurrency string, duration time.Duration) {
	ticker := time.NewTicker(duration)

	for range ticker.C {
		now := time.Now()
		validFrom := now.Truncate(time.Second).Add(3 * time.Second)
		validTo := validFrom.Add(duration).Add(-1 * time.Millisecond)

		dummyRate := bank.ExchangeRate{
			FromCurrency:       fromCurrency,
			ToCurrency:         ToCurrency,
			ValidFromTimestamp: validFrom,
			ValidToTimestamp:   validTo,
			Rate:               2000 + float64(rand.Intn(300)),
		}

		bs.CreateExchangeRate(dummyRate)
	}
}
