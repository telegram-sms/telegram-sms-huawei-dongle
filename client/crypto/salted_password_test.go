package crypto

import "testing"

func TestEncodeSaltedPassword(t *testing.T) {
	type args struct {
		username string
		password string
		token    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "quick sanity test",
			args: args{
				username: "admin",
				password: "admin",
				token:    "g_requestVerificationToken",
			},
			want: "MmI4NDg5YjBhNjI4ZDE5M2QwZTE2ZDYzNDljOGVmZGUwMDg1YjlmZjdlYzM4NTA1ZjBjNmQyNjI4ZDViOGIzMg==",
		},
		{
			name: "longer password test",
			args: args{
				username: "admin",
				password: "this is a very very long password",
				token:    "g_requestVerificationToken",
			},
			want: "ZTAyZGQ2N2E2ZGE3OTQyNDRhZDA4MDcwNDc2NWVhNGFlNWIyZjQ2MWI5MTFjODY3M2ZiNWMzZDUxM2RiZTNlZA==",
		},
		{
			name: "admin 2",
			args: args{
				username: "admin",
				password: "admin",
				token:    "XJxa4n/7SgTVnx9GeVNl6zxbTH4FfzQ9",
			},
			want: "NDQ4NDJkMmI2YjYyZmEyYzJlYTBhYTVjMGRmYzJkZjRhNGRhZDQzZGU3ZGNiMGU0OGJkMjMwZDMxZmM0ZjUzOA==",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeSaltedPassword(tt.args.username, tt.args.password, tt.args.token); got != tt.want {
				t.Errorf("EncodeSaltedPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
