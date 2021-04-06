package mac

import "testing"

func TestShowNotification(t *testing.T) {
	type args struct {
		title    string
		subtitle string
		message  string
		sound    string
	}
	tests := []struct {
		name     string
		title    string
		subtitle string
		message  string
		sound    string
		wantErr  bool
	}{
		{"No message", "", "", "", "", false},
		{"Title only", "I am a Title", "", "", "", false},
		{"SubTitle only", "", "I am a subtitle", "", "", false},
		{"Message only", "", "", "I am a message!", "", false},
		{"Sound only", "", "", "", "submarine.aiff", false},
		{"Full", "Title", "Subtitle", "This is a long message to show that text gets wrapped in a notification", "submarine.aiff", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ShowNotification(tt.title, tt.subtitle, tt.message, tt.sound); (err != nil) != tt.wantErr {
				t.Errorf("ShowNotification() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
