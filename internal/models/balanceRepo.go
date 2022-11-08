package models

import (
	"context"
	"time"
)

type BalanceRepo interface {
	GetBalanceById(ctx context.Context, userId int64) (*Balance, error)
	DepositOnBalance(ctx context.Context, userId int64, amount float64) error
	ReserveAmount(ctx context.Context, userId int64, serviceId int64, orderId int64, amount float64) error
	ApproveOrder(ctx context.Context, userId int64, serviceId int64, orderId int64, amount float64) error
	DisapproveOrder(ctx context.Context, userId int64, serviceId int64, orderId int64, amount float64) error
	GetServiceMonthRevenueReport(ctx context.Context, month time.Month, year int64) ([]ServiceMonthRevenueRecord, error)
}
