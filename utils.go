package gomirai

import (
	jsoniter "github.com/json-iterator/go"
)

var (
	// JSON -
	JSON = jsoniter.ConfigFastest
	// DefaultLogger *logrus.Logger = logrus.New()
)

// func init() {
// 	DefaultLogger.SetFormatter(&logrus.JSONFormatter{})
// }
