package integration

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/rtzgod/EWallet/internal/domain/entity"
	"github.com/rtzgod/EWallet/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
	"time"
)

var DB *sqlx.DB

type MyNewIntegrationSuite struct {
	suite.Suite
	db    *sqlx.DB
	repos *repository.Repository
}

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		logrus.Fatalf("Could not construct pool: %s", err)
	}
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("[1] Could not connect to docker: %s", err)
	}
	resource, err := pool.BuildAndRun("integration_tests_db", "fixtures/Dockerfile", []string{})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://testuser:123@%s/testdb?sslmode=disable", hostAndPort)

	if err = pool.Retry(func() error {
		var err error
		DB, err = sqlx.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return DB.Ping()
	}); err != nil {
		log.Fatalf("[2] Could not connect to docker: %s", err)
	}

	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	m.Run()
}

func TestMyNewIntegrationSuite(t *testing.T) {
	suite.Run(t, new(MyNewIntegrationSuite))
}

func (s *MyNewIntegrationSuite) SetupSuite() {
	s.db = DB
}

func (s *MyNewIntegrationSuite) TearDownSuite() {
	if err := s.db.Close(); err != nil {
		logrus.Fatal(err)
	}
}

func (s *MyNewIntegrationSuite) SetupTest() {
	s.repos = repository.NewRepository(s.db)
}

func (s *MyNewIntegrationSuite) TearDownTest() {
	_, _ = s.db.Exec("TRUNCATE TABLE wallets, transactions")
}

func (s *MyNewIntegrationSuite) Test_WalletCreate() {
	walletId := "1"
	wallet, err := s.repos.Wallet.Create(walletId)
	s.Require().NoError(err)
	s.Require().NotNil(wallet)

	dbWallet := s.dbGetWallet(wallet.Id)

	s.Require().Equal(wallet.Id, dbWallet.Id)
	s.Require().Equal(wallet.Balance, dbWallet.Balance)
}

func (s *MyNewIntegrationSuite) Test_WalletGetById() {
	walletId := "1"
	dbWallet := s.dbCreateWallet(walletId)
	wallet, err := s.repos.Wallet.GetById(walletId)
	s.Require().NoError(err)

	s.Require().Equal(wallet.Id, dbWallet.Id)
	s.Require().Equal(wallet.Balance, dbWallet.Balance)
}
func (s *MyNewIntegrationSuite) Test_WalletGetByIdNotFound() {
	_, err := s.repos.Wallet.GetById("this id doesn't exist")
	s.Require().Error(err)
}
func (s *MyNewIntegrationSuite) Test_WalletUpdate() {
	firstWalletId := "1"
	secondWalletId := "2"
	dbFirstWallet := s.dbCreateWallet(firstWalletId)
	dbSecondWallet := s.dbCreateWallet(secondWalletId)
	err := s.repos.Wallet.Update(dbFirstWallet.Id, dbSecondWallet.Id, 10)
	s.Require().NoError(err)
	dbFirstWallet = s.dbGetWallet(firstWalletId)
	dbSecondWallet = s.dbGetWallet(secondWalletId)
	s.Require().Equal(90.0, dbFirstWallet.Balance)
	s.Require().Equal(110.0, dbSecondWallet.Balance)
}
func (s *MyNewIntegrationSuite) Test_WalletUpdateNoReceiver() {
	firstWalletId := "1"
	dbFirstWallet := s.dbCreateWallet(firstWalletId)
	err := s.repos.Wallet.Update(dbFirstWallet.Id, "not existing walletId", 50)
	s.Require().Error(err)
	dbFirstWallet = s.dbGetWallet(firstWalletId)
	s.Require().Equal(100.0, dbFirstWallet.Balance)
}
func (s *MyNewIntegrationSuite) Test_WalletUpdateNoSender() {
	secondWalletId := "2"
	dbSecondWallet := s.dbCreateWallet(secondWalletId)
	err := s.repos.Wallet.Update("not existing walletId", dbSecondWallet.Id, 10)
	s.Require().Error(err)
	dbSecondWallet = s.dbGetWallet(secondWalletId)
	s.Require().Equal(100.0, dbSecondWallet.Balance)
}

func (s *MyNewIntegrationSuite) Test_TransactionCreate() {
	senderWalletId := "1"
	receiverWalletId := "2"
	Amount := 10.0
	err := s.repos.Transaction.Create(senderWalletId, receiverWalletId, Amount)
	s.Require().NoError(err)
	transaction := s.dbGetTransaction(senderWalletId)
	s.Require().Equal(Amount, transaction.Amount)
	s.Require().Equal(senderWalletId, transaction.SenderId)
	s.Require().Equal(receiverWalletId, transaction.ReceiverId)
}

func (s *MyNewIntegrationSuite) Test_TransactionGetAllById() {
	s.dbCreateTransaction("1", "2", 10.0)
	s.dbCreateTransaction("2", "1", 50.0)
	transactions, err := s.repos.Transaction.GetAllById("1")
	s.Require().NoError(err)
	s.Require().Len(transactions, 2)
	s.Require().Equal("1", transactions[0].SenderId)
	s.Require().Equal("1", transactions[1].ReceiverId)
	s.Require().Equal("2", transactions[1].SenderId)
	s.Require().Equal("2", transactions[0].ReceiverId)
	s.Require().Equal(10.0, transactions[0].Amount)
	s.Require().Equal(50.0, transactions[1].Amount)
}
func (s *MyNewIntegrationSuite) Test_TransactionGetAllByIdNoTransactions() {
	transactions, _ := s.repos.Transaction.GetAllById("3")
	s.Require().Len(transactions, 0)
}

// db prefix means direct sql query to db

func (s *MyNewIntegrationSuite) dbCreateWallet(walletId string) entity.Wallet {
	var wallet entity.Wallet
	query := fmt.Sprintf("insert into %s (id, balance) values ($1, $2)", "wallets")

	if _, err := s.db.Exec(query, walletId, 100.0); err != nil {
		s.Fail(err.Error())
	}

	wallet = entity.Wallet{
		Id:      walletId,
		Balance: 100.0,
	}
	return wallet
}

func (s *MyNewIntegrationSuite) dbGetWallet(walletId string) entity.Wallet {
	var wallet entity.Wallet
	query := fmt.Sprintf("select * from %s where id=$1", "wallets")
	err := s.db.Get(&wallet, query, walletId)
	if err != nil {
		s.Fail(err.Error())
	}
	return wallet
}

func (s *MyNewIntegrationSuite) dbCreateTransaction(senderId, receiverId string, amount float64) {
	query := fmt.Sprintf("insert into %s (time, sender_id, receiver_id, amount) values ($1, $2, $3, $4)", "transactions")
	_, err := s.db.Exec(query, time.Now().Format(time.RFC3339), senderId, receiverId, amount)
	if err != nil {
		s.Fail(err.Error())
	}
}

func (s *MyNewIntegrationSuite) dbGetTransaction(senderId string) entity.Transaction {
	var transaction entity.Transaction
	query := fmt.Sprintf("select * from %s where sender_id=$1", "transactions")
	err := s.db.Get(&transaction, query, senderId)
	if err != nil {
		s.Fail(err.Error())
	}
	return transaction
}
