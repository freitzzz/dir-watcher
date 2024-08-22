package internal

type Rules struct {
	Watch   []WatchDir `json:"watch"`
	Move    []MoveDir  `json:"move"`
	Unknown Path       `json:"unknown"`
}

type Path string

type WatchDir Path

type MoveDir struct {
	Path Path     `json:"path"`
	Ext  []string `json:"ext"`
}

// Example: "png" -> "/home/user/Pictures"
type FileExtToDirMap map[string]string
