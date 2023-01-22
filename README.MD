## Token Auth API

[![GitHub license](https://badgen.net/github/license/Naereen/Strapdown.js)](https://github.com/Naereen/StrapDown.js/blob/master/LICENSE)
[![Go](https://img.shields.io/badge/--00ADD8?logo=go&logoColor=ffffff)](https://golang.org/)
[![Twitter](https://badgen.net/badge/icon/twitter?icon=twitter&label)](https://twitter.com/jvckcallow)


This serves as a basic example of user based CRUD functionality, email and password login & token based authentication through middleware.

The idea is that these API endpoints are consumed by an SPA or similar front-end interface.

Endpoints

 - Register
 - Login
 - Get logged in user
 - Get user by ID
 - Update user
 - Delete user
 - Request password reset
 - Check password reset code valid
 - Do password reset
 - Admin only

Read the [PostMan Collection Docs](https://documenter.getpostman.com/view/5420516/2s8ZDa1gS4) for more information about the API endpoints and passing the token header.

____

The `pw-reset.html` file is the email template used for password reset notifications.

The following data variables can be used in the template by default:
- `{{.ActionUrl}}`
- `{{.Name}}`

[This code block](https://github.com/sjmc11/go-auth-api/blob/main/methods/authMethods.go#L136) dictates the data available in the email template.

____

This GoLang application *mainly* uses the following packages:

- Router - [Chi](https://github.com/go-chi/chi)
- Database - [PGX](github.com/jackc/pgx/v5)
- Mail - [Go Mail](github.com/ainsleyclark/go-mail)


#### TO-DO / Improvements

- Improve method/service initialisation
- Pool DB connections & close DB Conn after events