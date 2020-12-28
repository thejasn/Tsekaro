package model

import (
	"database/sql"
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


CREATE TABLE `flow` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `flow_UK` (`name`) USING HASH
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "id": 86}
*/

// Flow struct is a row record of the flow table in the tester database
type Flow struct {
	ID        int       `gorm:"AUTO_INCREMENT;column:id;type:INT;primary_key" json:"id"` //[ 0] id                                             int                  null: false  primary: true   auto: true   col: int             len: -1      default: []
	Name      string    `gorm:"column:name;type:TEXT;size:65535;" json:"name"`           //[ 1] name                                           text(65535)          null: false  primary: false  auto: false  col: text            len: 65535   default: []
	CreatedAt time.Time `gorm:"column:created_at;type:DATETIME;" json:"created_at"`      //[ 2] created_at                                     datetime             null: false  primary: false  auto: false  col: datetime        len: -1      default: [current_timestamp()]
	UpdatedAt time.Time `gorm:"column:updated_at;type:DATETIME;" json:"updated_at"`      //[ 3] updated_at                                     datetime             null: false  primary: false  auto: false  col: datetime        len: -1      default: [current_timestamp()]

}

// TableName sets the insert table name for this struct type
func (f *Flow) TableName() string {
	return "flow"
}
