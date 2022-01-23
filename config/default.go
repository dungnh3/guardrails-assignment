package config

var defaultValue = `
postgresql:
  host: 127.0.0.1
  port: 5432
  username: guardrails
  password: secret
  database: guard_db
  options: sslmode=disable
`
