//go:build windows

/*
 * Copyright (C) 2019 Tad Vizbaras. All Rights Reserved.
 * Copyright (C) 2010-2012 The W32 Authors. All Rights Reserved.
 */
package w32

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type CSIDL uint32

const (
	CSIDL_DESKTOP                 = 0x00
	CSIDL_INTERNET                = 0x01
	CSIDL_PROGRAMS                = 0x02
	CSIDL_CONTROLS                = 0x03
	CSIDL_PRINTERS                = 0x04
	CSIDL_PERSONAL                = 0x05
	CSIDL_FAVORITES               = 0x06
	CSIDL_STARTUP                 = 0x07
	CSIDL_RECENT                  = 0x08
	CSIDL_SENDTO                  = 0x09
	CSIDL_BITBUCKET               = 0x0A
	CSIDL_STARTMENU               = 0x0B
	CSIDL_MYDOCUMENTS             = 0x0C
	CSIDL_MYMUSIC                 = 0x0D
	CSIDL_MYVIDEO                 = 0x0E
	CSIDL_DESKTOPDIRECTORY        = 0x10
	CSIDL_DRIVES                  = 0x11
	CSIDL_NETWORK                 = 0x12
	CSIDL_NETHOOD                 = 0x13
	CSIDL_FONTS                   = 0x14
	CSIDL_TEMPLATES               = 0x15
	CSIDL_COMMON_STARTMENU        = 0x16
	CSIDL_COMMON_PROGRAMS         = 0x17
	CSIDL_COMMON_STARTUP          = 0x18
	CSIDL_COMMON_DESKTOPDIRECTORY = 0x19
	CSIDL_APPDATA                 = 0x1A
	CSIDL_PRINTHOOD               = 0x1B
	CSIDL_LOCAL_APPDATA           = 0x1C
	CSIDL_ALTSTARTUP              = 0x1D
	CSIDL_COMMON_ALTSTARTUP       = 0x1E
	CSIDL_COMMON_FAVORITES        = 0x1F
	CSIDL_INTERNET_CACHE          = 0x20
	CSIDL_COOKIES                 = 0x21
	CSIDL_HISTORY                 = 0x22
	CSIDL_COMMON_APPDATA          = 0x23
	CSIDL_WINDOWS                 = 0x24
	CSIDL_SYSTEM                  = 0x25
	CSIDL_PROGRAM_FILES           = 0x26
	CSIDL_MYPICTURES              = 0x27
	CSIDL_PROFILE                 = 0x28
	CSIDL_SYSTEMX86               = 0x29
	CSIDL_PROGRAM_FILESX86        = 0x2A
	CSIDL_PROGRAM_FILES_COMMON    = 0x2B
	CSIDL_PROGRAM_FILES_COMMONX86 = 0x2C
	CSIDL_COMMON_TEMPLATES        = 0x2D
	CSIDL_COMMON_DOCUMENTS        = 0x2E
	CSIDL_COMMON_ADMINTOOLS       = 0x2F
	CSIDL_ADMINTOOLS              = 0x30
	CSIDL_CONNECTIONS             = 0x31
	CSIDL_COMMON_MUSIC            = 0x35
	CSIDL_COMMON_PICTURES         = 0x36
	CSIDL_COMMON_VIDEO            = 0x37
	CSIDL_RESOURCES               = 0x38
	CSIDL_RESOURCES_LOCALIZED     = 0x39
	CSIDL_COMMON_OEM_LINKS        = 0x3A
	CSIDL_CDBURN_AREA             = 0x3B
	CSIDL_COMPUTERSNEARME         = 0x3D
	CSIDL_FLAG_CREATE             = 0x8000
	CSIDL_FLAG_DONT_VERIFY        = 0x4000
	CSIDL_FLAG_NO_ALIAS           = 0x1000
	CSIDL_FLAG_PER_USER_INIT      = 0x8000
	CSIDL_FLAG_MASK               = 0xFF00

	NOTIFYICON_VERSION = 4
)

var (
	FOLDERID_AccountPictures        = NewGUID("{008CA0B1-55B4-4C56-B8A8-4DE4B299D3BE}")
	FOLDERID_AddNewPrograms         = NewGUID("{DE61D971-5EBC-4F02-A3A9-6C82895E5C04}")
	FOLDERID_AdminTools             = NewGUID("{724EF170-A42D-4FEF-9F26-B60E846FBA4F}")
	FOLDERID_ApplicationShortcuts   = NewGUID("{A3918781-E5F2-4890-B3D9-A7E54332328C}")
	FOLDERID_AppsFolder             = NewGUID("{1E87508D-89C2-42F0-8A7E-645A0F50CA58}")
	FOLDERID_AppUpdates             = NewGUID("{A305CE99-F527-492B-8B1A-7E76FA98D6E4}")
	FOLDERID_CDBurning              = NewGUID("{9E52AB10-F80D-49DF-ACB8-4330F5687855}")
	FOLDERID_ChangeRemovePrograms   = NewGUID("{DF7266AC-9274-4867-8D55-3BD661DE872D}")
	FOLDERID_CommonAdminTools       = NewGUID("{D0384E7D-BAC3-4797-8F14-CBA229B392B5}")
	FOLDERID_CommonOEMLinks         = NewGUID("{C1BAE2D0-10DF-4334-BEDD-7AA20B227A9D}")
	FOLDERID_CommonPrograms         = NewGUID("{0139D44E-6AFE-49F2-8690-3DAFCAE6FFB8}")
	FOLDERID_CommonStartMenu        = NewGUID("{A4115719-D62E-491D-AA7C-E74B8BE3B067}")
	FOLDERID_CommonStartup          = NewGUID("{82A5EA35-D9CD-47C5-9629-E15D2F714E6E}")
	FOLDERID_CommonTemplates        = NewGUID("{B94237E7-57AC-4347-9151-B08C6C32D1F7}")
	FOLDERID_ComputerFolder         = NewGUID("{0AC0837C-BBF8-452A-850D-79D08E667CA7}")
	FOLDERID_ConflictFolder         = NewGUID("{4BFEFB45-347D-4006-A5BE-AC0CB0567192}")
	FOLDERID_ConnectionsFolder      = NewGUID("{6F0CD92B-2E97-45D1-88FF-B0D186B8DEDD}")
	FOLDERID_Contacts               = NewGUID("{56784854-C6CB-462B-8169-88E350ACB882}")
	FOLDERID_ControlPanelFolder     = NewGUID("{82A74AEB-AEB4-465C-A014-D097EE346D63}")
	FOLDERID_Cookies                = NewGUID("{2B0F765D-C0E9-4171-908E-08A611B84FF6}")
	FOLDERID_Desktop                = NewGUID("{B4BFCC3A-DB2C-424C-B029-7FE99A87C641}")
	FOLDERID_DeviceMetadataStore    = NewGUID("{5CE4A5E9-E4EB-479D-B89F-130C02886155}")
	FOLDERID_Documents              = NewGUID("{FDD39AD0-238F-46AF-ADB4-6C85480369C7}")
	FOLDERID_DocumentsLibrary       = NewGUID("{7B0DB17D-9CD2-4A93-9733-46CC89022E7C}")
	FOLDERID_Downloads              = NewGUID("{374DE290-123F-4565-9164-39C4925E467B}")
	FOLDERID_Favorites              = NewGUID("{1777F761-68AD-4D8A-87BD-30B759FA33DD}")
	FOLDERID_Fonts                  = NewGUID("{FD228CB7-AE11-4AE3-864C-16F3910AB8FE}")
	FOLDERID_Games                  = NewGUID("{CAC52C1A-B53D-4EDC-92D7-6B2E8AC19434}")
	FOLDERID_GameTasks              = NewGUID("{054FAE61-4DD8-4787-80B6-090220C4B700}")
	FOLDERID_History                = NewGUID("{D9DC8A3B-B784-432E-A781-5A1130A75963}")
	FOLDERID_HomeGroup              = NewGUID("{52528A6B-B9E3-4ADD-B60D-588C2DBA842D}")
	FOLDERID_HomeGroupCurrentUser   = NewGUID("{9B74B6A3-0DFD-4F11-9E78-5F7800F2E772}")
	FOLDERID_ImplicitAppShortcuts   = NewGUID("{BCB5256F-79F6-4CEE-B725-DC34E402FD46}")
	FOLDERID_InternetCache          = NewGUID("{352481E8-33BE-4251-BA85-6007CAEDCF9D}")
	FOLDERID_InternetFolder         = NewGUID("{4D9F7874-4E0C-4904-967B-40B0D20C3E4B}")
	FOLDERID_Libraries              = NewGUID("{1B3EA5DC-B587-4786-B4EF-BD1DC332AEAE}")
	FOLDERID_Links                  = NewGUID("{BFB9D5E0-C6A9-404C-B2B2-AE6DB6AF4968}")
	FOLDERID_LocalAppData           = NewGUID("{F1B32785-6FBA-4FCF-9D55-7B8E7F157091}")
	FOLDERID_LocalAppDataLow        = NewGUID("{A520A1A4-1780-4FF6-BD18-167343C5AF16}")
	FOLDERID_LocalizedResourcesDir  = NewGUID("{2A00375E-224C-49DE-B8D1-440DF7EF3DDC}")
	FOLDERID_Music                  = NewGUID("{4BD8D571-6D19-48D3-BE97-422220080E43}")
	FOLDERID_MusicLibrary           = NewGUID("{2112AB0A-C86A-4FFE-A368-0DE96E47012E}")
	FOLDERID_NetHood                = NewGUID("{C5ABBF53-E17F-4121-8900-86626FC2C973}")
	FOLDERID_NetworkFolder          = NewGUID("{D20BEEC4-5CA8-4905-AE3B-BF251EA09B53}")
	FOLDERID_OriginalImages         = NewGUID("{2C36C0AA-5812-4B87-BFD0-4CD0DFB19B39}")
	FOLDERID_PhotoAlbums            = NewGUID("{69D2CF90-FC33-4FB7-9A0C-EBB0F0FCB43C}")
	FOLDERID_Pictures               = NewGUID("{33E28130-4E1E-4676-835A-98395C3BC3BB}")
	FOLDERID_PicturesLibrary        = NewGUID("{A990AE9F-A03B-4E80-94BC-9912D7504104}")
	FOLDERID_Playlists              = NewGUID("{DE92C1C7-837F-4F69-A3BB-86E631204A23}")
	FOLDERID_PrintersFolder         = NewGUID("{76FC4E2D-D6AD-4519-A663-37BD56068185}")
	FOLDERID_PrintHood              = NewGUID("{9274BD8D-CFD1-41C3-B35E-B13F55A758F4}")
	FOLDERID_Profile                = NewGUID("{5E6C858F-0E22-4760-9AFE-EA3317B67173}")
	FOLDERID_ProgramData            = NewGUID("{62AB5D82-FDC1-4DC3-A9DD-070D1D495D97}")
	FOLDERID_ProgramFiles           = NewGUID("{905E63B6-C1BF-494E-B29C-65B732D3D21A}")
	FOLDERID_ProgramFilesCommon     = NewGUID("{F7F1ED05-9F6D-47A2-AAAE-29D317C6F066}")
	FOLDERID_ProgramFilesCommonX64  = NewGUID("{6365D5A7-0F0D-45E5-87F6-0DA56B6A4F7D}")
	FOLDERID_ProgramFilesCommonX86  = NewGUID("{DE974D24-D9C6-4D3E-BF91-F4455120B917}")
	FOLDERID_ProgramFilesX64        = NewGUID("{6D809377-6AF0-444B-8957-A3773F02200E}")
	FOLDERID_ProgramFilesX86        = NewGUID("{7C5A40EF-A0FB-4BFC-874A-C0F2E0B9FA8E}")
	FOLDERID_Programs               = NewGUID("{A77F5D77-2E2B-44C3-A6A2-ABA601054A51}")
	FOLDERID_Public                 = NewGUID("{DFDF76A2-C82A-4D63-906A-5644AC457385}")
	FOLDERID_PublicDesktop          = NewGUID("{C4AA340D-F20F-4863-AFEF-1F769F2BE730}")
	FOLDERID_PublicDocuments        = NewGUID("{ED4824AF-DCE4-45A8-81E2-FC7965083634}")
	FOLDERID_PublicDownloads        = NewGUID("{3D644C9B-1FB8-4F30-9B45-F670235F79C0}")
	FOLDERID_PublicGameTasks        = NewGUID("{DEBF2536-E1A8-4C59-B6A2-414586476AEA}")
	FOLDERID_PublicLibraries        = NewGUID("{48DAF80B-E6CF-4F4E-B800-0E69D84EE384}")
	FOLDERID_PublicMusic            = NewGUID("{3214FAB5-9757-4298-BB61-92A9DEAA44FF}")
	FOLDERID_PublicPictures         = NewGUID("{B6EBFB86-6907-413C-9AF7-4FC2ABF07CC5}")
	FOLDERID_PublicRingtones        = NewGUID("{E555AB60-153B-4D17-9F04-A5FE99FC15EC}")
	FOLDERID_PublicUserTiles        = NewGUID("{0482af6c-08f1-4c34-8c90-e17ec98b1e17}")
	FOLDERID_PublicVideos           = NewGUID("{2400183A-6185-49FB-A2D8-4A392A602BA3}")
	FOLDERID_QuickLaunch            = NewGUID("{52a4f021-7b75-48a9-9f6b-4b87a210bc8f}")
	FOLDERID_Recent                 = NewGUID("{AE50C081-EBD2-438A-8655-8A092E34987A}")
	FOLDERID_RecordedTVLibrary      = NewGUID("{1A6FDBA2-F42D-4358-A798-B74D745926C5}")
	FOLDERID_RecycleBinFolder       = NewGUID("{B7534046-3ECB-4C18-BE4E-64CD4CB7D6AC}")
	FOLDERID_ResourceDir            = NewGUID("{8AD10C31-2ADB-4296-A8F7-E4701232C972}")
	FOLDERID_Ringtones              = NewGUID("{C870044B-F49E-4126-A9C3-B52A1FF411E8}")
	FOLDERID_RoamingAppData         = NewGUID("{3EB685DB-65F9-4CF6-A03A-E3EF65729F3D}")
	FOLDERID_RoamingTiles           = NewGUID("{AAA8D5A5-F1D6-4259-BAA8-78E7EF60835E}")
	FOLDERID_SampleMusic            = NewGUID("{B250C668-F57D-4EE1-A63C-290EE7D1AA1F}")
	FOLDERID_SamplePictures         = NewGUID("{C4900540-2379-4C75-844B-64E6FAF8716B}")
	FOLDERID_SamplePlaylists        = NewGUID("{15CA69B3-30EE-49C1-ACE1-6B5EC372AFB5}")
	FOLDERID_SampleVideos           = NewGUID("{859EAD94-2E85-48AD-A71A-0969CB56A6CD}")
	FOLDERID_SavedGames             = NewGUID("{4C5C32FF-BB9D-43B0-B5B4-2D72E54EAAA4}")
	FOLDERID_SavedPictures          = NewGUID("{3B193882-D3AD-4EAB-965A-69829D1FB59F}")
	FOLDERID_SavedPicturesLibrary   = NewGUID("{E25B5812-BE88-4BD9-94B0-29233477B6C3}")
	FOLDERID_SavedSearches          = NewGUID("{7D1D3A04-DEBB-4115-95CF-2F29DA2920DA}")
	FOLDERID_SEARCH_CSC             = NewGUID("{ee32e446-31ca-4aba-814f-a5ebd2fd6d5e}")
	FOLDERID_SEARCH_MAPI            = NewGUID("{98ec0e18-2098-4d44-8644-66979315a281}")
	FOLDERID_SearchHome             = NewGUID("{190337d1-b8ca-4121-a639-6d472d16972a}")
	FOLDERID_SendTo                 = NewGUID("{8983036C-27C0-404B-8F08-102D10DCFD74}")
	FOLDERID_SidebarDefaultParts    = NewGUID("{7B396E54-9EC5-4300-BE0A-2482EBAE1A26}")
	FOLDERID_SidebarParts           = NewGUID("{A75D362E-50FC-4fb7-AC2C-A8BEAA314493}")
	FOLDERID_SkyDrive               = NewGUID("{A52BBA46-E9E1-435f-B3D9-28DAA648C0F6}")
	FOLDERID_SkyDriveCameraRoll     = NewGUID("{767E6811-49CB-4273-87C2-20F355E1085B}")
	FOLDERID_SkyDriveDocuments      = NewGUID("{24D89E24-2F19-4534-9DDE-6A6671FBB8FE}")
	FOLDERID_SkyDriveMusic          = NewGUID("{C3F2459E-80D6-45DC-BFEF-1F769F2BE730}")
	FOLDERID_SkyDrivePictures       = NewGUID("{339719B5-8C47-4894-94C2-D8F77ADD44A6}")
	FOLDERID_StartMenu              = NewGUID("{625B53C3-AB48-4EC1-BA1F-A1EF4146FC19}")
	FOLDERID_Startup                = NewGUID("{B97D20BB-F46A-4C97-BA10-5E3608430854}")
	FOLDERID_SyncManagerFolder      = NewGUID("{43668BF8-C14E-49B2-97C9-747784D784B7}")
	FOLDERID_SyncResultsFolder      = NewGUID("{289a9a43-be44-4057-a41b-587a76d7e7f9}")
	FOLDERID_SyncSetupFolder        = NewGUID("{0F214138-B1D3-4a90-BBA9-27CBC0C5389A}")
	FOLDERID_System                 = NewGUID("{1AC14E77-02E7-4E5D-B744-2EB1AE5198B7}")
	FOLDERID_SystemX86              = NewGUID("{D65231B0-B2F1-4857-A4CE-A8E7C6EA7D27}")
	FOLDERID_Templates              = NewGUID("{A63293E8-664E-48DB-A079-DF759E0509F7}")
	FOLDERID_UserPinned             = NewGUID("{9E3995AB-1F9C-4F13-B827-48B24B6C7174}")
	FOLDERID_UserProfiles           = NewGUID("{0762D272-C50A-4BB0-A382-697DCD729B80}")
	FOLDERID_UserProgramFiles       = NewGUID("{5CD7AEE2-2219-4A67-B85D-6C9CE15660CB}")
	FOLDERID_UserProgramFilesCommon = NewGUID("{BCBD3057-CA5C-4622-B42D-BC56DB0AE516}")
	FOLDERID_UsersFiles             = NewGUID("{F3CE0F7C-4901-4ACC-8648-D5D44B04EF8F}")
	FOLDERID_UsersLibraries         = NewGUID("{A302545D-DEFF-464b-ABE8-61C8648D939B}")
	FOLDERID_Videos                 = NewGUID("{18989B1D-99B5-455B-841C-AB7C74E4DDFC}")
	FOLDERID_VideosLibrary          = NewGUID("{491E922F-5643-4AF4-A7EB-4E7A138D8174}")
	FOLDERID_Windows                = NewGUID("{F38BF404-1D43-42F2-9305-67DE0B28FC23}")
)

var (
	modshell32 = syscall.NewLazyDLL("shell32.dll")

	procSHBrowseForFolder      = modshell32.NewProc("SHBrowseForFolderW")
	procSHGetPathFromIDList    = modshell32.NewProc("SHGetPathFromIDListW")
	procDragAcceptFiles        = modshell32.NewProc("DragAcceptFiles")
	procDragQueryFile          = modshell32.NewProc("DragQueryFileW")
	procDragQueryPoint         = modshell32.NewProc("DragQueryPoint")
	procDragFinish             = modshell32.NewProc("DragFinish")
	procShellExecute           = modshell32.NewProc("ShellExecuteW")
	procExtractIcon            = modshell32.NewProc("ExtractIconW")
	procGetSpecialFolderPath   = modshell32.NewProc("SHGetSpecialFolderPathW")
	procShellNotifyIcon        = modshell32.NewProc("Shell_NotifyIconW")
	procShellNotifyIconGetRect = modshell32.NewProc("Shell_NotifyIconGetRect")
	procSHGetKnownFolderPath   = modshell32.NewProc("SHGetKnownFolderPath")
	procSHAppBarMessage        = modshell32.NewProc("SHAppBarMessage")
)

type APPBARDATA struct {
	CbSize           uint32
	HWnd             HWND
	UCallbackMessage uint32
	UEdge            uint32
	Rc               RECT
	LParam           uintptr
}

func ShellNotifyIcon(cmd uintptr, nid *NOTIFYICONDATA) bool {
	ret, _, _ := procShellNotifyIcon.Call(cmd, uintptr(unsafe.Pointer(nid)))
	return ret == 1
}

func SHBrowseForFolder(bi *BROWSEINFO) uintptr {
	ret, _, _ := procSHBrowseForFolder.Call(uintptr(unsafe.Pointer(bi)))

	return ret
}

func SHGetKnownFolderPath(rfid *GUID, dwFlags uint32, hToken HANDLE) (string, error) {
	var path *uint16
	ret, _, _ := procSHGetKnownFolderPath.Call(uintptr(unsafe.Pointer(rfid)), uintptr(dwFlags), hToken, uintptr(unsafe.Pointer(path)))
	if ret != uintptr(windows.S_OK) {
		return "", fmt.Errorf("SHGetKnownFolderPath failed: %v", ret)
	}
	return windows.UTF16PtrToString(path), nil
}

func SHGetPathFromIDList(idl uintptr) string {
	buf := make([]uint16, 1024)
	procSHGetPathFromIDList.Call(
		idl,
		uintptr(unsafe.Pointer(&buf[0])))

	return syscall.UTF16ToString(buf)
}

func DragAcceptFiles(hwnd HWND, accept bool) {
	procDragAcceptFiles.Call(
		hwnd,
		uintptr(BoolToBOOL(accept)))
}

func DragQueryFile(hDrop HDROP, iFile uint) (fileName string, fileCount uint) {
	ret, _, _ := procDragQueryFile.Call(
		hDrop,
		uintptr(iFile),
		0,
		0)

	fileCount = uint(ret)

	if iFile != 0xFFFFFFFF {
		buf := make([]uint16, fileCount+1)

		ret, _, _ := procDragQueryFile.Call(
			hDrop,
			uintptr(iFile),
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(fileCount+1))

		if ret == 0 {
			panic("Invoke DragQueryFile error.")
		}

		fileName = syscall.UTF16ToString(buf)
	}

	return
}

func DragQueryPoint(hDrop HDROP) (x, y int, isClientArea bool) {
	var pt POINT
	ret, _, _ := procDragQueryPoint.Call(
		uintptr(hDrop),
		uintptr(unsafe.Pointer(&pt)))

	return int(pt.X), int(pt.Y), (ret == 1)
}

func DragFinish(hDrop HDROP) {
	procDragFinish.Call(uintptr(hDrop))
}

func ShellExecute(hwnd HWND, lpOperation, lpFile, lpParameters, lpDirectory string, nShowCmd int) error {
	var op, param, directory uintptr
	if len(lpOperation) != 0 {
		op = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpOperation)))
	}
	if len(lpParameters) != 0 {
		param = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpParameters)))
	}
	if len(lpDirectory) != 0 {
		directory = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpDirectory)))
	}

	ret, _, _ := procShellExecute.Call(
		uintptr(hwnd),
		op,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpFile))),
		param,
		directory,
		uintptr(nShowCmd))

	errorMsg := ""
	if ret != 0 && ret <= 32 {
		switch int(ret) {
		case ERROR_FILE_NOT_FOUND:
			errorMsg = "The specified file was not found."
		case ERROR_PATH_NOT_FOUND:
			errorMsg = "The specified path was not found."
		case ERROR_BAD_FORMAT:
			errorMsg = "The .exe file is invalid (non-Win32 .exe or error in .exe image)."
		case SE_ERR_ACCESSDENIED:
			errorMsg = "The operating system denied access to the specified file."
		case SE_ERR_ASSOCINCOMPLETE:
			errorMsg = "The file name association is incomplete or invalid."
		case SE_ERR_DDEBUSY:
			errorMsg = "The DDE transaction could not be completed because other DDE transactions were being processed."
		case SE_ERR_DDEFAIL:
			errorMsg = "The DDE transaction failed."
		case SE_ERR_DDETIMEOUT:
			errorMsg = "The DDE transaction could not be completed because the request timed out."
		case SE_ERR_DLLNOTFOUND:
			errorMsg = "The specified DLL was not found."
		case SE_ERR_NOASSOC:
			errorMsg = "There is no application associated with the given file name extension. This error will also be returned if you attempt to print a file that is not printable."
		case SE_ERR_OOM:
			errorMsg = "There was not enough memory to complete the operation."
		case SE_ERR_SHARE:
			errorMsg = "A sharing violation occurred."
		default:
			errorMsg = fmt.Sprintf("Unknown error occurred with error code %v", ret)
		}
	} else {
		return nil
	}

	return errors.New(errorMsg)
}

func ExtractIcon(lpszExeFileName string, nIconIndex int) HICON {
	ret, _, _ := procExtractIcon.Call(
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpszExeFileName))),
		uintptr(nIconIndex))

	return HICON(ret)
}

func SHGetSpecialFolderPath(hwndOwner HWND, lpszPath *uint16, csidl CSIDL, fCreate bool) bool {
	ret, _, _ := procGetSpecialFolderPath.Call(
		uintptr(hwndOwner),
		uintptr(unsafe.Pointer(lpszPath)),
		uintptr(csidl),
		uintptr(BoolToBOOL(fCreate)),
		0,
		0)

	return ret != 0
}

func GetSystrayBounds(hwnd HWND, uid uint32) (*RECT, error) {
	var rect RECT
	identifier := NOTIFYICONIDENTIFIER{
		CbSize: uint32(unsafe.Sizeof(NOTIFYICONIDENTIFIER{})),
		HWnd:   hwnd,
		UId:    uid,
	}
	ret, _, _ := procShellNotifyIconGetRect.Call(
		uintptr(unsafe.Pointer(&identifier)),
		uintptr(unsafe.Pointer(&rect)))

	if ret != S_OK {
		return nil, syscall.GetLastError()
	}

	return &rect, nil
}

// GetTaskbarPosition returns the location of the taskbar.
func GetTaskbarPosition() *APPBARDATA {
	var result APPBARDATA
	result.CbSize = uint32(unsafe.Sizeof(APPBARDATA{}))
	ret, _, _ := procSHAppBarMessage.Call(
		ABM_GETTASKBARPOS,
		uintptr(unsafe.Pointer(&result)))
	if ret == 0 {
		return nil
	}

	return &result
}
