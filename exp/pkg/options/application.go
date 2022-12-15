package options

type Application struct {
	Mac *Mac
}

var ApplicationDefaults = &Application{
	Mac: &Mac{
		ActivationPolicy: ActivationPolicyRegular,
	},
}
