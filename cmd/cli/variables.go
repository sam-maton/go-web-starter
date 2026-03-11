package main

var folderPaths = []string{
	"cmd/web",
	"internal/models",
	"internal/validator",
	"ui/html/pages",
	"ui/html/partials",
	"ui/static/css",
	"ui/static/js",
	"ui/static/img",
	"sql/schema",
}

var baseFiles = []string{
	"cmd/web/main.go",
	// "cmd/web/routes.go",
	// "cmd/web/handlers.go",
	// "cmd/web/middleware.go",
	// "cmd/web/helpers.go",
	// "cmd/web/templates.go",
	// "cmd/web/constants.go",

	// "internal/models/errors.go",
	// "internal/models/users.go",
	// "internal/validator/validator.go",

	// "ui/html/base.html",
	// "ui/html/pages/home.html",
	// "ui/html/pages/create.html",
	// "ui/html/pages/view.html",
}

var authFiles = []string{
	"ui/html/pages/login.html",
	"ui/html/pages/signup.html",
}
