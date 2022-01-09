package sense

import (
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type AuthRes struct {
	MfaToken    string `json:"mfa_token"`
	MfaType     string `json:"mfa_type"`
	Status      string `json:"status"`
	ErrorReason string `json:"error_reason"`
	Authorized  bool   `json:"authorized"`
	AccountId   int    `json:"account_id"`
	UserId      int    `json:"user_id"`
	AccessToken string `json:"access_token"`
	Settings    struct {
		UserId   int `json:"user_id"`
		Settings struct {
			Notifications struct {
				Field1 struct {
					NewNamedDevicePush   bool `json:"new_named_device_push"`
					NewNamedDeviceEmail  bool `json:"new_named_device_email"`
					MonitorOfflinePush   bool `json:"monitor_offline_push"`
					MonitorOfflineEmail  bool `json:"monitor_offline_email"`
					MonitorMonthlyEmail  bool `json:"monitor_monthly_email"`
					AlwaysOnChangePush   bool `json:"always_on_change_push"`
					ComparisonChangePush bool `json:"comparison_change_push"`
					NewPeakPush          bool `json:"new_peak_push"`
					NewPeakEmail         bool `json:"new_peak_email"`
					MonthlyChangePush    bool `json:"monthly_change_push"`
					WeeklyChangePush     bool `json:"weekly_change_push"`
					DailyChangePush      bool `json:"daily_change_push"`
					GeneratorOnPush      bool `json:"generator_on_push"`
					GeneratorOffPush     bool `json:"generator_off_push"`
					TimeOfUse            bool `json:"time_of_use"`
				} `json:"187899"`
				Field2 struct {
					NewNamedDevicePush   bool `json:"new_named_device_push"`
					NewNamedDeviceEmail  bool `json:"new_named_device_email"`
					MonitorOfflinePush   bool `json:"monitor_offline_push"`
					MonitorOfflineEmail  bool `json:"monitor_offline_email"`
					MonitorMonthlyEmail  bool `json:"monitor_monthly_email"`
					AlwaysOnChangePush   bool `json:"always_on_change_push"`
					ComparisonChangePush bool `json:"comparison_change_push"`
					NewPeakPush          bool `json:"new_peak_push"`
					NewPeakEmail         bool `json:"new_peak_email"`
					MonthlyChangePush    bool `json:"monthly_change_push"`
					WeeklyChangePush     bool `json:"weekly_change_push"`
					DailyChangePush      bool `json:"daily_change_push"`
					GeneratorOnPush      bool `json:"generator_on_push"`
					GeneratorOffPush     bool `json:"generator_off_push"`
					TimeOfUse            bool `json:"time_of_use"`
				} `json:"188074"`
			} `json:"notifications"`
			LabsEnabled bool `json:"labs_enabled"`
		} `json:"settings"`
		Version int `json:"version"`
	} `json:"settings"`
	Monitors []struct {
		Id              int    `json:"id"`
		SerialNumber    string `json:"serial_number"`
		TimeZone        string `json:"time_zone"`
		SolarConnected  bool   `json:"solar_connected"`
		SolarConfigured bool   `json:"solar_configured"`
		Online          bool   `json:"online"`
		Attributes      struct {
			Id                  int         `json:"id"`
			Name                string      `json:"name"`
			State               string      `json:"state"`
			Cost                float64     `json:"cost"`
			SellBackRate        float64     `json:"sell_back_rate"`
			UserSetCost         bool        `json:"user_set_cost"`
			CycleStart          int         `json:"cycle_start"`
			BasementType        string      `json:"basement_type"`
			HomeSizeType        string      `json:"home_size_type"`
			HomeType            string      `json:"home_type"`
			NumberOfOccupants   string      `json:"number_of_occupants"`
			OccupancyType       string      `json:"occupancy_type"`
			YearBuiltType       string      `json:"year_built_type"`
			PostalCode          string      `json:"postal_code"`
			ElectricityCost     interface{} `json:"electricity_cost"`
			ShowCost            bool        `json:"show_cost"`
			TouEnabled          bool        `json:"tou_enabled"`
			SolarTouEnabled     bool        `json:"solar_tou_enabled"`
			PowerRegion         interface{} `json:"power_region"`
			UserSetSellBackRate bool        `json:"user_set_sell_back_rate"`
		} `json:"attributes"`
		SignalCheckCompletedTime time.Time     `json:"signal_check_completed_time"`
		DataSharing              []interface{} `json:"data_sharing"`
		EthernetSupported        bool          `json:"ethernet_supported"`
		AuxIgnore                bool          `json:"aux_ignore"`
		AuxPort                  string        `json:"aux_port"`
		HardwareType             string        `json:"hardware_type"`
	} `json:"monitors"`
	BridgeServer string    `json:"bridge_server"`
	DateCreated  time.Time `json:"date_created"`
	TotpEnabled  bool      `json:"totp_enabled"`
	AbCohort     string    `json:"ab_cohort"`
	RefreshToken string    `json:"refresh_token"`
}

type SenseApi struct {
	ws           *websocket.Conn
	wssEndpoint  string
	refreshToken string
	authRes      AuthRes
	mutex        sync.RWMutex
	messages     []RealTime
	readingAsync bool
}

type AlwaysOn struct {
	Alerts struct {
		Allowed bool `json:"allowed"`
		Enabled bool `json:"enabled"`
	} `json:"alerts"`
	Device struct {
		Id        string `json:"id"`
		MonitorId int    `json:"monitorId"`
		Name      string `json:"name"`
		Icon      string `json:"icon"`
		Tags      struct {
			DefaultUserDeviceType       string `json:"DefaultUserDeviceType"`
			DeviceListAllowed           string `json:"DeviceListAllowed"`
			SSIEnabled                  string `json:"SSIEnabled"`
			TimelineAllowed             string `json:"TimelineAllowed"`
			UserEditable                string `json:"UserEditable"`
			UserDeviceTypeDisplayString string `json:"UserDeviceTypeDisplayString"`
			UserDeviceType              string `json:"UserDeviceType"`
		} `json:"tags"`
	} `json:"device"`
	LastState     interface{} `json:"last_state"`
	Notes         interface{} `json:"notes"`
	Info          string      `json:"info"`
	LastStateTime interface{} `json:"last_state_time"`
	Usage         struct {
		AvgMonthlyKWH    float64 `json:"avg_monthly_KWH"`
		AvgMonthlyPct    float64 `json:"avg_monthly_pct"`
		AvgWatts         float64 `json:"avg_watts"`
		YearlyKWH        float64 `json:"yearly_KWH"`
		YearlyText       string  `json:"yearly_text"`
		YearlyCost       int     `json:"yearly_cost"`
		AvgMonthlyCost   int     `json:"avg_monthly_cost"`
		CurrentAoWattage int     `json:"current_ao_wattage"`
		Comparison       struct {
			ComparisonText string   `json:"comparison_text"`
			TercileStrings []string `json:"tercile_strings"`
			CohortMarker   int      `json:"cohort_marker"`
			Cohort         struct {
				Id         int    `json:"id"`
				PostalCode string `json:"postal_code"`
				AreaCode   string `json:"area_code"`
				State      string `json:"state"`
				HomeSize   string `json:"home_size"`
				Location   string `json:"location"`
			} `json:"cohort"`
			Title      string  `json:"title"`
			W          float64 `json:"w"`
			CohortAvgW float64 `json:"cohort_avg_w"`
		} `json:"comparison"`
	} `json:"usage"`
	AlwaysOn struct {
		Description  string  `json:"description"`
		TotalWatts   float64 `json:"total_watts"`
		UnknownWatts float64 `json:"unknown_watts"`
		Devices      []struct {
			Id string  `json:"id"`
			W  float64 `json:"w"`
		} `json:"devices"`
	} `json:"always_on"`
	Timeline struct {
		Visible bool `json:"visible"`
		Allowed bool `json:"allowed"`
	} `json:"timeline"`
	Blurb struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	} `json:"blurb"`
}

type TimeLineRes struct {
	MoreItems   bool          `json:"more_items"`
	StickyItems []interface{} `json:"sticky_items"`
	UserId      int           `json:"user_id"`
	Items       []struct {
		Time           time.Time `json:"time"`
		Type           string    `json:"type"`
		Icon           string    `json:"icon"`
		Body           string    `json:"body"`
		Destination    string    `json:"destination"`
		DeviceId       string    `json:"device_id"`
		DeviceState    string    `json:"device_state"`
		ShowAction     bool      `json:"show_action"`
		AllowSticky    bool      `json:"allow_sticky"`
		UserDeviceType string    `json:"user_device_type"`
		StartTime      time.Time `json:"start_time,omitempty"`
		Children       []struct {
			Time        time.Time `json:"time"`
			Type        string    `json:"type"`
			Body        string    `json:"body"`
			DeviceId    string    `json:"device_id"`
			ShowAction  bool      `json:"show_action"`
			AllowSticky bool      `json:"allow_sticky"`
			StartTime   time.Time `json:"start_time,omitempty"`
			Children    []struct {
				Time        time.Time `json:"time"`
				Type        string    `json:"type"`
				DeviceId    string    `json:"device_id"`
				ShowAction  bool      `json:"show_action"`
				AllowSticky bool      `json:"allow_sticky"`
			} `json:"children,omitempty"`
		} `json:"children,omitempty"`
		Count int `json:"count,omitempty"`
	} `json:"items"`
}

type DevicesOverview struct {
	Devices []struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Make  string `json:"make,omitempty"`
		Model string `json:"model,omitempty"`
		Icon  string `json:"icon"`
		Tags  struct {
			Alertable             string    `json:"Alertable,omitempty"`
			AlwaysOn              string    `json:"AlwaysOn,omitempty"`
			DateCreated           time.Time `json:"DateCreated,omitempty"`
			DateFirstUsage        string    `json:"DateFirstUsage,omitempty"`
			DefaultUserDeviceType string    `json:"DefaultUserDeviceType,omitempty"`
			DeployToMonitor       string    `json:"DeployToMonitor,omitempty"`
			DeviceListAllowed     string    `json:"DeviceListAllowed"`
			ModelCreatedVersion   string    `json:"ModelCreatedVersion,omitempty"`
			ModelUpdatedVersion   string    `json:"ModelUpdatedVersion,omitempty"`
			NameUseredit          string    `json:"name_useredit,omitempty"`
			OriginalName          string    `json:"OriginalName,omitempty"`
			PeerNames             []struct {
				Name                        string  `json:"Name"`
				UserDeviceType              string  `json:"UserDeviceType"`
				Percent                     float64 `json:"Percent"`
				Icon                        string  `json:"Icon"`
				UserDeviceTypeDisplayString string  `json:"UserDeviceTypeDisplayString"`
				Make                        string  `json:"Make,omitempty"`
				Model                       string  `json:"Model,omitempty"`
			} `json:"PeerNames,omitempty"`
			Pending                     string `json:"Pending,omitempty"`
			Revoked                     string `json:"Revoked,omitempty"`
			SSIEnabled                  string `json:"SSIEnabled,omitempty"`
			TimelineAllowed             string `json:"TimelineAllowed"`
			TimelineDefault             string `json:"TimelineDefault,omitempty"`
			Type                        string `json:"Type,omitempty"`
			UserDeletable               string `json:"UserDeletable,omitempty"`
			UserDeviceType              string `json:"UserDeviceType"`
			UserDeviceTypeDisplayString string `json:"UserDeviceTypeDisplayString"`
			UserEditable                string `json:"UserEditable"`
			UserEditableMeta            string `json:"UserEditableMeta,omitempty"`
			UserMergeable               string `json:"UserMergeable,omitempty"`
			ExpectedAOWattage           int    `json:"ExpectedAOWattage,omitempty"`
			UserAdded                   string `json:"UserAdded,omitempty"`
			MergeId                     string `json:"MergeId,omitempty"`
			PreselectionIndex           int    `json:"PreselectionIndex,omitempty"`
			NameUserGuess               string `json:"NameUserGuess,omitempty"`
			MergedDevices               string `json:"MergedDevices,omitempty"`
			Virtual                     string `json:"Virtual,omitempty"`
			DefaultMake                 string `json:"DefaultMake,omitempty"`
			DefaultModel                string `json:"DefaultModel,omitempty"`
			UserDeleted                 string `json:"UserDeleted,omitempty"`
		} `json:"tags"`
		GivenMake     string `json:"given_make,omitempty"`
		GivenModel    string `json:"given_model,omitempty"`
		Location      string `json:"location,omitempty"`
		GivenLocation string `json:"given_location,omitempty"`
	} `json:"devices"`
	DeviceDataChecksum string `json:"device_data_checksum"`
}

type jwtClaims struct {
	Iss       string `json:"iss"`
	Exp       int    `json:"exp"`
	UserId    int    `json:"userId"`
	AccountId int    `json:"accountId"`
	Roles     string `json:"roles"`
	Dhash     string `json:"dhash"`
	jwt.StandardClaims
}

type TrendType struct {
	Steps       int       `json:"steps"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Consumption struct {
		Total   float64   `json:"total"`
		Totals  []float64 `json:"totals"`
		Devices []struct {
			Id        string `json:"id"`
			MonitorId int    `json:"monitorId"`
			Name      string `json:"name"`
			Icon      string `json:"icon"`
			Tags      struct {
				UserDeviceTypeDisplayString string `json:"UserDeviceTypeDisplayString"`
			} `json:"tags"`
			History     []float64 `json:"history"`
			Avgw        float64   `json:"avgw"`
			TotalKwh    float64   `json:"total_kwh"`
			TotalCost   int       `json:"total_cost"`
			CostHistory []int     `json:"cost_history"`
		} `json:"devices"`
		TotalCost  int   `json:"total_cost"`
		TotalCosts []int `json:"total_costs"`
	} `json:"consumption"`
	Production struct {
		Total      float64       `json:"total"`
		Totals     []float64     `json:"totals"`
		Devices    []interface{} `json:"devices"`
		TotalCost  int           `json:"total_cost"`
		TotalCosts []int         `json:"total_costs"`
	} `json:"production"`
	ToGrid                   float64     `json:"to_grid"`
	FromGrid                 float64     `json:"from_grid"`
	ConsumptionPercentChange interface{} `json:"consumption_percent_change"`
	ProductionPercentChange  interface{} `json:"production_percent_change"`
	ToGridCost               int         `json:"to_grid_cost"`
	FromGridCost             int         `json:"from_grid_cost"`
	TrendText                interface{} `json:"trend_text"`
	UsageText                interface{} `json:"usage_text"`
	Scale                    string      `json:"scale"`
	SolarPowered             int         `json:"solar_powered"`
	NetProduction            float64     `json:"net_production"`
	ProductionPct            int         `json:"production_pct"`
}

type PayloadType string

func (p PayloadType) String() string {
	return string(p)
}

const (
	PayloadRealTimeUpdate PayloadType = "realtime_update"
	PayloadDataChange     PayloadType = "data_change"
	PayloadMonitorInfo    PayloadType = "monitor_info"
	PayloadHello          PayloadType = "hello"
	PayloadErr            PayloadType = "error"
)

type RealTime struct {
	Payload struct {
		Online  bool      `json:"online"`
		Voltage []float64 `json:"voltage"`
		Frame   int       `json:"frame"`
		Devices []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
			Icon string `json:"icon"`
			Tags struct {
				DefaultUserDeviceType       string `json:"DefaultUserDeviceType"`
				DeviceListAllowed           string `json:"DeviceListAllowed"`
				SSIEnabled                  string `json:"SSIEnabled,omitempty"`
				TimelineAllowed             string `json:"TimelineAllowed"`
				UserDeviceType              string `json:"UserDeviceType"`
				UserDeviceTypeDisplayString string `json:"UserDeviceTypeDisplayString"`
				UserEditable                string `json:"UserEditable"`
				UserDeleted                 string `json:"UserDeleted,omitempty"`
				UserMergeable               string `json:"UserMergeable,omitempty"`
			} `json:"tags"`
			Attrs []interface{} `json:"attrs"`
			W     float64       `json:"w"`
		} `json:"devices"`
		Deltas      []interface{} `json:"deltas"`
		DefaultCost int           `json:"defaultCost"`
		Channels    []float64     `json:"channels"`
		Hz          float64       `json:"hz"`
		W           float64       `json:"w"`
		C           int           `json:"c"`
		TouAlert    struct {
			EndTime        time.Time `json:"end_time"`
			CostMultiplier float64   `json:"cost_multiplier"`
			TouCost        int       `json:"tou_cost"`
			Name           string    `json:"name"`
		} `json:"tou_alert"`
		SolarW float64 `json:"solar_w"`
		SolarC int     `json:"solar_c"`
		Stats  struct {
			Brcv float64 `json:"brcv"`
			Mrcv float64 `json:"mrcv"`
			Msnd float64 `json:"msnd"`
		} `json:"_stats"`
		Aux struct {
			Solar []float64 `json:"solar"`
		} `json:"aux"`
		DW       int `json:"d_w"`
		DSolarW  int `json:"d_solar_w"`
		GridW    int `json:"grid_w"`
		SolarPct int `json:"solar_pct"`
		Epoch    int `json:"epoch"`

		Features                string `json:"features"`
		UserVersion             int    `json:"user_version"`
		PartnerChecksum         string `json:"partner_checksum"`
		MonitorOverviewChecksum string `json:"monitor_overview_checksum"`
		DeviceDataChecksum      string `json:"device_data_checksum"`
		SettingsVersion         int    `json:"settings_version"`
		PendingEvents           struct {
			Type           string `json:"type"`
			NewDeviceFound struct {
				DeviceId  interface{} `json:"device_id"`
				Guid      string      `json:"guid"`
				Timestamp interface{} `json:"timestamp"`
			} `json:"new_device_found"`
			MonitorId int `json:"monitor_id"`
			Goal      struct {
				Guid           string      `json:"guid"`
				Timestamp      interface{} `json:"timestamp"`
				NotificationId interface{} `json:"notification_id"`
			} `json:"goal"`
		} `json:"pending_events"`
	} `json:"payload"`
	Type PayloadType `json:"type"`
}

func (r RealTime) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}
