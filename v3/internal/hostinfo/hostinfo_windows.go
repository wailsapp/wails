package hostinfo

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func Processors() ([]Processor, error) {
	rows, err := queryWMI("Win32_Processor", "Name", "NumberOfCores")
	if err != nil {
		return nil, err
	}
	result := make([]Processor, 0, len(rows))
	for _, row := range rows {
		cores, _ := strconv.Atoi(row["NumberOfCores"])
		result = append(result, Processor{Model: row["Name"], Cores: cores})
	}
	return result, nil
}

func Memory() (uint64, error) {
	rows, err := queryWMI("Win32_ComputerSystem", "TotalPhysicalMemory")
	if err != nil {
		return 0, err
	}
	if len(rows) == 0 {
		return 0, fmt.Errorf("no computer system returned by WMI")
	}
	return strconv.ParseUint(rows[0]["TotalPhysicalMemory"], 10, 64)
}

func Graphics() ([]GraphicsCard, error) {
	rows, err := queryWMI("Win32_VideoController", "Name", "AdapterCompatibility", "DriverVersion")
	if err != nil {
		return nil, err
	}
	result := make([]GraphicsCard, 0, len(rows))
	for _, row := range rows {
		result = append(result, GraphicsCard{
			Name: row["Name"], Vendor: row["AdapterCompatibility"], Driver: row["DriverVersion"],
		})
	}
	return result, nil
}

func queryWMI(class string, columns ...string) ([]map[string]string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
	if err != nil {
		code, ok := err.(*ole.OleError)
		if !ok || (code.Code() != 1 && code.Code() != 0x80010106) {
			return nil, err
		}
		// RPC_E_CHANGED_MODE means the caller already owns an apartment.
		if code.Code() == 1 {
			defer ole.CoUninitialize()
		}
	} else {
		defer ole.CoUninitialize()
	}
	object, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		return nil, err
	}
	defer object.Release()
	locator, err := object.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil, err
	}
	defer locator.Release()
	connection, err := oleutil.CallMethod(locator, "ConnectServer", ".", `root\cimv2`)
	if err != nil {
		return nil, err
	}
	defer connection.Clear()
	objects, err := oleutil.CallMethod(connection.ToIDispatch(), "ExecQuery", "SELECT "+strings.Join(columns, ",")+" FROM "+class)
	if err != nil {
		return nil, err
	}
	defer objects.Clear()
	var rows []map[string]string
	err = oleutil.ForEach(objects.ToIDispatch(), func(item *ole.VARIANT) error {
		defer item.Clear()
		row := make(map[string]string, len(columns))
		for _, column := range columns {
			value, err := oleutil.GetProperty(item.ToIDispatch(), column)
			if err != nil {
				return err
			}
			if value.Value() != nil {
				row[column] = fmt.Sprint(value.Value())
			}
			value.Clear()
		}
		rows = append(rows, row)
		return nil
	})
	return rows, err
}
