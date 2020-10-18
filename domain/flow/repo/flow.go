package repo

import (
	"context"

	"github.com/smallnest/gen/dbmeta"
	"github.com/thejasn/tester/cerrors"
	"github.com/thejasn/tester/domain/flow/model"
	"gorm.io/gorm"
)

type Flow interface {
	GetAll(ctx context.Context, page, pagesize int64, order string) ([]*model.Flow, int64, error)
	Get(context.Context, int) (model.Flow, error)
	Add(context.Context, *model.Flow) (*model.Flow, int64, error)
	Update(context.Context, int, *model.Flow) (*model.Flow, int64, error)
	Delete(context.Context, int) (int64, error)
}

func NewFlowRepo(db *gorm.DB) Flow {
	return flow{
		DB: db,
	}
}

type flow struct {
	DB *gorm.DB
}

// GetAll is a function to get a slice of record(s) from flow table in the tester database
// params - page     - page requested (defaults to 0)
// params - pagesize - number of records in a page  (defaults to 20)
// params - order    - db sort order column
// error - ErrNotFound, db Find error
func (f flow) GetAll(ctx context.Context, page, pagesize int64, order string) (flows []*model.Flow, totalRows int64, err error) {

	flows = []*model.Flow{}

	flowsOrm := f.DB.Model(&model.Flow{})
	flowsOrm.Count(&totalRows)

	if page > 0 {
		offset := (page - 1) * pagesize
		flowsOrm = flowsOrm.Offset(int(offset)).Limit(int(pagesize))
	} else {
		flowsOrm = flowsOrm.Limit(int(pagesize))
	}

	if order != "" {
		flowsOrm = flowsOrm.Order(order)
	}

	if err = flowsOrm.Find(&flows).Error; err != nil {
		err = cerrors.ErrNotFound
		return nil, -1, err
	}

	return flows, totalRows, nil
}

// GetFlow is a function to get a single record to flow table in the tester database
// error - ErrNotFound, db Find error
func (f flow) Get(ctx context.Context, id int) (record model.Flow, err error) {
	if err = f.DB.First(&record, id).Error; err != nil {
		err = cerrors.ErrNotFound
		return record, err
	}

	return record, nil
}

// AddFlow is a function to add a single record to flow table in the tester database
// error - ErrInsertFailed, db save call failed
func (f flow) Add(ctx context.Context, flow *model.Flow) (result *model.Flow, RowsAffected int64, err error) {
	db := f.DB.Save(flow)
	if err = db.Error; err != nil {
		return nil, -1, cerrors.ErrInsertFailed
	}

	return flow, db.RowsAffected, nil
}

// UpdateFlow is a function to update a single record from flow table in the tester database
// error - ErrNotFound, db record for id not found
// error - ErrUpdateFailed, db meta data copy failed or db.Save call failed
func (f flow) Update(ctx context.Context, id int, updated *model.Flow) (result *model.Flow, RowsAffected int64, err error) {

	result = &model.Flow{}
	db := f.DB.First(result, id)
	if err = db.Error; err != nil {
		return nil, -1, cerrors.ErrNotFound
	}

	if err = dbmeta.Copy(result, updated); err != nil {
		return nil, -1, cerrors.ErrUpdateFailed
	}

	db = db.Save(result)
	if err = db.Error; err != nil {
		return nil, -1, cerrors.ErrUpdateFailed
	}

	return result, db.RowsAffected, nil
}

// DeleteFlow is a function to delete a single record from flow table in the tester database
// error - ErrNotFound, db Find error
// error - ErrDeleteFailed, db Delete failed error
func (f flow) Delete(ctx context.Context, id int) (rowsAffected int64, err error) {

	flow := &model.Flow{}
	db := f.DB.First(flow, id)
	if db.Error != nil {
		return -1, cerrors.ErrNotFound
	}

	db = db.Delete(flow)
	if err = db.Error; err != nil {
		return -1, cerrors.ErrDeleteFailed
	}

	return db.RowsAffected, nil
}
