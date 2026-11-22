package config

import (
	"reflect"
	"testing"
)

func TestSetCloudLink(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	if err := Save(&Config{Routers: []RouterProfile{{ID: "r1", Name: "OKLM"}, {ID: "r2"}}}); err != nil {
		t.Fatal(err)
	}
	link := &CloudLink{OwnerPhone: "+22890000000", Profiles: []string{"4H", "24H"}}
	for _, tt := range []struct {
		name, id string
		link     *CloudLink
		wantErr  bool
	}{
		{"saved", "r1", link, false},
		{"unknown router", "nope", link, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetCloudLink(tt.id, tt.link); (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, want error %v", err, tt.wantErr)
			}
		})
	}
	cfg, _ := Load()
	if got := cfg.Routers[0]; !reflect.DeepEqual(got.Cloud, link) || got.Name != "OKLM" {
		t.Errorf("router = %+v, want its link saved and the rest kept", got)
	}
	if cfg.Routers[1].Cloud != nil {
		t.Error("the other router got a link")
	}

	// Disconnecting forgets it.
	if err := SetCloudLink("r1", nil); err != nil {
		t.Fatal(err)
	}
	if cfg, _ := Load(); cfg.Routers[0].Cloud != nil {
		t.Error("link kept after disconnect")
	}
}
