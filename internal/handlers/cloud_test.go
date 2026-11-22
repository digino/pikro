package handlers

import (
	"reflect"
	"testing"

	"github.com/digino/pikro/internal/config"
)

func TestCloudPlanOptions(t *testing.T) {
	metas := map[string]config.ProfileMeta{
		"1-heure": {Price: "100"},
		"1-jour":  {Price: "300"},
		"promo":   {Price: "1.5"}, // Cloud prices are whole units
		"free":    {Price: ""},
	}
	got := cloudPlanOptions([]string{"default", "1-heure", "1-jour", "promo", "free"}, metas)
	want := []cloudPlanOption{
		{Profile: "default"},
		{Profile: "1-heure", Price: 100},
		{Profile: "1-jour", Price: 300},
		{Profile: "promo"},
		{Profile: "free"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("cloudPlanOptions =\n %+v\nwant\n %+v", got, want)
	}
}

func TestWithoutCloudVouchers(t *testing.T) {
	users := []map[string]any{
		{"name": "abc12", "comment": "pikro-cloud"},                         // unsold pool voucher
		{"name": "fhp67", "comment": "exp:2026-11-06 10:00:00"},             // a normal voucher, in use
		{"name": "def34", "comment": "exp:2026-11-06 10:00:00 pikro-cloud"}, // a pool voucher sold and used: shown
		{"name": "admin"},
	}
	shown, hidden := withoutCloudVouchers(users)
	var names []string
	for _, u := range shown {
		names = append(names, u["name"].(string))
	}
	if hidden != 1 || !reflect.DeepEqual(names, []string{"fhp67", "def34", "admin"}) {
		t.Errorf("shown %v, hidden %d; want fhp67 def34 admin, 1 hidden", names, hidden)
	}
}
