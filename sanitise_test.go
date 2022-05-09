package cas

import (
	"net/url"
	"reflect"
	"testing"
)

func Test_sanitisedURL(t *testing.T) {
	type args struct {
		unclean *url.URL
	}
	tests := []struct {
		name    string
		args    args
		want    *url.URL
		wantErr bool
	}{
		{
			name: "Test the URL Scheme chaos value, cause be dealt with requestURL method",
			args: args{
				unclean: &url.URL{
					Scheme:      "chaos_input_from_header_X-Forwarded-Proto",
					Opaque:      "",
					User:        &url.Userinfo{},
					Host:        "a.b.c",
					Path:        "/",
					RawPath:     "/",
					ForceQuery:  false,
					RawQuery:    "",
					Fragment:    "",
					RawFragment: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sanitisedURL(tt.args.unclean)
			if (err != nil) != tt.wantErr {
				t.Errorf("sanitisedURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sanitisedURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
