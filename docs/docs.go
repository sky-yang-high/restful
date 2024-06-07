// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/task/all": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "GetAllTasks returns all tasks in ts.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/taskstore.Task"
                            }
                        }
                    },
                    "404": {
                        "description": "No tasks found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/task/create": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "CreateTask creates a new task in ts and returns its ID.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The text of the task",
                        "name": "text",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "array",
                        "description": "The tags of the task",
                        "name": "tags",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "The due time of the task in minutes",
                        "name": "due",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "The ID of the created task",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/task/{id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "GetTask returns the task with the given ID.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "The ID of the task to retrieve",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "The task with the given ID",
                        "schema": {
                            "$ref": "#/definitions/taskstore.Task"
                        }
                    },
                    "400": {
                        "description": "Invalid task ID",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "taskstore.Task": {
            "type": "object",
            "properties": {
                "due": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "text": {
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