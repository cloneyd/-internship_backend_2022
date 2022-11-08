package service

type balanceRequestJSON struct {
	UserId int64 `json:"user_id" validate:"required"`
}

type depositRequestJSON struct {
	UserId int64   `json:"user_id" validate:"required"`
	Amount float64 `json:"amount" validate:"required"`
}

type orderRequestJSON struct {
	UserId    int64   `json:"user_id" validate:"required"`
	ServiceId int64   `json:"service_id" validate:"required"`
	OrderId   int64   `json:"order_id" validate:"required"`
	Price     float64 `json:"price" validate:"required"`
}

type revenueReportRequestJson struct {
	Month int64 `json:"month" validate:"required"`
	Year  int64 `json:"year" validate:"required"`
}
