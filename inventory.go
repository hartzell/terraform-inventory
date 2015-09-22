package main

import (
//	"fmt"
//	"github.com/davecgh/go-spew/spew"
)

// interface{} is a bit too loose, really a map to one or more
// Group's and one Meta full of host variables
type Inventory map[string]interface{}

func (i *Inventory) SetHostVar(hostName string, varName string, varValue string) {
	if (*i)["_meta"] == nil {
		(*i)["_meta"] = make(map[string]interface{})
	}
	m := (*i)["_meta"].(map[string]interface{})

	if m["hostvars"] == nil {
		m["hostvars"] = make(map[string]map[string]string)
	}
	h := m["hostvars"].(map[string]map[string]string)

	if h[hostName] == nil {
		h[hostName] = make(map[string]string)
	}

	h[hostName][varName] = varValue
}

// untested, needs additions to parser.go and data in example state file
func (i *Inventory) SetGroupVar(groupName string, varName string, varValue string) {
	if (*i)[groupName] == nil {
		(*i)[groupName] = make(map[string]interface{})
	}
	g := (*i)[groupName].(map[string]interface{})

	if g["vars"] == nil {
		g["vars"] = make(map[string]string)
	}
	gv := g["vars"].(map[string]string)

	gv[varName] = varValue
}


func (i *Inventory) AddHostToGroup(hostName string, groupName string) {
	if (*i)[groupName] == nil {
		(*i)[groupName] = make(map[string]interface{})
	}
	g := (*i)[groupName].(map[string]interface{})

	if g["hosts"] == nil {
		g["hosts"] = make([]string, 0)
	}
	g["hosts"] = append(g["hosts"].([]string), hostName)
}
