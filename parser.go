package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
)

type state struct {
	Modules []moduleState `json:"modules"`
}

// keyNames contains the names of the keys to check for in each resource in the
// state file. This allows us to support multiple types of resource without too
// much fuss.
var keyNames []string

func init() {
	keyNames = []string{
		"ipv4_address", // DO
		"public_ip",    // AWS
		"private_ip",   // AWS
		"ipaddress",    // CS
		"floating_ip",  // OS, best choice (more likely to be reachable)
		"access_ip_v4", // OS, second choice (less likely to be reachable)
	}
}

// read populates the state object from a statefile.
func (s *state) read(stateFile io.Reader) error {

	// read statefile contents
	b, err := ioutil.ReadAll(stateFile)
	if err != nil {
		return err
	}

	// parse into struct
	err = json.Unmarshal(b, s)
	if err != nil {
		return err
	}

	return nil
}

// resources returns a map of name to resourceState, for any supported resources
// found in the statefile.
func (s *state) resources() map[string]resourceState {
	typeRemover := regexp.MustCompile(`^[\w_]+\.`)
	inst := make(map[string]resourceState)

	for _, m := range s.Modules {
		for k, r := range m.Resources {
			if r.isSupported() {
				name := typeRemover.ReplaceAllString(k, "")
				inst[name] = r
			}
		}
	}

	return inst
}

type moduleState struct {
	Resources map[string]resourceState `json:"resources"`
}

type resourceState struct {
	Type    string        `json:"type"`
	Primary instanceState `json:"primary"`
}

// isSupported returns true if terraform-inventory supports this resource.
func (s resourceState) isSupported() bool {
	return s.Address() != ""
}

// Address returns the IP address of this resource.
func (s resourceState) Address() string {
	for _, key := range keyNames {
		if ip := s.Primary.Attributes[key]; ip != "" {
			return ip
		}
	}

	return ""
}

// Attributes returns a map containing everything we know about this resource.
func (s resourceState) Attributes() map[string]string {
	return s.Primary.Attributes
}

type instanceState struct {
	ID         string            `json:"id"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// NovaMetadata returns a map[string]string of a resource's metadata attributes.
// Skips the count ("metadata.#)
func (s resourceState) NovaMetadata() map[string]string {
	metadata := make(map[string]string)
	meta_matcher := regexp.MustCompile(`^metadata.[^#]`)

	for attrname, attr := range s.Attributes() {
		if meta_matcher.MatchString(attrname) {
			metaname := strings.TrimPrefix(attrname, "metadata.")
			metadata[metaname] = attr
		}
	}

	return metadata
}

// AnsibleHostGroups looks for a piece of metadata named
// "ansible_host_groups" and returns the slice created by splitting its
// contents on a comma (will swallow whitespace around the comma).
func (s resourceState) AnsibleHostGroups() []string {
	for metaname, attr := range s.NovaMetadata() {
		if metaname == "ansible_host_groups" {
			return regexp.MustCompile("\\s*,\\s*").Split(attr, -1)
		}
	}
	return make([]string, 0)
}
