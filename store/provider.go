package store

import "github.com/google/wire"

var StoreProviderSet = wire.NewSet(
	wire.Bind(new(DBProviderInterface), new(*DBProvider)),
	NewDBProvider,

	wire.Bind(new(TxManagerInterface), new(*TxManager)),
	NewTxManager,

	NewUserStore,
	NewRoleStore,
	NewApiStore,
	NewCasbinStore,

	wire.Bind(new(CacheStorer), new(*CacheStore)),
	NewCacheStore,
)
