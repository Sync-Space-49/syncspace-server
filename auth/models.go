package auth

type CustomClaims struct {
	Scope string `json:"scope"`
}

type Role struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Permission struct {
	PermissionName           string `json:"permission_name"`
	Description              string `json:"description"`
	ResourceServerName       string `json:"resource_server_name"`
	ResourceServerIdentifier string `json:"resource_server_identifier"`
	Sources                  []struct {
		SourceId   string `json:"source_id"`
		SourceName string `json:"source_name"`
		SourceType string `json:"source_type"`
	} `json:"sources"`
}
