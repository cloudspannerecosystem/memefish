package analyzer

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/memefish/pkg/token"
)

type NameEnv map[string]*Name

func (env NameEnv) Lookup(text string) *Name {
	if env == nil {
		return nil
	}

	text = strings.ToUpper(text)
	if name, ok := env[text]; ok {
		return name
	}
	return nil
}

func (env NameEnv) Insert(name *Name) error {
	if name.Anonymous() {
		return nil
	}

	text := strings.ToUpper(name.Text)
	if oldName, ok := env[text]; ok {
		switch {
		case name.Kind == TableName && oldName.Kind == TableName:
			return fmt.Errorf("duplicate table name: %s", token.QuoteSQLIdent(name.Text))
		case name.Kind == TableName:
			env[text] = name
		case oldName.Kind == TableName:
			// nothing
		default:
			env[text] = makeAmbiguousName(oldName.Text, []*Name{oldName, name})
		}
	} else {
		env[text] = name
	}

	return nil
}

func (env NameEnv) InsertForce(name *Name) {
	if name.Anonymous() {
		return
	}

	text := strings.ToUpper(name.Text)
	env[text] = name
}
