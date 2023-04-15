package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	apiv0 "github.com/agwermann/dt-operator/apis/dtd/v0"
	dtdl "github.com/agwermann/dt-operator/cmd/cli/dtdl"
	pkg "github.com/agwermann/dt-operator/cmd/cli/pkg"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sJson "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

func main() {
	allArgs := os.Args
	args := allArgs[1:]

	if len(args) < 2 {
		log.Fatal("Inform DTDL input and output folders path")
	}

	inputFolderPath := args[0]
	outputFolderPath := args[1]

	files, err := ioutil.ReadDir(inputFolderPath)

	if err != nil {
		log.Fatal(err)
	}

	serializer := k8sJson.NewYAMLSerializer(k8sJson.DefaultMetaFactory, nil, nil)

	pkg.PrepareOutputFolder(outputFolderPath)

	for _, file := range files {
		if !file.IsDir() {
			inputFilename := filepath.Join(inputFolderPath, file.Name())
			outputFileName := strings.Split(file.Name(), ".")[0]
			outputFilename := filepath.Join(outputFolderPath, outputFileName+".yaml")
			if pkg.IsJsonFile(inputFilename) {
				fileContent, err := os.ReadFile(inputFilename)
				if err != nil {
					log.Fatal(err)
				}

				twinInterface := dtdl.Interface{}
				err = json.Unmarshal(fileContent, &twinInterface)
				if err != nil {
					log.Fatal(err)
				}
				//twinYaml, err := yaml.Marshal(twinInterface)

				ti := createTwinInterfaceK8sResource(twinInterface)
				tinstance := createTwinInstanceK8sResources(ti)

				yamlBuffer := new(bytes.Buffer)
				serializer.Encode(&ti, yamlBuffer)
				yamlBuffer.Write([]byte("---\n"))
				serializer.Encode(&tinstance, yamlBuffer)

				fmt.Printf(yamlBuffer.String())

				if err != nil {
					log.Fatal(err)
				}

				err = pkg.WriteToFile(outputFilename, yamlBuffer.Bytes())
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

}

// TODO: renew TwinComponent to TwinInstance
func createTwinInterfaceK8sResource(tInterface dtdl.Interface) apiv0.TwinComponent {
	var properties []apiv0.TwinProperty
	var relationships []apiv0.TwinRelationship
	var telemetries []apiv0.TwinTelemetry
	var commands []apiv0.TwinCommand
	var extendedComponent apiv0.TwinComponentExtendsSpec

	for _, content := range tInterface.Contents {
		if content.Property != nil {
			properties = processProperty(*content.Property, properties)
		}
		if content.Relationship != nil {
			relationships = processRelationship(*content.Relationship, relationships)
		}
		if content.Telemetry != nil {
			telemetries = processTelemetry(*content.Telemetry, telemetries)
		}
		if content.Command != nil {
			commands = processCommand(*content.Command, commands)
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

func processCommand(command dtdl.Command, commands []apiv0.TwinCommand) []apiv0.TwinCommand {
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

func processTelemetry(telemetry dtdl.Telemetry, telemetries []apiv0.TwinTelemetry) []apiv0.TwinTelemetry {
	twinSchema := createTwinSchema(telemetry.Schema)
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

func processRelationship(relationship dtdl.Relationship, relationships []apiv0.TwinRelationship) []apiv0.TwinRelationship {

	var relationshipProperties []apiv0.TwinProperty

	if relationship.Properties != nil {
		for _, property := range relationship.Properties {
			relationshipProperties = processProperty(property, relationshipProperties)
		}
	}

	twinSchema := createTwinSchema(relationship.Schema)
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

func processProperty(property dtdl.Property, properties []apiv0.TwinProperty) []apiv0.TwinProperty {
	twinSchema := createTwinSchema(property.Schema)
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

func createTwinInstanceK8sResources(twinInterface apiv0.TwinComponent) apiv0.TwinInstance {
	twinInstance := apiv0.TwinInstance{
		TypeMeta: v1.TypeMeta{
			Kind:       "TwinInstance",
			APIVersion: "dtd.digitaltwin/v0",
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

func createTwinSchema(schema dtdl.Schema) apiv0.TwinSchema {
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
