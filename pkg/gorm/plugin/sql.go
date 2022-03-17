// Copyright 2022 The imkuqin-zw Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plugin

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"gorm.io/gorm"
	"gorm.io/gorm/utils"
)

func StatementToString(statement *gorm.Statement, containArgs bool) string {
	if containArgs {
		return explainSQL(statement.SQL.String(), statement.Vars...)
	}
	return statement.SQL.String()
}

var convertableTypes = []reflect.Type{reflect.TypeOf(time.Time{}), reflect.TypeOf(false), reflect.TypeOf([]byte{})}

func isPrintable(s []byte) bool {
	for _, r := range s {
		if !unicode.IsPrint(rune(r)) {
			return false
		}
	}
	return true
}

func explainSQL(sql string, avars ...interface{}) string {
	var convertParams func(interface{}, int)
	var vars = make([]string, len(avars))
	var escaper = `"`

	convertParams = func(v interface{}, idx int) {
		switch v := v.(type) {
		case bool:
			vars[idx] = strconv.FormatBool(v)
		case time.Time:
			if v.IsZero() {
				vars[idx] = escaper + "0000-00-00 00:00:00" + escaper
			} else {
				vars[idx] = escaper + v.Format("2006-01-02 15:04:05.999") + escaper
			}
		case *time.Time:
			if v != nil {
				if v.IsZero() {
					vars[idx] = escaper + "0000-00-00 00:00:00" + escaper
				} else {
					vars[idx] = escaper + v.Format("2006-01-02 15:04:05.999") + escaper
				}
			} else {
				vars[idx] = "NULL"
			}
		case fmt.Stringer:
			vars[idx] = escaper + strings.Replace(fmt.Sprintf("%v", v), escaper, "\\"+escaper, -1) + escaper
		case driver.Valuer:
			reflectValue := reflect.ValueOf(v)
			if v != nil && reflectValue.IsValid() && (reflectValue.Kind() == reflect.Ptr && !reflectValue.IsNil()) {
				r, _ := v.Value()
				vars[idx] = fmt.Sprintf("%v", r)
			} else {
				vars[idx] = "NULL"
			}
		case []byte:
			if isPrintable(v) {
				vars[idx] = escaper + strings.Replace(string(v), escaper, "\\"+escaper, -1) + escaper
			} else {
				vars[idx] = escaper + "<binary>" + escaper
			}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			vars[idx] = utils.ToString(v)
		case float64, float32:
			vars[idx] = fmt.Sprintf("%.6f", v)
		case string:
			vars[idx] = escaper + strings.Replace(v, escaper, "\\"+escaper, -1) + escaper
		default:
			rv := reflect.ValueOf(v)
			if v == nil || !rv.IsValid() || rv.Kind() == reflect.Ptr && rv.IsNil() {
				vars[idx] = "NULL"
			} else if valuer, ok := v.(driver.Valuer); ok {
				v, _ = valuer.Value()
				convertParams(v, idx)
			} else if rv.Kind() == reflect.Ptr && !rv.IsZero() {
				convertParams(reflect.Indirect(rv).Interface(), idx)
			} else {
				for _, t := range convertableTypes {
					if rv.Type().ConvertibleTo(t) {
						convertParams(rv.Convert(t).Interface(), idx)
						return
					}
				}
				vars[idx] = escaper + strings.Replace(fmt.Sprint(v), escaper, "\\"+escaper, -1) + escaper
			}
		}
	}

	for idx, v := range avars {
		convertParams(v, idx)
	}

	for _, v := range vars {
		sql = strings.Replace(sql, "?", v, 1)
	}

	return sql
}
