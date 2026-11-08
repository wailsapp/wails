//go:build windows

package combridge

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"golang.org/x/sys/windows"
)

// Vtables are registered at package init, so the GUID helpers here are done in
// Go: windows.GUIDFromString and windows.GUID.String call ole32.dll, which
// loaded ole32 (and user32 and gdi32 with it) into every program that links
// this package.

// parseGUID parses a GUID written as {XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX}.
func parseGUID(s string) (windows.GUID, error) {
	var g windows.GUID
	if len(s) != 38 || s[0] != '{' || s[37] != '}' || s[9] != '-' || s[14] != '-' || s[19] != '-' || s[24] != '-' {
		return g, fmt.Errorf("invalid GUID %q", s)
	}
	b, err := hex.DecodeString(s[1:9] + s[10:14] + s[15:19] + s[20:24] + s[25:37])
	if err != nil {
		return g, fmt.Errorf("invalid GUID %q: %w", s, err)
	}
	g.Data1 = binary.BigEndian.Uint32(b[0:4])
	g.Data2 = binary.BigEndian.Uint16(b[4:6])
	g.Data3 = binary.BigEndian.Uint16(b[6:8])
	copy(g.Data4[:], b[8:16])
	return g, nil
}

// formatGUID writes g as {XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX}, as
// windows.GUID.String does, without ole32.dll (which that calls).
func formatGUID(g windows.GUID) string {
	return fmt.Sprintf("{%08X-%04X-%04X-%02X%02X-%02X%02X%02X%02X%02X%02X}",
		g.Data1, g.Data2, g.Data3, g.Data4[0], g.Data4[1],
		g.Data4[2], g.Data4[3], g.Data4[4], g.Data4[5], g.Data4[6], g.Data4[7])
}
