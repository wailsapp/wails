//go:build mcp && !ios && !android

package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"time"
)

// mcpTool is a single MCP tool: metadata plus its handler. Handlers return a
// value that is JSON-encoded into the tool result, or an error which becomes
// an isError result.
type mcpTool struct {
	Name        string
	Description string
	Schema      map[string]any
	Handler     func(args map[string]any) (any, error)
}

// Argument helpers. MCP arguments arrive as generic JSON; numbers are float64.

func mcpArgString(args map[string]any, key string) (string, bool) {
	value, ok := args[key].(string)
	return value, ok && value != ""
}

func mcpArgFloat(args map[string]any, key string) (float64, bool) {
	switch value := args[key].(type) {
	case float64:
		return value, true
	case json.Number:
		f, err := value.Float64()
		return f, err == nil
	default:
		return 0, false
	}
}

func mcpArgInt(args map[string]any, key string, fallback int) int {
	if value, ok := mcpArgFloat(args, key); ok {
		return int(value)
	}
	return fallback
}

func mcpArgBool(args map[string]any, key string) bool {
	value, _ := args[key].(bool)
	return value
}

func mcpArgStrings(args map[string]any, key string) []string {
	raw, ok := args[key].([]any)
	if !ok {
		return nil
	}
	result := make([]string, 0, len(raw))
	for _, item := range raw {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

// Schema helpers.

func mcpObjectSchema(required []string, properties map[string]any) map[string]any {
	s := map[string]any{
		"type":       "object",
		"properties": properties,
	}
	if len(required) > 0 {
		s["required"] = required
	}
	return s
}

func mcpProp(propType, description string) map[string]any {
	return map[string]any{"type": propType, "description": description}
}

func mcpWindowProp() map[string]any {
	return mcpProp("string", "Window name. Defaults to the focused window, or the first window.")
}

func mcpTimeoutProp() map[string]any {
	return mcpProp("number", "Timeout in milliseconds for the page to respond. Defaults to WAILS_MCP_TIMEOUT (30s).")
}

func mcpTargetProps(prefix, what string) map[string]any {
	key := func(name string) string {
		if prefix == "" {
			return name
		}
		return prefix + "_" + name
	}
	return map[string]any{
		key("x"):        mcpProp("number", "X coordinate of "+what+" in CSS pixels, relative to the window viewport. Ignored when a selector is given."),
		key("y"):        mcpProp("number", "Y coordinate of "+what+" in CSS pixels, relative to the window viewport. Ignored when a selector is given."),
		key("selector"): mcpProp("string", "CSS selector of "+what+". The first matching element is scrolled into view and its centre is used."),
	}
}

// mcpTarget builds the JSON target object consumed by the in-page library.
func mcpTarget(args map[string]any, prefix string) (string, error) {
	key := func(name string) string {
		if prefix == "" {
			return name
		}
		return prefix + "_" + name
	}
	result := map[string]any{}
	if selector, ok := mcpArgString(args, key("selector")); ok {
		result["selector"] = selector
	}
	x, hasX := mcpArgFloat(args, key("x"))
	y, hasY := mcpArgFloat(args, key("y"))
	if hasX && hasY {
		result["x"], result["y"] = x, y
	}
	if len(result) == 0 {
		return "", fmt.Errorf("provide either %q or both %q and %q", key("selector"), key("x"), key("y"))
	}
	data, err := json.Marshal(result)
	return string(data), err
}

func mcpMergeProps(maps ...map[string]any) map[string]any {
	result := map[string]any{}
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// mcpWindowDescriptor summarises a window for tool results.
func mcpWindowDescriptor(window Window) map[string]any {
	width, height := window.Size()
	x, y := window.Position()
	return map[string]any{
		"id":         window.ID(),
		"name":       window.Name(),
		"width":      width,
		"height":     height,
		"x":          x,
		"y":          y,
		"focused":    window.IsFocused(),
		"visible":    window.IsVisible(),
		"fullscreen": window.IsFullscreen(),
		"maximised":  window.IsMaximised(),
		"minimised":  window.IsMinimised(),
		"zoom":       window.GetZoom(),
	}
}

// mcpEvalTool is the common shape of tools that run JavaScript in a window
// and return its JSON result.
func (m *mcpServer) mcpEvalTool(args map[string]any, body string) (any, error) {
	name, _ := mcpArgString(args, "window")
	window, err := m.resolveWindow(name)
	if err != nil {
		return nil, err
	}
	var result any
	if err := m.evalInto(window, body, m.mcpEvalTimeout(args), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (m *mcpServer) registerTools() {
	m.tools = []*mcpTool{
		{
			Name:        "app_info",
			Description: "Get information about the running Wails application: platform, windows and MCP endpoint.",
			Schema:      mcpObjectSchema(nil, map[string]any{}),
			Handler: func(args map[string]any) (any, error) {
				windows := m.app.Window.GetAll()
				descriptors := make([]map[string]any, 0, len(windows))
				for _, window := range windows {
					descriptors = append(descriptors, mcpWindowDescriptor(window))
				}
				return map[string]any{
					"os":       runtime.GOOS,
					"arch":     runtime.GOARCH,
					"endpoint": fmt.Sprintf("http://%s/mcp", m.addr),
					"windows":  descriptors,
				}, nil
			},
		},
		{
			Name:        "windows_list",
			Description: "List all application windows with their geometry and state.",
			Schema:      mcpObjectSchema(nil, map[string]any{}),
			Handler: func(args map[string]any) (any, error) {
				windows := m.app.Window.GetAll()
				descriptors := make([]map[string]any, 0, len(windows))
				for _, window := range windows {
					descriptors = append(descriptors, mcpWindowDescriptor(window))
				}
				return descriptors, nil
			},
		},
		{
			Name: "window_control",
			Description: "Control a window. Actions: focus, show, hide, close, center, maximise, unmaximise, " +
				"minimise, unminimise, restore, fullscreen, unfullscreen, toggle_fullscreen, reload, force_reload, " +
				"open_devtools, set_size, set_position, set_title, set_url, set_zoom, set_always_on_top.",
			Schema: mcpObjectSchema([]string{"action"}, map[string]any{
				"action":        mcpProp("string", "The action to perform."),
				"window":        mcpWindowProp(),
				"width":         mcpProp("number", "New width for set_size."),
				"height":        mcpProp("number", "New height for set_size."),
				"x":             mcpProp("number", "New x position for set_position."),
				"y":             mcpProp("number", "New y position for set_position."),
				"title":         mcpProp("string", "New title for set_title."),
				"url":           mcpProp("string", "URL to load for set_url."),
				"zoom":          mcpProp("number", "Zoom factor for set_zoom (1.0 = 100%)."),
				"always_on_top": mcpProp("boolean", "Flag for set_always_on_top."),
			}),
			Handler: m.mcpWindowControl,
		},
		{
			Name: "js_eval",
			Description: "Evaluate JavaScript in a window. The code runs in an async function body, so use " +
				"`return` to produce a value and `await` freely. The result must be JSON-serialisable. " +
				"The Wails runtime is importable inside the page via `await import('/wails/runtime.js')`.",
			Schema: mcpObjectSchema([]string{"js"}, map[string]any{
				"js":         mcpProp("string", "JavaScript code to run (async function body)."),
				"window":     mcpWindowProp(),
				"timeout_ms": mcpTimeoutProp(),
			}),
			Handler: func(args map[string]any) (any, error) {
				js, ok := mcpArgString(args, "js")
				if !ok {
					return nil, errors.New("missing required argument: js")
				}
				return m.mcpEvalTool(args, js)
			},
		},
		{
			Name: "dom_html",
			Description: "Get the HTML of the page or of the first element matching a selector. " +
				"Useful for inspecting the UI before interacting with it.",
			Schema: mcpObjectSchema(nil, map[string]any{
				"selector":  mcpProp("string", "CSS selector. Defaults to the whole document."),
				"max_bytes": mcpProp("number", "Maximum HTML length to return. Defaults to 100000."),
				"window":    mcpWindowProp(),
			}),
			Handler: func(args map[string]any) (any, error) {
				selector, _ := mcpArgString(args, "selector")
				maxBytes := mcpArgInt(args, "max_bytes", 100_000)
				return m.mcpEvalTool(args, fmt.Sprintf(`
					const selector = %s, maxBytes = %d;
					const el = selector ? document.querySelector(selector) : document.documentElement;
					if (!el) throw new Error('no element matches selector: ' + selector);
					const html = el.outerHTML;
					return {
						html: html.length > maxBytes ? html.slice(0, maxBytes) : html,
						truncated: html.length > maxBytes,
						totalLength: html.length,
					};`,
					strconv.Quote(selector), maxBytes))
			},
		},
		{
			Name: "dom_query",
			Description: "Find elements by CSS selector and return a summary of each: tag, id, classes, text, " +
				"value, viewport bounds and visibility. Use this to discover what to click or type into.",
			Schema: mcpObjectSchema([]string{"selector"}, map[string]any{
				"selector": mcpProp("string", "CSS selector to query."),
				"limit":    mcpProp("number", "Maximum number of elements to return. Defaults to 25."),
				"window":   mcpWindowProp(),
			}),
			Handler: func(args map[string]any) (any, error) {
				selector, ok := mcpArgString(args, "selector")
				if !ok {
					return nil, errors.New("missing required argument: selector")
				}
				return m.mcpEvalTool(args, fmt.Sprintf("return mcp.query(%s, %d);",
					strconv.Quote(selector), mcpArgInt(args, "limit", 25)))
			},
		},
		{
			Name: "mouse_move",
			Description: "Move the animated mouse cursor to a point or element, firing pointer/mouse move and " +
				"hover events along the way. Returns a description of the hovered element.",
			Schema: mcpObjectSchema(nil, mcpMergeProps(mcpTargetProps("", "the destination"), map[string]any{
				"duration_ms": mcpProp("number", "Animation duration in ms. Defaults to a natural speed based on distance."),
				"window":      mcpWindowProp(),
				"timeout_ms":  mcpTimeoutProp(),
			})),
			Handler: func(args map[string]any) (any, error) {
				targetJSON, err := mcpTarget(args, "")
				if err != nil {
					return nil, err
				}
				return m.mcpEvalTool(args, fmt.Sprintf("return await mcp.move(%s, {duration: %d});",
					targetJSON, mcpArgInt(args, "duration_ms", 0)))
			},
		},
		{
			Name: "mouse_click",
			Description: "Click a point or element with the animated mouse cursor: the cursor visibly moves " +
				"there, presses with a ripple effect and dispatches the full pointer/mouse event sequence. " +
				"Returns a description of the clicked element.",
			Schema: mcpObjectSchema(nil, mcpMergeProps(mcpTargetProps("", "the click target"), map[string]any{
				"button":     mcpProp("string", "Mouse button: left (default), right or middle."),
				"count":      mcpProp("number", "Number of clicks: 1 (default) or 2 for a double-click."),
				"modifiers":  map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Held modifier keys: ctrl, shift, alt, meta."},
				"window":     mcpWindowProp(),
				"timeout_ms": mcpTimeoutProp(),
			})),
			Handler: func(args map[string]any) (any, error) {
				targetJSON, err := mcpTarget(args, "")
				if err != nil {
					return nil, err
				}
				button, _ := mcpArgString(args, "button")
				options := map[string]any{
					"button":    button,
					"count":     mcpArgInt(args, "count", 1),
					"modifiers": mcpArgStrings(args, "modifiers"),
				}
				optionsJSON, err := json.Marshal(options)
				if err != nil {
					return nil, err
				}
				return m.mcpEvalTool(args, fmt.Sprintf("return await mcp.click(%s, %s);", targetJSON, optionsJSON))
			},
		},
		{
			Name: "mouse_drag",
			Description: "Drag from one point/element to another with the animated cursor: press, animated move " +
				"with continuous pointer/mouse events (and HTML5 drag events for draggable elements), release.",
			Schema: mcpObjectSchema(nil, mcpMergeProps(
				mcpTargetProps("from", "the drag start"),
				mcpTargetProps("to", "the drop target"),
				map[string]any{
					"duration_ms": mcpProp("number", "Drag animation duration in ms. Defaults to a natural speed."),
					"window":      mcpWindowProp(),
					"timeout_ms":  mcpTimeoutProp(),
				})),
			Handler: func(args map[string]any) (any, error) {
				fromJSON, err := mcpTarget(args, "from")
				if err != nil {
					return nil, err
				}
				toJSON, err := mcpTarget(args, "to")
				if err != nil {
					return nil, err
				}
				return m.mcpEvalTool(args, fmt.Sprintf("return await mcp.drag(%s, %s, {duration: %d});",
					fromJSON, toJSON, mcpArgInt(args, "duration_ms", 0)))
			},
		},
		{
			Name: "mouse_scroll",
			Description: "Scroll at a point or element: the animated cursor moves there, a wheel event is " +
				"dispatched and the nearest scrollable container scrolls smoothly.",
			Schema: mcpObjectSchema(nil, mcpMergeProps(mcpTargetProps("", "the scroll position"), map[string]any{
				"delta_x":    mcpProp("number", "Horizontal scroll amount in pixels."),
				"delta_y":    mcpProp("number", "Vertical scroll amount in pixels. Positive scrolls down."),
				"window":     mcpWindowProp(),
				"timeout_ms": mcpTimeoutProp(),
			})),
			Handler: func(args map[string]any) (any, error) {
				targetJSON, err := mcpTarget(args, "")
				if err != nil {
					return nil, err
				}
				deltaX, _ := mcpArgFloat(args, "delta_x")
				deltaY, _ := mcpArgFloat(args, "delta_y")
				if deltaX == 0 && deltaY == 0 {
					deltaY = 120
				}
				return m.mcpEvalTool(args, fmt.Sprintf("return await mcp.scroll(%s, %g, %g);",
					targetJSON, deltaX, deltaY))
			},
		},
		{
			Name: "keyboard_type",
			Description: "Type text into the focused element (or an element given by selector, which is clicked " +
				"first) with per-character key, input and change events. Works with inputs, textareas and " +
				"contenteditable elements, including React-style controlled inputs.",
			Schema: mcpObjectSchema([]string{"text"}, map[string]any{
				"text":       mcpProp("string", "The text to type."),
				"selector":   mcpProp("string", "CSS selector of the element to focus first."),
				"delay_ms":   mcpProp("number", "Delay between characters in ms. Defaults to 25."),
				"window":     mcpWindowProp(),
				"timeout_ms": mcpTimeoutProp(),
			}),
			Handler: func(args map[string]any) (any, error) {
				text, ok := args["text"].(string)
				if !ok {
					return nil, errors.New("missing required argument: text")
				}
				textJSON, err := json.Marshal(text)
				if err != nil {
					return nil, err
				}
				selector, _ := mcpArgString(args, "selector")
				return m.mcpEvalTool(args, fmt.Sprintf("return await mcp.typeText(%s, %s, %d);",
					textJSON, strconv.Quote(selector), mcpArgInt(args, "delay_ms", 25)))
			},
		},
		{
			Name: "keyboard_press",
			Description: "Press a single key with optional modifiers, e.g. Enter, Tab, Escape, Backspace, " +
				"ArrowDown, F5, a. Dispatches keydown/keyup and applies the default editing action for " +
				"Enter, Backspace and Tab.",
			Schema: mcpObjectSchema([]string{"key"}, map[string]any{
				"key":        mcpProp("string", "Key value as in KeyboardEvent.key."),
				"modifiers":  map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Held modifier keys: ctrl, shift, alt, meta."},
				"window":     mcpWindowProp(),
				"timeout_ms": mcpTimeoutProp(),
			}),
			Handler: func(args map[string]any) (any, error) {
				key, ok := mcpArgString(args, "key")
				if !ok {
					return nil, errors.New("missing required argument: key")
				}
				modifiersJSON, err := json.Marshal(mcpArgStrings(args, "modifiers"))
				if err != nil {
					return nil, err
				}
				return m.mcpEvalTool(args, fmt.Sprintf("return await mcp.press(%s, %s);",
					strconv.Quote(key), modifiersJSON))
			},
		},
		{
			Name: "call_bound_method",
			Description: "Call a bound Go service method through the Wails runtime, e.g. " +
				"'main.GreetService.Greet' with args ['World']. Returns the method's result.",
			Schema: mcpObjectSchema([]string{"name"}, map[string]any{
				"name":       mcpProp("string", "Fully qualified method name: package.Service.Method."),
				"args":       map[string]any{"type": "array", "description": "Positional arguments for the method."},
				"window":     mcpWindowProp(),
				"timeout_ms": mcpTimeoutProp(),
			}),
			Handler: func(args map[string]any) (any, error) {
				name, ok := mcpArgString(args, "name")
				if !ok {
					return nil, errors.New("missing required argument: name")
				}
				callArgs, _ := args["args"].([]any)
				if callArgs == nil {
					callArgs = []any{}
				}
				argsJSON, err := json.Marshal(callArgs)
				if err != nil {
					return nil, err
				}
				return m.mcpEvalTool(args, fmt.Sprintf(`
					const runtime = await import('/wails/runtime.js');
					const result = await runtime.Call.ByName(%s, ...%s);
					return result === undefined ? null : result;`,
					strconv.Quote(name), argsJSON))
			},
		},
		{
			Name:        "emit_event",
			Description: "Emit a Wails application event that both Go and frontend listeners receive.",
			Schema: mcpObjectSchema([]string{"name"}, map[string]any{
				"name": mcpProp("string", "Event name."),
				"data": map[string]any{"description": "Optional event data (any JSON value)."},
			}),
			Handler: func(args map[string]any) (any, error) {
				name, ok := mcpArgString(args, "name")
				if !ok {
					return nil, errors.New("missing required argument: name")
				}
				if data, hasData := args["data"]; hasData {
					m.app.Event.Emit(name, data)
				} else {
					m.app.Event.Emit(name)
				}
				return "emitted " + name, nil
			},
		},
		{
			Name: "wait_for_event",
			Description: "Wait for a Wails application event to be emitted (by Go or the frontend) and return " +
				"its data. Useful for asserting that an interaction triggered the expected event.",
			Schema: mcpObjectSchema([]string{"name"}, map[string]any{
				"name":       mcpProp("string", "Event name to wait for."),
				"timeout_ms": mcpProp("number", "How long to wait in ms. Defaults to 30000."),
			}),
			Handler: func(args map[string]any) (any, error) {
				name, ok := mcpArgString(args, "name")
				if !ok {
					return nil, errors.New("missing required argument: name")
				}
				received := make(chan any, 1)
				cancel := m.app.Event.On(name, func(event *CustomEvent) {
					select {
					case received <- event.Data:
					default:
					}
				})
				defer cancel()
				select {
				case data := <-received:
					return map[string]any{"name": name, "data": data}, nil
				case <-time.After(m.mcpEvalTimeout(args)):
					return nil, fmt.Errorf("timed out waiting for event %q", name)
				}
			},
		},
		{
			Name: "screenshot_dom",
			Description: "Capture a structural snapshot of the visible page: the viewport size, scroll position, " +
				"focused element and a compact outline of visible elements with their bounds. This is a " +
				"DOM-based alternative to a pixel screenshot.",
			Schema: mcpObjectSchema(nil, map[string]any{
				"max_depth": mcpProp("number", "Maximum DOM depth to include. Defaults to 12."),
				"window":    mcpWindowProp(),
			}),
			Handler: func(args map[string]any) (any, error) {
				return m.mcpEvalTool(args, fmt.Sprintf("return mcp.snapshot(%d);", mcpArgInt(args, "max_depth", 12)))
			},
		},
	}
}

// mcpWindowControl implements the window_control tool.
func (m *mcpServer) mcpWindowControl(args map[string]any) (any, error) {
	action, ok := mcpArgString(args, "action")
	if !ok {
		return nil, errors.New("missing required argument: action")
	}
	name, _ := mcpArgString(args, "window")
	window, err := m.resolveWindow(name)
	if err != nil {
		return nil, err
	}

	requireNumbers := func(keys ...string) ([]int, error) {
		values := make([]int, len(keys))
		for i, key := range keys {
			value, ok := mcpArgFloat(args, key)
			if !ok {
				return nil, fmt.Errorf("action %s requires argument %q", action, key)
			}
			values[i] = int(value)
		}
		return values, nil
	}

	switch action {
	case "focus":
		window.Focus()
	case "show":
		window.Show()
	case "hide":
		window.Hide()
	case "close":
		window.Close()
	case "center":
		window.Center()
	case "maximise", "maximize":
		window.Maximise()
	case "unmaximise", "unmaximize":
		window.UnMaximise()
	case "minimise", "minimize":
		window.Minimise()
	case "unminimise", "unminimize":
		window.UnMinimise()
	case "restore":
		window.Restore()
	case "fullscreen":
		window.Fullscreen()
	case "unfullscreen":
		window.UnFullscreen()
	case "toggle_fullscreen":
		window.ToggleFullscreen()
	case "reload":
		window.Reload()
	case "force_reload":
		window.ForceReload()
	case "open_devtools":
		window.OpenDevTools()
	case "set_size":
		values, err := requireNumbers("width", "height")
		if err != nil {
			return nil, err
		}
		if values[0] <= 0 || values[1] <= 0 {
			return nil, errors.New("action set_size requires positive width and height")
		}
		window.SetSize(values[0], values[1])
	case "set_position":
		values, err := requireNumbers("x", "y")
		if err != nil {
			return nil, err
		}
		window.SetPosition(values[0], values[1])
	case "set_title":
		title, ok := mcpArgString(args, "title")
		if !ok {
			return nil, errors.New("action set_title requires argument \"title\"")
		}
		window.SetTitle(title)
	case "set_url":
		u, ok := mcpArgString(args, "url")
		if !ok {
			return nil, errors.New("action set_url requires argument \"url\"")
		}
		window.SetURL(u)
	case "set_zoom":
		zoom, ok := mcpArgFloat(args, "zoom")
		if !ok {
			return nil, errors.New("action set_zoom requires argument \"zoom\"")
		}
		window.SetZoom(zoom)
	case "set_always_on_top":
		val, ok := args["always_on_top"].(bool)
		if !ok {
			return nil, errors.New("action set_always_on_top requires a boolean argument \"always_on_top\"")
		}
		window.SetAlwaysOnTop(val)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}

	return mcpWindowDescriptor(window), nil
}
