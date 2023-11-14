package auth

type CustomClaims struct {
	Permissions []string `json:"permissions"`
	Scope       string   `json:"scope"`
}

type Role struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Permission struct {
	Name        string `json:"permission_name"`
	Description string `json:"description"`
}
