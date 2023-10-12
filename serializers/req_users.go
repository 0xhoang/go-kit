package serializers

type AuthByEmailReq struct {
	Email    string `json:"Email" binding:"required"`
	Password string `json:"Password" binding:"required"`
	IP       string
}
