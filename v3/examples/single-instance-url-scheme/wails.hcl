version = 3

project {
  icon = "build/appicon.png"
  name = "single-instance-url-scheme"
  product_name = "single-instance-url-scheme"
  identifier = "com.wails.examples.singleinstanceurlscheme"
  version = "0.1.0"
}

frontend {
  disabled = true
}

build {
  tags = ["devtools"]
}

protocol "wails-single-url" {
  description = "Wails example URL handler"
}
