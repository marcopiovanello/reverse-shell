package responses

import "time"

type DirectoryList struct {
	Current string `json:"current"`
	List    []Node `json:"list"`
}

type Node struct {
	IsDir bool      `json:"isDir"`
	Size  int64     `json:"size"`
	Name  string    `json:"name"`
	MTime time.Time `json:"mtime"`
}
