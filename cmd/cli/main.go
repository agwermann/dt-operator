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

				tw := createTwinInterfaceK8sResource(twinInterface)

				yamlBuffer := new(bytes.Buffer)
				serializer.Encode(&tw, yamlBuffer)

				fmt.Printf(yamlBuffer.String())

				if err != nil {
					log.Fatal(err)
				}

				err = pkg.WriteToFile(outputFilename, yamlBuffer.Bytes())
				if err != nil {
					log.Fatal(err)
				}
				return
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

	for _, content := range tInterface.Contents {

		if content.Property != nil {
			twinSchema := createTwinSchema(content.Property.Schema)
			property := apiv0.TwinProperty{
				Id:          string(content.Property.Id),
				Comment:     content.Property.Comment,
				Description: string(content.Property.Description),
				DisplayName: string(content.Property.DisplayName),
				Name:        content.Property.Name,
				Writeable:   content.Property.Writeable,
				Schema:      twinSchema,
			}
			properties = append(properties, property)
		}

		if content.Relationship != nil {
			twinSchema := createTwinSchema(content.Relationship.Schema)
			relationship := apiv0.TwinRelationship{
				Id:              string(content.Relationship.Id),
				Comment:         content.Relationship.Comment,
				Description:     string(content.Relationship.Description),
				DisplayName:     string(content.Relationship.DisplayName),
				Name:            content.Relationship.Name,
				Writeable:       content.Relationship.Writeable,
				MaxMultiplicity: content.Relationship.MaxMultiplicity,
				MinMultiplicity: content.Relationship.MinMultiplicity,
				Schema:          twinSchema,
			}
			relationships = append(relationships, relationship)
		}

		if content.Telemetry != nil {
			twinSchema := createTwinSchema(content.Telemetry.Schema)
			telemetry := apiv0.TwinTelemetry{
				Id:          string(content.Relationship.Id),
				Comment:     content.Relationship.Comment,
				Description: string(content.Relationship.Description),
				DisplayName: string(content.Relationship.DisplayName),
				Name:        content.Relationship.Name,
				Schema:      twinSchema,
			}
			telemetries = append(telemetries, telemetry)
		}

		if content.Command != nil {
			command := apiv0.TwinCommand{
				Id:          string(content.Relationship.Id),
				Comment:     content.Relationship.Comment,
				Description: string(content.Relationship.Description),
				DisplayName: string(content.Relationship.DisplayName),
				Name:        content.Relationship.Name,
			}
			commands = append(commands, command)
		}

	}

	twinInterface := apiv0.TwinComponent{
		TypeMeta: v1.TypeMeta{
			Kind:       "TwinComponent",
			APIVersion: "dtd.digitaltwin/v0",
		},
		Spec: apiv0.TwinComponentSpec{
			Id:            string(tInterface.Id),
			Properties:    properties,
			Relationships: relationships,
			Commands:      commands,
			Telemetries:   telemetries,
		},
	}

	return twinInterface
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
