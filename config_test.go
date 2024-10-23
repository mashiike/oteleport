package oteleport_test

import (
	"testing"

	"github.com/mashiike/oteleport"
	"github.com/stretchr/testify/assert"
)

func TestServerConfig_Load(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		opts    *oteleport.LoadOptions
		wantErr bool
	}{
		{
			name:    "valid config",
			path:    "testdata/default.jsonnet",
			opts:    nil,
			wantErr: false,
		},
		{
			name:    "with access key",
			path:    "testdata/with_access_key.jsonnet",
			opts:    nil,
			wantErr: false,
		},
		{
			name:    "invalid config path",
			path:    "testdata/invalid.jsonnet",
			opts:    nil,
			wantErr: true,
		},
		{
			name:    "unknwon_fields config",
			path:    "testdata/unknwon_fields.jsonnet",
			opts:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := oteleport.DefaultServerConfig()
			err := cfg.Load(tt.path, tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
