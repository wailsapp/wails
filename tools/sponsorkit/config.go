package main

// TierStyle controls how sponsors in a tier are drawn. Bigger tiers get
// bigger avatars and fancier ring treatments.
type TierStyle struct {
	// Title shown above the tier section.
	Title string
	// MinMonthly is the minimum monthly dollar amount for this tier.
	// Sponsors are placed in the highest tier whose minimum they meet.
	MinMonthly int
	// Avatar diameter in CSS pixels.
	Avatar float64
	// Box is the horizontal space reserved per sponsor.
	Box float64
	// RowGap is vertical space between rows within the tier.
	RowGap float64
	// ShowName renders the sponsor's name under the avatar.
	ShowName bool
	// NameSize is the font size for the name.
	NameSize float64
	// MaxName is the maximum rune count before the name is ellipsised.
	MaxName int
	// Ring selects the ring treatment: "animated", "static", "subtle" or "".
	Ring string
	// RingWidth is the ring stroke width.
	RingWidth float64
	// RingGradient is the id of the gradient used for the ring stroke.
	RingGradient string
	// Accent is the tier's accent colour (headers, glows).
	Accent string
	// Glow adds a soft blurred halo behind the ring.
	Glow bool
	// Badge renders a small pill with the tier name under the sponsor name.
	Badge string
}

// tiers is ordered from most to least generous. The last entry acts as the
// catch-all for anything below the previous tier's minimum (including
// sponsors whose tier details are not visible to the token).
var tiers = []TierStyle{
	{
		Title: "Partner", MinMonthly: 1000,
		Avatar: 168, Box: 280, RowGap: 30,
		ShowName: true, NameSize: 19, MaxName: 26,
		Ring: "animated", RingWidth: 5, RingGradient: "ringPartner", Accent: "#FF5C5C",
		Glow: true, Badge: "PARTNER",
	},
	{
		Title: "Champion", MinMonthly: 500,
		Avatar: 144, Box: 236, RowGap: 28,
		ShowName: true, NameSize: 17, MaxName: 24,
		Ring: "animated", RingWidth: 4.5, RingGradient: "ringChampion", Accent: "#FF8A5C",
		Glow: true, Badge: "CHAMPION",
	},
	{
		Title: "Gold Sponsors", MinMonthly: 200,
		Avatar: 108, Box: 172, RowGap: 24,
		ShowName: true, NameSize: 14, MaxName: 20,
		Ring: "animated", RingWidth: 4, RingGradient: "ringGold", Accent: "#FFD34D",
		Glow: true,
	},
	{
		Title: "Silver Sponsors", MinMonthly: 100,
		Avatar: 84, Box: 132, RowGap: 22,
		ShowName: true, NameSize: 12.5, MaxName: 16,
		Ring: "static", RingWidth: 3, RingGradient: "ringSilver", Accent: "#C7CEDC",
	},
	{
		Title: "Bronze Sponsors", MinMonthly: 50,
		Avatar: 66, Box: 106, RowGap: 20,
		ShowName: true, NameSize: 11.5, MaxName: 13,
		Ring: "static", RingWidth: 2.5, RingGradient: "ringBronze", Accent: "#E0A878",
	},
	{
		Title: "Covering Costs", MinMonthly: 20,
		Avatar: 52, Box: 66, RowGap: 14,
		Ring: "subtle", RingWidth: 1.5, Accent: "#8B95A9",
	},
	{
		Title: "Buying Breakfast", MinMonthly: 10,
		Avatar: 44, Box: 57, RowGap: 13,
		Ring: "subtle", RingWidth: 1.5, Accent: "#8B95A9",
	},
	{
		Title: "Buying Coffee", MinMonthly: 5,
		Avatar: 37, Box: 49, RowGap: 12,
		Ring: "subtle", RingWidth: 1.25, Accent: "#8B95A9",
	},
	{
		Title: "Helpers", MinMonthly: 0,
		Avatar: 30, Box: 41, RowGap: 11,
		Ring: "subtle", RingWidth: 1, Accent: "#8B95A9",
	},
}
