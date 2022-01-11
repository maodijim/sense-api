package sense

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"log"

	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	apiUrl               = "https://api.sense.com/apiservice/api/v1"
	formContentType      = "application/x-www-form-urlencoded"
	wssHost              = "clientrt.sense.com"
	senseProtocol        = "8"
	errMfaRequired       = "mfa_required"
	deviceId             = "dnccvg0pzp62lmdotpg2hho0q7333ihzccpr6ifijc5lh9kigzc7m6jcovy8avq4wb2c40u0tsygwzi3vcj9lkmn0v7wjar36r5tvdqjqcfbtcrqaawfemmvmtamahg2"
	defaultTimelineItems = 30
)

var (
	MaxMessageCache = 200
)

func (s *SenseApi) apiRequest(method, url, contentType, body string) (res *http.Response, err error) {
	if s.authRes.AccessToken != "" && isTokenExpired(s.authRes.AccessToken) {
		s.authRes.AccessToken = ""
		err = s.RenewToken()
		if err != nil {
			return res, errors.New("token expired")
		}
		res, err = s.apiRequest(method, url, contentType, body)
		return res, err
	}
	req, _ := http.NewRequest(method, url, strings.NewReader(body))
	headers := req.Header
	if contentType != "" {
		headers.Add("Content-Type", contentType)
	}
	headers.Add("x-sense-device-id", deviceId)
	headers.Add("authorization", "bearer "+s.authRes.AccessToken)
	client := http.Client{}
	res, err = client.Do(req)
	return res, err
}

func (s *SenseApi) mfaAuth(mfaToken, totp string) (err error) {
	u := apiUrl + "/authenticate/mfa"
	v := url.Values{}
	v.Add("mfaToken", mfaToken)
	v.Add("totp", totp)
	res, err := s.apiRequest(http.MethodPost, u, formContentType, v.Encode())
	if err != nil {
		return err
	}
	authRes := AuthRes{}
	err = parseRes(res, &authRes)
	if err != nil {
		return err
	}
	if authRes.Authorized {
		s.authSet(authRes)
	} else {
		return errors.New(authRes.ErrorReason)
	}
	return err
}

func (s *SenseApi) authSet(a AuthRes) {
	s.authRes = a
	s.wssEndpoint = "monitors/" + s.getMonitorId() + "/realtimefeed"
}

func (s *SenseApi) authenticate(username, password string) (err error) {
	authUrl := apiUrl + "/authenticate"
	v := url.Values{}
	v.Add("email", username)
	v.Add("password", password)
	res, err := s.apiRequest(http.MethodPost, authUrl, formContentType, v.Encode())
	if err != nil {
		return err
	}
	authRes := AuthRes{}
	err = parseRes(res, &authRes)
	if err != nil {
		return err
	}
	if authRes.Authorized {
		s.authSet(authRes)
	} else if authRes.Status == errMfaRequired {
		if authRes.MfaType != "totp" {
			return errors.New("only support totp mfa type but received unsupported mfa type: " + authRes.MfaType)
		}
		var totp string
		fmt.Println("Please enter two factor code: ")
		fmt.Scanln(&totp)
		err = s.mfaAuth(authRes.MfaToken, totp)
		if err != nil {
			return err
		}
	} else {
		return errors.New(authRes.ErrorReason)
	}
	return err
}

func (s *SenseApi) getMonitorId() string {
	if len(s.authRes.Monitors) == 0 {
		return ""
	}
	return strconv.FormatInt(int64(s.authRes.Monitors[0].Id), 10)
}

func (s *SenseApi) RenewToken() (err error) {
	u := fmt.Sprintf("%s/renew", apiUrl)
	v := url.Values{}
	v.Add("refresh_token", s.authRes.RefreshToken)
	v.Add("user_id", strconv.FormatInt(int64(s.authRes.UserId), 10))
	v.Add("is_access_token", "true")
	res, err := s.apiRequest(http.MethodPost, u, formContentType, v.Encode())
	if err != nil {
		return err
	}
	err = parseRes(res, &s.authRes)
	if err != nil {
		return err
	}
	return err
}

func (s *SenseApi) ListenWss() (err error) {
	q := url.Values{}
	q.Add("access_token", s.authRes.AccessToken)
	q.Add("sense_protocol", senseProtocol)
	q.Add("sense_client_type", "web")
	q.Add("sense_device_id", deviceId)
	u := url.URL{Scheme: "wss", Host: wssHost, Path: s.wssEndpoint, RawQuery: q.Encode()}
	s.ws, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	return err
}

func (s *SenseApi) AlwaysOn() (al *AlwaysOn, err error) {
	u := fmt.Sprintf("%s/app/monitors/%s/devices/always_on", apiUrl, s.getMonitorId())
	res, err := s.apiRequest("", u, formContentType, "")
	if err != nil {
		return al, err
	}
	al = &AlwaysOn{}
	err = parseRes(res, al)
	return al, err
}

func (s *SenseApi) DevicesOverview(includeMerged bool) (do *DevicesOverview, err error) {
	u := fmt.Sprintf("%s/app/monitors/%s/devices/overview?include_merged=%t", apiUrl, s.getMonitorId(), includeMerged)
	res, err := s.apiRequest("", u, "", "")
	if err != nil {
		return do, err
	}
	do = &DevicesOverview{}
	err = parseRes(res, do)
	if err != nil {
		return do, err
	}
	return do, err
}

func (s *SenseApi) TimeLine(items int) (tl *TimeLineRes, err error) {
	if items <= 0 {
		items = defaultTimelineItems
	}
	u := fmt.Sprintf("%s/users/%d/timeline?n_items=%d", apiUrl, s.authRes.UserId, items)
	res, err := s.apiRequest("", u, "", "")
	if err != nil {
		return tl, err
	}
	tl = &TimeLineRes{}
	err = parseRes(res, tl)
	if err != nil {
		return tl, err
	}
	return tl, err
}

type TrendScale string

const (
	TrendMonth TrendScale = "MONTH"
	TrendYear  TrendScale = "YEAR"
	TrendWeek  TrendScale = "WEEK"
	TrendDay   TrendScale = "DAY"
)

func (s *SenseApi) Trend(scale TrendScale, start time.Time) (trend *TrendType, err error) {
	v := url.Values{}
	v.Add("monitor_id", s.getMonitorId())
	v.Add("device_id", "")
	v.Add("scale", string(scale))
	v.Add("start", start.Format(time.RFC3339))
	u := fmt.Sprintf("%s/app/history/trends?%s", apiUrl, v.Encode())
	res, err := s.apiRequest("", u, "", "")
	if err != nil {
		return trend, err
	}
	trend = &TrendType{}
	err = parseRes(res, trend)
	if err != nil {
		return trend, err
	}
	return trend, err
}

func (s *SenseApi) reconnect() (err error) {
	if s.ws == nil {
		if s.authRes.AccessToken != "" && isTokenExpired(s.authRes.AccessToken) {
			s.authRes.AccessToken = ""
			err = s.RenewToken()
			if err != nil {
				return err
			}
			err = s.ListenWss()
			if err != nil {
				return err
			}
		} else {
			err = s.ListenWss()
			if err != nil {
				return err
			}
		}
	}
	return err
}

// ReadMessageAsync read message async and store messages in cache
// use ReadMessages() to retrieve cached messages
func (s *SenseApi) ReadMessageAsync(close <-chan bool) (err error) {
	s.readingAsync = true
	for {
		select {
		case <-close:
			s.readingAsync = false
			return
		default:
			err = s.reconnect()
			if err != nil {
				s.readingAsync = false
				return err
			}
			rt, err := s.ReadMessage()
			if err != nil {
				s.readingAsync = false
				return err
			}
			s.mutex.Lock()
			s.messages = append(s.messages, *rt)
			if len(s.messages) > MaxMessageCache {
				s.messages = s.messages[len(s.messages)-MaxMessageCache : len(s.messages)]
			}
			s.mutex.Unlock()
		}
	}
}

// ReadMessages Read Cached messages
// Cached messages are created async by ReadMessageAsync()
func (s *SenseApi) ReadMessages() (msgs []RealTime, err error) {
	if !s.readingAsync {
		return msgs, errors.New("reading async not start please run ReadMessageAsync() to start async reader")
	}
	err = s.reconnect()
	if err != nil {
		return msgs, err
	}
	s.mutex.Lock()
	msgs = s.messages
	s.messages = []RealTime{}
	s.mutex.Unlock()
	return msgs, err
}

// ReadMessage Read one real time message
func (s *SenseApi) ReadMessage() (msg *RealTime, err error) {
	err = s.reconnect()
	if err != nil {
		return msg, err
	}
	_, b, err := s.ws.ReadMessage()
	if err != nil {
		return msg, err
	}
	msg = &RealTime{}
	err = json.Unmarshal(b, &msg)
	if err != nil {
		return msg, err
	}
	return msg, err
}

// Close websocket connection
func (s *SenseApi) Close() (err error) {
	if s.ws == nil {
		return errors.New("websocket already closed")
	}
	err = s.ws.Close()
	if err != nil {
		return err
	}
	s.ws = nil
	return err
}

func unmarshallJWT(token string) (result *jwtClaims) {
	t, _, err := new(jwt.Parser).ParseUnverified(token, &jwtClaims{})
	if err != nil {
		log.Printf("failed to parse jwt token: %s", err)
	}
	if claims, ok := t.Claims.(*jwtClaims); ok {
		return claims
	}
	return result
}

func isTokenExpired(token string) (ok bool) {
	ok = true
	tPart := strings.Split(token, ".")
	if len(tPart) != 5 {
		return ok
	}
	j := unmarshallJWT(strings.Join(tPart[2:], "."))
	exp := j.Exp
	if time.Since(time.Unix(int64(exp), 0)) < 0 {
		return false
	}
	return ok
}

func parseRes(res *http.Response, parseType interface{}) (err error) {
	if reflect.TypeOf(parseType).Kind() != reflect.Ptr {
		return errors.New("parseType must be pointer reference")
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	err = json.Unmarshal(b, parseType)
	if err != nil {
		return err
	}
	return err
}

func NewSenseApi(username, password string) (s *SenseApi, err error) {
	s = &SenseApi{
		messages: []RealTime{},
	}
	err = s.authenticate(username, password)
	if err != nil {
		return s, err
	}
	err = s.ListenWss()
	return s, err
}
