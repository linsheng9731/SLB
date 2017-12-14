package config

import (
	"errors"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	lg "log"
	"strings"
)

// Generated at http://jsonschema.net/#/
const ConfigSchema = `
{
  "id": "/",
  "type": "object",
  "properties": {
    "general": {
      "id": "general",
      "type": "object",
      "properties": {
        "logLevel": {
          "id": "logLevel",
          "type": "string"
        },
        "logSize": {
          "id": "logSize",
          "type": "integer"
        },
        "apiHost": {
          "id": "apihost",
          "type": "string"
        },
        "apiPort": {
          "id": "apiport",
          "type": "integer"
        }
      }
    },
    "frontends": {
      "id": "frontends",
      "type": "array",
      "items": {
        "id": "0",
        "type": "object",
        "properties": {
          "name": {
            "id": "name",
            "type": "string"
          },
          "host": {
            "id": "host",
            "type": "string"
          },
          "port": {
            "id": "port",
            "type": "integer"
          },
		  "strategy": {
			"id": "strategy",
            "type": "string"
			},
          "timeout": {
            "id": "timeout",
            "type": "integer"
          },
           "heartbeatTime": {
              "id": "heartbeatTime",
               "type": "integer"
           },
           "heartbeat": {
                "id": "heartbeat",
                "type": "string"
           },
          "backends": {
            "id": "backends",
            "type": "array",
            "items": {
              "id": "0",
              "type": "object",
              "properties": {
                "name": {
                  "id": "name",
                  "type": "string"
                },
                "address": {
                  "id": "address",
                  "type": "string"
                },
                "hostname": {
                  "id": "hostname",
                  "type": "string"
                },
                "ignore_check": {
                  "id": "ignoreCheck",
                  "type": "boolean"
                },
                "weigth": {
                  "id": "weigth",
                  "type": "integer"
                }
              },
              "required": [
                "name",
                "address"
              ]
            }
          }
        },
        "required": [
          "name",
          "host",
          "port",
          "strategy"
        ]
      }
    }
  },
  "required": [
    "general",
    "frontends"
  ]
}
`

func Validate(file []byte) error {
	schemaLoader := gojsonschema.NewStringLoader(ConfigSchema)
	schema, err := gojsonschema.NewSchema(schemaLoader)

	documentLoader := gojsonschema.NewStringLoader(string(file))

	result, err := schema.Validate(documentLoader)
	if err != nil {
		lg.Println("[Error] Failed to validate", err.Error())
		return err
	}

	if !result.Valid() {
		errs := []string{}
		for _, desc := range result.Errors() {
			e := fmt.Sprintf("%s", desc)
			errs = append(errs, e)
		}
		res := strings.Join(errs, ",")
		return errors.New(res)
	}

	return nil
}
