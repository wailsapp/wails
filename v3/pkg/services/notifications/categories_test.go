package notifications

import "testing"

func TestDecodeCategories(t *testing.T) {
	for _, tc := range []struct{ input, existingID string }{
		{"null", ""},
		{" \nnull\t", ""},
		{"{}", ""},
		{`{"existing":{"id":"existing"}}`, "existing"},
	} {
		t.Run(tc.input, func(t *testing.T) {
			categories, err := decodeCategories([]byte(tc.input))
			if err != nil {
				t.Fatal(err)
			}
			if categories == nil {
				t.Fatal("decoded map is nil")
			}
			if got := categories["existing"].ID; got != tc.existingID {
				t.Errorf("existing ID=%q, want %q", got, tc.existingID)
			}
			categories["new"] = NotificationCategory{ID: "new"}
		})
	}
	for _, input := range []string{"", "{", "[]", `{"invalid":123}`} {
		if _, err := decodeCategories([]byte(input)); err == nil {
			t.Errorf("accepted invalid category data: %q", input)
		}
	}
}
