package app

import (
	"testing"
)

func TestEstimateTokens(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		method  string
		want    int
		wantErr bool
	}{
		{
			name:    "Empty text with default method",
			text:    "",
			method:  "",
			want:    0,
			wantErr: false,
		},
		{
			name:    "Simple text with default method (max)",
			text:    "Hello world, how are you?",
			method:  "",
			want:    7, // max(5/0.75, 26/4) = max(6.67, 6.5) = 7
			wantErr: false,
		},
		{
			name:    "Simple text with average method",
			text:    "Hello world, how are you?",
			method:  "average",
			want:    7, // (5/0.75 + 26/4)/2 = (6.67 + 6.5)/2 = 6.58 ≈ 7
			wantErr: false,
		},
		{
			name:    "Simple text with words method",
			text:    "Hello world, how are you?",
			method:  "words",
			want:    7, // 5/0.75 = 6.67 ≈ 7
			wantErr: false,
		},
		{
			name:    "Simple text with chars method",
			text:    "Hello world, how are you?",
			method:  "chars",
			want:    7, // 26/4 = 6.5, ceiling = 7
			wantErr: false,
		},
		{
			name:    "Simple text with min method",
			text:    "Hello world, how are you?",
			method:  "min",
			want:    7, // min(5/0.75, 26/4) = min(6.67, 6.5) = 6.5, ceiling = 7
			wantErr: false,
		},
		{
			name:    "Longer text with max method",
			text:    "This is a longer text that should have more tokens. It contains multiple sentences and words of varying lengths.",
			method:  "max",
			want:    28, // max(19/0.75, 109/4) = max(25.33, 27.25) = 28
			wantErr: false,
		},
		{
			name:    "Text with special characters",
			text:    "Special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?",
			method:  "chars",
			want:    11, // 42/4 = 10.5, ceiling = 11
			wantErr: false,
		},
		{
			name:    "Text with multiple spaces",
			text:    "This    has    multiple    spaces",
			method:  "words",
			want:    6, // 4/0.75 = 5.33, ceiling = 6
			wantErr: false,
		},
		{
			name:    "Invalid method",
			text:    "Some text here",
			method:  "invalid",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EstimateTokens(tt.text, tt.method)
			if (err != nil) != tt.wantErr {
				t.Errorf("EstimateTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EstimateTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}
