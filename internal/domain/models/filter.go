package models

type Filter struct {
	Field    string         `json:"field"`
	Operator FilterOperator `json:"operator"`
	Value    []string       `json:"value"`
}

type FilterOperator string

const (
	FilterOperator_Equals             FilterOperator = "="
	FilterOperator_NotEquals          FilterOperator = "!="
	FilterOperator_GreaterThan        FilterOperator = ">"
	FilterOperator_GreaterThanOrEqual FilterOperator = ">="
	FilterOperator_LessThan           FilterOperator = "<"
	FilterOperator_LessThanOrEqual    FilterOperator = "<="
	FilterOperator_Contains           FilterOperator = "contains"
	FilterOperator_NotContains        FilterOperator = "not_contains"
	FilterOperator_StartsWith         FilterOperator = "starts_with"
	FilterOperator_EndsWith           FilterOperator = "ends_with"
	FilterOperator_Between            FilterOperator = "between"
	FilterOperator_NotBetween         FilterOperator = "not_between"
	FilterOperator_In                 FilterOperator = "in"
	FilterOperator_NotIn              FilterOperator = "not_in"
)
