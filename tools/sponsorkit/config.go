package main

import "image/color"

// cardBackground is the card's midpoint colour. It is composited under
// transparent avatars before JPEG encoding so circle-clipped avatars blend
// into the card. Keep it in sync with the "bg" gradient in render.go
// (#131A2B at the top fading to #0B0F1A).
var cardBackground = color.RGBA{R: 0x11, G: 0x16, B: 0x24, A: 0xff}

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
	// ExtraSparkle adds a second, counter-rotating orbit sparkle.
	ExtraSparkle bool
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
		Glow: true, ExtraSparkle: true,
	},
	{
		Title: "Champion", MinMonthly: 500,
		Avatar: 144, Box: 236, RowGap: 28,
		ShowName: true, NameSize: 17, MaxName: 24,
		Ring: "animated", RingWidth: 4.5, RingGradient: "ringChampion", Accent: "#FF8A5C",
		Glow: true, ExtraSparkle: true,
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

// Band controls how contributors with at least MinCredit credits are drawn
// on the contributors card. Unlike sponsor tiers, bands have no headings:
// the card is one continuous mosaic of squircles graded by size.
type Band struct {
	// MinCredit is the minimum credit count (commits or changelog mentions)
	// for this band. Contributors land in the first band whose minimum
	// they meet.
	MinCredit int
	// Avatar is the squircle's edge length in CSS pixels.
	Avatar float64
	// Box is the horizontal space reserved per contributor.
	Box float64
	// RowGap is vertical space between rows within the band.
	RowGap float64
	// ShowName renders the contributor's login under the squircle.
	ShowName bool
	// NameSize is the font size for the login.
	NameSize float64
	// MaxName is the maximum rune count before the login is ellipsised.
	MaxName int
	// Ring selects the ring treatment: "animated", "static" or "subtle".
	Ring string
	// RingWidth is the ring stroke width.
	RingWidth float64
	// RingGradient is the id of the gradient used for the ring stroke.
	RingGradient string
	// Hover enables the hover power-up (lift, bloom, fast sweep, credit chip).
	Hover bool
}

// bands is ordered from most to least prolific. The last entry is the
// catch-all for first-time and drive-by contributors.
var bands = []Band{
	{
		MinCredit: 1000,
		Avatar:    116, Box: 190, RowGap: 26,
		ShowName: true, NameSize: 15, MaxName: 22,
		Ring: "animated", RingWidth: 4.5, RingGradient: "ringPartner", Hover: true,
	},
	{
		MinCredit: 150,
		Avatar:    94, Box: 152, RowGap: 24,
		ShowName: true, NameSize: 13.5, MaxName: 18,
		Ring: "animated", RingWidth: 4, RingGradient: "ringChampion", Hover: true,
	},
	{
		MinCredit: 60,
		Avatar:    74, Box: 116, RowGap: 22,
		ShowName: true, NameSize: 12, MaxName: 15,
		Ring: "static", RingWidth: 3, RingGradient: "ringGold", Hover: true,
	},
	{
		MinCredit: 20,
		Avatar:    58, Box: 74, RowGap: 16,
		Ring: "static", RingWidth: 2, RingGradient: "ringSilver", Hover: true,
	},
	{
		MinCredit: 8,
		Avatar:    46, Box: 59, RowGap: 13,
		Ring: "subtle", RingWidth: 1.5,
	},
	{
		MinCredit: 3,
		Avatar:    37, Box: 48, RowGap: 12,
		Ring: "subtle", RingWidth: 1.25,
	},
	{
		MinCredit: 0,
		Avatar:    30, Box: 40, RowGap: 11,
		Ring: "subtle", RingWidth: 1,
	},
}
