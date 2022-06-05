package proto

import "project/cms/internal/acl"

type AdminToken struct {
	ID        int           `json:"id"`
	Username  string        `json:"username"`
	Authority acl.Authority `json:"authority"`
}

type Pagination struct {
	Limit  int
	Offset int
}

type DateRange struct {
	Begin string
	End   string
}
