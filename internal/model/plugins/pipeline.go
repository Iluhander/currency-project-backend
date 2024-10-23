package plugins

type Pipeline struct {
	Auth []*Plugin
	Payment []*Plugin
	Statistics []*Plugin
}
