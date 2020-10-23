// +build server

package runtime

// AppType returns the application type, EG: desktop
func (r *system) AppType() string {
	return "server"
}
