package logar

import (
	"fmt"
	"slices"
	"time"

	"github.com/Lexographics/logar/models"
	"gorm.io/gorm"
)

type PaginationStrategy int

const (
	PaginationStatus_None PaginationStrategy = iota
	PaginationStatus_Cursor
	PaginationStatus_Offset
)

type QueryOptFunc func(*QueryOptions)

type QueryOptions struct {
	Model              string
	Category           string
	Filters            []string
	FilterStructs      []models.Filter
	Severity           models.Severity
	PaginationStrategy PaginationStrategy
	Limit              int
	Page               int
	Cursor             int
	From               *time.Time
	To                 *time.Time
	IDs                []uint
	IDGreaterThan      uint
}

func WithModel(model string) QueryOptFunc {
	return func(o *QueryOptions) {
		o.Model = model
	}
}

func WithCategory(category string) QueryOptFunc {
	return func(o *QueryOptions) {
		o.Category = category
	}
}

func WithFilterStruct(filter models.Filter) QueryOptFunc {
	return func(o *QueryOptions) {
		o.FilterStructs = append(o.FilterStructs, filter)
	}
}

func WithFilter(filter string) QueryOptFunc {
	return func(o *QueryOptions) {
		o.Filters = append(o.Filters, filter)
	}
}

func WithSeverity(severity models.Severity) QueryOptFunc {
	return func(o *QueryOptions) {
		o.Severity = severity
	}
}

func WithOffsetPagination(offset int, page int) QueryOptFunc {
	return func(o *QueryOptions) {
		o.PaginationStrategy = PaginationStatus_Offset
		o.Page = page
		o.Limit = page
	}
}

func WithCursorPagination(cursor int, limit int) QueryOptFunc {
	return func(o *QueryOptions) {
		o.PaginationStrategy = PaginationStatus_Cursor
		o.Cursor = cursor
		o.Limit = limit
	}
}

func WithTimeRange(from time.Time, to time.Time) QueryOptFunc {
	return func(o *QueryOptions) {
		o.From = &from
		o.To = &to
	}
}

func After(from time.Time) QueryOptFunc {
	return func(o *QueryOptions) {
		o.From = &from
	}
}

func Before(to time.Time) QueryOptFunc {
	return func(o *QueryOptions) {
		o.To = &to
	}
}

func WithIDs(ids ...uint) QueryOptFunc {
	return func(o *QueryOptions) {
		o.IDs = ids
	}
}

func WithIDGreaterThan(id uint) QueryOptFunc {
	return func(o *QueryOptions) {
		o.IDGreaterThan = id
	}
}

func (l *AppImpl) GetLogs(opts ...QueryOptFunc) ([]models.Log, error) {
	var logs []models.Log
	query := l.prepareQuery(opts...)
	err := query.Find(&logs).Error
	return logs, err
}

func (l *AppImpl) DeleteLogs(opts ...QueryOptFunc) error {
	query := l.prepareQuery(opts...)
	return query.Delete(&models.Log{}).Error
}

func (l *AppImpl) prepareQuery(opts ...QueryOptFunc) *gorm.DB {
	options := &QueryOptions{
		Model:              "",
		Category:           "",
		Filters:            []string{},
		Severity:           models.Severity_None,
		PaginationStrategy: PaginationStatus_None,
		Limit:              0,
		Page:               0,
		Cursor:             0,
		From:               nil,
		To:                 nil,
		IDs:                nil,
		IDGreaterThan:      0,
	}
	for _, opt := range opts {
		opt(options)
	}

	query := l.db
	if options.Model != "" {
		query = query.Where("`model` = ?", options.Model)
	}
	if options.Category != "" {
		query = query.Where("category = ?", options.Category)
	}
	if len(options.Filters) > 0 {
		for _, filter := range options.Filters {
			query = query.Where("message LIKE ?", "%"+filter+"%")
		}
	}
	if options.Severity != models.Severity_None {
		query = query.Where("severity = ?", options.Severity)
	}

	if options.PaginationStrategy == PaginationStatus_Offset {
		query = query.Offset(options.Page * options.Limit)
	}
	if options.PaginationStrategy == PaginationStatus_Cursor {
		if options.Cursor > 0 {
			query = query.Where("id < ?", options.Cursor)
		}
	}
	if options.Limit != 0 {
		query = query.Limit(options.Limit)
	}

	if options.From != nil {
		query = query.Where("created_at >= ?", options.From)
	}
	if options.To != nil {
		query = query.Where("created_at <= ?", options.To)
	}

	if options.IDs != nil {
		query = query.Where("id IN (?)", options.IDs)
	}

	if options.IDGreaterThan > 0 {
		query = query.Where("id > ?", options.IDGreaterThan)
	}

	if len(options.FilterStructs) > 0 {
		for _, filter := range options.FilterStructs {
			if !slices.Contains(models.Log{}.FieldNames(), filter.Field) {
				continue
			}

			values := []any{}
			for _, v := range filter.Value {
				values = append(values, v)
			}

			if filter.Field == "created_at" {
				for i, v := range filter.Value {
					layout := "02-01-2006 15:04:05.000"
					if len(v) == 19 {
						layout = "02-01-2006 15:04:05"
					}

					t, err := time.Parse(layout, v)
					if err != nil {
						continue
					}

					values[i] = t
				}
			}

			for len(values) < 2 {
				values = append(values, "")
			}

			switch filter.Operator {
			case models.FilterOperator_Equals:
				query = query.Where(filter.Field+" = ?", values[0])
			case models.FilterOperator_NotEquals:
				query = query.Where(filter.Field+" != ?", values[0])
			case models.FilterOperator_GreaterThan:
				query = query.Where(filter.Field+" > ?", values[0])
			case models.FilterOperator_GreaterThanOrEqual:
				query = query.Where(filter.Field+" >= ?", values[0])
			case models.FilterOperator_LessThan:
				query = query.Where(filter.Field+" < ?", values[0])
			case models.FilterOperator_LessThanOrEqual:
				query = query.Where(filter.Field+" <= ?", values[0])
			case models.FilterOperator_Contains:
				query = query.Where(filter.Field+" LIKE ?", "%"+fmt.Sprint(values[0])+"%")
			case models.FilterOperator_NotContains:
				query = query.Where(filter.Field+" NOT LIKE ?", "%"+fmt.Sprint(values[0])+"%")
			case models.FilterOperator_StartsWith:
				query = query.Where(filter.Field+" LIKE ?", fmt.Sprint(values[0])+"%")
			case models.FilterOperator_EndsWith:
				query = query.Where(filter.Field+" LIKE ?", "%"+fmt.Sprint(values[0])+"%")
			case models.FilterOperator_Between:
				query = query.Where(filter.Field+" BETWEEN ? AND ?", values[0], values[1])
			case models.FilterOperator_NotBetween:
				query = query.Where(filter.Field+" NOT BETWEEN ? AND ?", values[0], values[1])
			case models.FilterOperator_In:
				query = query.Where(filter.Field+" IN (?)", values)
			case models.FilterOperator_NotIn:
				query = query.Where(filter.Field+" NOT IN (?)", values)
			}
		}
	}

	query = query.Order("id DESC")
	return query
}
