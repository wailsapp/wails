# PROTOTYPE — representative generated file, not a published schema.

version = 3

project {
  name         = "hello"
  product_name = "Hello"
  identifier   = "com.example.hello"
  version      = "0.1.0"
  company      = "Example Ltd"
  binary_name  = "hello"
  icon         = "assets/appicon.png"
}

frontend {
  directory = "frontend"
  install   = ["npm", "install"]
  build     = ["npm", "run", "build"]
  dev       = ["npm", "run", "dev"]
  output    = "frontend/dist"
}

build {
  output = "bin"
}
