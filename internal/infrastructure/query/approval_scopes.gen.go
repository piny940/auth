// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"auth/internal/infrastructure/model"
)

func newApprovalScope(db *gorm.DB, opts ...gen.DOOption) approvalScope {
	_approvalScope := approvalScope{}

	_approvalScope.approvalScopeDo.UseDB(db, opts...)
	_approvalScope.approvalScopeDo.UseModel(&model.ApprovalScope{})

	tableName := _approvalScope.approvalScopeDo.TableName()
	_approvalScope.ALL = field.NewAsterisk(tableName)
	_approvalScope.ScopeID = field.NewInt32(tableName, "scope_id")
	_approvalScope.ApprovalID = field.NewInt64(tableName, "approval_id")

	_approvalScope.fillFieldMap()

	return _approvalScope
}

type approvalScope struct {
	approvalScopeDo

	ALL        field.Asterisk
	ScopeID    field.Int32
	ApprovalID field.Int64

	fieldMap map[string]field.Expr
}

func (a approvalScope) Table(newTableName string) *approvalScope {
	a.approvalScopeDo.UseTable(newTableName)
	return a.updateTableName(newTableName)
}

func (a approvalScope) As(alias string) *approvalScope {
	a.approvalScopeDo.DO = *(a.approvalScopeDo.As(alias).(*gen.DO))
	return a.updateTableName(alias)
}

func (a *approvalScope) updateTableName(table string) *approvalScope {
	a.ALL = field.NewAsterisk(table)
	a.ScopeID = field.NewInt32(table, "scope_id")
	a.ApprovalID = field.NewInt64(table, "approval_id")

	a.fillFieldMap()

	return a
}

func (a *approvalScope) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := a.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (a *approvalScope) fillFieldMap() {
	a.fieldMap = make(map[string]field.Expr, 2)
	a.fieldMap["scope_id"] = a.ScopeID
	a.fieldMap["approval_id"] = a.ApprovalID
}

func (a approvalScope) clone(db *gorm.DB) approvalScope {
	a.approvalScopeDo.ReplaceConnPool(db.Statement.ConnPool)
	return a
}

func (a approvalScope) replaceDB(db *gorm.DB) approvalScope {
	a.approvalScopeDo.ReplaceDB(db)
	return a
}

type approvalScopeDo struct{ gen.DO }

type IApprovalScopeDo interface {
	gen.SubQuery
	Debug() IApprovalScopeDo
	WithContext(ctx context.Context) IApprovalScopeDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IApprovalScopeDo
	WriteDB() IApprovalScopeDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IApprovalScopeDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IApprovalScopeDo
	Not(conds ...gen.Condition) IApprovalScopeDo
	Or(conds ...gen.Condition) IApprovalScopeDo
	Select(conds ...field.Expr) IApprovalScopeDo
	Where(conds ...gen.Condition) IApprovalScopeDo
	Order(conds ...field.Expr) IApprovalScopeDo
	Distinct(cols ...field.Expr) IApprovalScopeDo
	Omit(cols ...field.Expr) IApprovalScopeDo
	Join(table schema.Tabler, on ...field.Expr) IApprovalScopeDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IApprovalScopeDo
	RightJoin(table schema.Tabler, on ...field.Expr) IApprovalScopeDo
	Group(cols ...field.Expr) IApprovalScopeDo
	Having(conds ...gen.Condition) IApprovalScopeDo
	Limit(limit int) IApprovalScopeDo
	Offset(offset int) IApprovalScopeDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IApprovalScopeDo
	Unscoped() IApprovalScopeDo
	Create(values ...*model.ApprovalScope) error
	CreateInBatches(values []*model.ApprovalScope, batchSize int) error
	Save(values ...*model.ApprovalScope) error
	First() (*model.ApprovalScope, error)
	Take() (*model.ApprovalScope, error)
	Last() (*model.ApprovalScope, error)
	Find() ([]*model.ApprovalScope, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.ApprovalScope, err error)
	FindInBatches(result *[]*model.ApprovalScope, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.ApprovalScope) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IApprovalScopeDo
	Assign(attrs ...field.AssignExpr) IApprovalScopeDo
	Joins(fields ...field.RelationField) IApprovalScopeDo
	Preload(fields ...field.RelationField) IApprovalScopeDo
	FirstOrInit() (*model.ApprovalScope, error)
	FirstOrCreate() (*model.ApprovalScope, error)
	FindByPage(offset int, limit int) (result []*model.ApprovalScope, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IApprovalScopeDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (a approvalScopeDo) Debug() IApprovalScopeDo {
	return a.withDO(a.DO.Debug())
}

func (a approvalScopeDo) WithContext(ctx context.Context) IApprovalScopeDo {
	return a.withDO(a.DO.WithContext(ctx))
}

func (a approvalScopeDo) ReadDB() IApprovalScopeDo {
	return a.Clauses(dbresolver.Read)
}

func (a approvalScopeDo) WriteDB() IApprovalScopeDo {
	return a.Clauses(dbresolver.Write)
}

func (a approvalScopeDo) Session(config *gorm.Session) IApprovalScopeDo {
	return a.withDO(a.DO.Session(config))
}

func (a approvalScopeDo) Clauses(conds ...clause.Expression) IApprovalScopeDo {
	return a.withDO(a.DO.Clauses(conds...))
}

func (a approvalScopeDo) Returning(value interface{}, columns ...string) IApprovalScopeDo {
	return a.withDO(a.DO.Returning(value, columns...))
}

func (a approvalScopeDo) Not(conds ...gen.Condition) IApprovalScopeDo {
	return a.withDO(a.DO.Not(conds...))
}

func (a approvalScopeDo) Or(conds ...gen.Condition) IApprovalScopeDo {
	return a.withDO(a.DO.Or(conds...))
}

func (a approvalScopeDo) Select(conds ...field.Expr) IApprovalScopeDo {
	return a.withDO(a.DO.Select(conds...))
}

func (a approvalScopeDo) Where(conds ...gen.Condition) IApprovalScopeDo {
	return a.withDO(a.DO.Where(conds...))
}

func (a approvalScopeDo) Order(conds ...field.Expr) IApprovalScopeDo {
	return a.withDO(a.DO.Order(conds...))
}

func (a approvalScopeDo) Distinct(cols ...field.Expr) IApprovalScopeDo {
	return a.withDO(a.DO.Distinct(cols...))
}

func (a approvalScopeDo) Omit(cols ...field.Expr) IApprovalScopeDo {
	return a.withDO(a.DO.Omit(cols...))
}

func (a approvalScopeDo) Join(table schema.Tabler, on ...field.Expr) IApprovalScopeDo {
	return a.withDO(a.DO.Join(table, on...))
}

func (a approvalScopeDo) LeftJoin(table schema.Tabler, on ...field.Expr) IApprovalScopeDo {
	return a.withDO(a.DO.LeftJoin(table, on...))
}

func (a approvalScopeDo) RightJoin(table schema.Tabler, on ...field.Expr) IApprovalScopeDo {
	return a.withDO(a.DO.RightJoin(table, on...))
}

func (a approvalScopeDo) Group(cols ...field.Expr) IApprovalScopeDo {
	return a.withDO(a.DO.Group(cols...))
}

func (a approvalScopeDo) Having(conds ...gen.Condition) IApprovalScopeDo {
	return a.withDO(a.DO.Having(conds...))
}

func (a approvalScopeDo) Limit(limit int) IApprovalScopeDo {
	return a.withDO(a.DO.Limit(limit))
}

func (a approvalScopeDo) Offset(offset int) IApprovalScopeDo {
	return a.withDO(a.DO.Offset(offset))
}

func (a approvalScopeDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IApprovalScopeDo {
	return a.withDO(a.DO.Scopes(funcs...))
}

func (a approvalScopeDo) Unscoped() IApprovalScopeDo {
	return a.withDO(a.DO.Unscoped())
}

func (a approvalScopeDo) Create(values ...*model.ApprovalScope) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Create(values)
}

func (a approvalScopeDo) CreateInBatches(values []*model.ApprovalScope, batchSize int) error {
	return a.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (a approvalScopeDo) Save(values ...*model.ApprovalScope) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Save(values)
}

func (a approvalScopeDo) First() (*model.ApprovalScope, error) {
	if result, err := a.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.ApprovalScope), nil
	}
}

func (a approvalScopeDo) Take() (*model.ApprovalScope, error) {
	if result, err := a.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.ApprovalScope), nil
	}
}

func (a approvalScopeDo) Last() (*model.ApprovalScope, error) {
	if result, err := a.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.ApprovalScope), nil
	}
}

func (a approvalScopeDo) Find() ([]*model.ApprovalScope, error) {
	result, err := a.DO.Find()
	return result.([]*model.ApprovalScope), err
}

func (a approvalScopeDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.ApprovalScope, err error) {
	buf := make([]*model.ApprovalScope, 0, batchSize)
	err = a.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (a approvalScopeDo) FindInBatches(result *[]*model.ApprovalScope, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return a.DO.FindInBatches(result, batchSize, fc)
}

func (a approvalScopeDo) Attrs(attrs ...field.AssignExpr) IApprovalScopeDo {
	return a.withDO(a.DO.Attrs(attrs...))
}

func (a approvalScopeDo) Assign(attrs ...field.AssignExpr) IApprovalScopeDo {
	return a.withDO(a.DO.Assign(attrs...))
}

func (a approvalScopeDo) Joins(fields ...field.RelationField) IApprovalScopeDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Joins(_f))
	}
	return &a
}

func (a approvalScopeDo) Preload(fields ...field.RelationField) IApprovalScopeDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Preload(_f))
	}
	return &a
}

func (a approvalScopeDo) FirstOrInit() (*model.ApprovalScope, error) {
	if result, err := a.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.ApprovalScope), nil
	}
}

func (a approvalScopeDo) FirstOrCreate() (*model.ApprovalScope, error) {
	if result, err := a.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.ApprovalScope), nil
	}
}

func (a approvalScopeDo) FindByPage(offset int, limit int) (result []*model.ApprovalScope, count int64, err error) {
	result, err = a.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = a.Offset(-1).Limit(-1).Count()
	return
}

func (a approvalScopeDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = a.Count()
	if err != nil {
		return
	}

	err = a.Offset(offset).Limit(limit).Scan(result)
	return
}

func (a approvalScopeDo) Scan(result interface{}) (err error) {
	return a.DO.Scan(result)
}

func (a approvalScopeDo) Delete(models ...*model.ApprovalScope) (result gen.ResultInfo, err error) {
	return a.DO.Delete(models)
}

func (a *approvalScopeDo) withDO(do gen.Dao) *approvalScopeDo {
	a.DO = *do.(*gen.DO)
	return a
}
