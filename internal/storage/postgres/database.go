package postgres

import (
	"avito-balance-service/config"
	"avito-balance-service/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type DB struct {
	pool *pgxpool.Pool
}

var balanceColumns = []string{"user_id", "balance"}

func NewDB(ctx context.Context, cfg *config.Config) (*DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&pool_max_conns=%d", cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName, cfg.Postgres.SSLMode, cfg.Postgres.PoolMaxConns)
	pgxCfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, pgxCfg)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &DB{pool: pool}, nil
}

func (db *DB) GetBalanceById(ctx context.Context, id int64) (*models.Balance, error) {
	sqlQuery := "SELECT * FROM balances WHERE user_id=$1"
	var balance models.Balance

	err := db.pool.QueryRow(ctx, sqlQuery, id).Scan(&balance.UserId, &balance.Amount)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}

func (db *DB) DepositOnBalance(ctx context.Context, id int64, amount float64) error {
	sqlQuery := "CALL deposit_on_balance($1, $2)"

	if _, err := db.pool.Exec(ctx, sqlQuery, id, amount); err != nil {
		return err
	}

	return nil
}

func (db *DB) ReserveAmount(ctx context.Context, userId int64, serviceId int64, orderId int64, amount float64) error {
	sqlQuery := "CALL reserve_amount($1, $2, $3, $4)"

	if _, err := db.pool.Exec(ctx, sqlQuery, userId, serviceId, orderId, amount); err != nil {
		return err
	}

	return nil
}

func (db *DB) ApproveOrder(ctx context.Context, userId int64, serviceId int64, orderId int64, amount float64) error {
	sqlQuery := "CALL approve_order($1, $2, $3, $4)"

	if _, err := db.pool.Exec(ctx, sqlQuery, userId, serviceId, orderId, amount); err != nil {
		return err
	}

	return nil
}

func (db *DB) DisapproveOrder(ctx context.Context, userId int64, serviceId int64, orderId int64, amount float64) error {
	sqlQuery := "CALL disapprove_order($1, $2, $3, $4)"

	if _, err := db.pool.Exec(ctx, sqlQuery, userId, serviceId, orderId, amount); err != nil {
		return err
	}

	return nil
}

func (db *DB) GetServiceMonthRevenueReport(ctx context.Context, month time.Month, year int64) ([]models.ServiceMonthRevenueRecord, error) {
	sqlQuery := "SELECT * FROM get_service_month_revenue_report($1)"
	date := time.Date(int(year), month, 1, 0, 0, 0, 0, time.UTC)

	rows, err := db.pool.Query(ctx, sqlQuery, date.Format("01-02-2006"))
	if err != nil {
		return nil, err
	}

	var records []models.ServiceMonthRevenueRecord
	for rows.Next() {
		var record models.ServiceMonthRevenueRecord
		if err = rows.Scan(&record.Name, &record.Amount); err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}
