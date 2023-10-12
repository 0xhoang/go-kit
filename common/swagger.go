package common

import "gitlab.com/idolauncher/go-template-kit/docs"

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @description					Description for what is this security definition being used

//access docs: http://localhost:8080/docs/index.html

func SwaggerConfig() {
	docs.SwaggerInfo.Title = "Documentation API"
	docs.SwaggerInfo.Description = "Api Service"
	docs.SwaggerInfo.Host = "https://localhost"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"https", "http"}
}
