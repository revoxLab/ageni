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

	"github.com/readonme/open-studio/dal/model"
)

func newBotDraft(db *gorm.DB, opts ...gen.DOOption) botDraft {
	_botDraft := botDraft{}

	_botDraft.botDraftDo.UseDB(db, opts...)
	_botDraft.botDraftDo.UseModel(&model.BotDraft{})

	tableName := _botDraft.botDraftDo.TableName()
	_botDraft.ALL = field.NewAsterisk(tableName)
	_botDraft.ID = field.NewInt64(tableName, "id")
	_botDraft.BotID = field.NewInt64(tableName, "bot_id")
	_botDraft.CreatorID = field.NewInt64(tableName, "creator_id")
	_botDraft.Status = field.NewInt32(tableName, "status")
	_botDraft.Prompt = field.NewString(tableName, "prompt")
	_botDraft.Plugins = field.NewString(tableName, "plugins")
	_botDraft.WelcomeMsg = field.NewString(tableName, "welcome_msg")
	_botDraft.GuideInfo = field.NewString(tableName, "guide_info")
	_botDraft.ModelSettings = field.NewString(tableName, "model_settings")
	_botDraft.DebugConversationID = field.NewInt64(tableName, "debug_conversation_id")
	_botDraft.CreatedAt = field.NewTime(tableName, "created_at")
	_botDraft.UpdatedAt = field.NewTime(tableName, "updated_at")

	_botDraft.fillFieldMap()

	return _botDraft
}

type botDraft struct {
	botDraftDo botDraftDo

	ALL                 field.Asterisk
	ID                  field.Int64
	BotID               field.Int64
	CreatorID           field.Int64
	Status              field.Int32 // 0-default 1-publish
	Prompt              field.String
	Plugins             field.String
	WelcomeMsg          field.String
	GuideInfo           field.String
	ModelSettings       field.String
	DebugConversationID field.Int64
	CreatedAt           field.Time
	UpdatedAt           field.Time

	fieldMap map[string]field.Expr
}

func (b botDraft) Table(newTableName string) *botDraft {
	b.botDraftDo.UseTable(newTableName)
	return b.updateTableName(newTableName)
}

func (b botDraft) As(alias string) *botDraft {
	b.botDraftDo.DO = *(b.botDraftDo.As(alias).(*gen.DO))
	return b.updateTableName(alias)
}

func (b *botDraft) updateTableName(table string) *botDraft {
	b.ALL = field.NewAsterisk(table)
	b.ID = field.NewInt64(table, "id")
	b.BotID = field.NewInt64(table, "bot_id")
	b.CreatorID = field.NewInt64(table, "creator_id")
	b.Status = field.NewInt32(table, "status")
	b.Prompt = field.NewString(table, "prompt")
	b.Plugins = field.NewString(table, "plugins")
	b.WelcomeMsg = field.NewString(table, "welcome_msg")
	b.GuideInfo = field.NewString(table, "guide_info")
	b.ModelSettings = field.NewString(table, "model_settings")
	b.DebugConversationID = field.NewInt64(table, "debug_conversation_id")
	b.CreatedAt = field.NewTime(table, "created_at")
	b.UpdatedAt = field.NewTime(table, "updated_at")

	b.fillFieldMap()

	return b
}

func (b *botDraft) WithContext(ctx context.Context) *botDraftDo { return b.botDraftDo.WithContext(ctx) }

func (b botDraft) TableName() string { return b.botDraftDo.TableName() }

func (b botDraft) Alias() string { return b.botDraftDo.Alias() }

func (b botDraft) Columns(cols ...field.Expr) gen.Columns { return b.botDraftDo.Columns(cols...) }

func (b *botDraft) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := b.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (b *botDraft) fillFieldMap() {
	b.fieldMap = make(map[string]field.Expr, 12)
	b.fieldMap["id"] = b.ID
	b.fieldMap["bot_id"] = b.BotID
	b.fieldMap["creator_id"] = b.CreatorID
	b.fieldMap["status"] = b.Status
	b.fieldMap["prompt"] = b.Prompt
	b.fieldMap["plugins"] = b.Plugins
	b.fieldMap["welcome_msg"] = b.WelcomeMsg
	b.fieldMap["guide_info"] = b.GuideInfo
	b.fieldMap["model_settings"] = b.ModelSettings
	b.fieldMap["debug_conversation_id"] = b.DebugConversationID
	b.fieldMap["created_at"] = b.CreatedAt
	b.fieldMap["updated_at"] = b.UpdatedAt
}

func (b botDraft) clone(db *gorm.DB) botDraft {
	b.botDraftDo.ReplaceConnPool(db.Statement.ConnPool)
	return b
}

func (b botDraft) replaceDB(db *gorm.DB) botDraft {
	b.botDraftDo.ReplaceDB(db)
	return b
}

type botDraftDo struct{ gen.DO }

func (b botDraftDo) Debug() *botDraftDo {
	return b.withDO(b.DO.Debug())
}

func (b botDraftDo) WithContext(ctx context.Context) *botDraftDo {
	return b.withDO(b.DO.WithContext(ctx))
}

func (b botDraftDo) ReadDB() *botDraftDo {
	return b.Clauses(dbresolver.Read)
}

func (b botDraftDo) WriteDB() *botDraftDo {
	return b.Clauses(dbresolver.Write)
}

func (b botDraftDo) Session(config *gorm.Session) *botDraftDo {
	return b.withDO(b.DO.Session(config))
}

func (b botDraftDo) Clauses(conds ...clause.Expression) *botDraftDo {
	return b.withDO(b.DO.Clauses(conds...))
}

func (b botDraftDo) Returning(value interface{}, columns ...string) *botDraftDo {
	return b.withDO(b.DO.Returning(value, columns...))
}

func (b botDraftDo) Not(conds ...gen.Condition) *botDraftDo {
	return b.withDO(b.DO.Not(conds...))
}

func (b botDraftDo) Or(conds ...gen.Condition) *botDraftDo {
	return b.withDO(b.DO.Or(conds...))
}

func (b botDraftDo) Select(conds ...field.Expr) *botDraftDo {
	return b.withDO(b.DO.Select(conds...))
}

func (b botDraftDo) Where(conds ...gen.Condition) *botDraftDo {
	return b.withDO(b.DO.Where(conds...))
}

func (b botDraftDo) Order(conds ...field.Expr) *botDraftDo {
	return b.withDO(b.DO.Order(conds...))
}

func (b botDraftDo) Distinct(cols ...field.Expr) *botDraftDo {
	return b.withDO(b.DO.Distinct(cols...))
}

func (b botDraftDo) Omit(cols ...field.Expr) *botDraftDo {
	return b.withDO(b.DO.Omit(cols...))
}

func (b botDraftDo) Join(table schema.Tabler, on ...field.Expr) *botDraftDo {
	return b.withDO(b.DO.Join(table, on...))
}

func (b botDraftDo) LeftJoin(table schema.Tabler, on ...field.Expr) *botDraftDo {
	return b.withDO(b.DO.LeftJoin(table, on...))
}

func (b botDraftDo) RightJoin(table schema.Tabler, on ...field.Expr) *botDraftDo {
	return b.withDO(b.DO.RightJoin(table, on...))
}

func (b botDraftDo) Group(cols ...field.Expr) *botDraftDo {
	return b.withDO(b.DO.Group(cols...))
}

func (b botDraftDo) Having(conds ...gen.Condition) *botDraftDo {
	return b.withDO(b.DO.Having(conds...))
}

func (b botDraftDo) Limit(limit int) *botDraftDo {
	return b.withDO(b.DO.Limit(limit))
}

func (b botDraftDo) Offset(offset int) *botDraftDo {
	return b.withDO(b.DO.Offset(offset))
}

func (b botDraftDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *botDraftDo {
	return b.withDO(b.DO.Scopes(funcs...))
}

func (b botDraftDo) Unscoped() *botDraftDo {
	return b.withDO(b.DO.Unscoped())
}

func (b botDraftDo) Create(values ...*model.BotDraft) error {
	if len(values) == 0 {
		return nil
	}
	return b.DO.Create(values)
}

func (b botDraftDo) CreateInBatches(values []*model.BotDraft, batchSize int) error {
	return b.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (b botDraftDo) Save(values ...*model.BotDraft) error {
	if len(values) == 0 {
		return nil
	}
	return b.DO.Save(values)
}

func (b botDraftDo) First() (*model.BotDraft, error) {
	if result, err := b.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.BotDraft), nil
	}
}

func (b botDraftDo) Take() (*model.BotDraft, error) {
	if result, err := b.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.BotDraft), nil
	}
}

func (b botDraftDo) Last() (*model.BotDraft, error) {
	if result, err := b.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.BotDraft), nil
	}
}

func (b botDraftDo) Find() ([]*model.BotDraft, error) {
	result, err := b.DO.Find()
	return result.([]*model.BotDraft), err
}

func (b botDraftDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.BotDraft, err error) {
	buf := make([]*model.BotDraft, 0, batchSize)
	err = b.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (b botDraftDo) FindInBatches(result *[]*model.BotDraft, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return b.DO.FindInBatches(result, batchSize, fc)
}

func (b botDraftDo) Attrs(attrs ...field.AssignExpr) *botDraftDo {
	return b.withDO(b.DO.Attrs(attrs...))
}

func (b botDraftDo) Assign(attrs ...field.AssignExpr) *botDraftDo {
	return b.withDO(b.DO.Assign(attrs...))
}

func (b botDraftDo) Joins(fields ...field.RelationField) *botDraftDo {
	for _, _f := range fields {
		b = *b.withDO(b.DO.Joins(_f))
	}
	return &b
}

func (b botDraftDo) Preload(fields ...field.RelationField) *botDraftDo {
	for _, _f := range fields {
		b = *b.withDO(b.DO.Preload(_f))
	}
	return &b
}

func (b botDraftDo) FirstOrInit() (*model.BotDraft, error) {
	if result, err := b.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.BotDraft), nil
	}
}

func (b botDraftDo) FirstOrCreate() (*model.BotDraft, error) {
	if result, err := b.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.BotDraft), nil
	}
}

func (b botDraftDo) FindByPage(offset int, limit int) (result []*model.BotDraft, count int64, err error) {
	result, err = b.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = b.Offset(-1).Limit(-1).Count()
	return
}

func (b botDraftDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = b.Count()
	if err != nil {
		return
	}

	err = b.Offset(offset).Limit(limit).Scan(result)
	return
}

func (b botDraftDo) Scan(result interface{}) (err error) {
	return b.DO.Scan(result)
}

func (b botDraftDo) Delete(models ...*model.BotDraft) (result gen.ResultInfo, err error) {
	return b.DO.Delete(models)
}

func (b *botDraftDo) withDO(do gen.Dao) *botDraftDo {
	b.DO = *do.(*gen.DO)
	return b
}
