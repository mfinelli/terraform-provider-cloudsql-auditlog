terraform {
  required_providers {
    cloudsql-auditlog = {
      source = "facile.it/test/cloudsql-auditlog"
    }
  }
}

provider "cloudsql-auditlog" {
  endpoint = "127.0.0.1"
  password = ""
  username = "mario.finelli"
  engine   = "mysql"
}

# data "cloudsql-auditlog_audit_log_rules" "example" {}

# output "test" {
#   value = data.cloudsql-auditlog_audit_log_rules.example
# }

# resource "cloudsql-auditlog_audit_log_rule" "test" {
#   username = "`mario.finelli`@%"
#   dbname = "*"
#   object = "*"
#   operation = "*"
#   op_result = "E"
# }

resource "cloudsql-auditlog_audit_log_rule" "test2" {
  username  = "*"
  dbname    = "*"
  object    = "*"
  operation = "*"
  op_result = "E"
}

# import {
#   to = cloudsql-auditlog_audit_log_rule.test2
#   id = "7"
# }
