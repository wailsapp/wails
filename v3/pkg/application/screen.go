package application

type Screen struct {
	ID        string  // A unique identifier for the display
	Name      string  // The name of the display
	Scale     float32 // The scale factor of the display
	X         int     // The x-coordinate of the top-left corner of the rectangle
	Y         int     // The y-coordinate of the top-left corner of the rectangle
	Size      Size    // The size of the display
	Bounds    Rect    // The bounds of the display
	WorkArea  Rect    // The work area of the display
	IsPrimary bool    // Whether this is the primary display
	Rotation  float32 // The rotation of the display
}

type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

type Size struct {
	Width  int
	Height int
}
