package cc

type ConfigCenter interface {
	GetConfigClient() interface{}
}