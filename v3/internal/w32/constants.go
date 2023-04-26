package w32

// TODO Check that these messages are correct:
// WM_*
// TTM_*
// TBM_*

const (
	FALSE = 0
	TRUE  = 1
)

const (
	NO_ERROR                         = 0
	ERROR_SUCCESS                    = 0
	ERROR_FILE_NOT_FOUND             = 2
	ERROR_PATH_NOT_FOUND             = 3
	ERROR_ACCESS_DENIED              = 5
	ERROR_INVALID_HANDLE             = 6
	ERROR_BAD_FORMAT                 = 11
	ERROR_INVALID_NAME               = 123
	ERROR_ALREADY_EXISTS             = 183
	ERROR_MORE_DATA                  = 234
	ERROR_NO_MORE_ITEMS              = 259
	ERROR_INVALID_SERVICE_CONTROL    = 1052
	ERROR_SERVICE_REQUEST_TIMEOUT    = 1053
	ERROR_SERVICE_NO_THREAD          = 1054
	ERROR_SERVICE_DATABASE_LOCKED    = 1055
	ERROR_SERVICE_ALREADY_RUNNING    = 1056
	ERROR_SERVICE_DISABLED           = 1058
	ERROR_SERVICE_DOES_NOT_EXIST     = 1060
	ERROR_SERVICE_CANNOT_ACCEPT_CTRL = 1061
	ERROR_SERVICE_NOT_ACTIVE         = 1062
	ERROR_DATABASE_DOES_NOT_EXIST    = 1065
	ERROR_SERVICE_DEPENDENCY_FAIL    = 1068
	ERROR_SERVICE_LOGON_FAILED       = 1069
	ERROR_SERVICE_MARKED_FOR_DELETE  = 1072
	ERROR_SERVICE_DEPENDENCY_DELETED = 1075
)

const (
	SE_ERR_FNF             = 2
	SE_ERR_PNF             = 3
	SE_ERR_ACCESSDENIED    = 5
	SE_ERR_OOM             = 8
	SE_ERR_DLLNOTFOUND     = 32
	SE_ERR_SHARE           = 26
	SE_ERR_ASSOCINCOMPLETE = 27
	SE_ERR_DDETIMEOUT      = 28
	SE_ERR_DDEFAIL         = 29
	SE_ERR_DDEBUSY         = 30
	SE_ERR_NOASSOC         = 31
)

const (
	CW_USEDEFAULT = ^0x7FFFFFFF
)

// ShowWindow constants
const (
	SW_HIDE            = 0
	SW_NORMAL          = 1
	SW_SHOWNORMAL      = 1
	SW_SHOWMINIMIZED   = 2
	SW_MAXIMIZE        = 3
	SW_SHOWMAXIMIZED   = 3
	SW_SHOWNOACTIVATE  = 4
	SW_SHOW            = 5
	SW_MINIMIZE        = 6
	SW_SHOWMINNOACTIVE = 7
	SW_SHOWNA          = 8
	SW_RESTORE         = 9
	SW_SHOWDEFAULT     = 10
	SW_FORCEMINIMIZE   = 11
)

// Window class styles
const (
	CS_VREDRAW         = 0x00000001
	CS_HREDRAW         = 0x00000002
	CS_KEYCVTWINDOW    = 0x00000004
	CS_DBLCLKS         = 0x00000008
	CS_OWNDC           = 0x00000020
	CS_CLASSDC         = 0x00000040
	CS_PARENTDC        = 0x00000080
	CS_NOKEYCVT        = 0x00000100
	CS_NOCLOSE         = 0x00000200
	CS_SAVEBITS        = 0x00000800
	CS_BYTEALIGNCLIENT = 0x00001000
	CS_BYTEALIGNWINDOW = 0x00002000
	CS_GLOBALCLASS     = 0x00004000
	CS_IME             = 0x00010000
	CS_DROPSHADOW      = 0x00020000
)

// Predefined cursor constants
const (
	IDC_ARROW       = 32512
	IDC_IBEAM       = 32513
	IDC_WAIT        = 32514
	IDC_CROSS       = 32515
	IDC_UPARROW     = 32516
	IDC_SIZENWSE    = 32642
	IDC_SIZENESW    = 32643
	IDC_SIZEWE      = 32644
	IDC_SIZENS      = 32645
	IDC_SIZEALL     = 32646
	IDC_NO          = 32648
	IDC_HAND        = 32649
	IDC_APPSTARTING = 32650
	IDC_HELP        = 32651
	IDC_ICON        = 32641
	IDC_SIZE        = 32640
)

// Predefined icon constants
const (
	IDI_APPLICATION = 32512
	IDI_HAND        = 32513
	IDI_QUESTION    = 32514
	IDI_EXCLAMATION = 32515
	IDI_ASTERISK    = 32516
	IDI_WINLOGO     = 32517
	IDI_SHIELD      = 32518
	IDI_WARNING     = IDI_EXCLAMATION
	IDI_ERROR       = IDI_HAND
	IDI_INFORMATION = IDI_ASTERISK
)

// Button style constants
const (
	BS_3STATE          = 5
	BS_AUTO3STATE      = 6
	BS_AUTOCHECKBOX    = 3
	BS_AUTORADIOBUTTON = 9
	BS_BITMAP          = 128
	BS_BOTTOM          = 0x800
	BS_CENTER          = 0x300
	BS_CHECKBOX        = 2
	BS_DEFPUSHBUTTON   = 1
	BS_GROUPBOX        = 7
	BS_ICON            = 64
	BS_LEFT            = 256
	BS_LEFTTEXT        = 32
	BS_MULTILINE       = 0x2000
	BS_NOTIFY          = 0x4000
	BS_OWNERDRAW       = 0xB
	BS_PUSHBUTTON      = 0
	BS_PUSHLIKE        = 4096
	BS_RADIOBUTTON     = 4
	BS_RIGHT           = 512
	BS_RIGHTBUTTON     = 32
	BS_TEXT            = 0
	BS_TOP             = 0x400
	BS_USERBUTTON      = 8
	BS_VCENTER         = 0xC00
	BS_FLAT            = 0x8000
)

// Button state constants
const (
	BST_CHECKED       = 1
	BST_INDETERMINATE = 2
	BST_UNCHECKED     = 0
	BST_FOCUS         = 8
	BST_PUSHED        = 4
)

// Combo box style constants
const (
	CBS_SIMPLE            = 0x0001
	CBS_DROPDOWN          = 0x0002
	CBS_DROPDOWNLIST      = 0x0003
	CBS_OWNERDRAWFIXED    = 0x0010
	CBS_OWNERDRAWVARIABLE = 0x0020
	CBS_AUTOHSCROLL       = 0x0040
	CBS_OEMCONVERT        = 0x0080
	CBS_SORT              = 0x0100
	CBS_HASSTRINGS        = 0x0200
	CBS_NOINTEGRALHEIGHT  = 0x0400
	CBS_DISABLENOSCROLL   = 0x0800
	CBS_UPPERCASE         = 0x2000
	CBS_LOWERCASE         = 0x4000
)

// Combo box message constants
const (
	CB_GETEDITSEL            = 0x0140
	CB_LIMITTEXT             = 0x0141
	CB_SETEDITSEL            = 0x0142
	CB_ADDSTRING             = 0x0143
	CB_DELETESTRING          = 0x0144
	CB_DIR                   = 0x0145
	CB_GETCOUNT              = 0x0146
	CB_GETCURSEL             = 0x0147
	CB_GETLBTEXT             = 0x0148
	CB_GETLBTEXTLEN          = 0x0149
	CB_INSERTSTRING          = 0x014A
	CB_RESETCONTENT          = 0x014B
	CB_FINDSTRING            = 0x014C
	CB_SELECTSTRING          = 0x014D
	CB_SETCURSEL             = 0x014E
	CB_SHOWDROPDOWN          = 0x014F
	CB_GETITEMDATA           = 0x0150
	CB_SETITEMDATA           = 0x0151
	CB_GETDROPPEDCONTROLRECT = 0x0152
	CB_SETITEMHEIGHT         = 0x0153
	CB_GETITEMHEIGHT         = 0x0154
	CB_SETEXTENDEDUI         = 0x0155
	CB_GETEXTENDEDUI         = 0x0156
	CB_GETDROPPEDSTATE       = 0x0157
	CB_FINDSTRINGEXACT       = 0x0158
	CB_SETLOCALE             = 0x0159
	CB_GETLOCALE             = 0x015A
	CB_GETTOPINDEX           = 0x015B
	CB_SETTOPINDEX           = 0x015C
	CB_GETHORIZONTALEXTENT   = 0x015D
	CB_SETHORIZONTALEXTENT   = 0x015E
	CB_GETDROPPEDWIDTH       = 0x015F
	CB_SETDROPPEDWIDTH       = 0x0160
	CB_INITSTORAGE           = 0x0161
	CB_MULTIPLEADDSTRING     = 0x0163
	CB_GETCOMBOBOXINFO       = 0x0164
	CB_MSGMAX                = 0x0165
)

// Combo box return values
const (
	CB_OKAY     = 0
	CB_ERR      = -1
	CB_ERRSPACE = -2
)

// Combo box notification codes
const (
	CBN_ERRSPACE     = -1
	CBN_SELCHANGE    = 1
	CBN_DBLCLK       = 2
	CBN_SETFOCUS     = 3
	CBN_KILLFOCUS    = 4
	CBN_EDITCHANGE   = 5
	CBN_EDITUPDATE   = 6
	CBN_DROPDOWN     = 7
	CBN_CLOSEUP      = 8
	CBN_SELENDOK     = 9
	CBN_SELENDCANCEL = 10
)

// List box message constants
const (
	LB_ADDSTRING           = 384
	LB_INSERTSTRING        = 385
	LB_DELETESTRING        = 386
	LB_SELITEMRANGEEX      = 387
	LB_RESETCONTENT        = 388
	LB_SETSEL              = 389
	LB_SETCURSEL           = 390
	LB_GETSEL              = 391
	LB_GETCURSEL           = 392
	LB_GETTEXT             = 393
	LB_GETTEXTLEN          = 394
	LB_GETCOUNT            = 395
	LB_SELECTSTRING        = 396
	LB_DIR                 = 397
	LB_GETTOPINDEX         = 398
	LB_FINDSTRING          = 399
	LB_GETSELCOUNT         = 400
	LB_GETSELITEMS         = 401
	LB_SETTABSTOPS         = 402
	LB_GETHORIZONTALEXTENT = 403
	LB_SETHORIZONTALEXTENT = 404
	LB_SETCOLUMNWIDTH      = 405
	LB_ADDFILE             = 406
	LB_SETTOPINDEX         = 407
	LB_GETITEMRECT         = 408
	LB_GETITEMDATA         = 409
	LB_SETITEMDATA         = 410
	LB_SELITEMRANGE        = 411
	LB_SETANCHORINDEX      = 412
	LB_GETANCHORINDEX      = 413
	LB_SETCARETINDEX       = 414
	LB_GETCARETINDEX       = 415
	LB_SETITEMHEIGHT       = 416
	LB_GETITEMHEIGHT       = 417
	LB_FINDSTRINGEXACT     = 418
	LB_SETLOCALE           = 421
	LB_GETLOCALE           = 422
	LB_SETCOUNT            = 423
	LB_INITSTORAGE         = 424
	LB_ITEMFROMPOINT       = 425
	LB_SETTEXT             = 426
	LB_GETCHECKMARK        = 427
	LB_SETCHECKMARK        = 428
	LB_GETITEMADDDATA      = 429
	LB_SETITEMADDDATA      = 430
)

// List box styles
const (
	LBS_NOTIFY            = 0x0001
	LBS_SORT              = 0x0002
	LBS_NOREDRAW          = 0x0004
	LBS_MULTIPLESEL       = 0x0008
	LBS_OWNERDRAWFIXED    = 0x0010
	LBS_OWNERDRAWVARIABLE = 0x0020
	LBS_HASSTRINGS        = 0x0040
	LBS_USETABSTOPS       = 0x0080
	LBS_NOINTEGRALHEIGHT  = 0x0100
	LBS_MULTICOLUMN       = 0x0200
	LBS_WANTKEYBOARDINPUT = 0x0400
	LBS_EXTENDEDSEL       = 0x0800
	LBS_STANDARD          = LBS_NOTIFY | LBS_SORT | WS_VSCROLL | WS_BORDER
	LBS_CHECKBOX          = 0x1000
	LBS_USEICON           = 0x2000
	LBS_AUTOCHECK         = 0x4000
	LBS_AUTOCHECKBOX      = 0x5000
	LBS_PRELOADED         = 0x4000
	LBS_COMBOLBOX         = 0x8000
)

// List box notification messages
const (
	LBN_ERRSPACE       = -2
	LBN_SELCHANGE      = 1
	LBN_DBLCLK         = 2
	LBN_SELCANCEL      = 3
	LBN_SETFOCUS       = 4
	LBN_KILLFOCUS      = 5
	LBN_CLICKCHECKMARK = 6
)

// List box return values
const (
	LB_OKAY     = 0
	LB_ERR      = -1
	LB_ERRSPACE = -2
)

// Predefined color/brush constants.
const (
	// Scroll bar gray area.
	COLOR_SCROLLBAR = 0

	// Desktop.
	COLOR_BACKGROUND = 1

	// Desktop.
	COLOR_DESKTOP = 1

	// Active window title bar. The associated foreground color is
	// COLOR_CAPTIONTEXT. Specifies the left side color in the color gradient of
	// an active window's title bar if the gradient effect is enabled.
	COLOR_ACTIVECAPTION = 2

	// Inactive window caption.
	// The associated foreground color is COLOR_INACTIVECAPTIONTEXT.
	// Specifies the left side color in the color gradient of an inactive window's title bar if the gradient effect is enabled.
	COLOR_INACTIVECAPTION = 3

	// Menu background. The associated foreground color is COLOR_MENUTEXT.
	COLOR_MENU = 4

	// Window background. The associated foreground colors are COLOR_WINDOWTEXT
	// and COLOR_HOTLITE.
	COLOR_WINDOW = 5

	// Window frame.
	COLOR_WINDOWFRAME = 6

	// Text in menus. The associated background color is COLOR_MENU.
	COLOR_MENUTEXT = 7

	// Text in windows. The associated background color is COLOR_WINDOW.
	COLOR_WINDOWTEXT = 8

	// Text in caption, size box, and scroll bar arrow box. The associated
	// background color is COLOR_ACTIVECAPTION.
	COLOR_CAPTIONTEXT = 9

	// Active window border.
	COLOR_ACTIVEBORDER = 10

	// Inactive window border.
	COLOR_INACTIVEBORDER = 11

	// Background color of multiple document interface (MDI) applications.
	COLOR_APPWORKSPACE = 12

	// Item(s) selected in a control. The associated foreground color is
	// COLOR_HIGHLIGHTTEXT.
	COLOR_HIGHLIGHT = 13

	// Text of item(s) selected in a control. The associated background color is
	// COLOR_HIGHLIGHT.
	COLOR_HIGHLIGHTTEXT = 14

	// Face color for three-dimensional display elements and for dialog box
	// backgrounds.
	COLOR_3DFACE = 15

	// Face color for three-dimensional display elements and for dialog box
	// backgrounds. The associated foreground color is COLOR_BTNTEXT.
	COLOR_BTNFACE = 15

	// Shadow color for three-dimensional display elements (for edges facing
	// away from the light source).
	COLOR_3DSHADOW = 16

	// Shadow color for three-dimensional display elements (for edges facing
	// away from the light source).
	COLOR_BTNSHADOW = 16

	// Grayed (disabled) text. This color is set to 0 if the current display
	// driver does not support a solid gray color.
	COLOR_GRAYTEXT = 17

	// Text on push buttons. The associated background color is COLOR_BTNFACE.
	COLOR_BTNTEXT = 18

	// Color of text in an inactive caption. The associated background color is
	// COLOR_INACTIVECAPTION.
	COLOR_INACTIVECAPTIONTEXT = 19

	// Highlight color for three-dimensional display elements (for edges facing
	// the light source.)
	COLOR_3DHIGHLIGHT = 20

	// Highlight color for three-dimensional display elements (for edges facing
	// the light source.)
	COLOR_3DHILIGHT = 20

	// Highlight color for three-dimensional display elements (for edges facing
	// the light source.)
	COLOR_BTNHIGHLIGHT = 20

	// Highlight color for three-dimensional display elements (for edges facing
	// the light source.)
	COLOR_BTNHILIGHT = 20

	// Dark shadow for three-dimensional display elements.
	COLOR_3DDKSHADOW = 21

	// Light color for three-dimensional display elements (for edges facing the
	// light source.)
	COLOR_3DLIGHT = 22

	// Text color for tooltip controls. The associated background color is
	// COLOR_INFOBK.
	COLOR_INFOTEXT = 23

	// Background color for tooltip controls. The associated foreground color is
	// COLOR_INFOTEXT.
	COLOR_INFOBK = 24

	// Color for a hyperlink or hot-tracked item. The associated background
	// color is COLOR_WINDOW.
	COLOR_HOTLIGHT = 26

	// Right side color in the color gradient of an active window's title bar.
	// COLOR_ACTIVECAPTION specifies the left side color. Use
	// SPI_GETGRADIENTCAPTIONS with the SystemParametersInfo function to
	// determine whether the gradient effect is enabled.
	COLOR_GRADIENTACTIVECAPTION = 27

	// Right side color in the color gradient of an inactive window's title bar.
	// COLOR_INACTIVECAPTION specifies the left side color.
	COLOR_GRADIENTINACTIVECAPTION = 28

	// The color used to highlight menu items when the menu appears as a flat
	// menu (see SystemParametersInfo). The highlighted menu item is outlined
	// with COLOR_HIGHLIGHT. Windows 2000: This value is not supported.
	COLOR_MENUHILIGHT = 29

	// The background color for the menu bar when menus appear as flat menus
	// (see SystemParametersInfo). However, COLOR_MENU continues to specify the
	// background color of the menu popup. Windows 2000: This value is not
	// supported.
	COLOR_MENUBAR = 30
)

// Button message constants
const (
	BM_CLICK    = 245
	BM_GETCHECK = 240
	BM_GETIMAGE = 246
	BM_GETSTATE = 242
	BM_SETCHECK = 241
	BM_SETIMAGE = 247
	BM_SETSTATE = 243
	BM_SETSTYLE = 244
)

// Button notifications
const (
	BN_CLICKED       = 0
	BN_PAINT         = 1
	BN_HILITE        = 2
	BN_PUSHED        = BN_HILITE
	BN_UNHILITE      = 3
	BN_UNPUSHED      = BN_UNHILITE
	BN_DISABLE       = 4
	BN_DOUBLECLICKED = 5
	BN_DBLCLK        = BN_DOUBLECLICKED
	BN_SETFOCUS      = 6
	BN_KILLFOCUS     = 7
)

// GetWindowLong and GetWindowLongPtr constants
const (
	GWL_EXSTYLE     = -20
	GWL_STYLE       = -16
	GWL_WNDPROC     = -4
	GWLP_WNDPROC    = -4
	GWL_HINSTANCE   = -6
	GWLP_HINSTANCE  = -6
	GWL_HWNDPARENT  = -8
	GWLP_HWNDPARENT = -8
	GWL_ID          = -12
	GWLP_ID         = -12
	GWL_USERDATA    = -21
	GWLP_USERDATA   = -21
)

// Window style constants
const (
	WS_OVERLAPPED       = 0x00000000
	WS_POPUP            = 0x80000000
	WS_CHILD            = 0x40000000
	WS_MINIMIZE         = 0x20000000
	WS_VISIBLE          = 0x10000000
	WS_DISABLED         = 0x08000000
	WS_CLIPSIBLINGS     = 0x04000000
	WS_CLIPCHILDREN     = 0x02000000
	WS_MAXIMIZE         = 0x01000000
	WS_CAPTION          = 0x00C00000
	WS_BORDER           = 0x00800000
	WS_DLGFRAME         = 0x00400000
	WS_VSCROLL          = 0x00200000
	WS_HSCROLL          = 0x00100000
	WS_SYSMENU          = 0x00080000
	WS_THICKFRAME       = 0x00040000
	WS_GROUP            = 0x00020000
	WS_TABSTOP          = 0x00010000
	WS_MINIMIZEBOX      = 0x00020000
	WS_MAXIMIZEBOX      = 0x00010000
	WS_TILED            = 0x00000000
	WS_ICONIC           = 0x20000000
	WS_SIZEBOX          = 0x00040000
	WS_OVERLAPPEDWINDOW = 0x00000000 | 0x00C00000 | 0x00080000 | 0x00040000 | 0x00020000 | 0x00010000
	WS_POPUPWINDOW      = 0x80000000 | 0x00800000 | 0x00080000
	WS_CHILDWINDOW      = 0x40000000
)

// Extended window style constants
const (
	WS_EX_DLGMODALFRAME    = 0x00000001
	WS_EX_NOPARENTNOTIFY   = 0x00000004
	WS_EX_TOPMOST          = 0x00000008
	WS_EX_ACCEPTFILES      = 0x00000010
	WS_EX_TRANSPARENT      = 0x00000020
	WS_EX_MDICHILD         = 0x00000040
	WS_EX_TOOLWINDOW       = 0x00000080
	WS_EX_WINDOWEDGE       = 0x00000100
	WS_EX_CLIENTEDGE       = 0x00000200
	WS_EX_CONTEXTHELP      = 0x00000400
	WS_EX_RIGHT            = 0x00001000
	WS_EX_LEFT             = 0x00000000
	WS_EX_RTLREADING       = 0x00002000
	WS_EX_LTRREADING       = 0x00000000
	WS_EX_LEFTSCROLLBAR    = 0x00004000
	WS_EX_RIGHTSCROLLBAR   = 0x00000000
	WS_EX_CONTROLPARENT    = 0x00010000
	WS_EX_STATICEDGE       = 0x00020000
	WS_EX_APPWINDOW        = 0x00040000
	WS_EX_OVERLAPPEDWINDOW = 0x00000100 | 0x00000200
	WS_EX_PALETTEWINDOW    = 0x00000100 | 0x00000080 | 0x00000008
	WS_EX_LAYERED          = 0x00080000
	WS_EX_NOINHERITLAYOUT  = 0x00100000
	WS_EX_LAYOUTRTL        = 0x00400000
	WS_EX_COMPOSITED       = 0x02000000
	WS_EX_NOACTIVATE       = 0x08000000
)

// Window message constants
const (
	WM_APP                    = 32768
	WM_ACTIVATE               = 6
	WM_ACTIVATEAPP            = 28
	WM_AFXFIRST               = 864
	WM_AFXLAST                = 895
	WM_ASKCBFORMATNAME        = 780
	WM_CANCELJOURNAL          = 75
	WM_CANCELMODE             = 31
	WM_CAPTURECHANGED         = 533
	WM_CHANGECBCHAIN          = 781
	WM_CHAR                   = 258
	WM_CHARTOITEM             = 47
	WM_CHILDACTIVATE          = 34
	WM_CLEAR                  = 771
	WM_CLOSE                  = 16
	WM_COMMAND                = 273
	WM_COMMNOTIFY             = 68 // obsolete
	WM_COMPACTING             = 65
	WM_COMPAREITEM            = 57
	WM_CONTEXTMENU            = 123
	WM_COPY                   = 769
	WM_COPYDATA               = 74
	WM_CREATE                 = 1
	WM_CTLCOLORBTN            = 309
	WM_CTLCOLORDLG            = 310
	WM_CTLCOLOREDIT           = 307
	WM_CTLCOLORLISTBOX        = 308
	WM_CTLCOLORMSGBOX         = 306
	WM_CTLCOLORSCROLLBAR      = 311
	WM_CTLCOLORSTATIC         = 312
	WM_CUT                    = 768
	WM_DEADCHAR               = 259
	WM_DELETEITEM             = 45
	WM_DESTROY                = 2
	WM_DESTROYCLIPBOARD       = 775
	WM_DEVICECHANGE           = 537
	WM_DEVMODECHANGE          = 27
	WM_DISPLAYCHANGE          = 126
	WM_DRAWCLIPBOARD          = 776
	WM_DRAWITEM               = 43
	WM_DROPFILES              = 563
	WM_ENABLE                 = 10
	WM_ENDSESSION             = 22
	WM_ENTERIDLE              = 289
	WM_ENTERMENULOOP          = 529
	WM_ENTERSIZEMOVE          = 561
	WM_ERASEBKGND             = 20
	WM_EXITMENULOOP           = 530
	WM_EXITSIZEMOVE           = 562
	WM_FONTCHANGE             = 29
	WM_GETDLGCODE             = 135
	WM_GETFONT                = 49
	WM_GETHOTKEY              = 51
	WM_GETICON                = 127
	WM_GETMINMAXINFO          = 36
	WM_GETTEXT                = 13
	WM_GETTEXTLENGTH          = 14
	WM_HANDHELDFIRST          = 856
	WM_HANDHELDLAST           = 863
	WM_HELP                   = 83
	WM_HOTKEY                 = 786
	WM_HSCROLL                = 276
	WM_HSCROLLCLIPBOARD       = 782
	WM_ICONERASEBKGND         = 39
	WM_INITDIALOG             = 272
	WM_INITMENU               = 278
	WM_INITMENUPOPUP          = 279
	WM_INPUT                  = 0x00FF
	WM_INPUTLANGCHANGE        = 81
	WM_INPUTLANGCHANGEREQUEST = 80
	WM_KEYDOWN                = 256
	WM_KEYUP                  = 257
	WM_KILLFOCUS              = 8
	WM_MDIACTIVATE            = 546
	WM_MDICASCADE             = 551
	WM_MDICREATE              = 544
	WM_MDIDESTROY             = 545
	WM_MDIGETACTIVE           = 553
	WM_MDIICONARRANGE         = 552
	WM_MDIMAXIMIZE            = 549
	WM_MDINEXT                = 548
	WM_MDIREFRESHMENU         = 564
	WM_MDIRESTORE             = 547
	WM_MDISETMENU             = 560
	WM_MDITILE                = 550
	WM_MEASUREITEM            = 44
	WM_GETOBJECT              = 0x003D
	WM_CHANGEUISTATE          = 0x0127
	WM_UPDATEUISTATE          = 0x0128
	WM_QUERYUISTATE           = 0x0129
	WM_UNINITMENUPOPUP        = 0x0125
	WM_MENURBUTTONUP          = 290
	WM_MENUCOMMAND            = 0x0126
	WM_MENUGETOBJECT          = 0x0124
	WM_MENUDRAG               = 0x0123
	WM_APPCOMMAND             = 0x0319
	WM_MENUCHAR               = 288
	WM_MENUSELECT             = 287
	WM_MOVE                   = 3
	WM_MOVING                 = 534
	WM_NCACTIVATE             = 134
	WM_NCCALCSIZE             = 131
	WM_NCCREATE               = 129
	WM_NCDESTROY              = 130
	WM_NCHITTEST              = 132
	WM_NCLBUTTONDBLCLK        = 163
	WM_NCLBUTTONDOWN          = 161
	WM_NCLBUTTONUP            = 162
	WM_NCMBUTTONDBLCLK        = 169
	WM_NCMBUTTONDOWN          = 167
	WM_NCMBUTTONUP            = 168
	WM_NCXBUTTONDOWN          = 171
	WM_NCXBUTTONUP            = 172
	WM_NCXBUTTONDBLCLK        = 173
	WM_NCMOUSEHOVER           = 0x02A0
	WM_NCMOUSELEAVE           = 0x02A2
	WM_NCMOUSEMOVE            = 160
	WM_NCPAINT                = 133
	WM_NCRBUTTONDBLCLK        = 166
	WM_NCRBUTTONDOWN          = 164
	WM_NCRBUTTONUP            = 165
	WM_NEXTDLGCTL             = 40
	WM_NEXTMENU               = 531
	WM_NOTIFY                 = 78
	WM_NOTIFYFORMAT           = 85
	WM_NULL                   = 0
	WM_PAINT                  = 15
	WM_PAINTCLIPBOARD         = 777
	WM_PAINTICON              = 38
	WM_PALETTECHANGED         = 785
	WM_PALETTEISCHANGING      = 784
	WM_PARENTNOTIFY           = 528
	WM_PASTE                  = 770
	WM_PENWINFIRST            = 896
	WM_PENWINLAST             = 911
	WM_POWER                  = 72
	WM_POWERBROADCAST         = 536
	WM_PRINT                  = 791
	WM_PRINTCLIENT            = 792
	WM_QUERYDRAGICON          = 55
	WM_QUERYENDSESSION        = 17
	WM_QUERYNEWPALETTE        = 783
	WM_QUERYOPEN              = 19
	WM_QUEUESYNC              = 35
	WM_QUIT                   = 18
	WM_RENDERALLFORMATS       = 774
	WM_RENDERFORMAT           = 773
	WM_SETCURSOR              = 32
	WM_SETFOCUS               = 7
	WM_SETFONT                = 48
	WM_SETHOTKEY              = 50
	WM_SETICON                = 128
	WM_SETREDRAW              = 11
	WM_SETTEXT                = 12
	WM_SETTINGCHANGE          = 26
	WM_SHOWWINDOW             = 24
	WM_SIZE                   = 5
	WM_SIZECLIPBOARD          = 779
	WM_SIZING                 = 532
	WM_SPOOLERSTATUS          = 42
	WM_STYLECHANGED           = 125
	WM_STYLECHANGING          = 124
	WM_SYSCHAR                = 262
	WM_SYSCOLORCHANGE         = 21
	WM_SYSCOMMAND             = 274
	WM_SYSDEADCHAR            = 263
	WM_SYSKEYDOWN             = 260
	WM_SYSKEYUP               = 261
	WM_TCARD                  = 82
	WM_THEMECHANGED           = 794
	WM_TIMECHANGE             = 30
	WM_TIMER                  = 275
	WM_UNDO                   = 772
	WM_USER                   = 1024
	WM_USERCHANGED            = 84
	WM_VKEYTOITEM             = 46
	WM_VSCROLL                = 277
	WM_VSCROLLCLIPBOARD       = 778
	WM_WINDOWPOSCHANGED       = 71
	WM_WINDOWPOSCHANGING      = 70
	WM_WININICHANGE           = 26
	WM_KEYFIRST               = 256
	WM_KEYLAST                = 265
	WM_SYNCPAINT              = 136
	WM_MOUSEACTIVATE          = 33
	WM_MOUSEMOVE              = 512
	WM_LBUTTONDOWN            = 513
	WM_LBUTTONUP              = 514
	WM_LBUTTONDBLCLK          = 515
	WM_RBUTTONDOWN            = 516
	WM_RBUTTONUP              = 517
	WM_RBUTTONDBLCLK          = 518
	WM_MBUTTONDOWN            = 519
	WM_MBUTTONUP              = 520
	WM_MBUTTONDBLCLK          = 521
	WM_MOUSEWHEEL             = 522
	WM_MOUSEHWHEEL            = 526
	WM_MOUSEFIRST             = 512
	WM_XBUTTONDOWN            = 523
	WM_XBUTTONUP              = 524
	WM_XBUTTONDBLCLK          = 525
	WM_MOUSELAST              = 526
	WM_MOUSEHOVER             = 0x2A1
	WM_MOUSELEAVE             = 0x2A3
	WM_CLIPBOARDUPDATE        = 0x031D
)

// WM_ACTIVATE
const (
	WA_INACTIVE    = 0
	WA_ACTIVE      = 1
	WA_CLICKACTIVE = 2
)

const (
	LF_FACESIZE     = 32
	LF_FULLFACESIZE = 64
)

const (
	MM_MAX_NUMAXES      = 16
	MM_MAX_AXES_NAMELEN = 16
)

// Font weight constants
const (
	FW_DONTCARE   = 0
	FW_THIN       = 100
	FW_EXTRALIGHT = 200
	FW_ULTRALIGHT = FW_EXTRALIGHT
	FW_LIGHT      = 300
	FW_NORMAL     = 400
	FW_REGULAR    = 400
	FW_MEDIUM     = 500
	FW_SEMIBOLD   = 600
	FW_DEMIBOLD   = FW_SEMIBOLD
	FW_BOLD       = 700
	FW_EXTRABOLD  = 800
	FW_ULTRABOLD  = FW_EXTRABOLD
	FW_HEAVY      = 900
	FW_BLACK      = FW_HEAVY
)

// Charset constants
const (
	ANSI_CHARSET        = 0
	DEFAULT_CHARSET     = 1
	SYMBOL_CHARSET      = 2
	SHIFTJIS_CHARSET    = 128
	HANGEUL_CHARSET     = 129
	HANGUL_CHARSET      = 129
	GB2312_CHARSET      = 134
	CHINESEBIG5_CHARSET = 136
	GREEK_CHARSET       = 161
	TURKISH_CHARSET     = 162
	HEBREW_CHARSET      = 177
	ARABIC_CHARSET      = 178
	BALTIC_CHARSET      = 186
	RUSSIAN_CHARSET     = 204
	THAI_CHARSET        = 222
	EASTEUROPE_CHARSET  = 238
	OEM_CHARSET         = 255
	JOHAB_CHARSET       = 130
	VIETNAMESE_CHARSET  = 163
	MAC_CHARSET         = 77
)

// Font output precision constants
const (
	OUT_DEFAULT_PRECIS   = 0
	OUT_STRING_PRECIS    = 1
	OUT_CHARACTER_PRECIS = 2
	OUT_STROKE_PRECIS    = 3
	OUT_TT_PRECIS        = 4
	OUT_DEVICE_PRECIS    = 5
	OUT_RASTER_PRECIS    = 6
	OUT_TT_ONLY_PRECIS   = 7
	OUT_OUTLINE_PRECIS   = 8
	OUT_PS_ONLY_PRECIS   = 10
)

// Font clipping precision constants
const (
	CLIP_DEFAULT_PRECIS   = 0
	CLIP_CHARACTER_PRECIS = 1
	CLIP_STROKE_PRECIS    = 2
	CLIP_MASK             = 15
	CLIP_LH_ANGLES        = 16
	CLIP_TT_ALWAYS        = 32
	CLIP_EMBEDDED         = 128
)

// Font output quality constants
const (
	DEFAULT_QUALITY        = 0
	DRAFT_QUALITY          = 1
	PROOF_QUALITY          = 2
	NONANTIALIASED_QUALITY = 3
	ANTIALIASED_QUALITY    = 4
	CLEARTYPE_QUALITY      = 5
)

// Font pitch constants
const (
	DEFAULT_PITCH  = 0
	FIXED_PITCH    = 1
	VARIABLE_PITCH = 2
)

// Font family constants
const (
	FF_DECORATIVE = 80
	FF_DONTCARE   = 0
	FF_MODERN     = 48
	FF_ROMAN      = 16
	FF_SCRIPT     = 64
	FF_SWISS      = 32
)

// DeviceCapabilities capabilities
const (
	DC_FIELDS            = 1
	DC_PAPERS            = 2
	DC_PAPERSIZE         = 3
	DC_MINEXTENT         = 4
	DC_MAXEXTENT         = 5
	DC_BINS              = 6
	DC_DUPLEX            = 7
	DC_SIZE              = 8
	DC_EXTRA             = 9
	DC_VERSION           = 10
	DC_DRIVER            = 11
	DC_BINNAMES          = 12
	DC_ENUMRESOLUTIONS   = 13
	DC_FILEDEPENDENCIES  = 14
	DC_TRUETYPE          = 15
	DC_PAPERNAMES        = 16
	DC_ORIENTATION       = 17
	DC_COPIES            = 18
	DC_BINADJUST         = 19
	DC_EMF_COMPLIANT     = 20
	DC_DATATYPE_PRODUCED = 21
	DC_COLLATE           = 22
	DC_MANUFACTURER      = 23
	DC_MODEL             = 24
	DC_PERSONALITY       = 25
	DC_PRINTRATE         = 26
	DC_PRINTRATEUNIT     = 27
	DC_PRINTERMEM        = 28
	DC_MEDIAREADY        = 29
	DC_STAPLE            = 30
	DC_PRINTRATEPPM      = 31
	DC_COLORDEVICE       = 32
	DC_NUP               = 33
	DC_MEDIATYPENAMES    = 34
	DC_MEDIATYPES        = 35
)

// GetDeviceCaps index constants
const (
	DRIVERVERSION   = 0
	TECHNOLOGY      = 2
	HORZSIZE        = 4
	VERTSIZE        = 6
	HORZRES         = 8
	VERTRES         = 10
	LOGPIXELSX      = 88
	LOGPIXELSY      = 90
	BITSPIXEL       = 12
	PLANES          = 14
	NUMBRUSHES      = 16
	NUMPENS         = 18
	NUMFONTS        = 22
	NUMCOLORS       = 24
	NUMMARKERS      = 20
	ASPECTX         = 40
	ASPECTY         = 42
	ASPECTXY        = 44
	PDEVICESIZE     = 26
	CLIPCAPS        = 36
	SIZEPALETTE     = 104
	NUMRESERVED     = 106
	COLORRES        = 108
	PHYSICALWIDTH   = 110
	PHYSICALHEIGHT  = 111
	PHYSICALOFFSETX = 112
	PHYSICALOFFSETY = 113
	SCALINGFACTORX  = 114
	SCALINGFACTORY  = 115
	VREFRESH        = 116
	DESKTOPHORZRES  = 118
	DESKTOPVERTRES  = 117
	BLTALIGNMENT    = 119
	SHADEBLENDCAPS  = 120
	COLORMGMTCAPS   = 121
	RASTERCAPS      = 38
	CURVECAPS       = 28
	LINECAPS        = 30
	POLYGONALCAPS   = 32
	TEXTCAPS        = 34
)

// GetDeviceCaps TECHNOLOGY constants
const (
	DT_PLOTTER    = 0
	DT_RASDISPLAY = 1
	DT_RASPRINTER = 2
	DT_RASCAMERA  = 3
	DT_CHARSTREAM = 4
	DT_METAFILE   = 5
	DT_DISPFILE   = 6
)

// GetDeviceCaps SHADEBLENDCAPS constants
const (
	SB_NONE          = 0x00
	SB_CONST_ALPHA   = 0x01
	SB_PIXEL_ALPHA   = 0x02
	SB_PREMULT_ALPHA = 0x04
	SB_GRAD_RECT     = 0x10
	SB_GRAD_TRI      = 0x20
)

// GetDeviceCaps COLORMGMTCAPS constants
const (
	CM_NONE       = 0x00
	CM_DEVICE_ICM = 0x01
	CM_GAMMA_RAMP = 0x02
	CM_CMYK_COLOR = 0x04
)

// GetDeviceCaps RASTERCAPS constants
const (
	RC_BANDING      = 2
	RC_BITBLT       = 1
	RC_BITMAP64     = 8
	RC_DI_BITMAP    = 128
	RC_DIBTODEV     = 512
	RC_FLOODFILL    = 4096
	RC_GDI20_OUTPUT = 16
	RC_PALETTE      = 256
	RC_SCALING      = 4
	RC_STRETCHBLT   = 2048
	RC_STRETCHDIB   = 8192
	RC_DEVBITS      = 0x8000
	RC_OP_DX_OUTPUT = 0x4000
)

// GetDeviceCaps CURVECAPS constants
const (
	CC_NONE       = 0
	CC_CIRCLES    = 1
	CC_PIE        = 2
	CC_CHORD      = 4
	CC_ELLIPSES   = 8
	CC_WIDE       = 16
	CC_STYLED     = 32
	CC_WIDESTYLED = 64
	CC_INTERIORS  = 128
	CC_ROUNDRECT  = 256
)

// GetDeviceCaps LINECAPS constants
const (
	LC_NONE       = 0
	LC_POLYLINE   = 2
	LC_MARKER     = 4
	LC_POLYMARKER = 8
	LC_WIDE       = 16
	LC_STYLED     = 32
	LC_WIDESTYLED = 64
	LC_INTERIORS  = 128
)

// GetDeviceCaps POLYGONALCAPS constants
const (
	PC_NONE        = 0
	PC_POLYGON     = 1
	PC_POLYPOLYGON = 256
	PC_PATHS       = 512
	PC_RECTANGLE   = 2
	PC_WINDPOLYGON = 4
	PC_SCANLINE    = 8
	PC_TRAPEZOID   = 4
	PC_WIDE        = 16
	PC_STYLED      = 32
	PC_WIDESTYLED  = 64
	PC_INTERIORS   = 128
)

// GetDeviceCaps TEXTCAPS constants
const (
	TC_OP_CHARACTER = 1
	TC_OP_STROKE    = 2
	TC_CP_STROKE    = 4
	TC_CR_90        = 8
	TC_CR_ANY       = 16
	TC_SF_X_YINDEP  = 32
	TC_SA_DOUBLE    = 64
	TC_SA_INTEGER   = 128
	TC_SA_CONTIN    = 256
	TC_EA_DOUBLE    = 512
	TC_IA_ABLE      = 1024
	TC_UA_ABLE      = 2048
	TC_SO_ABLE      = 4096
	TC_RA_ABLE      = 8192
	TC_VA_ABLE      = 16384
	TC_RESERVED     = 32768
	TC_SCROLLBLT    = 65536
)

// Static control styles
const (
	SS_BITMAP          = 14
	SS_BLACKFRAME      = 7
	SS_BLACKRECT       = 4
	SS_CENTER          = 1
	SS_CENTERIMAGE     = 512
	SS_EDITCONTROL     = 0x2000
	SS_ENHMETAFILE     = 15
	SS_ETCHEDFRAME     = 18
	SS_ETCHEDHORZ      = 16
	SS_ETCHEDVERT      = 17
	SS_GRAYFRAME       = 8
	SS_GRAYRECT        = 5
	SS_ICON            = 3
	SS_LEFT            = 0
	SS_LEFTNOWORDWRAP  = 0xC
	SS_NOPREFIX        = 128
	SS_NOTIFY          = 256
	SS_OWNERDRAW       = 0xD
	SS_REALSIZECONTROL = 0x040
	SS_REALSIZEIMAGE   = 0x800
	SS_RIGHT           = 2
	SS_RIGHTJUST       = 0x400
	SS_SIMPLE          = 11
	SS_SUNKEN          = 4096
	SS_WHITEFRAME      = 9
	SS_WHITERECT       = 6
	SS_USERITEM        = 10
	SS_TYPEMASK        = 0x0000001F
	SS_ENDELLIPSIS     = 0x00004000
	SS_PATHELLIPSIS    = 0x00008000
	SS_WORDELLIPSIS    = 0x0000C000
	SS_ELLIPSISMASK    = 0x0000C000
)

// Edit styles
const (
	ES_LEFT        = 0x0000
	ES_CENTER      = 0x0001
	ES_RIGHT       = 0x0002
	ES_MULTILINE   = 0x0004
	ES_UPPERCASE   = 0x0008
	ES_LOWERCASE   = 0x0010
	ES_PASSWORD    = 0x0020
	ES_AUTOVSCROLL = 0x0040
	ES_AUTOHSCROLL = 0x0080
	ES_NOHIDESEL   = 0x0100
	ES_OEMCONVERT  = 0x0400
	ES_READONLY    = 0x0800
	ES_WANTRETURN  = 0x1000
	ES_NUMBER      = 0x2000
)

// Edit notifications
const (
	EN_SETFOCUS     = 0x0100
	EN_KILLFOCUS    = 0x0200
	EN_CHANGE       = 0x0300
	EN_UPDATE       = 0x0400
	EN_ERRSPACE     = 0x0500
	EN_MAXTEXT      = 0x0501
	EN_HSCROLL      = 0x0601
	EN_VSCROLL      = 0x0602
	EN_ALIGN_LTR_EC = 0x0700
	EN_ALIGN_RTL_EC = 0x0701
)

// Edit messages
const (
	EM_GETSEL               = 0x00B0
	EM_SETSEL               = 0x00B1
	EM_GETRECT              = 0x00B2
	EM_SETRECT              = 0x00B3
	EM_SETRECTNP            = 0x00B4
	EM_SCROLL               = 0x00B5
	EM_LINESCROLL           = 0x00B6
	EM_SCROLLCARET          = 0x00B7
	EM_GETMODIFY            = 0x00B8
	EM_SETMODIFY            = 0x00B9
	EM_GETLINECOUNT         = 0x00BA
	EM_LINEINDEX            = 0x00BB
	EM_SETHANDLE            = 0x00BC
	EM_GETHANDLE            = 0x00BD
	EM_GETTHUMB             = 0x00BE
	EM_LINELENGTH           = 0x00C1
	EM_REPLACESEL           = 0x00C2
	EM_GETLINE              = 0x00C4
	EM_LIMITTEXT            = 0x00C5
	EM_CANUNDO              = 0x00C6
	EM_UNDO                 = 0x00C7
	EM_FMTLINES             = 0x00C8
	EM_LINEFROMCHAR         = 0x00C9
	EM_SETTABSTOPS          = 0x00CB
	EM_SETPASSWORDCHAR      = 0x00CC
	EM_EMPTYUNDOBUFFER      = 0x00CD
	EM_GETFIRSTVISIBLELINE  = 0x00CE
	EM_SETREADONLY          = 0x00CF
	EM_SETWORDBREAKPROC     = 0x00D0
	EM_GETWORDBREAKPROC     = 0x00D1
	EM_GETPASSWORDCHAR      = 0x00D2
	EM_SETMARGINS           = 0x00D3
	EM_GETMARGINS           = 0x00D4
	EM_SETLIMITTEXT         = EM_LIMITTEXT
	EM_GETLIMITTEXT         = 0x00D5
	EM_POSFROMCHAR          = 0x00D6
	EM_CHARFROMPOS          = 0x00D7
	EM_SETIMESTATUS         = 0x00D8
	EM_GETIMESTATUS         = 0x00D9
	EM_SETCUEBANNER         = 0x1501
	EM_GETCUEBANNER         = 0x1502
	EM_AUTOURLDETECT        = 0x45B
	EM_CANPASTE             = 0x432
	EM_CANREDO              = 0x455
	EM_DISPLAYBAND          = 0x433
	EM_EXGETSEL             = 0x434
	EM_EXLIMITTEXT          = 0x435
	EM_EXLINEFROMCHAR       = 0x436
	EM_EXSETSEL             = 0x437
	EM_FINDTEXT             = 0x438
	EM_FINDTEXTEX           = 0x44F
	EM_FINDTEXTEXW          = 0x47C
	EM_FINDTEXTW            = 0x47B
	EM_FINDWORDBREAK        = 0x44C
	EM_FORMATRANGE          = 0x439
	EM_GETAUTOURLDETECT     = 0x45C
	EM_GETBIDIOPTIONS       = 0x4C9
	EM_GETCHARFORMAT        = 0x43A
	EM_GETEDITSTYLE         = 0x4CD
	EM_GETEVENTMASK         = 0x43B
	EM_GETIMECOLOR          = 0x469
	EM_GETIMECOMPMODE       = 0x47A
	EM_GETIMEOPTIONS        = 0x46B
	EM_GETLANGOPTIONS       = 0x479
	EM_GETOLEINTERFACE      = 0x43C
	EM_GETOPTIONS           = 0x44E
	EM_GETPARAFORMAT        = 0x43D
	EM_GETPUNCTUATION       = 0x465
	EM_GETREDONAME          = 0x457
	EM_GETSCROLLPOS         = 0x4DD
	EM_GETSELTEXT           = 0x43E
	EM_GETTEXTEX            = 0x45E
	EM_GETTEXTLENGTHEX      = 0x45F
	EM_GETTEXTMODE          = 0x45A
	EM_GETTEXTRANGE         = 0x44B
	EM_GETTYPOGRAPHYOPTIONS = 0x4CB
	EM_GETUNDONAME          = 0x456
	EM_GETWORDBREAKPROCEX   = 0x450
	EM_GETWORDWRAPMODE      = 0x467
	EM_GETZOOM              = 0x4E0
	EM_HIDESELECTION        = 0x43F
	EM_PASTESPECIAL         = 0x440
	EM_RECONVERSION         = 0x47D
	EM_REDO                 = 0x454
	EM_REQUESTRESIZE        = 0x441
	EM_SELECTIONTYPE        = 0x442
	EM_SETBIDIOPTIONS       = 0x4C8
	EM_SETBKGNDCOLOR        = 0x443
	EM_SETCHARFORMAT        = 0x444
	EM_SETEDITSTYLE         = 0x4CC
	EM_SETEVENTMASK         = 0x445
	EM_SETFONTSIZE          = 0x4DF
	EM_SETIMECOLOR          = 0x468
	EM_SETIMEOPTIONS        = 0x46A
	EM_SETLANGOPTIONS       = 0x478
	EM_SETOLECALLBACK       = 0x446
	EM_SETOPTIONS           = 0x44D
	EM_SETPALETTE           = 0x45D
	EM_SETPARAFORMAT        = 0x447
	EM_SETPUNCTUATION       = 0x464
	EM_SETSCROLLPOS         = 0x4DE
	EM_SETTARGETDEVICE      = 0x448
	EM_SETTEXTEX            = 0x461
	EM_SETTEXTMODE          = 0x459
	EM_SETTYPOGRAPHYOPTIONS = 0x4CA
	EM_SETUNDOLIMIT         = 0x452
	EM_SETWORDBREAKPROCEX   = 0x451
	EM_SETWORDWRAPMODE      = 0x466
	EM_SETZOOM              = 0x4E1
	EM_SHOWSCROLLBAR        = 0x460
	EM_STOPGROUPTYPING      = 0x458
	EM_STREAMIN             = 0x449
	EM_STREAMOUT            = 0x44A
)

const (
	CCM_FIRST            = 0x2000
	CCM_LAST             = CCM_FIRST + 0x0200
	CCM_SETBKCOLOR       = CCM_FIRST + 1
	CCM_SETCOLORSCHEME   = CCM_FIRST + 2
	CCM_GETCOLORSCHEME   = CCM_FIRST + 3
	CCM_GETDROPTARGET    = CCM_FIRST + 4
	CCM_SETUNICODEFORMAT = CCM_FIRST + 5
	CCM_GETUNICODEFORMAT = CCM_FIRST + 6
	CCM_SETVERSION       = CCM_FIRST + 7
	CCM_GETVERSION       = CCM_FIRST + 8
	CCM_SETNOTIFYWINDOW  = CCM_FIRST + 9
	CCM_SETWINDOWTHEME   = CCM_FIRST + 11
	CCM_DPISCALE         = CCM_FIRST + 12
)

// Common controls styles
const (
	CCS_TOP           = 1
	CCS_NOMOVEY       = 2
	CCS_BOTTOM        = 3
	CCS_NORESIZE      = 4
	CCS_NOPARENTALIGN = 8
	CCS_ADJUSTABLE    = 32
	CCS_NODIVIDER     = 64
	CCS_VERT          = 128
	CCS_LEFT          = 129
	CCS_NOMOVEX       = 130
	CCS_RIGHT         = 131
)

// ProgressBar messages
const (
	PROGRESS_CLASS  = "msctls_progress32"
	PBM_SETRANGE    = WM_USER + 1
	PBM_SETPOS      = WM_USER + 2
	PBM_DELTAPOS    = WM_USER + 3
	PBM_SETSTEP     = WM_USER + 4
	PBM_STEPIT      = WM_USER + 5
	PBM_SETRANGE32  = WM_USER + 6
	PBM_GETRANGE    = WM_USER + 7
	PBM_GETPOS      = WM_USER + 8
	PBM_SETBARCOLOR = WM_USER + 9
	PBM_SETMARQUEE  = WM_USER + 10
	PBM_SETBKCOLOR  = CCM_SETBKCOLOR
)

// Progress bar styles.
const (
	PBS_SMOOTH        = 0x01
	PBS_VERTICAL      = 0x04
	PBS_MARQUEE       = 0x08
	PBS_SMOOTHREVERSE = 0x10
)

// GetOpenFileName and GetSaveFileName extended flags
const (
	OFN_EX_NOPLACESBAR = 0x00000001
)

// GetOpenFileName and GetSaveFileName flags
const (
	OFN_ALLOWMULTISELECT     = 0x00000200
	OFN_CREATEPROMPT         = 0x00002000
	OFN_DONTADDTORECENT      = 0x02000000
	OFN_ENABLEHOOK           = 0x00000020
	OFN_ENABLEINCLUDENOTIFY  = 0x00400000
	OFN_ENABLESIZING         = 0x00800000
	OFN_ENABLETEMPLATE       = 0x00000040
	OFN_ENABLETEMPLATEHANDLE = 0x00000080
	OFN_EXPLORER             = 0x00080000
	OFN_EXTENSIONDIFFERENT   = 0x00000400
	OFN_FILEMUSTEXIST        = 0x00001000
	OFN_FORCESHOWHIDDEN      = 0x10000000
	OFN_HIDEREADONLY         = 0x00000004
	OFN_LONGNAMES            = 0x00200000
	OFN_NOCHANGEDIR          = 0x00000008
	OFN_NODEREFERENCELINKS   = 0x00100000
	OFN_NOLONGNAMES          = 0x00040000
	OFN_NONETWORKBUTTON      = 0x00020000
	OFN_NOREADONLYRETURN     = 0x00008000
	OFN_NOTESTFILECREATE     = 0x00010000
	OFN_NOVALIDATE           = 0x00000100
	OFN_OVERWRITEPROMPT      = 0x00000002
	OFN_PATHMUSTEXIST        = 0x00000800
	OFN_READONLY             = 0x00000001
	OFN_SHAREAWARE           = 0x00004000
	OFN_SHOWHELP             = 0x00000010
)

// SHBrowseForFolder flags
const (
	BIF_RETURNONLYFSDIRS    = 0x00000001
	BIF_DONTGOBELOWDOMAIN   = 0x00000002
	BIF_STATUSTEXT          = 0x00000004
	BIF_RETURNFSANCESTORS   = 0x00000008
	BIF_EDITBOX             = 0x00000010
	BIF_VALIDATE            = 0x00000020
	BIF_NEWDIALOGSTYLE      = 0x00000040
	BIF_BROWSEINCLUDEURLS   = 0x00000080
	BIF_USENEWUI            = BIF_EDITBOX | BIF_NEWDIALOGSTYLE
	BIF_UAHINT              = 0x00000100
	BIF_NONEWFOLDERBUTTON   = 0x00000200
	BIF_NOTRANSLATETARGETS  = 0x00000400
	BIF_BROWSEFORCOMPUTER   = 0x00001000
	BIF_BROWSEFORPRINTER    = 0x00002000
	BIF_BROWSEINCLUDEFILES  = 0x00004000
	BIF_SHAREABLE           = 0x00008000
	BIF_BROWSEFILEJUNCTIONS = 0x00010000
)

// MessageBox flags
const (
	MB_OK                = 0x00000000
	MB_OKCANCEL          = 0x00000001
	MB_ABORTRETRYIGNORE  = 0x00000002
	MB_YESNOCANCEL       = 0x00000003
	MB_YESNO             = 0x00000004
	MB_RETRYCANCEL       = 0x00000005
	MB_CANCELTRYCONTINUE = 0x00000006
	MB_ICONHAND          = 0x00000010
	MB_ICONQUESTION      = 0x00000020
	MB_ICONEXCLAMATION   = 0x00000030
	MB_ICONASTERISK      = 0x00000040
	MB_USERICON          = 0x00000080
	MB_ICONWARNING       = MB_ICONEXCLAMATION
	MB_ICONERROR         = MB_ICONHAND
	MB_ICONINFORMATION   = MB_ICONASTERISK
	MB_ICONSTOP          = MB_ICONHAND
	MB_DEFBUTTON1        = 0x00000000
	MB_DEFBUTTON2        = 0x00000100
	MB_DEFBUTTON3        = 0x00000200
	MB_DEFBUTTON4        = 0x00000300
	MB_TOPMOST           = 0x00040000
)

// COM
const (
	E_INVALIDARG  = 0x80070057
	E_OUTOFMEMORY = 0x8007000E
	E_UNEXPECTED  = 0x8000FFFF
)

const (
	S_OK               = 0
	S_FALSE            = 0x0001
	RPC_E_CHANGED_MODE = 0x80010106
)

// GetSystemMetrics constants
const (
	SM_CXSCREEN             = 0
	SM_CYSCREEN             = 1
	SM_CXVSCROLL            = 2
	SM_CYHSCROLL            = 3
	SM_CYCAPTION            = 4
	SM_CXBORDER             = 5
	SM_CYBORDER             = 6
	SM_CXDLGFRAME           = 7
	SM_CYDLGFRAME           = 8
	SM_CYVTHUMB             = 9
	SM_CXHTHUMB             = 10
	SM_CXICON               = 11
	SM_CYICON               = 12
	SM_CXCURSOR             = 13
	SM_CYCURSOR             = 14
	SM_CYMENU               = 15
	SM_CXFULLSCREEN         = 16
	SM_CYFULLSCREEN         = 17
	SM_CYKANJIWINDOW        = 18
	SM_MOUSEPRESENT         = 19
	SM_CYVSCROLL            = 20
	SM_CXHSCROLL            = 21
	SM_DEBUG                = 22
	SM_SWAPBUTTON           = 23
	SM_RESERVED1            = 24
	SM_RESERVED2            = 25
	SM_RESERVED3            = 26
	SM_RESERVED4            = 27
	SM_CXMIN                = 28
	SM_CYMIN                = 29
	SM_CXSIZE               = 30
	SM_CYSIZE               = 31
	SM_CXFRAME              = 32
	SM_CYFRAME              = 33
	SM_CXMINTRACK           = 34
	SM_CYMINTRACK           = 35
	SM_CXDOUBLECLK          = 36
	SM_CYDOUBLECLK          = 37
	SM_CXICONSPACING        = 38
	SM_CYICONSPACING        = 39
	SM_MENUDROPALIGNMENT    = 40
	SM_PENWINDOWS           = 41
	SM_DBCSENABLED          = 42
	SM_CMOUSEBUTTONS        = 43
	SM_CXFIXEDFRAME         = SM_CXDLGFRAME
	SM_CYFIXEDFRAME         = SM_CYDLGFRAME
	SM_CXSIZEFRAME          = SM_CXFRAME
	SM_CYSIZEFRAME          = SM_CYFRAME
	SM_SECURE               = 44
	SM_CXEDGE               = 45
	SM_CYEDGE               = 46
	SM_CXMINSPACING         = 47
	SM_CYMINSPACING         = 48
	SM_CXSMICON             = 49
	SM_CYSMICON             = 50
	SM_CYSMCAPTION          = 51
	SM_CXSMSIZE             = 52
	SM_CYSMSIZE             = 53
	SM_CXMENUSIZE           = 54
	SM_CYMENUSIZE           = 55
	SM_ARRANGE              = 56
	SM_CXMINIMIZED          = 57
	SM_CYMINIMIZED          = 58
	SM_CXMAXTRACK           = 59
	SM_CYMAXTRACK           = 60
	SM_CXMAXIMIZED          = 61
	SM_CYMAXIMIZED          = 62
	SM_NETWORK              = 63
	SM_CLEANBOOT            = 67
	SM_CXDRAG               = 68
	SM_CYDRAG               = 69
	SM_SHOWSOUNDS           = 70
	SM_CXMENUCHECK          = 71
	SM_CYMENUCHECK          = 72
	SM_SLOWMACHINE          = 73
	SM_MIDEASTENABLED       = 74
	SM_MOUSEWHEELPRESENT    = 75
	SM_XVIRTUALSCREEN       = 76
	SM_YVIRTUALSCREEN       = 77
	SM_CXVIRTUALSCREEN      = 78
	SM_CYVIRTUALSCREEN      = 79
	SM_CMONITORS            = 80
	SM_SAMEDISPLAYFORMAT    = 81
	SM_IMMENABLED           = 82
	SM_CXFOCUSBORDER        = 83
	SM_CYFOCUSBORDER        = 84
	SM_TABLETPC             = 86
	SM_MEDIACENTER          = 87
	SM_STARTER              = 88
	SM_SERVERR2             = 89
	SM_CMETRICS             = 97
	SM_REMOTESESSION        = 0x1000
	SM_SHUTTINGDOWN         = 0x2000
	SM_REMOTECONTROL        = 0x2001
	SM_CARETBLINKINGENABLED = 0x2002
)

const (
	CLSCTX_INPROC_SERVER   = 1
	CLSCTX_INPROC_HANDLER  = 2
	CLSCTX_LOCAL_SERVER    = 4
	CLSCTX_INPROC_SERVER16 = 8
	CLSCTX_REMOTE_SERVER   = 16
	CLSCTX_ALL             = CLSCTX_INPROC_SERVER | CLSCTX_INPROC_HANDLER | CLSCTX_LOCAL_SERVER | CLSCTX_REMOTE_SERVER
	CLSCTX_INPROC          = CLSCTX_INPROC_SERVER | CLSCTX_INPROC_HANDLER
	CLSCTX_SERVER          = CLSCTX_INPROC_SERVER | CLSCTX_LOCAL_SERVER | CLSCTX_REMOTE_SERVER
)

const (
	COINIT_APARTMENTTHREADED = 0x2
	COINIT_MULTITHREADED     = 0x0
	COINIT_DISABLE_OLE1DDE   = 0x4
	COINIT_SPEED_OVER_MEMORY = 0x8
)

const (
	DISPATCH_METHOD         = 1
	DISPATCH_PROPERTYGET    = 2
	DISPATCH_PROPERTYPUT    = 4
	DISPATCH_PROPERTYPUTREF = 8
)

const (
	CC_FASTCALL   = 0
	CC_CDECL      = 1
	CC_MSCPASCAL  = 2
	CC_PASCAL     = CC_MSCPASCAL
	CC_MACPASCAL  = 3
	CC_STDCALL    = 4
	CC_FPFASTCALL = 5
	CC_SYSCALL    = 6
	CC_MPWCDECL   = 7
	CC_MPWPASCAL  = 8
	CC_MAX        = 9
)

const (
	VT_EMPTY           = 0x0
	VT_NULL            = 0x1
	VT_I2              = 0x2
	VT_I4              = 0x3
	VT_R4              = 0x4
	VT_R8              = 0x5
	VT_CY              = 0x6
	VT_DATE            = 0x7
	VT_BSTR            = 0x8
	VT_DISPATCH        = 0x9
	VT_ERROR           = 0xA
	VT_BOOL            = 0xB
	VT_VARIANT         = 0xC
	VT_UNKNOWN         = 0xD
	VT_DECIMAL         = 0xE
	VT_I1              = 0x10
	VT_UI1             = 0x11
	VT_UI2             = 0x12
	VT_UI4             = 0x13
	VT_I8              = 0x14
	VT_UI8             = 0x15
	VT_INT             = 0x16
	VT_UINT            = 0x17
	VT_VOID            = 0x18
	VT_HRESULT         = 0x19
	VT_PTR             = 0x1A
	VT_SAFEARRAY       = 0x1B
	VT_CARRAY          = 0x1C
	VT_USERDEFINED     = 0x1D
	VT_LPSTR           = 0x1E
	VT_LPWSTR          = 0x1F
	VT_RECORD          = 0x24
	VT_INT_PTR         = 0x25
	VT_UINT_PTR        = 0x26
	VT_FILETIME        = 0x40
	VT_BLOB            = 0x41
	VT_STREAM          = 0x42
	VT_STORAGE         = 0x43
	VT_STREAMED_OBJECT = 0x44
	VT_STORED_OBJECT   = 0x45
	VT_BLOB_OBJECT     = 0x46
	VT_CF              = 0x47
	VT_CLSID           = 0x48
	VT_BSTR_BLOB       = 0xFFF
	VT_VECTOR          = 0x1000
	VT_ARRAY           = 0x2000
	VT_BYREF           = 0x4000
	VT_RESERVED        = 0x8000
	VT_ILLEGAL         = 0xFFFF
	VT_ILLEGALMASKED   = 0xFFF
	VT_TYPEMASK        = 0xFFF
)

const (
	DISPID_UNKNOWN     = -1
	DISPID_VALUE       = 0
	DISPID_PROPERTYPUT = -3
	DISPID_NEWENUM     = -4
	DISPID_EVALUATE    = -5
	DISPID_CONSTRUCTOR = -6
	DISPID_DESTRUCTOR  = -7
	DISPID_COLLECT     = -8
)

const (
	MONITOR_DEFAULTTONULL    = 0x00000000
	MONITOR_DEFAULTTOPRIMARY = 0x00000001
	MONITOR_DEFAULTTONEAREST = 0x00000002

	MONITORINFOF_PRIMARY = 0x00000001
)

const (
	CCHDEVICENAME = 32
	CCHFORMNAME   = 32
)

const (
	IDOK       = 1
	IDCANCEL   = 2
	IDABORT    = 3
	IDRETRY    = 4
	IDIGNORE   = 5
	IDYES      = 6
	IDNO       = 7
	IDCLOSE    = 8
	IDHELP     = 9
	IDTRYAGAIN = 10
	IDCONTINUE = 11
	IDTIMEOUT  = 32000
)

// Generic WM_NOTIFY notification codes
const (
	NM_FIRST           = 0
	NM_OUTOFMEMORY     = NM_FIRST - 1
	NM_CLICK           = NM_FIRST - 2
	NM_DBLCLK          = NM_FIRST - 3
	NM_RETURN          = NM_FIRST - 4
	NM_RCLICK          = NM_FIRST - 5
	NM_RDBLCLK         = NM_FIRST - 6
	NM_SETFOCUS        = NM_FIRST - 7
	NM_KILLFOCUS       = NM_FIRST - 8
	NM_CUSTOMDRAW      = NM_FIRST - 12
	NM_HOVER           = NM_FIRST - 13
	NM_NCHITTEST       = NM_FIRST - 14
	NM_KEYDOWN         = NM_FIRST - 15
	NM_RELEASEDCAPTURE = NM_FIRST - 16
	NM_SETCURSOR       = NM_FIRST - 17
	NM_CHAR            = NM_FIRST - 18
	NM_TOOLTIPSCREATED = NM_FIRST - 19
	NM_LAST            = NM_FIRST - 99
)

// ListView messages
const (
	// https://wiki.winehq.org/List_Of_Windows_Messages
	LVM_FIRST                    = 0x1000
	LVM_GETITEMCOUNT             = 0x1004
	LVM_SETIMAGELIST             = 0x1003
	LVM_GETIMAGELIST             = 0x1002
	LVM_GETITEM                  = 0x104B
	LVM_SETITEM                  = 0x104C
	LVM_INSERTITEM               = 0x104D
	LVM_DELETEITEM               = 0x1008
	LVM_DELETEALLITEMS           = 0x1009
	LVM_GETCALLBACKMASK          = 0x100A
	LVM_SETCALLBACKMASK          = 0x100B
	LVM_SETUNICODEFORMAT         = 0x2005
	LVM_GETNEXTITEM              = 0x100C
	LVM_FINDITEM                 = 0x1053
	LVM_GETITEMRECT              = 0x100E
	LVM_GETSTRINGWIDTH           = 0x1057
	LVM_HITTEST                  = 0x1012
	LVM_ENSUREVISIBLE            = 0x1013
	LVM_SCROLL                   = 0x1014
	LVM_REDRAWITEMS              = 0x1015
	LVM_ARRANGE                  = 0x1016
	LVM_EDITLABEL                = 0x1076
	LVM_GETEDITCONTROL           = 0x1018
	LVM_GETCOLUMN                = 0x105F
	LVM_SETCOLUMN                = 0x1060
	LVM_INSERTCOLUMN             = 0x1061
	LVM_DELETECOLUMN             = 0x101C
	LVM_GETCOLUMNWIDTH           = 0x101D
	LVM_SETCOLUMNWIDTH           = 0x101E
	LVM_GETHEADER                = 0x101F
	LVM_CREATEDRAGIMAGE          = 0x1021
	LVM_GETVIEWRECT              = 0x1022
	LVM_GETTEXTCOLOR             = 0x1023
	LVM_SETTEXTCOLOR             = 0x1024
	LVM_GETTEXTBKCOLOR           = 0x1025
	LVM_SETTEXTBKCOLOR           = 0x1026
	LVM_GETTOPINDEX              = 0x1027
	LVM_GETCOUNTPERPAGE          = 0x1028
	LVM_GETORIGIN                = 0x1029
	LVM_UPDATE                   = 0x102A
	LVM_SETITEMSTATE             = 0x102B
	LVM_GETITEMSTATE             = 0x102C
	LVM_GETITEMTEXT              = 0x1073
	LVM_SETITEMTEXT              = 0x1074
	LVM_SETITEMCOUNT             = 0x102F
	LVM_SORTITEMS                = 0x1030
	LVM_SETITEMPOSITION32        = 0x1031
	LVM_GETSELECTEDCOUNT         = 0x1032
	LVM_GETITEMSPACING           = 0x1033
	LVM_GETISEARCHSTRING         = 0x1075
	LVM_SETICONSPACING           = 0x1035
	LVM_SETEXTENDEDLISTVIEWSTYLE = 0x1036
	LVM_GETEXTENDEDLISTVIEWSTYLE = 0x1037
	LVM_GETSUBITEMRECT           = 0x1038
	LVM_SUBITEMHITTEST           = 0x1039
	LVM_SETCOLUMNORDERARRAY      = 0x103A
	LVM_GETCOLUMNORDERARRAY      = 0x103B
	LVM_SETHOTITEM               = 0x103C
	LVM_GETHOTITEM               = 0x103D
	LVM_SETHOTCURSOR             = 0x103E
	LVM_GETHOTCURSOR             = 0x103F
	LVM_APPROXIMATEVIEWRECT      = 0x1040
	LVM_SETWORKAREAS             = 0x1041
	LVM_GETWORKAREAS             = 0x1046
	LVM_GETNUMBEROFWORKAREAS     = 0x1049
	LVM_GETSELECTIONMARK         = 0x1042
	LVM_SETSELECTIONMARK         = 0x1043
	LVM_SETHOVERTIME             = 0x1047
	LVM_GETHOVERTIME             = 0x1048
	LVM_SETTOOLTIPS              = 0x104A
	LVM_GETTOOLTIPS              = 0x104E
	LVM_SORTITEMSEX              = 0x1051
	LVM_SETBKIMAGE               = 0x1044
	LVM_GETBKIMAGE               = 0x108B
	LVM_SETSELECTEDCOLUMN        = 0x108C
	LVM_SETVIEW                  = 0x108E
	LVM_GETVIEW                  = 0x108F
	LVM_INSERTGROUP              = 0x1091
	LVM_SETGROUPINFO             = 0x1093
	LVM_GETGROUPINFO             = 0x1095
	LVM_REMOVEGROUP              = 0x1096
	LVM_MOVEGROUP                = 0x1097
	LVM_GETGROUPCOUNT            = 0x1098
	LVM_GETGROUPINFOBYINDEX      = 0x1099
	LVM_MOVEITEMTOGROUP          = 0x109A
	LVM_GETGROUPRECT             = 0x1062
	LVM_SETGROUPMETRICS          = 0x109B
	LVM_GETGROUPMETRICS          = 0x109C
	LVM_ENABLEGROUPVIEW          = 0x109D
	LVM_SORTGROUPS               = 0x109E
	LVM_INSERTGROUPSORTED        = 0x109F
	LVM_REMOVEALLGROUPS          = 0x10A0
	LVM_HASGROUP                 = 0x10A1
	LVM_GETGROUPSTATE            = 0x105C
	LVM_GETFOCUSEDGROUP          = 0x105D
	LVM_SETTILEVIEWINFO          = 0x10A2
	LVM_GETTILEVIEWINFO          = 0x10A3
	LVM_SETTILEINFO              = 0x10A4
	LVM_GETTILEINFO              = 0x10A5
	LVM_SETINSERTMARK            = 0x10A6
	LVM_GETINSERTMARK            = 0x10A7
	LVM_INSERTMARKHITTEST        = 0x10A8
	LVM_GETINSERTMARKRECT        = 0x10A9
	LVM_SETINSERTMARKCOLOR       = 0x10AA
	LVM_GETINSERTMARKCOLOR       = 0x10AB
	LVM_SETINFOTIP               = 0x10AD
	LVM_GETSELECTEDCOLUMN        = 0x10AE
	LVM_ISGROUPVIEWENABLED       = 0x10AF
	LVM_GETOUTLINECOLOR          = 0x10B0
	LVM_SETOUTLINECOLOR          = 0x10B1
	LVM_CANCELEDITLABEL          = 0x10B3
	LVM_MAPINDEXTOID             = 0x10B4
	LVM_MAPIDTOINDEX             = 0x10B5
	LVM_ISITEMVISIBLE            = 0x10B6
	LVM_GETNEXTITEMINDEX         = 0x10D3
)

// ListView notifications
const (
	LVN_FIRST = -100

	LVN_ITEMCHANGING      = LVN_FIRST - 0
	LVN_ITEMCHANGED       = LVN_FIRST - 1
	LVN_INSERTITEM        = LVN_FIRST - 2
	LVN_DELETEITEM        = LVN_FIRST - 3
	LVN_DELETEALLITEMS    = LVN_FIRST - 4
	LVN_BEGINLABELEDITA   = LVN_FIRST - 5
	LVN_BEGINLABELEDITW   = LVN_FIRST - 75
	LVN_ENDLABELEDITA     = LVN_FIRST - 6
	LVN_ENDLABELEDITW     = LVN_FIRST - 76
	LVN_COLUMNCLICK       = LVN_FIRST - 8
	LVN_BEGINDRAG         = LVN_FIRST - 9
	LVN_BEGINRDRAG        = LVN_FIRST - 11
	LVN_ODCACHEHINT       = LVN_FIRST - 13
	LVN_ODFINDITEMA       = LVN_FIRST - 52
	LVN_ODFINDITEMW       = LVN_FIRST - 79
	LVN_ITEMACTIVATE      = LVN_FIRST - 14
	LVN_ODSTATECHANGED    = LVN_FIRST - 15
	LVN_HOTTRACK          = LVN_FIRST - 21
	LVN_GETDISPINFO       = LVN_FIRST - 50
	LVN_SETDISPINFO       = LVN_FIRST - 51
	LVN_KEYDOWN           = LVN_FIRST - 55
	LVN_MARQUEEBEGIN      = LVN_FIRST - 56
	LVN_GETINFOTIP        = LVN_FIRST - 57
	LVN_INCREMENTALSEARCH = LVN_FIRST - 62
	LVN_BEGINSCROLL       = LVN_FIRST - 80
	LVN_ENDSCROLL         = LVN_FIRST - 81
)

// ListView LVNI constants
const (
	LVNI_ALL         = 0
	LVNI_FOCUSED     = 1
	LVNI_SELECTED    = 2
	LVNI_CUT         = 4
	LVNI_DROPHILITED = 8
	LVNI_ABOVE       = 256
	LVNI_BELOW       = 512
	LVNI_TOLEFT      = 1024
	LVNI_TORIGHT     = 2048
)

// ListView styles
const (
	LVS_ICON            = 0x0000
	LVS_REPORT          = 0x0001
	LVS_SMALLICON       = 0x0002
	LVS_LIST            = 0x0003
	LVS_TYPEMASK        = 0x0003
	LVS_SINGLESEL       = 0x0004
	LVS_SHOWSELALWAYS   = 0x0008
	LVS_SORTASCENDING   = 0x0010
	LVS_SORTDESCENDING  = 0x0020
	LVS_SHAREIMAGELISTS = 0x0040
	LVS_NOLABELWRAP     = 0x0080
	LVS_AUTOARRANGE     = 0x0100
	LVS_EDITLABELS      = 0x0200
	LVS_OWNERDATA       = 0x1000
	LVS_NOSCROLL        = 0x2000
	LVS_TYPESTYLEMASK   = 0xFC00
	LVS_ALIGNTOP        = 0x0000
	LVS_ALIGNLEFT       = 0x0800
	LVS_ALIGNMASK       = 0x0C00
	LVS_OWNERDRAWFIXED  = 0x0400
	LVS_NOCOLUMNHEADER  = 0x4000
	LVS_NOSORTHEADER    = 0x8000
)

// ListView extended styles
const (
	LVS_EX_GRIDLINES        = 0x00000001
	LVS_EX_SUBITEMIMAGES    = 0x00000002
	LVS_EX_CHECKBOXES       = 0x00000004
	LVS_EX_TRACKSELECT      = 0x00000008
	LVS_EX_HEADERDRAGDROP   = 0x00000010
	LVS_EX_FULLROWSELECT    = 0x00000020
	LVS_EX_ONECLICKACTIVATE = 0x00000040
	LVS_EX_TWOCLICKACTIVATE = 0x00000080
	LVS_EX_FLATSB           = 0x00000100
	LVS_EX_REGIONAL         = 0x00000200
	LVS_EX_INFOTIP          = 0x00000400
	LVS_EX_UNDERLINEHOT     = 0x00000800
	LVS_EX_UNDERLINECOLD    = 0x00001000
	LVS_EX_MULTIWORKAREAS   = 0x00002000
	LVS_EX_LABELTIP         = 0x00004000
	LVS_EX_BORDERSELECT     = 0x00008000
	LVS_EX_DOUBLEBUFFER     = 0x00010000
	LVS_EX_HIDELABELS       = 0x00020000
	LVS_EX_SINGLEROW        = 0x00040000
	LVS_EX_SNAPTOGRID       = 0x00080000
	LVS_EX_SIMPLESELECT     = 0x00100000
)

// ListView column flags
const (
	LVCF_FMT     = 0x0001
	LVCF_WIDTH   = 0x0002
	LVCF_TEXT    = 0x0004
	LVCF_SUBITEM = 0x0008
	LVCF_IMAGE   = 0x0010
	LVCF_ORDER   = 0x0020
)

// ListView column format constants
const (
	LVCFMT_LEFT            = 0x0000
	LVCFMT_RIGHT           = 0x0001
	LVCFMT_CENTER          = 0x0002
	LVCFMT_JUSTIFYMASK     = 0x0003
	LVCFMT_IMAGE           = 0x0800
	LVCFMT_BITMAP_ON_RIGHT = 0x1000
	LVCFMT_COL_HAS_IMAGES  = 0x8000
)

// ListView item flags
const (
	LVIF_TEXT        = 0x00000001
	LVIF_IMAGE       = 0x00000002
	LVIF_PARAM       = 0x00000004
	LVIF_STATE       = 0x00000008
	LVIF_INDENT      = 0x00000010
	LVIF_NORECOMPUTE = 0x00000800
	LVIF_GROUPID     = 0x00000100
	LVIF_COLUMNS     = 0x00000200
)

// ListView item states
const (
	LVIS_FOCUSED        = 1
	LVIS_SELECTED       = 2
	LVIS_CUT            = 4
	LVIS_DROPHILITED    = 8
	LVIS_OVERLAYMASK    = 0xF00
	LVIS_STATEIMAGEMASK = 0xF000
)

// ListView hit test constants
const (
	LVHT_NOWHERE         = 0x00000001
	LVHT_ONITEMICON      = 0x00000002
	LVHT_ONITEMLABEL     = 0x00000004
	LVHT_ONITEMSTATEICON = 0x00000008
	LVHT_ONITEM          = LVHT_ONITEMICON | LVHT_ONITEMLABEL | LVHT_ONITEMSTATEICON

	LVHT_ABOVE   = 0x00000008
	LVHT_BELOW   = 0x00000010
	LVHT_TORIGHT = 0x00000020
	LVHT_TOLEFT  = 0x00000040
)

// ListView image list types
const (
	LVSIL_NORMAL      = 0
	LVSIL_SMALL       = 1
	LVSIL_STATE       = 2
	LVSIL_GROUPHEADER = 3
)

// InitCommonControlsEx flags
const (
	ICC_LISTVIEW_CLASSES   = 1
	ICC_TREEVIEW_CLASSES   = 2
	ICC_BAR_CLASSES        = 4
	ICC_TAB_CLASSES        = 8
	ICC_UPDOWN_CLASS       = 16
	ICC_PROGRESS_CLASS     = 32
	ICC_HOTKEY_CLASS       = 64
	ICC_ANIMATE_CLASS      = 128
	ICC_WIN95_CLASSES      = 255
	ICC_DATE_CLASSES       = 256
	ICC_USEREX_CLASSES     = 512
	ICC_COOL_CLASSES       = 1024
	ICC_INTERNET_CLASSES   = 2048
	ICC_PAGESCROLLER_CLASS = 4096
	ICC_NATIVEFNTCTL_CLASS = 8192
	INFOTIPSIZE            = 1024
	ICC_STANDARD_CLASSES   = 0x00004000
	ICC_LINK_CLASS         = 0x00008000
)

// Dialog Codes
const (
	DLGC_WANTARROWS      = 0x0001
	DLGC_WANTTAB         = 0x0002
	DLGC_WANTALLKEYS     = 0x0004
	DLGC_WANTMESSAGE     = 0x0004
	DLGC_HASSETSEL       = 0x0008
	DLGC_DEFPUSHBUTTON   = 0x0010
	DLGC_UNDEFPUSHBUTTON = 0x0020
	DLGC_RADIOBUTTON     = 0x0040
	DLGC_WANTCHARS       = 0x0080
	DLGC_STATIC          = 0x0100
	DLGC_BUTTON          = 0x2000
)

// Get/SetWindowWord/Long offsets for use with WC_DIALOG windows
const (
	DWL_MSGRESULT = 0
	DWL_DLGPROC   = 4
	DWL_USER      = 8
)

// Registry predefined keys
const (
	HKEY_CLASSES_ROOT     HKEY = 0x80000000
	HKEY_CURRENT_USER     HKEY = 0x80000001
	HKEY_LOCAL_MACHINE    HKEY = 0x80000002
	HKEY_USERS            HKEY = 0x80000003
	HKEY_PERFORMANCE_DATA HKEY = 0x80000004
	HKEY_CURRENT_CONFIG   HKEY = 0x80000005
	HKEY_DYN_DATA         HKEY = 0x80000006
)

// Registry Key Security and Access Rights
const (
	KEY_ALL_ACCESS         = 0xF003F
	KEY_CREATE_SUB_KEY     = 0x0004
	KEY_ENUMERATE_SUB_KEYS = 0x0008
	KEY_NOTIFY             = 0x0010
	KEY_QUERY_VALUE        = 0x0001
	KEY_SET_VALUE          = 0x0002
	KEY_READ               = 0x20019
	KEY_WRITE              = 0x20006
)

const (
	NFR_ANSI    = 1
	NFR_UNICODE = 2
	NF_QUERY    = 3
	NF_REQUERY  = 4
)

// Registry value types
const (
	RRF_RT_REG_NONE         = 0x00000001
	RRF_RT_REG_SZ           = 0x00000002
	RRF_RT_REG_EXPAND_SZ    = 0x00000004
	RRF_RT_REG_BINARY       = 0x00000008
	RRF_RT_REG_DWORD        = 0x00000010
	RRF_RT_REG_MULTI_SZ     = 0x00000020
	RRF_RT_REG_QWORD        = 0x00000040
	RRF_RT_DWORD            = RRF_RT_REG_BINARY | RRF_RT_REG_DWORD
	RRF_RT_QWORD            = RRF_RT_REG_BINARY | RRF_RT_REG_QWORD
	RRF_RT_ANY              = 0x0000FFFF
	RRF_NOEXPAND            = 0x10000000
	RRF_ZEROONFAILURE       = 0x20000000
	REG_PROCESS_APPKEY      = 0x00000001
	REG_MUI_STRING_TRUNCATE = 0x00000001
)

// PeekMessage wRemoveMsg value
const (
	PM_NOREMOVE = 0x000
	PM_REMOVE   = 0x001
	PM_NOYIELD  = 0x002
)

// ImageList flags
const (
	ILC_MASK             = 0x00000001
	ILC_COLOR            = 0x00000000
	ILC_COLORDDB         = 0x000000FE
	ILC_COLOR4           = 0x00000004
	ILC_COLOR8           = 0x00000008
	ILC_COLOR16          = 0x00000010
	ILC_COLOR24          = 0x00000018
	ILC_COLOR32          = 0x00000020
	ILC_PALETTE          = 0x00000800
	ILC_MIRROR           = 0x00002000
	ILC_PERITEMMIRROR    = 0x00008000
	ILC_ORIGINALSIZE     = 0x00010000
	ILC_HIGHQUALITYSCALE = 0x00020000
)

// Keystroke Message Flags
const (
	KF_EXTENDED = 0x0100
	KF_DLGMODE  = 0x0800
	KF_MENUMODE = 0x1000
	KF_ALTDOWN  = 0x2000
	KF_REPEAT   = 0x4000
	KF_UP       = 0x8000
)

// Virtual-Key Codes
const (
	VK_LBUTTON             = 0x01
	VK_RBUTTON             = 0x02
	VK_CANCEL              = 0x03
	VK_MBUTTON             = 0x04
	VK_XBUTTON1            = 0x05
	VK_XBUTTON2            = 0x06
	VK_BACK                = 0x08
	VK_TAB                 = 0x09
	VK_CLEAR               = 0x0C
	VK_RETURN              = 0x0D
	VK_SHIFT               = 0x10
	VK_CONTROL             = 0x11
	VK_MENU                = 0x12
	VK_PAUSE               = 0x13
	VK_CAPITAL             = 0x14
	VK_KANA                = 0x15
	VK_HANGEUL             = 0x15
	VK_HANGUL              = 0x15
	VK_IME_ON              = 0x16
	VK_JUNJA               = 0x17
	VK_FINAL               = 0x18
	VK_HANJA               = 0x19
	VK_KANJI               = 0x19
	VK_IME_OFF             = 0x1A
	VK_ESCAPE              = 0x1B
	VK_CONVERT             = 0x1C
	VK_NONCONVERT          = 0x1D
	VK_ACCEPT              = 0x1E
	VK_MODECHANGE          = 0x1F
	VK_SPACE               = 0x20
	VK_PRIOR               = 0x21
	VK_NEXT                = 0x22
	VK_END                 = 0x23
	VK_HOME                = 0x24
	VK_LEFT                = 0x25
	VK_UP                  = 0x26
	VK_RIGHT               = 0x27
	VK_DOWN                = 0x28
	VK_SELECT              = 0x29
	VK_PRINT               = 0x2A
	VK_EXECUTE             = 0x2B
	VK_SNAPSHOT            = 0x2C
	VK_INSERT              = 0x2D
	VK_DELETE              = 0x2E
	VK_HELP                = 0x2F
	VK_LWIN                = 0x5B
	VK_RWIN                = 0x5C
	VK_APPS                = 0x5D
	VK_SLEEP               = 0x5F
	VK_NUMPAD0             = 0x60
	VK_NUMPAD1             = 0x61
	VK_NUMPAD2             = 0x62
	VK_NUMPAD3             = 0x63
	VK_NUMPAD4             = 0x64
	VK_NUMPAD5             = 0x65
	VK_NUMPAD6             = 0x66
	VK_NUMPAD7             = 0x67
	VK_NUMPAD8             = 0x68
	VK_NUMPAD9             = 0x69
	VK_MULTIPLY            = 0x6A
	VK_ADD                 = 0x6B
	VK_SEPARATOR           = 0x6C
	VK_SUBTRACT            = 0x6D
	VK_DECIMAL             = 0x6E
	VK_DIVIDE              = 0x6F
	VK_F1                  = 0x70
	VK_F2                  = 0x71
	VK_F3                  = 0x72
	VK_F4                  = 0x73
	VK_F5                  = 0x74
	VK_F6                  = 0x75
	VK_F7                  = 0x76
	VK_F8                  = 0x77
	VK_F9                  = 0x78
	VK_F10                 = 0x79
	VK_F11                 = 0x7A
	VK_F12                 = 0x7B
	VK_F13                 = 0x7C
	VK_F14                 = 0x7D
	VK_F15                 = 0x7E
	VK_F16                 = 0x7F
	VK_F17                 = 0x80
	VK_F18                 = 0x81
	VK_F19                 = 0x82
	VK_F20                 = 0x83
	VK_F21                 = 0x84
	VK_F22                 = 0x85
	VK_F23                 = 0x86
	VK_F24                 = 0x87
	VK_NUMLOCK             = 0x90
	VK_SCROLL              = 0x91
	VK_OEM_NEC_EQUAL       = 0x92
	VK_OEM_FJ_JISHO        = 0x92
	VK_OEM_FJ_MASSHOU      = 0x93
	VK_OEM_FJ_TOUROKU      = 0x94
	VK_OEM_FJ_LOYA         = 0x95
	VK_OEM_FJ_ROYA         = 0x96
	VK_LSHIFT              = 0xA0
	VK_RSHIFT              = 0xA1
	VK_LCONTROL            = 0xA2
	VK_RCONTROL            = 0xA3
	VK_LMENU               = 0xA4
	VK_RMENU               = 0xA5
	VK_BROWSER_BACK        = 0xA6
	VK_BROWSER_FORWARD     = 0xA7
	VK_BROWSER_REFRESH     = 0xA8
	VK_BROWSER_STOP        = 0xA9
	VK_BROWSER_SEARCH      = 0xAA
	VK_BROWSER_FAVORITES   = 0xAB
	VK_BROWSER_HOME        = 0xAC
	VK_VOLUME_MUTE         = 0xAD
	VK_VOLUME_DOWN         = 0xAE
	VK_VOLUME_UP           = 0xAF
	VK_MEDIA_NEXT_TRACK    = 0xB0
	VK_MEDIA_PREV_TRACK    = 0xB1
	VK_MEDIA_STOP          = 0xB2
	VK_MEDIA_PLAY_PAUSE    = 0xB3
	VK_LAUNCH_MAIL         = 0xB4
	VK_LAUNCH_MEDIA_SELECT = 0xB5
	VK_LAUNCH_APP1         = 0xB6
	VK_LAUNCH_APP2         = 0xB7
	VK_OEM_1               = 0xBA
	VK_OEM_PLUS            = 0xBB
	VK_OEM_COMMA           = 0xBC
	VK_OEM_MINUS           = 0xBD
	VK_OEM_PERIOD          = 0xBE
	VK_OEM_2               = 0xBF
	VK_OEM_3               = 0xC0
	VK_OEM_4               = 0xDB
	VK_OEM_5               = 0xDC
	VK_OEM_6               = 0xDD
	VK_OEM_7               = 0xDE
	VK_OEM_8               = 0xDF
	VK_OEM_AX              = 0xE1
	VK_OEM_102             = 0xE2
	VK_ICO_HELP            = 0xE3
	VK_ICO_00              = 0xE4
	VK_PROCESSKEY          = 0xE5
	VK_ICO_CLEAR           = 0xE6
	VK_PACKET              = 0xE7
	VK_OEM_RESET           = 0xE9
	VK_OEM_JUMP            = 0xEA
	VK_OEM_PA1             = 0xEB
	VK_OEM_PA2             = 0xEC
	VK_OEM_PA3             = 0xED
	VK_OEM_WSCTRL          = 0xEE
	VK_OEM_CUSEL           = 0xEF
	VK_OEM_ATTN            = 0xF0
	VK_OEM_FINISH          = 0xF1
	VK_OEM_COPY            = 0xF2
	VK_OEM_AUTO            = 0xF3
	VK_OEM_ENLW            = 0xF4
	VK_OEM_BACKTAB         = 0xF5
	VK_ATTN                = 0xF6
	VK_CRSEL               = 0xF7
	VK_EXSEL               = 0xF8
	VK_EREOF               = 0xF9
	VK_PLAY                = 0xFA
	VK_ZOOM                = 0xFB
	VK_NONAME              = 0xFC
	VK_PA1                 = 0xFD
	VK_OEM_CLEAR           = 0xFE
)

// Registry Value Types
const (
	REG_NONE                       = 0
	REG_SZ                         = 1
	REG_EXPAND_SZ                  = 2
	REG_BINARY                     = 3
	REG_DWORD                      = 4
	REG_DWORD_LITTLE_ENDIAN        = 4
	REG_DWORD_BIG_ENDIAN           = 5
	REG_LINK                       = 6
	REG_MULTI_SZ                   = 7
	REG_RESOURCE_LIST              = 8
	REG_FULL_RESOURCE_DESCRIPTOR   = 9
	REG_RESOURCE_REQUIREMENTS_LIST = 10
	REG_QWORD                      = 11
	REG_QWORD_LITTLE_ENDIAN        = 11
)

// Tooltip styles
const (
	TTS_ALWAYSTIP      = 0x01
	TTS_NOPREFIX       = 0x02
	TTS_NOANIMATE      = 0x10
	TTS_NOFADE         = 0x20
	TTS_BALLOON        = 0x40
	TTS_CLOSE          = 0x80
	TTS_USEVISUALSTYLE = 0x100
)

// Tooltip messages
const (
	TTM_ACTIVATE        = WM_USER + 1
	TTM_SETDELAYTIME    = WM_USER + 3
	TTM_ADDTOOL         = WM_USER + 4
	TTM_DELTOOL         = WM_USER + 5
	TTM_NEWTOOLRECT     = WM_USER + 6
	TTM_RELAYEVENT      = WM_USER + 7
	TTM_GETTOOLINFO     = WM_USER + 8
	TTM_SETTOOLINFO     = WM_USER + 9
	TTM_HITTEST         = WM_USER + 10
	TTM_GETTEXT         = WM_USER + 11
	TTM_UPDATETIPTEXT   = WM_USER + 12
	TTM_GETTOOLCOUNT    = WM_USER + 13
	TTM_ENUMTOOLS       = WM_USER + 14
	TTM_GETCURRENTTOOL  = WM_USER + 15
	TTM_WINDOWFROMPOINT = WM_USER + 16
	TTM_TRACKACTIVATE   = WM_USER + 17
	TTM_TRACKPOSITION   = WM_USER + 18
	TTM_SETTIPBKCOLOR   = WM_USER + 19
	TTM_SETTIPTEXTCOLOR = WM_USER + 20
	TTM_GETDELAYTIME    = WM_USER + 21
	TTM_GETTIPBKCOLOR   = WM_USER + 22
	TTM_GETTIPTEXTCOLOR = WM_USER + 23
	TTM_SETMAXTIPWIDTH  = WM_USER + 24
	TTM_GETMAXTIPWIDTH  = WM_USER + 25
	TTM_SETMARGIN       = WM_USER + 26
	TTM_GETMARGIN       = WM_USER + 27
	TTM_POP             = WM_USER + 28
	TTM_UPDATE          = WM_USER + 29
	TTM_GETBUBBLESIZE   = WM_USER + 30
	TTM_ADJUSTRECT      = WM_USER + 31
	TTM_SETTITLE        = WM_USER + 32
	TTM_POPUP           = WM_USER + 34
	TTM_GETTITLE        = WM_USER + 35
)

// Tooltip icons
const (
	TTI_NONE          = 0
	TTI_INFO          = 1
	TTI_WARNING       = 2
	TTI_ERROR         = 3
	TTI_INFO_LARGE    = 4
	TTI_WARNING_LARGE = 5
	TTI_ERROR_LARGE   = 6
)

// Tooltip notifications
const (
	TTN_FIRST       = -520
	TTN_LAST        = -549
	TTN_GETDISPINFO = TTN_FIRST
	TTN_SHOW        = TTN_FIRST - 1
	TTN_POP         = TTN_FIRST - 2
	TTN_LINKCLICK   = TTN_FIRST - 3
	TTN_NEEDTEXT    = TTN_GETDISPINFO
)

const (
	TTF_IDISHWND    = 0x0001
	TTF_CENTERTIP   = 0x0002
	TTF_RTLREADING  = 0x0004
	TTF_SUBCLASS    = 0x0010
	TTF_TRACK       = 0x0020
	TTF_ABSOLUTE    = 0x0080
	TTF_TRANSPARENT = 0x0100
	TTF_PARSELINKS  = 0x1000
	TTF_DI_SETITEM  = 0x8000
)

const (
	SWP_NOSIZE         = 0x0001
	SWP_NOMOVE         = 0x0002
	SWP_NOZORDER       = 0x0004
	SWP_NOREDRAW       = 0x0008
	SWP_NOACTIVATE     = 0x0010
	SWP_FRAMECHANGED   = 0x0020
	SWP_SHOWWINDOW     = 0x0040
	SWP_HIDEWINDOW     = 0x0080
	SWP_NOCOPYBITS     = 0x0100
	SWP_NOOWNERZORDER  = 0x0200
	SWP_NOSENDCHANGING = 0x0400
	SWP_DRAWFRAME      = SWP_FRAMECHANGED
	SWP_NOREPOSITION   = SWP_NOOWNERZORDER
	SWP_DEFERERASE     = 0x2000
	SWP_ASYNCWINDOWPOS = 0x4000
)

// Predefined window handles
const (
	HWND_BROADCAST = HWND(0xFFFF)
	HWND_BOTTOM    = HWND(1)
	HWND_NOTOPMOST = ^HWND(1) // -2
	HWND_TOP       = HWND(0)
	HWND_TOPMOST   = ^HWND(0) // -1
	HWND_DESKTOP   = HWND(0)
	HWND_MESSAGE   = ^HWND(2) // -3
)

// Pen types
const (
	PS_COSMETIC  = 0x00000000
	PS_GEOMETRIC = 0x00010000
	PS_TYPE_MASK = 0x000F0000
)

// Pen styles
const (
	PS_SOLID       = 0
	PS_DASH        = 1
	PS_DOT         = 2
	PS_DASHDOT     = 3
	PS_DASHDOTDOT  = 4
	PS_NULL        = 5
	PS_INSIDEFRAME = 6
	PS_USERSTYLE   = 7
	PS_ALTERNATE   = 8
	PS_STYLE_MASK  = 0x0000000F
)

// Pen cap types
const (
	PS_ENDCAP_ROUND  = 0x00000000
	PS_ENDCAP_SQUARE = 0x00000100
	PS_ENDCAP_FLAT   = 0x00000200
	PS_ENDCAP_MASK   = 0x00000F00
)

// Pen join types
const (
	PS_JOIN_ROUND = 0x00000000
	PS_JOIN_BEVEL = 0x00001000
	PS_JOIN_MITER = 0x00002000
	PS_JOIN_MASK  = 0x0000F000
)

// Hatch styles
const (
	HS_HORIZONTAL = 0
	HS_VERTICAL   = 1
	HS_FDIAGONAL  = 2
	HS_BDIAGONAL  = 3
	HS_CROSS      = 4
	HS_DIAGCROSS  = 5
)

// Stock Logical Objects
const (
	WHITE_BRUSH         = 0
	LTGRAY_BRUSH        = 1
	GRAY_BRUSH          = 2
	DKGRAY_BRUSH        = 3
	BLACK_BRUSH         = 4
	NULL_BRUSH          = 5
	HOLLOW_BRUSH        = NULL_BRUSH
	WHITE_PEN           = 6
	BLACK_PEN           = 7
	NULL_PEN            = 8
	OEM_FIXED_FONT      = 10
	ANSI_FIXED_FONT     = 11
	ANSI_VAR_FONT       = 12
	SYSTEM_FONT         = 13
	DEVICE_DEFAULT_FONT = 14
	DEFAULT_PALETTE     = 15
	SYSTEM_FIXED_FONT   = 16
	DEFAULT_GUI_FONT    = 17
	DC_BRUSH            = 18
	DC_PEN              = 19
)

// Brush styles
const (
	BS_SOLID         = 0
	BS_NULL          = 1
	BS_HOLLOW        = BS_NULL
	BS_HATCHED       = 2
	BS_PATTERN       = 3
	BS_INDEXED       = 4
	BS_DIBPATTERN    = 5
	BS_DIBPATTERNPT  = 6
	BS_PATTERN8X8    = 7
	BS_DIBPATTERN8X8 = 8
	BS_MONOPATTERN   = 9
)

// TRACKMOUSEEVENT flags
const (
	TME_HOVER     = 0x00000001
	TME_LEAVE     = 0x00000002
	TME_NONCLIENT = 0x00000010
	TME_QUERY     = 0x40000000
	TME_CANCEL    = 0x80000000

	HOVER_DEFAULT = 0xFFFFFFFF
)

// WM_NCHITTEST and MOUSEHOOKSTRUCT Mouse Position Codes
const (
	HTERROR       = -2
	HTTRANSPARENT = -1
	HTNOWHERE     = 0
	HTCLIENT      = 1
	HTCAPTION     = 2
	HTSYSMENU     = 3
	HTGROWBOX     = 4
	HTSIZE        = HTGROWBOX
	HTMENU        = 5
	HTHSCROLL     = 6
	HTVSCROLL     = 7
	HTMINBUTTON   = 8
	HTMAXBUTTON   = 9
	HTLEFT        = 10
	HTRIGHT       = 11
	HTTOP         = 12
	HTTOPLEFT     = 13
	HTTOPRIGHT    = 14
	HTBOTTOM      = 15
	HTBOTTOMLEFT  = 16
	HTBOTTOMRIGHT = 17
	HTBORDER      = 18
	HTREDUCE      = HTMINBUTTON
	HTZOOM        = HTMAXBUTTON
	HTSIZEFIRST   = HTLEFT
	HTSIZELAST    = HTBOTTOMRIGHT
	HTOBJECT      = 19
	HTCLOSE       = 20
	HTHELP        = 21
)

// DrawText[Ex] format flags
const (
	DT_TOP                  = 0x00000000
	DT_LEFT                 = 0x00000000
	DT_CENTER               = 0x00000001
	DT_RIGHT                = 0x00000002
	DT_VCENTER              = 0x00000004
	DT_BOTTOM               = 0x00000008
	DT_WORDBREAK            = 0x00000010
	DT_SINGLELINE           = 0x00000020
	DT_EXPANDTABS           = 0x00000040
	DT_TABSTOP              = 0x00000080
	DT_NOCLIP               = 0x00000100
	DT_EXTERNALLEADING      = 0x00000200
	DT_CALCRECT             = 0x00000400
	DT_NOPREFIX             = 0x00000800
	DT_INTERNAL             = 0x00001000
	DT_EDITCONTROL          = 0x00002000
	DT_PATH_ELLIPSIS        = 0x00004000
	DT_END_ELLIPSIS         = 0x00008000
	DT_MODIFYSTRING         = 0x00010000
	DT_RTLREADING           = 0x00020000
	DT_WORD_ELLIPSIS        = 0x00040000
	DT_NOFULLWIDTHCHARBREAK = 0x00080000
	DT_HIDEPREFIX           = 0x00100000
	DT_PREFIXONLY           = 0x00200000
)

const CLR_INVALID = 0xFFFFFFFF

// Background Modes
const (
	TRANSPARENT = 1
	OPAQUE      = 2
	BKMODE_LAST = 2
)

// Global Memory Flags
const (
	GMEM_FIXED          = 0x0000
	GMEM_MOVEABLE       = 0x0002
	GMEM_NOCOMPACT      = 0x0010
	GMEM_NODISCARD      = 0x0020
	GMEM_ZEROINIT       = 0x0040
	GMEM_MODIFY         = 0x0080
	GMEM_DISCARDABLE    = 0x0100
	GMEM_NOT_BANKED     = 0x1000
	GMEM_SHARE          = 0x2000
	GMEM_DDESHARE       = 0x2000
	GMEM_NOTIFY         = 0x4000
	GMEM_LOWER          = GMEM_NOT_BANKED
	GMEM_VALID_FLAGS    = 0x7F72
	GMEM_INVALID_HANDLE = 0x8000
	GHND                = GMEM_MOVEABLE | GMEM_ZEROINIT
	GPTR                = GMEM_FIXED | GMEM_ZEROINIT
)

// Ternary raster operations
const (
	SRCCOPY        = 0x00CC0020
	SRCPAINT       = 0x00EE0086
	SRCAND         = 0x008800C6
	SRCINVERT      = 0x00660046
	SRCERASE       = 0x00440328
	NOTSRCCOPY     = 0x00330008
	NOTSRCERASE    = 0x001100A6
	MERGECOPY      = 0x00C000CA
	MERGEPAINT     = 0x00BB0226
	PATCOPY        = 0x00F00021
	PATPAINT       = 0x00FB0A09
	PATINVERT      = 0x005A0049
	DSTINVERT      = 0x00550009
	BLACKNESS      = 0x00000042
	WHITENESS      = 0x00FF0062
	NOMIRRORBITMAP = 0x80000000
	CAPTUREBLT     = 0x40000000
)

// Clipboard formats
const (
	CF_TEXT            = 1
	CF_BITMAP          = 2
	CF_METAFILEPICT    = 3
	CF_SYLK            = 4
	CF_DIF             = 5
	CF_TIFF            = 6
	CF_OEMTEXT         = 7
	CF_DIB             = 8
	CF_PALETTE         = 9
	CF_PENDATA         = 10
	CF_RIFF            = 11
	CF_WAVE            = 12
	CF_UNICODETEXT     = 13
	CF_ENHMETAFILE     = 14
	CF_HDROP           = 15
	CF_LOCALE          = 16
	CF_DIBV5           = 17
	CF_MAX             = 18
	CF_OWNERDISPLAY    = 0x0080
	CF_DSPTEXT         = 0x0081
	CF_DSPBITMAP       = 0x0082
	CF_DSPMETAFILEPICT = 0x0083
	CF_DSPENHMETAFILE  = 0x008E
	CF_PRIVATEFIRST    = 0x0200
	CF_PRIVATELAST     = 0x02FF
	CF_GDIOBJFIRST     = 0x0300
	CF_GDIOBJLAST      = 0x03FF
)

// Bitmap compression formats
const (
	BI_RGB       = 0
	BI_RLE8      = 1
	BI_RLE4      = 2
	BI_BITFIELDS = 3
	BI_JPEG      = 4
	BI_PNG       = 5
)

// SetDIBitsToDevice fuColorUse
const (
	DIB_PAL_COLORS = 1
	DIB_RGB_COLORS = 0
)

// Service Control Manager object specific access types
const (
	SC_MANAGER_CONNECT            = 0x0001
	SC_MANAGER_CREATE_SERVICE     = 0x0002
	SC_MANAGER_ENUMERATE_SERVICE  = 0x0004
	SC_MANAGER_LOCK               = 0x0008
	SC_MANAGER_QUERY_LOCK_STATUS  = 0x0010
	SC_MANAGER_MODIFY_BOOT_CONFIG = 0x0020
	SC_MANAGER_ALL_ACCESS         = STANDARD_RIGHTS_REQUIRED | SC_MANAGER_CONNECT | SC_MANAGER_CREATE_SERVICE | SC_MANAGER_ENUMERATE_SERVICE | SC_MANAGER_LOCK | SC_MANAGER_QUERY_LOCK_STATUS | SC_MANAGER_MODIFY_BOOT_CONFIG
)

// Service Types (Bit Mask)
const (
	SERVICE_KERNEL_DRIVER       = 0x00000001
	SERVICE_FILE_SYSTEM_DRIVER  = 0x00000002
	SERVICE_ADAPTER             = 0x00000004
	SERVICE_RECOGNIZER_DRIVER   = 0x00000008
	SERVICE_DRIVER              = SERVICE_KERNEL_DRIVER | SERVICE_FILE_SYSTEM_DRIVER | SERVICE_RECOGNIZER_DRIVER
	SERVICE_WIN32_OWN_PROCESS   = 0x00000010
	SERVICE_WIN32_SHARE_PROCESS = 0x00000020
	SERVICE_WIN32               = SERVICE_WIN32_OWN_PROCESS | SERVICE_WIN32_SHARE_PROCESS
	SERVICE_INTERACTIVE_PROCESS = 0x00000100
	SERVICE_TYPE_ALL            = 1023
)

// Service State -- for CurrentState
const (
	SERVICE_STOPPED          = 0x00000001
	SERVICE_START_PENDING    = 0x00000002
	SERVICE_STOP_PENDING     = 0x00000003
	SERVICE_RUNNING          = 0x00000004
	SERVICE_CONTINUE_PENDING = 0x00000005
	SERVICE_PAUSE_PENDING    = 0x00000006
	SERVICE_PAUSED           = 0x00000007
)

// Controls Accepted  (Bit Mask)
const (
	SERVICE_ACCEPT_STOP                  = 0x00000001
	SERVICE_ACCEPT_PAUSE_CONTINUE        = 0x00000002
	SERVICE_ACCEPT_SHUTDOWN              = 0x00000004
	SERVICE_ACCEPT_PARAMCHANGE           = 0x00000008
	SERVICE_ACCEPT_NETBINDCHANGE         = 0x00000010
	SERVICE_ACCEPT_HARDWAREPROFILECHANGE = 0x00000020
	SERVICE_ACCEPT_POWEREVENT            = 0x00000040
	SERVICE_ACCEPT_SESSIONCHANGE         = 0x00000080
	SERVICE_ACCEPT_PRESHUTDOWN           = 0x00000100
	SERVICE_ACCEPT_TIMECHANGE            = 0x00000200
	SERVICE_ACCEPT_TRIGGEREVENT          = 0x00000400
)

// Service object specific access type
const (
	SERVICE_QUERY_CONFIG         = 0x0001
	SERVICE_CHANGE_CONFIG        = 0x0002
	SERVICE_QUERY_STATUS         = 0x0004
	SERVICE_ENUMERATE_DEPENDENTS = 0x0008
	SERVICE_START                = 0x0010
	SERVICE_STOP                 = 0x0020
	SERVICE_PAUSE_CONTINUE       = 0x0040
	SERVICE_INTERROGATE          = 0x0080
	SERVICE_USER_DEFINED_CONTROL = 0x0100

	SERVICE_ALL_ACCESS = STANDARD_RIGHTS_REQUIRED |
		SERVICE_QUERY_CONFIG |
		SERVICE_CHANGE_CONFIG |
		SERVICE_QUERY_STATUS |
		SERVICE_ENUMERATE_DEPENDENTS |
		SERVICE_START |
		SERVICE_STOP |
		SERVICE_PAUSE_CONTINUE |
		SERVICE_INTERROGATE |
		SERVICE_USER_DEFINED_CONTROL
)

// MapVirtualKey maptypes
const (
	MAPVK_VK_TO_CHAR   = 2
	MAPVK_VK_TO_VSC    = 0
	MAPVK_VSC_TO_VK    = 1
	MAPVK_VSC_TO_VK_EX = 3
)

// ReadEventLog Flags
const (
	EVENTLOG_SEEK_READ       = 0x0002
	EVENTLOG_SEQUENTIAL_READ = 0x0001
	EVENTLOG_FORWARDS_READ   = 0x0004
	EVENTLOG_BACKWARDS_READ  = 0x0008
)

// CreateToolhelp32Snapshot flags
const (
	TH32CS_SNAPHEAPLIST = 0x00000001
	TH32CS_SNAPPROCESS  = 0x00000002
	TH32CS_SNAPTHREAD   = 0x00000004
	TH32CS_SNAPMODULE   = 0x00000008
	TH32CS_SNAPMODULE32 = 0x00000010
	TH32CS_INHERIT      = 0x80000000
	TH32CS_SNAPALL      = TH32CS_SNAPHEAPLIST | TH32CS_SNAPMODULE | TH32CS_SNAPPROCESS | TH32CS_SNAPTHREAD
)

const (
	MAX_MODULE_NAME32 = 255
	MAX_PATH          = 260
)

const (
	FOREGROUND_BLUE            = 0x0001
	FOREGROUND_GREEN           = 0x0002
	FOREGROUND_RED             = 0x0004
	FOREGROUND_INTENSITY       = 0x0008
	BACKGROUND_BLUE            = 0x0010
	BACKGROUND_GREEN           = 0x0020
	BACKGROUND_RED             = 0x0040
	BACKGROUND_INTENSITY       = 0x0080
	COMMON_LVB_LEADING_BYTE    = 0x0100
	COMMON_LVB_TRAILING_BYTE   = 0x0200
	COMMON_LVB_GRID_HORIZONTAL = 0x0400
	COMMON_LVB_GRID_LVERTICAL  = 0x0800
	COMMON_LVB_GRID_RVERTICAL  = 0x1000
	COMMON_LVB_REVERSE_VIDEO   = 0x4000
	COMMON_LVB_UNDERSCORE      = 0x8000
)

// Flags used by the DWM_BLURBEHIND structure to indicate
// which of its members contain valid information.
const (
	DWM_BB_ENABLE                = 0x00000001 //     A value for the fEnable member has been specified.
	DWM_BB_BLURREGION            = 0x00000002 //     A value for the hRgnBlur member has been specified.
	DWM_BB_TRANSITIONONMAXIMIZED = 0x00000004 //     A value for the fTransitionOnMaximized member has been specified.
)

// Flags used by the DwmEnableComposition  function
// to change the state of Desktop Window Manager (DWM) composition.
const (
	DWM_EC_DISABLECOMPOSITION = 0 //     Disable composition
	DWM_EC_ENABLECOMPOSITION  = 1 //     Enable composition
)

// enum-lite implementation for the following constant structure
type DWM_SHOWCONTACT int32

const (
	DWMSC_DOWN      = 0x00000001
	DWMSC_UP        = 0x00000002
	DWMSC_DRAG      = 0x00000004
	DWMSC_HOLD      = 0x00000008
	DWMSC_PENBARREL = 0x00000010
	DWMSC_NONE      = 0x00000000
	DWMSC_ALL       = 0xFFFFFFFF
)

// enum-lite implementation for the following constant structure
type DWM_SOURCE_FRAME_SAMPLING int32

// Flags used by the DwmSetPresentParameters function to specify the frame
// sampling type
const (
	DWM_SOURCE_FRAME_SAMPLING_POINT    = 0
	DWM_SOURCE_FRAME_SAMPLING_COVERAGE = 1
	DWM_SOURCE_FRAME_SAMPLING_LAST     = 2
)

// Flags used by the DWM_THUMBNAIL_PROPERTIES structure to
// indicate which of its members contain valid information.
const (
	DWM_TNP_RECTDESTINATION      = 0x00000001 //    A value for the rcDestination member has been specified
	DWM_TNP_RECTSOURCE           = 0x00000002 //    A value for the rcSource member has been specified
	DWM_TNP_OPACITY              = 0x00000004 //    A value for the opacity member has been specified
	DWM_TNP_VISIBLE              = 0x00000008 //    A value for the fVisible member has been specified
	DWM_TNP_SOURCECLIENTAREAONLY = 0x00000010 //    A value for the fSourceClientAreaOnly member has been specified
)

// enum-lite implementation for the following constant structure
type DWMFLIP3DWINDOWPOLICY int32

// Flags used by the DwmSetWindowAttribute function to specify the Flip3D window
// policy
const (
	DWMFLIP3D_DEFAULT      = 0
	DWMFLIP3D_EXCLUDEBELOW = 1
	DWMFLIP3D_EXCLUDEABOVE = 2
	DWMFLIP3D_LAST         = 3
)

// enum-lite implementation for the following constant structure
type DWMNCRENDERINGPOLICY int32

// Flags used by the DwmSetWindowAttribute function to specify the non-client
// area rendering policy
const (
	DWMNCRP_USEWINDOWSTYLE = 0
	DWMNCRP_DISABLED       = 1
	DWMNCRP_ENABLED        = 2
	DWMNCRP_LAST           = 3
)

// enum-lite implementation for the following constant structure
type DWMTRANSITION_OWNEDWINDOW_TARGET int32

const (
	DWMTRANSITION_OWNEDWINDOW_NULL       = -1
	DWMTRANSITION_OWNEDWINDOW_REPOSITION = 0
)

// enum-lite implementation for the following constant structure
type DWMWINDOWATTRIBUTE int32

// Flags used by the DwmGetWindowAttribute and DwmSetWindowAttribute functions
// to specify window attributes for non-client rendering
const (
	DWMWA_NCRENDERING_ENABLED         = 1
	DWMWA_NCRENDERING_POLICY          = 2
	DWMWA_TRANSITIONS_FORCEDISABLED   = 3
	DWMWA_ALLOW_NCPAINT               = 4
	DWMWA_CAPTION_BUTTON_BOUNDS       = 5
	DWMWA_NONCLIENT_RTL_LAYOUT        = 6
	DWMWA_FORCE_ICONIC_REPRESENTATION = 7
	DWMWA_FLIP3D_POLICY               = 8
	DWMWA_EXTENDED_FRAME_BOUNDS       = 9
	DWMWA_HAS_ICONIC_BITMAP           = 10
	DWMWA_DISALLOW_PEEK               = 11
	DWMWA_EXCLUDED_FROM_PEEK          = 12
	DWMWA_CLOAK                       = 13
	DWMWA_CLOAKED                     = 14
	DWMWA_FREEZE_REPRESENTATION       = 15
	DWMWA_LAST                        = 16
)

const (
	DWM_CLOAKED_APP       = 1
	DWM_CLOAKED_SHELL     = 2
	DWM_CLOAKED_INHERITED = 4
)

// enum-lite implementation for the following constant structure
type GESTURE_TYPE int32

// Identifies the gesture type
const (
	GT_PEN_TAP                 = 0
	GT_PEN_DOUBLETAP           = 1
	GT_PEN_RIGHTTAP            = 2
	GT_PEN_PRESSANDHOLD        = 3
	GT_PEN_PRESSANDHOLDABORT   = 4
	GT_TOUCH_TAP               = 5
	GT_TOUCH_DOUBLETAP         = 6
	GT_TOUCH_RIGHTTAP          = 7
	GT_TOUCH_PRESSANDHOLD      = 8
	GT_TOUCH_PRESSANDHOLDABORT = 9
	GT_TOUCH_PRESSANDTAP       = 10
)

// Icons
const (
	ICON_SMALL  = 0
	ICON_BIG    = 1
	ICON_SMALL2 = 2
)

const (
	SIZE_RESTORED  = 0
	SIZE_MINIMIZED = 1
	SIZE_MAXIMIZED = 2
	SIZE_MAXSHOW   = 3
	SIZE_MAXHIDE   = 4
)

// XButton values
const (
	XBUTTON1 = 1
	XBUTTON2 = 2
)

// Devmode
const (
	DM_SPECVERSION = 0x0401

	DM_ORIENTATION        = 0x00000001
	DM_PAPERSIZE          = 0x00000002
	DM_PAPERLENGTH        = 0x00000004
	DM_PAPERWIDTH         = 0x00000008
	DM_SCALE              = 0x00000010
	DM_POSITION           = 0x00000020
	DM_NUP                = 0x00000040
	DM_DISPLAYORIENTATION = 0x00000080
	DM_COPIES             = 0x00000100
	DM_DEFAULTSOURCE      = 0x00000200
	DM_PRINTQUALITY       = 0x00000400
	DM_COLOR              = 0x00000800
	DM_DUPLEX             = 0x00001000
	DM_YRESOLUTION        = 0x00002000
	DM_TTOPTION           = 0x00004000
	DM_COLLATE            = 0x00008000
	DM_FORMNAME           = 0x00010000
	DM_LOGPIXELS          = 0x00020000
	DM_BITSPERPEL         = 0x00040000
	DM_PELSWIDTH          = 0x00080000
	DM_PELSHEIGHT         = 0x00100000
	DM_DISPLAYFLAGS       = 0x00200000
	DM_DISPLAYFREQUENCY   = 0x00400000
	DM_ICMMETHOD          = 0x00800000
	DM_ICMINTENT          = 0x01000000
	DM_MEDIATYPE          = 0x02000000
	DM_DITHERTYPE         = 0x04000000
	DM_PANNINGWIDTH       = 0x08000000
	DM_PANNINGHEIGHT      = 0x10000000
	DM_DISPLAYFIXEDOUTPUT = 0x20000000
)

// ChangeDisplaySettings
const (
	CDS_UPDATEREGISTRY  = 0x00000001
	CDS_TEST            = 0x00000002
	CDS_FULLSCREEN      = 0x00000004
	CDS_GLOBAL          = 0x00000008
	CDS_SET_PRIMARY     = 0x00000010
	CDS_VIDEOPARAMETERS = 0x00000020
	CDS_RESET           = 0x40000000
	CDS_NORESET         = 0x10000000

	DISP_CHANGE_SUCCESSFUL  = 0
	DISP_CHANGE_RESTART     = 1
	DISP_CHANGE_FAILED      = -1
	DISP_CHANGE_BADMODE     = -2
	DISP_CHANGE_NOTUPDATED  = -3
	DISP_CHANGE_BADFLAGS    = -4
	DISP_CHANGE_BADPARAM    = -5
	DISP_CHANGE_BADDUALVIEW = -6
)

const (
	ENUM_CURRENT_SETTINGS  = 0xFFFFFFFF
	ENUM_REGISTRY_SETTINGS = 0xFFFFFFFE
)

// PIXELFORMATDESCRIPTOR
const (
	PFD_TYPE_RGBA       = 0
	PFD_TYPE_COLORINDEX = 1

	PFD_MAIN_PLANE     = 0
	PFD_OVERLAY_PLANE  = 1
	PFD_UNDERLAY_PLANE = -1

	PFD_DOUBLEBUFFER         = 0x00000001
	PFD_STEREO               = 0x00000002
	PFD_DRAW_TO_WINDOW       = 0x00000004
	PFD_DRAW_TO_BITMAP       = 0x00000008
	PFD_SUPPORT_GDI          = 0x00000010
	PFD_SUPPORT_OPENGL       = 0x00000020
	PFD_GENERIC_FORMAT       = 0x00000040
	PFD_NEED_PALETTE         = 0x00000080
	PFD_NEED_SYSTEM_PALETTE  = 0x00000100
	PFD_SWAP_EXCHANGE        = 0x00000200
	PFD_SWAP_COPY            = 0x00000400
	PFD_SWAP_LAYER_BUFFERS   = 0x00000800
	PFD_GENERIC_ACCELERATED  = 0x00001000
	PFD_SUPPORT_DIRECTDRAW   = 0x00002000
	PFD_DIRECT3D_ACCELERATED = 0x00004000
	PFD_SUPPORT_COMPOSITION  = 0x00008000

	PFD_DEPTH_DONTCARE        = 0x20000000
	PFD_DOUBLEBUFFER_DONTCARE = 0x40000000
	PFD_STEREO_DONTCARE       = 0x80000000
)

const (
	INPUT_MOUSE    = 0
	INPUT_KEYBOARD = 1
	INPUT_HARDWARE = 2
)

const (
	MOUSEEVENTF_ABSOLUTE        = 0x8000
	MOUSEEVENTF_HWHEEL          = 0x01000
	MOUSEEVENTF_MOVE            = 0x0001
	MOUSEEVENTF_MOVE_NOCOALESCE = 0x2000
	MOUSEEVENTF_LEFTDOWN        = 0x0002
	MOUSEEVENTF_LEFTUP          = 0x0004
	MOUSEEVENTF_RIGHTDOWN       = 0x0008
	MOUSEEVENTF_RIGHTUP         = 0x0010
	MOUSEEVENTF_MIDDLEDOWN      = 0x0020
	MOUSEEVENTF_MIDDLEUP        = 0x0040
	MOUSEEVENTF_VIRTUALDESK     = 0x4000
	MOUSEEVENTF_WHEEL           = 0x0800
	MOUSEEVENTF_XDOWN           = 0x0080
	MOUSEEVENTF_XUP             = 0x0100
)

const (
	KEYEVENTF_EXTENDEDKEY = 0x0001
	KEYEVENTF_KEYUP       = 0x0002
	KEYEVENTF_SCANCODE    = 0x0008
	KEYEVENTF_UNICODE     = 0x0004
)

// Windows Hooks (WH_*)
// http://msdn.microsoft.com/en-us/library/windows/desktop/ms644990(v=vs.85).aspx
const (
	WH_CALLWNDPROC     = 4
	WH_CALLWNDPROCRET  = 12
	WH_CBT             = 5
	WH_DEBUG           = 9
	WH_FOREGROUNDIDLE  = 11
	WH_GETMESSAGE      = 3
	WH_JOURNALPLAYBACK = 1
	WH_JOURNALRECORD   = 0
	WH_KEYBOARD        = 2
	WH_KEYBOARD_LL     = 13
	WH_MOUSE           = 7
	WH_MOUSE_LL        = 14
	WH_MSGFILTER       = -1
	WH_SHELL           = 10
	WH_SYSMSGFILTER    = 6
)

// Process Security and Access Rights
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms684880(v=vs.85).aspx
const (
	PROCESS_TERMINATE                 = 0x0001
	PROCESS_CREATE_THREAD             = 0x0002
	PROCESS_SET_SESSIONID             = 0x0004
	PROCESS_VM_OPERATION              = 0x0008
	PROCESS_VM_READ                   = 0x0010
	PROCESS_VM_WRITE                  = 0x0020
	PROCESS_DUP_HANDLE                = 0x0040
	PROCESS_CREATE_PROCESS            = 0x0080
	PROCESS_SET_QUOTA                 = 0x0100
	PROCESS_SET_INFORMATION           = 0x0200
	PROCESS_QUERY_INFORMATION         = 0x0400
	PROCESS_SUSPEND_RESUME            = 0x0800
	PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
	PROCESS_SET_LIMITED_INFORMATION   = 0x2000
)

const (
	IMAGE_BITMAP = 0
	IMAGE_ICON   = 1
	IMAGE_CURSOR = 2
)

const (
	LR_DEFAULTCOLOR     = 0x0000
	LR_MONOCHROME       = 0x0001
	LR_COLOR            = 0x0002
	LR_COPYRETURNORG    = 0x0004
	LR_COPYDELETEORG    = 0x0008
	LR_LOADFROMFILE     = 0x0010
	LR_LOADTRANSPARENT  = 0x0020
	LR_DEFAULTSIZE      = 0x0040
	LR_VGACOLOR         = 0x0080
	LR_LOADMAP3DCOLORS  = 0x1000
	LR_CREATEDIBSECTION = 0x2000
	LR_COPYFROMRESOURCE = 0x4000
	LR_SHARED           = 0x8000
)

const (
	MF_BITMAP       = 0x0004
	MF_CHECKED      = 0x0008
	MF_DISABLED     = 0x0002
	MF_ENABLED      = 0x0000
	MF_GRAYED       = 0x0001
	MF_MENUBARBREAK = 0x0020
	MF_MENUBREAK    = 0x0040
	MF_OWNERDRAW    = 0x0100
	MF_POPUP        = 0x0010
	MF_SEPARATOR    = 0x0800
	MF_STRING       = 0x0000
	MF_UNCHECKED    = 0x0000
	MF_BYCOMMAND    = 0x0000
	MF_BYPOSITION   = 0x0400
)

const (
	RID_INPUT  = 0x10000003
	RID_HEADER = 0x10000005
)

const (
	RIM_TYPEMOUSE    = 0
	RIM_TYPEKEYBOARD = 1
	RIM_TYPEHID      = 2
)

const (
	RI_KEY_MAKE            = 0 // key is down
	RI_KEY_BREAK           = 1 // key is up
	RI_KEY_E0              = 2 // scan code has e0 prefix
	RI_KEY_E1              = 4 // scan code has e1 prefix
	RI_KEY_TERMSRV_SET_LED = 8
	RI_KEY_TERMSRV_SHADOW  = 0x10
)

const (
	MOUSE_MOVE_RELATIVE      = 0
	MOUSE_MOVE_ABSOLUTE      = 1
	MOUSE_VIRTUAL_DESKTOP    = 0x02
	MOUSE_ATTRIBUTES_CHANGED = 0x04
	MOUSE_MOVE_NOCOALESCE    = 0x08
)

const (
	RI_MOUSE_LEFT_BUTTON_DOWN   = 0x0001
	RI_MOUSE_LEFT_BUTTON_UP     = 0x0002
	RI_MOUSE_RIGHT_BUTTON_DOWN  = 0x0004
	RI_MOUSE_RIGHT_BUTTON_UP    = 0x0008
	RI_MOUSE_MIDDLE_BUTTON_DOWN = 0x0010
	RI_MOUSE_MIDDLE_BUTTON_UP   = 0x0020
	RI_MOUSE_BUTTON_1_DOWN      = RI_MOUSE_LEFT_BUTTON_DOWN
	RI_MOUSE_BUTTON_1_UP        = RI_MOUSE_LEFT_BUTTON_UP
	RI_MOUSE_BUTTON_2_DOWN      = RI_MOUSE_RIGHT_BUTTON_DOWN
	RI_MOUSE_BUTTON_2_UP        = RI_MOUSE_RIGHT_BUTTON_UP
	RI_MOUSE_BUTTON_3_DOWN      = RI_MOUSE_MIDDLE_BUTTON_DOWN
	RI_MOUSE_BUTTON_3_UP        = RI_MOUSE_MIDDLE_BUTTON_UP
	RI_MOUSE_BUTTON_4_DOWN      = 0x0040 // XBUTTON1 changed to down
	RI_MOUSE_BUTTON_4_UP        = 0x0080 // XBUTTON1 changed to up
	RI_MOUSE_BUTTON_5_DOWN      = 0x100  // XBUTTON2 changed to down
	RI_MOUSE_BUTTON_5_UP        = 0x0200 // XBUTTON2 changed to up
	RI_MOUSE_WHEEL              = 0x0400
)

const (
	RIDEV_REMOVE       = 0x00000001
	RIDEV_EXCLUDE      = 0x00000010
	RIDEV_PAGEONLY     = 0x00000020
	RIDEV_NOLEGACY     = 0x00000030
	RIDEV_INPUTSINK    = 0x00000100
	RIDEV_CAPTUREMOUSE = 0x00000200
	RIDEV_NOHOTKEYS    = 0x00000200
	RIDEV_APPKEYS      = 0x00000400
	RIDEV_EXINPUTSINK  = 0x00001000
	RIDEV_DEVNOTIFY    = 0x00002000
)

// GDI+ constants
const (
	Ok                        = 0
	GenericError              = 1
	InvalidParameter          = 2
	OutOfMemory               = 3
	ObjectBusy                = 4
	InsufficientBuffer        = 5
	NotImplemented            = 6
	Win32Error                = 7
	WrongState                = 8
	Aborted                   = 9
	FileNotFound              = 10
	ValueOverflow             = 11
	AccessDenied              = 12
	UnknownImageFormat        = 13
	FontFamilyNotFound        = 14
	FontStyleNotFound         = 15
	NotTrueTypeFont           = 16
	UnsupportedGdiplusVersion = 17
	GdiplusNotInitialized     = 18
	PropertyNotFound          = 19
	PropertyNotSupported      = 20
	ProfileNotFound           = 21
)

var (
	IID_NULL                      = &GUID{0x00000000, 0x0000, 0x0000, [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}
	IID_IUnknown                  = &GUID{0x00000000, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	IID_IDispatch                 = &GUID{0x00020400, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	IID_IConnectionPointContainer = &GUID{0xB196B284, 0xBAB4, 0x101A, [8]byte{0xB6, 0x9C, 0x00, 0xAA, 0x00, 0x34, 0x1D, 0x07}}
	IID_IConnectionPoint          = &GUID{0xB196B286, 0xBAB4, 0x101A, [8]byte{0xB6, 0x9C, 0x00, 0xAA, 0x00, 0x34, 0x1D, 0x07}}
)

const (
	VS_FF_DEBUG        = 0x00000001
	VS_FF_INFOINFERRED = 0x00000010
	VS_FF_PATCHED      = 0x00000004
	VS_FF_PRERELEASE   = 0x00000002
	VS_FF_PRIVATEBUILD = 0x00000008
	VS_FF_SPECIALBUILD = 0x00000020

	VOS_DOS        = 0x00010000
	VOS_NT         = 0x00040000
	VOS__WINDOWS16 = 0x00000001
	VOS__WINDOWS32 = 0x00000004
	VOS_OS216      = 0x00020000
	VOS_OS232      = 0x00030000
	VOS__PM16      = 0x00000002
	VOS__PM32      = 0x00000003
	VOS_UNKNOWN    = 0x00000000

	VOS_DOS_WINDOWS16 = 0x00010001
	VOS_DOS_WINDOWS32 = 0x00010004
	VOS_NT_WINDOWS32  = 0x00040004
	VOS_OS216_PM16    = 0x00020002
	VOS_OS232_PM32    = 0x00030003

	VFT_APP        = 0x00000001
	VFT_DLL        = 0x00000002
	VFT_DRV        = 0x00000003
	VFT_FONT       = 0x00000004
	VFT_STATIC_LIB = 0x00000007
	VFT_UNKNOWN    = 0x00000000
	VFT_VXD        = 0x00000005

	VFT2_DRV_COMM              = 0x0000000A
	VFT2_DRV_DISPLAY           = 0x00000004
	VFT2_DRV_INSTALLABLE       = 0x00000008
	VFT2_DRV_KEYBOARD          = 0x00000002
	VFT2_DRV_LANGUAGE          = 0x00000003
	VFT2_DRV_MOUSE             = 0x00000005
	VFT2_DRV_NETWORK           = 0x00000006
	VFT2_DRV_PRINTER           = 0x00000001
	VFT2_DRV_SOUND             = 0x00000009
	VFT2_DRV_SYSTEM            = 0x00000007
	VFT2_DRV_VERSIONED_PRINTER = 0x0000000C
	VFT2_UNKNOWN               = 0x00000000

	VFT2_FONT_RASTER   = 0x00000001
	VFT2_FONT_TRUETYPE = 0x00000003
	VFT2_FONT_VECTOR   = 0x00000002
)

const (
	SND_SYNC        = 0
	SND_ASYNC       = 1
	SND_NODEFAULT   = 2
	SND_MEMORY      = 4
	SND_LOOP        = 8
	SND_NOSTOP      = 16
	SND_NOWAIT      = 0x2000
	SND_ALIAS       = 0x10000
	SND_ALIAS_ID    = 0x110000
	SND_FILENAME    = 0x20000
	SND_RESOURCE    = 0x40004
	SND_PURGE       = 0x40
	SND_APPLICATION = 0x80
	SND_ALIAS_START = 0
)

const (
	BLACKONWHITE        = 1
	WHITEONBLACK        = 2
	COLORONCOLOR        = 3
	HALFTONE            = 4
	STRETCH_ANDSCANS    = BLACKONWHITE
	STRETCH_DELETESCANS = COLORONCOLOR
	STRETCH_HALFTONE    = HALFTONE
	STRETCH_ORSCANS     = WHITEONBLACK
)

const (
	GW_HWNDFIRST    = 0
	GW_HWNDLAST     = 1
	GW_HWNDNEXT     = 2
	GW_HWNDPREV     = 3
	GW_OWNER        = 4
	GW_CHILD        = 5
	GW_ENABLEDPOPUP = 6
)

// ACCEL behavior flags
const (
	FVIRTKEY  = 0x01
	FNOINVERT = 0x02
	FSHIFT    = 0x04
	FCONTROL  = 0x08
	FALT      = 0x10
)

// Mouse key flags
const (
	MK_LBUTTON  = 0x0001
	MK_RBUTTON  = 0x0002
	MK_SHIFT    = 0x0004
	MK_CONTROL  = 0x0008
	MK_MBUTTON  = 0x0010
	MK_XBUTTON1 = 0x0020
	MK_XBUTTON2 = 0x0040
)

const UPDOWN_CLASS = "msctls_updown32"

const (
	UDS_WRAP        = 1
	UDS_SETBUDDYINT = 2
	UDS_ALIGNRIGHT  = 4
	UDS_ALIGNLEFT   = 8
	UDS_AUTOBUDDY   = 16
	UDS_ARROWKEYS   = 32
	UDS_HORZ        = 64
	UDS_NOTHOUSANDS = 128
)

const (
	UDN_FIRST    = 4294966575
	UDN_LAST     = 4294966567
	UDN_DELTAPOS = UDN_FIRST - 1
)

const (
	UDM_SETRANGE         = WM_USER + 101
	UDM_GETRANGE         = WM_USER + 102
	UDM_SETPOS           = WM_USER + 103
	UDM_GETPOS           = WM_USER + 104
	UDM_SETBUDDY         = WM_USER + 105
	UDM_GETBUDDY         = WM_USER + 106
	UDM_SETACCEL         = WM_USER + 107
	UDM_GETACCEL         = WM_USER + 108
	UDM_SETBASE          = WM_USER + 109
	UDM_GETBASE          = WM_USER + 110
	UDM_SETRANGE32       = WM_USER + 111
	UDM_GETRANGE32       = WM_USER + 112
	UDM_SETPOS32         = WM_USER + 113
	UDM_GETPOS32         = WM_USER + 114
	UDM_SETUNICODEFORMAT = 0x2005
	UDM_GETUNICODEFORMAT = 0x2006
)

const (
	GCL_CBCLSEXTRA     = -20
	GCL_CBWNDEXTRA     = -18
	GCLP_HBRBACKGROUND = -10
	GCLP_HCURSOR       = -12
	GCLP_HICON         = -14
	GCLP_HICONSM       = -34
	GCLP_HMODULE       = -16
	GCLP_MENUNAME      = -8
	GCL_STYLE          = -26
	GCLP_WNDPROC       = -24
)

// system commands
const (
	SC_CLOSE        = 0xF060
	SC_CONTEXTHELP  = 0xF180
	SC_DEFAULT      = 0xF160
	SC_HOTKEY       = 0xF150
	SC_HSCROLL      = 0xF080
	SCF_ISSECURE    = 0x0001
	SC_KEYMENU      = 0xF100
	SC_MAXIMIZE     = 0xF030
	SC_MINIMIZE     = 0xF020
	SC_MONITORPOWER = 0xF170
	SC_MOUSEMENU    = 0xF090
	SC_MOVE         = 0xF010
	SC_NEXTWINDOW   = 0xF040
	SC_PREVWINDOW   = 0xF050
	SC_RESTORE      = 0xF120
	SC_SCREENSAVE   = 0xF140
	SC_SIZE         = 0xF000
	SC_TASKLIST     = 0xF130
	SC_VSCROLL      = 0xF070
)

const AC_SRC_OVER = 0

const AC_SRC_ALPHA = 1

const (
	RESOURCE_CONNECTED  = 0x00000001
	RESOURCE_GLOBALNET  = 0x00000002
	RESOURCE_REMEMBERED = 0x00000003
	RESOURCE_RECENT     = 0x00000004
	RESOURCE_CONTEXT    = 0x00000005
)

const (
	RESOURCETYPE_ANY      = 0x00000000
	RESOURCETYPE_DISK     = 0x00000001
	RESOURCETYPE_PRINT    = 0x00000002
	RESOURCETYPE_RESERVED = 0x00000008
	RESOURCETYPE_UNKNOWN  = 0xFFFFFFFF
)

const (
	RESOURCEUSAGE_CONNECTABLE   = 0x00000001
	RESOURCEUSAGE_CONTAINER     = 0x00000002
	RESOURCEUSAGE_NOLOCALDEVICE = 0x00000004
	RESOURCEUSAGE_SIBLING       = 0x00000008
	RESOURCEUSAGE_ATTACHED      = 0x00000010
	RESOURCEUSAGE_ALL           = RESOURCEUSAGE_CONNECTABLE | RESOURCEUSAGE_CONTAINER | RESOURCEUSAGE_ATTACHED
	RESOURCEUSAGE_RESERVED      = 0x80000000
)

const (
	RESOURCEDISPLAYTYPE_GENERIC      = 0x00000000
	RESOURCEDISPLAYTYPE_DOMAIN       = 0x00000001
	RESOURCEDISPLAYTYPE_SERVER       = 0x00000002
	RESOURCEDISPLAYTYPE_SHARE        = 0x00000003
	RESOURCEDISPLAYTYPE_FILE         = 0x00000004
	RESOURCEDISPLAYTYPE_GROUP        = 0x00000005
	RESOURCEDISPLAYTYPE_NETWORK      = 0x00000006
	RESOURCEDISPLAYTYPE_ROOT         = 0x00000007
	RESOURCEDISPLAYTYPE_SHAREADMIN   = 0x00000008
	RESOURCEDISPLAYTYPE_DIRECTORY    = 0x00000009
	RESOURCEDISPLAYTYPE_TREE         = 0x0000000A
	RESOURCEDISPLAYTYPE_NDSCONTAINER = 0x0000000B
)

const NETPROPERTY_PERSISTENT = 1

const (
	CONNECT_UPDATE_PROFILE = 0x00000001
	CONNECT_UPDATE_RECENT  = 0x00000002
	CONNECT_TEMPORARY      = 0x00000004
	CONNECT_INTERACTIVE    = 0x00000008
	CONNECT_PROMPT         = 0x00000010
	CONNECT_NEED_DRIVE     = 0x00000020
	CONNECT_REFCOUNT       = 0x00000040
	CONNECT_REDIRECT       = 0x00000080
	CONNECT_LOCALDRIVE     = 0x00000100
	CONNECT_CURRENT_MEDIA  = 0x00000200
	CONNECT_DEFERRED       = 0x00000400
	CONNECT_RESERVED       = 0xFF000000
	CONNECT_COMMANDLINE    = 0x00000800
	CONNECT_CMD_SAVECRED   = 0x00001000
)

const PW_CLIENTONLY = 1

const (
	SEM_FAILCRITICALERRORS     = 0x0001
	SEM_NOALIGNMENTFAULTEXCEPT = 0x0004
	SEM_NOGPFAULTERRORBOX      = 0x0002
	SEM_NOOPENFILEERRORBOX     = 0x8000
)

const (
	GENERIC_ALL     = 0x10000000
	GENERIC_EXECUTE = 0x20000000
	GENERIC_WRITE   = 0x40000000
	GENERIC_READ    = 0x80000000
)

const SYNCHRONIZE = 0x100000

const (
	STANDARD_RIGHTS_REQUIRED = 0xF0000
	STANDARD_RIGHTS_READ     = 0x20000
	STANDARD_RIGHTS_WRITE    = 0x20000
	STANDARD_RIGHTS_EXECUTE  = 0x20000
	STANDARD_RIGHTS_ALL      = 0x1F0000
	SPECIFIC_RIGHTS_ALL      = 0xFFFF
	ACCESS_SYSTEM_SECURITY   = 0x1000000
)

const (
	FILE_READ_DATA                     = 1
	FILE_LIST_DIRECTORY                = 1
	FILE_WRITE_DATA                    = 2
	FILE_ADD_FILE                      = 2
	FILE_APPEND_DATA                   = 4
	FILE_ADD_SUBDIRECTORY              = 4
	FILE_CREATE_PIPE_INSTANCE          = 4
	FILE_READ_EA                       = 8
	FILE_READ_PROPERTIES               = 8
	FILE_WRITE_EA                      = 16
	FILE_WRITE_PROPERTIES              = 16
	FILE_EXECUTE                       = 32
	FILE_TRAVERSE                      = 32
	FILE_DELETE_CHILD                  = 64
	FILE_READ_ATTRIBUTES               = 128
	FILE_WRITE_ATTRIBUTES              = 256
	FILE_ALL_ACCESS                    = STANDARD_RIGHTS_REQUIRED | SYNCHRONIZE | 0x1FF
	FILE_GENERIC_READ                  = STANDARD_RIGHTS_READ | FILE_READ_DATA | FILE_READ_ATTRIBUTES | FILE_READ_EA | SYNCHRONIZE
	FILE_GENERIC_WRITE                 = STANDARD_RIGHTS_WRITE | FILE_WRITE_DATA | FILE_WRITE_ATTRIBUTES | FILE_WRITE_EA | FILE_APPEND_DATA | SYNCHRONIZE
	FILE_GENERIC_EXECUTE               = STANDARD_RIGHTS_EXECUTE | FILE_READ_ATTRIBUTES | FILE_EXECUTE | SYNCHRONIZE
	FILE_SHARE_READ                    = 1
	FILE_SHARE_WRITE                   = 2
	FILE_SHARE_DELETE                  = 4
	FILE_ATTRIBUTE_READONLY            = 1
	FILE_ATTRIBUTE_HIDDEN              = 2
	FILE_ATTRIBUTE_SYSTEM              = 4
	FILE_ATTRIBUTE_DIRECTORY           = 16
	FILE_ATTRIBUTE_ARCHIVE             = 32
	FILE_ATTRIBUTE_ENCRYPTED           = 16384
	FILE_ATTRIBUTE_NORMAL              = 128
	FILE_ATTRIBUTE_TEMPORARY           = 256
	FILE_ATTRIBUTE_SPARSE_FILE         = 512
	FILE_ATTRIBUTE_REPARSE_POINT       = 1024
	FILE_ATTRIBUTE_COMPRESSED          = 2048
	FILE_ATTRIBUTE_OFFLINE             = 0x1000
	FILE_ATTRIBUTE_NOT_CONTENT_INDEXED = 0x2000
	FILE_NOTIFY_CHANGE_FILE_NAME       = 1
	FILE_NOTIFY_CHANGE_DIR_NAME        = 2
	FILE_NOTIFY_CHANGE_ATTRIBUTES      = 4
	FILE_NOTIFY_CHANGE_SIZE            = 8
	FILE_NOTIFY_CHANGE_LAST_WRITE      = 16
	FILE_NOTIFY_CHANGE_LAST_ACCESS     = 32
	FILE_NOTIFY_CHANGE_CREATION        = 64
	FILE_NOTIFY_CHANGE_SECURITY        = 256
	FILE_CASE_SENSITIVE_SEARCH         = 1
	FILE_CASE_PRESERVED_NAMES          = 2
	FILE_UNICODE_ON_DISK               = 4
	FILE_PERSISTENT_ACLS               = 8
	FILE_FILE_COMPRESSION              = 16
	FILE_VOLUME_QUOTAS                 = 32
	FILE_SUPPORTS_SPARSE_FILES         = 64
	FILE_SUPPORTS_REPARSE_POINTS       = 128
)

const (
	CREATE_NEW        = 1
	CREATE_ALWAYS     = 2
	OPEN_EXISTING     = 3
	OPEN_ALWAYS       = 4
	TRUNCATE_EXISTING = 5
)

const (
	SECURITY_CONTEXT_TRACKING = 0x00040000
	SECURITY_EFFECTIVE_ONLY   = 0x00080000
	SECURITY_SQOS_PRESENT     = 0x00100000
	SECURITY_VALID_SQOS_FLAGS = 0x001F0000
)

const INVALID_HANDLE_VALUE = ^HANDLE(0)

// STORAGE_BUS_TYPE values
const (
	BusTypeUnknown           = 0
	BusTypeScsi              = 1
	BusTypeAtapi             = 2
	BusTypeAta               = 3
	BusType1394              = 4
	BusTypeSsa               = 5
	BusTypeFibre             = 6
	BusTypeUsb               = 7
	BusTypeRAID              = 8
	BusTypeiScsi             = 9
	BusTypeSas               = 10
	BusTypeSata              = 11
	BusTypeSd                = 12
	BusTypeMmc               = 13
	BusTypeVirtual           = 14
	BusTypeFileBackedVirtual = 15
	BusTypeSpaces            = 16
	BusTypeNvme              = 17
	BusTypeMax               = 20
	BusTypeMaxReserved       = 0x7F
)

// STORAGE_QUERY_TYPE values
const (
	PropertyStandardQuery   = 0
	PropertyExistsQuery     = 1
	PropertyMaskQuery       = 2
	PropertyQueryMaxDefined = 3
)

// STORAGE_PROPERTY_ID values
const (
	StorageDeviceProperty                 = 0
	StorageAdapterProperty                = 1
	StorageDeviceIdProperty               = 2
	StorageDeviceUniqueIdProperty         = 3
	StorageDeviceWriteCacheProperty       = 4
	StorageMiniportProperty               = 5
	StorageAccessAlignmentProperty        = 6
	StorageDeviceSeekPenaltyProperty      = 7
	StorageDeviceTrimProperty             = 8
	StorageDeviceWriteAggregationProperty = 9
	StorageDeviceDeviceTelemetryProperty  = 10
	StorageDeviceLBProvisioningProperty   = 11
	StorageDevicePowerProperty            = 12
	StorageDeviceCopyOffloadProperty      = 13
	StorageDeviceResiliencyProperty       = 14
	StorageDeviceMediumProductType        = 15
)

// STREAM_INFO_LEVELS
const (
	FindStreamInfoStandard = 0
)

const (
	MIIM_STATE      = 1
	MIIM_ID         = 2
	MIIM_SUBMENU    = 4
	MIIM_CHECKMARKS = 8
	MIIM_TYPE       = 16
	MIIM_DATA       = 32
	MIIM_STRING     = 64
	MIIM_BITMAP     = 128
	MIIM_FTYPE      = 256
)

const (
	MFT_BITMAP       = 4
	MFT_MENUBARBREAK = 32
	MFT_MENUBREAK    = 64
	MFT_OWNERDRAW    = 256
	MFT_RADIOCHECK   = 512
	MFT_RIGHTJUSTIFY = 0x4000
	MFT_SEPARATOR    = 0x800
	MFT_RIGHTORDER   = 0x2000
	MFT_STRING       = 0
)

const (
	MFS_CHECKED   = 8
	MFS_DEFAULT   = 4096
	MFS_DISABLED  = 3
	MFS_ENABLED   = 0
	MFS_GRAYED    = 3
	MFS_HILITE    = 128
	MFS_UNCHECKED = 0
	MFS_UNHILITE  = 0
)

const (
	// DEVICE_NOTIFY_WINDOW_HANDLE
	// Notifications are sent using WM_POWERBROADCAST messages with a wParam
	// parameter of PBT_POWERSETTINGCHANGE.
	DEVICE_NOTIFY_WINDOW_HANDLE = 0

	// DEVICE_NOTIFY_SERVICE_HANDLE
	// Notifications are sent to the HandlerEx callback function with a
	// dwControl parameter of SERVICE_CONTROL_POWEREVENT and a dwEventType of
	// PBT_POWERSETTINGCHANGE.
	DEVICE_NOTIFY_SERVICE_HANDLE = 1
)

// Power setting GUIDs identify power change events. This topic lists power
// setting GUIDs for notifications that are most useful to applications. An
// application should register for each power change event that might impact its
// behavior. Notification is sent each time a setting changes, through
// WM_POWERBROADCAST.
var (
	// GUID_ACDC_POWER_SOURCE
	//
	// The system power source has changed. The Data member is a uint32 with
	// values from the SYSTEM_POWER_CONDITION enumeration that indicates the
	// current power source.
	//
	// PoAc (0) - The computer is powered by an AC power source (or similar,
	// such as a laptop powered by a 12V automotive adapter).
	//
	// PoDc (1) - The computer is powered by an onboard battery power source.
	//
	// PoHot (2) - The computer is powered by a short-term power source such as
	// a UPS device.
	GUID_ACDC_POWER_SOURCE = &GUID{0x5d3e9a59, 0xe9D5, 0x4b00, [8]byte{0xa6, 0xbd, 0xff, 0x34, 0xff, 0x51, 0x65, 0x48}}

	// GUID_BATTERY_PERCENTAGE_REMAINING
	//
	// The remaining battery capacity has changed. The granularity varies from
	// system to system but the finest granularity is 1 percent. The Data member
	// is a uint32 that indicates the current battery capacity remaining as a
	// percentage from 0 through 100.
	GUID_BATTERY_PERCENTAGE_REMAINING = &GUID{0xa7ad8041, 0xb45a, 0x4cae, [8]byte{0x87, 0xa3, 0xee, 0xcb, 0xb4, 0x68, 0xa9, 0xe1}}

	// GUID_CONSOLE_DISPLAY_STATE
	//
	// The current monitor's display state has changed.
	//
	// Windows 7, Windows Server 2008 R2, Windows Vista and Windows Server 2008:
	// This notification is available starting with Windows 8 and Windows Server
	// 2012.
	//
	// The Data member is a uint32 with one of the following values.
	//
	// 0x0 - The display is off.
	//
	// 0x1 - The display is on.
	//
	// 0x2 - The display is dimmed.
	GUID_CONSOLE_DISPLAY_STATE = &GUID{0x6fe69556, 0x704a, 0x47a0, [8]byte{0x8f, 0x24, 0xc2, 0x8d, 0x93, 0x6f, 0xda, 0x47}}

	// GUID_GLOBAL_USER_PRESENCE
	//
	// The user status associated with any session has changed. This represents
	// the combined status of user presence across all local and remote sessions
	// on the system.
	//
	// This notification is sent only to services and other programs running in
	// session 0. User-mode applications should register for
	// GUID_SESSION_USER_PRESENCE instead.
	//
	// Windows 7, Windows Server 2008 R2, Windows Vista and Windows Server 2008:
	// This notification is available starting with Windows 8 and Windows Server
	// 2012.
	//
	// The Data member is a DWORD with one of the following values.
	//
	// PowerUserPresent (0) - The user is present in any local or remote session
	// on the system.
	//
	// PowerUserInactive (2) - The user is not present in any local or remote
	// session on the system.
	GUID_GLOBAL_USER_PRESENCE = &GUID{0x786E8A1D, 0xB427, 0x4344, [8]byte{0x92, 0x07, 0x09, 0xE7, 0x0B, 0xDC, 0xBE, 0xA9}}

	// GUID_IDLE_BACKGROUND_TASK
	//
	// The system is busy. This indicates that the system will not be moving
	// into an idle state in the near future and that the current time is a good
	// time for components to perform background or idle tasks that would
	// otherwise prevent the computer from entering an idle state.
	//
	// There is no notification when the system is able to move into an idle
	// state. The idle background task notification does not indicate whether a
	// user is present at the computer. The Data member has no information and
	// can be ignored.
	GUID_IDLE_BACKGROUND_TASK = &GUID{0x515c31d8, 0xf734, 0x163d, [8]byte{0xa0, 0xfd, 0x11, 0xa0, 0x8c, 0x91, 0xe8, 0xf1}}

	// GUID_MONITOR_POWER_ON
	//
	// The primary system monitor has been powered on or off. This notification
	// is useful for components that actively render content to the display
	// device, such as media visualization. These applications should register
	// for this notification and stop rendering graphics content when the
	// monitor is off to reduce system power consumption. The Data member is a
	// DWORD that indicates the current monitor state.
	//
	// 0x0 - The monitor is off.
	//
	// 0x1 - The monitor is on.
	//
	// Windows 8 and Windows Server 2012: New applications should use
	// GUID_CONSOLE_DISPLAY_STATE instead of this notification.
	GUID_MONITOR_POWER_ON = &GUID{0x02731015, 0x4510, 0x4526, [8]byte{0x99, 0xe6, 0xe5, 0xa1, 0x7e, 0xbd, 0x1a, 0xea}}

	// GUID_POWER_SAVING_STATUS
	//
	// Battery saver has been turned off or on in response to changing power
	// conditions. This notification is useful for components that participate
	// in energy conservation. These applications should register for this
	// notification and save power when battery saver is on.
	//
	// The Data member is a DWORD that indicates battery saver state.
	//
	// 0x0 - Battery saver is off.
	//
	// 0x1 - Battery saver is on. Save energy where possible.
	//
	// For general information about battery saver, see battery saver (in the
	// hardware component guidelines).
	GUID_POWER_SAVING_STATUS = &GUID{0xE00958C0, 0xC213, 0x4ACE, [8]byte{0xAC, 0x77, 0xFE, 0xCC, 0xED, 0x2E, 0xEE, 0xA5}}

	// GUID_POWERSCHEME_PERSONALITY
	//
	// The active power scheme personality has changed. All power schemes map to
	// one of these personalities. The Data member is a GUID that indicates the
	// new active power scheme personality (GUID_MIN_POWER_SAVINGS,
	// GUID_MAX_POWER_SAVINGS or GUID_TYPICAL_POWER_SAVINGS).
	GUID_POWERSCHEME_PERSONALITY = &GUID{0x245d8541, 0x3943, 0x4422, [8]byte{0xb0, 0x25, 0x13, 0xA7, 0x84, 0xF6, 0x79, 0xB7}}

	// GUID_MIN_POWER_SAVINGS
	//
	// High Performance - The scheme is designed to deliver maximum performance
	// at the expense of power consumption savings.
	GUID_MIN_POWER_SAVINGS = &GUID{0x8c5e7fda, 0xe8bf, 0x4a96, [8]byte{0x9a, 0x85, 0xa6, 0xe2, 0x3a, 0x8c, 0x63, 0x5c}}

	// GUID_MAX_POWER_SAVINGS
	//
	// Power Saver - The scheme is designed to deliver maximum power consumption
	// savings at the expense of system performance and responsiveness.
	GUID_MAX_POWER_SAVINGS = &GUID{0xa1841308, 0x3541, 0x4fab, [8]byte{0xbc, 0x81, 0xf7, 0x15, 0x56, 0xf2, 0x0b, 0x4a}}

	// GUID_TYPICAL_POWER_SAVINGS
	//
	// Automatic - The scheme is designed to automatically balance performance
	// and power consumption savings.
	GUID_TYPICAL_POWER_SAVINGS = &GUID{0x381b4222, 0xf694, 0x41f0, [8]byte{0x96, 0x85, 0xff, 0x5b, 0xb2, 0x60, 0xdf, 0x2e}}

	// GUID_SESSION_DISPLAY_STATUS
	//
	// The display associated with the application's session has been powered on
	// or off.
	//
	// Windows 7, Windows Server 2008 R2, Windows Vista and Windows Server 2008:
	// This notification is available starting with Windows 8 and Windows Server
	// 2012.
	//
	// This notification is sent only to user-mode applications. Services and
	// other programs running in session 0 do not receive this notification. The
	// Data member is a DWORD with one of the following values.
	//
	// 0x0 - The display is off.
	//
	// 0x1 - The display is on.
	//
	// 0x2 - The display is dimmed.
	GUID_SESSION_DISPLAY_STATUS = &GUID{0x2B84C20E, 0xAD23, 0x4ddf, [8]byte{0x93, 0xDB, 0x05, 0xFF, 0xBD, 0x7E, 0xFC, 0xA5}}

	// GUID_SESSION_USER_PRESENCE
	//
	// The user status associated with the application's session has changed.
	//
	// Windows 7, Windows Server 2008 R2, Windows Vista and Windows Server 2008:
	// This notification is available starting with Windows 8 and Windows Server
	// 2012.
	//
	// This notification is sent only to user-mode applications running in an
	// interactive session. Services and other programs running in session 0
	// should register for GUID_GLOBAL_USER_PRESENCE. The Data member is a DWORD
	// with one of the following values.
	//
	// PowerUserPresent (0) - The user is providing input to the session.
	//
	// PowerUserInactive (2) - The user activity timeout has elapsed with no
	// interaction from the user.
	//
	// NOTE All applications that run in an interactive user-mode session should
	// use this setting. When kernel mode applications register for monitor
	// status they should use GUID_CONSOLE_DISPLAY_STATUS instead.
	GUID_SESSION_USER_PRESENCE = &GUID{0x3C0F4548, 0xC03F, 0x4c4d, [8]byte{0xB9, 0xF2, 0x23, 0x7E, 0xDE, 0x68, 0x63, 0x76}}

	// GUID_SYSTEM_AWAYMODE
	//
	// The system is entering or exiting away mode. The Data member is a DWORD
	// that indicates the current away mode state.
	//
	// 0x0 - The computer is exiting away mode.
	//
	// 0x1 - The computer is entering away mode.
	GUID_SYSTEM_AWAYMODE = &GUID{0x98a7f580, 0x01f7, 0x48aa, [8]byte{0x9c, 0x0f, 0x44, 0x35, 0x2c, 0x29, 0xe5, 0xC0}}

	// GUID_DEVINTERFACE_COMPORT is defined for COM ports.
	GUID_DEVINTERFACE_COMPORT = &GUID{0x86E0D1E0, 0x8089, 0x11D0, [8]byte{0x9C, 0xE4, 0x08, 0x00, 0x3E, 0x30, 0x1F, 0x73}}
)

const (
	// Power status has changed.
	PBT_APMPOWERSTATUSCHANGE = 0x000A

	// Operation is resuming automatically from a low-power state. This message
	// is sent every time the system resumes.
	PBT_APMRESUMEAUTOMATIC = 0x0012

	// Operation is resuming from a low-power state. This message is sent after
	// PBT_APMRESUMEAUTOMATIC if the resume is triggered by user input, such as
	// pressing a key.
	PBT_APMRESUMESUSPEND = 0x0007

	// System is suspending operation.
	PBT_APMSUSPEND = 0x0004

	// A power setting change event has been received.
	PBT_POWERSETTINGCHANGE = 0x8013
)

// SYSTEM_POWER_CONDITION enumeration
const (
	// PoAc: the computer is powered by an AC power source (or similar, such as
	// a laptop powered by a 12V automotive adapter).
	PoAc = 0
	// PoDc: the system is receiving power from built-in batteries.
	PoDc = 1
	// PoHot: the computer is powered by a short-term power source such as a UPS
	// device.
	PoHot = 2
	// PoConditionMaximum: values equal to or greater than this value indicate
	// an out of range value.
	PoConditionMaximum = 3
)

const (
	CSIDL_DESKTOP                 = 0x0000         // <desktop>
	CSIDL_INTERNET                = 0x0001         // Internet Explorer (icon on desktop)
	CSIDL_PROGRAMS                = 0x0002         // Start Menu\Programs
	CSIDL_CONTROLS                = 0x0003         // My Computer\Control Panel
	CSIDL_PRINTERS                = 0x0004         // My Computer\Printers
	CSIDL_PERSONAL                = 0x0005         // My Documents
	CSIDL_FAVORITES               = 0x0006         // <user name>\Favorites
	CSIDL_STARTUP                 = 0x0007         // Start Menu\Programs\Startup
	CSIDL_RECENT                  = 0x0008         // <user name>\Recent
	CSIDL_SENDTO                  = 0x0009         // <user name>\SendTo
	CSIDL_BITBUCKET               = 0x000A         // <desktop>\Recycle Bin
	CSIDL_STARTMENU               = 0x000B         // <user name>\Start Menu
	CSIDL_MYDOCUMENTS             = CSIDL_PERSONAL //  Personal was just a silly name for My Documents
	CSIDL_MYMUSIC                 = 0x000D         // "My Music" folder
	CSIDL_MYVIDEO                 = 0x000E         // "My Videos" folder
	CSIDL_DESKTOPDIRECTORY        = 0x0010         // <user name>\Desktop
	CSIDL_DRIVES                  = 0x0011         // My Computer
	CSIDL_NETWORK                 = 0x0012         // Network Neighborhood (My Network Places)
	CSIDL_NETHOOD                 = 0x0013         // <user name>\nethood
	CSIDL_FONTS                   = 0x0014         // windows\fonts
	CSIDL_TEMPLATES               = 0x0015
	CSIDL_COMMON_STARTMENU        = 0x0016 // All Users\Start Menu
	CSIDL_COMMON_PROGRAMS         = 0x0017 // All Users\Start Menu\Programs
	CSIDL_COMMON_STARTUP          = 0x0018 // All Users\Startup
	CSIDL_COMMON_DESKTOPDIRECTORY = 0x0019 // All Users\Desktop
	CSIDL_APPDATA                 = 0x001A // <user name>\Application Data
	CSIDL_PRINTHOOD               = 0x001B // <user name>\PrintHood
	CSIDL_LOCAL_APPDATA           = 0x001C // <user name>\Local Settings\Applicaiton Data (non roaming)
	CSIDL_ALTSTARTUP              = 0x001D // non localized startup
	CSIDL_COMMON_ALTSTARTUP       = 0x001E // non localized common startup
	CSIDL_COMMON_FAVORITES        = 0x001F
	CSIDL_INTERNET_CACHE          = 0x0020
	CSIDL_COOKIES                 = 0x0021
	CSIDL_HISTORY                 = 0x0022
	CSIDL_COMMON_APPDATA          = 0x0023 // All Users\Application Data
	CSIDL_WINDOWS                 = 0x0024 // GetWindowsDirectory()
	CSIDL_SYSTEM                  = 0x0025 // GetSystemDirectory()
	CSIDL_PROGRAM_FILES           = 0x0026 // C:\Program Files
	CSIDL_MYPICTURES              = 0x0027 // C:\Program Files\My Pictures
	CSIDL_PROFILE                 = 0x0028 // USERPROFILE
	CSIDL_SYSTEMX86               = 0x0029 // x86 system directory on RISC
	CSIDL_PROGRAM_FILESX86        = 0x002A // x86 C:\Program Files on RISC
	CSIDL_PROGRAM_FILES_COMMON    = 0x002B // C:\Program Files\Common
	CSIDL_PROGRAM_FILES_COMMONX86 = 0x002C // x86 Program Files\Common on RISC
	CSIDL_COMMON_TEMPLATES        = 0x002D // All Users\Templates
	CSIDL_COMMON_DOCUMENTS        = 0x002E // All Users\Documents
	CSIDL_COMMON_ADMINTOOLS       = 0x002F // All Users\Start Menu\Programs\Administrative Tools
	CSIDL_ADMINTOOLS              = 0x0030 // <user name>\Start Menu\Programs\Administrative Tools
	CSIDL_CONNECTIONS             = 0x0031 // Network and Dial-up Connections
	CSIDL_COMMON_MUSIC            = 0x0035 // All Users\My Music
	CSIDL_COMMON_PICTURES         = 0x0036 // All Users\My Pictures
	CSIDL_COMMON_VIDEO            = 0x0037 // All Users\My Video
	CSIDL_RESOURCES               = 0x0038 // Resource Direcotry
)

// Drive types to use with GetDriveType.
const (
	DRIVE_UNKNOWN     = 0
	DRIVE_NO_ROOT_DIR = 1
	DRIVE_REMOVABLE   = 2
	DRIVE_FIXED       = 3
	DRIVE_REMOTE      = 4
	DRIVE_CDROM       = 5
	DRIVE_RAMDISK     = 6
)

// https://docs.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-osversioninfoexw
const (
	VER_PLATFORM_WIN32_NT = 2
)

// https://docs.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-osversioninfoexw
const (
	VER_NT_SERVER            = 0x0000003
	VER_NT_DOMAIN_CONTROLLER = 0x0000002
	VER_NT_WORKSTATION       = 0x0000001
)

// https://docs.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-osversioninfoexw
const (
	VER_SUITE_BACKOFFICE               = 0x00000004
	VER_SUITE_BLADE                    = 0x00000400
	VER_SUITE_COMPUTE_SERVER           = 0x00004000
	VER_SUITE_DATACENTER               = 0x00000080
	VER_SUITE_ENTERPRISE               = 0x00000002
	VER_SUITE_EMBEDDEDNT               = 0x00000040
	VER_SUITE_PERSONAL                 = 0x00000200
	VER_SUITE_SINGLEUSERTS             = 0x00000100
	VER_SUITE_SMALLBUSINESS            = 0x00000001
	VER_SUITE_SMALLBUSINESS_RESTRICTED = 0x00000020
	VER_SUITE_STORAGE_SERVER           = 0x00002000
	VER_SUITE_TERMINAL                 = 0x00000010
	VER_SUITE_WH_SERVER                = 0x00008000
	VER_SUITE_MULTIUSERTS              = 0x00020000
)

// https://docs.microsoft.com/en-us/windows/win32/api/sysinfoapi/ns-sysinfoapi-system_info
const (
	PROCESSOR_ARCHITECTURE_AMD64   = 9
	PROCESSOR_ARCHITECTURE_ARM     = 5
	PROCESSOR_ARCHITECTURE_ARM64   = 12
	PROCESSOR_ARCHITECTURE_IA64    = 6
	PROCESSOR_ARCHITECTURE_INTEL   = 0
	PROCESSOR_ARCHITECTURE_UNKNOWN = 0xFFFF
)

// Flags for SetLayeredWindowAttributes.
const (
	LWA_COLORKEY = 0x1
	LWA_ALPHA    = 0x2
)

// Flags for RedrawWindow.
const (
	RDW_INVALIDATE      = 0x0001
	RDW_INTERNALPAINT   = 0x0002
	RDW_ERASE           = 0x0004
	RDW_VALIDATE        = 0x0008
	RDW_NOINTERNALPAINT = 0x0010
	RDW_NOERASE         = 0x0020
	RDW_NOCHILDREN      = 0x0040
	RDW_ALLCHILDREN     = 0x0080
	RDW_UPDATENOW       = 0x0100
	RDW_ERASENOW        = 0x0200
	RDW_FRAME           = 0x0400
	RDW_NOFRAME         = 0x0800
)

// DrawIconEx flags.
const (
	DI_COMPAT      = 0x0004
	DI_DEFAULTSIZE = 0x0008
	DI_IMAGE       = 0x0002
	DI_MASK        = 0x0001
	DI_NOMIRROR    = 0x0010
	DI_NORMAL      = DI_IMAGE | DI_MASK
)

// Track bar messages.
const (
	TBM_GETPOS         = WM_USER
	TBM_GETRANGEMIN    = WM_USER + 1
	TBM_GETRANGEMAX    = WM_USER + 2
	TBM_GETTIC         = WM_USER + 3
	TBM_SETTIC         = WM_USER + 4
	TBM_SETPOS         = WM_USER + 5
	TBM_SETRANGE       = WM_USER + 6
	TBM_SETRANGEMIN    = WM_USER + 7
	TBM_SETRANGEMAX    = WM_USER + 8
	TBM_CLEARTICS      = WM_USER + 9
	TBM_SETSEL         = WM_USER + 10
	TBM_SETSELSTART    = WM_USER + 11
	TBM_SETSELEND      = WM_USER + 12
	TBM_GETPTICS       = WM_USER + 14
	TBM_GETTICPOS      = WM_USER + 15
	TBM_GETNUMTICS     = WM_USER + 16
	TBM_GETSELSTART    = WM_USER + 17
	TBM_GETSELEND      = WM_USER + 18
	TBM_CLEARSEL       = WM_USER + 19
	TBM_SETTICFREQ     = WM_USER + 20
	TBM_SETPAGESIZE    = WM_USER + 21
	TBM_GETPAGESIZE    = WM_USER + 22
	TBM_SETLINESIZE    = WM_USER + 23
	TBM_GETLINESIZE    = WM_USER + 24
	TBM_GETTHUMBRECT   = WM_USER + 25
	TBM_GETCHANNELRECT = WM_USER + 26
	TBM_SETTHUMBLENGTH = WM_USER + 27
	TBM_GETTHUMBLENGTH = WM_USER + 28
)

// Track bar styles.
const (
	TBS_AUTOTICKS      = 1
	TBS_VERT           = 2
	TBS_HORZ           = 0
	TBS_TOP            = 4
	TBS_BOTTOM         = 0
	TBS_LEFT           = 4
	TBS_RIGHT          = 0
	TBS_BOTH           = 8
	TBS_NOTICKS        = 16
	TBS_ENABLESELRANGE = 32
	TBS_FIXEDLENGTH    = 64
	TBS_NOTHUMB        = 128
)

// Track bar scroll notifications.
const (
	TB_LINEUP        = 0
	TB_LINEDOWN      = 1
	TB_PAGEUP        = 2
	TB_PAGEDOWN      = 3
	TB_THUMBPOSITION = 4
	TB_THUMBTRACK    = 5
	TB_TOP           = 6
	TB_BOTTOM        = 7
	TB_ENDTRACK      = 8
)

// Clip region types, e.g. returned from IntersectClipRect.
const (
	ERROR         = 0
	NULLREGION    = 1
	SIMPLEREGION  = 2
	COMPLEXREGION = 3
)

// Modes for CombineRgn.
const (
	RGN_AND  = 1
	RGN_OR   = 2
	RGN_XOR  = 3
	RGN_DIFF = 4
	RGN_COPY = 5
	RGN_MIN  = RGN_AND
	RGN_MAX  = RGN_COPY
)

const (
	DIGCF_PRESENT         = 0x0002
	DIGCF_ALLCLASSES      = 0x0004
	DIGCF_PROFILE         = 0x0008
	DIGCF_DEVICEINTERFACE = 0x0010
	DIGCF_INTERFACEDEVICE = DIGCF_DEVICEINTERFACE
)

const (
	DICS_FLAG_GLOBAL         = 0x00000001
	DICS_FLAG_CONFIGSPECIFIC = 0x00000002
	DICS_FLAG_CONFIGGENERAL  = 0x00000004
)

const (
	DIREG_DEV  = 0x00000001
	DIREG_DRV  = 0x00000002
	DIREG_BOTH = 0x00000004
)

const (
	ABM_NEW              = 0x0
	ABM_REMOVE           = 0x1
	ABM_QUERYPOS         = 0x2
	ABM_SETPOS           = 0x3
	ABM_GETSTATE         = 0x4
	ABM_GETTASKBARPOS    = 0x5
	ABM_ACTIVATE         = 0x6
	ABM_GETAUTOHIDEBAR   = 0x7
	ABM_SETAUTOHIDEBAR   = 0x8
	ABM_WINDOWPOSCHANGED = 0x9
	ABM_SETSTATE         = 0xA
)

const (
	ACM_OPEN = 0x464
	ACM_PLAY = 0x465
	ACM_STOP = 0x466
)

const (
	CBEM_GETCOMBOCONTROL  = 0x406
	CBEM_GETEDITCONTROL   = 0x407
	CBEM_GETEXSTYLE       = 0x409
	CBEM_GETEXTENDEDSTYLE = 0x409
	CBEM_GETIMAGELIST     = 0x403
	CBEM_GETITEMA         = 0x404
	CBEM_GETITEMW         = 0x40D
	CBEM_HASEDITCHANGED   = 0x40A
	CBEM_INSERTITEMA      = 0x401
	CBEM_INSERTITEMW      = 0x40B
	CBEM_SETEXSTYLE       = 0x408
	CBEM_SETEXTENDEDSTYLE = 0x40E
	CBEM_SETIMAGELIST     = 0x402
	CBEM_SETITEMA         = 0x405
	CBEM_SETITEMW         = 0x40C
	CBEM_GETUNICODEFORMAT = 0x2006
	CBEM_SETUNICODEFORMAT = 0x2005
)

const (
	CDM_GETFILEPATH     = 0x465
	CDM_GETFOLDERIDLIST = 0x467
	CDM_GETFOLDERPATH   = 0x466
	CDM_GETSPEC         = 0x464
	CDM_HIDECONTROL     = 0x469
	CDM_SETCONTROLTEXT  = 0x468
	CDM_SETDEFEXT       = 0x46A
)

const (
	DL_BEGINDRAG  = 0x485
	DL_CANCELDRAG = 0x488
	DL_COPYCURSOR = 0x002
	DL_CURSORSET  = 0x000
	DL_DRAGGING   = 0x486
	DL_DROPPED    = 0x487
	DL_MOVECURSOR = 0x003
	DL_STOPCURSOR = 0x001
)

const (
	DM_GETDEFID   = 0x400
	DM_REPOSITION = 0x402
	DM_SETDEFID   = 0x401
)

const (
	DTM_GETMCCOLOR    = 0x1007
	DTM_GETMCFONT     = 0x100A
	DTM_GETMONTHCAL   = 0x1008
	DTM_GETRANGE      = 0x1003
	DTM_GETSYSTEMTIME = 0x1001
	DTM_SETFORMATA    = 0x1005
	DTM_SETFORMATW    = 0x1032
	DTM_SETMCCOLOR    = 0x1006
	DTM_SETMCFONT     = 0x1009
	DTM_SETRANGE      = 0x1004
	DTM_SETSYSTEMTIME = 0x1002
)

const (
	HDM_CLEARFILTER            = 0x1218
	HDM_CREATEDRAGIMAGE        = 0x1210
	HDM_DELETEITEM             = 0x1202
	HDM_EDITFILTER             = 0x1217
	HDM_GETBITMAPMARGIN        = 0x1215
	HDM_GETIMAGELIST           = 0x1209
	HDM_GETITEM                = 0x1203
	HDM_GETITEMCOUNT           = 0x1200
	HDM_GETITEMRECT            = 0x1207
	HDM_GETORDERARRAY          = 0x1211
	HDM_GETUNICODEFORMAT       = 0x2006
	HDM_HITTEST                = 0x1206
	HDM_INSERTITEM             = 0x1201
	HDM_LAYOUT                 = 0x1205
	HDM_ORDERTOINDEX           = 0x120F
	HDM_SETBITMAPMARGIN        = 0x1214
	HDM_SETFILTERCHANGETIMEOUT = 0x1216
	HDM_SETHOTDIVIDER          = 0x1213
	HDM_SETIMAGELIST           = 0x1208
	HDM_SETITEM                = 0x1204
	HDM_SETORDERARRAY          = 0x1212
	HDM_SETUNICODEFORMAT       = 0x2005
)

const (
	HKM_GETHOTKEY = 0x402
	HKM_SETHOTKEY = 0x401
	HKM_SETRULES  = 0x403
)

const (
	IPM_CLEARADDRESS = 0x464
	IPM_GETADDRESS   = 0x466
	IPM_ISBLANK      = 0x469
	IPM_SETADDRESS   = 0x465
	IPM_SETFOCUS     = 0x468
	IPM_SETRANGE     = 0x467
)

const (
	LBCB_CARETOFF = 0x01A4
	LBCB_CARETON  = 0x01A3
)

const (
	MCM_GETCOLOR          = 0x100B
	MCM_GETCURSEL         = 0x1001
	MCM_GETFIRSTDAYOFWEEK = 0x1010
	MCM_GETMAXSELCOUNT    = 0x1003
	MCM_GETMAXTODAYWIDTH  = 0x1015
	MCM_GETMINREQRECT     = 0x1009
	MCM_GETMONTHDELTA     = 0x1013
	MCM_GETMONTHRANGE     = 0x1007
	MCM_GETRANGE          = 0x1011
	MCM_GETSELRANGE       = 0x1005
	MCM_GETTODAY          = 0x100D
	MCM_GETUNICODEFORMAT  = 0x2006
	MCM_HITTEST           = 0x100E
	MCM_SETCOLOR          = 0x100A
	MCM_SETCURSEL         = 0x1002
	MCM_SETDAYSTATE       = 0x1008
	MCM_SETFIRSTDAYOFWEEK = 0x100F
	MCM_SETMAXSELCOUNT    = 0x1004
	MCM_SETMONTHDELTA     = 0x1014
	MCM_SETRANGE          = 0x1012
	MCM_SETSELRANGE       = 0x1006
	MCM_SETTODAY          = 0x100C
	MCM_SETUNICODEFORMAT  = 0x2005
)

const (
	MN_BUTTONDOWN              = 0x01ED
	MN_BUTTONUP                = 0x01EF
	MN_CANCELMENUS             = 0x01E6
	MN_CLOSEHIERARCHY          = 0x01E4
	MN_DBLCLK                  = 0x01F1
	MN_FINDMENUWINDOWFROMPOINT = 0x01EB
	MN_GETHMENU                = 0x01E1
	MN_GETPPOPUPMENU           = 0x01EA
	MN_MOUSEMOVE               = 0x01EE
	MN_OPENHIERARCHY           = 0x01E3
	MN_SELECTFIRSTVALIDITEM    = 0x01E7
	MN_SELECTITEM              = 0x01E5
	MN_SETHMENU                = 0x01E0
	MN_SETTIMERTOOPENHIERARCHY = 0x01F0
	MN_SHOWPOPUPWINDOW         = 0x01EC
	MN_SIZEWINDOW              = 0x01E2
)

const (
	OCM_CHARTOITEM        = 0x202F
	OCM_COMMAND           = 0x2111
	OCM_COMPAREITEM       = 0x2039
	OCM_CTLCOLORBTN       = 0x2135
	OCM_CTLCOLORDLG       = 0x2136
	OCM_CTLCOLOREDIT      = 0x2133
	OCM_CTLCOLORLISTBOX   = 0x2134
	OCM_CTLCOLORMSGBOX    = 0x2132
	OCM_CTLCOLORSCROLLBAR = 0x2137
	OCM_CTLCOLORSTATIC    = 0x2138
	OCM_DELETEITEM        = 0x202D
	OCM_DRAWITEM          = 0x202B
	OCM_HSCROLL           = 0x2114
	OCM_MEASUREITEM       = 0x202C
	OCM_NOTIFY            = 0x204E
	OCM_PARENTNOTIFY      = 0x2210
	OCM_VKEYTOITEM        = 0x202E
	OCM_VSCROLL           = 0x2115
)

const (
	PGM_FORWARDMOUSE   = 0x1403
	PGM_GETBKCOLOR     = 0x1405
	PGM_GETBORDER      = 0x1407
	PGM_GETBUTTONSIZE  = 0x140B
	PGM_GETBUTTONSTATE = 0x140C
	PGM_GETDROPTARGET  = 0x2004
	PGM_GETPOS         = 0x1409
	PGM_RECALCSIZE     = 0x1402
	PGM_SETBKCOLOR     = 0x1404
	PGM_SETBORDER      = 0x1406
	PGM_SETBUTTONSIZE  = 0x140A
	PGM_SETCHILD       = 0x1401
	PGM_SETPOS         = 0x1408
)

const (
	PSM_ADDPAGE            = 0x467
	PSM_APPLY              = 0x46E
	PSM_CANCELTOCLOSE      = 0x46B
	PSM_CHANGED            = 0x468
	PSM_GETCURRENTPAGEHWND = 0x476
	PSM_GETRESULT          = 0x487
	PSM_GETTABCONTROL      = 0x474
	PSM_HWNDTOINDEX        = 0x481
	PSM_IDTOINDEX          = 0x485
	PSM_INDEXTOHWND        = 0x482
	PSM_INDEXTOID          = 0x486
	PSM_INDEXTOPAGE        = 0x484
	PSM_INSERTPAGE         = 0x477
	PSM_ISDIALOGMESSAGE    = 0x475
	PSM_PAGETOINDEX        = 0x483
	PSM_PRESSBUTTON        = 0x471
	PSM_QUERYSIBLINGS      = 0x46C
	PSM_REBOOTSYSTEM       = 0x46A
	PSM_RECALCPAGESIZES    = 0x488
	PSM_REMOVEPAGE         = 0x466
	PSM_RESTARTWINDOWS     = 0x469
	PSM_SETCURSEL          = 0x465
	PSM_SETCURSELID        = 0x472
	PSM_SETFINISHTEXT      = 0x473
	PSM_SETHEADERSUBTITLEA = 0x47F
	PSM_SETHEADERSUBTITLEW = 0x480
	PSM_SETHEADERTITLEA    = 0x47D
	PSM_SETHEADERTITLEW    = 0x47E
	PSM_SETTITLE           = 0x46F
	PSM_SETWIZBUTTONS      = 0x470
	PSM_UNCHANGED          = 0x46D
)

const (
	RB_BEGINDRAG        = 0x418
	RB_DELETEBAND       = 0x402
	RB_DRAGMOVE         = 0x41A
	RB_ENDDRAG          = 0x419
	RB_GETBANDBORDERS   = 0x422
	RB_GETBANDCOUNT     = 0x40C
	RB_GETBANDINFO      = 0x41D
	RB_GETBANDINFOA     = 0x41D
	RB_GETBANDINFOW     = 0x41C
	RB_GETBARHEIGHT     = 0x41B
	RB_GETBARINFO       = 0x403
	RB_GETBKCOLOR       = 0x414
	RB_GETPALETTE       = 0x426
	RB_GETRECT          = 0x409
	RB_GETROWCOUNT      = 0x40D
	RB_GETROWHEIGHT     = 0x40E
	RB_GETTEXTCOLOR     = 0x416
	RB_GETTOOLTIPS      = 0x411
	RB_HITTEST          = 0x408
	RB_IDTOINDEX        = 0x410
	RB_INSERTBAND       = 0x401
	RB_INSERTBANDA      = 0x401
	RB_INSERTBANDW      = 0x40A
	RB_MAXIMIZEBAND     = 0x41F
	RB_MINIMIZEBAND     = 0x41E
	RB_MOVEBAND         = 0x427
	RB_PUSHCHEVRON      = 0x42B
	RB_SETBANDINFO      = 0x406
	RB_SETBANDINFOA     = 0x406
	RB_SETBANDINFOW     = 0x40B
	RB_SETBARINFO       = 0x404
	RB_SETBKCOLOR       = 0x413
	RB_SETPALETTE       = 0x425
	RB_SETPARENT        = 0x407
	RB_SETTEXTCOLOR     = 0x415
	RB_SETTOOLTIPS      = 0x412
	RB_SHOWBAND         = 0x423
	RB_SIZETORECT       = 0x417
	RB_SETCOLORSCHEME   = 0x2002
	RB_GETCOLORSCHEME   = 0x2003
	RB_GETDROPTARGET    = 0x2004
	RB_SETUNICODEFORMAT = 0x2005
	RB_GETUNICODEFORMAT = 0x2006
)

const (
	SBM_ENABLE_ARROWS  = 0xE4
	SBM_GETPOS         = 0xE1
	SBM_GETRANGE       = 0xE3
	SBM_GETSCROLLINFO  = 0xEA
	SBM_SETPOS         = 0xE0
	SBM_SETRANGE       = 0xE2
	SBM_SETRANGEREDRAW = 0xE6
	SBM_SETSCROLLINFO  = 0xE9
)

const (
	SB_GETBORDERS       = 0x407
	SB_GETICON          = 0x414
	SB_GETPARTS         = 0x406
	SB_GETRECT          = 0x40A
	SB_GETTEXTA         = 0x402
	SB_GETTEXTLENGTHA   = 0x403
	SB_GETTEXTLENGTHW   = 0x40C
	SB_GETTEXTW         = 0x40D
	SB_GETTIPTEXTA      = 0x412
	SB_GETTIPTEXTW      = 0x413
	SB_ISSIMPLE         = 0x40E
	SB_SETICON          = 0x40F
	SB_SETMINHEIGHT     = 0x408
	SB_SETPARTS         = 0x404
	SB_SETTEXTA         = 0x401
	SB_SETTEXTW         = 0x40B
	SB_SETTIPTEXTA      = 0x410
	SB_SETTIPTEXTW      = 0x411
	SB_SIMPLE           = 0x409
	SB_SETBKCOLOR       = 0x2001
	SB_SETUNICODEFORMAT = 0x2005
	SB_GETUNICODEFORMAT = 0x2006
)

const (
	STM_GETICON  = 0x171
	STM_GETIMAGE = 0x173
	STM_SETICON  = 0x170
	STM_SETIMAGE = 0x172
)

const (
	TB_ADDBITMAP             = 0x413
	TB_ADDBUTTONS            = 0x414
	TB_ADDSTRING             = 0x41C
	TB_AUTOSIZE              = 0x421
	TB_BUTTONCOUNT           = 0x418
	TB_BUTTONSTRUCTSIZE      = 0x41E
	TB_CHANGEBITMAP          = 0x42B
	TB_CHECKBUTTON           = 0x402
	TB_COMMANDTOINDEX        = 0x419
	TB_CUSTOMIZE             = 0x41B
	TB_DELETEBUTTON          = 0x416
	TB_ENABLEBUTTON          = 0x401
	TB_GETANCHORHIGHLIGHT    = 0x44A
	TB_GETBITMAP             = 0x42C
	TB_GETBITMAPFLAGS        = 0x429
	TB_GETBUTTON             = 0x417
	TB_GETBUTTONINFOA        = 0x441
	TB_GETBUTTONINFOW        = 0x43F
	TB_GETBUTTONSIZE         = 0x43A
	TB_GETBUTTONTEXT         = 0x42D
	TB_GETCOLORSCHEME        = 0x2003
	TB_GETDISABLEDIMAGELIST  = 0x437
	TB_GETEXTENDEDSTYLE      = 0x455
	TB_GETHOTIMAGELIST       = 0x435
	TB_GETHOTITEM            = 0x447
	TB_GETIMAGELIST          = 0x431
	TB_GETINSERTMARK         = 0x44F
	TB_GETINSERTMARKCOLOR    = 0x459
	TB_GETITEMRECT           = 0x41D
	TB_GETMAXSIZE            = 0x453
	TB_GETOBJECT             = 0x43E
	TB_GETPADDING            = 0x456
	TB_GETRECT               = 0x433
	TB_GETROWS               = 0x428
	TB_GETSTATE              = 0x412
	TB_GETSTRING             = 0x45C
	TB_GETSTYLE              = 0x439
	TB_GETTEXTROWS           = 0x43D
	TB_GETTOOLTIPS           = 0x423
	TB_GETUNICODEFORMAT      = 0x2006
	TB_HIDEBUTTON            = 0x404
	TB_HITTEST               = 0x445
	TB_INDETERMINATE         = 0x405
	TB_INSERTBUTTON          = 0x415
	TB_INSERTMARKHITTEST     = 0x451
	TB_ISBUTTONCHECKED       = 0x40A
	TB_ISBUTTONENABLED       = 0x409
	TB_ISBUTTONHIDDEN        = 0x40C
	TB_ISBUTTONHIGHLIGHTED   = 0x40E
	TB_ISBUTTONINDETERMINATE = 0x40D
	TB_ISBUTTONPRESSED       = 0x40B
	TB_LOADIMAGES            = 0x432
	TB_MAPACCELERATORA       = 0x44E
	TB_MAPACCELERATORW       = 0x45A
	TB_MARKBUTTON            = 0x406
	TB_MOVEBUTTON            = 0x452
	TB_PRESSBUTTON           = 0x403
	TB_REPLACEBITMAP         = 0x42E
	TB_SAVERESTORE           = 0x41A
	TB_SETANCHORHIGHLIGHT    = 0x449
	TB_SETBITMAPSIZE         = 0x420
	TB_SETBUTTONINFOA        = 0x442
	TB_SETBUTTONINFOW        = 0x440
	TB_SETBUTTONSIZE         = 0x41F
	TB_SETBUTTONWIDTH        = 0x43B
	TB_SETCMDID              = 0x42A
	TB_SETCOLORSCHEME        = 0x2002
	TB_SETDISABLEDIMAGELIST  = 0x436
	TB_SETDRAWTEXTFLAGS      = 0x446
	TB_SETEXTENDEDSTYLE      = 0x454
	TB_SETHOTIMAGELIST       = 0x434
	TB_SETHOTITEM            = 0x448
	TB_SETIMAGELIST          = 0x430
	TB_SETINDENT             = 0x42F
	TB_SETINSERTMARK         = 0x450
	TB_SETINSERTMARKCOLOR    = 0x458
	TB_SETMAXTEXTROWS        = 0x43C
	TB_SETPADDING            = 0x457
	TB_SETPARENT             = 0x425
	TB_SETROWS               = 0x427
	TB_SETSTATE              = 0x411
	TB_SETSTYLE              = 0x438
	TB_SETTOOLTIPS           = 0x424
	TB_SETUNICODEFORMAT      = 0x2005
)

const (
	TCM_ADJUSTRECT       = 0x1328
	TCM_DELETEALLITEMS   = 0x1309
	TCM_DELETEITEM       = 0x1308
	TCM_DESELECTALL      = 0x1332
	TCM_GETCURFOCUS      = 0x132F
	TCM_GETCURSEL        = 0x130B
	TCM_GETEXTENDEDSTYLE = 0x1335
	TCM_GETIMAGELIST     = 0x1302
	TCM_GETITEM          = 0x1305
	TCM_GETITEMCOUNT     = 0x1304
	TCM_GETITEMRECT      = 0x130A
	TCM_GETROWCOUNT      = 0x132C
	TCM_GETTOOLTIPS      = 0x132D
	TCM_GETUNICODEFORMAT = 0x2006
	TCM_HIGHLIGHTITEM    = 0x1333
	TCM_HITTEST          = 0x130D
	TCM_INSERTITEM       = 0x1307
	TCM_REMOVEIMAGE      = 0x132A
	TCM_SETCURFOCUS      = 0x1330
	TCM_SETCURSEL        = 0x130C
	TCM_SETEXTENDEDSTYLE = 0x1334
	TCM_SETIMAGELIST     = 0x1303
	TCM_SETITEM          = 0x1306
	TCM_SETITEMEXTRA     = 0x130E
	TCM_SETITEMSIZE      = 0x1329
	TCM_SETMINTABWIDTH   = 0x1331
	TCM_SETPADDING       = 0x132B
	TCM_SETTOOLTIPS      = 0x132E
	TCM_SETUNICODEFORMAT = 0x2005
)

const (
	TVM_CREATEDRAGIMAGE    = 0x1112
	TVM_DELETEITEM         = 0x1101
	TVM_EDITLABEL          = 0x110E
	TVM_ENDEDITLABELNOW    = 0x1116
	TVM_ENSUREVISIBLE      = 0x1114
	TVM_EXPAND             = 0x1102
	TVM_GETBKCOLOR         = 0x111F
	TVM_GETCOUNT           = 0x1105
	TVM_GETEDITCONTROL     = 0x110F
	TVM_GETIMAGELIST       = 0x1108
	TVM_GETINDENT          = 0x1106
	TVM_GETINSERTMARKCOLOR = 0x1126
	TVM_GETISEARCHSTRING   = 0x1117
	TVM_GETITEM            = 0x110C
	TVM_GETITEMHEIGHT      = 0x111C
	TVM_GETITEMRECT        = 0x1104
	TVM_GETITEMSTATE       = 0x1127
	TVM_GETLINECOLOR       = 0x1129
	TVM_GETNEXTITEM        = 0x110A
	TVM_GETSCROLLTIME      = 0x1122
	TVM_GETTEXTCOLOR       = 0x1120
	TVM_GETTOOLTIPS        = 0x1119
	TVM_GETUNICODEFORMAT   = 0x2006
	TVM_GETVISIBLECOUNT    = 0x1110
	TVM_HITTEST            = 0x1111
	TVM_INSERTITEM         = 0x1100
	TVM_SELECTITEM         = 0x110B
	TVM_SETBKCOLOR         = 0x111D
	TVM_SETIMAGELIST       = 0x1109
	TVM_SETINDENT          = 0x1107
	TVM_SETINSERTMARK      = 0x111A
	TVM_SETINSERTMARKCOLOR = 0x1125
	TVM_SETITEM            = 0x110D
	TVM_SETITEMHEIGHT      = 0x111B
	TVM_SETLINECOLOR       = 0x1128
	TVM_SETSCROLLTIME      = 0x1121
	TVM_SETTEXTCOLOR       = 0x111E
	TVM_SETTOOLTIPS        = 0x1118
	TVM_SETUNICODEFORMAT   = 0x2005
	TVM_SORTCHILDREN       = 0x1113
	TVM_SORTCHILDRENCB     = 0x1115
)

// PROCESS_DPI_AWARENESS
const (
	PROCESS_DPI_UNAWARE           = 0
	PROCESS_SYSTEM_DPI_AWARE      = 1
	PROCESS_PER_MONITOR_DPI_AWARE = 2
)

// WM_DEVICECHANGE WPARAM options
const (
	DBT_DEVNODES_CHANGED        = 0x0007
	DBT_QUERYCHANGECONFIG       = 0x0017
	DBT_CONFIGCHANGED           = 0x0018
	DBT_CONFIGCHANGECANCELED    = 0x0019
	DBT_DEVICEARRIVAL           = 0x8000
	DBT_DEVICEQUERYREMOVE       = 0x8001
	DBT_DEVICEQUERYREMOVEFAILED = 0x8002
	DBT_DEVICEREMOVEPENDING     = 0x8003
	DBT_DEVICEREMOVECOMPLETE    = 0x8004
	DBT_DEVICETYPESPECIFIC      = 0x8005
	DBT_CUSTOMEVENT             = 0x8006
	DBT_USERDEFINED             = 0xFFFF
)

const (
	SPI_GETACCESSTIMEOUT            = 0x003C
	SPI_GETAUDIODESCRIPTION         = 0x0074
	SPI_GETCLIENTAREAANIMATION      = 0x1042
	SPI_GETDISABLEOVERLAPPEDCONTENT = 0x1040
	SPI_GETFILTERKEYS               = 0x0032
	SPI_GETFOCUSBORDERHEIGHT        = 0x2010
	SPI_GETFOCUSBORDERWIDTH         = 0x200E
	SPI_GETHIGHCONTRAST             = 0x0042
	SPI_GETLOGICALDPIOVERRIDE       = 0x009E
	SPI_GETMESSAGEDURATION          = 0x2016
	SPI_GETMOUSECLICKLOCK           = 0x101E
	SPI_GETMOUSECLICKLOCKTIME       = 0x2008
	SPI_GETMOUSEKEYS                = 0x0036
	SPI_GETMOUSESONAR               = 0x101C
	SPI_GETMOUSEVANISH              = 0x1020
	SPI_GETSCREENREADER             = 0x0046
	SPI_GETSERIALKEYS               = 0x003E
	SPI_GETSHOWSOUNDS               = 0x0038
	SPI_GETSOUNDSENTRY              = 0x0040
	SPI_GETSTICKYKEYS               = 0x003A
	SPI_GETTOGGLEKEYS               = 0x0034
	SPI_SETACCESSTIMEOUT            = 0x003D
	SPI_SETAUDIODESCRIPTION         = 0x0075
	SPI_SETCLIENTAREAANIMATION      = 0x1043
	SPI_SETDISABLEOVERLAPPEDCONTENT = 0x1041
	SPI_SETFILTERKEYS               = 0x0033
	SPI_SETFOCUSBORDERHEIGHT        = 0x2011
	SPI_SETFOCUSBORDERWIDTH         = 0x200F
	SPI_SETHIGHCONTRAST             = 0x0043
	SPI_SETLOGICALDPIOVERRIDE       = 0x009F
	SPI_SETMESSAGEDURATION          = 0x2017
	SPI_SETMOUSECLICKLOCK           = 0x101F
	SPI_SETMOUSECLICKLOCKTIME       = 0x2009
	SPI_SETMOUSEKEYS                = 0x0037
	SPI_SETMOUSESONAR               = 0x101D
	SPI_SETMOUSEVANISH              = 0x1021
	SPI_SETSCREENREADER             = 0x0047
	SPI_SETSERIALKEYS               = 0x003F
	SPI_SETSHOWSOUNDS               = 0x0039
	SPI_SETSOUNDSENTRY              = 0x0041
	SPI_SETSTICKYKEYS               = 0x003B
	SPI_SETTOGGLEKEYS               = 0x0035
	SPI_GETCLEARTYPE                = 0x1048
	SPI_GETDESKWALLPAPER            = 0x0073
	SPI_GETDROPSHADOW               = 0x1024
	SPI_GETFLATMENU                 = 0x1022
	SPI_GETFONTSMOOTHING            = 0x004A
	SPI_GETFONTSMOOTHINGCONTRAST    = 0x200C
	SPI_GETFONTSMOOTHINGORIENTATION = 0x2012
	SPI_GETFONTSMOOTHINGTYPE        = 0x200A
	SPI_GETWORKAREA                 = 0x0030
	SPI_SETCLEARTYPE                = 0x1049
	SPI_SETCURSORS                  = 0x0057
	SPI_SETDESKPATTERN              = 0x0015
	SPI_SETDESKWALLPAPER            = 0x0014
	SPI_SETDROPSHADOW               = 0x1025
	SPI_SETFLATMENU                 = 0x1023
	SPI_SETFONTSMOOTHING            = 0x004B
	SPI_SETFONTSMOOTHINGCONTRAST    = 0x200D
	SPI_SETFONTSMOOTHINGORIENTATION = 0x2013
	SPI_SETFONTSMOOTHINGTYPE        = 0x200B
	SPI_SETWORKAREA                 = 0x002F
	SPI_GETICONMETRICS              = 0x002D
	SPI_GETICONTITLELOGFONT         = 0x001F
	SPI_GETICONTITLEWRAP            = 0x0019
	SPI_ICONHORIZONTALSPACING       = 0x000D
	SPI_ICONVERTICALSPACING         = 0x0018
	SPI_SETICONMETRICS              = 0x002E
	SPI_SETICONS                    = 0x0058
	SPI_SETICONTITLELOGFONT         = 0x0022
	SPI_SETICONTITLEWRAP            = 0x001A
	SPI_GETBEEP                     = 0x0001
	SPI_GETBLOCKSENDINPUTRESETS     = 0x1026
	SPI_GETCONTACTVISUALIZATION     = 0x2018
	SPI_GETDEFAULTINPUTLANG         = 0x0059
	SPI_GETGESTUREVISUALIZATION     = 0x201A
	SPI_GETKEYBOARDCUES             = 0x100A
	SPI_GETKEYBOARDDELAY            = 0x0016
	SPI_GETKEYBOARDPREF             = 0x0044
	SPI_GETKEYBOARDSPEED            = 0x000A
	SPI_GETMOUSE                    = 0x0003
	SPI_GETMOUSEHOVERHEIGHT         = 0x0064
	SPI_GETMOUSEHOVERTIME           = 0x0066
	SPI_GETMOUSEHOVERWIDTH          = 0x0062
	SPI_GETMOUSESPEED               = 0x0070
	SPI_GETMOUSETRAILS              = 0x005E
	SPI_GETMOUSEWHEELROUTING        = 0x201C
	SPI_GETPENVISUALIZATION         = 0x201E
	SPI_GETSNAPTODEFBUTTON          = 0x005F
	SPI_GETSYSTEMLANGUAGEBAR        = 0x1050
	SPI_GETTHREADLOCALINPUTSETTINGS = 0x104E
	SPI_GETWHEELSCROLLCHARS         = 0x006C
	SPI_GETWHEELSCROLLLINES         = 0x0068
	SPI_SETBEEP                     = 0x0002
	SPI_SETBLOCKSENDINPUTRESETS     = 0x1027
	SPI_SETCONTACTVISUALIZATION     = 0x2019
	SPI_SETDEFAULTINPUTLANG         = 0x005A
	SPI_SETDOUBLECLICKTIME          = 0x0020
	SPI_SETDOUBLECLKHEIGHT          = 0x001E
	SPI_SETDOUBLECLKWIDTH           = 0x001D
	SPI_SETGESTUREVISUALIZATION     = 0x201B
	SPI_SETKEYBOARDCUES             = 0x100B
	SPI_SETKEYBOARDDELAY            = 0x0017
	SPI_SETKEYBOARDPREF             = 0x0045
	SPI_SETKEYBOARDSPEED            = 0x000B
	SPI_SETLANGTOGGLE               = 0x005B
	SPI_SETMOUSE                    = 0x0004
	SPI_SETMOUSEBUTTONSWAP          = 0x0021
	SPI_SETMOUSEHOVERHEIGHT         = 0x0065
	SPI_SETMOUSEHOVERTIME           = 0x0067
	SPI_SETMOUSEHOVERWIDTH          = 0x0063
	SPI_SETMOUSESPEED               = 0x0071
	SPI_SETMOUSETRAILS              = 0x005D
	SPI_SETMOUSEWHEELROUTING        = 0x201D
	SPI_SETPENVISUALIZATION         = 0x201F
	SPI_SETSNAPTODEFBUTTON          = 0x0060
	SPI_SETSYSTEMLANGUAGEBAR        = 0x1051
	SPI_SETTHREADLOCALINPUTSETTINGS = 0x104F
	SPI_SETWHEELSCROLLCHARS         = 0x006D
	SPI_SETWHEELSCROLLLINES         = 0x0069
	SPI_GETMENUDROPALIGNMENT        = 0x001B
	SPI_GETMENUFADE                 = 0x1012
	SPI_GETMENUSHOWDELAY            = 0x006A
	SPI_SETMENUDROPALIGNMENT        = 0x001C
	SPI_SETMENUFADE                 = 0x1013
	SPI_SETMENUSHOWDELAY            = 0x006B
	SPI_GETLOWPOWERACTIVE           = 0x0053
	SPI_GETLOWPOWERTIMEOUT          = 0x004F
	SPI_GETPOWEROFFACTIVE           = 0x0054
	SPI_GETPOWEROFFTIMEOUT          = 0x0050
	SPI_SETLOWPOWERACTIVE           = 0x0055
	SPI_SETLOWPOWERTIMEOUT          = 0x0051
	SPI_SETPOWEROFFACTIVE           = 0x0056
	SPI_SETPOWEROFFTIMEOUT          = 0x0052
	SPI_GETSCREENSAVEACTIVE         = 0x0010
	SPI_GETSCREENSAVERRUNNING       = 0x0072
	SPI_GETSCREENSAVESECURE         = 0x0076
	SPI_GETSCREENSAVETIMEOUT        = 0x000E
	SPI_SETSCREENSAVEACTIVE         = 0x0011
	SPI_SETSCREENSAVESECURE         = 0x0077
	SPI_SETSCREENSAVETIMEOUT        = 0x000F
	SPI_GETHUNGAPPTIMEOUT           = 0x0078
	SPI_GETWAITTOKILLTIMEOUT        = 0x007A
	SPI_GETWAITTOKILLSERVICETIMEOUT = 0x007C
	SPI_SETHUNGAPPTIMEOUT           = 0x0079
	SPI_SETWAITTOKILLTIMEOUT        = 0x007B
	SPI_SETWAITTOKILLSERVICETIMEOUT = 0x007D
	SPI_GETCOMBOBOXANIMATION        = 0x1004
	SPI_GETCURSORSHADOW             = 0x101A
	SPI_GETGRADIENTCAPTIONS         = 0x1008
	SPI_GETHOTTRACKING              = 0x100E
	SPI_GETLISTBOXSMOOTHSCROLLING   = 0x1006
	SPI_GETMENUANIMATION            = 0x1002
	SPI_GETMENUUNDERLINES           = 0x100A
	SPI_GETSELECTIONFADE            = 0x1014
	SPI_GETTOOLTIPANIMATION         = 0x1016
	SPI_GETTOOLTIPFADE              = 0x1018
	SPI_GETUIEFFECTS                = 0x103E
	SPI_SETCOMBOBOXANIMATION        = 0x1005
	SPI_SETCURSORSHADOW             = 0x101B
	SPI_SETGRADIENTCAPTIONS         = 0x1009
	SPI_SETHOTTRACKING              = 0x100F
	SPI_SETLISTBOXSMOOTHSCROLLING   = 0x1007
	SPI_SETMENUANIMATION            = 0x1003
	SPI_SETMENUUNDERLINES           = 0x100B
	SPI_SETSELECTIONFADE            = 0x1015
	SPI_SETTOOLTIPANIMATION         = 0x1017
	SPI_SETTOOLTIPFADE              = 0x1019
	SPI_SETUIEFFECTS                = 0x103F
	SPI_GETACTIVEWINDOWTRACKING     = 0x1000
	SPI_GETACTIVEWNDTRKZORDER       = 0x100C
	SPI_GETACTIVEWNDTRKTIMEOUT      = 0x2002
	SPI_GETANIMATION                = 0x0048
	SPI_GETBORDER                   = 0x0005
	SPI_GETCARETWIDTH               = 0x2006
	SPI_GETDOCKMOVING               = 0x0090
	SPI_GETDRAGFROMMAXIMIZE         = 0x008C
	SPI_GETDRAGFULLWINDOWS          = 0x0026
	SPI_GETFOREGROUNDFLASHCOUNT     = 0x2004
	SPI_GETFOREGROUNDLOCKTIMEOUT    = 0x2000
	SPI_GETMINIMIZEDMETRICS         = 0x002B
	SPI_GETMOUSEDOCKTHRESHOLD       = 0x007E
	SPI_GETMOUSEDRAGOUTTHRESHOLD    = 0x0084
	SPI_GETMOUSESIDEMOVETHRESHOLD   = 0x0088
	SPI_GETNONCLIENTMETRICS         = 0x0029
	SPI_GETPENDOCKTHRESHOLD         = 0x0080
	SPI_GETPENDRAGOUTTHRESHOLD      = 0x0086
	SPI_GETPENSIDEMOVETHRESHOLD     = 0x008A
	SPI_GETSHOWIMEUI                = 0x006E
	SPI_GETSNAPSIZING               = 0x008E
	SPI_GETWINARRANGING             = 0x0082
	SPI_SETACTIVEWINDOWTRACKING     = 0x1001
	SPI_SETACTIVEWNDTRKZORDER       = 0x100D
	SPI_SETACTIVEWNDTRKTIMEOUT      = 0x2003
	SPI_SETANIMATION                = 0x0049
	SPI_SETBORDER                   = 0x0006
	SPI_SETCARETWIDTH               = 0x2007
	SPI_SETDOCKMOVING               = 0x0091
	SPI_SETDRAGFROMMAXIMIZE         = 0x008D
	SPI_SETDRAGFULLWINDOWS          = 0x0025
	SPI_SETDRAGHEIGHT               = 0x004D
	SPI_SETDRAGWIDTH                = 0x004C
	SPI_SETFOREGROUNDFLASHCOUNT     = 0x2005
	SPI_SETFOREGROUNDLOCKTIMEOUT    = 0x2001
	SPI_SETMINIMIZEDMETRICS         = 0x002C
	SPI_SETMOUSEDOCKTHRESHOLD       = 0x007F
	SPI_SETMOUSEDRAGOUTTHRESHOLD    = 0x0085
	SPI_SETMOUSESIDEMOVETHRESHOLD   = 0x0089
	SPI_SETNONCLIENTMETRICS         = 0x002A
	SPI_SETPENDOCKTHRESHOLD         = 0x0081
	SPI_SETPENDRAGOUTTHRESHOLD      = 0x0087
	SPI_SETPENSIDEMOVETHRESHOLD     = 0x008B
	SPI_SETSHOWIMEUI                = 0x006F
	SPI_SETSNAPSIZING               = 0x008F
	SPI_SETWINARRANGING             = 0x0083
)

const (
	SPIF_UPDATEINIFILE    = 0x1
	SPIF_SENDCHANGE       = 0x2
	SPIF_SENDWININICHANGE = SPIF_SENDCHANGE
)
