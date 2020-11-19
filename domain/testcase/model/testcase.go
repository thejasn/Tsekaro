package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

/*
DB Table Details
-------------------------------------


CREATE TABLE `testcase` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `expected` blob DEFAULT NULL,
  `actual` text COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `operation` enum('EQUAL','NOT_EQUAL') COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `status` tinyint(4) NOT NULL DEFAULT 1,
  `flow_id` int(11) NOT NULL,
  `test_case_id` int(11) NOT NULL,
  `scheme` enum('https','http') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'http',
  `host` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `port` int(11) NOT NULL DEFAULT 8080,
  `headers` blob DEFAULT NULL,
  `method` text COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `path` text COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '/',
  `body` text COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `mapping_test_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `testcase_UN` (`flow_id`,`test_case_id`),
  CONSTRAINT `testcase_FK` FOREIGN KEY (`flow_id`) REFERENCES `flow` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "id": 14}
*/

type JSON json.RawMessage

// Testcase struct is a row record of the testcase table in the tester database
type Testcase struct {
	ID            int         `gorm:"AUTO_INCREMENT;column:id;type:INT;primary_key" json:"id"`      //[ 0] id                                             int                  null: false  primary: true   auto: true   col: int             len: -1      default: []
	Name          string      `gorm:"column:name;type:TEXT;size:65535;" json:"name"`                //[ 1] name                                           text(65535)          null: false  primary: false  auto: false  col: text            len: 65535   default: []
	Expected      JSON        `gorm:"column:expected;" json:"expected"`                             //[ 2] expected                                       blob                 null: true   primary: false  auto: false  col: blob            len: -1      default: [NULL]
	Actual        null.String `gorm:"column:actual;type:TEXT;size:65535;" json:"actual"`            //[ 3] actual                                         text(65535)          null: true   primary: false  auto: false  col: text            len: 65535   default: [NULL]
	Operation     string      `gorm:"column:operation;type:CHAR;size:9;" json:"operation"`          //[ 4] operation                                      char(9)              null: false  primary: false  auto: false  col: char            len: 9       default: []
	CreatedAt     time.Time   `gorm:"column:created_at;type:DATETIME;" json:"created_at"`           //[ 5] created_at                                     datetime             null: false  primary: false  auto: false  col: datetime        len: -1      default: [current_timestamp()]
	UpdatedAt     time.Time   `gorm:"column:updated_at;type:DATETIME;" json:"updated_at"`           //[ 6] updated_at                                     datetime             null: false  primary: false  auto: false  col: datetime        len: -1      default: [current_timestamp()]
	Status        int         `gorm:"column:status;type:TINYINT;default:1;" json:"status"`          //[ 7] status                                         tinyint              null: false  primary: false  auto: false  col: tinyint         len: -1      default: [1]
	FlowID        int         `gorm:"column:flow_id;type:INT;" json:"flow_id"`                      //[ 8] flow_id                                        int                  null: false  primary: false  auto: false  col: int             len: -1      default: []
	TestCaseID    int         `gorm:"column:test_case_id;type:INT;" json:"test_case_id"`            //[ 9] test_case_id                                   int                  null: false  primary: false  auto: false  col: int             len: -1      default: []
	Scheme        string      `gorm:"column:scheme;type:CHAR;size:5;default:'http';" json:"scheme"` //[10] scheme                                         char(5)              null: false  primary: false  auto: false  col: char            len: 5       default: ['http']
	Host          string      `gorm:"column:host;type:TEXT;size:65535;" json:"host"`                //[11] host                                           text(65535)          null: false  primary: false  auto: false  col: text            len: 65535   default: []
	Port          int         `gorm:"column:port;type:INT;default:8080;" json:"port"`               //[12] port                                           int                  null: false  primary: false  auto: false  col: int             len: -1      default: [8080]
	Headers       []byte      `gorm:"column:headers;type:BLOB;" json:"headers"`                     //[13] headers                                        blob                 null: true   primary: false  auto: false  col: blob            len: -1      default: [NULL]
	Method        null.String `gorm:"column:method;type:TEXT;size:65535;" json:"method"`            //[14] method                                         text(65535)          null: true   primary: false  auto: false  col: text            len: 65535   default: [NULL]
	Path          string      `gorm:"column:path;type:TEXT;size:65535;default:'/';" json:"path"`    //[15] path                                           text(65535)          null: false  primary: false  auto: false  col: text            len: 65535   default: ['/']
	Body          null.String `gorm:"column:body;type:TEXT;size:65535;" json:"body"`                //[16] body                                           text(65535)          null: true   primary: false  auto: false  col: text            len: 65535   default: [NULL]
	MappingTestID null.Int    `gorm:"column:mapping_test_id;type:INT;" json:"mapping_test_id"`      //[17] mapping_test_id                                int                  null: true   primary: false  auto: false  col: int             len: -1      default: [NULL]
	API           string      `gorm:"column:api;type:CHAR;size:4;default:'REST';" json:"api"`       //[18] api                                            char(4)              null: false  primary: false  auto: false  col: char            len: 4       default: ['REST']
}

type Result struct {
	Data interface{}
}

func (JSON) GormDataType() string {
	return "json"
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

// Value return json value, implement driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

func (r *JSON) UnmarshalJSON(b []byte) error {
	result := json.RawMessage{}
	err := json.Unmarshal(b, &result)
	*r = JSON(result)
	return err
}

func (r JSON) MarshalJSON() ([]byte, error) {
	if len(r) == 0 {
		return nil, nil
	}
	return json.RawMessage(r).MarshalJSON()
}

// TableName sets the insert table name for this struct type
func (t *Testcase) TableName() string {
	return "testcase"
}
