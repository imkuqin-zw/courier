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

package trace

import (
	"strings"

	"github.com/imkuqin-zw/courier/pkg/gorm/plugin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"gorm.io/gorm"
)

type trace struct {
	addr        string
	user        string
	dbName      string
	containArgs bool
}

func New(user, dbName, addr string, containArgs bool) gorm.Plugin {
	return &trace{
		addr:        addr,
		user:        user,
		dbName:      dbName,
		containArgs: containArgs,
	}
}

func (t *trace) Name() string {
	return "gorm:trace"
}

func (t *trace) Initialize(db *gorm.DB) error {
	t.registerCallbacks(db)
	return nil
}

func (t *trace) tracing(next func(*gorm.DB)) func(*gorm.DB) {
	return func(db *gorm.DB) {
		statement := plugin.StatementToString(db.Statement, t.containArgs)
		op := (strings.Split(statement, " "))[0]
		span, _ := opentracing.StartSpanFromContext(db.Statement.Context, op)
		defer span.Finish()

		// 延迟执行 scope.CombinedConditionSql() 避免sqlVar被重复追加
		next(db)

		ext.DBInstance.Set(span, t.dbName)
		ext.DBType.Set(span, "sql")
		ext.DBUser.Set(span, t.user)
		ext.DBStatement.Set(span, statement)
		ext.SpanKind.Set(span, "client")
		ext.PeerService.Set(span, "mysql")
		ext.PeerAddress.Set(span, t.addr)
		return
	}
}

func (t *trace) registerCallbacks(db *gorm.DB) {
	_ = db.Callback().Create().Replace("gorm:create", t.tracing(db.Callback().Create().Get("gorm:create")))
	_ = db.Callback().Query().Replace("gorm:query", t.tracing(db.Callback().Query().Get("gorm:query")))
	_ = db.Callback().Delete().Replace("gorm:delete", t.tracing(db.Callback().Delete().Get("gorm:delete")))
	_ = db.Callback().Update().Replace("gorm:update", t.tracing(db.Callback().Update().Get("gorm:update")))
	_ = db.Callback().Row().Replace("gorm:row", t.tracing(db.Callback().Row().Get("gorm:row")))
	_ = db.Callback().Raw().Replace("gorm:raw", t.tracing(db.Callback().Raw().Get("gorm:raw")))
}
