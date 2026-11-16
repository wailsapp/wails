package windowsign

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"context"
	"encoding/binary"
	"errors"
	"github.com/wailsapp/wails/v3/internal/windowspackage"
	"io"
	"os"
	"sort"
)

type zipRecord struct {
	name        string
	offset, end int64
	header      int
	compressed  uint64
	offsetIndex int
	wideOffset  bool
	directory   []byte
}

func packageRecords(source *os.File, size int64) ([]zipRecord, error) {
	start, err := zipDirectory(source, size)
	if err != nil {
		return nil, err
	}
	length := size - start - 22
	footer := make([]byte, 22)
	if _, err := source.ReadAt(footer, size-22); err != nil {
		return nil, err
	}
	if binary.LittleEndian.Uint32(footer[16:]) == 0xffffffff {
		length -= 76
	}
	if length < 0 || length > 64<<20 {
		return nil, errors.New("MSIX ZIP directory exceeds 64 MiB")
	}
	data := make([]byte, length)
	if _, err := source.ReadAt(data, start); err != nil {
		return nil, err
	}
	var records []zipRecord
	for len(data) > 0 {
		if len(data) < 46 || binary.LittleEndian.Uint32(data) != 0x02014b50 {
			return nil, errors.New("invalid MSIX central directory record")
		}
		nameLen := int(binary.LittleEndian.Uint16(data[28:]))
		entryLen := 46 + nameLen + int(binary.LittleEndian.Uint16(data[30:])) + int(binary.LittleEndian.Uint16(data[32:]))
		if entryLen > len(data) {
			return nil, errors.New("truncated MSIX central directory record")
		}
		record, err := readZIPRecord(bytes.Clone(data[:entryLen]), nameLen)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
		data = data[entryLen:]
	}
	sort.Slice(records, func(i, j int) bool { return records[i].offset < records[j].offset })
	for i := range records {
		record := &records[i]
		record.end = start
		if i+1 < len(records) {
			record.end = records[i+1].offset
		}
		if record.offset < 0 || record.end <= record.offset || record.end > start || (i == 0 && record.offset != 0) {
			return nil, errors.New("invalid MSIX local record offsets")
		}
		header := make([]byte, 30)
		if _, err := source.ReadAt(header, record.offset); err != nil {
			return nil, err
		}
		if binary.LittleEndian.Uint32(header) != 0x04034b50 {
			return nil, errors.New("invalid MSIX local file header")
		}
		nameLen := int(binary.LittleEndian.Uint16(header[26:]))
		extraLen := int(binary.LittleEndian.Uint16(header[28:]))
		headerSize := 30 + nameLen + extraLen
		record.header = headerSize
		if record.offset+int64(headerSize)+int64(record.compressed) > record.end {
			return nil, errors.New("MSIX local file exceeds its record")
		}
		name := make([]byte, nameLen)
		if _, err := source.ReadAt(name, record.offset+30); err != nil {
			return nil, err
		}
		if string(name) != record.name {
			return nil, errors.New("MSIX local and central file names differ")
		}
	}
	return records, nil
}
func unsignedMSIX(ctx context.Context, target, source *os.File, records []zipRecord, contentTypes []byte, replaceTypes bool) (int64, []byte, error) {
	var directory []byte
	count := 0
	for _, record := range records {
		if record.name == signatureName {
			continue
		}
		offset, err := target.Seek(0, io.SeekCurrent)
		if err != nil {
			return 0, nil, err
		}
		entry := record.directory
		if record.name == "[Content_Types].xml" && replaceTypes {
			local, central, err := footprintRecord(record.name, contentTypes)
			if err != nil {
				return 0, nil, err
			}
			if _, err := target.Write(local); err != nil {
				return 0, nil, err
			}
			entry = central
		} else {
			if _, err := io.Copy(target, contextReader{ctx, io.NewSectionReader(source, record.offset, record.end-record.offset)}); err != nil {
				return 0, nil, err
			}
		}
		if offset >= 1<<32 {
			return 0, nil, errors.New("MSIX signed records exceed 4 GiB")
		}
		if entry == nil {
			return 0, nil, errors.New("missing ZIP central directory entry")
		}
		if record.name == "[Content_Types].xml" && replaceTypes {
			binary.LittleEndian.PutUint32(entry[42:], uint32(offset))
		} else if record.wideOffset {
			binary.LittleEndian.PutUint64(entry[record.offsetIndex:], uint64(offset))
		} else {
			binary.LittleEndian.PutUint32(entry[42:], uint32(offset))
		}
		directory = append(directory, entry...)
		count++
	}
	start, err := target.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, nil, err
	}
	if err := finishMSIX(target, start, directory, count); err != nil {
		return 0, nil, err
	}
	return start, directory, nil
}
func footprintRecord(name string, data []byte) ([]byte, []byte, error) {
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	header, compressed, err := windowspackage.EncodeFootprint(name, data, flate.DefaultCompression)
	if err != nil {
		return nil, nil, err
	}
	// Native APPX verification requires this to agree with the ZIP64 footer.
	header.CreatorVersion = 45
	header.Flags = 0
	entry, err := writer.CreateRaw(header)
	if err != nil {
		return nil, nil, err
	}
	if _, err := entry.Write(compressed); err != nil {
		return nil, nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, nil, err
	}
	archive := buffer.Bytes()
	end := archive[len(archive)-22:]
	start := int(binary.LittleEndian.Uint32(end[16:]))
	return bytes.Clone(archive[:start]), bytes.Clone(archive[start : len(archive)-22]), nil
}
func finishMSIX(file *os.File, start int64, directory []byte, count int) error {
	if count > 65534 || start < 0 || start+int64(len(directory))+98 >= 1<<32 {
		return errors.New("MSIX ZIP exceeds portable signing size limits")
	}
	if _, err := file.Write(directory); err != nil {
		return err
	}
	footer := make([]byte, 98)
	binary.LittleEndian.PutUint32(footer, 0x06064b50)
	binary.LittleEndian.PutUint64(footer[4:], 44)
	binary.LittleEndian.PutUint16(footer[12:], 45)
	binary.LittleEndian.PutUint16(footer[14:], 45)
	binary.LittleEndian.PutUint64(footer[24:], uint64(count))
	binary.LittleEndian.PutUint64(footer[32:], uint64(count))
	binary.LittleEndian.PutUint64(footer[40:], uint64(len(directory)))
	binary.LittleEndian.PutUint64(footer[48:], uint64(start))
	binary.LittleEndian.PutUint32(footer[56:], 0x07064b50)
	binary.LittleEndian.PutUint64(footer[64:], uint64(start)+uint64(len(directory)))
	binary.LittleEndian.PutUint32(footer[72:], 1)
	binary.LittleEndian.PutUint32(footer[76:], 0x06054b50)
	binary.LittleEndian.PutUint16(footer[84:], 0xffff)
	binary.LittleEndian.PutUint16(footer[86:], 0xffff)
	binary.LittleEndian.PutUint32(footer[88:], 0xffffffff)
	binary.LittleEndian.PutUint32(footer[92:], 0xffffffff)
	_, err := file.Write(footer)
	return err
}

func readZIPRecord(data []byte, nameLen int) (zipRecord, error) {
	record := zipRecord{name: string(data[46 : 46+nameLen]), offset: int64(binary.LittleEndian.Uint32(data[42:])), compressed: uint64(binary.LittleEndian.Uint32(data[20:])), directory: data, offsetIndex: 42}
	uncompressed := uint64(binary.LittleEndian.Uint32(data[24:]))
	needsSize := uncompressed == 0xffffffff
	needsCompressed := record.compressed == 0xffffffff
	needsOffset := record.offset == 0xffffffff
	if !needsSize && !needsCompressed && !needsOffset {
		return record, nil
	}
	extraStart := 46 + nameLen
	extraEnd := extraStart + int(binary.LittleEndian.Uint16(data[30:]))
	found := false
	for position := extraStart; position < extraEnd; {
		if position+4 > extraEnd {
			return record, errors.New("truncated ZIP extra header")
		}
		kind := binary.LittleEndian.Uint16(data[position:])
		length := int(binary.LittleEndian.Uint16(data[position+2:]))
		position += 4
		if position+length > extraEnd {
			return record, errors.New("truncated ZIP extra data")
		}
		if kind == 1 {
			if found {
				return record, errors.New("duplicate ZIP64 metadata")
			}
			found = true
			cursor := position
			end := position + length
			read := func() (uint64, error) {
				if cursor+8 > end {
					return 0, errors.New("missing ZIP64 value")
				}
				value := binary.LittleEndian.Uint64(data[cursor:])
				cursor += 8
				return value, nil
			}
			var err error
			if needsSize {
				uncompressed, err = read()
				if err != nil {
					return record, err
				}
			}
			if needsCompressed {
				record.compressed, err = read()
				if err != nil {
					return record, err
				}
			}
			if needsOffset {
				record.offsetIndex = cursor
				record.wideOffset = true
				value, err := read()
				if err != nil {
					return record, err
				}
				if value >= 1<<32 {
					return record, errors.New("ZIP64 offset exceeds 4 GiB")
				}
				record.offset = int64(value)
			}
		}
		position += length
	}
	if !found || uncompressed >= 1<<32 || record.compressed >= 1<<32 {
		return record, errors.New("missing or oversized ZIP64 entry metadata")
	}
	return record, nil
}
