package pkg

import (
	"log"
	"regexp"
	"strings"

	apiv0 "github.com/agwermann/dt-operator/apis/dtd/v0"
	dtdl "github.com/agwermann/dt-operator/cmd/cli/dtdl"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ResourceBuilder interface {
	CreateTwinInterface(tInterface dtdl.Interface) apiv0.TwinComponent
	CreateTwinInstance(twinInterface apiv0.TwinComponent) apiv0.TwinInstance
}

func NewResourceBuilder() ResourceBuilder {
	return &resourceBuilder{}
}

type resourceBuilder struct {
}

// TODO: renew TwinComponent to TwinInstance
func (r *resourceBuilder) CreateTwinInterface(tInterface dtdl.Interface) apiv0.TwinComponent {
	var properties []apiv0.TwinProperty
	var relationships []apiv0.TwinRelationship
	var telemetries []apiv0.TwinTelemetry
	var commands []apiv0.TwinCommand
	var extendedComponent apiv0.TwinComponentExtendsSpec

	for _, content := range tInterface.Contents {
		if content.Property != nil {
			properties = r.processProperty(*content.Property, properties)
		}
		if content.Relationship != nil {
			relationships = r.processRelationship(*content.Relationship, relationships)
		}
		if content.Telemetry != nil {
			telemetries = r.processTelemetry(*content.Telemetry, telemetries)
		}
		if content.Command != nil {
			commands = r.processCommand(*content.Command, commands)
		}
	}

	// Only supports one parent interface
	if len(tInterface.Extends) > 0 {
		extendedComponent = apiv0.TwinComponentExtendsSpec{
			Id: tInterface.Extends[0],
		}
	}

	twinInterface := apiv0.TwinComponent{
		TypeMeta: v1.TypeMeta{
			Kind:       "TwinInterface",
			APIVersion: "dtd.digitaltwin/v0",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      r.parseHostName(string(tInterface.Id)),
			Namespace: "default",
		},
		Spec: apiv0.TwinComponentSpec{
			Id:            string(tInterface.Id),
			DisplayName:   string(tInterface.DisplayName),
			Description:   string(tInterface.Description),
			Comment:       string(tInterface.Comment),
			Properties:    properties,
			Relationships: relationships,
			Commands:      commands,
			Telemetries:   telemetries,
			Extends:       extendedComponent,
		},
	}

	return twinInterface
}

func (r *resourceBuilder) CreateTwinInstance(twinInterface apiv0.TwinComponent) apiv0.TwinInstance {
	twinInstance := apiv0.TwinInstance{
		TypeMeta: v1.TypeMeta{
			Kind:       "TwinInstance",
			APIVersion: "dtd.digitaltwin/v0",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      r.parseHostName(string(twinInterface.Spec.Id)),
			Namespace: "default",
		},
		Spec: apiv0.TwinInstanceSpec{
			Id: twinInterface.Spec.Id + "-instance",
			Component: apiv0.TwinComponentSpec{
				Id: twinInterface.Spec.Id,
			},
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{Containers: []corev1.Container{
					{
						Name:            "ktwin/" + twinInterface.Spec.Id,
						Image:           "ktwin/" + twinInterface.Spec.Id + ":0.0.1",
						ImagePullPolicy: corev1.PullIfNotPresent,
						Args:            []string{"http://mqtt-response-handler", "80"},
					},
				}},
			},
		},
	}

	return twinInstance
}

// Parse the string and make it compliant with RFC 1123 host names, by removing invalid characters
func (r *resourceBuilder) parseHostName(name string) string {
	newName := strings.ToLower(name)
	invalidCharacters := []string{":", ";", "_"}

	for _, invalidString := range invalidCharacters {
		newName = strings.Replace(newName, invalidString, "-", -1)
	}

	_, err := regexp.MatchString("[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*", newName)

	if err != nil {
		log.Fatal("Error matching host name:", err.Error())
	}

	return newName
}

func (r *resourceBuilder) processCommand(command dtdl.Command, commands []apiv0.TwinCommand) []apiv0.TwinCommand {
	newCommand := apiv0.TwinCommand{
		Id:          string(command.Id),
		Comment:     command.Comment,
		Description: string(command.Description),
		DisplayName: string(command.DisplayName),
		Name:        command.Name,
	}
	commands = append(commands, newCommand)
	return commands
}

func (r *resourceBuilder) processTelemetry(telemetry dtdl.Telemetry, telemetries []apiv0.TwinTelemetry) []apiv0.TwinTelemetry {
	twinSchema := r.createTwinSchema(telemetry.Schema)
	newTelemetry := apiv0.TwinTelemetry{
		Id:          string(telemetry.Id),
		Comment:     telemetry.Comment,
		Description: string(telemetry.Description),
		DisplayName: string(telemetry.DisplayName),
		Name:        telemetry.Name,
		Schema:      twinSchema,
	}
	telemetries = append(telemetries, newTelemetry)
	return telemetries
}

func (r *resourceBuilder) processRelationship(relationship dtdl.Relationship, relationships []apiv0.TwinRelationship) []apiv0.TwinRelationship {

	var relationshipProperties []apiv0.TwinProperty

	if relationship.Properties != nil {
		for _, property := range relationship.Properties {
			relationshipProperties = r.processProperty(property, relationshipProperties)
		}
	}

	twinSchema := r.createTwinSchema(relationship.Schema)
	newRelationship := apiv0.TwinRelationship{
		Id:              string(relationship.Id),
		Comment:         relationship.Comment,
		Description:     string(relationship.Description),
		DisplayName:     string(relationship.DisplayName),
		Name:            relationship.Name,
		Writeable:       relationship.Writeable,
		MaxMultiplicity: relationship.MaxMultiplicity,
		MinMultiplicity: relationship.MinMultiplicity,
		Schema:          twinSchema,
		Properties:      relationshipProperties,
		Target:          string(relationship.Target),
	}
	relationships = append(relationships, newRelationship)
	return relationships
}

func (r *resourceBuilder) processProperty(property dtdl.Property, properties []apiv0.TwinProperty) []apiv0.TwinProperty {
	twinSchema := r.createTwinSchema(property.Schema)
	newProperty := apiv0.TwinProperty{
		Id:          string(property.Id),
		Comment:     property.Comment,
		Description: string(property.Description),
		DisplayName: string(property.DisplayName),
		Name:        property.Name,
		Writeable:   property.Writeable,
		Schema:      twinSchema,
	}
	properties = append(properties, newProperty)
	return properties
}

func (r *resourceBuilder) createTwinSchema(schema dtdl.Schema) apiv0.TwinSchema {
	var twinEnumSchemaValues []apiv0.TwinEnumSchemaValues
	var twinEnumSchema apiv0.TwinEnumSchema

	for _, enumValue := range schema.EnumSchema.EnumValues {
		twinEnumValue := apiv0.TwinEnumSchemaValues{
			Name:        enumValue.Name,
			DisplayName: enumValue.DisplayName,
			EnumValue:   enumValue.EnumValue,
		}
		twinEnumSchemaValues = append(twinEnumSchemaValues, twinEnumValue)
	}

	twinEnumSchema = apiv0.TwinEnumSchema{
		ValueSchema: apiv0.PrimitiveType(schema.EnumSchema.ValueSchema),
		EnumValues:  twinEnumSchemaValues,
	}

	twinSchema := apiv0.TwinSchema{
		PrimitiveType: apiv0.PrimitiveType(schema.DefaultSchemaValue),
		EnumType:      twinEnumSchema,
	}

	return twinSchema
}
