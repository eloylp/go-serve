package server

var (
	Name      string
	Version   string
	Build     string
	BuildTime string
)

type Info struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Build     string `json:"build"`
	BuildTime string `json:"build_time"`
}
