package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const exampleStateFile = `
{
	"version": 1,
	"serial": 1,
	"modules": [
		{
			"path": [
				"root"
			],
			"outputs": {},
			"resources": {
				"aws_instance.one": {
					"type": "aws_instance",
					"primary": {
						"id": "i-aaaaaaaa",
						"attributes": {
							"id": "i-aaaaaaaa",
							"private_ip": "10.0.0.1"
						}
					}
				},
				"aws_instance.two": {
					"type": "aws_instance",
					"primary": {
						"id": "i-bbbbbbbb",
						"attributes": {
							"id": "i-bbbbbbbb",
							"private_ip": "10.0.0.2",
							"public_ip": "50.0.0.1"
						}
					}
				},
				"aws_security_group.example": {
					"type": "aws_security_group",
					"primary": {
						"id": "sg-cccccccc",
						"attributes": {
							"id": "sg-cccccccc",
							"description": "Whatever"
						}
					}
				},
				"digitalocean_droplet.three": {
					"type": "digitalocean_droplet",
					"primary": {
						"id": "ddddddd",
						"attributes": {
							"id": "ddddddd",
							"ipv4_address": "192.168.0.3"
						}
					}
				},
				"cloudstack_instance.four": {
					"type": "cloudstack_instance",
					"primary": {
						"id": "aaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
						"attributes": {
							"id": "aaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
							"ipaddress": "10.2.1.5"
						}
					}
				},
        "openstack_compute_instance_v2.five": {
          "type": "openstack_compute_instance_v2",
          "primary": {
            "id": "92dbe904-a987-4ad4-963d-0f9ba0cb2b96",
            "attributes": {
              "access_ip_v4": "10.33.147.228",
              "access_ip_v6": "",
              "floating_ip": "10.33.147.228"
            }
          }
        },
        "openstack_compute_instance_v2.six": {
          "type": "openstack_compute_instance_v2",
          "primary": {
            "id": "92dbe904-a987-4ad4-963d-0f9ba0cb2b95",
            "attributes": {
              "access_ip_v4": "192.168.1.1",
              "access_ip_v6": ""
            }
          }
        }
			}
		}
	]
}
`

func TestStateRead(t *testing.T) {
	var s state
	r := strings.NewReader(exampleStateFile)
	err := s.read(r)
	assert.Nil(t, err)
	assert.Equal(t, "aws_instance", s.Modules[0].Resources["aws_instance.one"].Type)
}

func TestResources(t *testing.T) {
	r := strings.NewReader(exampleStateFile)

	var s state
	err := s.read(r)
	assert.Nil(t, err)

	inst := s.resources()
	assert.Equal(t, 6, len(inst))
	assert.Equal(t, "aws_instance", inst["one"].Type)
	assert.Equal(t, "aws_instance", inst["two"].Type)
	assert.Equal(t, "digitalocean_droplet", inst["three"].Type)
	assert.Equal(t, "cloudstack_instance", inst["four"].Type)
	assert.Equal(t, "openstack_compute_instance_v2", inst["five"].Type)
	assert.Equal(t, "openstack_compute_instance_v2", inst["six"].Type)
}

func TestAddress(t *testing.T) {
	r := strings.NewReader(exampleStateFile)

	var s state
	err := s.read(r)
	assert.Nil(t, err)

	inst := s.resources()
	assert.Equal(t, 6, len(inst))
	assert.Equal(t, "10.0.0.1", inst["one"].Address())
	assert.Equal(t, "50.0.0.1", inst["two"].Address())
	assert.Equal(t, "192.168.0.3", inst["three"].Address())
	assert.Equal(t, "10.2.1.5", inst["four"].Address())
	assert.Equal(t, "10.33.147.228", inst["five"].Address())
	assert.Equal(t, "192.168.1.1", inst["six"].Address())
}

func TestIsSupported(t *testing.T) {
	r := resourceState{
		Type: "something",
	}
	assert.Equal(t, false, r.isSupported())

	r = resourceState{
		Type: "aws_instance",
		Primary: instanceState{
			Attributes: map[string]string{
				"private_ip": "10.0.0.2",
			},
		},
	}
	assert.Equal(t, true, r.isSupported())

	r = resourceState{
		Type: "digitalocean_droplet",
		Primary: instanceState{
			Attributes: map[string]string{
				"ipv4_address": "192.168.0.3",
			},
		},
	}
	assert.Equal(t, true, r.isSupported())

	r = resourceState{
		Type: "cloudstack_instance",
		Primary: instanceState{
			Attributes: map[string]string{
				"ipaddress": "10.2.1.5",
			},
		},
	}
	assert.Equal(t, true, r.isSupported())

	r = resourceState{
		Type: "openstack_compute_instance_v2",
		Primary: instanceState{
			Attributes: map[string]string{
				"floating_ip": "10.2.1.5",
			},
		},
	}
	assert.Equal(t, true, r.isSupported())

	r = resourceState{
		Type: "openstack_compute_instance_v2",
		Primary: instanceState{
			Attributes: map[string]string{
				"access_ip_v4": "10.2.1.5",
			},
		},
	}
	assert.Equal(t, true, r.isSupported())
}
