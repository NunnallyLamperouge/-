package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// StringArray 用于存储字符串数组
type StringArray []string

// Value 实现 driver.Valuer 接口，将 StringArray 转换为数据库可存储的值
func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}

// Scan 实现 sql.Scanner 接口，将数据库值转换为 StringArray
func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = StringArray{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言为[]byte失败")
	}

	return json.Unmarshal(bytes, a)
}

// StringMatrix 用于存储二维字符串数组
type StringMatrix [][]string

// Value 实现 driver.Valuer 接口，将 StringMatrix 转换为数据库可存储的值
func (m StringMatrix) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}
	return json.Marshal(m)
}

// Scan 实现 sql.Scanner 接口，将数据库值转换为 StringMatrix
func (m *StringMatrix) Scan(value interface{}) error {
	if value == nil {
		*m = StringMatrix{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言为[]byte失败")
	}

	return json.Unmarshal(bytes, m)
}

// IntArray 用于存储整数数组
type IntArray []int

// Value 实现 driver.Valuer 接口，将 IntArray 转换为数据库可存储的值
func (a IntArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}

// Scan 实现 sql.Scanner 接口，将数据库值转换为 IntArray
func (a *IntArray) Scan(value interface{}) error {
	if value == nil {
		*a = IntArray{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言为[]byte失败")
	}

	return json.Unmarshal(bytes, a)
}

// TableFieldArray 用于存储TableField数组
type TableFieldArray []TableField

// Value 实现 driver.Valuer 接口，将 TableFieldArray 转换为数据库可存储的值
func (a TableFieldArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}

// Scan 实现 sql.Scanner 接口，将数据库值转换为 TableFieldArray
func (a *TableFieldArray) Scan(value interface{}) error {
	if value == nil {
		*a = TableFieldArray{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言为[]byte失败")
	}

	return json.Unmarshal(bytes, a)
}

// TocItemArray 用于存储TocItem数组
type TocItemArray []TocItem

// Value 实现 driver.Valuer 接口，将 TocItemArray 转换为数据库可存储的值
func (a TocItemArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}

// Scan 实现 sql.Scanner 接口，将数据库值转换为 TocItemArray
func (a *TocItemArray) Scan(value interface{}) error {
	if value == nil {
		*a = TocItemArray{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言为[]byte失败")
	}

	return json.Unmarshal(bytes, a)
}

// ListItemArray 用于存储ListItem数组
type ListItemArray []ListItem

// Value 实现 driver.Valuer 接口，将 ListItemArray 转换为数据库可存储的值
func (a ListItemArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}

// Scan 实现 sql.Scanner 接口，将数据库值转换为 ListItemArray
func (a *ListItemArray) Scan(value interface{}) error {
	if value == nil {
		*a = ListItemArray{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言为[]byte失败")
	}

	return json.Unmarshal(bytes, a)
}

// ReferenceItemArray 用于存储ReferenceItem数组
type ReferenceItemArray []ReferenceItem

// Value 实现 driver.Valuer 接口，将 ReferenceItemArray 转换为数据库可存储的值
func (a ReferenceItemArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}

// Scan 实现 sql.Scanner 接口，将数据库值转换为 ReferenceItemArray
func (a *ReferenceItemArray) Scan(value interface{}) error {
	if value == nil {
		*a = ReferenceItemArray{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言为[]byte失败")
	}

	return json.Unmarshal(bytes, a)
}

// ScheduleItemArray 用于存储ScheduleItem数组
type ScheduleItemArray []ScheduleItem

// Value 实现 driver.Valuer 接口，将 ScheduleItemArray 转换为数据库可存储的值
func (a ScheduleItemArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}

// Scan 实现 sql.Scanner 接口，将数据库值转换为 ScheduleItemArray
func (a *ScheduleItemArray) Scan(value interface{}) error {
	if value == nil {
		*a = ScheduleItemArray{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言为[]byte失败")
	}

	return json.Unmarshal(bytes, a)
}

// CoverTableFieldArray 用于存储CoverTableField数组
type CoverTableFieldArray []CoverTableField

// Value 实现 driver.Valuer 接口，将 CoverTableFieldArray 转换为数据库可存储的值
func (a CoverTableFieldArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}

// Scan 实现 sql.Scanner 接口，将数据库值转换为 CoverTableFieldArray
func (a *CoverTableFieldArray) Scan(value interface{}) error {
	if value == nil {
		*a = CoverTableFieldArray{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言为[]byte失败")
	}

	return json.Unmarshal(bytes, a)
}
