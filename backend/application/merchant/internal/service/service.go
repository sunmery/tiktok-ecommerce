package service

import (
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewInventoryService, NewProductService, NewOrderService, NewAddressService)
