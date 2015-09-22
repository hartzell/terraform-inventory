package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func cmdList(stdout io.Writer, stderr io.Writer, s *state) int {
	//	groups := make(map[string][]string, 0)
	//
	//	// add each instance as a pseudo-group, so they can be provisioned
	//	// individually where necessary.
	//	for name, res := range s.resources() {
	//		groups[name] = []string{res.Address()}
	//
	//		for _, host_group := range res.AnsibleHostGroups() {
	//			groups[host_group] = append(groups[host_group], res.Address())
	//		}
	//	}
	//
	//	return output(stdout, stderr, groups)
	//

	i := Inventory{}

	groups := make(map[string]*GroupInfo)

	for name, res := range s.resources() {
		//		fmt.Println("host groups", name, res.Address(), res.AnsibleHostGroups())
		gi := GroupInfo{}
		gi.AddHosts(res.Address())
		i.SetGroupInfo(name, gi)
		for _, group := range res.AnsibleHostGroups() {
			if groups[group] == nil {
				groups[group] = &GroupInfo{}
			}
			groups[group].AddHosts(res.Address())
		}
	}
	for name, value := range groups {
		i.SetGroupInfo(name, *value)
	}

	m := MetaInfo{}

	for _, res := range s.resources() {
		res_name := res.Address()
		//attrs := res.Attributes()		// start off with what Terraform tells us
		attrs := make(map[string]string)
		for name, value := range res.AnsibleHostVars() {
			attrs[name] = value
		}
		hv := HostVars{}
		for name, value := range attrs {
			hv.setHostVar(name, value)
		}
		m.setHostVars(res_name, hv)
	}

	i.SetMetaInfo(m)

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
