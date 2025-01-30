package common

type Status string

const (
	Travelling   Status = "Travelling"
	Transferring        = "Transferring"
	Unused              = "Unused"
	Emergency           = "Emergency"
)

type Train struct {
	Name string `json:"name"`
	// meters / second
	TopSpeed int `json:"top_speed"`
	Coordinates
	Status `json:"status"`
}
