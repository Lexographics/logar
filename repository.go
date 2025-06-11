package logar

import (
	"fmt"
	"slices"
	"time"

	"gorm.io/gorm"
	"sadk.dev/logar/models"
)

type PaginationStrategy int

const (
	PaginationStatus_None PaginationStrategy = iota
	PaginationStatus_Cursor
	PaginationStatus_Offset
)

type QueryOptions struct {
	Model              string
	Category           string
	MessageContains    []string
	Filters            []models.Filter
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

type Query struct {
	Options *QueryOptions
}

func NewQuery() *Query {
	return &Query{
		Options: &QueryOptions{},
	}
}

func (q *Query) WithModel(model string) *Query {
	q.Options.Model = model
	return q
}

func (q *Query) WithCategory(category string) *Query {
	q.Options.Category = category
	return q
}

func (q *Query) WithFilter(filter models.Filter) *Query {
	q.Options.Filters = append(q.Options.Filters, filter)
	return q
}

func (q *Query) MessageContaints(text string) *Query {
	q.Options.MessageContains = append(q.Options.MessageContains, text)
	return q
}

func (q *Query) WithSeverity(severity models.Severity) *Query {
	q.Options.Severity = severity
	return q
}

func (q *Query) WithOffsetPagination(offset int, page int) *Query {
	q.Options.PaginationStrategy = PaginationStatus_Offset
	q.Options.Page = page
	q.Options.Limit = page
	return q
}

func (q *Query) WithCursorPagination(cursor int, limit int) *Query {
	q.Options.PaginationStrategy = PaginationStatus_Cursor
	q.Options.Cursor = cursor
	q.Options.Limit = limit
	return q
}

func (q *Query) WithTimeRange(from time.Time, to time.Time) *Query {
	q.Options.From = &from
	q.Options.To = &to
	return q
}

func (q *Query) After(from time.Time) *Query {
	q.Options.From = &from
	return q
}

func (q *Query) Before(to time.Time) *Query {
	q.Options.To = &to
	return q
}

func (q *Query) WithIDs(ids ...uint) *Query {
	q.Options.IDs = ids
	return q
}

func (q *Query) WithIDGreaterThan(id uint) *Query {
	q.Options.IDGreaterThan = id
	return q
}

func (l *AppImpl) GetLogs(q *Query) ([]models.Log, error) {
	var logs []models.Log
	query := l.prepareQuery(q.Options)
	err := query.Find(&logs).Error
	return logs, err
}

func (l *AppImpl) DeleteLogs(q *Query) error {
	query := l.prepareQuery(q.Options)
	return query.Delete(&models.Log{}).Error
}

func (l *AppImpl) prepareQuery(options *QueryOptions) *gorm.DB {
	query := l.db
	if options.Model != "" {
		query = query.Where("`model` = ?", options.Model)
	}
	if options.Category != "" {
		query = query.Where("category = ?", options.Category)
	}
	if len(options.MessageContains) > 0 {
		for _, filter := range options.MessageContains {
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

	if len(options.Filters) > 0 {
		for _, filter := range options.Filters {
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
