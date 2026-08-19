# PROTOTYPE — representative file, not a published schema.

version = 3

project {
  name         = "ledger"
  product_name = "Ledger"
  identifier   = "io.acme.ledger"
  version      = "1.4.0"
  company      = "Acme Ltd"
  binary_name  = "ledger"
  icon         = "assets/appicon.png"
}

frontend {
  directory = "frontend"
  install   = ["pnpm", "install", "--frozen-lockfile"]
  build     = ["./scripts/build-frontend"]
  dev       = ["pnpm", "run", "dev"]
  output    = "frontend/dist"

  environment = {
    PUBLIC_EDITION = "community"
  }
}

build {
  output = "dist"
  tags   = ["sqlite_fts5"]
}

dev {
  debounce_ms    = 500
  log_level      = "info"
  watch          = ["**/*.go", "frontend/src/**"]
  exclude        = ["dist/**", "frontend/node_modules/**"]
  use_git_ignore = true
}

target "windows/amd64" {
  toolchain = "zig"
  tags      = ["enterprise"]
}

package "nsis" {
  install_scope = "user"
}

profile "release" {
  target "windows/amd64" {
    formats = ["nsis"]
  }
}
