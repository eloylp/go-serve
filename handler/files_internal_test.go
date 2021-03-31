package handler

import (
	"testing"
)

func Test_checkPath(t *testing.T) {
	type args struct {
		docRoot string
		path    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Should pass when both paths are absolutes and path is inside the document root",
			args{"/var/www", "/var/www/uploaded"},
			false,
		},
		{
			"Should pass when both paths are relatives and path is inside the document root",
			args{"images", "images/uploaded"},
			false,
		},
		{
			"Should pass when both paths are relatives and path is inside the document root (dir traversing)",
			args{"images", "images/uploaded/.."},
			false,
		},
		{
			"Should pass when both paths are relatives and path is inside the document root (dir traversing)",
			args{"images", "images/uploaded/.."},
			false,
		},

		{
			"Should NOT pass when both paths are absolutes and path is outside the document root",
			args{"/var/www", "/var/non-allowed-dir/uploaded"},
			true,
		},
		{
			"Should NOT pass when both paths are relatives and path is outside the document root",
			args{"images", "photos/upload"},
			true,
		},
		{
			"Should NOT pass when both paths are relatives and path is outside the document root (dir traversing)",
			args{"images", "images/upload/../.."},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkPath(tt.args.docRoot, tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("checkPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
