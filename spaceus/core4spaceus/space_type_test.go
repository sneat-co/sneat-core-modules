package core4spaceus

import "testing"

func TestIsValidSpaceType(t *testing.T) {
	type args struct {
		v SpaceType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"SpaceTypePrivate", args{SpaceTypePrivate}, true},
		{"EmptySpaceType", args{""}, false},
		{"InvalidSpaceType", args{"Foo"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidSpaceType(tt.args.v); got != tt.want {
				t.Errorf("IsValidSpaceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSpaceRef(t *testing.T) {
	type args struct {
		spaceType SpaceType
		spaceID   string
	}
	tests := []struct {
		name string
		args args
		want SpaceRef
	}{
		{"ShouldPass", args{SpaceTypePrivate, "foo"}, "private!foo"},
		{"EmptySpaceType", args{"", "foo"}, ""},
		{"ShouldPass", args{SpaceTypePrivate, ""}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want == "" {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("NewSpaceRef() did not panic")
					}
				}()
			}
			if got := NewSpaceRef(tt.args.spaceType, tt.args.spaceID); got != tt.want {
				t.Errorf("NewSpaceRef() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpaceRef_SpaceID(t *testing.T) {
	tests := []struct {
		name string
		v    SpaceRef
		want string
	}{
		{"full", "private!foo", "foo"},
		{"no_spaceType", "!foo", "foo"},
		{"no_spaceID", "private!", ""},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.SpaceID(); got != tt.want {
				t.Errorf("SpaceID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpaceRef_SpaceType(t *testing.T) {
	tests := []struct {
		name string
		v    SpaceRef
		want SpaceType
	}{
		{"full", "private!foo", "private"},
		{"no_spaceID", "private!", "private"},
		{"no_spaceType", "!foo", ""},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.SpaceType(); got != tt.want {
				t.Errorf("SpaceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpaceRef_UrlPath(t *testing.T) {
	tests := []struct {
		name string
		v    SpaceRef
		want string
	}{
		{"ShouldPass", "private!foo", "private/foo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.UrlPath(); got != tt.want {
				t.Errorf("UrlPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWeakSpaceRef(t *testing.T) {
	type args struct {
		spaceType SpaceType
	}
	tests := []struct {
		name string
		args args
		want SpaceRef
	}{
		{
			name: "ShouldPass",
			args: args{SpaceTypePrivate},
			want: "private",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWeakSpaceRef(tt.args.spaceType); got != tt.want {
				t.Errorf("NewWeakSpaceRef() = %v, want %v", got, tt.want)
			}
		})
	}
}
