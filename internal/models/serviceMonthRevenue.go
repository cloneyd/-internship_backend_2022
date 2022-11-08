package models

import (
	"fmt"
)

type ServiceMonthRevenueRecord struct {
	Name   string
	Amount float64
}

func (smrr *ServiceMonthRevenueRecord) ToCSV() []string {
	return []string{smrr.Name, fmt.Sprintf("%f", smrr.Amount)}
}
