package trx

import (
	"context"
	"database/sql"
	"github.com/maczh/mgin/db"
	"github.com/maczh/mgin/logs"
	"github.com/maczh/mgin/models"
	"gorm.io/gorm"
	"numbers-allocator/errcode"
	"numbers-allocator/multidb"
	"runtime/debug"
)

type txKey struct{}

func Transaction[T any](ctx context.Context, t *T, handler func(ctx context.Context, t *T) models.Result[any]) (ret models.Result[any]) {

	if _, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		//已经有事务，直接执行
		ret = handler(ctx, t)
		return
	}

	//还没有开启事务，开启事务
	conn, err := db.Mysql.GetConnection(multidb.GetPartnerId())
	if err != nil {
		return errcode.DbConnectErr.MGError()
	}

	defer func() {
		if err := recover(); err != nil {
			conn.Rollback()
			logs.Error("panic error: {} stack: {}", err, string(debug.Stack()))
			ret = errcode.SystemError.MGError()
		}
	}()
	opt := sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	}

	conn = conn.Begin(&opt)
	goCtx := context.WithValue(ctx, txKey{}, conn)
	ret = handler(goCtx, t)

	if ret.Status != 1 {
		conn.Rollback()
		return ret
	}
	conn.Commit()

	return ret
}

func NoTransaction[T any](ctx context.Context, t *T, handler func(ctx context.Context, t *T) models.Result[any]) models.Result[any] {
	return handler(ctx, t)
}

//
//func NoTransaction(ginCtx *gin.Context, handler func(ctx context.Context, c *gin.Context) models.Result[any]) models.Result[any] {
//	backCtx := context.Background()
//	return handler(backCtx, ginCtx)
//}

func ExtractDb(ctx context.Context) (*gorm.DB, error) {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx, nil
	}

	conn, err := db.Mysql.GetConnection(multidb.GetPartnerId())
	if err != nil {
		logs.Error("获取数据库连接出错 {}", err.Error())
		return nil, errcode.DbConnectErr.Error()
	}
	return conn, nil
}
