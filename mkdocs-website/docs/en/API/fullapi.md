# application

```go
import "github.com/wailsapp/wails/v3/pkg/application"
```

## Index

- [Constants](#constants)
- [Variables](#variables)
- [func DefaultLogger\(level slog.Level\) \*slog.Logger](#DefaultLogger)
- [func Fatal\(message string, args ...interface\{\}\)](#Fatal)
- [func InvokeAsync\(fn func\(\)\)](#InvokeAsync)
- [func InvokeSync\(fn func\(\)\)](#InvokeSync)
- [func InvokeSyncWithError\(fn func\(\) error\) \(err error\)](#InvokeSyncWithError)
- [func InvokeSyncWithResult\[T any\]\(fn func\(\) T\) \(res T\)](#InvokeSyncWithResult)
- [func InvokeSyncWithResultAndError\[T any\]\(fn func\(\) \(T, error\)\) \(res T, err error\)](#InvokeSyncWithResultAndError)
- [func NewIconFromResource\(instance w32.HINSTANCE, resId uint16\) \(w32.HICON, error\)](#NewIconFromResource)
- [func ScaleToDefaultDPI\(pixels int, dpi uint\) int](#ScaleToDefaultDPI)
- [func ScaleWithDPI\(pixels int, dpi uint\) int](#ScaleWithDPI)
- [type ActivationPolicy](#ActivationPolicy)
- [type App](#App)
  - [func Get\(\) \*App](#Get)
  - [func New\(appOptions Options\) \*App](#New)
  - [func \(a \*App\) Capabilities\(\) capabilities.Capabilities](#App.Capabilities)
  - [func \(a \*App\) Clipboard\(\) \*Clipboard](#App.Clipboard)
  - [func \(a \*App\) CurrentWindow\(\) \*WebviewWindow](#App.CurrentWindow)
  - [func \(a \*App\) GetPID\(\) int](#App.GetPID)
  - [func \(a \*App\) GetPrimaryScreen\(\) \(\*Screen, error\)](#App.GetPrimaryScreen)
  - [func \(a \*App\) GetScreens\(\) \(\[\]\*Screen, error\)](#App.GetScreens)
  - [func \(a \*App\) GetWindowByName\(name string\) \*WebviewWindow](#App.GetWindowByName)
  - [func \(a \*App\) Hide\(\)](#App.Hide)
  - [func \(a \*App\) IsDarkMode\(\) bool](#App.IsDarkMode)
  - [func \(a \*App\) NewMenu\(\) \*Menu](#App.NewMenu)
  - [func \(a \*App\) NewSystemTray\(\) \*SystemTray](#App.NewSystemTray)
  - [func \(a \*App\) NewWebviewWindow\(\) \*WebviewWindow](#App.NewWebviewWindow)
  - [func \(a \*App\) NewWebviewWindowWithOptions\(windowOptions WebviewWindowOptions\) \*WebviewWindow](#App.NewWebviewWindowWithOptions)
  - [func \(a \*App\) On\(eventType events.ApplicationEventType, callback func\(event \*Event\)\) func\(\)](#App.On)
  - [func \(a \*App\) OnWindowCreation\(callback func\(window \*WebviewWindow\)\)](#App.OnWindowCreation)
  - [func \(a \*App\) Quit\(\)](#App.Quit)
  - [func \(a \*App\) RegisterContextMenu\(name string, menu \*Menu\)](#App.RegisterContextMenu)
  - [func \(a \*App\) RegisterHook\(eventType events.ApplicationEventType, callback func\(event \*Event\)\) func\(\)](#App.RegisterHook)
  - [func \(a \*App\) Run\(\) error](#App.Run)
  - [func \(a \*App\) SetMenu\(menu \*Menu\)](#App.SetMenu)
  - [func \(a \*App\) Show\(\)](#App.Show)
  - [func \(a \*App\) ShowAboutDialog\(\)](#App.ShowAboutDialog)
- [type ApplicationEventContext](#ApplicationEventContext)
  - [func \(c ApplicationEventContext\) IsDarkMode\(\) bool](#ApplicationEventContext.IsDarkMode)
  - [func \(c ApplicationEventContext\) OpenedFiles\(\) \[\]string](#ApplicationEventContext.OpenedFiles)
- [type Args](#Args)
  - [func \(a \*Args\) Bool\(s string\) \*bool](#Args.Bool)
  - [func \(a \*Args\) Float64\(s string\) \*float64](#Args.Float64)
  - [func \(a \*Args\) Int\(s string\) \*int](#Args.Int)
  - [func \(a \*Args\) String\(key string\) \*string](#Args.String)
  - [func \(a \*Args\) UInt\(s string\) \*uint](#Args.UInt)
  - [func \(a \*Args\) UInt8\(s string\) \*uint8](#Args.UInt8)
- [type AssetOptions](#AssetOptions)
- [type BackdropType](#BackdropType)
- [type BackgroundType](#BackgroundType)
- [type Bindings](#Bindings)
  - [func NewBindings\(structs \[\]any, aliases map\[uint32\]uint32\) \(\*Bindings, error\)](#NewBindings)
  - [func \(b \*Bindings\) Add\(structPtr interface\{\}\) error](#Bindings.Add)
  - [func \(b \*Bindings\) AddPlugins\(plugins map\[string\]Plugin\) error](#Bindings.AddPlugins)
  - [func \(b \*Bindings\) GenerateID\(name string\) \(uint32, error\)](#Bindings.GenerateID)
  - [func \(b \*Bindings\) Get\(options \*CallOptions\) \*BoundMethod](#Bindings.Get)
  - [func \(b \*Bindings\) GetByID\(id uint32\) \*BoundMethod](#Bindings.GetByID)
- [type BoundMethod](#BoundMethod)
  - [func \(b \*BoundMethod\) Call\(args \[\]interface\{\}\) \(returnValue interface\{\}, err error\)](#BoundMethod.Call)
  - [func \(b \*BoundMethod\) String\(\) string](#BoundMethod.String)
- [type Button](#Button)
  - [func \(b \*Button\) OnClick\(callback func\(\)\) \*Button](#Button.OnClick)
  - [func \(b \*Button\) SetAsCancel\(\) \*Button](#Button.SetAsCancel)
  - [func \(b \*Button\) SetAsDefault\(\) \*Button](#Button.SetAsDefault)
- [type CallOptions](#CallOptions)
  - [func \(c CallOptions\) Name\(\) string](#CallOptions.Name)
- [type Clipboard](#Clipboard)
  - [func \(c \*Clipboard\) SetText\(text string\) bool](#Clipboard.SetText)
  - [func \(c \*Clipboard\) Text\(\) \(string, bool\)](#Clipboard.Text)
- [type Context](#Context)
  - [func \(c \*Context\) ClickedMenuItem\(\) \*MenuItem](#Context.ClickedMenuItem)
  - [func \(c \*Context\) ContextMenuData\(\) any](#Context.ContextMenuData)
  - [func \(c \*Context\) IsChecked\(\) bool](#Context.IsChecked)
- [type ContextMenuData](#ContextMenuData)
- [type DialogType](#DialogType)
- [type Event](#Event)
  - [func \(w \*Event\) Cancel\(\)](#Event.Cancel)
  - [func \(w \*Event\) Context\(\) \*ApplicationEventContext](#Event.Context)
- [type EventListener](#EventListener)
- [type EventProcessor](#EventProcessor)
  - [func NewWailsEventProcessor\(dispatchEventToWindows func\(\*WailsEvent\)\) \*EventProcessor](#NewWailsEventProcessor)
  - [func \(e \*EventProcessor\) Emit\(thisEvent \*WailsEvent\)](#EventProcessor.Emit)
  - [func \(e \*EventProcessor\) Off\(eventName string\)](#EventProcessor.Off)
  - [func \(e \*EventProcessor\) OffAll\(\)](#EventProcessor.OffAll)
  - [func \(e \*EventProcessor\) On\(eventName string, callback func\(event \*WailsEvent\)\) func\(\)](#EventProcessor.On)
  - [func \(e \*EventProcessor\) OnMultiple\(eventName string, callback func\(event \*WailsEvent\), counter int\) func\(\)](#EventProcessor.OnMultiple)
  - [func \(e \*EventProcessor\) Once\(eventName string, callback func\(event \*WailsEvent\)\) func\(\)](#EventProcessor.Once)
  - [func \(e \*EventProcessor\) RegisterHook\(eventName string, callback func\(\*WailsEvent\)\) func\(\)](#EventProcessor.RegisterHook)
- [type FileFilter](#FileFilter)
- [type IconPosition](#IconPosition)
- [type MacAppearanceType](#MacAppearanceType)
- [type MacBackdrop](#MacBackdrop)
- [type MacOptions](#MacOptions)
- [type MacTitleBar](#MacTitleBar)
- [type MacToolbarStyle](#MacToolbarStyle)
- [type MacWindow](#MacWindow)
- [type Menu](#Menu)
  - [func NewMenu\(\) \*Menu](#NewMenu)
  - [func \(m \*Menu\) Add\(label string\) \*MenuItem](#Menu.Add)
  - [func \(m \*Menu\) AddCheckbox\(label string, enabled bool\) \*MenuItem](#Menu.AddCheckbox)
  - [func \(m \*Menu\) AddRadio\(label string, enabled bool\) \*MenuItem](#Menu.AddRadio)
  - [func \(m \*Menu\) AddRole\(role Role\) \*Menu](#Menu.AddRole)
  - [func \(m \*Menu\) AddSeparator\(\)](#Menu.AddSeparator)
  - [func \(m \*Menu\) AddSubmenu\(s string\) \*Menu](#Menu.AddSubmenu)
  - [func \(m \*Menu\) SetLabel\(label string\)](#Menu.SetLabel)
  - [func \(m \*Menu\) Update\(\)](#Menu.Update)
- [type MenuItem](#MenuItem)
  - [func \(m \*MenuItem\) Checked\(\) bool](#MenuItem.Checked)
  - [func \(m \*MenuItem\) Enabled\(\) bool](#MenuItem.Enabled)
  - [func \(m \*MenuItem\) Hidden\(\) bool](#MenuItem.Hidden)
  - [func \(m \*MenuItem\) IsCheckbox\(\) bool](#MenuItem.IsCheckbox)
  - [func \(m \*MenuItem\) IsRadio\(\) bool](#MenuItem.IsRadio)
  - [func \(m \*MenuItem\) IsSeparator\(\) bool](#MenuItem.IsSeparator)
  - [func \(m \*MenuItem\) IsSubmenu\(\) bool](#MenuItem.IsSubmenu)
  - [func \(m \*MenuItem\) Label\(\) string](#MenuItem.Label)
  - [func \(m \*MenuItem\) OnClick\(f func\(\*Context\)\) \*MenuItem](#MenuItem.OnClick)
  - [func \(m \*MenuItem\) SetAccelerator\(shortcut string\) \*MenuItem](#MenuItem.SetAccelerator)
  - [func \(m \*MenuItem\) SetChecked\(checked bool\) \*MenuItem](#MenuItem.SetChecked)
  - [func \(m \*MenuItem\) SetEnabled\(enabled bool\) \*MenuItem](#MenuItem.SetEnabled)
  - [func \(m \*MenuItem\) SetHidden\(hidden bool\) \*MenuItem](#MenuItem.SetHidden)
  - [func \(m \*MenuItem\) SetLabel\(s string\) \*MenuItem](#MenuItem.SetLabel)
  - [func \(m \*MenuItem\) SetTooltip\(s string\) \*MenuItem](#MenuItem.SetTooltip)
  - [func \(m \*MenuItem\) Tooltip\(\) string](#MenuItem.Tooltip)
- [type MessageDialog](#MessageDialog)
  - [func ErrorDialog\(\) \*MessageDialog](#ErrorDialog)
  - [func InfoDialog\(\) \*MessageDialog](#InfoDialog)
  - [func OpenDirectoryDialog\(\) \*MessageDialog](#OpenDirectoryDialog)
  - [func QuestionDialog\(\) \*MessageDialog](#QuestionDialog)
  - [func WarningDialog\(\) \*MessageDialog](#WarningDialog)
  - [func \(d \*MessageDialog\) AddButton\(s string\) \*Button](#MessageDialog.AddButton)
  - [func \(d \*MessageDialog\) AddButtons\(buttons \[\]\*Button\) \*MessageDialog](#MessageDialog.AddButtons)
  - [func \(d \*MessageDialog\) AttachToWindow\(window \*WebviewWindow\) \*MessageDialog](#MessageDialog.AttachToWindow)
  - [func \(d \*MessageDialog\) SetCancelButton\(button \*Button\) \*MessageDialog](#MessageDialog.SetCancelButton)
  - [func \(d \*MessageDialog\) SetDefaultButton\(button \*Button\) \*MessageDialog](#MessageDialog.SetDefaultButton)
  - [func \(d \*MessageDialog\) SetIcon\(icon \[\]byte\) \*MessageDialog](#MessageDialog.SetIcon)
  - [func \(d \*MessageDialog\) SetMessage\(message string\) \*MessageDialog](#MessageDialog.SetMessage)
  - [func \(d \*MessageDialog\) SetTitle\(title string\) \*MessageDialog](#MessageDialog.SetTitle)
  - [func \(d \*MessageDialog\) Show\(\)](#MessageDialog.Show)
- [type MessageDialogOptions](#MessageDialogOptions)
- [type MessageProcessor](#MessageProcessor)
  - [func NewMessageProcessor\(logger \*slog.Logger\) \*MessageProcessor](#NewMessageProcessor)
  - [func \(m \*MessageProcessor\) Error\(message string, args ...any\)](#MessageProcessor.Error)
  - [func \(m \*MessageProcessor\) HandleRuntimeCall\(rw http.ResponseWriter, r \*http.Request\)](#MessageProcessor.HandleRuntimeCall)
  - [func \(m \*MessageProcessor\) HandleRuntimeCallWithIDs\(rw http.ResponseWriter, r \*http.Request\)](#MessageProcessor.HandleRuntimeCallWithIDs)
  - [func \(m \*MessageProcessor\) Info\(message string, args ...any\)](#MessageProcessor.Info)
- [type Middleware](#Middleware)
  - [func ChainMiddleware\(middleware ...Middleware\) Middleware](#ChainMiddleware)
- [type OpenFileDialogOptions](#OpenFileDialogOptions)
- [type OpenFileDialogStruct](#OpenFileDialogStruct)
  - [func OpenFileDialog\(\) \*OpenFileDialogStruct](#OpenFileDialog)
  - [func OpenFileDialogWithOptions\(options \*OpenFileDialogOptions\) \*OpenFileDialogStruct](#OpenFileDialogWithOptions)
  - [func \(d \*OpenFileDialogStruct\) AddFilter\(displayName, pattern string\) \*OpenFileDialogStruct](#OpenFileDialogStruct.AddFilter)
  - [func \(d \*OpenFileDialogStruct\) AllowsOtherFileTypes\(allowsOtherFileTypes bool\) \*OpenFileDialogStruct](#OpenFileDialogStruct.AllowsOtherFileTypes)
  - [func \(d \*OpenFileDialogStruct\) AttachToWindow\(window \*WebviewWindow\) \*OpenFileDialogStruct](#OpenFileDialogStruct.AttachToWindow)
  - [func \(d \*OpenFileDialogStruct\) CanChooseDirectories\(canChooseDirectories bool\) \*OpenFileDialogStruct](#OpenFileDialogStruct.CanChooseDirectories)
  - [func \(d \*OpenFileDialogStruct\) CanChooseFiles\(canChooseFiles bool\) \*OpenFileDialogStruct](#OpenFileDialogStruct.CanChooseFiles)
  - [func \(d \*OpenFileDialogStruct\) CanCreateDirectories\(canCreateDirectories bool\) \*OpenFileDialogStruct](#OpenFileDialogStruct.CanCreateDirectories)
  - [func \(d \*OpenFileDialogStruct\) CanSelectHiddenExtension\(canSelectHiddenExtension bool\) \*OpenFileDialogStruct](#OpenFileDialogStruct.CanSelectHiddenExtension)
  - [func \(d \*OpenFileDialogStruct\) HideExtension\(hideExtension bool\) \*OpenFileDialogStruct](#OpenFileDialogStruct.HideExtension)
  - [func \(d \*OpenFileDialogStruct\) PromptForMultipleSelection\(\) \(\[\]string, error\)](#OpenFileDialogStruct.PromptForMultipleSelection)
  - [func \(d \*OpenFileDialogStruct\) PromptForSingleSelection\(\) \(string, error\)](#OpenFileDialogStruct.PromptForSingleSelection)
  - [func \(d \*OpenFileDialogStruct\) ResolvesAliases\(resolvesAliases bool\) \*OpenFileDialogStruct](#OpenFileDialogStruct.ResolvesAliases)
  - [func \(d \*OpenFileDialogStruct\) SetButtonText\(text string\) \*OpenFileDialogStruct](#OpenFileDialogStruct.SetButtonText)
  - [func \(d \*OpenFileDialogStruct\) SetDirectory\(directory string\) \*OpenFileDialogStruct](#OpenFileDialogStruct.SetDirectory)
  - [func \(d \*OpenFileDialogStruct\) SetMessage\(message string\) \*OpenFileDialogStruct](#OpenFileDialogStruct.SetMessage)
  - [func \(d \*OpenFileDialogStruct\) SetOptions\(options \*OpenFileDialogOptions\)](#OpenFileDialogStruct.SetOptions)
  - [func \(d \*OpenFileDialogStruct\) SetTitle\(title string\) \*OpenFileDialogStruct](#OpenFileDialogStruct.SetTitle)
  - [func \(d \*OpenFileDialogStruct\) ShowHiddenFiles\(showHiddenFiles bool\) \*OpenFileDialogStruct](#OpenFileDialogStruct.ShowHiddenFiles)
  - [func \(d \*OpenFileDialogStruct\) TreatsFilePackagesAsDirectories\(treatsFilePackagesAsDirectories bool\) \*OpenFileDialogStruct](#OpenFileDialogStruct.TreatsFilePackagesAsDirectories)
- [type Options](#Options)
- [type Parameter](#Parameter)
  - [func \(p \*Parameter\) IsError\(\) bool](#Parameter.IsError)
  - [func \(p \*Parameter\) IsType\(typename string\) bool](#Parameter.IsType)
- [type Plugin](#Plugin)
- [type PluginCallOptions](#PluginCallOptions)
- [type PluginManager](#PluginManager)
  - [func NewPluginManager\(plugins map\[string\]Plugin, assetServer \*assetserver.AssetServer\) \*PluginManager](#NewPluginManager)
  - [func \(p \*PluginManager\) Init\(\) error](#PluginManager.Init)
  - [func \(p \*PluginManager\) Shutdown\(\)](#PluginManager.Shutdown)
- [type PositionOptions](#PositionOptions)
- [type QueryParams](#QueryParams)
  - [func \(qp QueryParams\) Args\(\) \(\*Args, error\)](#QueryParams.Args)
  - [func \(qp QueryParams\) Bool\(key string\) \*bool](#QueryParams.Bool)
  - [func \(qp QueryParams\) Float64\(key string\) \*float64](#QueryParams.Float64)
  - [func \(qp QueryParams\) Int\(key string\) \*int](#QueryParams.Int)
  - [func \(qp QueryParams\) String\(key string\) \*string](#QueryParams.String)
  - [func \(qp QueryParams\) ToStruct\(str any\) error](#QueryParams.ToStruct)
  - [func \(qp QueryParams\) UInt\(key string\) \*uint](#QueryParams.UInt)
  - [func \(qp QueryParams\) UInt8\(key string\) \*uint8](#QueryParams.UInt8)
- [type RGBA](#RGBA)
- [type RadioGroup](#RadioGroup)
  - [func \(r \*RadioGroup\) Add\(id int, item \*MenuItem\)](#RadioGroup.Add)
  - [func \(r \*RadioGroup\) Bounds\(\) \(int, int\)](#RadioGroup.Bounds)
  - [func \(r \*RadioGroup\) MenuID\(item \*MenuItem\) int](#RadioGroup.MenuID)
- [type RadioGroupMember](#RadioGroupMember)
- [type Rect](#Rect)
- [type Role](#Role)
- [type SaveFileDialogOptions](#SaveFileDialogOptions)
- [type SaveFileDialogStruct](#SaveFileDialogStruct)
  - [func SaveFileDialog\(\) \*SaveFileDialogStruct](#SaveFileDialog)
  - [func SaveFileDialogWithOptions\(s \*SaveFileDialogOptions\) \*SaveFileDialogStruct](#SaveFileDialogWithOptions)
  - [func \(d \*SaveFileDialogStruct\) AddFilter\(displayName, pattern string\) \*SaveFileDialogStruct](#SaveFileDialogStruct.AddFilter)
  - [func \(d \*SaveFileDialogStruct\) AllowsOtherFileTypes\(allowOtherFileTypes bool\) \*SaveFileDialogStruct](#SaveFileDialogStruct.AllowsOtherFileTypes)
  - [func \(d \*SaveFileDialogStruct\) AttachToWindow\(window \*WebviewWindow\) \*SaveFileDialogStruct](#SaveFileDialogStruct.AttachToWindow)
  - [func \(d \*SaveFileDialogStruct\) CanCreateDirectories\(canCreateDirectories bool\) \*SaveFileDialogStruct](#SaveFileDialogStruct.CanCreateDirectories)
  - [func \(d \*SaveFileDialogStruct\) CanSelectHiddenExtension\(canSelectHiddenExtension bool\) \*SaveFileDialogStruct](#SaveFileDialogStruct.CanSelectHiddenExtension)
  - [func \(d \*SaveFileDialogStruct\) HideExtension\(hideExtension bool\) \*SaveFileDialogStruct](#SaveFileDialogStruct.HideExtension)
  - [func \(d \*SaveFileDialogStruct\) PromptForSingleSelection\(\) \(string, error\)](#SaveFileDialogStruct.PromptForSingleSelection)
  - [func \(d \*SaveFileDialogStruct\) SetButtonText\(text string\) \*SaveFileDialogStruct](#SaveFileDialogStruct.SetButtonText)
  - [func \(d \*SaveFileDialogStruct\) SetDirectory\(directory string\) \*SaveFileDialogStruct](#SaveFileDialogStruct.SetDirectory)
  - [func \(d \*SaveFileDialogStruct\) SetFilename\(filename string\) \*SaveFileDialogStruct](#SaveFileDialogStruct.SetFilename)
  - [func \(d \*SaveFileDialogStruct\) SetMessage\(message string\) \*SaveFileDialogStruct](#SaveFileDialogStruct.SetMessage)
  - [func \(d \*SaveFileDialogStruct\) SetOptions\(options \*SaveFileDialogOptions\)](#SaveFileDialogStruct.SetOptions)
  - [func \(d \*SaveFileDialogStruct\) ShowHiddenFiles\(showHiddenFiles bool\) \*SaveFileDialogStruct](#SaveFileDialogStruct.ShowHiddenFiles)
  - [func \(d \*SaveFileDialogStruct\) TreatsFilePackagesAsDirectories\(treatsFilePackagesAsDirectories bool\) \*SaveFileDialogStruct](#SaveFileDialogStruct.TreatsFilePackagesAsDirectories)
- [type Screen](#Screen)
- [type Size](#Size)
- [type SystemTray](#SystemTray)
  - [func \(s \*SystemTray\) AttachWindow\(window \*WebviewWindow\) \*SystemTray](#SystemTray.AttachWindow)
  - [func \(s \*SystemTray\) Destroy\(\)](#SystemTray.Destroy)
  - [func \(s \*SystemTray\) Label\(\) string](#SystemTray.Label)
  - [func \(s \*SystemTray\) OnClick\(handler func\(\)\) \*SystemTray](#SystemTray.OnClick)
  - [func \(s \*SystemTray\) OnDoubleClick\(handler func\(\)\) \*SystemTray](#SystemTray.OnDoubleClick)
  - [func \(s \*SystemTray\) OnMouseEnter\(handler func\(\)\) \*SystemTray](#SystemTray.OnMouseEnter)
  - [func \(s \*SystemTray\) OnMouseLeave\(handler func\(\)\) \*SystemTray](#SystemTray.OnMouseLeave)
  - [func \(s \*SystemTray\) OnRightClick\(handler func\(\)\) \*SystemTray](#SystemTray.OnRightClick)
  - [func \(s \*SystemTray\) OnRightDoubleClick\(handler func\(\)\) \*SystemTray](#SystemTray.OnRightDoubleClick)
  - [func \(s \*SystemTray\) OpenMenu\(\)](#SystemTray.OpenMenu)
  - [func \(s \*SystemTray\) PositionWindow\(window \*WebviewWindow, offset int\) error](#SystemTray.PositionWindow)
  - [func \(s \*SystemTray\) SetDarkModeIcon\(icon \[\]byte\) \*SystemTray](#SystemTray.SetDarkModeIcon)
  - [func \(s \*SystemTray\) SetIcon\(icon \[\]byte\) \*SystemTray](#SystemTray.SetIcon)
  - [func \(s \*SystemTray\) SetIconPosition\(iconPosition int\) \*SystemTray](#SystemTray.SetIconPosition)
  - [func \(s \*SystemTray\) SetLabel\(label string\)](#SystemTray.SetLabel)
  - [func \(s \*SystemTray\) SetMenu\(menu \*Menu\) \*SystemTray](#SystemTray.SetMenu)
  - [func \(s \*SystemTray\) SetTemplateIcon\(icon \[\]byte\) \*SystemTray](#SystemTray.SetTemplateIcon)
  - [func \(s \*SystemTray\) WindowDebounce\(debounce time.Duration\) \*SystemTray](#SystemTray.WindowDebounce)
  - [func \(s \*SystemTray\) WindowOffset\(offset int\) \*SystemTray](#SystemTray.WindowOffset)
- [type Theme](#Theme)
- [type ThemeSettings](#ThemeSettings)
- [type WailsEvent](#WailsEvent)
  - [func \(e \*WailsEvent\) Cancel\(\)](#WailsEvent.Cancel)
- [type WebviewWindow](#WebviewWindow)
  - [func \(w \*WebviewWindow\) AbsolutePosition\(\) \(int, int\)](#WebviewWindow.AbsolutePosition)
  - [func \(w \*WebviewWindow\) Center\(\)](#WebviewWindow.Center)
  - [func \(w \*WebviewWindow\) Close\(\)](#WebviewWindow.Close)
  - [func \(w \*WebviewWindow\) Destroy\(\)](#WebviewWindow.Destroy)
  - [func \(w \*WebviewWindow\) ExecJS\(js string\)](#WebviewWindow.ExecJS)
  - [func \(w \*WebviewWindow\) Flash\(enabled bool\)](#WebviewWindow.Flash)
  - [func \(w \*WebviewWindow\) Focus\(\)](#WebviewWindow.Focus)
  - [func \(w \*WebviewWindow\) ForceReload\(\)](#WebviewWindow.ForceReload)
  - [func \(w \*WebviewWindow\) Fullscreen\(\) \*WebviewWindow](#WebviewWindow.Fullscreen)
  - [func \(w \*WebviewWindow\) GetScreen\(\) \(\*Screen, error\)](#WebviewWindow.GetScreen)
  - [func \(w \*WebviewWindow\) GetZoom\(\) float64](#WebviewWindow.GetZoom)
  - [func \(w \*WebviewWindow\) Height\(\) int](#WebviewWindow.Height)
  - [func \(w \*WebviewWindow\) Hide\(\) \*WebviewWindow](#WebviewWindow.Hide)
  - [func \(w \*WebviewWindow\) IsFocused\(\) bool](#WebviewWindow.IsFocused)
  - [func \(w \*WebviewWindow\) IsFullscreen\(\) bool](#WebviewWindow.IsFullscreen)
  - [func \(w \*WebviewWindow\) IsMaximised\(\) bool](#WebviewWindow.IsMaximised)
  - [func \(w \*WebviewWindow\) IsMinimised\(\) bool](#WebviewWindow.IsMinimised)
  - [func \(w \*WebviewWindow\) IsVisible\(\) bool](#WebviewWindow.IsVisible)
  - [func \(w \*WebviewWindow\) Maximise\(\) \*WebviewWindow](#WebviewWindow.Maximise)
  - [func \(w \*WebviewWindow\) Minimise\(\) \*WebviewWindow](#WebviewWindow.Minimise)
  - [func \(w \*WebviewWindow\) Name\(\) string](#WebviewWindow.Name)
  - [func \(w \*WebviewWindow\) NativeWindowHandle\(\) \(uintptr, error\)](#WebviewWindow.NativeWindowHandle)
  - [func \(w \*WebviewWindow\) On\(eventType events.WindowEventType, callback func\(event \*WindowEvent\)\) func\(\)](#WebviewWindow.On)
  - [func \(w \*WebviewWindow\) Print\(\) error](#WebviewWindow.Print)
  - [func \(w \*WebviewWindow\) RegisterContextMenu\(name string, menu \*Menu\)](#WebviewWindow.RegisterContextMenu)
  - [func \(w \*WebviewWindow\) RegisterHook\(eventType events.WindowEventType, callback func\(event \*WindowEvent\)\) func\(\)](#WebviewWindow.RegisterHook)
  - [func \(w \*WebviewWindow\) RelativePosition\(\) \(int, int\)](#WebviewWindow.RelativePosition)
  - [func \(w \*WebviewWindow\) Reload\(\)](#WebviewWindow.Reload)
  - [func \(w \*WebviewWindow\) Resizable\(\) bool](#WebviewWindow.Resizable)
  - [func \(w \*WebviewWindow\) Restore\(\)](#WebviewWindow.Restore)
  - [func \(w \*WebviewWindow\) SetAbsolutePosition\(x int, y int\)](#WebviewWindow.SetAbsolutePosition)
  - [func \(w \*WebviewWindow\) SetAlwaysOnTop\(b bool\) \*WebviewWindow](#WebviewWindow.SetAlwaysOnTop)
  - [func \(w \*WebviewWindow\) SetBackgroundColour\(colour RGBA\) \*WebviewWindow](#WebviewWindow.SetBackgroundColour)
  - [func \(w \*WebviewWindow\) SetEnabled\(enabled bool\)](#WebviewWindow.SetEnabled)
  - [func \(w \*WebviewWindow\) SetFrameless\(frameless bool\) \*WebviewWindow](#WebviewWindow.SetFrameless)
  - [func \(w \*WebviewWindow\) SetFullscreenButtonEnabled\(enabled bool\) \*WebviewWindow](#WebviewWindow.SetFullscreenButtonEnabled)
  - [func \(w \*WebviewWindow\) SetHTML\(html string\) \*WebviewWindow](#WebviewWindow.SetHTML)
  - [func \(w \*WebviewWindow\) SetMaxSize\(maxWidth, maxHeight int\) \*WebviewWindow](#WebviewWindow.SetMaxSize)
  - [func \(w \*WebviewWindow\) SetMinSize\(minWidth, minHeight int\) \*WebviewWindow](#WebviewWindow.SetMinSize)
  - [func \(w \*WebviewWindow\) SetRelativePosition\(x, y int\) \*WebviewWindow](#WebviewWindow.SetRelativePosition)
  - [func \(w \*WebviewWindow\) SetResizable\(b bool\) \*WebviewWindow](#WebviewWindow.SetResizable)
  - [func \(w \*WebviewWindow\) SetSize\(width, height int\) \*WebviewWindow](#WebviewWindow.SetSize)
  - [func \(w \*WebviewWindow\) SetTitle\(title string\) \*WebviewWindow](#WebviewWindow.SetTitle)
  - [func \(w \*WebviewWindow\) SetURL\(s string\) \*WebviewWindow](#WebviewWindow.SetURL)
  - [func \(w \*WebviewWindow\) SetZoom\(magnification float64\) \*WebviewWindow](#WebviewWindow.SetZoom)
  - [func \(w \*WebviewWindow\) Show\(\) \*WebviewWindow](#WebviewWindow.Show)
  - [func \(w \*WebviewWindow\) Size\(\) \(int, int\)](#WebviewWindow.Size)
  - [func \(w \*WebviewWindow\) ToggleDevTools\(\)](#WebviewWindow.ToggleDevTools)
  - [func \(w \*WebviewWindow\) ToggleFullscreen\(\)](#WebviewWindow.ToggleFullscreen)
  - [func \(w \*WebviewWindow\) ToggleMaximise\(\)](#WebviewWindow.ToggleMaximise)
  - [func \(w \*WebviewWindow\) UnFullscreen\(\)](#WebviewWindow.UnFullscreen)
  - [func \(w \*WebviewWindow\) UnMaximise\(\)](#WebviewWindow.UnMaximise)
  - [func \(w \*WebviewWindow\) UnMinimise\(\)](#WebviewWindow.UnMinimise)
  - [func \(w \*WebviewWindow\) Width\(\) int](#WebviewWindow.Width)
  - [func \(w \*WebviewWindow\) Zoom\(\)](#WebviewWindow.Zoom)
  - [func \(w \*WebviewWindow\) ZoomIn\(\)](#WebviewWindow.ZoomIn)
  - [func \(w \*WebviewWindow\) ZoomOut\(\)](#WebviewWindow.ZoomOut)
  - [func \(w \*WebviewWindow\) ZoomReset\(\) \*WebviewWindow](#WebviewWindow.ZoomReset)
- [type WebviewWindowOptions](#WebviewWindowOptions)
- [type Win32Menu](#Win32Menu)
  - [func NewApplicationMenu\(parent w32.HWND, inputMenu \*Menu\) \*Win32Menu](#NewApplicationMenu)
  - [func NewPopupMenu\(parent w32.HWND, inputMenu \*Menu\) \*Win32Menu](#NewPopupMenu)
  - [func \(p \*Win32Menu\) Destroy\(\)](#Win32Menu.Destroy)
  - [func \(p \*Win32Menu\) OnMenuClose\(fn func\(\)\)](#Win32Menu.OnMenuClose)
  - [func \(p \*Win32Menu\) OnMenuOpen\(fn func\(\)\)](#Win32Menu.OnMenuOpen)
  - [func \(p \*Win32Menu\) ProcessCommand\(cmdMsgID int\) bool](#Win32Menu.ProcessCommand)
  - [func \(p \*Win32Menu\) ShowAt\(x int, y int\)](#Win32Menu.ShowAt)
  - [func \(p \*Win32Menu\) ShowAtCursor\(\)](#Win32Menu.ShowAtCursor)
  - [func \(p \*Win32Menu\) Update\(\)](#Win32Menu.Update)
  - [func \(p \*Win32Menu\) UpdateMenuItem\(item \*MenuItem\)](#Win32Menu.UpdateMenuItem)
- [type WindowAttachConfig](#WindowAttachConfig)
- [type WindowEvent](#WindowEvent)
  - [func NewWindowEvent\(\) \*WindowEvent](#NewWindowEvent)
  - [func \(w \*WindowEvent\) Cancel\(\)](#WindowEvent.Cancel)
  - [func \(w \*WindowEvent\) Context\(\) \*WindowEventContext](#WindowEvent.Context)
- [type WindowEventContext](#WindowEventContext)
  - [func \(c WindowEventContext\) DroppedFiles\(\) \[\]string](#WindowEventContext.DroppedFiles)
- [type WindowEventListener](#WindowEventListener)
- [type WindowState](#WindowState)
- [type WindowsOptions](#WindowsOptions)
- [type WindowsWindow](#WindowsWindow)

## Constants

<a name="ApplicationHide"></a>

```go
const (
    ApplicationHide = 0
    ApplicationShow = 1
    ApplicationQuit = 2
)
```

<a name="ClipboardSetText"></a>

```go
const (
    ClipboardSetText = 0
    ClipboardText    = 1
)
```

<a name="DialogInfo"></a>

```go
const (
    DialogInfo     = 0
    DialogWarning  = 1
    DialogError    = 2
    DialogQuestion = 3
    DialogOpenFile = 4
    DialogSaveFile = 5
)
```

<a name="ScreensGetAll"></a>

```go
const (
    ScreensGetAll     = 0
    ScreensGetPrimary = 1
    ScreensGetCurrent = 2
)
```

<a name="WindowCenter"></a>

```go
const (
    WindowCenter              = 0
    WindowSetTitle            = 1
    WindowFullscreen          = 2
    WindowUnFullscreen        = 3
    WindowSetSize             = 4
    WindowSize                = 5
    WindowSetMaxSize          = 6
    WindowSetMinSize          = 7
    WindowSetAlwaysOnTop      = 8
    WindowSetRelativePosition = 9
    WindowRelativePosition    = 10
    WindowScreen              = 11
    WindowHide                = 12
    WindowMaximise            = 13
    WindowUnMaximise          = 14
    WindowToggleMaximise      = 15
    WindowMinimise            = 16
    WindowUnMinimise          = 17
    WindowRestore             = 18
    WindowShow                = 19
    WindowClose               = 20
    WindowSetBackgroundColour = 21
    WindowSetResizable        = 22
    WindowWidth               = 23
    WindowHeight              = 24
    WindowZoomIn              = 25
    WindowZoomOut             = 26
    WindowZoomReset           = 27
    WindowGetZoomLevel        = 28
    WindowSetZoomLevel        = 29
)
```

<a name="NSImageNone"></a>

```go
const (
    NSImageNone = iota
    NSImageOnly
    NSImageLeft
    NSImageRight
    NSImageBelow
    NSImageAbove
    NSImageOverlaps
    NSImageLeading
    NSImageTrailing
)
```

<a name="CallBinding"></a>

```go
const (
    CallBinding = 0
)
```

<a name="ContextMenuOpen"></a>

```go
const (
    ContextMenuOpen = 0
)
```

<a name="EventsEmit"></a>

```go
const (
    EventsEmit = 0
)
```

<a name="MenuItemMsgID"></a>

```go
const (
    MenuItemMsgID = w32.WM_APP + 1024
)
```

<a name="SystemIsDarkMode"></a>

```go
const (
    SystemIsDarkMode = 0
)
```

<a name="WM_USER_SYSTRAY"></a>

```go
const (
    WM_USER_SYSTRAY = w32.WM_USER + 1
)
```

## Variables

<a name="BuildInfo"></a>BuildInfo contains the build info for the application

```go
var BuildInfo *debug.BuildInfo
```

<a name="BuildSettings"></a>BuildSettings contains the build settings for the
application

```go
var BuildSettings map[string]string
```

<a name="MacTitleBarDefault"></a>MacTitleBarDefault results in the default Mac
MacTitleBar

```go
var MacTitleBarDefault = MacTitleBar{
    AppearsTransparent:   false,
    Hide:                 false,
    HideTitle:            false,
    FullSizeContent:      false,
    UseToolbar:           false,
    HideToolbarSeparator: false,
}
```

<a name="MacTitleBarHidden"></a>MacTitleBarHidden results in a hidden title bar
and a full size content window, yet the title bar still has the standard window
controls \(“traffic lights”\) in the top left.

```go
var MacTitleBarHidden = MacTitleBar{
    AppearsTransparent:   true,
    Hide:                 false,
    HideTitle:            true,
    FullSizeContent:      true,
    UseToolbar:           false,
    HideToolbarSeparator: false,
}
```

<a name="MacTitleBarHiddenInset"></a>MacTitleBarHiddenInset results in a hidden
title bar with an alternative look where the traffic light buttons are slightly
more inset from the window edge.

```go
var MacTitleBarHiddenInset = MacTitleBar{
    AppearsTransparent:   true,
    Hide:                 false,
    HideTitle:            true,
    FullSizeContent:      true,
    UseToolbar:           true,
    HideToolbarSeparator: true,
}
```

<a name="MacTitleBarHiddenInsetUnified"></a>MacTitleBarHiddenInsetUnified
results in a hidden title bar with an alternative look where the traffic light
buttons are even more inset from the window edge.

```go
var MacTitleBarHiddenInsetUnified = MacTitleBar{
    AppearsTransparent:   true,
    Hide:                 false,
    HideTitle:            true,
    FullSizeContent:      true,
    UseToolbar:           true,
    HideToolbarSeparator: true,
    ToolbarStyle:         MacToolbarStyleUnified,
}
```

<a name="VirtualKeyCodes"></a>

```go
var VirtualKeyCodes = map[uint]string{
    0x01: "lbutton",
    0x02: "rbutton",
    0x03: "cancel",
    0x04: "mbutton",
    0x05: "xbutton1",
    0x06: "xbutton2",
    0x08: "back",
    0x09: "tab",
    0x0C: "clear",
    0x0D: "return",
    0x10: "shift",
    0x11: "control",
    0x12: "menu",
    0x13: "pause",
    0x14: "capital",
    0x15: "kana",
    0x17: "junja",
    0x18: "final",
    0x19: "hanja",
    0x1B: "escape",
    0x1C: "convert",
    0x1D: "nonconvert",
    0x1E: "accept",
    0x1F: "modechange",
    0x20: "space",
    0x21: "prior",
    0x22: "next",
    0x23: "end",
    0x24: "home",
    0x25: "left",
    0x26: "up",
    0x27: "right",
    0x28: "down",
    0x29: "select",
    0x2A: "print",
    0x2B: "execute",
    0x2C: "snapshot",
    0x2D: "insert",
    0x2E: "delete",
    0x2F: "help",
    0x30: "0",
    0x31: "1",
    0x32: "2",
    0x33: "3",
    0x34: "4",
    0x35: "5",
    0x36: "6",
    0x37: "7",
    0x38: "8",
    0x39: "9",
    0x41: "a",
    0x42: "b",
    0x43: "c",
    0x44: "d",
    0x45: "e",
    0x46: "f",
    0x47: "g",
    0x48: "h",
    0x49: "i",
    0x4A: "j",
    0x4B: "k",
    0x4C: "l",
    0x4D: "m",
    0x4E: "n",
    0x4F: "o",
    0x50: "p",
    0x51: "q",
    0x52: "r",
    0x53: "s",
    0x54: "t",
    0x55: "u",
    0x56: "v",
    0x57: "w",
    0x58: "x",
    0x59: "y",
    0x5A: "z",
    0x5B: "lwin",
    0x5C: "rwin",
    0x5D: "apps",
    0x5F: "sleep",
    0x60: "numpad0",
    0x61: "numpad1",
    0x62: "numpad2",
    0x63: "numpad3",
    0x64: "numpad4",
    0x65: "numpad5",
    0x66: "numpad6",
    0x67: "numpad7",
    0x68: "numpad8",
    0x69: "numpad9",
    0x6A: "multiply",
    0x6B: "add",
    0x6C: "separator",
    0x6D: "subtract",
    0x6E: "decimal",
    0x6F: "divide",
    0x70: "f1",
    0x71: "f2",
    0x72: "f3",
    0x73: "f4",
    0x74: "f5",
    0x75: "f6",
    0x76: "f7",
    0x77: "f8",
    0x78: "f9",
    0x79: "f10",
    0x7A: "f11",
    0x7B: "f12",
    0x7C: "f13",
    0x7D: "f14",
    0x7E: "f15",
    0x7F: "f16",
    0x80: "f17",
    0x81: "f18",
    0x82: "f19",
    0x83: "f20",
    0x84: "f21",
    0x85: "f22",
    0x86: "f23",
    0x87: "f24",
    0x88: "navigation_view",
    0x89: "navigation_menu",
    0x8A: "navigation_up",
    0x8B: "navigation_down",
    0x8C: "navigation_left",
    0x8D: "navigation_right",
    0x8E: "navigation_accept",
    0x8F: "navigation_cancel",
    0x90: "numlock",
    0x91: "scroll",
    0x92: "oem_nec_equal",
    0x93: "oem_fj_masshou",
    0x94: "oem_fj_touroku",
    0x95: "oem_fj_loya",
    0x96: "oem_fj_roya",
    0xA0: "lshift",
    0xA1: "rshift",
    0xA2: "lcontrol",
    0xA3: "rcontrol",
    0xA4: "lmenu",
    0xA5: "rmenu",
    0xA6: "browser_back",
    0xA7: "browser_forward",
    0xA8: "browser_refresh",
    0xA9: "browser_stop",
    0xAA: "browser_search",
    0xAB: "browser_favorites",
    0xAC: "browser_home",
    0xAD: "volume_mute",
    0xAE: "volume_down",
    0xAF: "volume_up",
    0xB0: "media_next_track",
    0xB1: "media_prev_track",
    0xB2: "media_stop",
    0xB3: "media_play_pause",
    0xB4: "launch_mail",
    0xB5: "launch_media_select",
    0xB6: "launch_app1",
    0xB7: "launch_app2",
    0xBA: "oem_1",
    0xBB: "oem_plus",
    0xBC: "oem_comma",
    0xBD: "oem_minus",
    0xBE: "oem_period",
    0xBF: "oem_2",
    0xC0: "oem_3",
    0xC3: "gamepad_a",
    0xC4: "gamepad_b",
    0xC5: "gamepad_x",
    0xC6: "gamepad_y",
    0xC7: "gamepad_right_shoulder",
    0xC8: "gamepad_left_shoulder",
    0xC9: "gamepad_left_trigger",
    0xCA: "gamepad_right_trigger",
    0xCB: "gamepad_dpad_up",
    0xCC: "gamepad_dpad_down",
    0xCD: "gamepad_dpad_left",
    0xCE: "gamepad_dpad_right",
    0xCF: "gamepad_menu",
    0xD0: "gamepad_view",
    0xD1: "gamepad_left_thumbstick_button",
    0xD2: "gamepad_right_thumbstick_button",
    0xD3: "gamepad_left_thumbstick_up",
    0xD4: "gamepad_left_thumbstick_down",
    0xD5: "gamepad_left_thumbstick_right",
    0xD6: "gamepad_left_thumbstick_left",
    0xD7: "gamepad_right_thumbstick_up",
    0xD8: "gamepad_right_thumbstick_down",
    0xD9: "gamepad_right_thumbstick_right",
    0xDA: "gamepad_right_thumbstick_left",
    0xDB: "oem_4",
    0xDC: "oem_5",
    0xDD: "oem_6",
    0xDE: "oem_7",
    0xDF: "oem_8",
    0xE1: "oem_ax",
    0xE2: "oem_102",
    0xE3: "ico_help",
    0xE4: "ico_00",
    0xE5: "processkey",
    0xE6: "ico_clear",
    0xE7: "packet",
    0xE9: "oem_reset",
    0xEA: "oem_jump",
    0xEB: "oem_pa1",
    0xEC: "oem_pa2",
    0xED: "oem_pa3",
    0xEE: "oem_wsctrl",
    0xEF: "oem_cusel",
    0xF0: "oem_attn",
    0xF1: "oem_finish",
    0xF2: "oem_copy",
    0xF3: "oem_auto",
    0xF4: "oem_enlw",
    0xF5: "oem_backtab",
    0xF6: "attn",
    0xF7: "crsel",
    0xF8: "exsel",
    0xF9: "ereof",
    0xFA: "play",
    0xFB: "zoom",
    0xFC: "noname",
    0xFD: "pa1",
    0xFE: "oem_clear",
}
```

<a name="WebviewWindowDefaults"></a>

```go
var WebviewWindowDefaults = &WebviewWindowOptions{
    Title:  "",
    Width:  800,
    Height: 600,
    URL:    "",
    BackgroundColour: RGBA{
        Red:   255,
        Green: 255,
        Blue:  255,
        Alpha: 255,
    },
}
```

<a name="DefaultLogger"></a>

## func [DefaultLogger](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/logger_windows.go#L14)

```go
func DefaultLogger(level slog.Level) *slog.Logger
```

<a name="Fatal"></a>

## func [Fatal](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/errors.go#L8)

```go
func Fatal(message string, args ...interface{})
```

<a name="InvokeAsync"></a>

## func [InvokeAsync](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/mainthread.go#L70)

```go
func InvokeAsync(fn func())
```

<a name="InvokeSync"></a>

## func [InvokeSync](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/mainthread.go#L23)

```go
func InvokeSync(fn func())
```

<a name="InvokeSyncWithError"></a>

## func [InvokeSyncWithError](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/mainthread.go#L46)

```go
func InvokeSyncWithError(fn func() error) (err error)
```

<a name="InvokeSyncWithResult"></a>

## func [InvokeSyncWithResult](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/mainthread.go#L34)

```go
func InvokeSyncWithResult[T any](fn func() T) (res T)
```

<a name="InvokeSyncWithResultAndError"></a>

## func [InvokeSyncWithResultAndError](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/mainthread.go#L58)

```go
func InvokeSyncWithResultAndError[T any](fn func() (T, error)) (res T, err error)
```

<a name="NewIconFromResource"></a>

## func [NewIconFromResource](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window_windows.go#L1603)

```go
func NewIconFromResource(instance w32.HINSTANCE, resId uint16) (w32.HICON, error)
```

<a name="ScaleToDefaultDPI"></a>

## func [ScaleToDefaultDPI](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window_windows.go#L1599)

```go
func ScaleToDefaultDPI(pixels int, dpi uint) int
```

<a name="ScaleWithDPI"></a>

## func [ScaleWithDPI](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window_windows.go#L1595)

```go
func ScaleWithDPI(pixels int, dpi uint) int
```

<a name="ActivationPolicy"></a>

## type [ActivationPolicy](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/options_application_mac.go#L4)

ActivationPolicy is the activation policy for the application.

```go
type ActivationPolicy int
```

<a name="ActivationPolicyRegular"></a>

```go
const (
    // ActivationPolicyRegular is used for applications that have a user interface,
    ActivationPolicyRegular ActivationPolicy = iota
    // ActivationPolicyAccessory is used for applications that do not have a main window,
    // such as system tray applications or background applications.
    ActivationPolicyAccessory
    ActivationPolicyProhibited
)
```

<a name="App"></a>

## type [App](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L214-L269)

```go
type App struct {

    // The main application menu
    ApplicationMenu *Menu

    Events *EventProcessor
    Logger *slog.Logger
    // contains filtered or unexported fields
}
```

<a name="Get"></a>

### func [Get](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L36)

```go
func Get() *App
```

<a name="New"></a>

### func [New](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L40)

```go
func New(appOptions Options) *App
```

<a name="App.Capabilities"></a>

### func \(\*App\) [Capabilities](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L300)

```go
func (a *App) Capabilities() capabilities.Capabilities
```

<a name="App.Clipboard"></a>

### func \(\*App\) [Clipboard](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L639)

```go
func (a *App) Clipboard() *Clipboard
```

<a name="App.CurrentWindow"></a>

### func \(\*App\) [CurrentWindow](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L563)

```go
func (a *App) CurrentWindow() *WebviewWindow
```

<a name="App.GetPID"></a>

### func \(\*App\) [GetPID](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L347)

```go
func (a *App) GetPID() int
```

<a name="App.GetPrimaryScreen"></a>

### func \(\*App\) [GetPrimaryScreen](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L631)

```go
func (a *App) GetPrimaryScreen() (*Screen, error)
```

<a name="App.GetScreens"></a>

### func \(\*App\) [GetScreens](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L635)

```go
func (a *App) GetScreens() ([]*Screen, error)
```

<a name="App.GetWindowByName"></a>

### func \(\*App\) [GetWindowByName](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L716)

```go
func (a *App) GetWindowByName(name string) *WebviewWindow
```

<a name="App.Hide"></a>

### func \(\*App\) [Hide](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L686)

```go
func (a *App) Hide()
```

<a name="App.IsDarkMode"></a>

### func \(\*App\) [IsDarkMode](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L679)

```go
func (a *App) IsDarkMode() bool
```

<a name="App.NewMenu"></a>

### func \(\*App\) [NewMenu](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menu.go#L98)

```go
func (a *App) NewMenu() *Menu
```

<a name="App.NewSystemTray"></a>

### func \(\*App\) [NewSystemTray](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L397)

```go
func (a *App) NewSystemTray() *SystemTray
```

<a name="App.NewWebviewWindow"></a>

### func \(\*App\) [NewWebviewWindow](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L343)

```go
func (a *App) NewWebviewWindow() *WebviewWindow
```

<a name="App.NewWebviewWindowWithOptions"></a>

### func \(\*App\) [NewWebviewWindowWithOptions](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L379)

```go
func (a *App) NewWebviewWindowWithOptions(windowOptions WebviewWindowOptions) *WebviewWindow
```

<a name="App.On"></a>

### func \(\*App\) [On](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L304)

```go
func (a *App) On(eventType events.ApplicationEventType, callback func(event *Event)) func()
```

<a name="App.OnWindowCreation"></a>

### func \(\*App\) [OnWindowCreation](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L712)

```go
func (a *App) OnWindowCreation(callback func(window *WebviewWindow))
```

<a name="App.Quit"></a>

### func \(\*App\) [Quit](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L573)

```go
func (a *App) Quit()
```

<a name="App.RegisterContextMenu"></a>

### func \(\*App\) [RegisterContextMenu](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L698)

```go
func (a *App) RegisterContextMenu(name string, menu *Menu)
```

<a name="App.RegisterHook"></a>

### func \(\*App\) [RegisterHook](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L327)

```go
func (a *App) RegisterHook(eventType events.ApplicationEventType, callback func(event *Event)) func()
```

RegisterHook registers a hook for the given event type. Hooks are called before
the event listeners and can cancel the event. The returned function can be
called to remove the hook.

<a name="App.Run"></a>

### func \(\*App\) [Run](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L410)

```go
func (a *App) Run() error
```

<a name="App.SetMenu"></a>

### func \(\*App\) [SetMenu](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L591)

```go
func (a *App) SetMenu(menu *Menu)
```

<a name="App.Show"></a>

### func \(\*App\) [Show](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L692)

```go
func (a *App) Show()
```

<a name="App.ShowAboutDialog"></a>

### func \(\*App\) [ShowAboutDialog](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L597)

```go
func (a *App) ShowAboutDialog()
```

<a name="ApplicationEventContext"></a>

## type [ApplicationEventContext](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/context_application_event.go#L10-L13)

```go
type ApplicationEventContext struct {
    // contains filtered or unexported fields
}
```

<a name="ApplicationEventContext.IsDarkMode"></a>

### func \(ApplicationEventContext\) [IsDarkMode](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/context_application_event.go#L35)

```go
func (c ApplicationEventContext) IsDarkMode() bool
```

<a name="ApplicationEventContext.OpenedFiles"></a>

### func \(ApplicationEventContext\) [OpenedFiles](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/context_application_event.go#L15)

```go
func (c ApplicationEventContext) OpenedFiles() []string
```

<a name="Args"></a>

## type [Args](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor_params.go#L109-L111)

```go
type Args struct {
    // contains filtered or unexported fields
}
```

<a name="Args.Bool"></a>

### func \(\*Args\) [Bool](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor_params.go#L178)

```go
func (a *Args) Bool(s string) *bool
```

<a name="Args.Float64"></a>

### func \(\*Args\) [Float64](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor_params.go#L167)

```go
func (a *Args) Float64(s string) *float64
```

<a name="Args.Int"></a>

### func \(\*Args\) [Int](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor_params.go#L124)

```go
func (a *Args) Int(s string) *int
```

<a name="Args.String"></a>

### func \(\*Args\) [String](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor_params.go#L113)

```go
func (a *Args) String(key string) *string
```

<a name="Args.UInt"></a>

### func \(\*Args\) [UInt](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor_params.go#L157)

```go
func (a *Args) UInt(s string) *uint
```

<a name="Args.UInt8"></a>

### func \(\*Args\) [UInt8](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor_params.go#L148)

```go
func (a *Args) UInt8(s string) *uint8
```

<a name="AssetOptions"></a>

## type [AssetOptions](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/options_application.go#L57-L87)

AssetOptions defines the configuration of the AssetServer.

```go
type AssetOptions struct {
    // Handler which serves all the content to the WebView.
	Handler http.Handler

	// Middleware is a HTTP Middleware which allows to hook into the AssetServer request chain. It allows to skip the default
	// request handler dynamically, e.g. implement specialized Routing etc.
	// The Middleware is called to build a new `http.Handler` used by the AssetSever and it also receives the default
	// handler used by the AssetServer as an argument.
	//
	// This middleware injects itself before any of Wails internal middlewares.
	//
	// If not defined, the default AssetServer request chain is executed.
	//
	// Multiple Middlewares can be chained together with:
	//   ChainMiddleware(middleware ...Middleware) Middleware
	Middleware Middleware

    // External URL can be set to a development server URL so that all requests are forwarded to it. This is useful
    // when using a development server like `vite` or `snowpack` which serves the assets on a different port.
    ExternalURL string
}
```

<a name="BackdropType"></a>

## type [BackdropType](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/options_webview_window_win.go#L5)

```go
type BackdropType int32
```

<a name="Auto"></a>

```go
const (
    Auto    BackdropType = 0
    None    BackdropType = 1
    Mica    BackdropType = 2
    Acrylic BackdropType = 3
    Tabbed  BackdropType = 4
)
```

<a name="BackgroundType"></a>

## type [BackgroundType](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/options_webview_window.go#L137)

```go
type BackgroundType int
```

<a name="BackgroundTypeSolid"></a>

```go
const (
    BackgroundTypeSolid BackgroundType = iota
    BackgroundTypeTransparent
    BackgroundTypeTranslucent
)
```

<a name="Bindings"></a>

## type [Bindings](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/bindings.go#L76-L80)

```go
type Bindings struct {
    // contains filtered or unexported fields
}
```

<a name="NewBindings"></a>

### func [NewBindings](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/bindings.go#L82)

```go
func NewBindings(structs []any, aliases map[uint32]uint32) (*Bindings, error)
```

<a name="Bindings.Add"></a>

### func \(\*Bindings\) [Add](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/bindings.go#L98)

```go
func (b *Bindings) Add(structPtr interface{}) error
```

Add the given struct methods to the Bindings

<a name="Bindings.AddPlugins"></a>

### func \(\*Bindings\) [AddPlugins](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/bindings.go#L123)

```go
func (b *Bindings) AddPlugins(plugins map[string]Plugin) error
```

<a name="Bindings.GenerateID"></a>

### func \(\*Bindings\) [GenerateID](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/bindings.go#L190)

```go
func (b *Bindings) GenerateID(name string) (uint32, error)
```

GenerateID generates a unique ID for a binding

<a name="Bindings.Get"></a>

### func \(\*Bindings\) [Get](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/bindings.go#L161)

```go
func (b *Bindings) Get(options *CallOptions) *BoundMethod
```

Get returns the bound method with the given name

<a name="Bindings.GetByID"></a>

### func \(\*Bindings\) [GetByID](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/bindings.go#L178)

```go
func (b *Bindings) GetByID(id uint32) *BoundMethod
```

GetByID returns the bound method with the given ID

<a name="BoundMethod"></a>

## type [BoundMethod](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/bindings.go#L64-L74)

BoundMethod defines all the data related to a Go method that is bound to the
Wails application

```go
type BoundMethod struct {
    ID          uint32        `json:"id"`
    Name        string        `json:"name"`
    Inputs      []*Parameter  `json:"inputs,omitempty"`
    Outputs     []*Parameter  `json:"outputs,omitempty"`
    Comments    string        `json:"comments,omitempty"`
    Method      reflect.Value `json:"-"`
    PackageName string
    StructName  string
    PackagePath string
}
```

<a name="BoundMethod.Call"></a>

### func \(\*BoundMethod\) [Call](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/bindings.go#L298)

```go
func (b *BoundMethod) Call(args []interface{}) (returnValue interface{}, err error)
```

Call will attempt to call this bound method with the given args

<a name="BoundMethod.String"></a>

### func \(\*BoundMethod\) [String](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/bindings.go#L203)

```go
func (b *BoundMethod) String() string
```

<a name="Button"></a>

## type [Button](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L47-L52)

```go
type Button struct {
    Label     string
    IsCancel  bool
    IsDefault bool
    Callback  func()
}
```

<a name="Button.OnClick"></a>

### func \(\*Button\) [OnClick](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L54)

```go
func (b *Button) OnClick(callback func()) *Button
```

<a name="Button.SetAsCancel"></a>

### func \(\*Button\) [SetAsCancel](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L64)

```go
func (b *Button) SetAsCancel() *Button
```

<a name="Button.SetAsDefault"></a>

### func \(\*Button\) [SetAsDefault](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L59)

```go
func (b *Button) SetAsDefault() *Button
```

<a name="CallOptions"></a>

## type [CallOptions](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/bindings.go#L13-L19)

```go
type CallOptions struct {
    MethodID    uint32 `json:"methodID"`
    PackageName string `json:"packageName"`
    StructName  string `json:"structName"`
    MethodName  string `json:"methodName"`
    Args        []any  `json:"args"`
}
```

<a name="CallOptions.Name"></a>

### func \(CallOptions\) [Name](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/bindings.go#L21)

```go
func (c CallOptions) Name() string
```

<a name="Clipboard"></a>

## type [Clipboard](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/clipboard.go#L8-L10)

```go
type Clipboard struct {
    // contains filtered or unexported fields
}
```

<a name="Clipboard.SetText"></a>

### func \(\*Clipboard\) [SetText](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/clipboard.go#L18)

```go
func (c *Clipboard) SetText(text string) bool
```

<a name="Clipboard.Text"></a>

### func \(\*Clipboard\) [Text](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/clipboard.go#L22)

```go
func (c *Clipboard) Text() (string, bool)
```

<a name="Context"></a>

## type [Context](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/context.go#L3-L6)

```go
type Context struct {
    // contains filtered or unexported fields
}
```

<a name="Context.ClickedMenuItem"></a>

### func \(\*Context\) [ClickedMenuItem](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/context.go#L20)

```go
func (c *Context) ClickedMenuItem() *MenuItem
```

<a name="Context.ContextMenuData"></a>

### func \(\*Context\) [ContextMenuData](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/context.go#L35)

```go
func (c *Context) ContextMenuData() any
```

<a name="Context.IsChecked"></a>

### func \(\*Context\) [IsChecked](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/context.go#L28)

```go
func (c *Context) IsChecked() bool
```

<a name="ContextMenuData"></a>

## type [ContextMenuData](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor_contextmenu.go#L7-L12)

```go
type ContextMenuData struct {
    Id   string `json:"id"`
    X    int    `json:"x"`
    Y    int    `json:"y"`
    Data any    `json:"data"`
}
```

<a name="DialogType"></a>

## type [DialogType](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L8)

```go
type DialogType int
```

<a name="InfoDialogType"></a>

```go
const (
    InfoDialogType DialogType = iota
    QuestionDialogType
    WarningDialogType
    ErrorDialogType
    OpenDirectoryDialogType
)
```

<a name="Event"></a>

## type [Event](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/events.go#L10-L14)

```go
type Event struct {
    Id  uint

    Cancelled bool
    // contains filtered or unexported fields
}
```

<a name="Event.Cancel"></a>

### func \(\*Event\) [Cancel](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/events.go#L27)

```go
func (w *Event) Cancel()
```

<a name="Event.Context"></a>

### func \(\*Event\) [Context](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/events.go#L16)

```go
func (w *Event) Context() *ApplicationEventContext
```

<a name="EventListener"></a>

## type [EventListener](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L32-L34)

```go
type EventListener struct {
    // contains filtered or unexported fields
}
```

<a name="EventProcessor"></a>

## type [EventProcessor](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/events.go#L77-L84)

EventProcessor handles custom events

```go
type EventProcessor struct {
    // contains filtered or unexported fields
}
```

<a name="NewWailsEventProcessor"></a>

### func [NewWailsEventProcessor](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/events.go#L86)

```go
func NewWailsEventProcessor(dispatchEventToWindows func(*WailsEvent)) *EventProcessor
```

<a name="EventProcessor.Emit"></a>

### func \(\*EventProcessor\) [Emit](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/events.go#L110)

```go
func (e *EventProcessor) Emit(thisEvent *WailsEvent)
```

Emit sends an event to all listeners

<a name="EventProcessor.Off"></a>

### func \(\*EventProcessor\) [Off](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/events.go#L131)

```go
func (e *EventProcessor) Off(eventName string)
```

<a name="EventProcessor.OffAll"></a>

### func \(\*EventProcessor\) [OffAll](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/events.go#L135)

```go
func (e *EventProcessor) OffAll()
```

<a name="EventProcessor.On"></a>

### func \(\*EventProcessor\) [On](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/events.go#L95)

```go
func (e *EventProcessor) On(eventName string, callback func(event *WailsEvent)) func()
```

On is the equivalent of Javascript's \`addEventListener\`

<a name="EventProcessor.OnMultiple"></a>

### func \(\*EventProcessor\) [OnMultiple](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/events.go#L100)

```go
func (e *EventProcessor) OnMultiple(eventName string, callback func(event *WailsEvent), counter int) func()
```

OnMultiple is the same as \`On\` but will unregister after \`count\` events

<a name="EventProcessor.Once"></a>

### func \(\*EventProcessor\) [Once](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/events.go#L105)

```go
func (e *EventProcessor) Once(eventName string, callback func(event *WailsEvent)) func()
```

Once is the same as \`On\` but will unregister after the first event

<a name="EventProcessor.RegisterHook"></a>

### func \(\*EventProcessor\) [RegisterHook](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/events.go#L167)

```go
func (e *EventProcessor) RegisterHook(eventName string, callback func(*WailsEvent)) func()
```

RegisterHook provides a means of registering methods to be called before
emitting the event

<a name="FileFilter"></a>

## type [FileFilter](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L165-L168)

```go
type FileFilter struct {
    DisplayName string // Filter information EG: "Image Files (*.jpg, *.png)"
    Pattern     string // semicolon separated list of extensions, EG: "*.jpg;*.png"
}
```

<a name="IconPosition"></a>

## type [IconPosition](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L12)

```go
type IconPosition int
```

<a name="MacAppearanceType"></a>

## type [MacAppearanceType](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/options_mac.go#L119)

MacAppearanceType is a type of Appearance for Cocoa windows

```go
type MacAppearanceType string
```

<a name="DefaultAppearance"></a>

```go
const (
    // DefaultAppearance uses the default system value
    DefaultAppearance MacAppearanceType = ""
    // NSAppearanceNameAqua - The standard light system appearance.
    NSAppearanceNameAqua MacAppearanceType = "NSAppearanceNameAqua"
    // NSAppearanceNameDarkAqua - The standard dark system appearance.
    NSAppearanceNameDarkAqua MacAppearanceType = "NSAppearanceNameDarkAqua"
    // NSAppearanceNameVibrantLight - The light vibrant appearance
    NSAppearanceNameVibrantLight MacAppearanceType = "NSAppearanceNameVibrantLight"
    // NSAppearanceNameAccessibilityHighContrastAqua - A high-contrast version of the standard light system appearance.
    NSAppearanceNameAccessibilityHighContrastAqua MacAppearanceType = "NSAppearanceNameAccessibilityHighContrastAqua"
    // NSAppearanceNameAccessibilityHighContrastDarkAqua - A high-contrast version of the standard dark system appearance.
    NSAppearanceNameAccessibilityHighContrastDarkAqua MacAppearanceType = "NSAppearanceNameAccessibilityHighContrastDarkAqua"
    // NSAppearanceNameAccessibilityHighContrastVibrantLight - A high-contrast version of the light vibrant appearance.
    NSAppearanceNameAccessibilityHighContrastVibrantLight MacAppearanceType = "NSAppearanceNameAccessibilityHighContrastVibrantLight"
    // NSAppearanceNameAccessibilityHighContrastVibrantDark - A high-contrast version of the dark vibrant appearance.
    NSAppearanceNameAccessibilityHighContrastVibrantDark MacAppearanceType = "NSAppearanceNameAccessibilityHighContrastVibrantDark"
)
```

<a name="MacBackdrop"></a>

## type [MacBackdrop](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/options_mac.go#L6)

MacBackdrop is the backdrop type for macOS

```go
type MacBackdrop int
```

<a name="MacBackdropNormal"></a>

```go
const (
    // MacBackdropNormal - The default value. The window will have a normal opaque background.
    MacBackdropNormal MacBackdrop = iota
    // MacBackdropTransparent - The window will have a transparent background, with the content underneath it being visible
    MacBackdropTransparent
    // MacBackdropTranslucent - The window will have a translucent background, with the content underneath it being "fuzzy" or "frosted"
    MacBackdropTranslucent
)
```

<a name="MacOptions"></a>

## type [MacOptions](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/options_application_mac.go#L16-L22)

MacOptions contains options for macOS applications.

```go
type MacOptions struct {
    // ActivationPolicy is the activation policy for the application. Defaults to
    // applicationActivationPolicyRegular.
    ActivationPolicy ActivationPolicy
    // If set to true, the application will terminate when the last window is closed.
    ApplicationShouldTerminateAfterLastWindowClosed bool
}
```

<a name="MacTitleBar"></a>

## type [MacTitleBar](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/options_mac.go#L54-L69)

MacTitleBar contains options for the Mac titlebar

```go
type MacTitleBar struct {
    // AppearsTransparent will make the titlebar transparent
    AppearsTransparent bool
    // Hide will hide the titlebar
    Hide bool
    // HideTitle will hide the title
    HideTitle bool
    // FullSizeContent will extend the window content to the full size of the window
    FullSizeContent bool
    // UseToolbar will use a toolbar instead of a titlebar
    UseToolbar bool
    // HideToolbarSeparator will hide the toolbar separator
    HideToolbarSeparator bool
    // ShowToolbarWhenFullscreen will keep the toolbar visible when the window is in fullscreen mode
	ShowToolbarWhenFullscreen bool
    // ToolbarStyle is the style of toolbar to use
    ToolbarStyle MacToolbarStyle
}
```

<a name="MacToolbarStyle"></a>

## type [MacToolbarStyle](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/options_mac.go#L18)

MacToolbarStyle is the style of toolbar for macOS

```go
type MacToolbarStyle int
```

<a name="MacToolbarStyleAutomatic"></a>

```go
const (
    // MacToolbarStyleAutomatic - The default value. The style will be determined by the window's given configuration
    MacToolbarStyleAutomatic MacToolbarStyle = iota
    // MacToolbarStyleExpanded - The toolbar will appear below the window title
    MacToolbarStyleExpanded
    // MacToolbarStylePreference - The toolbar will appear below the window title and the items in the toolbar will attempt to have equal widths when possible
    MacToolbarStylePreference
    // MacToolbarStyleUnified - The window title will appear inline with the toolbar when visible
    MacToolbarStyleUnified
    // MacToolbarStyleUnifiedCompact - Same as MacToolbarStyleUnified, but with reduced margins in the toolbar allowing more focus to be on the contents of the window
    MacToolbarStyleUnifiedCompact
)
```

<a name="MacWindow"></a>

## type [MacWindow](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/options_mac.go#L34-L51)

MacWindow contains macOS specific options for Webview Windows

```go
type MacWindow struct {
    // Backdrop is the backdrop type for the window
    Backdrop MacBackdrop
    // DisableShadow will disable the window shadow
    DisableShadow bool
    // TitleBar contains options for the Mac titlebar
    TitleBar MacTitleBar
    // Appearance is the appearance type for the window
    Appearance MacAppearanceType
    // InvisibleTitleBarHeight defines the height of an invisible titlebar which responds to dragging
    InvisibleTitleBarHeight int
    // Maps events from platform specific to common event types
    EventMapping map[events.WindowEventType]events.WindowEventType

    // EnableFraudulentWebsiteWarnings will enable warnings for fraudulent websites.
    // Default: false
    EnableFraudulentWebsiteWarnings bool
}
```

<a name="Menu"></a>

## type [Menu](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menu.go#L7-L12)

```go
type Menu struct {
    // contains filtered or unexported fields
}
```

<a name="NewMenu"></a>

### func [NewMenu](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menu.go#L14)

```go
func NewMenu() *Menu
```

<a name="Menu.Add"></a>

### func \(\*Menu\) [Add](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menu.go#L18)

```go
func (m *Menu) Add(label string) *MenuItem
```

<a name="Menu.AddCheckbox"></a>

### func \(\*Menu\) [AddCheckbox](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menu.go#L29)

```go
func (m *Menu) AddCheckbox(label string, enabled bool) *MenuItem
```

<a name="Menu.AddRadio"></a>

### func \(\*Menu\) [AddRadio](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menu.go#L35)

```go
func (m *Menu) AddRadio(label string, enabled bool) *MenuItem
```

<a name="Menu.AddRole"></a>

### func \(\*Menu\) [AddRole](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menu.go#L55)

```go
func (m *Menu) AddRole(role Role) *Menu
```

<a name="Menu.AddSeparator"></a>

### func \(\*Menu\) [AddSeparator](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menu.go#L24)

```go
func (m *Menu) AddSeparator()
```

<a name="Menu.AddSubmenu"></a>

### func \(\*Menu\) [AddSubmenu](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menu.go#L49)

```go
func (m *Menu) AddSubmenu(s string) *Menu
```

<a name="Menu.SetLabel"></a>

### func \(\*Menu\) [SetLabel](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menu.go#L88)

```go
func (m *Menu) SetLabel(label string)
```

<a name="Menu.Update"></a>

### func \(\*Menu\) [Update](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menu.go#L41)

```go
func (m *Menu) Update()
```

<a name="MenuItem"></a>

## type [MenuItem](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menuitem.go#L45-L61)

```go
type MenuItem struct {
    // contains filtered or unexported fields
}
```

<a name="MenuItem.Checked"></a>

### func \(\*MenuItem\) [Checked](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menuitem.go#L273)

```go
func (m *MenuItem) Checked() bool
```

<a name="MenuItem.Enabled"></a>

### func \(\*MenuItem\) [Enabled](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menuitem.go#L310)

```go
func (m *MenuItem) Enabled() bool
```

<a name="MenuItem.Hidden"></a>

### func \(\*MenuItem\) [Hidden](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menuitem.go#L293)

```go
func (m *MenuItem) Hidden() bool
```

<a name="MenuItem.IsCheckbox"></a>

### func \(\*MenuItem\) [IsCheckbox](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menuitem.go#L285)

```go
func (m *MenuItem) IsCheckbox() bool
```

<a name="MenuItem.IsRadio"></a>

### func \(\*MenuItem\) [IsRadio](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menuitem.go#L289)

```go
func (m *MenuItem) IsRadio() bool
```

<a name="MenuItem.IsSeparator"></a>

### func \(\*MenuItem\) [IsSeparator](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menuitem.go#L277)

```go
func (m *MenuItem) IsSeparator() bool
```

<a name="MenuItem.IsSubmenu"></a>

### func \(\*MenuItem\) [IsSubmenu](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menuitem.go#L281)

```go
func (m *MenuItem) IsSubmenu() bool
```

<a name="MenuItem.Label"></a>

### func \(\*MenuItem\) [Label](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menuitem.go#L302)

```go
func (m *MenuItem) Label() string
```

<a name="MenuItem.OnClick"></a>

### func \(\*MenuItem\) [OnClick](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menuitem.go#L297)

```go
func (m *MenuItem) OnClick(f func(*Context)) *MenuItem
```

<a name="MenuItem.SetAccelerator"></a>

### func \(\*MenuItem\) [SetAccelerator](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menuitem.go#L220)

```go
func (m *MenuItem) SetAccelerator(shortcut string) *MenuItem
```

<a name="MenuItem.SetChecked"></a>

### func \(\*MenuItem\) [SetChecked](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menuitem.go#L257)

```go
func (m *MenuItem) SetChecked(checked bool) *MenuItem
```

<a name="MenuItem.SetEnabled"></a>

### func \(\*MenuItem\) [SetEnabled](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menuitem.go#L249)

```go
func (m *MenuItem) SetEnabled(enabled bool) *MenuItem
```

<a name="MenuItem.SetHidden"></a>

### func \(\*MenuItem\) [SetHidden](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menuitem.go#L265)

```go
func (m *MenuItem) SetHidden(hidden bool) *MenuItem
```

<a name="MenuItem.SetLabel"></a>

### func \(\*MenuItem\) [SetLabel](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menuitem.go#L241)

```go
func (m *MenuItem) SetLabel(s string) *MenuItem
```

<a name="MenuItem.SetTooltip"></a>

### func \(\*MenuItem\) [SetTooltip](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menuitem.go#L233)

```go
func (m *MenuItem) SetTooltip(s string) *MenuItem
```

<a name="MenuItem.Tooltip"></a>

### func \(\*MenuItem\) [Tooltip](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/menuitem.go#L306)

```go
func (m *MenuItem) Tooltip() string
```

<a name="MessageDialog"></a>

## type [MessageDialog](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L82-L87)

```go
type MessageDialog struct {
    MessageDialogOptions
    // contains filtered or unexported fields
}
```

<a name="ErrorDialog"></a>

### func [ErrorDialog](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L615)

```go
func ErrorDialog() *MessageDialog
```

<a name="InfoDialog"></a>

### func [InfoDialog](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L603)

```go
func InfoDialog() *MessageDialog
```

<a name="OpenDirectoryDialog"></a>

### func [OpenDirectoryDialog](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L619)

```go
func OpenDirectoryDialog() *MessageDialog
```

<a name="QuestionDialog"></a>

### func [QuestionDialog](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L607)

```go
func QuestionDialog() *MessageDialog
```

<a name="WarningDialog"></a>

### func [WarningDialog](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L611)

```go
func WarningDialog() *MessageDialog
```

<a name="MessageDialog.AddButton"></a>

### func \(\*MessageDialog\) [AddButton](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L122)

```go
func (d *MessageDialog) AddButton(s string) *Button
```

<a name="MessageDialog.AddButtons"></a>

### func \(\*MessageDialog\) [AddButtons](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L130)

```go
func (d *MessageDialog) AddButtons(buttons []*Button) *MessageDialog
```

<a name="MessageDialog.AttachToWindow"></a>

### func \(\*MessageDialog\) [AttachToWindow](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L135)

```go
func (d *MessageDialog) AttachToWindow(window *WebviewWindow) *MessageDialog
```

<a name="MessageDialog.SetCancelButton"></a>

### func \(\*MessageDialog\) [SetCancelButton](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L148)

```go
func (d *MessageDialog) SetCancelButton(button *Button) *MessageDialog
```

<a name="MessageDialog.SetDefaultButton"></a>

### func \(\*MessageDialog\) [SetDefaultButton](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L140)

```go
func (d *MessageDialog) SetDefaultButton(button *Button) *MessageDialog
```

<a name="MessageDialog.SetIcon"></a>

### func \(\*MessageDialog\) [SetIcon](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L117)

```go
func (d *MessageDialog) SetIcon(icon []byte) *MessageDialog
```

<a name="MessageDialog.SetMessage"></a>

### func \(\*MessageDialog\) [SetMessage](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L156)

```go
func (d *MessageDialog) SetMessage(message string) *MessageDialog
```

<a name="MessageDialog.SetTitle"></a>

### func \(\*MessageDialog\) [SetTitle](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L105)

```go
func (d *MessageDialog) SetTitle(title string) *MessageDialog
```

<a name="MessageDialog.Show"></a>

### func \(\*MessageDialog\) [Show](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L110)

```go
func (d *MessageDialog) Show()
```

<a name="MessageDialogOptions"></a>

## type [MessageDialogOptions](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L73-L80)

```go
type MessageDialogOptions struct {
    DialogType DialogType
    Title      string
    Message    string
    Buttons    []*Button
    Icon       []byte
    // contains filtered or unexported fields
}
```

<a name="MessageProcessor"></a>

## type [MessageProcessor](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor.go#L26-L29)

```go
type MessageProcessor struct {
    // contains filtered or unexported fields
}
```

<a name="NewMessageProcessor"></a>

### func [NewMessageProcessor](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor.go#L31)

```go
func NewMessageProcessor(logger *slog.Logger) *MessageProcessor
```

<a name="MessageProcessor.Error"></a>

### func \(\*MessageProcessor\) [Error](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor.go#L147)

```go
func (m *MessageProcessor) Error(message string, args ...any)
```

<a name="MessageProcessor.HandleRuntimeCall"></a>

### func \(\*MessageProcessor\) [HandleRuntimeCall](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor.go#L65)

```go
func (m *MessageProcessor) HandleRuntimeCall(rw http.ResponseWriter, r *http.Request)
```

<a name="MessageProcessor.HandleRuntimeCallWithIDs"></a>

### func \(\*MessageProcessor\) [HandleRuntimeCallWithIDs](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor.go#L104)

```go
func (m *MessageProcessor) HandleRuntimeCallWithIDs(rw http.ResponseWriter, r *http.Request)
```

<a name="MessageProcessor.Info"></a>

### func \(\*MessageProcessor\) [Info](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor.go#L151)

```go
func (m *MessageProcessor) Info(message string, args ...any)
```

<a name="Middleware"></a>

## type [Middleware](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/options_application.go#L92)

Middleware defines HTTP middleware that can be applied to the AssetServer. The
handler passed as next is the next handler in the chain. One can decide to call
the next handler or implement a specialized handling.

```go
type Middleware func(next http.Handler) http.Handler
```

<a name="ChainMiddleware"></a>

### func [ChainMiddleware](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/options_application.go#L95)

```go
func ChainMiddleware(middleware ...Middleware) Middleware
```

ChainMiddleware allows chaining multiple middlewares to one middleware.

<a name="OpenFileDialogOptions"></a>

## type [OpenFileDialogOptions](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L170-L188)

```go
type OpenFileDialogOptions struct {
    CanChooseDirectories            bool
    CanChooseFiles                  bool
    CanCreateDirectories            bool
    ShowHiddenFiles                 bool
    ResolvesAliases                 bool
    AllowsMultipleSelection         bool
    HideExtension                   bool
    CanSelectHiddenExtension        bool
    TreatsFilePackagesAsDirectories bool
    AllowsOtherFileTypes            bool
    Filters                         []FileFilter
    Window                          *WebviewWindow

    Title      string
    Message    string
    ButtonText string
    Directory  string
}
```

<a name="OpenFileDialogStruct"></a>

## type [OpenFileDialogStruct](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L190-L211)

```go
type OpenFileDialogStruct struct {
    // contains filtered or unexported fields
}
```

<a name="OpenFileDialog"></a>

### func [OpenFileDialog](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L623)

```go
func OpenFileDialog() *OpenFileDialogStruct
```

<a name="OpenFileDialogWithOptions"></a>

### func [OpenFileDialogWithOptions](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L661)

```go
func OpenFileDialogWithOptions(options *OpenFileDialogOptions) *OpenFileDialogStruct
```

<a name="OpenFileDialogStruct.AddFilter"></a>

### func \(\*OpenFileDialogStruct\) [AddFilter](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L279)

```go
func (d *OpenFileDialogStruct) AddFilter(displayName, pattern string) *OpenFileDialogStruct
```

AddFilter adds a filter to the dialog. The filter is a display name and a
semicolon separated list of extensions. EG: AddFilter\("Image Files",
"\*.jpg;\*.png"\)

<a name="OpenFileDialogStruct.AllowsOtherFileTypes"></a>

### func \(\*OpenFileDialogStruct\) [AllowsOtherFileTypes](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L228)

```go
func (d *OpenFileDialogStruct) AllowsOtherFileTypes(allowsOtherFileTypes bool) *OpenFileDialogStruct
```

<a name="OpenFileDialogStruct.AttachToWindow"></a>

### func \(\*OpenFileDialogStruct\) [AttachToWindow](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L248)

```go
func (d *OpenFileDialogStruct) AttachToWindow(window *WebviewWindow) *OpenFileDialogStruct
```

<a name="OpenFileDialogStruct.CanChooseDirectories"></a>

### func \(\*OpenFileDialogStruct\) [CanChooseDirectories](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L218)

```go
func (d *OpenFileDialogStruct) CanChooseDirectories(canChooseDirectories bool) *OpenFileDialogStruct
```

<a name="OpenFileDialogStruct.CanChooseFiles"></a>

### func \(\*OpenFileDialogStruct\) [CanChooseFiles](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L213)

```go
func (d *OpenFileDialogStruct) CanChooseFiles(canChooseFiles bool) *OpenFileDialogStruct
```

<a name="OpenFileDialogStruct.CanCreateDirectories"></a>

### func \(\*OpenFileDialogStruct\) [CanCreateDirectories](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L223)

```go
func (d *OpenFileDialogStruct) CanCreateDirectories(canCreateDirectories bool) *OpenFileDialogStruct
```

<a name="OpenFileDialogStruct.CanSelectHiddenExtension"></a>

### func \(\*OpenFileDialogStruct\) [CanSelectHiddenExtension](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L310)

```go
func (d *OpenFileDialogStruct) CanSelectHiddenExtension(canSelectHiddenExtension bool) *OpenFileDialogStruct
```

<a name="OpenFileDialogStruct.HideExtension"></a>

### func \(\*OpenFileDialogStruct\) [HideExtension](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L238)

```go
func (d *OpenFileDialogStruct) HideExtension(hideExtension bool) *OpenFileDialogStruct
```

<a name="OpenFileDialogStruct.PromptForMultipleSelection"></a>

### func \(\*OpenFileDialogStruct\) [PromptForMultipleSelection](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L287)

```go
func (d *OpenFileDialogStruct) PromptForMultipleSelection() ([]string, error)
```

<a name="OpenFileDialogStruct.PromptForSingleSelection"></a>

### func \(\*OpenFileDialogStruct\) [PromptForSingleSelection](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L263)

```go
func (d *OpenFileDialogStruct) PromptForSingleSelection() (string, error)
```

<a name="OpenFileDialogStruct.ResolvesAliases"></a>

### func \(\*OpenFileDialogStruct\) [ResolvesAliases](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L253)

```go
func (d *OpenFileDialogStruct) ResolvesAliases(resolvesAliases bool) *OpenFileDialogStruct
```

<a name="OpenFileDialogStruct.SetButtonText"></a>

### func \(\*OpenFileDialogStruct\) [SetButtonText](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L300)

```go
func (d *OpenFileDialogStruct) SetButtonText(text string) *OpenFileDialogStruct
```

<a name="OpenFileDialogStruct.SetDirectory"></a>

### func \(\*OpenFileDialogStruct\) [SetDirectory](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L305)

```go
func (d *OpenFileDialogStruct) SetDirectory(directory string) *OpenFileDialogStruct
```

<a name="OpenFileDialogStruct.SetMessage"></a>

### func \(\*OpenFileDialogStruct\) [SetMessage](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L295)

```go
func (d *OpenFileDialogStruct) SetMessage(message string) *OpenFileDialogStruct
```

<a name="OpenFileDialogStruct.SetOptions"></a>

### func \(\*OpenFileDialogStruct\) [SetOptions](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L315)

```go
func (d *OpenFileDialogStruct) SetOptions(options *OpenFileDialogOptions)
```

<a name="OpenFileDialogStruct.SetTitle"></a>

### func \(\*OpenFileDialogStruct\) [SetTitle](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L258)

```go
func (d *OpenFileDialogStruct) SetTitle(title string) *OpenFileDialogStruct
```

<a name="OpenFileDialogStruct.ShowHiddenFiles"></a>

### func \(\*OpenFileDialogStruct\) [ShowHiddenFiles](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L233)

```go
func (d *OpenFileDialogStruct) ShowHiddenFiles(showHiddenFiles bool) *OpenFileDialogStruct
```

<a name="OpenFileDialogStruct.TreatsFilePackagesAsDirectories"></a>

### func \(\*OpenFileDialogStruct\) [TreatsFilePackagesAsDirectories](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L243)

```go
func (d *OpenFileDialogStruct) TreatsFilePackagesAsDirectories(treatsFilePackagesAsDirectories bool) *OpenFileDialogStruct
```

<a name="Options"></a>

## type [Options](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/options_application.go#L9-L54)

```go
type Options struct {
    // Name is the name of the application
    Name string

    // Description is the description of the application (used in the default about box)
    Description string

    // Icon is the icon of the application (used in the default about box)
    Icon []byte

    // Mac is the Mac specific configuration for Mac builds
    Mac MacOptions

    // Windows is the Windows specific configuration for Windows builds
    Windows WindowsOptions

    // Bind allows you to bind Go methods to the frontend.
    Bind []any

    // BindAliases allows you to specify alias IDs for your bound methods.
    // Example: `BindAliases: map[uint32]uint32{1: 1411160069}` states that alias ID 1 maps to the Go method with ID 1411160069.
    BindAliases map[uint32]uint32

    // Logger i a slog.Logger instance used for logging Wails system messages (not application messages).
    // If not defined, a default logger is used.
    Logger *slog.Logger

    // LogLevel defines the log level of the Wails system logger.
    LogLevel slog.Level

    // Assets are the application assets to be used.
    Assets AssetOptions

    // Plugins is a map of plugins used by the application
    Plugins map[string]Plugin

    // Flags are key value pairs that are available to the frontend.
    // This is also used by Wails to provide information to the frontend.
    Flags map[string]any

    // PanicHandler is a way to register a custom panic handler
    PanicHandler func(any)

    // KeyBindings is a map of key bindings to functions
    KeyBindings map[string]func(window *WebviewWindow)
}
```

<a name="Parameter"></a>

## type [Parameter](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/bindings.go#L38-L42)

Parameter defines a Go method parameter

```go
type Parameter struct {
    Name        string `json:"name,omitempty"`
    TypeName    string `json:"type"`
    ReflectType reflect.Type
}
```

<a name="Parameter.IsError"></a>

### func \(\*Parameter\) [IsError](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/bindings.go#L58)

```go
func (p *Parameter) IsError() bool
```

IsError returns true if the parameter type is an error

<a name="Parameter.IsType"></a>

### func \(\*Parameter\) [IsType](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/bindings.go#L53)

```go
func (p *Parameter) IsType(typename string) bool
```

IsType returns true if the given

<a name="Plugin"></a>

## type [Plugin](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/plugins.go#L5-L11)

```go
type Plugin interface {
    Name() string
    Init() error
    Shutdown()
    CallableByJS() []string
    InjectJS() string
}
```

<a name="PluginCallOptions"></a>

## type [PluginCallOptions](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/bindings.go#L25-L28)

```go
type PluginCallOptions struct {
    Name string `json:"name"`
    Args []any  `json:"args"`
}
```

<a name="PluginManager"></a>

## type [PluginManager](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/plugins.go#L13-L17)

```go
type PluginManager struct {
    // contains filtered or unexported fields
}
```

<a name="NewPluginManager"></a>

### func [NewPluginManager](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/plugins.go#L19)

```go
func NewPluginManager(plugins map[string]Plugin, assetServer *assetserver.AssetServer) *PluginManager
```

<a name="PluginManager.Init"></a>

### func \(\*PluginManager\) [Init](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/plugins.go#L27)

```go
func (p *PluginManager) Init() error
```

<a name="PluginManager.Shutdown"></a>

### func \(\*PluginManager\) [Shutdown](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/plugins.go#L45)

```go
func (p *PluginManager) Shutdown()
```

<a name="PositionOptions"></a>

## type [PositionOptions](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L41-L43)

```go
type PositionOptions struct {
    Buffer int
}
```

<a name="QueryParams"></a>

## type [QueryParams](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor_params.go#L9)

```go
type QueryParams map[string][]string
```

<a name="QueryParams.Args"></a>

### func \(QueryParams\) [Args](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor_params.go#L189)

```go
func (qp QueryParams) Args() (*Args, error)
```

<a name="QueryParams.Bool"></a>

### func \(QueryParams\) [Bool](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor_params.go#L77)

```go
func (qp QueryParams) Bool(key string) *bool
```

<a name="QueryParams.Float64"></a>

### func \(QueryParams\) [Float64](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor_params.go#L89)

```go
func (qp QueryParams) Float64(key string) *float64
```

<a name="QueryParams.Int"></a>

### func \(QueryParams\) [Int](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor_params.go#L22)

```go
func (qp QueryParams) Int(key string) *int
```

<a name="QueryParams.String"></a>

### func \(QueryParams\) [String](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor_params.go#L11)

```go
func (qp QueryParams) String(key string) *string
```

<a name="QueryParams.ToStruct"></a>

### func \(QueryParams\) [ToStruct](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor_params.go#L101)

```go
func (qp QueryParams) ToStruct(str any) error
```

<a name="QueryParams.UInt"></a>

### func \(QueryParams\) [UInt](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor_params.go#L55)

```go
func (qp QueryParams) UInt(key string) *uint
```

<a name="QueryParams.UInt8"></a>

### func \(QueryParams\) [UInt8](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/messageprocessor_params.go#L34)

```go
func (qp QueryParams) UInt8(key string) *uint8
```

<a name="RGBA"></a>

## type [RGBA](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/options_webview_window.go#L133-L135)

```go
type RGBA struct {
    Red, Green, Blue, Alpha uint8
}
```

<a name="RadioGroup"></a>

## type [RadioGroup](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/popupmenu_windows.go#L17)

```go
type RadioGroup []*RadioGroupMember
```

<a name="RadioGroup.Add"></a>

### func \(\*RadioGroup\) [Add](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/popupmenu_windows.go#L19)

```go
func (r *RadioGroup) Add(id int, item *MenuItem)
```

<a name="RadioGroup.Bounds"></a>

### func \(\*RadioGroup\) [Bounds](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/popupmenu_windows.go#L26)

```go
func (r *RadioGroup) Bounds() (int, int)
```

<a name="RadioGroup.MenuID"></a>

### func \(\*RadioGroup\) [MenuID](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/popupmenu_windows.go#L31)

```go
func (r *RadioGroup) MenuID(item *MenuItem) int
```

<a name="RadioGroupMember"></a>

## type [RadioGroupMember](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/popupmenu_windows.go#L12-L15)

```go
type RadioGroupMember struct {
    ID       int
    MenuItem *MenuItem
}
```

<a name="Rect"></a>

## type [Rect](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/screen.go#L16-L21)

```go
type Rect struct {
    X      int
    Y      int
    Width  int
    Height int
}
```

<a name="Role"></a>

## type [Role](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/roles.go#L9)

Role is a type to identify menu roles

```go
type Role uint
```

<a name="NoRole"></a>These constants need to be kept in sync with
\`v2/internal/frontend/desktop/darwin/Role.h\`

```go
const (
    NoRole       Role = iota
    AppMenu      Role = iota
    EditMenu     Role = iota
    ViewMenu     Role = iota
    WindowMenu   Role = iota
    ServicesMenu Role = iota
    HelpMenu     Role = iota

    Hide               Role = iota
    HideOthers         Role = iota
    UnHide             Role = iota
    About              Role = iota
    Undo               Role = iota
    Redo               Role = iota
    Cut                Role = iota
    Copy               Role = iota
    Paste              Role = iota
    PasteAndMatchStyle Role = iota
    SelectAll          Role = iota
    Delete             Role = iota
    SpeechMenu         Role = iota
    Quit               Role = iota
    FileMenu           Role = iota
    Close              Role = iota
    Reload             Role = iota
    ForceReload        Role = iota
    ToggleDevTools     Role = iota
    ResetZoom          Role = iota
    ZoomIn             Role = iota
    ZoomOut            Role = iota
    ToggleFullscreen   Role = iota

    Minimize   Role = iota
    Zoom       Role = iota
    FullScreen Role = iota
)
```

<a name="SaveFileDialogOptions"></a>

## type [SaveFileDialogOptions](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L351-L365)

```go
type SaveFileDialogOptions struct {
    CanCreateDirectories            bool
    ShowHiddenFiles                 bool
    CanSelectHiddenExtension        bool
    AllowOtherFileTypes             bool
    HideExtension                   bool
    TreatsFilePackagesAsDirectories bool
    Title                           string
    Message                         string
    Directory                       string
    Filename                        string
    ButtonText                      string
    Filters                         []FileFilter
    Window                          *WebviewWindow
}
```

<a name="SaveFileDialogStruct"></a>

## type [SaveFileDialogStruct](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L367-L385)

```go
type SaveFileDialogStruct struct {
    // contains filtered or unexported fields
}
```

<a name="SaveFileDialog"></a>

### func [SaveFileDialog](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L627)

```go
func SaveFileDialog() *SaveFileDialogStruct
```

<a name="SaveFileDialogWithOptions"></a>

### func [SaveFileDialogWithOptions](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/application.go#L667)

```go
func SaveFileDialogWithOptions(s *SaveFileDialogOptions) *SaveFileDialogStruct
```

<a name="SaveFileDialogStruct.AddFilter"></a>

### func \(\*SaveFileDialogStruct\) [AddFilter](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L409)

```go
func (d *SaveFileDialogStruct) AddFilter(displayName, pattern string) *SaveFileDialogStruct
```

AddFilter adds a filter to the dialog. The filter is a display name and a
semicolon separated list of extensions. EG: AddFilter\("Image Files",
"\*.jpg;\*.png"\)

<a name="SaveFileDialogStruct.AllowsOtherFileTypes"></a>

### func \(\*SaveFileDialogStruct\) [AllowsOtherFileTypes](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L464)

```go
func (d *SaveFileDialogStruct) AllowsOtherFileTypes(allowOtherFileTypes bool) *SaveFileDialogStruct
```

<a name="SaveFileDialogStruct.AttachToWindow"></a>

### func \(\*SaveFileDialogStruct\) [AttachToWindow](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L442)

```go
func (d *SaveFileDialogStruct) AttachToWindow(window *WebviewWindow) *SaveFileDialogStruct
```

<a name="SaveFileDialogStruct.CanCreateDirectories"></a>

### func \(\*SaveFileDialogStruct\) [CanCreateDirectories](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L417)

```go
func (d *SaveFileDialogStruct) CanCreateDirectories(canCreateDirectories bool) *SaveFileDialogStruct
```

<a name="SaveFileDialogStruct.CanSelectHiddenExtension"></a>

### func \(\*SaveFileDialogStruct\) [CanSelectHiddenExtension](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L422)

```go
func (d *SaveFileDialogStruct) CanSelectHiddenExtension(canSelectHiddenExtension bool) *SaveFileDialogStruct
```

<a name="SaveFileDialogStruct.HideExtension"></a>

### func \(\*SaveFileDialogStruct\) [HideExtension](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L469)

```go
func (d *SaveFileDialogStruct) HideExtension(hideExtension bool) *SaveFileDialogStruct
```

<a name="SaveFileDialogStruct.PromptForSingleSelection"></a>

### func \(\*SaveFileDialogStruct\) [PromptForSingleSelection](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L447)

```go
func (d *SaveFileDialogStruct) PromptForSingleSelection() (string, error)
```

<a name="SaveFileDialogStruct.SetButtonText"></a>

### func \(\*SaveFileDialogStruct\) [SetButtonText](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L454)

```go
func (d *SaveFileDialogStruct) SetButtonText(text string) *SaveFileDialogStruct
```

<a name="SaveFileDialogStruct.SetDirectory"></a>

### func \(\*SaveFileDialogStruct\) [SetDirectory](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L437)

```go
func (d *SaveFileDialogStruct) SetDirectory(directory string) *SaveFileDialogStruct
```

<a name="SaveFileDialogStruct.SetFilename"></a>

### func \(\*SaveFileDialogStruct\) [SetFilename](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L459)

```go
func (d *SaveFileDialogStruct) SetFilename(filename string) *SaveFileDialogStruct
```

<a name="SaveFileDialogStruct.SetMessage"></a>

### func \(\*SaveFileDialogStruct\) [SetMessage](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L432)

```go
func (d *SaveFileDialogStruct) SetMessage(message string) *SaveFileDialogStruct
```

<a name="SaveFileDialogStruct.SetOptions"></a>

### func \(\*SaveFileDialogStruct\) [SetOptions](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L391)

```go
func (d *SaveFileDialogStruct) SetOptions(options *SaveFileDialogOptions)
```

<a name="SaveFileDialogStruct.ShowHiddenFiles"></a>

### func \(\*SaveFileDialogStruct\) [ShowHiddenFiles](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L427)

```go
func (d *SaveFileDialogStruct) ShowHiddenFiles(showHiddenFiles bool) *SaveFileDialogStruct
```

<a name="SaveFileDialogStruct.TreatsFilePackagesAsDirectories"></a>

### func \(\*SaveFileDialogStruct\) [TreatsFilePackagesAsDirectories](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/dialogs.go#L474)

```go
func (d *SaveFileDialogStruct) TreatsFilePackagesAsDirectories(treatsFilePackagesAsDirectories bool) *SaveFileDialogStruct
```

<a name="Screen"></a>

## type [Screen](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/screen.go#L3-L14)

```go
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
```

<a name="Size"></a>

## type [Size](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/screen.go#L23-L26)

```go
type Size struct {
    Width  int
    Height int
}
```

<a name="SystemTray"></a>

## type [SystemTray](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L45-L66)

```go
type SystemTray struct {
    // contains filtered or unexported fields
}
```

<a name="SystemTray.AttachWindow"></a>

### func \(\*SystemTray\) [AttachWindow](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L247)

```go
func (s *SystemTray) AttachWindow(window *WebviewWindow) *SystemTray
```

AttachWindow attaches a window to the system tray. The window will be shown when
the system tray icon is clicked. The window will be hidden when the system tray
icon is clicked again, or when the window loses focus.

<a name="SystemTray.Destroy"></a>

### func \(\*SystemTray\) [Destroy](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L187)

```go
func (s *SystemTray) Destroy()
```

<a name="SystemTray.Label"></a>

### func \(\*SystemTray\) [Label](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L93)

```go
func (s *SystemTray) Label() string
```

<a name="SystemTray.OnClick"></a>

### func \(\*SystemTray\) [OnClick](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L194)

```go
func (s *SystemTray) OnClick(handler func()) *SystemTray
```

<a name="SystemTray.OnDoubleClick"></a>

### func \(\*SystemTray\) [OnDoubleClick](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L204)

```go
func (s *SystemTray) OnDoubleClick(handler func()) *SystemTray
```

<a name="SystemTray.OnMouseEnter"></a>

### func \(\*SystemTray\) [OnMouseEnter](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L214)

```go
func (s *SystemTray) OnMouseEnter(handler func()) *SystemTray
```

<a name="SystemTray.OnMouseLeave"></a>

### func \(\*SystemTray\) [OnMouseLeave](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L219)

```go
func (s *SystemTray) OnMouseLeave(handler func()) *SystemTray
```

<a name="SystemTray.OnRightClick"></a>

### func \(\*SystemTray\) [OnRightClick](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L199)

```go
func (s *SystemTray) OnRightClick(handler func()) *SystemTray
```

<a name="SystemTray.OnRightDoubleClick"></a>

### func \(\*SystemTray\) [OnRightDoubleClick](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L209)

```go
func (s *SystemTray) OnRightDoubleClick(handler func()) *SystemTray
```

<a name="SystemTray.OpenMenu"></a>

### func \(\*SystemTray\) [OpenMenu](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L290)

```go
func (s *SystemTray) OpenMenu()
```

<a name="SystemTray.PositionWindow"></a>

### func \(\*SystemTray\) [PositionWindow](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L122)

```go
func (s *SystemTray) PositionWindow(window *WebviewWindow, offset int) error
```

<a name="SystemTray.SetDarkModeIcon"></a>

### func \(\*SystemTray\) [SetDarkModeIcon](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L142)

```go
func (s *SystemTray) SetDarkModeIcon(icon []byte) *SystemTray
```

<a name="SystemTray.SetIcon"></a>

### func \(\*SystemTray\) [SetIcon](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L131)

```go
func (s *SystemTray) SetIcon(icon []byte) *SystemTray
```

<a name="SystemTray.SetIconPosition"></a>

### func \(\*SystemTray\) [SetIconPosition](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L164)

```go
func (s *SystemTray) SetIconPosition(iconPosition int) *SystemTray
```

<a name="SystemTray.SetLabel"></a>

### func \(\*SystemTray\) [SetLabel](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L83)

```go
func (s *SystemTray) SetLabel(label string)
```

<a name="SystemTray.SetMenu"></a>

### func \(\*SystemTray\) [SetMenu](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L153)

```go
func (s *SystemTray) SetMenu(menu *Menu) *SystemTray
```

<a name="SystemTray.SetTemplateIcon"></a>

### func \(\*SystemTray\) [SetTemplateIcon](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L175)

```go
func (s *SystemTray) SetTemplateIcon(icon []byte) *SystemTray
```

<a name="SystemTray.WindowDebounce"></a>

### func \(\*SystemTray\) [WindowDebounce](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L262)

```go
func (s *SystemTray) WindowDebounce(debounce time.Duration) *SystemTray
```

WindowDebounce is used by Windows to indicate how long to wait before responding
to a mouse up event on the notification icon. This prevents the window from
being hidden and then immediately shown when the user clicks on the system tray
icon. See
https://stackoverflow.com/questions/4585283/alternate-showing-hiding-window-when-notify-icon-is-clicked

<a name="SystemTray.WindowOffset"></a>

### func \(\*SystemTray\) [WindowOffset](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L253)

```go
func (s *SystemTray) WindowOffset(offset int) *SystemTray
```

WindowOffset sets the gap in pixels between the system tray and the window

<a name="Theme"></a>

## type [Theme](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/options_webview_window_win.go#L79)

```go
type Theme int
```

<a name="SystemDefault"></a>

```go
const (
    // SystemDefault will use whatever the system theme is. The application will follow system theme changes.
    SystemDefault Theme = 0
    // Dark Mode
    Dark Theme = 1
    // Light Mode
    Light Theme = 2
)
```

<a name="ThemeSettings"></a>

## type [ThemeSettings](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/options_webview_window_win.go#L92-L105)

ThemeSettings defines custom colours to use in dark or light mode. They may be
set using the hex values: 0x00BBGGRR

```go
type ThemeSettings struct {
    DarkModeTitleBar           int32
    DarkModeTitleBarInactive   int32
    DarkModeTitleText          int32
    DarkModeTitleTextInactive  int32
    DarkModeBorder             int32
    DarkModeBorderInactive     int32
    LightModeTitleBar          int32
    LightModeTitleBarInactive  int32
    LightModeTitleText         int32
    LightModeTitleTextInactive int32
    LightModeBorder            int32
    LightModeBorderInactive    int32
}
```

<a name="WailsEvent"></a>

## type [WailsEvent](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/events.go#L42-L47)

```go
type WailsEvent struct {
    Name      string `json:"name"`
    Data      any    `json:"data"`
    Sender    string `json:"sender"`
    Cancelled bool
}
```

<a name="WailsEvent.Cancel"></a>

### func \(\*WailsEvent\) [Cancel](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/events.go#L49)

```go
func (e *WailsEvent) Cancel()
```

<a name="WebviewWindow"></a>

## type [WebviewWindow](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L97-L116)

```go
type WebviewWindow struct {
    // contains filtered or unexported fields
}
```

<a name="WebviewWindow.AbsolutePosition"></a>

### func \(\*WebviewWindow\) [AbsolutePosition](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L679)

```go
func (w *WebviewWindow) AbsolutePosition() (int, int)
```

AbsolutePosition returns the absolute position of the window to the screen

<a name="WebviewWindow.Center"></a>

### func \(\*WebviewWindow\) [Center](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L580)

```go
func (w *WebviewWindow) Center()
```

Center centers the window on the screen

<a name="WebviewWindow.Close"></a>

### func \(\*WebviewWindow\) [Close](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L770)

```go
func (w *WebviewWindow) Close()
```

Close closes the window

<a name="WebviewWindow.Destroy"></a>

### func \(\*WebviewWindow\) [Destroy](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L690)

```go
func (w *WebviewWindow) Destroy()
```

<a name="WebviewWindow.ExecJS"></a>

### func \(\*WebviewWindow\) [ExecJS](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L434)

```go
func (w *WebviewWindow) ExecJS(js string)
```

ExecJS executes the given javascript in the context of the window.

<a name="WebviewWindow.Flash"></a>

### func \(\*WebviewWindow\) [Flash](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L465)

```go
func (w *WebviewWindow) Flash(enabled bool)
```

Flash flashes the window's taskbar button/icon. Windows only.

<a name="WebviewWindow.Focus"></a>

### func \(\*WebviewWindow\) [Focus](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L1004)

```go
func (w *WebviewWindow) Focus()
```

<a name="WebviewWindow.ForceReload"></a>

### func \(\*WebviewWindow\) [ForceReload](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L712)

```go
func (w *WebviewWindow) ForceReload()
```

ForceReload forces the window to reload the page assets

<a name="WebviewWindow.Fullscreen"></a>

### func \(\*WebviewWindow\) [Fullscreen](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L442)

```go
func (w *WebviewWindow) Fullscreen() *WebviewWindow
```

Fullscreen sets the window to fullscreen mode. Min/Max size constraints are
disabled.

<a name="WebviewWindow.GetScreen"></a>

### func \(\*WebviewWindow\) [GetScreen](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L916)

```go
func (w *WebviewWindow) GetScreen() (*Screen, error)
```

GetScreen returns the screen that the window is on

<a name="WebviewWindow.GetZoom"></a>

### func \(\*WebviewWindow\) [GetZoom](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L346)

```go
func (w *WebviewWindow) GetZoom() float64
```

GetZoom returns the current zoom level of the window.

<a name="WebviewWindow.Height"></a>

### func \(\*WebviewWindow\) [Height](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L659)

```go
func (w *WebviewWindow) Height() int
```

Height returns the height of the window

<a name="WebviewWindow.Hide"></a>

### func \(\*WebviewWindow\) [Hide](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L315)

```go
func (w *WebviewWindow) Hide() *WebviewWindow
```

Hide hides the window.

<a name="WebviewWindow.IsFocused"></a>

### func \(\*WebviewWindow\) [IsFocused](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L511)

```go
func (w *WebviewWindow) IsFocused() bool
```

IsFocused returns true if the window is currently focused

<a name="WebviewWindow.IsFullscreen"></a>

### func \(\*WebviewWindow\) [IsFullscreen](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L519)

```go
func (w *WebviewWindow) IsFullscreen() bool
```

IsFullscreen returns true if the window is fullscreen

<a name="WebviewWindow.IsMaximised"></a>

### func \(\*WebviewWindow\) [IsMaximised](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L491)

```go
func (w *WebviewWindow) IsMaximised() bool
```

IsMaximised returns true if the window is maximised

<a name="WebviewWindow.IsMinimised"></a>

### func \(\*WebviewWindow\) [IsMinimised](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L475)

```go
func (w *WebviewWindow) IsMinimised() bool
```

IsMinimised returns true if the window is minimised

<a name="WebviewWindow.IsVisible"></a>

### func \(\*WebviewWindow\) [IsVisible](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L483)

```go
func (w *WebviewWindow) IsVisible() bool
```

IsVisible returns true if the window is visible

<a name="WebviewWindow.Maximise"></a>

### func \(\*WebviewWindow\) [Maximise](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L822)

```go
func (w *WebviewWindow) Maximise() *WebviewWindow
```

Maximise maximises the window. Min/Max size constraints are disabled.

<a name="WebviewWindow.Minimise"></a>

### func \(\*WebviewWindow\) [Minimise](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L809)

```go
func (w *WebviewWindow) Minimise() *WebviewWindow
```

Minimise minimises the window.

<a name="WebviewWindow.Name"></a>

### func \(\*WebviewWindow\) [Name](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L234)

```go
func (w *WebviewWindow) Name() string
```

Name returns the name of the window

<a name="WebviewWindow.NativeWindowHandle"></a>

### func \(\*WebviewWindow\) [NativeWindowHandle](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L997)

```go
func (w *WebviewWindow) NativeWindowHandle() (uintptr, error)
```

NativeWindowHandle returns the platform native window handle for the window.

<a name="WebviewWindow.On"></a>

### func \(\*WebviewWindow\) [On](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L589)

```go
func (w *WebviewWindow) On(eventType events.WindowEventType, callback func(event *WindowEvent)) func()
```

On registers a callback for the given window event

<a name="WebviewWindow.Print"></a>

### func \(\*WebviewWindow\) [Print](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L1027)

```go
func (w *WebviewWindow) Print() error
```

<a name="WebviewWindow.RegisterContextMenu"></a>

### func \(\*WebviewWindow\) [RegisterContextMenu](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L990)

```go
func (w *WebviewWindow) RegisterContextMenu(name string, menu *Menu)
```

RegisterContextMenu registers a context menu and assigns it the given name.

<a name="WebviewWindow.RegisterHook"></a>

### func \(\*WebviewWindow\) [RegisterHook](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L609)

```go
func (w *WebviewWindow) RegisterHook(eventType events.WindowEventType, callback func(event *WindowEvent)) func()
```

RegisterHook registers a hook for the given window event

<a name="WebviewWindow.RelativePosition"></a>

### func \(\*WebviewWindow\) [RelativePosition](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L667)

```go
func (w *WebviewWindow) RelativePosition() (int, int)
```

RelativePosition returns the relative position of the window to the screen

<a name="WebviewWindow.Reload"></a>

### func \(\*WebviewWindow\) [Reload](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L704)

```go
func (w *WebviewWindow) Reload()
```

Reload reloads the page assets

<a name="WebviewWindow.Resizable"></a>

### func \(\*WebviewWindow\) [Resizable](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L365)

```go
func (w *WebviewWindow) Resizable() bool
```

Resizable returns true if the window is resizable.

<a name="WebviewWindow.Restore"></a>

### func \(\*WebviewWindow\) [Restore](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L871)

```go
func (w *WebviewWindow) Restore()
```

Restore restores the window to its previous state if it was previously
minimised, maximised or fullscreen.

<a name="WebviewWindow.SetAbsolutePosition"></a>

### func \(\*WebviewWindow\) [SetAbsolutePosition](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L1043)

```go
func (w *WebviewWindow) SetAbsolutePosition(x int, y int)
```

<a name="WebviewWindow.SetAlwaysOnTop"></a>

### func \(\*WebviewWindow\) [SetAlwaysOnTop](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L290)

```go
func (w *WebviewWindow) SetAlwaysOnTop(b bool) *WebviewWindow
```

SetAlwaysOnTop sets the window to be always on top.

<a name="WebviewWindow.SetBackgroundColour"></a>

### func \(\*WebviewWindow\) [SetBackgroundColour](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L527)

```go
func (w *WebviewWindow) SetBackgroundColour(colour RGBA) *WebviewWindow
```

SetBackgroundColour sets the background colour of the window

<a name="WebviewWindow.SetEnabled"></a>

### func \(\*WebviewWindow\) [SetEnabled](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L1034)

```go
func (w *WebviewWindow) SetEnabled(enabled bool)
```

<a name="WebviewWindow.SetFrameless"></a>

### func \(\*WebviewWindow\) [SetFrameless](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L924)

```go
func (w *WebviewWindow) SetFrameless(frameless bool) *WebviewWindow
```

SetFrameless removes the window frame and title bar

<a name="WebviewWindow.SetFullscreenButtonEnabled"></a>

### func \(\*WebviewWindow\) [SetFullscreenButtonEnabled](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L454)

```go
func (w *WebviewWindow) SetFullscreenButtonEnabled(enabled bool) *WebviewWindow
```

<a name="WebviewWindow.SetHTML"></a>

### func \(\*WebviewWindow\) [SetHTML](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L786)

```go
func (w *WebviewWindow) SetHTML(html string) *WebviewWindow
```

SetHTML sets the HTML of the window to the given html string.

<a name="WebviewWindow.SetMaxSize"></a>

### func \(\*WebviewWindow\) [SetMaxSize](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L402)

```go
func (w *WebviewWindow) SetMaxSize(maxWidth, maxHeight int) *WebviewWindow
```

SetMaxSize sets the maximum size of the window.

<a name="WebviewWindow.SetMinSize"></a>

### func \(\*WebviewWindow\) [SetMinSize](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L370)

```go
func (w *WebviewWindow) SetMinSize(minWidth, minHeight int) *WebviewWindow
```

SetMinSize sets the minimum size of the window.

<a name="WebviewWindow.SetRelativePosition"></a>

### func \(\*WebviewWindow\) [SetRelativePosition](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L797)

```go
func (w *WebviewWindow) SetRelativePosition(x, y int) *WebviewWindow
```

SetRelativePosition sets the position of the window.

<a name="WebviewWindow.SetResizable"></a>

### func \(\*WebviewWindow\) [SetResizable](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L354)

```go
func (w *WebviewWindow) SetResizable(b bool) *WebviewWindow
```

SetResizable sets whether the window is resizable.

<a name="WebviewWindow.SetSize"></a>

### func \(\*WebviewWindow\) [SetSize](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L239)

```go
func (w *WebviewWindow) SetSize(width, height int) *WebviewWindow
```

SetSize sets the size of the window

<a name="WebviewWindow.SetTitle"></a>

### func \(\*WebviewWindow\) [SetTitle](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L223)

```go
func (w *WebviewWindow) SetTitle(title string) *WebviewWindow
```

SetTitle sets the title of the window

<a name="WebviewWindow.SetURL"></a>

### func \(\*WebviewWindow\) [SetURL](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L324)

```go
func (w *WebviewWindow) SetURL(s string) *WebviewWindow
```

<a name="WebviewWindow.SetZoom"></a>

### func \(\*WebviewWindow\) [SetZoom](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L335)

```go
func (w *WebviewWindow) SetZoom(magnification float64) *WebviewWindow
```

SetZoom sets the zoom level of the window.

<a name="WebviewWindow.Show"></a>

### func \(\*WebviewWindow\) [Show](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L301)

```go
func (w *WebviewWindow) Show() *WebviewWindow
```

Show shows the window.

<a name="WebviewWindow.Size"></a>

### func \(\*WebviewWindow\) [Size](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L499)

```go
func (w *WebviewWindow) Size() (int, int)
```

Size returns the size of the window

<a name="WebviewWindow.ToggleDevTools"></a>

### func \(\*WebviewWindow\) [ToggleDevTools](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L733)

```go
func (w *WebviewWindow) ToggleDevTools()
```

<a name="WebviewWindow.ToggleFullscreen"></a>

### func \(\*WebviewWindow\) [ToggleFullscreen](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L720)

```go
func (w *WebviewWindow) ToggleFullscreen()
```

ToggleFullscreen toggles the window between fullscreen and normal

<a name="WebviewWindow.ToggleMaximise"></a>

### func \(\*WebviewWindow\) [ToggleMaximise](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L720)

```go
func (w *WebviewWindow) ToggleMaximise()
```

ToggleMaximise toggles the window between maximised and normal

<a name="WebviewWindow.UnFullscreen"></a>

### func \(\*WebviewWindow\) [UnFullscreen](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L859)

```go
func (w *WebviewWindow) UnFullscreen()
```

UnFullscreen un\-fullscreens the window.

<a name="WebviewWindow.UnMaximise"></a>

### func \(\*WebviewWindow\) [UnMaximise](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L847)

```go
func (w *WebviewWindow) UnMaximise()
```

UnMaximise un\-maximises the window.

<a name="WebviewWindow.UnMinimise"></a>

### func \(\*WebviewWindow\) [UnMinimise](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L836)

```go
func (w *WebviewWindow) UnMinimise()
```

UnMinimise un\-minimises the window. Min/Max size constraints are re\-enabled.

<a name="WebviewWindow.Width"></a>

### func \(\*WebviewWindow\) [Width](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L651)

```go
func (w *WebviewWindow) Width() int
```

Width returns the width of the window

<a name="WebviewWindow.Zoom"></a>

### func \(\*WebviewWindow\) [Zoom](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L777)

```go
func (w *WebviewWindow) Zoom()
```

<a name="WebviewWindow.ZoomIn"></a>

### func \(\*WebviewWindow\) [ZoomIn](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L751)

```go
func (w *WebviewWindow) ZoomIn()
```

ZoomIn increases the zoom level of the webview content

<a name="WebviewWindow.ZoomOut"></a>

### func \(\*WebviewWindow\) [ZoomOut](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L761)

```go
func (w *WebviewWindow) ZoomOut()
```

ZoomOut decreases the zoom level of the webview content

<a name="WebviewWindow.ZoomReset"></a>

### func \(\*WebviewWindow\) [ZoomReset](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L741)

```go
func (w *WebviewWindow) ZoomReset() *WebviewWindow
```

ZoomReset resets the zoom level of the webview content to 100%

<a name="WebviewWindowOptions"></a>

## type [WebviewWindowOptions](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/options_webview_window.go#L12-L118)

```go
type WebviewWindowOptions struct {
    // Name is a unique identifier that can be given to a window.
    Name string

    // Title is the title of the window.
    Title string

    // Width is the starting width of the window.
    Width int

    // Height is the starting height of the window.
    Height int

    // AlwaysOnTop will make the window float above other windows.
    AlwaysOnTop bool

    // URL is the URL to load in the window.
    URL string

    // DisableResize will disable the ability to resize the window.
    DisableResize bool

    // Frameless will remove the window frame.
    Frameless bool

    // MinWidth is the minimum width of the window.
    MinWidth int

    // MinHeight is the minimum height of the window.
    MinHeight int

    // MaxWidth is the maximum width of the window.
    MaxWidth int

    // MaxHeight is the maximum height of the window.
    MaxHeight int

    // StartState indicates the state of the window when it is first shown.
    // Default: WindowStateNormal
    StartState WindowState

    // Centered will center the window on the screen.
    Centered bool

    // BackgroundType is the type of background to use for the window.
    // Default: BackgroundTypeSolid
    BackgroundType BackgroundType

    // BackgroundColour is the colour to use for the window background.
    BackgroundColour RGBA

    // HTML is the HTML to load in the window.
    HTML string

    // JS is the JavaScript to load in the window.
    JS  string

    // CSS is the CSS to load in the window.
    CSS string

    // X is the starting X position of the window.
    X   int

    // Y is the starting Y position of the window.
    Y   int

    // TransparentTitlebar will make the titlebar transparent.
    // TODO: Move to mac window options
    FullscreenButtonEnabled bool

    // Hidden will hide the window when it is first created.
    Hidden bool

    // Zoom is the zoom level of the window.
    Zoom float64

    // ZoomControlEnabled will enable the zoom control.
    ZoomControlEnabled bool

    // EnableDragAndDrop will enable drag and drop.
    EnableDragAndDrop bool

    // OpenInspectorOnStartup will open the inspector when the window is first shown.
    OpenInspectorOnStartup bool

    // Mac options
    Mac MacWindow

    // Windows options
    Windows WindowsWindow

    // Focused indicates the window should be focused when initially shown
    Focused bool

    // ShouldClose is called when the window is about to close.
    // Return true to allow the window to close, or false to prevent it from closing.
    ShouldClose func(window *WebviewWindow) bool

    // If true, the window's devtools will be available (default true in builds without the `production` build tag)
    DevToolsEnabled bool

    // If true, the window's default context menu will be disabled (default false)
    DefaultContextMenuDisabled bool

    // KeyBindings is a map of key bindings to functions
    KeyBindings map[string]func(window *WebviewWindow)
}
```

<a name="Win32Menu"></a>

## type [Win32Menu](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/popupmenu_windows.go#L40-L51)

```go
type Win32Menu struct {
    // contains filtered or unexported fields
}
```

<a name="NewApplicationMenu"></a>

### func [NewApplicationMenu](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/popupmenu_windows.go#L143)

```go
func NewApplicationMenu(parent w32.HWND, inputMenu *Menu) *Win32Menu
```

<a name="NewPopupMenu"></a>

### func [NewPopupMenu](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/popupmenu_windows.go#L132)

```go
func NewPopupMenu(parent w32.HWND, inputMenu *Menu) *Win32Menu
```

<a name="Win32Menu.Destroy"></a>

### func \(\*Win32Menu\) [Destroy](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/popupmenu_windows.go#L200)

```go
func (p *Win32Menu) Destroy()
```

<a name="Win32Menu.OnMenuClose"></a>

### func \(\*Win32Menu\) [OnMenuClose](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/popupmenu_windows.go#L241)

```go
func (p *Win32Menu) OnMenuClose(fn func())
```

<a name="Win32Menu.OnMenuOpen"></a>

### func \(\*Win32Menu\) [OnMenuOpen](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/popupmenu_windows.go#L237)

```go
func (p *Win32Menu) OnMenuOpen(fn func())
```

<a name="Win32Menu.ProcessCommand"></a>

### func \(\*Win32Menu\) [ProcessCommand](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/popupmenu_windows.go#L185)

```go
func (p *Win32Menu) ProcessCommand(cmdMsgID int) bool
```

<a name="Win32Menu.ShowAt"></a>

### func \(\*Win32Menu\) [ShowAt](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/popupmenu_windows.go#L154)

```go
func (p *Win32Menu) ShowAt(x int, y int)
```

<a name="Win32Menu.ShowAtCursor"></a>

### func \(\*Win32Menu\) [ShowAtCursor](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/popupmenu_windows.go#L176)

```go
func (p *Win32Menu) ShowAtCursor()
```

<a name="Win32Menu.Update"></a>

### func \(\*Win32Menu\) [Update](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/popupmenu_windows.go#L124)

```go
func (p *Win32Menu) Update()
```

<a name="Win32Menu.UpdateMenuItem"></a>

### func \(\*Win32Menu\) [UpdateMenuItem](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/popupmenu_windows.go#L204)

```go
func (p *Win32Menu) UpdateMenuItem(item *MenuItem)
```

<a name="WindowAttachConfig"></a>

## type [WindowAttachConfig](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/systemtray.go#L224-L243)

```go
type WindowAttachConfig struct {
    // Window is the window to attach to the system tray. If it's null, the request to attach will be ignored.
    Window *WebviewWindow

    // Offset indicates the gap in pixels between the system tray and the window
    Offset int

    // Debounce is used by Windows to indicate how long to wait before responding to a mouse
    // up event on the notification icon. See https://stackoverflow.com/questions/4585283/alternate-showing-hiding-window-when-notify-icon-is-clicked
    Debounce time.Duration
    // contains filtered or unexported fields
}
```

<a name="WindowEvent"></a>

## type [WindowEvent](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L76-L79)

```go
type WindowEvent struct {
    Cancelled bool
    // contains filtered or unexported fields
}
```

<a name="NewWindowEvent"></a>

### func [NewWindowEvent](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L85)

```go
func NewWindowEvent() *WindowEvent
```

<a name="WindowEvent.Cancel"></a>

### func \(\*WindowEvent\) [Cancel](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L89)

```go
func (w *WindowEvent) Cancel()
```

<a name="WindowEvent.Context"></a>

### func \(\*WindowEvent\) [Context](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L81)

```go
func (w *WindowEvent) Context() *WindowEventContext
```

<a name="WindowEventContext"></a>

## type [WindowEventContext](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/context_window_event.go#L10-L13)

```go
type WindowEventContext struct {
    // contains filtered or unexported fields
}
```

<a name="WindowEventContext.DroppedFiles"></a>

### func \(WindowEventContext\) [DroppedFiles](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/context_window_event.go#L15)

```go
func (c WindowEventContext) DroppedFiles() []string
```

<a name="WindowEventListener"></a>

## type [WindowEventListener](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window.go#L93-L95)

```go
type WindowEventListener struct {
    // contains filtered or unexported fields
}
```

<a name="WindowState"></a>

## type [WindowState](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/options_webview_window.go#L3)

```go
type WindowState int
```

<a name="WindowStateNormal"></a>

```go
const (
    WindowStateNormal WindowState = iota
    WindowStateMinimised
    WindowStateMaximised
    WindowStateFullscreen
)
```

<a name="WindowsOptions"></a>

## type [WindowsOptions](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/options_application_win.go#L4-L21)

WindowsOptions contains options for Windows applications.

```go
type WindowsOptions struct {

    // WndProcInterceptor is a function that will be called for every message sent in the application.
    // Use this to hook into the main message loop. This is useful for handling custom window messages.
    // If `shouldReturn` is `true` then `returnCode` will be returned by the main message loop.
    // If `shouldReturn` is `false` then returnCode will be ignored and the message will be processed by the main message loop.
    WndProcInterceptor func(hwnd uintptr, msg uint32, wParam, lParam uintptr) (returnCode uintptr, shouldReturn bool)

    // DisableQuitOnLastWindowClosed disables the auto quit of the application if the last window has been closed.
    DisableQuitOnLastWindowClosed bool

    // Path where the WebView2 stores the user data. If empty %APPDATA%\[BinaryName.exe] will be used.
    // If the path is not valid, a messagebox will be displayed with the error and the app will exit with error code.
    WebviewUserDataPath string

    // Path to the directory with WebView2 executables. If empty WebView2 installed in the system will be used.
    WebviewBrowserPath string
}
```

<a name="WindowsWindow"></a>

## type [WindowsWindow](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/options_webview_window_win.go#L15-L77)

```go
type WindowsWindow struct {
    // Select the type of translucent backdrop. Requires Windows 11 22621 or later.
    // Only used when window's `BackgroundType` is set to `BackgroundTypeTranslucent`.
    // Default: Auto
    BackdropType BackdropType

    // Disable the icon in the titlebar
    // Default: false
    DisableIcon bool

    // Theme (Dark / Light / SystemDefault)
    // Default: SystemDefault - The application will follow system theme changes.
    Theme Theme

    // Specify custom colours to use for dark/light mode
    // Default: nil
    CustomTheme *ThemeSettings

    // Disable all window decorations in Frameless mode, which means no "Aero Shadow" and no "Rounded Corner" will be shown.
    // "Rounded Corners" are only available on Windows 11.
    // Default: false
    DisableFramelessWindowDecorations bool

    // WindowMask is used to set the window shape. Use a PNG with an alpha channel to create a custom shape.
    // Default: nil
    WindowMask []byte

    // WindowMaskDraggable is used to make the window draggable by clicking on the window mask.
    // Default: false
    WindowMaskDraggable bool

    // WebviewGpuIsDisabled is used to enable / disable GPU acceleration for the webview
    // Default: false
    WebviewGpuIsDisabled bool

    // ResizeDebounceMS is the amount of time to debounce redraws of webview2
    // when resizing the window
    // Default: 0
    ResizeDebounceMS uint16

    // Disable the menu bar for this window
    // Default: false
    DisableMenu bool

    // Event mapping for the window. This allows you to define a translation from one event to another.
    // Default: nil
    EventMapping map[events.WindowEventType]events.WindowEventType

    // HiddenOnTaskbar hides the window from the taskbar
    // Default: false
    HiddenOnTaskbar bool

    // EnableSwipeGestures enables swipe gestures for the window
    // Default: false
    EnableSwipeGestures bool

    // EnableFraudulentWebsiteWarnings will enable warnings for fraudulent websites.
    // Default: false
    EnableFraudulentWebsiteWarnings bool

    // Menu is the menu to use for the window.
    Menu *Menu
}
```
