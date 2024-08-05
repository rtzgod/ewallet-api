package integration

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rtzgod/EWallet/internal/domain/entity"
	"github.com/rtzgod/EWallet/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MyNewIntegrationSuite struct {
	suite.Suite
	db    *sqlx.DB
	repos *repository.Repository
}

func TestMyNewIntegrationSuite(t *testing.T) {
	suite.Run(t, new(MyNewIntegrationSuite))
}

func (s *MyNewIntegrationSuite) SetupSuite() {

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		"localhost", "5432", "postgres", "1", "postgres", "disable")
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		logrus.Fatal(err)
	}
	s.db = db
}

func (s *MyNewIntegrationSuite) TearDownTest() {
	_, _ = s.db.Exec("TRUNCATE TABLE wallets, transactions")
}

func (s *MyNewIntegrationSuite) SetupTest() {
	s.repos = repository.NewRepository(s.db)
}

func (s *MyNewIntegrationSuite) TestWalletCreate() {
	walletId := "1"
	wallet, err := s.repos.Wallet.Create(walletId)
	s.Require().NoError(err)
	s.Require().NotNil(wallet)

	dbWallet := s.dbGetWallet(wallet.Id)

	s.Require().Equal(wallet.Id, dbWallet.Id)
	s.Require().Equal(wallet.Balance, dbWallet.Balance)
}

func (s *MyNewIntegrationSuite) TestGetWalletById() {
	walletId := "1"
	dbWallet := s.dbCreateWallet(walletId)
	wallet, err := s.repos.Wallet.GetById(walletId)
	s.Require().NoError(err)

	s.Require().Equal(wallet.Id, dbWallet.Id)
	s.Require().Equal(wallet.Balance, dbWallet.Balance)
}
func (s *MyNewIntegrationSuite) TestGetWalletNotFound() {
	_, err := s.repos.GetById("this id doesn't exist")
	s.Require().Error(err)
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
