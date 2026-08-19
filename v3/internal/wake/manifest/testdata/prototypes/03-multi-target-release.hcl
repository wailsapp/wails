# PROTOTYPE — representative file, not a published schema.

version = 3

project {
  name         = "orbit"
  product_name = "Orbit"
  identifier   = "com.acme.orbit"
  version      = "2.3.1"
  company      = "Acme Ltd"
  binary_name  = "orbit"
  icon         = "assets/appicon.png"
}

frontend {
  directory = "frontend"
  install   = ["pnpm", "install", "--frozen-lockfile"]
  build     = ["pnpm", "run", "build"]
  dev       = ["pnpm", "run", "dev"]
  output    = "frontend/dist"
}

build {
  output    = "bin"
  trim_path = true
  strip     = true
}

windows {
  publisher = "CN=Acme Ltd"

  signing {
    credential       = "windows-release"
    timestamp_server = "http://timestamp.digicert.com"
  }
}

darwin {
  minimum_version = "12.0"

  signing {
    credential = "apple-release"
    identity   = "Developer ID Application: Acme Ltd"
  }

  notarization {
    credential = "apple-notary"
  }
}

linux {
  signing {
    credential = "linux-packages"
    identity   = "release@acme.example"
  }
}

package "nsis" {
  install_scope = "machine"
}

package "dmg" {
  background    = "assets/dmg-background.png"
  window_width  = 600
  window_height = 420
}

package "appimage" {
  categories = ["Office", "Utility"]
}

package "deb" {
  maintainer   = "Acme Release Team <release@acme.example>"
  section      = "utils"
  dependencies = ["libgtk-4-1", "libwebkitgtk-6.0-4"]
}

profile "release" {
  target "windows/amd64" {
    formats = ["nsis"]
    sign    = true
  }

  target "darwin/universal" {
    formats  = ["dmg"]
    sign     = true
    notarize = true
  }

  target "linux/amd64" {
    formats = ["deb"]
    sign    = true
  }
}
