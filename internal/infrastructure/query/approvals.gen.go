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

func newApproval(db *gorm.DB, opts ...gen.DOOption) approval {
	_approval := approval{}

	_approval.approvalDo.UseDB(db, opts...)
	_approval.approvalDo.UseModel(&model.Approval{})

	tableName := _approval.approvalDo.TableName()
	_approval.ALL = field.NewAsterisk(tableName)
	_approval.ID = field.NewInt64(tableName, "id")
	_approval.ClientID = field.NewString(tableName, "client_id")
	_approval.UserID = field.NewInt64(tableName, "user_id")
	_approval.CreatedAt = field.NewTime(tableName, "created_at")
	_approval.UpdatedAt = field.NewTime(tableName, "updated_at")
	_approval.AuthTime = field.NewTime(tableName, "auth_time")

	_approval.fillFieldMap()

	return _approval
}

type approval struct {
	approvalDo

	ALL       field.Asterisk
	ID        field.Int64
	ClientID  field.String
	UserID    field.Int64
	CreatedAt field.Time
	UpdatedAt field.Time
	AuthTime  field.Time

	fieldMap map[string]field.Expr
}

func (a approval) Table(newTableName string) *approval {
	a.approvalDo.UseTable(newTableName)
	return a.updateTableName(newTableName)
}

func (a approval) As(alias string) *approval {
	a.approvalDo.DO = *(a.approvalDo.As(alias).(*gen.DO))
	return a.updateTableName(alias)
}

func (a *approval) updateTableName(table string) *approval {
	a.ALL = field.NewAsterisk(table)
	a.ID = field.NewInt64(table, "id")
	a.ClientID = field.NewString(table, "client_id")
	a.UserID = field.NewInt64(table, "user_id")
	a.CreatedAt = field.NewTime(table, "created_at")
	a.UpdatedAt = field.NewTime(table, "updated_at")
	a.AuthTime = field.NewTime(table, "auth_time")

	a.fillFieldMap()

	return a
}

func (a *approval) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := a.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (a *approval) fillFieldMap() {
	a.fieldMap = make(map[string]field.Expr, 6)
	a.fieldMap["id"] = a.ID
	a.fieldMap["client_id"] = a.ClientID
	a.fieldMap["user_id"] = a.UserID
	a.fieldMap["created_at"] = a.CreatedAt
	a.fieldMap["updated_at"] = a.UpdatedAt
	a.fieldMap["auth_time"] = a.AuthTime
}

func (a approval) clone(db *gorm.DB) approval {
	a.approvalDo.ReplaceConnPool(db.Statement.ConnPool)
	return a
}

func (a approval) replaceDB(db *gorm.DB) approval {
	a.approvalDo.ReplaceDB(db)
	return a
}

type approvalDo struct{ gen.DO }

type IApprovalDo interface {
	gen.SubQuery
	Debug() IApprovalDo
	WithContext(ctx context.Context) IApprovalDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IApprovalDo
	WriteDB() IApprovalDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IApprovalDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IApprovalDo
	Not(conds ...gen.Condition) IApprovalDo
	Or(conds ...gen.Condition) IApprovalDo
	Select(conds ...field.Expr) IApprovalDo
	Where(conds ...gen.Condition) IApprovalDo
	Order(conds ...field.Expr) IApprovalDo
	Distinct(cols ...field.Expr) IApprovalDo
	Omit(cols ...field.Expr) IApprovalDo
	Join(table schema.Tabler, on ...field.Expr) IApprovalDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IApprovalDo
	RightJoin(table schema.Tabler, on ...field.Expr) IApprovalDo
	Group(cols ...field.Expr) IApprovalDo
	Having(conds ...gen.Condition) IApprovalDo
	Limit(limit int) IApprovalDo
	Offset(offset int) IApprovalDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IApprovalDo
	Unscoped() IApprovalDo
	Create(values ...*model.Approval) error
	CreateInBatches(values []*model.Approval, batchSize int) error
	Save(values ...*model.Approval) error
	First() (*model.Approval, error)
	Take() (*model.Approval, error)
	Last() (*model.Approval, error)
	Find() ([]*model.Approval, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Approval, err error)
	FindInBatches(result *[]*model.Approval, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Approval) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IApprovalDo
	Assign(attrs ...field.AssignExpr) IApprovalDo
	Joins(fields ...field.RelationField) IApprovalDo
	Preload(fields ...field.RelationField) IApprovalDo
	FirstOrInit() (*model.Approval, error)
	FirstOrCreate() (*model.Approval, error)
	FindByPage(offset int, limit int) (result []*model.Approval, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IApprovalDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (a approvalDo) Debug() IApprovalDo {
	return a.withDO(a.DO.Debug())
}

func (a approvalDo) WithContext(ctx context.Context) IApprovalDo {
	return a.withDO(a.DO.WithContext(ctx))
}

func (a approvalDo) ReadDB() IApprovalDo {
	return a.Clauses(dbresolver.Read)
}

func (a approvalDo) WriteDB() IApprovalDo {
	return a.Clauses(dbresolver.Write)
}

func (a approvalDo) Session(config *gorm.Session) IApprovalDo {
	return a.withDO(a.DO.Session(config))
}

func (a approvalDo) Clauses(conds ...clause.Expression) IApprovalDo {
	return a.withDO(a.DO.Clauses(conds...))
}

func (a approvalDo) Returning(value interface{}, columns ...string) IApprovalDo {
	return a.withDO(a.DO.Returning(value, columns...))
}

func (a approvalDo) Not(conds ...gen.Condition) IApprovalDo {
	return a.withDO(a.DO.Not(conds...))
}

func (a approvalDo) Or(conds ...gen.Condition) IApprovalDo {
	return a.withDO(a.DO.Or(conds...))
}

func (a approvalDo) Select(conds ...field.Expr) IApprovalDo {
	return a.withDO(a.DO.Select(conds...))
}

func (a approvalDo) Where(conds ...gen.Condition) IApprovalDo {
	return a.withDO(a.DO.Where(conds...))
}

func (a approvalDo) Order(conds ...field.Expr) IApprovalDo {
	return a.withDO(a.DO.Order(conds...))
}

func (a approvalDo) Distinct(cols ...field.Expr) IApprovalDo {
	return a.withDO(a.DO.Distinct(cols...))
}

func (a approvalDo) Omit(cols ...field.Expr) IApprovalDo {
	return a.withDO(a.DO.Omit(cols...))
}

func (a approvalDo) Join(table schema.Tabler, on ...field.Expr) IApprovalDo {
	return a.withDO(a.DO.Join(table, on...))
}

func (a approvalDo) LeftJoin(table schema.Tabler, on ...field.Expr) IApprovalDo {
	return a.withDO(a.DO.LeftJoin(table, on...))
}

func (a approvalDo) RightJoin(table schema.Tabler, on ...field.Expr) IApprovalDo {
	return a.withDO(a.DO.RightJoin(table, on...))
}

func (a approvalDo) Group(cols ...field.Expr) IApprovalDo {
	return a.withDO(a.DO.Group(cols...))
}

func (a approvalDo) Having(conds ...gen.Condition) IApprovalDo {
	return a.withDO(a.DO.Having(conds...))
}

func (a approvalDo) Limit(limit int) IApprovalDo {
	return a.withDO(a.DO.Limit(limit))
}

func (a approvalDo) Offset(offset int) IApprovalDo {
	return a.withDO(a.DO.Offset(offset))
}

func (a approvalDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IApprovalDo {
	return a.withDO(a.DO.Scopes(funcs...))
}

func (a approvalDo) Unscoped() IApprovalDo {
	return a.withDO(a.DO.Unscoped())
}

func (a approvalDo) Create(values ...*model.Approval) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Create(values)
}

func (a approvalDo) CreateInBatches(values []*model.Approval, batchSize int) error {
	return a.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (a approvalDo) Save(values ...*model.Approval) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Save(values)
}

func (a approvalDo) First() (*model.Approval, error) {
	if result, err := a.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Approval), nil
	}
}

func (a approvalDo) Take() (*model.Approval, error) {
	if result, err := a.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Approval), nil
	}
}

func (a approvalDo) Last() (*model.Approval, error) {
	if result, err := a.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Approval), nil
	}
}

func (a approvalDo) Find() ([]*model.Approval, error) {
	result, err := a.DO.Find()
	return result.([]*model.Approval), err
}

func (a approvalDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Approval, err error) {
	buf := make([]*model.Approval, 0, batchSize)
	err = a.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (a approvalDo) FindInBatches(result *[]*model.Approval, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return a.DO.FindInBatches(result, batchSize, fc)
}

func (a approvalDo) Attrs(attrs ...field.AssignExpr) IApprovalDo {
	return a.withDO(a.DO.Attrs(attrs...))
}

func (a approvalDo) Assign(attrs ...field.AssignExpr) IApprovalDo {
	return a.withDO(a.DO.Assign(attrs...))
}

func (a approvalDo) Joins(fields ...field.RelationField) IApprovalDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Joins(_f))
	}
	return &a
}

func (a approvalDo) Preload(fields ...field.RelationField) IApprovalDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Preload(_f))
	}
	return &a
}

func (a approvalDo) FirstOrInit() (*model.Approval, error) {
	if result, err := a.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Approval), nil
	}
}

func (a approvalDo) FirstOrCreate() (*model.Approval, error) {
	if result, err := a.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Approval), nil
	}
}

func (a approvalDo) FindByPage(offset int, limit int) (result []*model.Approval, count int64, err error) {
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

func (a approvalDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = a.Count()
	if err != nil {
		return
	}

	err = a.Offset(offset).Limit(limit).Scan(result)
	return
}

func (a approvalDo) Scan(result interface{}) (err error) {
	return a.DO.Scan(result)
}

func (a approvalDo) Delete(models ...*model.Approval) (result gen.ResultInfo, err error) {
	return a.DO.Delete(models)
}

func (a *approvalDo) withDO(do gen.Dao) *approvalDo {
	a.DO = *do.(*gen.DO)
	return a
}
