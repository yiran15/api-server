package store_test

import (
	"github.com/yiran15/api-server/base/conf"
	"github.com/yiran15/api-server/base/data"
	"github.com/yiran15/api-server/base/log"
	"github.com/yiran15/api-server/model"
	"github.com/yiran15/api-server/store"
)

var (
	userRepo store.Repository[model.User]
	// roleRepo  store.Repository[model.Role]
	txManager store.TxManagerInterface
)

func init() {
	conf.LoadConfig("../../config.yaml")
	db, _, err := data.NewDB()
	if err != nil {
		panic(err)
	}
	log.NewLogger()
	provider := store.NewDBProvider(db)
	userRepo = store.NewRepository[model.User](provider)
	// roleRepo = store.NewRepository[model.Role](provider)
	txManager = store.NewTxManager(db)
}
