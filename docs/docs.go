// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/auth/sign-up/": {
            "post": {
                "description": "Create a new student with the input payload",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Students"
                ],
                "summary": "Create a new student",
                "parameters": [
                    {
                        "description": "Student data",
                        "name": "student",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Student"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.Student"
                        }
                    }
                }
            }
        },
        "/api/sign-in": {
            "post": {
                "responses": {}
            }
        },
        "/api/students/": {
            "post": {
                "description": "Create a new student with the input payload",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Students"
                ],
                "summary": "Create a new student",
                "parameters": [
                    {
                        "description": "Student data",
                        "name": "student",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Student"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.Student"
                        }
                    }
                }
            }
        },
        "/students/{id}": {
            "get": {
                "description": "Get student by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "students"
                ],
                "summary": "Get student by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Student ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Student"
                        }
                    }
                }
            },
            "put": {
                "responses": {}
            },
            "delete": {
                "responses": {}
            }
        }
    },
    "definitions": {
        "model.Student": {
            "type": "object",
            "properties": {
                "age": {
                    "type": "string"
                },
                "courses": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
