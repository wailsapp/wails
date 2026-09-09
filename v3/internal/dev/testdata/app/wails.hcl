version = 3
project {
 name = "wails-dev-acceptance"
 product_name = "Wails Dev Acceptance"
 identifier = "org.wails.devacceptance"
 version = "1.0.0"
}
frontend {
 directory = "frontend"
 install = ["npm", "install"]
 build = ["npm", "run", "build"]
 dev = ["npm", "run", "dev"]
 output = "frontend/dist"
}
dev {
 debounce_ms = 250
 watch = ["**/*.go", "wails.hcl"]
 exclude = [".wails", "frontend", "bin", ".git"]
}
