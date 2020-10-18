package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/thejasn/tester/cerrors"
	"github.com/thejasn/tester/domain/testcase/model"
	"github.com/thejasn/tester/pkg/log"
	"github.com/thejasn/tester/service"
)

type testcasehandler struct {
	svc service.Testcase
}

func NewTestcaseHandler(ts service.Testcase) testcasehandler {
	return testcasehandler{
		svc: ts,
	}
}

func (t testcasehandler) ConfigTestcasesRouter(router chi.Router) {
	router.Get("/testcases", t.GetAllTestcases)
	router.Post("/testcases", t.AddTestcase)
	router.Get("/testcases/{id}", t.GetTestcase)
	router.Put("/testcases/{id}", t.UpdateTestcase)
	router.Delete("/testcases/{id}", t.DeleteTestcase)
	router.Get("/testcases/execute/{id}", t.ExecuteTestcase)
}

// GetAllTestcases is a function to get a slice of record(s) from testcase table in the tester database
// @Summary Get list of Testcase
// @Tags Testcase
// @Description GetAllTestcase is a handler to get a slice of record(s) from testcase table in the tester database
// @Accept  json
// @Produce  json
// @Param   page     query    int     false        "page requested (defaults to 0)"
// @Param   pagesize query    int     false        "number of records in a page  (defaults to 20)"
// @Param   order    query    string  false        "db sort order column"
// @Success 200 {object} api.PagedResults{data=[]model.Testcase}
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /testcases [get]
// http http://localhost:8080/testcases?page=0&pagesize=20
func (t testcasehandler) GetAllTestcases(w http.ResponseWriter, r *http.Request) {
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

	records, totalRows, err := t.svc.GetAll(log.WithLogger(r.Context(), log.Init()), page, pagesize, order)
	if err != nil {
		returnError(w, r, err)
		return
	}

	result := &PagedResults{Page: page, PageSize: pagesize, Data: records, TotalRecords: totalRows}
	writeJSON(w, result)
}

// GetTestcase is a function to get a single record to testcase table in the tester database
// @Summary Get record from table Testcase by id
// @Tags Testcase
// @ID record id
// @Description GetTestcase is a function to get a single record to testcase table in the tester database
// @Accept  json
// @Produce  json
// @Param  id path int true "record id"
// @Success 200 {object} model.Testcase
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError "ErrNotFound, db record for id not found - returns NotFound HTTP 404 not found error"
// @Router /testcases/{id} [get]
// http http://localhost:8080/testcases/1
func (t testcasehandler) GetTestcase(w http.ResponseWriter, r *http.Request) {

	id, err := parseInt(chi.URLParam(r, "id"))
	if err != nil {
		returnError(w, r, err)
		return
	}

	record, err := t.svc.Get(log.WithLogger(r.Context(), log.Init()), id)
	if err != nil {
		returnError(w, r, err)
		return
	}

	writeJSON(w, record)
}

// AddTestcase add to add a single record to testcase table in the tester database
// @Summary Add an record to testcase table
// @Description add to add a single record to testcase table in the tester database
// @Tags Testcase
// @Accept  json
// @Produce  json
// @Param Testcase body model.Testcase true "Add Testcase"
// @Success 200 {object} model.Testcase
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /testcases [post]
// echo '{"id": 42}' | http POST http://localhost:8080/testcases
func (t testcasehandler) AddTestcase(w http.ResponseWriter, r *http.Request) {
	testcase := &model.Testcase{}

	if err := readJSON(r, testcase); err != nil {
		returnError(w, r, cerrors.ErrBadParams)
		return
	}

	var err error
	testcase, _, err = t.svc.Add(log.WithLogger(r.Context(), log.Init()), testcase)
	if err != nil {
		returnError(w, r, err)
		return
	}

	writeJSON(w, testcase)
}

// UpdateTestcase Update a single record from testcase table in the tester database
// @Summary Update an record in table testcase
// @Description Update a single record from testcase table in the tester database
// @Tags Testcase
// @Accept  json
// @Produce  json
// @Param  id path int true "Account ID"
// @Param  Testcase body model.Testcase true "Update Testcase record"
// @Success 200 {object} model.Testcase
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /testcases/{id} [patch]
// echo '{"id": 42}' | http PATCH http://localhost:8080/testcases/1
func (t testcasehandler) UpdateTestcase(w http.ResponseWriter, r *http.Request) {
	id, err := parseInt(chi.URLParam(r, "id"))
	if err != nil {
		returnError(w, r, err)
		return
	}

	testcase := &model.Testcase{}
	if err := readJSON(r, testcase); err != nil {
		returnError(w, r, cerrors.ErrBadParams)
		return
	}

	testcase, _, err = t.svc.Update(log.WithLogger(r.Context(), log.Init()), id, testcase)
	if err != nil {
		returnError(w, r, err)
		return
	}

	writeJSON(w, testcase)
}

// DeleteTestcase Delete a single record from testcase table in the tester database
// @Summary Delete a record from testcase
// @Description Delete a single record from testcase table in the tester database
// @Tags Testcase
// @Accept  json
// @Produce  json
// @Param  id path int true "ID" Format(int64)
// @Success 204 {object} model.Testcase
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /testcases/{id} [delete]
// http DELETE http://localhost:8080/testcases/1
func (t testcasehandler) DeleteTestcase(w http.ResponseWriter, r *http.Request) {
	id, err := parseInt(chi.URLParam(r, "id"))
	if err != nil {
		returnError(w, r, err)
		return
	}

	rowsAffected, err := t.svc.Delete(log.WithLogger(r.Context(), log.Init()), id)
	if err != nil {
		returnError(w, r, err)
		return
	}

	writeRowsAffected(w, rowsAffected)
}

func (t testcasehandler) ExecuteTestcase(w http.ResponseWriter, r *http.Request) {

	id, err := parseInt(chi.URLParam(r, "id"))
	if err != nil {
		returnError(w, r, err)
		return
	}

	record, err := t.svc.Execute(log.WithLogger(r.Context(), log.Init()), id)
	if err != nil {
		returnError(w, r, err)
		return
	}

	writeJSON(w, record)
}
