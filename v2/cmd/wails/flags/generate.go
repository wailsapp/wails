package flags

type GenerateModule struct {
	Common
	Tags      string `description:"Build tags to pass to Go compiler. Must be quoted. Space or comma (but not both) separated"`
	Verbosity int    `name:"v" description:"Verbosity level (0 = quiet, 1 = normal, 2 = verbose)"`
}

type GenerateTemplate struct {
	Common
	Name     string `description:"Name of the template to generate"`
	Frontend string `description:"Frontend to use for the template"`
	Quiet    bool   `description:"Suppress output"`
}
