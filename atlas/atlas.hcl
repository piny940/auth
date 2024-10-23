data "external" "dotenv" {
  program = [
    "go",
    "run",
    "scripts/load_dotenv/main.go",
    "../.env",
  ]
}
data "external" "dotenv_test" {
  program = [
    "go",
    "run",
    "scripts/load_dotenv/main.go",
    "../.env.test",
  ]
}

locals {
  dotenv      = jsondecode(data.external.dotenv)
  dotenv_test = jsondecode(data.external.dotenv_test)
}
env "local" {
  src = "file://schema.hcl"
  url = "postgres://${local.dotenv.DB_USER}:${local.dotenv.DB_PASSWORD}@${local.dotenv.DB_HOST}:${local.dotenv.DB_PORT}/${local.dotenv.DB_NAME}?sslmode=${local.dotenv.DB_SSLMODE}"
  dev = "docker://postgres/15/dev"
  migration {
    dir = "file://migrations"
  }
}
env "test" {
  src = "file://schema.hcl"
  url = "postgres://${local.dotenv_test.DB_USER}:${local.dotenv_test.DB_PASSWORD}@${local.dotenv_test.DB_HOST}:${local.dotenv_test.DB_PORT}/${local.dotenv_test.DB_NAME}?sslmode=${local.dotenv_test.DB_SSLMODE}"
  dev = "docker://postgres/15/dev"
  migration {
    dir = "file://migrations"
  }
}
