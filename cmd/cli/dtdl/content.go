package dtdl

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
)

var (
	ContentCommandType      = "Command"
	ContentComponentType    = "Component"
	ContentPropertyType     = "Property"
	ContentRelationshipType = "Relationship"
	ContentTelemetryType    = "Telemetry"

	ErrContentUnmarshalInvalidType = errors.New("Invalid content @type")
)

func ErrContentUnmarshalTypeNotSupported(typeValue string) error {
	return errors.New(fmt.Sprintf("Content @type %s not supported", typeValue))
}

type Content struct {
	command      Command
	component    Component
	property     Property
	relationship Relationship
	telemetry    Telemetry
}

func (c *Content) UnmarshalJSON(data []byte) error {
	var jsonObject interface{}
	err := json.Unmarshal(data, &jsonObject)

	if err != nil {
		return err
	}

	objectMap := jsonObject.(map[string]interface{})

	objectType := objectMap["@type"].(string)

	switch objectMap["@type"] {
	case ContentPropertyType:
		c.property = c.newProperty(objectMap)
		return nil
	case ContentRelationshipType:
		c.relationship = c.newRelationship(objectMap)
		return nil
	case ContentCommandType:
		return ErrContentUnmarshalTypeNotSupported(objectType)
	case ContentComponentType:
		return ErrContentUnmarshalTypeNotSupported(objectType)
	case ContentTelemetryType:
		return ErrContentUnmarshalTypeNotSupported(objectType)
	default:
		return ErrContentUnmarshalInvalidType
	}
}

func (c Content) MarshalYAML() (interface{}, error) {
	if !reflect.DeepEqual(c.command, Command{}) {
		return c.command, nil
	}

	if !reflect.DeepEqual(c.relationship, Relationship{}) {
		return c.relationship, nil
	}

	if !reflect.DeepEqual(c.property, Property{}) {
		return c.property, nil
	}

	if !reflect.DeepEqual(c.telemetry, Telemetry{}) {
		return c.telemetry, nil
	}

	if !reflect.DeepEqual(c.component, Component{}) {
		return c.telemetry, nil
	}

	return nil, errors.New("Not possible to marshal Yaml")
}

func (s *Content) newProperty(data interface{}) Property {
	property := Property{}

	dataByte, err := json.Marshal(data)

	if err != nil {
		log.Fatal("Error in marshaling")
	}

	err = json.Unmarshal(dataByte, &property)

	if err != nil {
		log.Fatal("Error in unmarshaling")
	}

	return property
}

func (s *Content) newRelationship(data interface{}) Relationship {
	relationship := Relationship{}

	dataByte, err := json.Marshal(data)

	if err != nil {
		log.Fatal("Error in marshaling")
	}

	err = json.Unmarshal(dataByte, &relationship)

	if err != nil {
		log.Fatal("Error in unmarshaling")
	}

	return relationship
}
