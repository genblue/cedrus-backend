// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-12-13 12:56:30.082367 +0100 CET m=+1.910769929

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "email": "email@ded.fr"
        },
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/claims": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Get all Claims",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Claim"
                        }
                    }
                }
            },
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "Create a Claim",
                "parameters": [
                    {
                        "description": "New Claim",
                        "name": "new-claim",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/inputs.NewClaim"
                        }
                    }
                ],
                "responses": {
                    "201": {},
                    "400": {}
                }
            }
        },
        "/claims/{id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Get a Claim",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Claim ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Claim"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "summary": "get application health",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Health"
                        }
                    }
                }
            }
        },
        "/transfer": {
            "post": {
                "summary": "Transfer Cedar coins",
                "parameters": [
                    {
                        "description": "Transfer",
                        "name": "transfer",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/inputs.NewTransfer"
                        }
                    }
                ],
                "responses": {
                    "202": {}
                }
            }
        }
    },
    "definitions": {
        "inputs.NewClaim": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "treeCount": {
                    "type": "integer"
                }
            }
        },
        "inputs.NewTransfer": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "claim-code": {
                    "type": "string"
                }
            }
        },
        "model.Claim": {
            "type": "object",
            "properties": {
                "_id": {
                    "type": "string"
                },
                "claimCode": {
                    "type": "string"
                },
                "creationDate": {
                    "type": "integer"
                },
                "email": {
                    "type": "string"
                },
                "emailSent": {
                    "type": "boolean"
                },
                "emailSentDate": {
                    "type": "integer"
                },
                "memo": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "settlementDate": {
                    "type": "integer"
                },
                "status": {
                    "type": "integer"
                },
                "transferAddress": {
                    "type": "string"
                },
                "treeCount": {
                    "type": "integer"
                }
            }
        },
        "model.Health": {
            "type": "object",
            "properties": {
                "connectedToDb": {
                    "type": "boolean"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "v1",
	Host:        "",
	BasePath:    "/api/v1",
	Schemes:     []string{},
	Title:       "Cedrus service API",
	Description: "For managing claims",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
