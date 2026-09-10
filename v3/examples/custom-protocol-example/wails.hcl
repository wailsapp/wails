version = 3

project {
  icon = "build/appicon.png"
  name = "custom-protocol-example"
  product_name = "My Product"
  identifier = "com.wails.examples.customprotocolexample"
  version = "0.0.1"
}

build {
  tags = ["devtools"]
}

protocol "wailsexample" {
  description = "Wails example URL handler"
}
