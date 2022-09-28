package sap_api_wrapper

import (
	"strconv"
	"strings"
)

type SapApiQueryParams struct {
	// Fields to select for the returned result
	Select []string
	// Field to use for ordering the returned result
	OrderBy []string
	// SQL `WHERE` statement to use for filtering the result
	Filter string
	// The offset to use for pagination
	Skip int
	// The limit to use for pagination
	Top int
}

func (p *SapApiQueryParams) AsReqParams() map[string]string {
	queryParams := make(map[string]string)
	if p.Select != nil {
		queryParams["$select"] = strings.Join(p.Select, ",")
	}
	if p.OrderBy != nil {
		queryParams["$orderby"] = strings.Join(p.OrderBy, ",")
	}
	if p.Filter != "" {
		queryParams["$filter"] = p.Filter
	}
	if p.Skip != 0 {
		queryParams["$skip"] = strconv.Itoa(p.Skip)
	}
	if p.Top != 0 {
		queryParams["$top"] = strconv.Itoa(p.Top)
	}

	return queryParams
}
