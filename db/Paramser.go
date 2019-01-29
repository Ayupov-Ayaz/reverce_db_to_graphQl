package db

type Paramser interface {
	GetParams() *Params
	GetSupportedVersions() []string
}