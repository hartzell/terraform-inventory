package main

import (
	"encoding/json"
	"fmt"
//		"github.com/davecgh/go-spew/spew"
	"io"
)

func cmdList(stdout io.Writer, stderr io.Writer, s *state) int {
	i := Inventory{}

	for hostName, res := range s.resources() {
		// a group named w/ address just for this instance
		i.AddHostToGroup(hostName, res.Address())

		// add this host to each of the groups it wants to be part of
		for _, groupName := range res.AnsibleHostGroups() {
			i.AddHostToGroup(hostName, groupName)
		}

		// set each host variable as specified
		for varName, varValue := range res.AnsibleHostVars() {
			// TODO, should this use res.Address() instead of hostName?
			i.SetHostVar(hostName, varName, varValue)
		}
	}
//	spew.Dump(i)

	return output(stdout, stderr, i)
}

func cmdHost(stdout io.Writer, stderr io.Writer, s *state, hostname string) int {
	//	for name, res := range s.resources() {
	for name, _ := range s.resources() {
		if hostname == name {
			//			attrs := res.Attributes()
			//			for name, value := range res.AnsibleHostVars() {
			//				attrs[name] = value
			//			}
			attrs := make(map[string]string)
			return output(stdout, stderr, attrs)
		}
	}

	fmt.Fprintf(stderr, "No such host: %s\n", hostname)
	return 1
}

// output marshals an arbitrary JSON object and writes it to stdout, or writes
// an error to stderr, then returns the appropriate exit code.
func output(stdout io.Writer, stderr io.Writer, whatever interface{}) int {
	b, err := json.Marshal(whatever)
	if err != nil {
		fmt.Fprintf(stderr, "Error encoding JSON: %s\n", err)
		return 1
	}

	_, err = stdout.Write(b)
	if err != nil {
		fmt.Fprintf(stderr, "Error writing JSON: %s\n", err)
		return 1
	}

	return 0
}
