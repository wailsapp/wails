package applesign

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"os"
)

// SignDMG embeds a signature between a UDIF image's payload and trailing header.
// The payload is unchanged. The header is sealed in special slot 6 with its
// signature size zeroed, matching Apple's independently generated disk images.
func SignDMG(ctx context.Context, input, output string, o Options) (id CodeIdentity, err error) {
	if err := validateOptions(o); err != nil {
		return id, err
	}
	if len(o.Entitlements) != 0 || o.HardenedRuntime {
		return id, errors.New("DMG signatures do not take entitlements or hardened runtime")
	}
	f, err := os.Open(input)
	if err != nil {
		return id, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return id, err
	}
	if !st.Mode().IsRegular() || st.Size() < 512 {
		return id, errors.New("DMG must be a regular UDIF image")
	}
	var footer [512]byte
	if _, err := f.ReadAt(footer[:], st.Size()-512); err != nil {
		return id, err
	}
	be := binary.BigEndian
	if string(footer[:4]) != "koly" || be.Uint32(footer[4:]) != 4 || be.Uint32(footer[8:]) != 512 {
		return id, errors.New("unsupported UDIF trailer")
	}
	limit := uint64(st.Size() - 512)
	offset, length := be.Uint64(footer[296:]), be.Uint64(footer[304:])
	if offset != 0 || length != 0 {
		if offset == 0 || offset > limit || length != limit-offset {
			return id, errors.New("invalid existing DMG signature bounds")
		}
		limit = offset
	}
	if limit > 0xffffffff {
		return id, errors.New("DMG signing currently supports payloads below 4 GiB")
	}
	for _, pair := range [][2]int{{24, 32}, {40, 48}, {216, 224}} {
		start, span := be.Uint64(footer[pair[0]:]), be.Uint64(footer[pair[1]:])
		if start > limit || span > limit-start {
			return id, errors.New("UDIF data extent exceeds payload")
		}
	}
	be.PutUint64(footer[296:], limit)
	be.PutUint64(footer[304:], 0)
	o.representation = footer[:]
	components, err := signatureComponents(o)
	if err != nil {
		return id, err
	}
	out, err := os.OpenFile(output, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)
	if err != nil {
		return id, err
	}
	defer func() {
		err = errors.Join(err, out.Close())
		if err != nil {
			_ = os.Remove(output)
		}
	}()
	if _, err = io.CopyN(out, &contextReader{ctx, f}, int64(limit)); err != nil {
		return id, err
	}
	directory, err := codeDirectory(ctx, out, int64(limit), 0, machImage{}, o, components)
	if err != nil {
		return id, err
	}
	signature, err := signatureBlob(ctx, directory, components, o)
	if err != nil {
		return id, err
	}
	be.PutUint64(footer[304:], uint64(len(signature)))
	if _, err = out.Write(signature); err != nil {
		return id, err
	}
	if _, err = out.Write(footer[:]); err != nil {
		return id, err
	}
	if err = ctx.Err(); err != nil {
		return id, err
	}
	return identity(0, o.Identifier, directory), out.Chmod(st.Mode().Perm())
}
