// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"auth/internal/infrastructure/model"
)

func newRedirectURI(db *gorm.DB, opts ...gen.DOOption) redirectURI {
	_redirectURI := redirectURI{}

	_redirectURI.redirectURIDo.UseDB(db, opts...)
	_redirectURI.redirectURIDo.UseModel(&model.RedirectURI{})

	tableName := _redirectURI.redirectURIDo.TableName()
	_redirectURI.ALL = field.NewAsterisk(tableName)
	_redirectURI.ID = field.NewInt64(tableName, "id")
	_redirectURI.ClientID = field.NewString(tableName, "client_id")
	_redirectURI.URI = field.NewString(tableName, "uri")
	_redirectURI.CreatedAt = field.NewTime(tableName, "created_at")
	_redirectURI.UpdatedAt = field.NewTime(tableName, "updated_at")

	_redirectURI.fillFieldMap()

	return _redirectURI
}

type redirectURI struct {
	redirectURIDo

	ALL       field.Asterisk
	ID        field.Int64
	ClientID  field.String
	URI       field.String
	CreatedAt field.Time
	UpdatedAt field.Time

	fieldMap map[string]field.Expr
}

func (r redirectURI) Table(newTableName string) *redirectURI {
	r.redirectURIDo.UseTable(newTableName)
	return r.updateTableName(newTableName)
}

func (r redirectURI) As(alias string) *redirectURI {
	r.redirectURIDo.DO = *(r.redirectURIDo.As(alias).(*gen.DO))
	return r.updateTableName(alias)
}

func (r *redirectURI) updateTableName(table string) *redirectURI {
	r.ALL = field.NewAsterisk(table)
	r.ID = field.NewInt64(table, "id")
	r.ClientID = field.NewString(table, "client_id")
	r.URI = field.NewString(table, "uri")
	r.CreatedAt = field.NewTime(table, "created_at")
	r.UpdatedAt = field.NewTime(table, "updated_at")

	r.fillFieldMap()

	return r
}

func (r *redirectURI) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := r.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (r *redirectURI) fillFieldMap() {
	r.fieldMap = make(map[string]field.Expr, 5)
	r.fieldMap["id"] = r.ID
	r.fieldMap["client_id"] = r.ClientID
	r.fieldMap["uri"] = r.URI
	r.fieldMap["created_at"] = r.CreatedAt
	r.fieldMap["updated_at"] = r.UpdatedAt
}

func (r redirectURI) clone(db *gorm.DB) redirectURI {
	r.redirectURIDo.ReplaceConnPool(db.Statement.ConnPool)
	return r
}

func (r redirectURI) replaceDB(db *gorm.DB) redirectURI {
	r.redirectURIDo.ReplaceDB(db)
	return r
}

type redirectURIDo struct{ gen.DO }

type IRedirectURIDo interface {
	gen.SubQuery
	Debug() IRedirectURIDo
	WithContext(ctx context.Context) IRedirectURIDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IRedirectURIDo
	WriteDB() IRedirectURIDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IRedirectURIDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IRedirectURIDo
	Not(conds ...gen.Condition) IRedirectURIDo
	Or(conds ...gen.Condition) IRedirectURIDo
	Select(conds ...field.Expr) IRedirectURIDo
	Where(conds ...gen.Condition) IRedirectURIDo
	Order(conds ...field.Expr) IRedirectURIDo
	Distinct(cols ...field.Expr) IRedirectURIDo
	Omit(cols ...field.Expr) IRedirectURIDo
	Join(table schema.Tabler, on ...field.Expr) IRedirectURIDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IRedirectURIDo
	RightJoin(table schema.Tabler, on ...field.Expr) IRedirectURIDo
	Group(cols ...field.Expr) IRedirectURIDo
	Having(conds ...gen.Condition) IRedirectURIDo
	Limit(limit int) IRedirectURIDo
	Offset(offset int) IRedirectURIDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IRedirectURIDo
	Unscoped() IRedirectURIDo
	Create(values ...*model.RedirectURI) error
	CreateInBatches(values []*model.RedirectURI, batchSize int) error
	Save(values ...*model.RedirectURI) error
	First() (*model.RedirectURI, error)
	Take() (*model.RedirectURI, error)
	Last() (*model.RedirectURI, error)
	Find() ([]*model.RedirectURI, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.RedirectURI, err error)
	FindInBatches(result *[]*model.RedirectURI, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.RedirectURI) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IRedirectURIDo
	Assign(attrs ...field.AssignExpr) IRedirectURIDo
	Joins(fields ...field.RelationField) IRedirectURIDo
	Preload(fields ...field.RelationField) IRedirectURIDo
	FirstOrInit() (*model.RedirectURI, error)
	FirstOrCreate() (*model.RedirectURI, error)
	FindByPage(offset int, limit int) (result []*model.RedirectURI, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Rows() (*sql.Rows, error)
	Row() *sql.Row
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IRedirectURIDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (r redirectURIDo) Debug() IRedirectURIDo {
	return r.withDO(r.DO.Debug())
}

func (r redirectURIDo) WithContext(ctx context.Context) IRedirectURIDo {
	return r.withDO(r.DO.WithContext(ctx))
}

func (r redirectURIDo) ReadDB() IRedirectURIDo {
	return r.Clauses(dbresolver.Read)
}

func (r redirectURIDo) WriteDB() IRedirectURIDo {
	return r.Clauses(dbresolver.Write)
}

func (r redirectURIDo) Session(config *gorm.Session) IRedirectURIDo {
	return r.withDO(r.DO.Session(config))
}

func (r redirectURIDo) Clauses(conds ...clause.Expression) IRedirectURIDo {
	return r.withDO(r.DO.Clauses(conds...))
}

func (r redirectURIDo) Returning(value interface{}, columns ...string) IRedirectURIDo {
	return r.withDO(r.DO.Returning(value, columns...))
}

func (r redirectURIDo) Not(conds ...gen.Condition) IRedirectURIDo {
	return r.withDO(r.DO.Not(conds...))
}

func (r redirectURIDo) Or(conds ...gen.Condition) IRedirectURIDo {
	return r.withDO(r.DO.Or(conds...))
}

func (r redirectURIDo) Select(conds ...field.Expr) IRedirectURIDo {
	return r.withDO(r.DO.Select(conds...))
}

func (r redirectURIDo) Where(conds ...gen.Condition) IRedirectURIDo {
	return r.withDO(r.DO.Where(conds...))
}

func (r redirectURIDo) Order(conds ...field.Expr) IRedirectURIDo {
	return r.withDO(r.DO.Order(conds...))
}

func (r redirectURIDo) Distinct(cols ...field.Expr) IRedirectURIDo {
	return r.withDO(r.DO.Distinct(cols...))
}

func (r redirectURIDo) Omit(cols ...field.Expr) IRedirectURIDo {
	return r.withDO(r.DO.Omit(cols...))
}

func (r redirectURIDo) Join(table schema.Tabler, on ...field.Expr) IRedirectURIDo {
	return r.withDO(r.DO.Join(table, on...))
}

func (r redirectURIDo) LeftJoin(table schema.Tabler, on ...field.Expr) IRedirectURIDo {
	return r.withDO(r.DO.LeftJoin(table, on...))
}

func (r redirectURIDo) RightJoin(table schema.Tabler, on ...field.Expr) IRedirectURIDo {
	return r.withDO(r.DO.RightJoin(table, on...))
}

func (r redirectURIDo) Group(cols ...field.Expr) IRedirectURIDo {
	return r.withDO(r.DO.Group(cols...))
}

func (r redirectURIDo) Having(conds ...gen.Condition) IRedirectURIDo {
	return r.withDO(r.DO.Having(conds...))
}

func (r redirectURIDo) Limit(limit int) IRedirectURIDo {
	return r.withDO(r.DO.Limit(limit))
}

func (r redirectURIDo) Offset(offset int) IRedirectURIDo {
	return r.withDO(r.DO.Offset(offset))
}

func (r redirectURIDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IRedirectURIDo {
	return r.withDO(r.DO.Scopes(funcs...))
}

func (r redirectURIDo) Unscoped() IRedirectURIDo {
	return r.withDO(r.DO.Unscoped())
}

func (r redirectURIDo) Create(values ...*model.RedirectURI) error {
	if len(values) == 0 {
		return nil
	}
	return r.DO.Create(values)
}

func (r redirectURIDo) CreateInBatches(values []*model.RedirectURI, batchSize int) error {
	return r.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (r redirectURIDo) Save(values ...*model.RedirectURI) error {
	if len(values) == 0 {
		return nil
	}
	return r.DO.Save(values)
}

func (r redirectURIDo) First() (*model.RedirectURI, error) {
	if result, err := r.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.RedirectURI), nil
	}
}

func (r redirectURIDo) Take() (*model.RedirectURI, error) {
	if result, err := r.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.RedirectURI), nil
	}
}

func (r redirectURIDo) Last() (*model.RedirectURI, error) {
	if result, err := r.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.RedirectURI), nil
	}
}

func (r redirectURIDo) Find() ([]*model.RedirectURI, error) {
	result, err := r.DO.Find()
	return result.([]*model.RedirectURI), err
}

func (r redirectURIDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.RedirectURI, err error) {
	buf := make([]*model.RedirectURI, 0, batchSize)
	err = r.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (r redirectURIDo) FindInBatches(result *[]*model.RedirectURI, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return r.DO.FindInBatches(result, batchSize, fc)
}

func (r redirectURIDo) Attrs(attrs ...field.AssignExpr) IRedirectURIDo {
	return r.withDO(r.DO.Attrs(attrs...))
}

func (r redirectURIDo) Assign(attrs ...field.AssignExpr) IRedirectURIDo {
	return r.withDO(r.DO.Assign(attrs...))
}

func (r redirectURIDo) Joins(fields ...field.RelationField) IRedirectURIDo {
	for _, _f := range fields {
		r = *r.withDO(r.DO.Joins(_f))
	}
	return &r
}

func (r redirectURIDo) Preload(fields ...field.RelationField) IRedirectURIDo {
	for _, _f := range fields {
		r = *r.withDO(r.DO.Preload(_f))
	}
	return &r
}

func (r redirectURIDo) FirstOrInit() (*model.RedirectURI, error) {
	if result, err := r.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.RedirectURI), nil
	}
}

func (r redirectURIDo) FirstOrCreate() (*model.RedirectURI, error) {
	if result, err := r.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.RedirectURI), nil
	}
}

func (r redirectURIDo) FindByPage(offset int, limit int) (result []*model.RedirectURI, count int64, err error) {
	result, err = r.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = r.Offset(-1).Limit(-1).Count()
	return
}

func (r redirectURIDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = r.Count()
	if err != nil {
		return
	}

	err = r.Offset(offset).Limit(limit).Scan(result)
	return
}

func (r redirectURIDo) Scan(result interface{}) (err error) {
	return r.DO.Scan(result)
}

func (r redirectURIDo) Delete(models ...*model.RedirectURI) (result gen.ResultInfo, err error) {
	return r.DO.Delete(models)
}

func (r *redirectURIDo) withDO(do gen.Dao) *redirectURIDo {
	r.DO = *do.(*gen.DO)
	return r
}
