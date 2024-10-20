data "external" "dotenv" {
  program = [
    "go",
    "run",
    "scripts/load_dotenv/main.go",
    "../.env",
  ]
}

locals {
  dotenv      = jsondecode(data.external.dotenv)
}
env "local" {
  src = "file://schema.hcl"
  url = "postgres://${local.dotenv.DB_USER}:${local.dotenv.DB_PASSWORD}@${local.dotenv.DB_HOST}:${local.dotenv.DB_PORT}/${local.dotenv.DB_NAME}"
  dev = "docker://postgres/15/dev"
  migration {
    dir = "file://migrations"
  }
}
