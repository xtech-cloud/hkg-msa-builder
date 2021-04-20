package config

const defaultYAML string = `
service:
    name: xtc.api.hkg.builder
    address: :19903
    ttl: 15
    interval: 10
logger:
    level: info
    dir: /var/log/hkg/
database:
    mongodb:
        address: localhost:27017
        timeout: 10
        user: root
        password: mongodb@OMO
        db: hkg_builder
`
