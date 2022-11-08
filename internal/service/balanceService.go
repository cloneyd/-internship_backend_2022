package service

import (
	"avito-balance-service/internal/utils"
	"encoding/json"
	"net/http"
	"time"

	"avito-balance-service/internal/models"

	"github.com/go-playground/validator/v10"
)

type BalanceService struct {
	validator *validator.Validate
	storage   models.BalanceRepo
}

func NewBalanceService(storage models.BalanceRepo) *BalanceService {
	return &BalanceService{
		storage:   storage,
		validator: validator.New(),
	}
}

func (bs *BalanceService) HandleGetBalance(w http.ResponseWriter, r *http.Request) {
	var balanceRequest balanceRequestJSON

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&balanceRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(newErrorResponse(err))
		return
	}

	if err := bs.validator.Struct(balanceRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(newErrorResponse(err))
		return
	}

	balance, err := bs.storage.GetBalanceById(r.Context(), balanceRequest.UserId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(newErrorResponse(err))
		return
	}

	json.NewEncoder(w).Encode(balance)
}

func (bs *BalanceService) HandleDepositOnBalance(w http.ResponseWriter, r *http.Request) {
	var depositRequest depositRequestJSON

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&depositRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(newErrorResponse(err))
		return
	}

	if err := bs.validator.Struct(depositRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(newErrorResponse(err))
		return
	}

	err := bs.storage.DepositOnBalance(r.Context(), depositRequest.UserId, depositRequest.Amount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(newErrorResponse(err))
		return
	}

	balance, err := bs.storage.GetBalanceById(r.Context(), depositRequest.UserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(newErrorResponse(err))
		return
	}

	json.NewEncoder(w).Encode(balance)
}

func (bs *BalanceService) HandleReserveAmount(w http.ResponseWriter, r *http.Request) {
	var orderRequest orderRequestJSON

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(newErrorResponse(err))
		return
	}

	if err := bs.validator.Struct(orderRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(newErrorResponse(err))
		return
	}

	err := bs.storage.ReserveAmount(r.Context(), orderRequest.UserId, orderRequest.ServiceId, orderRequest.OrderId, orderRequest.Price)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(newErrorResponse(err))
		return
	}

	balance, err := bs.storage.GetBalanceById(r.Context(), orderRequest.UserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(newErrorResponse(err))
		return
	}

	json.NewEncoder(w).Encode(balance)
}

func (bs *BalanceService) HandleApproveOrder(w http.ResponseWriter, r *http.Request) {
	var orderRequest orderRequestJSON

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(newErrorResponse(err))
		return
	}

	if err := bs.validator.Struct(orderRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(newErrorResponse(err))
		return
	}

	err := bs.storage.ApproveOrder(r.Context(), orderRequest.UserId, orderRequest.ServiceId, orderRequest.OrderId, orderRequest.Price)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(newErrorResponse(err))
		return
	}

	json.NewEncoder(w).Encode(newResultResponse("order approved"))
}

func (bs *BalanceService) HandleDisapproveOrder(w http.ResponseWriter, r *http.Request) {
	var orderRequest orderRequestJSON

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(newErrorResponse(err))
		return
	}

	if err := bs.validator.Struct(orderRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(newErrorResponse(err))
		return
	}

	err := bs.storage.DisapproveOrder(r.Context(), orderRequest.UserId, orderRequest.ServiceId, orderRequest.OrderId, orderRequest.Price)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(newErrorResponse(err))
		return
	}

	json.NewEncoder(w).Encode(newResultResponse("order disapproved"))
}

func (bs *BalanceService) HandleServiceMonthRevenueReport(w http.ResponseWriter, r *http.Request) {
	var reportRequest revenueReportRequestJson

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&reportRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(newErrorResponse(err))
		return
	}

	if err := bs.validator.Struct(reportRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(newErrorResponse(err))
		return
	}

	serviceMonthRevenueRecords, err := bs.storage.GetServiceMonthRevenueReport(r.Context(), time.Month(reportRequest.Month), reportRequest.Year)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(newErrorResponse(err))
		return
	}

	link, err := utils.SaveCSV(reportRequest.Month, reportRequest.Year, serviceMonthRevenueRecords)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(newErrorResponse(err))
		return
	}

	json.NewEncoder(w).Encode(newResultResponse(link))
}
