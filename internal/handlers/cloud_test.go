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
