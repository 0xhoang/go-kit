package common

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	ENV_DEVELOPMENT = "development"
	ENV_PRODUCTION  = "production"
)

const (
	CustomerKey = "customer"
)
