package repo

import (
	"context"

	"github.com/smallnest/gen/dbmeta"
	"github.com/thejasn/tester/cerrors"
	"github.com/thejasn/tester/domain/testcase/model"
	"gorm.io/gorm"
)

type Testcase interface {
	GetAll(ctx context.Context, page, pagesize int64, order string) ([]*model.Testcase, int64, error)
	Get(context.Context, int) (model.Testcase, error)
	Add(context.Context, *model.Testcase) (*model.Testcase, int64, error)
	Update(context.Context, int, *model.Testcase) (*model.Testcase, int64, error)
	Delete(context.Context, int) (int64, error)
	GetAllOrderedWhere(context.Context, map[string]interface{}) ([]*model.Testcase, int64, error)
}

func NewTestcaseRepo(db *gorm.DB) Testcase {
	return testcase{
		DB: db,
	}
}

type testcase struct {
	DB *gorm.DB
}

// GetAllTestcases is a function to get a slice of record(s) from testcase table in the tester database
// params - page     - page requested (defaults to 0)
// params - pagesize - number of records in a page  (defaults to 20)
// params - order    - db sort order column
// error - ErrNotFound, db Find error
func (t testcase) GetAll(ctx context.Context, page, pagesize int64, order string) (testcases []*model.Testcase, totalRows int64, err error) {

	testcases = []*model.Testcase{}

	testcasesOrm := t.DB.Model(&model.Testcase{})
	testcasesOrm.Count(&totalRows)

	if page > 0 {
		offset := (page - 1) * pagesize
		testcasesOrm = testcasesOrm.Offset(int(offset)).Limit(int(pagesize))
	} else {
		testcasesOrm = testcasesOrm.Limit(int(pagesize))
	}

	if order != "" {
		testcasesOrm = testcasesOrm.Order(order)
	}

	if err = testcasesOrm.Find(&testcases).Error; err != nil {
		err = cerrors.ErrNotFound
		return nil, -1, err
	}

	return testcases, totalRows, nil
}

// GetTestcase is a function to get a single record to testcase table in the tester database
// error - ErrNotFound, db Find error
func (t testcase) Get(ctx context.Context, id int) (record model.Testcase, err error) {
	if err = t.DB.Where("id = ?", id).First(&record).Error; err != nil {
		err = cerrors.ErrNotFound
		return record, err
	}
	return record, nil
}

// AddTestcase is a function to add a single record to testcase table in the tester database
// error - ErrInsertFailed, db save call failed
func (t testcase) Add(ctx context.Context, testcase *model.Testcase) (result *model.Testcase, RowsAffected int64, err error) {
	db := t.DB.Save(testcase)
	if err = db.Error; err != nil {
		return nil, -1, cerrors.ErrInsertFailed
	}

	return testcase, db.RowsAffected, nil
}

// UpdateTestcase is a function to update a single record from testcase table in the tester database
// error - ErrNotFound, db record for id not found
// error - ErrUpdateFailed, db meta data copy failed or db.Save call failed
func (t testcase) Update(ctx context.Context, id int, updated *model.Testcase) (result *model.Testcase, RowsAffected int64, err error) {

	result = &model.Testcase{}
	db := t.DB.First(result, id)
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

// DeleteTestcase is a function to delete a single record from testcase table in the tester database
// error - ErrNotFound, db Find error
// error - ErrDeleteFailed, db Delete failed error
func (t testcase) Delete(ctx context.Context, id int) (rowsAffected int64, err error) {

	testcase := &model.Testcase{}
	db := t.DB.First(testcase, id)
	if db.Error != nil {
		return -1, cerrors.ErrNotFound
	}

	db = db.Delete(testcase)
	if err = db.Error; err != nil {
		return -1, cerrors.ErrDeleteFailed
	}

	return db.RowsAffected, nil
}

func (t testcase) GetAllOrderedWhere(ctx context.Context, conditions map[string]interface{}) (testcases []*model.Testcase, totalRows int64, err error) {
	testcases = []*model.Testcase{}

	testcasesOrm := t.DB.Model(&model.Testcase{})
	testcasesOrm.Count(&totalRows)

	testcasesOrm = testcasesOrm.Order("test_case_id")

	if err = testcasesOrm.Where(conditions).Find(&testcases).Error; err != nil {
		err = cerrors.ErrNotFound
		return nil, -1, err
	}

	return testcases, totalRows, nil
}
