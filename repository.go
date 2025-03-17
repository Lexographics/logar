package logar

import (
	"time"

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
	Filter             string
	Severity           Severity
	PaginationStrategy PaginationStrategy
	Limit              int
	Page               int
	Cursor             int
	From               *time.Time
	To                 *time.Time
	IDs                []uint
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

func WithFilter(filter string) QueryOptFunc {
	return func(o *QueryOptions) {
		o.Filter = filter
	}
}

func WithSeverity(severity Severity) QueryOptFunc {
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

func (l *Logger) GetLogs(opts ...QueryOptFunc) ([]Log, error) {
	var logs []Log
	query := l.prepareQuery(opts...)
	err := query.Find(&logs).Error
	return logs, err
}

func (l *Logger) DeleteLogs(opts ...QueryOptFunc) error {
	query := l.prepareQuery(opts...)
	return query.Delete(&Log{}).Error
}

func (l *Logger) prepareQuery(opts ...QueryOptFunc) *gorm.DB {
	options := &QueryOptions{
		Model:              "",
		Category:           "",
		Filter:             "",
		Severity:           Severity_None,
		PaginationStrategy: PaginationStatus_None,
		Limit:              0,
		Page:               0,
		Cursor:             0,
		From:               nil,
		To:                 nil,
		IDs:                nil,
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
	if options.Filter != "" {
		query = query.Where("message LIKE ?", "%"+options.Filter+"%")
	}
	if options.Severity != Severity_None {
		query = query.Where("severity = ?", options.Severity)
	}

	if options.PaginationStrategy == PaginationStatus_Offset {
		query = query.Offset(options.Page * options.Limit).Limit(options.Limit)
	}
	if options.PaginationStrategy == PaginationStatus_Cursor {
		query = query.Limit(options.Limit)
		if options.Cursor > 0 {
			query = query.Where("id < ?", options.Cursor)
		}
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

	query = query.Order("id DESC")
	return query
}
