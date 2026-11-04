package service

import "testing"

func TestValidateSiteURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "Valid HTTPS",
			url:     "https://google.com",
			wantErr: false,
		},
		{
			name:    "Valid HTTP",
			url:     "http://example.com",
			wantErr: false,
		},
		{
			name:    "Valid HTTPS with path",
			url:     "https://example.com/api/health",
			wantErr: false,
		},
		{
			name:    "Invalid protocol",
			url:     "ftp://example.com",
			wantErr: true,
		},
		{
			name:    "Missing hostname",
			url:     "https://",
			wantErr: true,
		},
		{
			name:    "Empty URL",
			url:     "",
			wantErr: true,
		},
		{
			name:    "Forbidden port",
			url:     "https://example.com:5432",
			wantErr: true,
		},
		{
			name:    "Allowed port",
			url:     "https://example.com:443",
			wantErr: false,
		},
		{
			name:    "Credentials in URL",
			url:     "https://user:pass@example.com",
			wantErr: true,
		},
		{
			name:    "Fragment in URL",
			url:     "https://example.com/#section",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSiteURL(tt.url)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"ValidateSiteURL(%q): err=%v, wantErr=%v",
					tt.url,
					err,
					tt.wantErr,
				)
			}
		})
	}
}
