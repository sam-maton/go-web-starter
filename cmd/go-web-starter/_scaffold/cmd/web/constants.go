package main

// Session keys
var (
[[- if .EnableAuth]]
AUTH_USER_KEY = "authenticatedUserID"
[[- end]]
FLASH_KEY = "flash"
)

// Form validation messages
var (
NOT_BLANK_ERROR = "This field cannot be blank"
)
[[if .EnableAuth]]
// Page template file names
var (
LOGIN_PAGE     = "login.html"
SIGNUP_PAGE    = "signup.html"
DASHBOARD_PAGE = "dashboard.html"
)
[[end -]]
