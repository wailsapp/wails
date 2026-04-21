package commands

func Delve(options *DevOptions) error {
	options.UseDelve = true
	return Dev(options)
}
