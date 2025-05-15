terraform {
  required_providers {
    cloudsql-auditlog = {
      source = "facile.it/test/cloudsql-auditlog"
    }
  }
}

provider "cloudsql-auditlog" {
  endpoint = "test"
  password = "testp"
  username = "testu"
}

data "cloudsql-auditlog_coffees" "example" {}
