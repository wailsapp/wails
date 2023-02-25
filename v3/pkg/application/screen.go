package application

type Screen struct {
	ID        string  `json:"id,omitempty"`         // A unique identifier for the display
	Name      string  `json:"name,omitempty"`       // The name of the display
	Scale     float32 `json:"scale,omitempty"`      // The scale factor of the display
	X         int     `json:"x,omitempty"`          // The x-coordinate of the top-left corner of the rectangle
	Y         int     `json:"y,omitempty"`          // The y-coordinate of the top-left corner of the rectangle
	Size      Size    `json:"size"`                 // The size of the display
	Bounds    Rect    `json:"bounds"`               // The bounds of the display
	WorkArea  Rect    `json:"work_area"`            // The work area of the display
	IsPrimary bool    `json:"is_primary,omitempty"` // Whether this is the primary display
	Rotation  float32 `json:"rotation,omitempty"`   // The rotation of the display
}

type Rect struct {
	X      int `json:"x,omitempty"`
	Y      int `json:"y,omitempty"`
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
}

type Size struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
}
