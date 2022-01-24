package rule

import "testing"

func TestDetectG101(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		want    bool
		wantErr bool
	}{
		{
			name:    "Detect G101 rule success with public_key content",
			line:    "SECRET_KEY start with prefix public_key",
			want:    true,
			wantErr: false,
		},
		{
			name:    "Detect G101 rule success with private_key content",
			line:    "SECRET_KEY start with prefix private_key",
			want:    true,
			wantErr: false,
		},
		{
			name:    "Detect G101 rule success with public_key content",
			line:    "SECRET_KEY start with prefix key",
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DetectG101(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectG101() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DetectG101() got = %v, want %v", got, tt.want)
			}
		})
	}
}
