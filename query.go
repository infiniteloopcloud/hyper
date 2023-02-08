package hyper

import (
	"net/url"
	"strconv"
	"time"

	"github.com/volatiletech/null/v9"
)

// GetQueryStringParam getting a query param if exist or return with null.String{Valid: false}
func GetQueryStringParam(query url.Values, param string) null.String {
	if query.Has(param) {
		return null.StringFrom(query.Get(param))
	}
	return null.String{}
}

// GetQueryUint64Param getting a query param if exist or return with null.Uint64{Valid: false}
func GetQueryUint64Param(query url.Values, param string) null.Uint64 {
	if query.Has(param) {
		res, err := strconv.Atoi(query.Get(param))
		if err != nil {
			return null.Uint64{}
		}
		return null.Uint64From(uint64(res))
	}
	return null.Uint64{}
}

// GetQueryTimeParam getting a query param if exist or return with null.Time{Valid: false}
func GetQueryTimeParam(query url.Values, param, format string) null.Time {
	if query.Has(param) {
		res, err := time.Parse(format, query.Get(param))
		if err != nil {
			return null.Time{}
		}
		return null.TimeFrom(res)
	}
	return null.Time{}
}
