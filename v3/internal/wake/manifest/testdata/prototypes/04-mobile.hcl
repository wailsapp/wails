# PROTOTYPE — representative file, not a published schema.

version = 3

project {
  name         = "field-notes"
  product_name = "Field Notes"
  identifier   = "com.acme.fieldnotes"
  version      = "3.2.0"
  company      = "Acme Ltd"
  binary_name  = "field-notes"
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
  output = "bin"
}

ios {
  bundle_id        = "com.acme.fieldnotes"
  display_name     = "Field Notes"
  build_number     = 87
  minimum_version  = "15.0"
  background_modes = ["fetch", "remote-notification"]

  signing {
    credential           = "apple-development"
    identity             = "Apple Distribution: Acme Ltd"
    provisioning_profile = "assets/FieldNotes.mobileprovision"
    entitlements         = "assets/ios-entitlements.plist"
  }
}

android {
  application_id     = "com.acme.fieldnotes"
  display_name       = "Field Notes"
  version_name       = "3.2.0"
  version_code       = 87
  minimum_sdk        = 26
  target_sdk         = 35

  signing {
    credential = "google-play-upload"
    key_alias  = "field-notes-upload"
  }
}

profile "app-store" {
  target "ios/arm64" {
    destination = "device"
    formats     = ["ipa"]
    sign        = true
  }
}

profile "play-store" {
  target "android/universal" {
    formats = ["aab"]
    sign    = true
  }
}
