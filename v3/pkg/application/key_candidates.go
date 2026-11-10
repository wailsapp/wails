package application

import "slices"

// punctuationAcceleratorCandidates lists the accelerator strings a press of a
// punctuation key can mean, most specific first, for platforms that report
// keys by virtual key (Windows):
//
//   - the key's own character with the modifiers as pressed: Ctrl+Shift+=
//     when Shift is held, Ctrl+= when not;
//   - when Shift is held and changes the character, the shifted character
//     without Shift: Ctrl++, which is what CmdOrCtrl+plus means (the "+"
//     character is Shift+= on a US keyboard);
//   - legacy, the name the key was reported by before (OEM_PLUS), so a
//     binding written that way keeps working.
//
// unshifted and shifted are 0 when the layout gives the key no character.
func punctuationAcceleratorCandidates(mods []modifier, unshifted, shifted rune, legacy string) []string {
	var out []string
	add := func(a accelerator) {
		if s := a.String(); !slices.Contains(out, s) {
			out = append(out, s)
		}
	}
	shift := slices.Contains(mods, ShiftKey)
	if unshifted != 0 {
		add(accelerator{Key: string(unshifted), Modifiers: mods})
	}
	if shift && shifted != 0 && shifted != unshifted {
		var rest []modifier
		for _, m := range mods {
			if m != ShiftKey {
				rest = append(rest, m)
			}
		}
		add(accelerator{Key: string(shifted), Modifiers: rest})
	}
	if legacy != "" {
		add(accelerator{Key: legacy, Modifiers: mods})
	}
	return out
}
