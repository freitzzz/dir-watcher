package internal

type Rules struct {
	Watch   []WatchDir `json:"watch"`
	Move    []MoveDir  `json:"move"`
	Unknown GlobPath   `json:"unknown"`
}

type GlobPath string

type WatchDir GlobPath

type MoveDir struct {
	Path GlobPath `json:"path"`
	Ext  []string `json:"ext"`
}
