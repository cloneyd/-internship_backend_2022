package server

import (
	"avito-balance-service/internal/service"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) MapHandlers() {
	balanceService := service.NewBalanceService(s.db)

	router := mux.NewRouter().StrictSlash(true)
	{
		api := router.PathPrefix("/api").Subrouter()
		{
			v1 := api.PathPrefix("/v1").Subrouter()
			{
				balanceApiV1 := v1.PathPrefix("/balance").Subrouter()

				balanceApiV1.HandleFunc("/", balanceService.HandleGetBalance).Methods(http.MethodGet)
				balanceApiV1.HandleFunc("/deposit", balanceService.HandleDepositOnBalance).Methods(http.MethodPost)
				balanceApiV1.HandleFunc("/reserve", balanceService.HandleReserveAmount).Methods(http.MethodPost)
				balanceApiV1.HandleFunc("/reserve/approve", balanceService.HandleApproveOrder).Methods(http.MethodPost)
				balanceApiV1.HandleFunc("/reserve/disapprove", balanceService.HandleDisapproveOrder).Methods(http.MethodPost)
				balanceApiV1.HandleFunc("/month_report", balanceService.HandleServiceMonthRevenueReport).Methods(http.MethodPost)
				balanceApiV1.HandleFunc("/user_report", nil).Methods(http.MethodPost)
				fileServer := http.FileServer(http.Dir("/month_reports/"))
				balanceApiV1.Handle("/month_reports/{rest}", http.StripPrefix("/api/v1/balance/month_reports/", fileServer)).Methods(http.MethodGet)
			}
			v1.Headers("Content-Type", "application/json")
		}
	}

	s.httpServer.Handler = router
}
