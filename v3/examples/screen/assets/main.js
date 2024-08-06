setExamplesType(document.getElementById('examples-type').value, 0)

function setExamplesType(type, autoSelectLayout = 1) {
  window.examples_type = parseInt(type)
  document.getElementById('examples-list').innerHTML = examples[examples_type].map((layout, i) => {
    return `<span class="radio-btn" data-value="${i + 1}" title="${layout.name}">${i + 1}</span>`
  }).join("\n")
  if (autoSelectLayout != null) setLayout(autoSelectLayout)
}

async function setLayout(indexOrLayout, physicalCoordinate = true) {
  if (typeof indexOrLayout == 'number') {
    await radioBtnClick(null, `#layout-selector [data-value="${indexOrLayout}"]`)
  } else {
    document.querySelectorAll('#layout-selector .active').forEach(el => el.classList.remove('active'))
    window.layout = indexOrLayout
    window.point = null
    window.rect = null
    await processLayout()
    await draw()
  }

  const physical = !parseInt(document.querySelector('#coordinate-selector .active').dataset.value)
  if (physical != physicalCoordinate) {
    await setCoordinateType(physicalCoordinate)
  }
}

async function setCoordinateType(physicalCoordinate = true) {
  await radioBtnClick(null, `#coordinate-selector [data-value="${physicalCoordinate ? 0 : 1}"]`)
}

async function radioBtnClick(e, selector) {
  if (e == null) {
    e = new Event("mousedown")
    document.querySelector(selector).dispatchEvent(e)
  }
  if (!e.target.classList.contains('radio-btn')) return
  const btnGroup = e.target.closest('.radio-btn-group')
  btnGroup.querySelectorAll('.radio-btn.active').forEach(el => el.classList.remove('active'))
  e.target.classList.add('active')

  if (btnGroup.id == 'layout-selector') {
    window.point = null
    window.rect = null
    await processLayout()
  }

  await draw()
}

async function processLayout() {
  const layoutBtn = document.querySelector('#layout-selector .active')
  const i = layoutBtn ? parseInt(layoutBtn.dataset.value) : -1
  if (i == 0) {
    // system screens
    window.layout = {
      name: '',
      screens: await callBinding('main.ScreenService.GetSystemScreens'),
    }
  } else {
    if (i > 0) {
      // example layouts
      window.layout = structuredClone(examples[examples_type][i - 1])
    }
    layout.screens = await callBinding('main.ScreenService.ProcessExampleScreens', layout.screens)
  }
  document.getElementById('example-name').textContent = layout.name
}

async function draw() {
  console.log(layout)
  let minX = 0, minY = 0, maxX = 0, maxY = 0;
  let html = '';

  const physical = !parseInt(document.querySelector('#coordinate-selector .active').dataset.value)
  const retainViewbox = document.querySelector('#retain-viewbox').checked

  layout.screens.forEach(screen => {
    const b = physical ? screen.PhysicalBounds : screen.Bounds
    const wa = physical ? screen.PhysicalWorkArea : screen.WorkArea
    const vbBounds = retainViewbox ? [screen.Bounds, screen.PhysicalBounds] : [b]

    minX = Math.min(minX, ...vbBounds.map(b => b.X))
    minY = Math.min(minY, ...vbBounds.map(b => b.Y))
    maxX = Math.max(maxX, ...vbBounds.map(b => b.X + b.Width))
    maxY = Math.max(maxY, ...vbBounds.map(b => b.Y + b.Height))

    html += `
      <rect x="${b.X}" y="${b.Y}" width="${b.Width}" height="${b.Height}" fill="#00ceff" />
      <rect x="${wa.X}" y="${wa.Y}" width="${wa.Width}" height="${wa.Height}" fill="#def9ff" />
      <rect x="${b.X + 1}" y="${b.Y + 1}" width="${b.Width - 2}" height="${b.Height - 2}" stroke="black" stroke-width="2" />

      <g transform="translate(${b.X}, ${b.Y})" fill="black">
        <text x="10" y="10" text-anchor="start" dominant-baseline="hanging" font-size="50">(${b.X}, ${b.Y})</text>
        <g transform="translate(${b.Width / 2}, ${b.Height / 2})" text-anchor="middle" font-size="100">
          <text x="0" y="-1.5em" fill="${screen.IsPrimary ? '#006bff' : '#5881ba'}" font-weight="bold">${screen.Name}</text>
          <text x="0" y="0">${b.Width} x ${b.Height}</text>
          <text x="0" y="1.2em">Scale factor: ${screen.ScaleFactor}</text>
        </g>
      </g>
    `
  })

  const svg = document.getElementById('svg')
  svg.innerHTML = `
    ${svg.querySelector('& > defs').outerHTML}
    <rect x="${minX}" y="${minY}" width="${maxX - minX}" height="${maxY - minY}" fill="antiquewhite" />
    ${html}
    <g id="rects"></g>
    <g id="points"></g>
  `

  svg.setAttribute('viewBox', `${minX} ${minY} ${maxX - minX} ${maxY - minY}`)

  if (window.point) await probePoint()
  if (window.rect) await drawRect()

  svg.onmousedown = async function(e) {
    let pt = new DOMPoint(e.clientX, e.clientY)
    pt = pt.matrixTransform(svg.getScreenCTM().inverse())
    pt.x = parseInt(pt.x)
    pt.y = parseInt(pt.y)
    if (e.buttons == 1) {
      await probePoint({X: pt.x, Y: pt.y})
    } else if (e.buttons == 2) {
      if (e.ctrlKey) {
        if (!window.rect) {
          window.rect = {X: pt.x, Y: pt.y, Width: 0, Height: 0}
        }
        if (!window.rectCursor) {
          window.rectAnchor = {x: window.rect.X, y: window.rect.Y}
          window.rectCursor = {x: window.rectAnchor.x + window.rect.Width, y: window.rectAnchor.y + window.rect.Height}
        }
        window.rectCursorOffset = {
          x: pt.x - window.rectCursor.x,
          y: pt.y - window.rectCursor.y,
        }
      } else {
        window.rectAnchor = pt
        window.rectCursorOffset = {x: 0, y: 0}
        window.probing = true
        drawRect({X: pt.x, Y: pt.y, Width: 0, Height: 0})
        window.probing = false
      }
    } else if (e.buttons == 4) {
      drawRect({X: pt.x, Y: pt.y, Width: 50, Height: 50})
    }
  }
  svg.onmousemove = async function(e) {
    if (window.probing) return
    window.probing = true
    if (e.buttons == 1) {
      await svg.onmousedown(e)
    } else if (e.buttons == 2) {
      let pt = new DOMPoint(e.clientX, e.clientY)
      pt = pt.matrixTransform(svg.getScreenCTM().inverse())
      if (e.ctrlKey) {
        window.rectAnchor.x += pt.x - rectCursor.x - window.rectCursorOffset.x
        window.rectAnchor.y += pt.y - rectCursor.y - window.rectCursorOffset.y
      }
      window.rectCursor = {
        x: pt.x - window.rectCursorOffset.x,
        y: pt.y - window.rectCursorOffset.y,
      }
      await drawRect({
        X: parseInt(Math.min(window.rectAnchor.x, window.rectCursor.x)),
        Y: parseInt(Math.min(window.rectAnchor.y, window.rectCursor.y)),
        Width: parseInt(Math.abs(window.rectCursor.x - window.rectAnchor.x)),
        Height: parseInt(Math.abs(window.rectCursor.y - window.rectAnchor.y)),
      })
    }
    window.probing = false
  }
  svg.oncontextmenu = function(e) {
    e.preventDefault()
  }
}

async function probePoint(p = null) {
  const svg = document.getElementById('svg');
  const physical = !parseInt(document.querySelector('#coordinate-selector .active').dataset.value)

  if (p == null) {
    if (window.pointIsPhysical == physical) {
      p = window.point
    } else {
      p = (await callBinding('main.ScreenService.TransformPoint', window.point, window.pointIsPhysical))[0]
    }
  }

  window.point = p
  window.pointIsPhysical = physical
  const [ptTransformed, ptDblTransformed] = await callBinding('main.ScreenService.TransformPoint', p, physical)

  svg.getElementById('points').innerHTML = `
    <circle cx="${p.X}" cy="${p.Y}" r="15" fill="red" />
    <circle cx="${ptTransformed.X}" cy="${ptTransformed.Y}" r="5" fill="green" />
    <circle cx="${ptTransformed.X}" cy="${ptTransformed.Y}" r="35" stroke="green" stroke-width="4" />
    <circle cx="${ptDblTransformed.X}" cy="${ptDblTransformed.Y}" r="25" stroke="red" stroke-width="4" />
  `
  // await new Promise((resolve) => setTimeout(resolve, 200)) // delay
  return ptDblTransformed
}

async function drawRect(r = null) {
  const svg = document.getElementById('svg');
  const physical = !parseInt(document.querySelector('#coordinate-selector .active').dataset.value)

  if (r == null) {
    if (window.rectIsPhysical == physical) {
      r = window.rect
    } else {
      r = await callBinding('main.ScreenService.TransformRect', window.rect, window.rectIsPhysical)
    }
  }

  if (!window.probing) {
    window.rectAnchor = null
    window.rectCursor = null
  }

  document.getElementById('x').value = r.X
  document.getElementById('y').value = r.Y
  document.getElementById('w').value = r.Width
  document.getElementById('h').value = r.Height

  window.rect = r
  window.rectIsPhysical = physical
  window.rTransformed = await callBinding('main.ScreenService.TransformRect', r, physical)
  window.rDblTransformed = await callBinding('main.ScreenService.TransformRect', rTransformed, !physical)
  window.rTransformed = rTransformed

  await rectLayers()
  return rDblTransformed
}

async function rectLayers() {
  const s = document.getElementById('slider').value
  if (window.rect == null) await test1()

  const r = await callBinding('main.ScreenService.TransformRect', rectIsPhysical ? rect : rTransformed, true)
  const rShifted = {...r, X: r.X+50}
  const rShiftedPhysical = await callBinding('main.ScreenService.TransformRect', rShifted, false)

  svg.getElementById('rects').innerHTML = [
    [window.rect, 'rgb(255 255 255 / 100%)'],     // w
    [window.rTransformed, 'rgb(0 255 0 / 25%)'],  // g
    [window.rDblTransformed, 'none'],             // none
    [rShifted, 'rgb(255 0 0 / 15%)'],             // r
    [rShiftedPhysical, 'rgb(0 0 255 / 15%)'],     // b
  ].filter((_,i) => i<s).map(([r, color], i) => {
    let lines = ''
    if (i == 0) {
      const center = {X: r.X + (r.Width-1)/2, Y: r.Y + (r.Height-1)/2}
      lines += `
        <line x1="${center.X}" x2="${center.X}" y1="${r.Y}" y2="${r.Y + r.Height-1}" stroke="gray" stroke-width="1" />
        <line x1="${r.X}" x2="${r.X + r.Width-1}" y1="${center.Y}" y2="${center.Y}" stroke="gray" stroke-width="1" />
      `
    }
    return `<rect x="${r.X}" y="${r.Y}" width="${r.Width}" height="${r.Height}" fill="${color}" stroke="${color == 'none' ? 'red' : 'black'}" stroke-width="${color == 'none' ? 5 : 1}" stroke-dasharray="${color == 'none' ? 5 : 'none'}" />${lines}`
  }).join('/n')
}

async function updateDipRect(x, y=0, w=0, h=0) {
  if (rect == null) {
    await drawRect({
      X: +document.getElementById('x').value,
      Y: +document.getElementById('y').value,
      Width: +document.getElementById('w').value,
      Height: +document.getElementById('h').value,
    })
  }
  // Simulate real window by first retrieving the physical bounds then transforming it to dip
  // then updating the bounds and transforming it back to physical
  let rPhysical = rectIsPhysical ? rect : rTransformed
  const r = await callBinding('main.ScreenService.TransformRect', rPhysical, true)
  r.X += x
  r.Y += y
  r.Width += w
  r.Height += h
  rPhysical = await callBinding('main.ScreenService.TransformRect', r, false)
  drawRect(rectIsPhysical ? rPhysical : r)
}

function arrowMove(e) {
  let x = 0, y = 0
  if (e.key == 'ArrowLeft') x = -step.value
  if (e.key == 'ArrowRight') x = +step.value
  if (e.key == 'ArrowUp') y = -step.value
  if (e.key == 'ArrowDown') y = +step.value
  if (!(x || y)) return
  e.preventDefault()
  updateDipRect(x, y)
}

async function test1() {
  // Edge case 1: invalid dip rect: no physical rect can produce it
  await setLayout(parseLayout({screens: [
    {id: 1, w: 1200, h: 1200, s: 1},
    {id: 2, w: 1200, h: 1100, s: 1.5, parent: {id: 1, align: "r", offset: 0}},
  ]}), false)
  await drawRect({X: 1050, Y: 700, Width: 400, Height: 300})
}

async function test2() {
  // Edge case 2: physical rect that changes when double transformed (2 physical rects produce the same dip rect)
  await setLayout(parseLayout({screens: [
    {id: 1, w: 1200, h: 1200, s: 1.5},
    {id: 2, w: 1200, h: 900, s: 1, parent: {id: 1, align: "r", offset: 0}},
  ]}), true)
  await drawRect({X: 1050, Y: 890, Width: 400, Height: 300})
}

async function probeLayout(finishup = true) {
  const probeButtons = document.getElementById('probe-buttons')
  const svg = document.getElementById('svg')
  const threshold = 1

  const physical = !parseInt(document.querySelector('#coordinate-selector .active').dataset.value)
  window.cancelProbing = false
  probeButtons.classList.add('active')

  const steps = 3
  let failed = false
  for (const screen of layout.screens) {
    if (window.cancelProbing) break
    const b = physical ? screen.PhysicalBounds : screen.Bounds
    const xStep = parseInt(b.Width / steps) || 1
    const yStep = parseInt(b.Height / steps) || 1
    let x = b.X, y = b.Y
    let xDone = false, yDone = false

    while (!(yDone || window.cancelProbing)) {
      if (y >= b.Y + b.Height - 1) {
        y = b.Y + b.Height - 1
        yDone = true
      }
      x = b.X
      xDone = false
      while (!(xDone || window.cancelProbing)) {
        if (x >= b.X + b.Width - 1) {
          x = b.X + b.Width - 1
          xDone = true
        }
        const pt = {X: x, Y: y}
        let ptDblTransformed, err
        try {
          ptDblTransformed = await probePoint(pt)
        } catch (e) {
          err = e
        }
        if (err || Math.abs(pt.X - ptDblTransformed.X) > threshold || Math.abs(pt.Y - ptDblTransformed.Y) > threshold) {
          failed = true
          console.log(pt, ptDblTransformed)
          window.cancelProbing = true
          setTimeout(() => {
            alert(err ?? `**FAILED**\nProbing failed at point: {X: ${pt.X}, Y: ${pt.Y}}\nDouble transformed point: {X: ${ptDblTransformed.X}, Y: ${ptDblTransformed.Y}}\n(Exceeded threshold of ${threshold} pixels)`)
          }, 50)
        }
        x += xStep
      }
      y += yStep
    }
  }

  if (finishup || window.cancelProbing) probeButtons.classList.remove('active')
  if (!(failed || window.cancelProbing)) {
    window.point = null
    if (finishup) {
      setTimeout(() => {
        svg.getElementById('points').innerHTML = ''
        alert(`Successfully probed all points!, All within threshold of ${threshold} pixels.`)
      }, 50)
    }
    return true
  }
}

async function probeAllExamples() {
  console.time('probeAllExamples')
loop1:
  for (let typeI = 0; typeI < examples.length; typeI++) {
    document.getElementById('examples-type').value = typeI
    setExamplesType(typeI, null)

    for (let layoutI = (typeI ? 0 : -1); layoutI < examples[typeI].length; layoutI++) {
      await radioBtnClick(null, `#layout-selector [data-value="${layoutI + 1}"]`)
      for (let i = 0; i < 2; i++) {
        const lastLayout = (typeI == examples.length - 1 && layoutI == examples[typeI].length - 1 && i == 1)
        if (!await probeLayout(lastLayout)) break loop1
        if (i == 0) await setCoordinateType(!pointIsPhysical)
      }
    }
  }
  console.timeEnd('probeAllExamples')
}

async function callBinding(name, ...params) {
  return wails.Call.ByName(name, ...params)
}

function showAdvanced(e) {
  e.target.style.display = 'none'
  document.querySelectorAll('.advanced').forEach(el => el.style.display = 'initial')
}
