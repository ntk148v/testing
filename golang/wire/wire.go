//+build wireinject
package main

import (
	"github.com/google/wire"
)

func CreateConcatService() *ConcatService {
	wire.Build(
		wire.Struct(new(Logger), "*"),
		NewHttpClient,
		NewConcatService,
	)
	return nil
}
