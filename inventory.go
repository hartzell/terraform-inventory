package main

// interface{} is a bit too loose, really a map to one or more
// Group's and one Meta full of host variables
type Inventory map[string]interface{}

type GroupInfo struct {
	Hosts []string          `json:"hosts"`
	Vars  map[string]string `json:"vars,omitempty"`
}

func (i *Inventory) SetGroupInfo(group_name string, info GroupInfo) {
	(*i)[group_name] = info
}

func (i *Inventory) SetMetaInfo(meta MetaInfo) {
	(*i)["_meta"] = meta
}

type MetaInfo struct {
	HostVars map[string]HostVars `json:"hostvars,omitempty"`
}

type HostVars map[string]string

func (h *HostVars) setHostVar(name string, value string) {
	(*h)[name] = value
}

func (m *MetaInfo) setHostVars(host_name string, hv HostVars) {
	if m.HostVars == nil {
		m.HostVars = make(map[string]HostVars)
	}
	m.HostVars[host_name] = hv
}

// GroupInfo support

func (g *GroupInfo) AddHosts(hostnames ...string) {
	(*g).Hosts = append((*g).Hosts, hostnames...)
}

func (g *GroupInfo) SetVar(name string, value string) {
	if (*g).Vars == nil {
		(*g).Vars = make(map[string]string)
	}
	(*g).Vars[name] = value
}
