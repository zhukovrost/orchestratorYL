package handlers

import (
	"encoding/json"
	"goCalculatorYL/internal/service"
	"goCalculatorYL/pkg/util"
	"net/http"
)

type OrchestratorHandler struct {
	srv *service.Service
}

func NewHandler(srv *service.Service) *OrchestratorHandler {
	srv.Logger.Debug("Setting up orchestrator handlers...")
	return &OrchestratorHandler{
		srv: srv,
	}
}

// AddExpressionHandler выполняет добавление вычисления арифметического выражения
func (h *OrchestratorHandler) AddExpressionHandler(w http.ResponseWriter, r *http.Request) {
	h.srv.Logger.Debug("new request")

	var calculationRequest service.CalculationRequest
	err := json.NewDecoder(r.Body).Decode(&calculationRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.srv.Logger.Error(err.Error())
		return
	}

	if calculationRequest.Expression == "" {
		http.Error(w, "Missing required fields", http.StatusUnprocessableEntity)
		return
	}

	if calculationRequest.ID == "" {
		calculationRequest.ID = util.GenerateId()
	}

	h.srv.AddExpression(&calculationRequest)

	// Формирование ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	h.srv.Logger.Debug("successful response (201)")
}

// GetExpressionsHandler выполняет получение списка выражений
func (h *OrchestratorHandler) GetExpressionsHandler(w http.ResponseWriter, r *http.Request) {

}

// GetExpressionByIDHandler выполняет получение списка выражений
func (h *OrchestratorHandler) GetExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {

}

// GetTaskHandler выполняет получение списка выражений
func (h *OrchestratorHandler) GetTaskHandler(w http.ResponseWriter, r *http.Request) {

}

// ResultHandler выполняет прием результата обработки данных
func (h *OrchestratorHandler) ResultHandler(w http.ResponseWriter, r *http.Request) {

}
