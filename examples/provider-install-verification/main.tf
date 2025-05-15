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
}

data "cloudsql-auditlog_coffees" "example" {}
