package server

var (
	Name        string
	Version     string
	Build       string
	BuildTime   string
	Information Info
)

func init() {
	Information = Info{
		Name:      Name,
		Version:   Version,
		Build:     Build,
		BuildTime: BuildTime,
	}
}

type Info struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Build     string `json:"build"`
	BuildTime string `json:"build_time"`
}
