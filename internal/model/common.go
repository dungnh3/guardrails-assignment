package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// JSON is a customized data type. The customized data type has to implement
// the Scanner and Valuer interfaces, so GORM knowns to how to receive/save
// it into the database
type JSON json.RawMessage

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

func (JSON) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql", "sqlite":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

func (j JSON) ConvertToStructPb() (*structpb.Struct, error) {
	if j == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if err := json.Unmarshal(j, &m); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unmarshal json failed, value: %v", string(j)))
	}

	data, err := structpb.NewStruct(m)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("convert to struct pb failed, value: %v", m))
	}
	return data, nil
}
