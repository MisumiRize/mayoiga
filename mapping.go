package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

const mappingPath = "./.mayoiga/mapping.json"

type mapping struct {
	Variables []string
	Aliases   map[string]string
}

func (m *mapping) getVariables() []string {
	variables := m.Variables
	if variables == nil {
		return []string{}
	}
	return variables
}

func (m *mapping) getAliases() map[string]string {
	aliases := m.Aliases
	if aliases == nil {
		return map[string]string{}
	}
	return aliases
}

func (m *mapping) isEmpty() bool {
	return len(m.getVariables()) == 0 && len(m.getAliases()) == 0
}

func (m *mapping) remove(variable string) mapping {
	variables := []string{}
	for _, v := range m.getVariables() {
		if v != variable {
			variables = append(variables, v)
		}
	}

	aliases := map[string]string{}
	for k, v := range m.getAliases() {
		if v != variable {
			aliases[k] = v
		}
	}

	return mapping{variables, aliases}
}

func (m *mapping) merge(merge mapping) mapping {
	variables := m.getVariables()
	for _, variable := range merge.Variables {
		variables = addIfNotExist(variables, variable)
	}

	aliases := m.getAliases()
	for variable, alias := range merge.Aliases {
		aliases[variable] = alias
	}

	return mapping{
		Variables: variables,
		Aliases:   aliases,
	}
}

func (m *mapping) buildMappedEnvBody(env map[string]string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)

	for _, k := range m.Variables {
		if v, ok := env[k]; ok {
			if _, err := buf.WriteString(fmt.Sprintln(fmt.Sprintf("%s=%s", k, v))); err != nil {
				return nil, err
			}
		}
	}

	for k, a := range m.Aliases {
		if v, ok := env[k]; ok {
			if _, err := buf.WriteString(fmt.Sprintln(fmt.Sprintf("%s=%s", a, v))); err != nil {
				return nil, err
			}
		}
	}

	return buf, nil
}

type mappingsWrapper struct {
	Version  int
	Mappings map[string]mapping
}

func fetchMappings() (map[string]mapping, error) {
	config, err := readConfig()
	if err != nil {
		return nil, err
	}

	buf, err := s3GetObject(config.MappingS3Key)
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() == "NoSuchKey" {
			return map[string]mapping{}, nil
		}
		return nil, err
	} else if err != nil {
		return nil, err
	}

	var wrapper mappingsWrapper
	if err = json.Unmarshal(buf.Bytes(), &wrapper); err != nil {
		return nil, err
	}
	return wrapper.Mappings, nil
}

func writeMappingsFile(mappings map[string]mapping) (err error) {
	j, err := json.MarshalIndent(mappingsWrapper{
		Version:  version,
		Mappings: mappings,
	}, "", "  ")
	if err != nil {
		return
	}

	_, err = os.Stat(configDir)
	if os.IsNotExist(err) {
		if err = os.Mkdir(configDir, 0755); err != nil {
			return
		}
	} else if err != nil {
		return
	}

	file, err := os.Create(mappingPath)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.Write(j)
	return
}

func putMappingsToS3(mappings map[string]mapping) (err error) {
	if len(mappings) == 0 {
		return
	}

	j, err := json.MarshalIndent(mappingsWrapper{
		Version:  version,
		Mappings: mappings,
	}, "", "  ")
	if err != nil {
		return
	}

	config, err := readConfig()
	if err != nil {
		return
	}

	return s3PutObject(config.MappingS3Key, bytes.NewReader(j))
}

func generateMappedEnvFiles(envBody *bytes.Buffer, mappings map[string]mapping) (err error) {
	env := parseEnv(bufio.NewScanner(envBody))

	for fileName, mapping := range mappings {
		file, err := os.Create(fileName)
		if err != nil {
			return err
		}

		buf, err := mapping.buildMappedEnvBody(env)
		if err != nil {
			return err
		}

		if _, err = file.Write(buf.Bytes()); err != nil {
			return err
		}
	}

	return nil
}

func removeVariable(mappings map[string]mapping, variable string) map[string]mapping {
	removed := map[string]mapping{}

	for file, m := range mappings {
		mapping := m.remove(variable)

		if !mapping.isEmpty() {
			removed[file] = mapping
		}
	}

	return removed
}
