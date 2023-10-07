package kind

import (
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v3"
)

type Me struct {
	Profile map[string]any `json:"profile"`
	Allows  []string       `json:"allows"`
}

func (m Me) Allowed(check string) bool {
	//if check == "hello" {
	//	return true
	//}
	//for _, allow := range m.Allows {
	//	if checkParts := strings.Split(check, "/"); allow == checkParts[0] {
	//		return true
	//	}
	//}
	//return false
	return true
}

func (m Me) IsRoot(check string) bool {
	return len(strings.Split(check, "/")) == 1
}

func (m Me) YAML() []byte {
	raw, _ := yaml.Marshal(m)
	return raw
}

func (m Me) JSON() []byte {
	raw, _ := json.Marshal(m)
	return raw
}
