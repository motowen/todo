package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

var (
	instance *Manager
)

func GetInstance() *Manager {
	return instance
}

type Manager struct {
	conn    *pgxpool.Pool
	context context.Context
}

type Config struct {
	Username                string
	Password                string
	Host                    string
	Port                    string
	TableName               string
	MinConnSize             int32
	MaxConnSize             int32
	MaxConnIdleTimeBySecond time.Duration
	MaxConnLifetimeBySecond time.Duration
}

func (manager *Manager) Setup(config Config) error {

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.TableName)

	pgxConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		fmt.Printf("postgres parse config fail, %+v\n", err)
		return err
	}

	background := context.Background()

	pgxConfig.MinConns = config.MinConnSize
	pgxConfig.MaxConns = config.MaxConnSize
	pgxConfig.MaxConnIdleTime = config.MaxConnIdleTimeBySecond * time.Second
	pgxConfig.MaxConnLifetime = config.MaxConnLifetimeBySecond * time.Second

	pool, err := pgxpool.ConnectConfig(background, pgxConfig)
	if err != nil {
		fmt.Printf("postgres connect fail, %+v\n", err)
		return err
	}

	err = pool.Ping(background)
	if err != nil {
		fmt.Printf("postgres ping fail, %+v\n", err)
		return err
	}

	instance = &Manager{
		conn:    pool,
		context: background,
	}

	return nil
}

func (manager *Manager) Query(sql string, args ...interface{}) (pgx.Rows, error) {
	return manager.conn.Query(manager.context, sql, args...)
}

func (manager *Manager) Exec(sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return manager.conn.Exec(manager.context, sql, arguments...)
}

func (manager *Manager) TxQuery(tx *pgx.Tx, sql string, args ...interface{}) (pgx.Rows, error) {
	return (*tx).Query(manager.context, sql, args...)
}

func (manager *Manager) TxExec(tx *pgx.Tx, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return (*tx).Exec(manager.context, sql, arguments...)
}

func (manager *Manager) DoInTx(fn func(*pgx.Tx) error) error {

	tx, err := manager.conn.Begin(manager.context)
	if err != nil {
		return err
	}

	defer tx.Rollback(manager.context)

	err = fn(&tx)
	if err != nil {
		return err
	}

	return tx.Commit(manager.context)
}
