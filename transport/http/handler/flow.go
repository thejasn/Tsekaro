package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/thejasn/tester/cerrors"
	"github.com/thejasn/tester/domain/flow/model"
	"github.com/thejasn/tester/pkg/log"
	"github.com/thejasn/tester/service"
)

type flowhandler struct {
	svc service.Flow
}

func NewFlowHandler(fs service.Flow) flowhandler {
	return flowhandler{
		svc: fs,
	}
}

func (f flowhandler) ConfigFlowsRouter(router chi.Router) {
	router.Get("/flows", f.GetAllFlows)
	router.Post("/flows", f.AddFlow)
	router.Get("/flows/{id}", f.GetFlow)
	router.Put("/flows/{id}", f.UpdateFlow)
	router.Delete("/flows/{id}", f.DeleteFlow)
	router.Get("/flows/execute/{id}", f.ExecuteFlow)
}

// GetAllFlows is a function to get a slice of record(s) from flow table in the tester database
// @Summary Get list of Flow
// @Tags Flow
// @Description GetAllFlow is a handler to get a slice of record(s) from flow table in the tester database
// @Accept  json
// @Produce  json
// @Param   page     query    int     false        "page requested (defaults to 0)"
// @Param   pagesize query    int     false        "number of records in a page  (defaults to 20)"
// @Param   order    query    string  false        "db sort order column"
// @Success 200 {object} api.PagedResults{data=[]model.Flow}
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /flows [get]
// http http://localhost:8080/flows?page=0&pagesize=20
func (f flowhandler) GetAllFlows(w http.ResponseWriter, r *http.Request) {
	page, err := readInt(r, "page", 0)
	if err != nil || page < 0 {
		returnError(w, r, cerrors.ErrBadParams)
		return
	}

	pagesize, err := readInt(r, "pagesize", 20)
	if err != nil || pagesize <= 0 {
		returnError(w, r, cerrors.ErrBadParams)
		return
	}

	order := r.FormValue("order")

	records, totalRows, err := f.svc.GetAll(log.WithLogger(r.Context(), log.Init()), page, pagesize, order)
	if err != nil {
		returnError(w, r, err)
		return
	}

	result := &PagedResults{Page: page, PageSize: pagesize, Data: records, TotalRecords: totalRows}
	writeJSON(w, result)
}

// GetFlow is a function to get a single record to flow table in the tester database
// @Summary Get record from table Flow by id
// @Tags Flow
// @ID record id
// @Description GetFlow is a function to get a single record to flow table in the tester database
// @Accept  json
// @Produce  json
// @Param  id path int true "record id"
// @Success 200 {object} model.Flow
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError "ErrNotFound, db record for id not found - returns NotFound HTTP 404 not found error"
// @Router /flows/{id} [get]
// http http://localhost:8080/flows/1
func (f flowhandler) GetFlow(w http.ResponseWriter, r *http.Request) {

	id, err := parseInt(chi.URLParam(r, "id"))
	if err != nil {
		returnError(w, r, err)
		return
	}

	record, err := f.svc.Get(log.WithLogger(r.Context(), log.Init()), id)
	if err != nil {
		returnError(w, r, err)
		return
	}

	writeJSON(w, record)
}

// AddFlow add to add a single record to flow table in the tester database
// @Summary Add an record to flow table
// @Description add to add a single record to flow table in the tester database
// @Tags Flow
// @Accept  json
// @Produce  json
// @Param Flow body model.Flow true "Add Flow"
// @Success 200 {object} model.Flow
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /flows [post]
// echo '{"id": 86}' | http POST http://localhost:8080/flows
func (f flowhandler) AddFlow(w http.ResponseWriter, r *http.Request) {
	flow := &model.Flow{}

	if err := readJSON(r, flow); err != nil {
		returnError(w, r, cerrors.ErrBadParams)
		return
	}

	var err error
	flow, _, err = f.svc.Add(log.WithLogger(r.Context(), log.Init()), flow)
	if err != nil {
		returnError(w, r, err)
		return
	}

	writeJSON(w, flow)
}

// UpdateFlow Update a single record from flow table in the tester database
// @Summary Update an record in table flow
// @Description Update a single record from flow table in the tester database
// @Tags Flow
// @Accept  json
// @Produce  json
// @Param  id path int true "Account ID"
// @Param  Flow body model.Flow true "Update Flow record"
// @Success 200 {object} model.Flow
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /flows/{id} [patch]
// echo '{"id": 86}' | http PATCH http://localhost:8080/flows/1
func (f flowhandler) UpdateFlow(w http.ResponseWriter, r *http.Request) {
	id, err := parseInt(chi.URLParam(r, "id"))
	if err != nil {
		returnError(w, r, err)
		return
	}

	flow := &model.Flow{}
	if err := readJSON(r, flow); err != nil {
		returnError(w, r, cerrors.ErrBadParams)
		return
	}

	flow, _, err = f.svc.Update(log.WithLogger(r.Context(), log.Init()), id, flow)
	if err != nil {
		returnError(w, r, err)
		return
	}

	writeJSON(w, flow)
}

// DeleteFlow Delete a single record from flow table in the tester database
// @Summary Delete a record from flow
// @Description Delete a single record from flow table in the tester database
// @Tags Flow
// @Accept  json
// @Produce  json
// @Param  id path int true "ID" Format(int64)
// @Success 204 {object} model.Flow
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /flows/{id} [delete]
// http DELETE http://localhost:8080/flows/1
func (f flowhandler) DeleteFlow(w http.ResponseWriter, r *http.Request) {
	id, err := parseInt(chi.URLParam(r, "id"))
	if err != nil {
		returnError(w, r, err)
		return
	}

	rowsAffected, err := f.svc.Delete(log.WithLogger(r.Context(), log.Init()), id)
	if err != nil {
		returnError(w, r, err)
		return
	}

	writeRowsAffected(w, rowsAffected)
}

func (f flowhandler) ExecuteFlow(w http.ResponseWriter, r *http.Request) {

	id, err := parseInt(chi.URLParam(r, "id"))
	if err != nil {
		returnError(w, r, err)
		return
	}

	record, err := f.svc.Execute(log.WithLogger(r.Context(), log.Init()), id)
	if err != nil {
		returnError(w, r, err)
		return
	}

	writeJSON(w, record)
}
