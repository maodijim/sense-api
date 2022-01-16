package sense

import (
	"testing"
	"time"
)

func TestNewSenseApi(t *testing.T) {
	type cred struct {
		username string
		password string
	}
	tests := []struct {
		name string
		cred cred
	}{
		{
			name: "Test Panic",
			cred: cred{
				username: "",
				password: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewSenseApi(tt.cred.username, tt.cred.password)
			_, _ = s.ReadMessage()
			_ = s.RenewToken()
			_, _ = s.DevicesOverview(true)
			_, _ = s.TimeLine(30)
			_, _ = s.AlwaysOn()
			_, _ = s.Trend(TrendMonth, time.Now())
			_, _ = s.GetHistoryComparison()
		})
	}
}
