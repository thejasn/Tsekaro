package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thejasn/tester/cerrors"
	"github.com/thejasn/tester/utils/httputil"
)

// PagedResults results for pages GetAll results.
type PagedResults struct {
	Page         int64       `json:"page"`
	PageSize     int64       `json:"page_size"`
	Data         interface{} `json:"data"`
	TotalRecords int64       `json:"total_records"`
}

// HTTPError example
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

func readInt(r *http.Request, param string, v int64) (int64, error) {
	p := r.FormValue(param)
	if p == "" {
		return v, nil
	}

	return strconv.ParseInt(p, 10, 64)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	data, _ := json.Marshal(v)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(data)
}

func writeRowsAffected(w http.ResponseWriter, rowsAffected int64) {
	data, _ := json.Marshal(rowsAffected)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(data)
}

func readJSON(r *http.Request, v interface{}) error {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(buf, v)
}

func returnError(w http.ResponseWriter, r *http.Request, err error) {
	status := 0
	switch err {
	case cerrors.ErrNotFound:
		status = http.StatusBadRequest
	case cerrors.ErrUnableToMarshalJSON:
		status = http.StatusBadRequest
	case cerrors.ErrUpdateFailed:
		status = http.StatusBadRequest
	case cerrors.ErrInsertFailed:
		status = http.StatusBadRequest
	case cerrors.ErrDeleteFailed:
		status = http.StatusBadRequest
	case cerrors.ErrBadParams:
		status = http.StatusBadRequest
	default:
		status = http.StatusBadRequest
	}
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}

	httputil.SendJSON(w, r, er.Code, er)
}

// NewError example
func NewError(ctx *gin.Context, status int, err error) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}

func parseInt(key string) (int, error) {
	id, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return -1, err
	}
	return int(id), err
}

func parseInt32(query url.Values, key string) (int32, error) {
	idStr := query.Get(key)
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return -1, err
	}
	return int32(id), err
}

func parseInt64(query url.Values, key string) (int64, error) {
	idStr := query.Get(key)
	id, err := strconv.ParseInt(idStr, 10, 54)
	if err != nil {
		return -1, err
	}
	return int64(id), err
}

func parseInterface(query url.Values, key string) (interface{}, error) {
	idStr := query.Get(key)
	return idStr, nil
}

func parseString(query url.Values, key string) (string, error) {
	idStr := query.Get(key)
	return idStr, nil
}
