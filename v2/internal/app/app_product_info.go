package app

import "context"

var (
	ProductName    = "Unset"
	ProductVersion = "Unset"
	CompanyName    = "Unset"
	Copyright      = "Unset"
	Comments       = "Unset"
)

func injectProductInfo(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, "productName", ProductName)
	ctx = context.WithValue(ctx, "productVersion", ProductVersion)
	ctx = context.WithValue(ctx, "companyName", CompanyName)
	ctx = context.WithValue(ctx, "copyright", Copyright)
	return context.WithValue(ctx, "comments", Comments)
}
