package serializers

import (
	"time"
)

type ProjectsList struct {
	Count         int             `json:"count"`
	ProjectsValue []ProjectsValue `json:"value"`
}

type ProjectsValue struct {
	Id             string    `json:"id"`
	Url            string    `json:"url"`
	Name           string    `json:"name"`
	State          string    `json:"state"`
	Revision       int       `json:"revision"`
	Visibility     string    `json:"visibility"`
	LastUpdateTime time.Time `json:"lastUpdateTime"`
}
