# PROTOTYPE — intentionally heavy case, not a published schema.

version = 3

project {
  name         = "studio"
  product_name = "Studio"
  identifier   = "com.acme.studio"
  version      = "5.0.0"
  company      = "Acme Ltd"
  binary_name  = "studio"
  icon         = "assets/appicon.png"
}

frontend {
  directory = "frontend"
  install   = ["pnpm", "install", "--frozen-lockfile"]
  build     = ["./scripts/build-frontend", "--production"]
  dev       = ["pnpm", "run", "dev"]
  output    = "frontend/dist"
}

build {
  output = "artifacts"
}

windows {
  icon     = "assets/windows/app.ico"
  manifest = "assets/windows/app.manifest"

  signing {
    credential       = "windows-release"
    thumbprint       = "0123456789ABCDEF"
    timestamp_server = "http://timestamp.digicert.com"
  }
}

darwin {
  icon         = "assets/macos/app.icns"
  assets_car   = "assets/macos/Assets.car"
  info_plist   = "assets/macos/Info.plist"

  signing {
    entitlements = "assets/macos/entitlements.plist"
  }
}

linux {
  icon          = "assets/linux/studio.png"
  desktop_entry = "assets/linux/studio.desktop"
}

ios {
  info_plist  = "assets/ios/Info.plist"
  assets_car  = "assets/ios/Assets.car"

  signing {
    entitlements = "assets/ios/entitlements.plist"
  }
}

android {
  manifest = "assets/android/AndroidManifest.xml"
}

package "nsis" {
  template = "packaging/windows/installer.nsi"
}

package "msix" {
  manifest = "packaging/windows/AppxManifest.xml"
}

package "dmg" {
  background  = "packaging/macos/background.png"
  volume_icon = "packaging/macos/volume.icns"
  file_icon   = "packaging/macos/file.icns"

  files = {
    "License.pdf" = "LICENSE.pdf"
  }
}

package "appimage" {
  desktop_entry = "assets/linux/studio.desktop"
}

package "deb" {
  pre_install  = "packaging/linux/preinstall.sh"
  post_install = "packaging/linux/postinstall.sh"
  pre_remove   = "packaging/linux/preremove.sh"
  post_remove  = "packaging/linux/postremove.sh"
}

file_association "studio_project" {
  extensions  = ["studio"]
  name        = "Studio Project"
  description = "Studio project file"
  icon        = "assets/filetypes/studio-project.png"
  role        = "editor"
  mime_type   = "application/x-studio-project"
  platforms   = ["windows", "darwin", "linux"]
}

protocol "studio" {
  description = "Open Studio links"
  platforms   = ["windows", "darwin", "linux", "ios", "android"]
}
