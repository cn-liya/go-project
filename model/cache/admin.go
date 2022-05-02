package cache

import "strconv"

const (
	keyAdminSSO   = "asso:" // +id
	keyAdminToken = "atk:"  // +token
)

func AdminSSOKey(id int) string { return keyAdminSSO + strconv.Itoa(id) }

func AdminTokenKey(token string) string { return keyAdminToken + token }

type AdminToken struct {
	ID       int      `json:"i"`
	Username string   `json:"u"`
	IsSuper  bool     `json:"s,omitempty"`
	Actions  []string `json:"a,omitempty"`
}
