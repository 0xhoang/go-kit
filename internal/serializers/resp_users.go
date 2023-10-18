package serializers

type UserByEmailResp struct {
	ID              uint   `json:"ID"`
	UserName        string `json:"UserName"`
	FirstName       string `json:"FirstName"`
	LastName        string `json:"LastName"`
	Email           string `json:"Email"`
	Address         string `json:"Address"`
	IsActive        bool   `json:"IsActive"`
	IsVerifiedEmail bool   `json:"IsVerifiedEmail"`
}

type UserLoginResp struct {
	Token   string `json:"Token"`
	Expired string `json:"Expired"`
}
