package gomodprivate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExtractTag(t *testing.T) {
	type args struct {
		packageName string
	}
	tests := []struct {
		name        string
		args        args
		wantPackage string
		wantTag     string
		wantErr     bool
	}{
		{
			"Valid",
			args{
				packageName: "rendoru.com/module/sync-mq@v1.0.0",
			},
			"rendoru.com/module/sync-mq",
			"v1.0.0",
			false,
		},
		{
			"Valid_NoTag",
			args{
				packageName: "rendoru.com/module/sync-mq",
			},
			"rendoru.com/module/sync-mq",
			"",
			false,
		},
		{
			"Invalid_Empty",
			args{
				packageName: "",
			},
			"",
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := _ExtractTag(tt.args.packageName)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, got, tt.wantPackage)
			assert.Equal(t, got1, tt.wantTag)
		})
	}
}
