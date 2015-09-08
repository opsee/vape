package api

type messages struct {
	Ok                         string
	InternalServerError        string
	BadRequest                 string
	IdRequired                 string
	EmailRequired              string
	NameRequired               string
	PasswordRequired           string
	TokenRequired              string
	CustomerIdRequired         string
	AdminRequired              string
	UserOrAdminRequired        string
	CredentialsMismatch        string
	BastionCredentialsMismatch string
	CustomerNotAuthorized      string
	InvalidToken               string
	EmailConflict              string
	UserConflict               string
	SignupNotFound             string
	UserNotFound               string
	UserDeleted                string
}

var Messages = &messages{
	Ok:                         "ok",
	InternalServerError:        "an unexpected error happened!",
	BadRequest:                 "malformed request",
	IdRequired:                 "id is required",
	EmailRequired:              "email is required",
	NameRequired:               "name is required",
	PasswordRequired:           "password is required",
	TokenRequired:              "a valid token is required",
	CustomerIdRequired:         "customer id is required",
	AdminRequired:              "an administrator is required to access this resource",
	UserOrAdminRequired:        "an authorized user or administrator is required to access this resource",
	CredentialsMismatch:        "credentials don't match an active user",
	BastionCredentialsMismatch: "credentials don't match an active bastion",
	CustomerNotAuthorized:      "this customer is not authorized",
	InvalidToken:               "token is invalid or expired",
	EmailConflict:              "that email has already been taken",
	UserConflict:               "that user has already been claimed",
	SignupNotFound:             "signup not found",
	UserNotFound:               "user not found",
	UserDeleted:                "user has been deleted",
}
