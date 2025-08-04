//go:build windows

package application

import (
	"fmt"
	"github.com/wailsapp/wails/v3/pkg/w32"
	"syscall"
	"unsafe"
)

type JumpListItemType int

const (
	JumpListItemTypeTask JumpListItemType = iota
	JumpListItemTypeSeparator
)

type JumpListItem struct {
	Type        JumpListItemType
	Title       string
	Description string
	FilePath    string
	Arguments   string
	IconPath    string
	IconIndex   int
}

type JumpListCategory struct {
	Name  string
	Items []JumpListItem
}

type JumpList struct {
	app        *windowsApp
	categories []JumpListCategory
}

var (
	modole32                          = syscall.NewLazyDLL("ole32.dll")
	modshell32                        = syscall.NewLazyDLL("shell32.dll")
	procCoCreateInstance              = modole32.NewProc("CoCreateInstance")
	procSHCreateItemFromParsingName   = modshell32.NewProc("SHCreateItemFromParsingName")
)

const (
	CLSID_DestinationList = "{77F10CF0-3DB5-4966-B520-B7C54FD35ED6}"
	IID_ICustomDestinationList = "{6332DEBF-87B5-4670-90C0-5E57B408A49E}"
	CLSID_ShellLink = "{00021401-0000-0000-C000-000000000046}"
	IID_IShellLink = "{000214F9-0000-0000-C000-000000000046}"
	IID_IPropertyStore = "{886D8EEB-8CF2-4446-8D02-CDBA1DBDCF99}"
	IID_IObjectArray = "{92CA9DCD-5622-4BBA-A805-5E9F541BD8C9}"
	IID_IObjectCollection = "{5632B1A4-E38A-400A-928A-D4CD63230295}"
	IID_IShellItem = "{43826D1E-E718-42EE-BC55-A1E261C37BFE}"
)

var (
	CLSID_DestinationListGUID = w32.NewGUID(CLSID_DestinationList)
	IID_ICustomDestinationListGUID = w32.NewGUID(IID_ICustomDestinationList)
	CLSID_ShellLinkGUID = w32.NewGUID(CLSID_ShellLink)
	IID_IShellLinkGUID = w32.NewGUID(IID_IShellLink)
	IID_IPropertyStoreGUID = w32.NewGUID(IID_IPropertyStore)
	IID_IObjectArrayGUID = w32.NewGUID(IID_IObjectArray)
	IID_IObjectCollectionGUID = w32.NewGUID(IID_IObjectCollection)
	IID_IShellItemGUID = w32.NewGUID(IID_IShellItem)
	CLSID_EnumerableObjectCollectionGUID = w32.NewGUID("{2D3468C1-36A7-43B6-AC24-D3F02FD9607A}")
)

type ICustomDestinationListVtbl struct {
	QueryInterface          uintptr
	AddRef                  uintptr
	Release                 uintptr
	SetAppID                uintptr
	BeginList               uintptr
	AppendCategory          uintptr
	AppendKnownCategory     uintptr
	AddUserTasks            uintptr
	CommitList              uintptr
	GetRemovedDestinations  uintptr
	DeleteList              uintptr
	AbortList               uintptr
}

type ICustomDestinationList struct {
	lpVtbl *ICustomDestinationListVtbl
}

type IShellLinkVtbl struct {
	QueryInterface          uintptr
	AddRef                  uintptr
	Release                 uintptr
	GetPath                 uintptr
	GetIDList               uintptr
	SetIDList               uintptr
	GetDescription          uintptr
	SetDescription          uintptr
	GetWorkingDirectory     uintptr
	SetWorkingDirectory     uintptr
	GetArguments            uintptr
	SetArguments            uintptr
	GetHotkey               uintptr
	SetHotkey               uintptr
	GetShowCmd              uintptr
	SetShowCmd              uintptr
	GetIconLocation         uintptr
	SetIconLocation         uintptr
	SetRelativePath         uintptr
	Resolve                 uintptr
	SetPath                 uintptr
}

type IShellLink struct {
	lpVtbl *IShellLinkVtbl
}

type IPropertyStoreVtbl struct {
	QueryInterface          uintptr
	AddRef                  uintptr
	Release                 uintptr
	GetCount                uintptr
	GetAt                   uintptr
	GetValue                uintptr
	SetValue                uintptr
	Commit                  uintptr
}

type IPropertyStore struct {
	lpVtbl *IPropertyStoreVtbl
}

type IObjectCollectionVtbl struct {
	QueryInterface          uintptr
	AddRef                  uintptr
	Release                 uintptr
	GetCount                uintptr
	GetAt                   uintptr
	AddObject               uintptr
	AddFromArray            uintptr
	RemoveObjectAt          uintptr
	Clear                   uintptr
}

type IObjectCollection struct {
	lpVtbl *IObjectCollectionVtbl
}

type PROPERTYKEY struct {
	Fmtid w32.GUID
	Pid   uint32
}

var PKEY_Title = PROPERTYKEY{
	Fmtid: *w32.NewGUID("{F29F85E0-4FF9-1068-AB91-08002B27B3D9}"),
	Pid:   2,
}

type PROPVARIANT struct {
	Vt         uint16
	Reserved1  uint16
	Reserved2  uint16
	Reserved3  uint16
	Val        [16]byte
}

func (app *windowsApp) CreateJumpList() *JumpList {
	return &JumpList{
		app:        app,
		categories: []JumpListCategory{},
	}
}

func (j *JumpList) AddCategory(category JumpListCategory) {
	j.categories = append(j.categories, category)
}

func (j *JumpList) ClearCategories() {
	j.categories = []JumpListCategory{}
}

func (j *JumpList) Apply() error {
	hr := w32.CoInitializeEx(0, w32.COINIT_APARTMENTTHREADED)
	if hr != w32.S_OK && hr != w32.S_FALSE {
		return fmt.Errorf("CoInitializeEx failed: %v", hr)
	}
	defer w32.CoUninitialize()

	var pDestList *ICustomDestinationList
	hr = CoCreateInstance(
		CLSID_DestinationListGUID,
		nil,
		w32.CLSCTX_INPROC_SERVER,
		IID_ICustomDestinationListGUID,
		&pDestList,
	)
	if hr != w32.S_OK {
		return fmt.Errorf("CoCreateInstance failed: %v", hr)
	}
	defer pDestList.Release()

	appID := w32.MustStringToUTF16Ptr(j.app.parent.options.Name)
	
	hr = pDestList.SetAppID(appID)
	if hr != w32.S_OK {
		return fmt.Errorf("SetAppID failed: %v", hr)
	}

	var cMinSlots uint32
	var pRemovedItems uintptr
	hr = pDestList.BeginList(&cMinSlots, IID_IObjectArrayGUID, &pRemovedItems)
	if hr != w32.S_OK {
		return fmt.Errorf("BeginList failed: %v", hr)
	}

	hasItems := false
	for _, category := range j.categories {
		if len(category.Items) > 0 {
			if category.Name == "" {
				// Add as tasks
				err := j.addTasks(pDestList, category.Items)
				if err != nil {
					pDestList.AbortList()
					return err
				}
			} else {
				// Add as custom category
				err := j.addCategory(pDestList, category)
				if err != nil {
					pDestList.AbortList()
					return err
				}
			}
			hasItems = true
		}
	}

	if !hasItems {
		// Clear the jump list if no items
		pDestList.DeleteList(appID)
		return nil
	}

	hr = pDestList.CommitList()
	if hr != w32.S_OK {
		return fmt.Errorf("CommitList failed: %v", hr)
	}

	return nil
}

func (j *JumpList) addTasks(pDestList *ICustomDestinationList, items []JumpListItem) error {
	var pObjectCollection *IObjectCollection
	hr := CoCreateInstance(
		CLSID_EnumerableObjectCollectionGUID,
		nil,
		w32.CLSCTX_INPROC_SERVER,
		IID_IObjectCollectionGUID,
		&pObjectCollection,
	)
	if hr != w32.S_OK {
		return fmt.Errorf("CoCreateInstance for IObjectCollection failed: %v", hr)
	}
	defer pObjectCollection.Release()

	for _, item := range items {
		if item.Type == JumpListItemTypeSeparator {
			// Skip separators in tasks
			continue
		}

		shellLink, err := j.createShellLink(item)
		if err != nil {
			return err
		}

		hr = pObjectCollection.AddObject(shellLink)
		shellLink.Release()
		if hr != w32.S_OK {
			return fmt.Errorf("AddObject failed: %v", hr)
		}
	}

	hr = pDestList.AddUserTasks(pObjectCollection)
	if hr != w32.S_OK {
		return fmt.Errorf("AddUserTasks failed: %v", hr)
	}

	return nil
}

func (j *JumpList) addCategory(pDestList *ICustomDestinationList, category JumpListCategory) error {
	var pObjectCollection *IObjectCollection
	hr := CoCreateInstance(
		CLSID_EnumerableObjectCollectionGUID,
		nil,
		w32.CLSCTX_INPROC_SERVER,
		IID_IObjectCollectionGUID,
		&pObjectCollection,
	)
	if hr != w32.S_OK {
		return fmt.Errorf("CoCreateInstance for IObjectCollection failed: %v", hr)
	}
	defer pObjectCollection.Release()

	for _, item := range category.Items {
		if item.Type == JumpListItemTypeSeparator {
			// Skip separators in custom categories
			continue
		}

		shellLink, err := j.createShellLink(item)
		if err != nil {
			return err
		}

		hr = pObjectCollection.AddObject(shellLink)
		shellLink.Release()
		if hr != w32.S_OK {
			return fmt.Errorf("AddObject failed: %v", hr)
		}
	}

	categoryName := w32.MustStringToUTF16Ptr(category.Name)

	hr = pDestList.AppendCategory(categoryName, pObjectCollection)
	if hr != w32.S_OK {
		return fmt.Errorf("AppendCategory failed: %v", hr)
	}

	return nil
}

func (j *JumpList) createShellLink(item JumpListItem) (*IShellLink, error) {
	var pShellLink *IShellLink
	hr := CoCreateInstance(
		CLSID_ShellLinkGUID,
		nil,
		w32.CLSCTX_INPROC_SERVER,
		IID_IShellLinkGUID,
		&pShellLink,
	)
	if hr != w32.S_OK {
		return nil, fmt.Errorf("CoCreateInstance for IShellLink failed: %v", hr)
	}

	// Set path
	path := w32.MustStringToUTF16Ptr(item.FilePath)
	hr = pShellLink.SetPath(path)
	if hr != w32.S_OK {
		pShellLink.Release()
		return nil, fmt.Errorf("SetPath failed: %v", hr)
	}

	// Set arguments
	if item.Arguments != "" {
		args := w32.MustStringToUTF16Ptr(item.Arguments)
		hr = pShellLink.SetArguments(args)
		if hr != w32.S_OK {
			pShellLink.Release()
			return nil, fmt.Errorf("SetArguments failed: %v", hr)
		}
	}

	// Set description
	if item.Description != "" {
		desc := w32.MustStringToUTF16Ptr(item.Description)
		hr = pShellLink.SetDescription(desc)
		if hr != w32.S_OK {
			pShellLink.Release()
			return nil, fmt.Errorf("SetDescription failed: %v", hr)
		}
	}

	// Set icon
	if item.IconPath != "" {
		iconPath := w32.MustStringToUTF16Ptr(item.IconPath)
		hr = pShellLink.SetIconLocation(iconPath, item.IconIndex)
		if hr != w32.S_OK {
			pShellLink.Release()
			return nil, fmt.Errorf("SetIconLocation failed: %v", hr)
		}
	}

	// Set title through property store
	if item.Title != "" {
		var pPropertyStore *IPropertyStore
		hr = pShellLink.QueryInterface(IID_IPropertyStoreGUID, &pPropertyStore)
		if hr == w32.S_OK {
			defer pPropertyStore.Release()

			var propVar PROPVARIANT
			propVar.Vt = 31 // VT_LPWSTR
			titlePtr := w32.MustStringToUTF16Ptr(item.Title)
			*(*uintptr)(unsafe.Pointer(&propVar.Val[0])) = uintptr(unsafe.Pointer(titlePtr))
			hr = pPropertyStore.SetValue(&PKEY_Title, &propVar)
			if hr == w32.S_OK {
				pPropertyStore.Commit()
			}
		}
	}

	return pShellLink, nil
}

func CoCreateInstance(rclsid *w32.GUID, pUnkOuter unsafe.Pointer, dwClsContext uint32, riid *w32.GUID, ppv interface{}) w32.HRESULT {
	var ret uintptr
	switch v := ppv.(type) {
	case **ICustomDestinationList:
		ret, _, _ = procCoCreateInstance.Call(
			uintptr(unsafe.Pointer(rclsid)),
			uintptr(pUnkOuter),
			uintptr(dwClsContext),
			uintptr(unsafe.Pointer(riid)),
			uintptr(unsafe.Pointer(v)),
		)
	case **IShellLink:
		ret, _, _ = procCoCreateInstance.Call(
			uintptr(unsafe.Pointer(rclsid)),
			uintptr(pUnkOuter),
			uintptr(dwClsContext),
			uintptr(unsafe.Pointer(riid)),
			uintptr(unsafe.Pointer(v)),
		)
	case **IObjectCollection:
		ret, _, _ = procCoCreateInstance.Call(
			uintptr(unsafe.Pointer(rclsid)),
			uintptr(pUnkOuter),
			uintptr(dwClsContext),
			uintptr(unsafe.Pointer(riid)),
			uintptr(unsafe.Pointer(v)),
		)
	default:
		panic("invalid type for CoCreateInstance")
	}
	return w32.HRESULT(ret)
}

// ICustomDestinationList methods
func (p *ICustomDestinationList) SetAppID(pszAppID *uint16) w32.HRESULT {
	ret, _, _ := syscall.Syscall(p.lpVtbl.SetAppID, 2, uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(pszAppID)), 0)
	return w32.HRESULT(ret)
}

func (p *ICustomDestinationList) BeginList(pcMinSlots *uint32, riid *w32.GUID, ppv *uintptr) w32.HRESULT {
	ret, _, _ := syscall.Syscall6(p.lpVtbl.BeginList, 4, uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(pcMinSlots)), uintptr(unsafe.Pointer(riid)), uintptr(unsafe.Pointer(ppv)), 0, 0)
	return w32.HRESULT(ret)
}

func (p *ICustomDestinationList) AppendCategory(pszCategory *uint16, poa *IObjectCollection) w32.HRESULT {
	ret, _, _ := syscall.Syscall(p.lpVtbl.AppendCategory, 3, uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(pszCategory)), uintptr(unsafe.Pointer(poa)))
	return w32.HRESULT(ret)
}

func (p *ICustomDestinationList) AddUserTasks(poa *IObjectCollection) w32.HRESULT {
	ret, _, _ := syscall.Syscall(p.lpVtbl.AddUserTasks, 2, uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(poa)), 0)
	return w32.HRESULT(ret)
}

func (p *ICustomDestinationList) CommitList() w32.HRESULT {
	ret, _, _ := syscall.Syscall(p.lpVtbl.CommitList, 1, uintptr(unsafe.Pointer(p)), 0, 0)
	return w32.HRESULT(ret)
}

func (p *ICustomDestinationList) DeleteList(pszAppID *uint16) w32.HRESULT {
	ret, _, _ := syscall.Syscall(p.lpVtbl.DeleteList, 2, uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(pszAppID)), 0)
	return w32.HRESULT(ret)
}

func (p *ICustomDestinationList) AbortList() w32.HRESULT {
	ret, _, _ := syscall.Syscall(p.lpVtbl.AbortList, 1, uintptr(unsafe.Pointer(p)), 0, 0)
	return w32.HRESULT(ret)
}

func (p *ICustomDestinationList) Release() uint32 {
	ret, _, _ := syscall.Syscall(p.lpVtbl.Release, 1, uintptr(unsafe.Pointer(p)), 0, 0)
	return uint32(ret)
}

// IShellLink methods
func (p *IShellLink) QueryInterface(riid *w32.GUID, ppvObject interface{}) w32.HRESULT {
	var ret uintptr
	switch v := ppvObject.(type) {
	case **IPropertyStore:
		ret, _, _ = syscall.Syscall(p.lpVtbl.QueryInterface, 3, uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(riid)), uintptr(unsafe.Pointer(v)))
	default:
		panic("invalid type for QueryInterface")
	}
	return w32.HRESULT(ret)
}

func (p *IShellLink) SetPath(pszFile *uint16) w32.HRESULT {
	ret, _, _ := syscall.Syscall(p.lpVtbl.SetPath, 2, uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(pszFile)), 0)
	return w32.HRESULT(ret)
}

func (p *IShellLink) SetArguments(pszArgs *uint16) w32.HRESULT {
	ret, _, _ := syscall.Syscall(p.lpVtbl.SetArguments, 2, uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(pszArgs)), 0)
	return w32.HRESULT(ret)
}

func (p *IShellLink) SetDescription(pszName *uint16) w32.HRESULT {
	ret, _, _ := syscall.Syscall(p.lpVtbl.SetDescription, 2, uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(pszName)), 0)
	return w32.HRESULT(ret)
}

func (p *IShellLink) SetIconLocation(pszIconPath *uint16, iIcon int) w32.HRESULT {
	ret, _, _ := syscall.Syscall(p.lpVtbl.SetIconLocation, 3, uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(pszIconPath)), uintptr(iIcon))
	return w32.HRESULT(ret)
}

func (p *IShellLink) Release() uint32 {
	ret, _, _ := syscall.Syscall(p.lpVtbl.Release, 1, uintptr(unsafe.Pointer(p)), 0, 0)
	return uint32(ret)
}

// IPropertyStore methods
func (p *IPropertyStore) SetValue(key *PROPERTYKEY, propvar *PROPVARIANT) w32.HRESULT {
	ret, _, _ := syscall.Syscall(p.lpVtbl.SetValue, 3, uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(key)), uintptr(unsafe.Pointer(propvar)))
	return w32.HRESULT(ret)
}

func (p *IPropertyStore) Commit() w32.HRESULT {
	ret, _, _ := syscall.Syscall(p.lpVtbl.Commit, 1, uintptr(unsafe.Pointer(p)), 0, 0)
	return w32.HRESULT(ret)
}

func (p *IPropertyStore) Release() uint32 {
	ret, _, _ := syscall.Syscall(p.lpVtbl.Release, 1, uintptr(unsafe.Pointer(p)), 0, 0)
	return uint32(ret)
}

// IObjectCollection methods
func (p *IObjectCollection) AddObject(punk interface{}) w32.HRESULT {
	var punkPtr uintptr
	switch v := punk.(type) {
	case *IShellLink:
		punkPtr = uintptr(unsafe.Pointer(v))
	default:
		panic("invalid type for AddObject")
	}
	ret, _, _ := syscall.Syscall(p.lpVtbl.AddObject, 2, uintptr(unsafe.Pointer(p)), punkPtr, 0)
	return w32.HRESULT(ret)
}

func (p *IObjectCollection) Release() uint32 {
	ret, _, _ := syscall.Syscall(p.lpVtbl.Release, 1, uintptr(unsafe.Pointer(p)), 0, 0)
	return uint32(ret)
}