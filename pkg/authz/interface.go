package authz

//Authorization interface for all security schemes
type Authorization interface {
	Authorize() bool
}
