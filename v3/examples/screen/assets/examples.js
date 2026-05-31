window.examples = [
  [
    // Normal examples (demonstrate real life scenarios)
    {
      name: "Single 4k monitor",
      screens: [
        {id: 1, w: 3840, h: 2160, s: 163.0 / 96, name: `27" 4K UHD 163DPI`},
      ]
    },
    {
      name: "Two monitors",
      screens: [
        {id: 1, w: 3840, h: 2160, s: 163.0 / 96, name: `27" 4K UHD 163DPI`},
        {id: 2, w: 1920, h: 1080, s: 1, parent: {id: 1, align: "r", offset: 0}, name: `23" FHD 96DPI`},
      ]
    },
    {
      name: "Two monitors (2)",
      screens: [
        {id: 1, w: 1920, h: 1080, s: 1, name: `23" FHD 96DPI`},
        {id: 2, w: 1920, h: 1080, s: 1.25, parent: {id: 1, align: "r", offset: 0}, name: `23" FHD 96DPI (125%)`},
      ]
    },
    {
      name: "Three monitors",
      screens: [
        {id: 1, w: 3840, h: 2160, s: 163.0 / 96, name: `27" 4K UHD 163DPI`},
        {id: 2, w: 1920, h: 1080, s: 1, parent: {id: 1, align: "r", offset: 0}, name: `23" FHD 96DPI`},
        {id: 3, w: 1920, h: 1080, s: 1.25, parent: {id: 1, align: "l", offset: 0}, name: `23" FHD 96DPI (125%)`},
      ]
    },
    {
      name: "Four monitors",
      screens: [
        {id: 1, w: 3840, h: 2160, s: 163.0 / 96, name: `27" 4K UHD 163DPI`},
        {id: 2, w: 1920, h: 1080, s: 1, parent: {id: 1, align: "r", offset: 0}, name: `23" FHD 96DPI`},
        {id: 3, w: 1920, h: 1080, s: 1.25, parent: {id: 2, align: "b", offset: 0}, name: `23" FHD 96DPI (125%)`},
        {id: 4, w: 1080, h: 1920, s: 1, parent: {id: 1, align: "l", offset: 0}, name: `23" FHD (90deg)`},
      ]
    },
  ],
  [
    // Test cases examples (demonstrate the algorithm basics)
    {
      name: "Child scaled, Start offset",
      screens: [
        {id: 1, w: 1200, h: 1200, s: 1, name: "Parent"},
        {id: 2, w: 1200, h: 1200, s: 1.5, parent: {id: 1, align: "r", offset: 600}, name: "Child"},
      ]
    },
    {
      name: "Child scaled, End offset",
      screens: [
        {id: 1, w: 1200, h: 1200, s: 1, name: "Parent"},
        {id: 2, w: 1200, h: 1200, s: 1.5, parent: {id: 1, align: "r", offset: -600}, name: "Child"},
      ]
    },
    {
      name: "Parent scaled, Start offset percent",
      screens: [
        {id: 1, w: 1200, h: 1200, s: 1.5, name: "Parent"},
        {id: 2, w: 1200, h: 1200, s: 1, parent: {id: 1, align: "r", offset: 600}, name: "Child"},
      ]
    },
    {
      name: "Parent scaled, End offset percent",
      screens: [
        {id: 1, w: 1200, h: 1200, s: 1.5, name: "Parent"},
        {id: 2, w: 1200, h: 1200, s: 1, parent: {id: 1, align: "r", offset: -600}, name: "Child"},
      ]
    },
    {
      name: "Parent scaled, Start align",
      screens: [
        {id: 1, w: 1200, h: 1200, s: 1.5, name: "Parent"},
        {id: 2, w: 1200, h: 1100, s: 1, parent: {id: 1, align: "r", offset: 0}, name: "Child"},
      ]
    },
    {
      name: "Parent scaled, End align",
      screens: [
        {id: 1, w: 1200, h: 1200, s: 1.5, name: "Parent"},
        {id: 2, w: 1200, h: 1200, s: 1, parent: {id: 1, align: "r", offset: 0}, name: "Child"},
      ]
    },
    {
      name: "Parent scaled, in-between",
      screens: [
        {id: 1, w: 1200, h: 1200, s: 1.5, name: "Parent"},
        {id: 2, w: 1200, h: 1500, s: 1, parent: {id: 1, align: "r", offset: -200}, name: "Child"},
      ]
    },
  ],
  [
    // Edge cases examples
    {
      name: "Parent order (5 is parent of 4)",
      screens: [
        {id: 1, w: 1920, h: 1080, s: 1},
        {id: 2, w: 1024, h: 600, s: 1.25, parent: {id: 1, align: "r", offset: -200}},
        {id: 3, w: 800, h: 800, s: 1.25, parent: {id: 2, align: "b", offset: 0}},
        {id: 4, w: 800, h: 1080, s: 1.5, parent: {id: 2, align: "re", offset: 100}},
        {id: 5, w: 600, h: 600, s: 1, parent: {id: 3, align: "r", offset: 100}},
      ]
    },
    {
      name: "de-intersection reparent",
      screens: [
        {id: 1, w: 1920, h: 1080, s: 1},
        {id: 2, w: 1680, h: 1050, s: 1.25, parent: {id: 1, align: "r", offset: 10}},
        {id: 3, w: 1440, h: 900, s: 1.5, parent: {id: 1, align: "le", offset: 150}},
        {id: 4, w: 1024, h: 768, s: 1, parent: {id: 3, align: "bc", offset: -200}},
        {id: 5, w: 1024, h: 768, s: 1.25, parent: {id: 4, align: "r", offset: 400}},
      ]
    },
    {
      name: "de-intersection (unattached child)",
      screens: [
        {id: 1, w: 1920, h: 1080, s: 1},
        {id: 2, w: 1024, h: 768, s: 1.5, parent: {id: 1, align: "le", offset: 10}},
        {id: 3, w: 1024, h: 768, s: 1.25, parent: {id: 2, align: "b", offset: 100}},
        {id: 4, w: 1024, h: 768, s: 1, parent: {id: 3, align: "r", offset: 500}},
      ]
    },
    {
      name: "Multiple de-intersection",
      screens: [
        {id: 1, w: 1920, h: 1080, s: 1},
        {id: 2, w: 1024, h: 768, s: 1, parent: {id: 1, align: "be", offset: 0}},
        {id: 3, w: 1024, h: 768, s: 1, parent: {id: 2, align: "b", offset: 300}},
        {id: 4, w: 1024, h: 768, s: 1.5, parent: {id: 2, align: "le", offset: 100}},
        {id: 5, w: 1024, h: 768, s: 1, parent: {id: 4, align: "be", offset: 100}},
      ]
    },
    {
      name: "Multiple de-intersection (left-side)",
      screens: [
        {id: 1, w: 1920, h: 1080, s: 1},
        {id: 2, w: 1024, h: 768, s: 1, parent: {id: 1, align: "le", offset: 0}},
        {id: 3, w: 1024, h: 768, s: 1, parent: {id: 2, align: "b", offset: 300}},
        {id: 4, w: 1024, h: 768, s: 1.5, parent: {id: 2, align: "le", offset: 100}},
        {id: 5, w: 1024, h: 768, s: 1, parent: {id: 4, align: "be", offset: 100}},
      ]
    },
    {
      name: "Parent de-intersection child offset",
      screens: [
        {id: 1, w: 1600, h: 1600, s: 1.5},
        {id: 2, w: 800, h: 800, s: 1, parent: {id: 1, align: "r", offset: 0}},
        {id: 3, w: 800, h: 800, s: 1, parent: {id: 1, align: "r", offset: 800}},
        {id: 4, w: 800, h: 1600, s: 1, parent: {id: 2, align: "r", offset: 0}},
      ]
    },
  ],
].map(sections => sections.map(layout => {
  return parseLayout(layout)
}))

function parseLayout(layout) {
  const screens = []

  for (const screen of layout.screens) {
    let x = 0, y = 0
    const {w, h} = screen

    if (screen.parent) {
      const parent = screens.find(s => s.ID == screen.parent.id).Bounds
      const offset = screen.parent.offset
      let align = screen.parent.align
      let align2 = ""

      if (align.length == 2) {
        align2 = align.charAt(1)
        align = align.charAt(0)
      }

      x = parent.X
      y = parent.Y
			// t: top, b: bottom, l: left, r: right, e: edge, c: corner
      if (align == "t" || align == "b") {
        x += offset + (align2 == "e" || align2 == "c" ? parent.Width : 0) - (align2 == "e" ? w : 0)
        y += (align == "t" ? -h : parent.Height)
      } else {
        y += offset + (align2 == "e" || align2 == "c" ? parent.Height : 0) - (align2 == "e" ? h : 0)
        x += (align == "l" ? -w : parent.Width)
      }
    }

    screens.push({
      ID: `${screen.id}`,
      Name: screen.name ?? `Display${screen.id}`,
      ScaleFactor: Math.round(screen.s * 100) / 100,
      X: x,
      Y: y,
      Size: {Width: w, Height: h},
      Bounds: {X: x, Y: y, Width: w, Height: h},
      PhysicalBounds: {X: x, Y: y, Width: w, Height: h},
      WorkArea: {X: x, Y: y, Width: w, Height: h-Math.round(40*screen.s)},
      PhysicalWorkArea: {X: x, Y: y, Width: w, Height: h-Math.round(40*screen.s)},
      IsPrimary: screen.id == 1,
      Rotation: 0
    })
  }

  return {name: layout.name, screens}
}
