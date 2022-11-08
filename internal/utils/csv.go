package utils

import (
	"avito-balance-service/internal/models"
	"encoding/csv"
	"errors"
	"os"
	"strconv"
	"strings"
)

func getReportPath(month int64, year int64) string {
	var sb strings.Builder

	sb.WriteString("/month_reports/report_")
	sb.WriteString(strconv.Itoa(int(month)))
	sb.WriteRune('_')
	sb.WriteString(strconv.Itoa(int(year)))
	sb.WriteString(".csv")

	return sb.String()
}

func getReportLink(path string) string {
	var sb strings.Builder

	sb.WriteString("/api/v1/balance")
	sb.WriteString(path)

	return sb.String()
}

func SaveCSV(month int64, year int64, records []models.ServiceMonthRevenueRecord) (string, error) {
	path := getReportPath(month, year)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir("month_reports", os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if err = os.Truncate(path, 0); err != nil {
		return "", err
	}

	writer := csv.NewWriter(f)
	defer writer.Flush()

	csvRecords := make([][]string, len(records))
	for i := 0; i < len(records); i++ {
		csvRecords[i] = records[i].ToCSV()
		if err = writer.Write(csvRecords[i]); err != nil {
			return "", err
		}
	}

	return getReportLink(path), nil
}
