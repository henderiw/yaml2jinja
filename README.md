# yaml2jinja

Converts YAML specs into Jinja3 files

## Installation

### Binary Installation

Pre-compiled binaries are available on the releases page. You can download the correct binary depending on your system arch, put it into `$PATH` and hit `yaml2go help`

### Install From Source

Build binary using go build

```bash
$ go get -u github.com/henderiw/yaml2go
$ go build -o yaml2go github.com/henderiw/yaml2go/*.go
```

## Usage

#### Show help

```bash
yaml2jinja -h
yaml2jinja converts YAML specs to Jinja templates

Usage:
    yaml2jinja < /path/to/yamlspec.yaml

Examples:
    yaml2jinja < test/example1.yaml
    yaml2jinja < test/example1.yaml > example1.go
```

#### Convert yaml spec to Go struct

```bash
$ yaml2go < srl-interface-gen.yml
```
e.g

```bash
$ cat srl-interface-gen.yml
---
srl_nokia-interfaces: 
  interface: 
    - name:
      description:
      admin-state:
      mtu:
      transceiver: 
        admin-state:
        ddm-events:
        forward-error-correction:
      ethernet: 
        flow-control: 
          receive:
      subinterface: 
        - index:
          description:
          admin-state:
          ip-mtu:
          ipv4: 
            address: 
              - ip-prefix:
            allow-directed-broadcast:
          ipv6: 
            address: 
              - ip-prefix:
```

```bash
{% if interface is defined %}
  interface: 
{% for interface in interface %}
{% if interface.mtu is defined %}
    - mtu: {{ interface.mtu }}
{% endif %}
{% if interface.transceiver is defined %}
      transceiver: 
{% if interface.transceiver.admin_state is defined %}
        admin-state: {{ interface.transceiver.admin_state }}
{% endif %}
{% if interface.transceiver.ddm_events is defined %}
        ddm-events: {{ interface.transceiver.ddm_events }}
{% endif %}
{% if interface.transceiver.forward_error_correction is defined %}
        forward-error-correction: {{ interface.transceiver.forward_error_correction }}
{% endif %}
{% endif %}
{% if interface.ethernet is defined %}
      ethernet: 
{% if interface.ethernet.flow_control is defined %}
        flow-control: 
{% if interface.ethernet.flow_control.receive is defined %}
          receive: {{ interface.ethernet.flow_control.receive }}
{% endif %}
{% endif %}
{% endif %}
{% if interface.subinterface is defined %}
      subinterface: 
{% for subinterface in interface.subinterface %}
{% if subinterface.index is defined %}
        - index: {{ subinterface.index }}
{% endif %}
{% if subinterface.description is defined %}
          description: {{ subinterface.description }}
{% endif %}
{% if subinterface.admin_state is defined %}
          admin-state: {{ subinterface.admin_state }}
{% endif %}
{% if subinterface.ip_mtu is defined %}
          ip-mtu: {{ subinterface.ip_mtu }}
{% endif %}
{% if subinterface.ipv4 is defined %}
          ipv4: 
{% if subinterface.ipv4.address is defined %}
            address: 
{% for address in subinterface.ipv4.address %}
{% if address.ip_prefix is defined %}
              - ip-prefix: {{ address.ip_prefix }}
{% endif %}
{% endfor %}
{% endif %}
{% if subinterface.ipv4.allow_directed_broadcast is defined %}
            allow-directed-broadcast: {{ subinterface.ipv4.allow_directed_broadcast }}
{% endif %}
{% endif %}
{% if subinterface.ipv6 is defined %}
          ipv6: 
{% if subinterface.ipv6.address is defined %}
            address: 
{% for address in subinterface.ipv6.address %}
{% if address.ip_prefix is defined %}
              - ip-prefix: {{ address.ip_prefix }}
{% endif %}
{% endfor %}
{% endif %}
{% endif %}
{% endfor %}
{% endif %}
{% if interface.name is defined %}
        name: {{ interface.name }}
{% endif %}
{% if interface.description is defined %}
      description: {{ interface.description }}
{% endif %}
{% if interface.admin_state is defined %}
      admin-state: {{ interface.admin_state }}
{% endif %}
{% endfor %}
{% endif %}
```

## Contributing

We love your input! We want to make contributing to this project as easy and transparent as possible, whether it's:
- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features

## Inspired by
yaml2go project
