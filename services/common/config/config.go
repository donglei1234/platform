package config

import (
	"bytes"
	"errors"
	"log"
	"reflect"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/fx"
)

// EnvironmentBlock is used to compose an environment config block.
type EnvironmentBlock struct {
	fx.Out
	HelpOut
}

// Load fills in all of a structure's variables from the environment with prefix.  If spec has a provided Help then
// it will also be filled in.
func LoadWithPrefix(prefix string, spec interface{}) error {
	if err := envconfig.Process(prefix, spec); err != nil {
		return err
	} else {
		if helpValue := reflect.ValueOf(spec).Elem().FieldByNameFunc(func(s string) bool {
			return s == "Help"
		}); helpValue.IsValid() && helpValue.Type() == reflect.TypeOf(Help{}) {
			if err := extractHelp(prefix, spec, helpValue.Addr().Interface().(*Help)); err != nil {
				return err
			}
		}
	}

	return nil
}

func LoadFromVault(spec interface{}) error {
	if !reflect.ValueOf(spec).Elem().FieldByName("VaultAddr").IsValid() {
		return nil
	}
	conf := &api.Config{
		Address: reflect.ValueOf(spec).Elem().FieldByName("VaultAddr").String(),
	}
	client, err := api.NewClient(conf)
	if err != nil {
		return err
	}
	client.SetToken(reflect.ValueOf(spec).Elem().FieldByName("VaultToken").String())
	secret, err := client.Logical().Read(reflect.ValueOf(spec).Elem().FieldByName("VaultPath").String())
	if err != nil {
		if strings.Contains(err.Error(), "connect: connection refused") ||
			strings.Contains(err.Error(), "connect: connection timed out") {
			log.Println("vault address cannot be connected, only environment configuration available")
			return nil
		} else if strings.Contains(err.Error(), "permission denied") {
			log.Println("vault token is invalid, only environment configuration available")
			return nil
		}
		return err
	}
	m, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return errors.New("can't read from vault")
	}
	t := reflect.TypeOf(spec).Elem()
	num := t.NumField()
	for i := 0; i < num; i++ {
		key := t.Field(i).Tag.Get("vault")
		if key != "" && m[key] != nil {
			reflect.ValueOf(spec).Elem().FieldByName(t.Field(i).Name).SetString(m[key].(string))
		}
	}
	return nil
}

// Load fills in all of a structure's variables from the environment.  If spec has a provided Help then
// it will also be filled in.
func Load(spec interface{}) error {
	return LoadWithPrefix("", spec)
	//if err := LoadWithPrefix("", spec); err != nil {
	//	return err
	//}
	//return LoadFromVault(spec)
}

// NewHelp fills in an existing help structure by reflecting over spec.
func extractHelp(prefix string, spec interface{}, help *Help) error {
	buf := new(bytes.Buffer)
	if err := envconfig.Usagef(prefix, spec, buf, usageFormat); err != nil {
		return err
	} else {
		help.Usage = buf.String()
		return nil
	}
}
