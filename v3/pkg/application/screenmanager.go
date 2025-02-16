package application

import (
	"errors"
	"math"
	"sort"
)

// Heavily inspired by the Chromium project (Copyright 2015 The Chromium Authors)
// Chromium License: https://chromium.googlesource.com/chromium/src/+/HEAD/LICENSE

type ScreenManager struct {
	screens       []*Screen
	primaryScreen *Screen
}

type Screen struct {
	ID               string  // A unique identifier for the display
	Name             string  // The name of the display
	ScaleFactor      float32 // The scale factor of the display (DPI/96)
	X                int     // The x-coordinate of the top-left corner of the rectangle
	Y                int     // The y-coordinate of the top-left corner of the rectangle
	Size             Size    // The size of the display
	Bounds           Rect    // The bounds of the display
	PhysicalBounds   Rect    // The physical bounds of the display (before scaling)
	WorkArea         Rect    // The work area of the display
	PhysicalWorkArea Rect    // The physical work area of the display (before scaling)
	IsPrimary        bool    // Whether this is the primary display
	Rotation         float32 // The rotation of the display
}

type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

type Point struct {
	X int
	Y int
}
type Size struct {
	Width  int
	Height int
}

type Alignment int
type OffsetReference int

const (
	TOP Alignment = iota
	RIGHT
	BOTTOM
	LEFT
)

const (
	BEGIN OffsetReference = iota // TOP or LEFT
	END                          // BOTTOM or RIGHT
)

// ScreenPlacement specifies where the screen (S) is placed relative to
// parent (P) screen. In the following example, (S) is RIGHT aligned to (P)
// with a positive offset and a BEGIN (top) offset reference.
//
// .           +------------+   +
// .           |            |   | offset
// .           |     P      |   v
// .           |            +--------+
// .           |            |        |
// .           +------------+   S    |
// .                        |        |
// .                        +--------+
type ScreenPlacement struct {
	screen          *Screen
	parent          *Screen
	alignment       Alignment
	offset          int
	offsetReference OffsetReference
}

func (r Rect) Origin() Point {
	return Point{
		X: r.X,
		Y: r.Y,
	}
}

func (s Screen) Origin() Point {
	return Point{
		X: s.X,
		Y: s.Y,
	}
}

func (r Rect) Corner() Point {
	return Point{
		X: r.right(),
		Y: r.bottom(),
	}
}

func (r Rect) InsideCorner() Point {
	return Point{
		X: r.right() - 1,
		Y: r.bottom() - 1,
	}
}

func (r Rect) right() int {
	return r.X + r.Width
}

func (r Rect) bottom() int {
	return r.Y + r.Height
}

func (s Screen) right() int {
	return s.Bounds.right()
}

func (s Screen) bottom() int {
	return s.Bounds.bottom()
}

func (s Screen) scale(value int, toDip bool) int {
	// Round up when scaling down and round down when scaling up.
	// This mix rounding strategy prevents drift over time when applying multiple scaling back and forth.
	// In addition, It has been shown that using this approach minimized rounding issues and improved overall
	// precision when converting between DIP and physical coordinates.
	if toDip {
		return int(math.Ceil(float64(value) / float64(s.ScaleFactor)))
	} else {
		return int(math.Floor(float64(value) * float64(s.ScaleFactor)))
	}
}

func (r Rect) Size() Size {
	return Size{
		Width:  r.Width,
		Height: r.Height,
	}
}

func (r Rect) IsEmpty() bool {
	return r.Width <= 0 || r.Height <= 0
}

func (r Rect) Contains(pt Point) bool {
	return pt.X >= r.X && pt.X < r.X+r.Width && pt.Y >= r.Y && pt.Y < r.Y+r.Height
}

// Get intersection with another rect
func (r Rect) Intersect(otherRect Rect) Rect {
	if r.IsEmpty() || otherRect.IsEmpty() {
		return Rect{}
	}

	maxLeft := max(r.X, otherRect.X)
	maxTop := max(r.Y, otherRect.Y)
	minRight := min(r.right(), otherRect.right())
	minBottom := min(r.bottom(), otherRect.bottom())

	if minRight > maxLeft && minBottom > maxTop {
		return Rect{
			X:      maxLeft,
			Y:      maxTop,
			Width:  minRight - maxLeft,
			Height: minBottom - maxTop,
		}
	}
	return Rect{}
}

// Check if screens intersects another screen
func (s *Screen) intersects(otherScreen *Screen) bool {
	maxLeft := max(s.X, otherScreen.X)
	maxTop := max(s.Y, otherScreen.Y)
	minRight := min(s.right(), otherScreen.right())
	minBottom := min(s.bottom(), otherScreen.bottom())

	return minRight > maxLeft && minBottom > maxTop
}

// Get distance from another rect (squared)
func (r Rect) distanceFromRectSquared(otherRect Rect) int {
	// If they intersect, return negative area of intersection
	intersection := r.Intersect(otherRect)
	if !intersection.IsEmpty() {
		return -(intersection.Width * intersection.Height)
	}

	dX := max(0, max(r.X-otherRect.right(), otherRect.X-r.right()))
	dY := max(0, max(r.Y-otherRect.bottom(), otherRect.Y-r.bottom()))

	// Distance squared
	return dX*dX + dY*dY
}

// Apply screen placement
func (p ScreenPlacement) apply() {
	parentBounds := p.parent.Bounds
	screenBounds := p.screen.Bounds

	newX := parentBounds.X
	newY := parentBounds.Y
	offset := p.offset

	if p.alignment == TOP || p.alignment == BOTTOM {
		if p.offsetReference == END {
			offset = parentBounds.Width - offset - screenBounds.Width
		}
		offset = min(offset, parentBounds.Width)
		offset = max(offset, -screenBounds.Width)
		newX += offset
		if p.alignment == TOP {
			newY -= screenBounds.Height
		} else {
			newY += parentBounds.Height
		}
	} else {
		if p.offsetReference == END {
			offset = parentBounds.Height - offset - screenBounds.Height
		}
		offset = min(offset, parentBounds.Height)
		offset = max(offset, -screenBounds.Height)
		newY += offset
		if p.alignment == LEFT {
			newX -= screenBounds.Width
		} else {
			newX += parentBounds.Width
		}
	}

	p.screen.move(newX, newY)
}

func (s *Screen) absoluteToRelativeDipPoint(dipPoint Point) Point {
	return Point{
		X: dipPoint.X - s.Bounds.X,
		Y: dipPoint.Y - s.Bounds.Y,
	}
}

func (s *Screen) relativeToAbsoluteDipPoint(dipPoint Point) Point {
	return Point{
		X: dipPoint.X + s.Bounds.X,
		Y: dipPoint.Y + s.Bounds.Y,
	}
}

func (s *Screen) absoluteToRelativePhysicalPoint(physicalPoint Point) Point {
	return Point{
		X: physicalPoint.X - s.PhysicalBounds.X,
		Y: physicalPoint.Y - s.PhysicalBounds.Y,
	}
}

func (s *Screen) relativeToAbsolutePhysicalPoint(physicalPoint Point) Point {
	return Point{
		X: physicalPoint.X + s.PhysicalBounds.X,
		Y: physicalPoint.Y + s.PhysicalBounds.Y,
	}
}

func (s *Screen) move(newX, newY int) {
	workAreaOffsetX := s.WorkArea.X - s.X
	workAreaOffsetY := s.WorkArea.Y - s.Y

	s.X = newX
	s.Y = newY
	s.Bounds.X = newX
	s.Bounds.Y = newY
	s.WorkArea.X = newX + workAreaOffsetX
	s.WorkArea.Y = newY + workAreaOffsetY
}

func (s *Screen) applyDPIScaling() {
	if s.ScaleFactor == 1 {
		return
	}
	workAreaOffsetX := s.WorkArea.X - s.Bounds.X
	workAreaOffsetY := s.WorkArea.Y - s.Bounds.Y

	s.WorkArea.X = s.Bounds.X + s.scale(workAreaOffsetX, true)
	s.WorkArea.Y = s.Bounds.Y + s.scale(workAreaOffsetY, true)

	s.Bounds.Width = s.scale(s.PhysicalBounds.Width, true)
	s.Bounds.Height = s.scale(s.PhysicalBounds.Height, true)
	s.WorkArea.Width = s.scale(s.PhysicalWorkArea.Width, true)
	s.WorkArea.Height = s.scale(s.PhysicalWorkArea.Height, true)

	s.Size.Width = s.Bounds.Width
	s.Size.Height = s.Bounds.Height
}

func (s *Screen) dipToPhysicalPoint(dipPoint Point, isCorner bool) Point {
	relativePoint := s.absoluteToRelativeDipPoint(dipPoint)
	scaledRelativePoint := Point{
		X: s.scale(relativePoint.X, false),
		Y: s.scale(relativePoint.Y, false),
	}
	// Align edge points (fixes rounding issues)
	edgeOffset := 1
	if isCorner {
		edgeOffset = 0
	}
	if relativePoint.X == s.Bounds.Width-edgeOffset {
		scaledRelativePoint.X = s.PhysicalBounds.Width - edgeOffset
	}
	if relativePoint.Y == s.Bounds.Height-edgeOffset {
		scaledRelativePoint.Y = s.PhysicalBounds.Height - edgeOffset
	}
	return s.relativeToAbsolutePhysicalPoint(scaledRelativePoint)
}

func (s *Screen) physicalToDipPoint(physicalPoint Point, isCorner bool) Point {
	relativePoint := s.absoluteToRelativePhysicalPoint(physicalPoint)
	scaledRelativePoint := Point{
		X: s.scale(relativePoint.X, true),
		Y: s.scale(relativePoint.Y, true),
	}
	// Align edge points (fixes rounding issues)
	edgeOffset := 1
	if isCorner {
		edgeOffset = 0
	}
	if relativePoint.X == s.PhysicalBounds.Width-edgeOffset {
		scaledRelativePoint.X = s.Bounds.Width - edgeOffset
	}
	if relativePoint.Y == s.PhysicalBounds.Height-edgeOffset {
		scaledRelativePoint.Y = s.Bounds.Height - edgeOffset
	}
	return s.relativeToAbsoluteDipPoint(scaledRelativePoint)
}

func (s *Screen) dipToPhysicalRect(dipRect Rect) Rect {
	origin := s.dipToPhysicalPoint(dipRect.Origin(), false)
	corner := s.dipToPhysicalPoint(dipRect.Corner(), true)

	return Rect{
		X:      origin.X,
		Y:      origin.Y,
		Width:  corner.X - origin.X,
		Height: corner.Y - origin.Y,
	}
}

func (s *Screen) physicalToDipRect(physicalRect Rect) Rect {
	origin := s.physicalToDipPoint(physicalRect.Origin(), false)
	corner := s.physicalToDipPoint(physicalRect.Corner(), true)

	return Rect{
		X:      origin.X,
		Y:      origin.Y,
		Width:  corner.X - origin.X,
		Height: corner.Y - origin.Y,
	}
}

// Layout screens in the virtual space with DIP calculations and cache the screens
// for future coordinate transformation between the physical and logical (DIP) space
func (m *ScreenManager) LayoutScreens(screens []*Screen) error {
	if screens == nil || len(screens) == 0 {
		return errors.New("screens parameter is nil or empty")
	}
	m.screens = screens

	err := m.calculateScreensDipCoordinates()
	if err != nil {
		return err
	}

	return nil
}

func (m *ScreenManager) Screens() []*Screen {
	return m.screens
}

func (m *ScreenManager) PrimaryScreen() *Screen {
	return m.primaryScreen
}

// Reference: https://source.chromium.org/chromium/chromium/src/+/main:ui/display/win/screen_win.cc;l=317
func (m *ScreenManager) calculateScreensDipCoordinates() error {
	remainingScreens := []*Screen{}

	// Find the primary screen
	m.primaryScreen = nil
	for _, screen := range m.screens {
		if screen.IsPrimary {
			m.primaryScreen = screen
		} else {
			remainingScreens = append(remainingScreens, screen)
		}
	}
	if m.primaryScreen == nil {
		return errors.New("no primary screen found")
	} else if len(remainingScreens) != len(m.screens)-1 {
		return errors.New("invalid primary screen found")
	}

	// Build screens tree using the primary screen as root
	screensPlacements := []ScreenPlacement{}
	availableParents := []*Screen{m.primaryScreen}
	for len(availableParents) > 0 {
		// Pop a parent
		end := len(availableParents) - 1
		parent := availableParents[end]
		availableParents = availableParents[:end]
		// Find touching screens
		for _, child := range m.findAndRemoveTouchingScreens(parent, &remainingScreens) {
			screenPlacement := m.calculateScreenPlacement(child, parent)
			screensPlacements = append(screensPlacements, screenPlacement)
			availableParents = append(availableParents, child)
		}
	}

	// Apply screens DPI scaling and placement starting with
	// the primary screen and then dependent screens
	m.primaryScreen.applyDPIScaling()
	for _, placement := range screensPlacements {
		placement.screen.applyDPIScaling()
		placement.apply()
	}

	// Now that all the placements have been applied,
	// we must detect and fix any overlapping screens.
	m.deIntersectScreens(screensPlacements)

	return nil
}

// Returns a ScreenPlacement for |screen| relative to |parent|.
// Note that ScreenPlacement's are always in DIPs, so this also performs the
// required scaling.
// References:
//   - https://github.com/chromium/chromium/blob/main/ui/display/win/scaling_util.h#L25
//   - https://github.com/chromium/chromium/blob/main/ui/display/win/scaling_util.cc#L142
func (m *ScreenManager) calculateScreenPlacement(screen, parent *Screen) ScreenPlacement {
	// Examples (The offset is indicated by the arrow.):
	// Scaled and Unscaled Coordinates
	// +--------------+    +          Since both screens are of the same scale
	// |              |    |          factor, relative positions remain the same.
	// |    Parent    |    V
	// |      1x      +----------+
	// |              |          |
	// +--------------+  Screen  |
	//                |    1x    |
	//                +----------+
	//
	// Unscaled Coordinates
	// +--------------+               The 2x screen is offset to maintain a
	// |              |               similar neighboring relationship with the 1x
	// |    Parent    |               parent. Screen's position is based off of the
	// |      1x      +----------+    percentage position along its parent. This
	// |              |          |    percentage position is preserved in the scaled
	// +--------------+  Screen  |    coordinates.
	//                |    2x    |
	//                +----------+
	// Scaled Coordinates
	// +--------------+  +
	// |              |  |
	// |    Parent    |  V
	// |      1x      +-----+
	// |              | S 2x|
	// +--------------+-----+
	//
	//
	// Unscaled Coordinates
	// +--------------+               The parent screen has a 2x scale factor.
	// |              |               The offset is adjusted to maintain the
	// |              |               relative positioning of the 1x screen in
	// |    Parent    +----------+    the scaled coordinate space. Screen's
	// |      2x      |          |    position is based off of the percentage
	// |              |  Screen  |    position along its parent. This percentage
	// |              |    1x    |    position is preserved in the scaled
	// +--------------+          |    coordinates.
	//                |          |
	//                +----------+
	// Scaled Coordinates
	// +-------+    +
	// |       |    V
	// | Parent+----------+
	// |   2x  |          |
	// +-------+  Screen  |
	//         |    1x    |
	//         |          |
	//         |          |
	//         +----------+
	//
	// Unscaled Coordinates
	//         +----------+           In this case, parent lies between the top and
	//         |          |           bottom of parent. The roles are reversed when
	// +-------+          |           this occurs, and screen is placed to maintain
	// |       |  Screen  |           parent's relative position along screen.
	// | Parent|    1x    |
	// |   2x  |          |
	// +-------+          |
	//         +----------+
	// Scaled Coordinates
	//  ^      +----------+
	//  |      |          |
	//  + +----+          |
	//    |Prnt|  Screen  |
	//    | 2x |    1x    |
	//    +----+          |
	//         |          |
	//         +----------+
	//
	// Scaled and Unscaled Coordinates
	// +--------+                     If the two screens are bottom aligned or
	// |        |                     right aligned, the ScreenPlacement will
	// |        +--------+            have an offset of 0 relative to the
	// |        |        |            end of the screen.
	// |        |        |
	// +--------+--------+

	placement := ScreenPlacement{
		screen:          screen,
		parent:          parent,
		alignment:       m.getScreenAlignment(screen, parent),
		offset:          0,
		offsetReference: BEGIN,
	}

	screenBegin, screenEnd := 0, 0
	parentBegin, parentEnd := 0, 0

	switch placement.alignment {
	case TOP, BOTTOM:
		screenBegin = screen.X
		screenEnd = screen.right()
		parentBegin = parent.X
		parentEnd = parent.right()
	case LEFT, RIGHT:
		screenBegin = screen.Y
		screenEnd = screen.bottom()
		parentBegin = parent.Y
		parentEnd = parent.bottom()
	}

	// Since we're calculating offsets, make everything relative to parentBegin
	parentEnd -= parentBegin
	screenBegin -= parentBegin
	screenEnd -= parentBegin
	parentBegin = 0

	// There are a few ways lines can intersect:
	// End Aligned
	// SCREEN's offset is relative to the END (BOTTOM or RIGHT).
	//                 +-PARENT----------------+
	//                    +-SCREEN-------------+
	//
	// Positioning based off of |screenBegin|.
	// SCREEN's offset is simply a percentage of its position on PARENT.
	//                 +-PARENT----------------+
	//                        ^+-SCREEN------------+
	//
	// Positioning based off of |screenEnd|.
	// SCREEN's offset is dependent on the percentage of its end position on PARENT.
	//                 +-PARENT----------------+
	//           +-SCREEN------------+^
	//
	// Positioning based off of |parentBegin| on SCREEN.
	// SCREEN's offset is dependent on the percentage of its position on PARENT.
	//                 +-PARENT----------------+
	//          ^+-SCREEN--------------------------+

	if screenEnd == parentEnd {
		placement.offsetReference = END
		placement.offset = 0
	} else if screenBegin >= parentBegin {
		placement.offsetReference = BEGIN
		placement.offset = m.scaleOffset(parentEnd, parent.ScaleFactor, screenBegin)
	} else if screenEnd <= parentEnd {
		placement.offsetReference = END
		placement.offset = m.scaleOffset(parentEnd, parent.ScaleFactor, parentEnd-screenEnd)
	} else {
		placement.offsetReference = BEGIN
		placement.offset = m.scaleOffset(screenEnd-screenBegin, screen.ScaleFactor, screenBegin)
	}

	return placement
}

// Get screen alignment relative to parent (TOP, RIGHT, BOTTOM, LEFT)
func (m *ScreenManager) getScreenAlignment(screen, parent *Screen) Alignment {
	maxLeft := max(screen.X, parent.X)
	maxTop := max(screen.Y, parent.Y)
	minRight := min(screen.right(), parent.right())
	minBottom := min(screen.bottom(), parent.bottom())

	// Corners touching
	if maxLeft == minRight && maxTop == minBottom {
		if screen.Y == maxTop {
			return BOTTOM
		} else if parent.X == maxLeft {
			return LEFT
		}
		return TOP
	}

	// Vertical edge touching
	if maxLeft == minRight {
		if screen.X == maxLeft {
			return RIGHT
		} else {
			return LEFT
		}
	}

	// Horizontal edge touching
	if maxTop == minBottom {
		if screen.Y == maxTop {
			return BOTTOM
		} else {
			return TOP
		}
	}

	return -1 // Shouldn't be reached
}

func (m *ScreenManager) deIntersectScreens(screensPlacements []ScreenPlacement) {
	parentIDMap := make(map[string]string)
	for _, placement := range screensPlacements {
		parentIDMap[placement.screen.ID] = placement.parent.ID
	}

	treeDepthMap := make(map[string]int)
	for _, screen := range m.screens {
		id, ok, depth := screen.ID, true, 0
		const maxDepth = 100
		for id != m.primaryScreen.ID && depth < maxDepth {
			depth++
			id, ok = parentIDMap[id]
			if !ok {
				depth = maxDepth
			}
		}
		treeDepthMap[screen.ID] = depth
	}

	sortedScreens := make([]*Screen, len(m.screens))
	copy(sortedScreens, m.screens)

	// Sort the screens first by their depth in the screen hierarchy tree,
	// and then by distance from screen origin to primary origin. This way we
	// process the screens starting at the root (the primary screen), in the
	// order of their descendance spanning out from the primary screen.
	sort.Slice(sortedScreens, func(i, j int) bool {
		s1, s2 := m.screens[i], m.screens[j]
		s1_depth := treeDepthMap[s1.ID]
		s2_depth := treeDepthMap[s2.ID]

		if s1_depth != s2_depth {
			return s1_depth < s2_depth
		}

		// Distance squared
		s1_distance := s1.X*s1.X + s1.Y*s1.Y
		s2_distance := s2.X*s2.X + s2.Y*s2.Y
		if s1_distance != s2_distance {
			return s1_distance < s2_distance
		}

		return s1.ID < s2.ID
	})

	for i := 1; i < len(sortedScreens); i++ {
		targetScreen := sortedScreens[i]
		for j := 0; j < i; j++ {
			sourceScreen := sortedScreens[j]
			if targetScreen.intersects(sourceScreen) {
				m.fixScreenIntersection(targetScreen, sourceScreen)
			}
		}
	}
}

// Offset the target screen along either X or Y axis away from the origin
// so that it removes the intersection with the source screen
// This function assume both screens already intersect.
func (m *ScreenManager) fixScreenIntersection(targetScreen, sourceScreen *Screen) {
	offsetX, offsetY := 0, 0

	if targetScreen.X >= 0 {
		offsetX = sourceScreen.right() - targetScreen.X
	} else {
		offsetX = -(targetScreen.right() - sourceScreen.X)
	}

	if targetScreen.Y >= 0 {
		offsetY = sourceScreen.bottom() - targetScreen.Y
	} else {
		offsetY = -(targetScreen.bottom() - sourceScreen.Y)
	}

	// Choose the smaller offset (X or Y)
	if math.Abs(float64(offsetX)) <= math.Abs(float64(offsetY)) {
		offsetY = 0
	} else {
		offsetX = 0
	}

	// Apply the offset
	newX := targetScreen.X + offsetX
	newY := targetScreen.Y + offsetY
	targetScreen.move(newX, newY)
}

func (m *ScreenManager) findAndRemoveTouchingScreens(parent *Screen, screens *[]*Screen) []*Screen {
	touchingScreens := []*Screen{}
	remainingScreens := []*Screen{}

	for _, screen := range *screens {
		if m.areScreensTouching(parent, screen) {
			touchingScreens = append(touchingScreens, screen)
		} else {
			remainingScreens = append(remainingScreens, screen)
		}
	}
	*screens = remainingScreens
	return touchingScreens
}

func (m *ScreenManager) areScreensTouching(a, b *Screen) bool {
	maxLeft := max(a.X, b.X)
	maxTop := max(a.Y, b.Y)
	minRight := min(a.right(), b.right())
	minBottom := min(a.bottom(), b.bottom())
	return (maxLeft == minRight && maxTop <= minBottom) || (maxTop == minBottom && maxLeft <= minRight)
}

// Scale |unscaledOffset| to the same relative position on |unscaledLength|
// based off of |unscaledLength|'s |scaleFactor|
func (m *ScreenManager) scaleOffset(unscaledLength int, scaleFactor float32, unscaledOffset int) int {
	scaledLength := float32(unscaledLength) / scaleFactor
	percent := float32(unscaledOffset) / float32(unscaledLength)
	return int(math.Floor(float64(scaledLength * percent)))
}

func (m *ScreenManager) screenNearestPoint(point Point, isPhysical bool) *Screen {
	for _, screen := range m.screens {
		if isPhysical {
			if screen.PhysicalBounds.Contains(point) {
				return screen
			}
		} else {
			if screen.Bounds.Contains(point) {
				return screen
			}
		}
	}
	return m.primaryScreen
}

func (m *ScreenManager) screenNearestRect(rect Rect, isPhysical bool, excludedScreens map[string]bool) *Screen {
	var nearestScreen *Screen
	var distance, nearestScreenDistance int
	for _, screen := range m.screens {
		if excludedScreens[screen.ID] {
			continue
		}
		if isPhysical {
			distance = rect.distanceFromRectSquared(screen.PhysicalBounds)
		} else {
			distance = rect.distanceFromRectSquared(screen.Bounds)
		}
		if nearestScreen == nil || distance < nearestScreenDistance {
			nearestScreen = screen
			nearestScreenDistance = distance
		}
	}
	if !isPhysical && len(excludedScreens) < len(m.screens)-1 {
		// Make sure to give the same screen that would be given by the physical rect
		// of this dip rect so transforming back and forth always gives the same result.
		// This is important because it could happen that a dip rect intersects Screen1
		// more than Screen2 but in the physical layout Screen2 will scale up or Screen1
		// will scale down causing the intersection area to change so transforming back
		// would give a different rect.
		physicalRect := nearestScreen.dipToPhysicalRect(rect)
		physicalRectScreen := m.screenNearestRect(physicalRect, true, nil)
		if nearestScreen != physicalRectScreen {
			if excludedScreens == nil {
				excludedScreens = make(map[string]bool)
			}
			excludedScreens[nearestScreen.ID] = true
			return m.screenNearestRect(rect, isPhysical, excludedScreens)
		}
	}
	return nearestScreen
}

func (m *ScreenManager) DipToPhysicalPoint(dipPoint Point) Point {
	screen := m.ScreenNearestDipPoint(dipPoint)
	return screen.dipToPhysicalPoint(dipPoint, false)
}

func (m *ScreenManager) PhysicalToDipPoint(physicalPoint Point) Point {
	screen := m.ScreenNearestPhysicalPoint(physicalPoint)
	return screen.physicalToDipPoint(physicalPoint, false)
}

func (m *ScreenManager) DipToPhysicalRect(dipRect Rect) Rect {
	screen := m.ScreenNearestDipRect(dipRect)
	return screen.dipToPhysicalRect(dipRect)
}

func (m *ScreenManager) PhysicalToDipRect(physicalRect Rect) Rect {
	screen := m.ScreenNearestPhysicalRect(physicalRect)
	return screen.physicalToDipRect(physicalRect)
}

func (m *ScreenManager) ScreenNearestPhysicalPoint(physicalPoint Point) *Screen {
	return m.screenNearestPoint(physicalPoint, true)
}

func (m *ScreenManager) ScreenNearestDipPoint(dipPoint Point) *Screen {
	return m.screenNearestPoint(dipPoint, false)
}

func (m *ScreenManager) ScreenNearestPhysicalRect(physicalRect Rect) *Screen {
	return m.screenNearestRect(physicalRect, true, nil)
}

func (m *ScreenManager) ScreenNearestDipRect(dipRect Rect) *Screen {
	return m.screenNearestRect(dipRect, false, nil)
}

// ================================================================================================
// Exported application-level methods for internal convenience and availability to application devs

func DipToPhysicalPoint(dipPoint Point) Point {
	return globalApplication.screenManager.DipToPhysicalPoint(dipPoint)
}

func PhysicalToDipPoint(physicalPoint Point) Point {
	return globalApplication.screenManager.PhysicalToDipPoint(physicalPoint)
}

func DipToPhysicalRect(dipRect Rect) Rect {
	return globalApplication.screenManager.DipToPhysicalRect(dipRect)
}

func PhysicalToDipRect(physicalRect Rect) Rect {
	return globalApplication.screenManager.PhysicalToDipRect(physicalRect)
}

func ScreenNearestPhysicalPoint(physicalPoint Point) *Screen {
	return globalApplication.screenManager.ScreenNearestPhysicalPoint(physicalPoint)
}

func ScreenNearestDipPoint(dipPoint Point) *Screen {
	return globalApplication.screenManager.ScreenNearestDipPoint(dipPoint)
}

func ScreenNearestPhysicalRect(physicalRect Rect) *Screen {
	return globalApplication.screenManager.ScreenNearestPhysicalRect(physicalRect)
}

func ScreenNearestDipRect(dipRect Rect) *Screen {
	return globalApplication.screenManager.ScreenNearestDipRect(dipRect)
}
