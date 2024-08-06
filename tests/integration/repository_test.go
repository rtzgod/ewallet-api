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
