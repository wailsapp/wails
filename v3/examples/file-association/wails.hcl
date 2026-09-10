version = 3

project {
  icon = "build/appicon.png"
  name = "file-association"
  product_name = "My Product"
  identifier = "com.wails.examples.fileassociation"
  version = "0.0.1"
}

build {
  tags = ["devtools"]
}

file_association "wails" {
  extensions = ["wails"]
  description = "Wails Application File"
  role = "Editor"
}
