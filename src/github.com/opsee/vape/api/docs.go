package api

var swaggerJson = `
{
  "basePath": "/%7B%7B.%7D%7D",
  "swagger": "2.0",
  "info": {
    "title": "Vape API",
    "version": "0.0.1",
    "description": "API for user/customer management and authentication"
  },
  "tags": [
    {
      "name": "authenticate",
      "description": "Authentication API"
    },
    {
      "name": "signups",
      "description": "Signup API"
    },
    {
      "name": "users",
      "description": "User API"
    }
  ],
  "paths": {
    "/authenticate/password": {
      "post": {
        "tags": [
          "authenticate"
        ],
        "operationId": "authenticateFromPassword",
        "summary": "Authenticates a user with email and password.",
        "parameters": [
          {
            "in": "body",
            "description": "A user's email",
            "name": "email",
            "required": true,
            "schema": {
              "type": "string",
              "minimum": 0,
              "maximum": 0
            }
          },
          {
            "in": "body",
            "description": "A user's password",
            "name": "password",
            "required": true,
            "schema": {
              "type": "string",
              "minimum": 0,
              "maximum": 0
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Description was not specified",
            "schema": {
              "$ref": "#/definitions/github.com.opsee.vape.api.UserTokenResponse",
              "items": {}
            }
          },
          "401": {
            "description": "Description was not specified"
          }
        }
      }
    },
    "/authenticate/token": {
      "post": {
        "tags": [
          "authenticate"
        ],
        "operationId": "authenticateFromToken",
        "summary": "Authenticates a user by emailing a Bearer token.",
        "parameters": [
          {
            "in": "body",
            "description": "A user's email",
            "name": "email",
            "required": true,
            "schema": {
              "type": "string",
              "minimum": 0,
              "maximum": 0
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Description was not specified",
            "schema": {
              "$ref": "#/definitions/github.com.opsee.vape.api.MessageResponse",
              "items": {}
            }
          },
          "401": {
            "description": "Description was not specified"
          }
        }
      }
    },
    "/authenticate/echo": {
      "get": {
        "tags": [
          "authenticate"
        ],
        "operationId": "echoSession",
        "summary": "Echos a user session given an authentication token.",
        "parameters": [
          {
            "in": "header",
            "description": "The Bearer token",
            "name": "Authorization",
            "required": true,
            "type": "string",
            "minimum": 0,
            "maximum": 0
          }
        ],
        "responses": {
          "200": {
            "description": "Description was not specified",
            "schema": {
              "$ref": "#/definitions/github.com.opsee.vape.model.User",
              "items": {}
            }
          }
        }
      }
    },
    "/signups": {
      "get": {
        "tags": [
          "signups"
        ],
        "operationId": "listSignups",
        "summary": "List all signups.",
        "parameters": [
          {
            "in": "header",
            "description": "The Bearer token - an admin user token is required",
            "name": "Authorization",
            "required": true,
            "type": "string",
            "minimum": 0,
            "maximum": 0
          },
          {
            "in": "query",
            "description": "Pagination - number of records per page",
            "name": "per_page",
            "required": false,
            "type": "integer",
            "format": "int32",
            "minimum": 0,
            "maximum": 0
          },
          {
            "in": "query",
            "description": "Pagination - which page",
            "name": "page",
            "required": false,
            "type": "integer",
            "format": "int32",
            "minimum": 0,
            "maximum": 0
          }
        ],
        "responses": {
          "200": {
            "description": "Description was not specified",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/github.com.opsee.vape.model.Signup"
              }
            }
          },
          "401": {
            "description": "Description was not specified"
          }
        }
      },
      "post": {
        "tags": [
          "signups"
        ],
        "operationId": "createSignup",
        "summary": "Create a new signup.",
        "parameters": [
          {
            "in": "body",
            "description": "The user's name",
            "name": "name",
            "required": true,
            "schema": {
              "type": "string",
              "minimum": 0,
              "maximum": 0
            }
          },
          {
            "in": "body",
            "description": "The user's email",
            "name": "email",
            "required": true,
            "schema": {
              "type": "string",
              "minimum": 0,
              "maximum": 0
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Description was not specified",
            "schema": {
              "$ref": "#/definitions/github.com.opsee.vape.model.Signup",
              "items": {}
            }
          },
          "409": {
            "description": "Email was already used to sign up"
          }
        }
      }
    },
    "/signups/{id}/activate": {
      "put": {
        "tags": [
          "signups"
        ],
        "operationId": "activateSignup",
        "summary": "Sends the activation email for a signup. Can be called multiple times to send multiple emails.",
        "parameters": [
          {
            "in": "path",
            "description": "The signup's id",
            "name": "id",
            "required": true,
            "type": "integer",
            "format": "int32",
            "minimum": 0,
            "maximum": 0
          }
        ],
        "responses": {
          "200": {
            "description": "An object with the claim token used to verify the signup (sent in email)",
            "schema": {
              "$ref": "#/definitions/github.com.opsee.vape.api.SignupActivationResponse",
              "items": {}
            }
          },
          "401": {
            "description": "Description was not specified"
          }
        }
      }
    },
    "/signups/{id}": {
      "get": {
        "tags": [
          "signups"
        ],
        "operationId": "getSignup",
        "summary": "Get a single signup.",
        "parameters": [
          {
            "in": "header",
            "description": "The Bearer token - an admin user token is required",
            "name": "Authorization",
            "required": true,
            "type": "string",
            "minimum": 0,
            "maximum": 0
          },
          {
            "in": "path",
            "description": "The signup id",
            "name": "id",
            "required": true,
            "type": "integer",
            "format": "int32",
            "minimum": 0,
            "maximum": 0
          }
        ],
        "responses": {
          "200": {
            "description": "Description was not specified",
            "schema": {
              "$ref": "#/definitions/github.com.opsee.vape.model.Signup",
              "items": {}
            }
          },
          "401": {
            "description": "Description was not specified"
          }
        }
      }
    },
    "/signups/{id}/claim": {
      "post": {
        "tags": [
          "signups"
        ],
        "operationId": "claimSignup",
        "summary": "Claim a signup and turn it into a user (usually from a url in an activation email).",
        "parameters": [
          {
            "in": "path",
            "description": "The signup id",
            "name": "id",
            "required": true,
            "type": "integer",
            "format": "int32",
            "minimum": 0,
            "maximum": 0
          },
          {
            "in": "body",
            "description": "The signup verification token",
            "name": "token",
            "required": true,
            "schema": {
              "type": "string",
              "minimum": 0,
              "maximum": 0
            }
          },
          {
            "in": "body",
            "description": "The desired plaintext password for the new user",
            "name": "password",
            "required": true,
            "schema": {
              "type": "string",
              "minimum": 0,
              "maximum": 0
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Description was not specified",
            "schema": {
              "$ref": "#/definitions/github.com.opsee.vape.api.UserTokenResponse",
              "items": {}
            }
          },
          "401": {
            "description": "Description was not specified"
          },
          "409": {
            "description": "Description was not specified"
          }
        }
      }
    },
    "/users/{id}": {
      "get": {
        "tags": [
          "users"
        ],
        "operationId": "getUser",
        "summary": "Get a single user.",
        "parameters": [
          {
            "in": "header",
            "description": "The Bearer token - an admin user token or a token with matching id is required",
            "name": "Authorization",
            "required": true,
            "type": "string",
            "minimum": 0,
            "maximum": 0
          },
          {
            "in": "path",
            "description": "The user id",
            "name": "id",
            "required": true,
            "type": "integer",
            "format": "int32",
            "minimum": 0,
            "maximum": 0
          }
        ],
        "responses": {
          "200": {
            "description": "Description was not specified",
            "schema": {
              "$ref": "#/definitions/github.com.opsee.vape.model.User",
              "items": {}
            }
          },
          "401": {
            "description": "Description was not specified"
          }
        }
      },
      "put": {
        "tags": [
          "users"
        ],
        "operationId": "updateUser",
        "summary": "Update a single user.",
        "parameters": [
          {
            "in": "header",
            "description": "The Bearer token - an admin user token or a token with matching id is required",
            "name": "Authorization",
            "required": true,
            "type": "string",
            "minimum": 0,
            "maximum": 0
          },
          {
            "in": "path",
            "description": "The user id",
            "name": "id",
            "required": true,
            "type": "integer",
            "format": "int32",
            "minimum": 0,
            "maximum": 0
          },
          {
            "in": "body",
            "description": "A new email address",
            "name": "email",
            "required": false,
            "schema": {
              "type": "string",
              "minimum": 0,
              "maximum": 0
            }
          },
          {
            "in": "body",
            "description": "A new name",
            "name": "name",
            "required": false,
            "schema": {
              "type": "string",
              "minimum": 0,
              "maximum": 0
            }
          },
          {
            "in": "body",
            "description": "A new password",
            "name": "password",
            "required": false,
            "schema": {
              "type": "string",
              "minimum": 0,
              "maximum": 0
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Description was not specified",
            "schema": {
              "$ref": "#/definitions/github.com.opsee.vape.api.UserTokenResponse",
              "items": {}
            }
          },
          "401": {
            "description": "Description was not specified"
          }
        }
      },
      "delete": {
        "tags": [
          "users"
        ],
        "operationId": "deleteUser",
        "summary": "Update a single user.",
        "parameters": [
          {
            "in": "header",
            "description": "The Bearer token - an admin user token or a token with matching id is required",
            "name": "Authorization",
            "required": true,
            "type": "string",
            "minimum": 0,
            "maximum": 0
          },
          {
            "in": "path",
            "description": "The user id",
            "name": "id",
            "required": true,
            "type": "integer",
            "format": "int32",
            "minimum": 0,
            "maximum": 0
          }
        ],
        "responses": {
          "200": {
            "description": "Description was not specified",
            "schema": {
              "$ref": "#/definitions/github.com.opsee.vape.api.MessageResponse",
              "items": {}
            }
          },
          "401": {
            "description": "Description was not specified"
          }
        }
      }
    },
    "/users/{id}/data": {
      "put": {
        "tags": [
          "users"
        ],
        "operationId": "updateUserData",
        "summary": "Update a single user.",
        "parameters": [
          {
            "in": "header",
            "description": "The Bearer token - an admin user token or a token with matching id is required",
            "name": "Authorization",
            "required": true,
            "type": "string",
            "minimum": 0,
            "maximum": 0
          },
          {
            "in": "path",
            "description": "The user id",
            "name": "id",
            "required": true,
            "type": "integer",
            "format": "int32",
            "minimum": 0,
            "maximum": 0
          }
        ],
        "responses": {
          "200": {
            "description": "Description was not specified",
            "schema": {
              "$ref": "#/definitions/github.com.opsee.vape.api.UserDataResponse",
              "items": {}
            }
          },
          "401": {
            "description": "Description was not specified"
          }
        }
      }
    }
  },
  "definitions": {
    "github.com.opsee.vape.api.MessageResponse": {
      "properties": {
        "message": {
          "type": "string",
          "items": {}
        }
      }
    },
    "github.com.opsee.vape.api.UserTokenResponse": {
      "properties": {
        "token": {
          "type": "string",
          "items": {}
        },
        "user": {
          "$ref": "#/definitions/github.com.opsee.vape.model.User",
          "items": {}
        }
      }
    },
    "github.com.opsee.vape.model.User": {
      "properties": {
        "active": {
          "$ref": "#/definitions/bool",
          "items": {}
        },
        "admin": {
          "$ref": "#/definitions/bool",
          "items": {}
        },
        "created_at": {
          "$ref": "#/definitions/Time",
          "items": {}
        },
        "customer_id": {
          "type": "string",
          "items": {}
        },
        "email": {
          "type": "string",
          "items": {}
        },
        "id": {
          "type": "integer",
          "format": "int32",
          "items": {}
        },
        "name": {
          "type": "string",
          "items": {}
        },
        "updated_at": {
          "$ref": "#/definitions/Time",
          "items": {}
        },
        "verified": {
          "$ref": "#/definitions/bool",
          "items": {}
        }
      }
    },
    "github.com.opsee.vape.api.SignupActivationResponse": {
      "properties": {
        "token": {
          "type": "string",
          "items": {}
        }
      }
    },
    "github.com.opsee.vape.model.Signup": {
      "properties": {
        "activated": {
          "$ref": "#/definitions/bool",
          "items": {}
        },
        "claimed": {
          "$ref": "#/definitions/bool",
          "items": {}
        },
        "created_at": {
          "$ref": "#/definitions/Time",
          "items": {}
        },
        "email": {
          "type": "string",
          "items": {}
        },
        "id": {
          "type": "integer",
          "format": "int32",
          "items": {}
        },
        "name": {
          "type": "string",
          "items": {}
        },
        "updated_at": {
          "$ref": "#/definitions/Time",
          "items": {}
        }
      }
    },
    "github.com.opsee.vape.api.UserDataResponse": {}
  }
}
`
