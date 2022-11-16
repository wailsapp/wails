package flags

type commonFlag struct {
	Flag string
	Description string
}

var TsPrefix = commonFlag{
	Flag: "tsprefix",
	Description: "prefix for generated typescript entities",
}

var TsSuffix = commonFlag{
	Flag: "tssuffix",
	Description: "suffix for generated typescript entities",
}
