package serializers

type UserInfo struct {
	ID        uint   `json:"ID"`
	UserName  string `json:"UserName"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Email     string `json:"Email"`
}
