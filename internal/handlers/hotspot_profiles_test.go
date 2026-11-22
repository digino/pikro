package handlers

import "testing"

func TestProfileBodyValidate(t *testing.T) {
	valid := profileBody{Name: "1-jour", SharedUsers: float64(1), RateLimit: "2M/4M", Validity: "1d", Price: "300"}
	tests := []struct {
		name   string
		change func(*profileBody)
		ok     bool
	}{
		{"valid", func(*profileBody) {}, true},
		{"all optional fields empty", func(b *profileBody) { *b = profileBody{Name: "free"} }, true},
		{"price with cents", func(b *profileBody) { b.Price = "1.50" }, true},
		{"shared users as text", func(b *profileBody) { b.SharedUsers = "3" }, true},
		{"shared users unlimited", func(b *profileBody) { b.SharedUsers = "unlimited" }, true},
		{"rate limit with burst", func(b *profileBody) { b.RateLimit = "1M/2M 2M/4M 512k/1M 10/10" }, true},
		{"validity 2h30m", func(b *profileBody) { b.Validity = "2h30m" }, true},

		{"no name", func(b *profileBody) { b.Name = " " }, false},
		{"name with ; (breaks Cloud sync)", func(b *profileBody) { b.Name = "1;jour" }, false},
		{"name with ,", func(b *profileBody) { b.Name = "1,jour" }, false},
		{"name with a quote", func(b *profileBody) { b.Name = `1"jour` }, false},
		{"price as text", func(b *profileBody) { b.Price = "500 F" }, false},
		{"price zero", func(b *profileBody) { b.Price = "0" }, false},
		{"price negative", func(b *profileBody) { b.Price = "-5" }, false},
		{"price with 3 decimals", func(b *profileBody) { b.Price = "1.505" }, false},
		{"validity as words", func(b *profileBody) { b.Validity = "one day" }, false},
		{"validity zero", func(b *profileBody) { b.Validity = "0d" }, false},
		{"shared users zero", func(b *profileBody) { b.SharedUsers = float64(0) }, false},
		{"shared users fraction", func(b *profileBody) { b.SharedUsers = float64(1.5) }, false},
		{"shared users text", func(b *profileBody) { b.SharedUsers = "many" }, false},
		{"rate limit as words", func(b *profileBody) { b.RateLimit = "fast" }, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := valid
			tt.change(&b)
			if err := b.validate(); (err == nil) != tt.ok {
				t.Errorf("validate() = %v, want ok=%v", err, tt.ok)
			}
		})
	}
}
