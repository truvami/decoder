package loracloud

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/solver"
	"go.uber.org/zap"
)

type LoracloudClient struct {
	accessToken string
	logger      *zap.Logger
	BaseUrl     string
	timeNow     func() time.Time
}

const (
	SemtechLoRaCloudBaseUrl  = "https://mgs.loracloud.com"
	TraxmateLoRaCloudBaseUrl = "https://lw.traxmate.io"
)

var _ solver.SolverV1 = &LoracloudClient{}

func (m LoracloudClient) isSemtechLoRaCloudShutdown() error {
	if m.BaseUrl != SemtechLoRaCloudBaseUrl {
		return nil
	}

	if m.timeNow().After(time.Date(2025, 7, 31, 0, 0, 0, 0, time.UTC)) {
		return ErrSemtechLoRaCloudShutdown
	}
	m.logger.Warn("LoRa Cloud is Sunsetting on 31.07.2025", zap.String("url", "https://www.semtech.com/loracloud-shutdown"))

	return nil
}

func NewLoracloudClient(ctx context.Context, accessToken string, logger *zap.Logger, options ...LoracloudClientOptions) (LoracloudClient, error) {
	client := LoracloudClient{
		accessToken: accessToken,
		BaseUrl:     TraxmateLoRaCloudBaseUrl,
		logger:      logger,
		timeNow:     time.Now,
	}

	for _, option := range options {
		option(&client)
	}
	// Check after all the options have been applied
	if err := client.isSemtechLoRaCloudShutdown(); err != nil {
		return LoracloudClient{}, err
	}

	return client, nil
}

type LoracloudClientOptions func(*LoracloudClient)

func WithBaseUrl(baseUrl string) LoracloudClientOptions {
	return func(c *LoracloudClient) {
		c.BaseUrl = baseUrl
	}
}

func WithTimeNow(timeNow func() time.Time) LoracloudClientOptions {
	return func(c *LoracloudClient) {
		c.timeNow = timeNow
	}
}

func validateContext(ctx context.Context) error {
	_, ok := ctx.Value(decoder.PORT_CONTEXT_KEY).(uint8)
	if !ok {
		return ErrContextPortNotFound
	}
	devEui, ok := ctx.Value(decoder.DEVEUI_CONTEXT_KEY).(string)
	if !ok {
		return ErrContextDevEuiNotFound
	}
	fCount, ok := ctx.Value(decoder.FCNT_CONTEXT_KEY).(int)
	if !ok {
		return ErrContextFCountNotFound
	}
	if len(devEui) != 16 {
		return ErrContextInvalidDevEui
	}
	// check if devEui is a valid hex string
	hexCheck, err := hex.DecodeString(devEui)
	if err != nil || len(hexCheck) != 8 {
		return ErrContextInvalidDevEui
	}
	if fCount < 0 {
		return ErrContextInvalidFCount
	}
	return nil
}

func (m LoracloudClient) Solve(ctx context.Context, payload string) (*decoder.DecodedUplink, error) {
	if err := validateContext(ctx); err != nil {
		return nil, fmt.Errorf("context validation failed: %v", err)
	}

	port, ok := ctx.Value(decoder.PORT_CONTEXT_KEY).(uint8)
	if !ok {
		return nil, fmt.Errorf("port not found in context")
	}
	devEui, ok := ctx.Value(decoder.DEVEUI_CONTEXT_KEY).(string)
	if !ok {
		return nil, fmt.Errorf("devEui not found in context")
	}
	fCount, ok := ctx.Value(decoder.FCNT_CONTEXT_KEY).(int)
	if !ok {
		return nil, fmt.Errorf("fCount not found in context")
	}

	var timestamp *float64 = nil
	if m.BaseUrl == TraxmateLoRaCloudBaseUrl {
		// NOTE: Traxmate requires us to set the timestamp in order to solve the position,
		// Semtech used to do this on its own.
		unixTime := float64(m.timeNow().Unix())
		timestamp = &unixTime
	}

	decodedData, err := m.DeliverUplinkMessage(devEui, UplinkMsg{
		MsgType:   "updf",
		Port:      uint8(port),
		Payload:   payload,
		FCount:    uint32(fCount),
		Timestamp: timestamp,
	})

	if err != nil {
		return nil, fmt.Errorf("error delivering uplink message: %v", err)
	}
	return decoder.NewDecodedUplink([]decoder.Feature{decoder.FeatureTimestamp, decoder.FeatureGNSS, decoder.FeatureBuffered}, decodedData), err
}

func (m LoracloudClient) post(url string, body []byte) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creating loracloud request: %v", err)
	}

	request.Header.Set("Authorization", m.accessToken)
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	return client.Do(request)
}

// POST /api/v1/device/send
//
// Similar to the uplink/send API endpoint, but accepting a single uplink message of a device.
// The source device is denoted by the deveui field, the message by the uplink field.
//
// Request Headers
// Authorization – Required. Access token
//
// Request Body:
//
// The request body must contain the deveui and uplink fields.
//
//	{
//	  "deveui":  DEVEUI,     // Required. Source device EUI.
//	  "uplink":  UPLINK_MSG  // Required.
//
//	  ..
//	}
//
// deveui: Required, a valid device EUI.
//
// uplink: Required, UplinkMsg object.
//
// Response:
//
// Status Codes
// 200 OK – OK
//
// 401 Unauthorized – Authentication failed
//
// Response Headers
// Content-Type – ‘application/json’
//
// Response JSON:
//
// The response adheres to the Base Response Format. If successful, the result field keeps an UplinkResponse object. If errors were encountered, they are signaled in the errors field.
//
//	{
//	  "result": UPLINK_RESPONSE,  // Uplink response object for this EUI
//	  "errors": [ STRING, .. ]    // Error messages in case of error
//	}
//
// result: UplinkResponse instance detailing device state including information such as completed requests, files, stream records, and pending downlink messages.
//
// errors: If set and non-empty, error message in case the operation did not succeed.
func (m LoracloudClient) DeliverUplinkMessage(devEui string, uplinkMsg UplinkMsg) (*UplinkMsgResponse, error) {
	// validate uplinkMsg
	validate := validator.New()
	err := validate.Struct(uplinkMsg)
	if err != nil {
		return nil, fmt.Errorf("error validating uplink message: %v", err)
	}

	url := fmt.Sprintf("%v/api/v1/device/send", m.BaseUrl)

	// format devEui to match ^([0-9a-fA-F]){2}(-([0-9a-fA-F]){2}){7}$
	devEui = strings.ToUpper(devEui)
	if !strings.Contains(devEui, "-") {
		devEui = strings.Join([]string{
			devEui[0:2],
			devEui[2:4],
			devEui[4:6],
			devEui[6:8],
			devEui[8:10],
			devEui[10:12],
			devEui[12:14],
			devEui[14:16],
		}, "-")
	}

	body := map[string]any{
		"deveui": devEui,
		"uplink": uplinkMsg,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	response, err := m.post(url, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSendingRequest, err)
	}

	if response.StatusCode != http.StatusOK {
		responseJson := map[string]any{}
		err = json.NewDecoder(response.Body).Decode(&responseJson)
		if err != nil {
			return nil, fmt.Errorf("%w: HTTP %v", ErrUnexpectedStatusCode, response.StatusCode)
		}
		return nil, fmt.Errorf("%w: HTTP %v, %v", ErrUnexpectedStatusCode, response.StatusCode, responseJson)
	}

	var uplinkResponse UplinkMsgResponse

	if m.BaseUrl == TraxmateLoRaCloudBaseUrl {
		// NOTE: Traxmate LoRaCloud returns a nested response structure
		// {
		// 	"result": {
		// 		"10-CE-45-FF-FE-01-CE-53": {
		// 			UplinkMsgResponse
		// 		}
		// 	}
		// }
		var traxmateResponse struct {
			Result map[string]UplinkMsgResponse `json:"result"`
		}

		err = json.NewDecoder(response.Body).Decode(&traxmateResponse)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDecodingResponse, err)
		}

		if len(traxmateResponse.Result) != 1 {
			return nil, fmt.Errorf("%w: expected exactly one device EUI in the Traxmate LoRaCloud response, got %d", ErrMultipleDevicesInResponse, len(traxmateResponse.Result))
		}

		nestedUplinkResponse, ok := traxmateResponse.Result[devEui]
		if !ok {
			return nil, fmt.Errorf("%w: device EUI %s not found in Traxmate LoRaCloud response", ErrDeviceEuiNotInResponse, devEui)
		}

		uplinkResponse = nestedUplinkResponse
	} else {
		err = json.NewDecoder(response.Body).Decode(&uplinkResponse)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDecodingResponse, err)
		}
	}

	// remove the '-' from the devEui
	uplinkResponse.Result.Deveui = strings.ReplaceAll(uplinkResponse.Result.Deveui, "-", "")

	return &uplinkResponse, nil
}

type UplinkMsg struct {
	MsgType                 string     `json:"msgtype" validate:"required"`
	FCount                  uint32     `json:"fcnt" validate:"required"`
	Port                    uint8      `json:"port" validate:"required"`
	Payload                 string     `json:"payload" validate:"required"` // HEX string with LoRaWAN message payload
	DR                      *uint8     `json:"dr,omitempty"`
	Frequency               *uint32    `json:"freq,omitempty"`
	Timestamp               *float64   `json:"timestamp,omitempty"` // RX timestamp in seconds, UTC
	DNMTU                   *uint32    `json:"dn_mtu,omitempty"`
	GNSSCaptureTime         *float64   `json:"gnss_capture_time,omitempty"`
	GNSSCaptureTimeAccuracy *float64   `json:"gnss_capture_time_accuracy,omitempty"`
	GNSSAssistPosition      *[]float64 `json:"gnss_assist_position,omitempty"`
	GNSSAssistAltitude      *float64   `json:"gnss_assist_altitude,omitempty"`
	GNSSUse2DSolver         *bool      `json:"gnss_use_2D_solver,omitempty"`
}

// An “Uplink Response” object reflects the current device state as well as new items resulting from the submitted uplink message.
// In addition to the info_fields and log_messages fields of the “DeviceInfo” object (DeviceInfo) it contains fields which signal the following state changes and
// completed items that are due to an uplink message:
//
// operation: Required, one of “gnss”, “wifi”, “modem”, “other”.
// In case of “other”, the message is application specific and can be forwarded to the application.
// In case of “modem”, “gnss” and “wifi” the message may not be handled by the application.
//
// file: Optional. FileObject. Contains the data of a file upload if the uplink message led to the completion of a file upload session.
//
// stream_records: Optional. StreamUpdate. Contains a data stream update entry which signals re-assembled stream records. The array might be empty.
//
// position_solution: Optional. PositionSolution. Contains the solution of a valid position calculation.
//
// fulfilled_requests: Optional. List of PendingRequest objects. Contains all pending requests which have been completed by this uplink.
//
// dnlink: Optional. LoRaDnlink. If set, the downlink which has to be scheduled with the network server.
//
// fports: Required. Current DevicePorts settings.
//
// info_fields: Required. List of InfoFields objects
//
// pending_requests: Required. List of PendingRequest objects
//
// log_messages: Required. List of log messages, LogMessage
type UplinkMsgResponse struct {
	Result struct {
		Deveui          string `json:"deveui"`
		PendingRequests struct {
			Requests []any `json:"requests"`
			ID       int   `json:"id"`
			Updelay  int   `json:"updelay"`
			Upcount  int   `json:"upcount"`
		} `json:"pending_requests"`
		InfoFields struct {
			Rfu     any `json:"rfu"`
			Temp    any `json:"temp"`
			Charge  any `json:"charge"`
			Deveui  any `json:"deveui"`
			Region  any `json:"region"`
			Rxtime  any `json:"rxtime"`
			Signal  any `json:"signal"`
			Status  any `json:"status"`
			Uptime  any `json:"uptime"`
			Adrmode any `json:"adrmode"`
			Alcsync struct {
				Value struct {
					Time  int `json:"time"`
					Token int `json:"token"`
				} `json:"value"`
				Timestamp float64 `json:"timestamp"`
			} `json:"alcsync"`
			Chipeui any `json:"chipeui"`
			Joineui any `json:"joineui"`
			Session struct {
				Value     int     `json:"value"`
				Timestamp float64 `json:"timestamp"`
			} `json:"session"`
			Voltage  any `json:"voltage"`
			Crashlog struct {
				Value     string  `json:"value"`
				Timestamp float64 `json:"timestamp"`
			} `json:"crashlog"`
			Firmware struct {
				Value struct {
					Fwcrc       string `json:"fwcrc"`
					Fwtotal     int    `json:"fwtotal"`
					Fwcompleted int    `json:"fwcompleted"`
				} `json:"value"`
				Timestamp float64 `json:"timestamp"`
			} `json:"firmware"`
			Interval any `json:"interval"`
			Rstcount struct {
				Value     int     `json:"value"`
				Timestamp float64 `json:"timestamp"`
			} `json:"rstcount"`
			Appstatus any `json:"appstatus"`
			Streampar any `json:"streampar"`
		} `json:"info_fields"`
		LogMessages []any `json:"log_messages"`
		Fports      struct {
			Dmport     int `json:"dmport"`
			Gnssport   int `json:"gnssport"`
			Wifiport   int `json:"wifiport"`
			Fragport   int `json:"fragport"`
			Streamport int `json:"streamport"`
			Gnssngport int `json:"gnssngport"`
		} `json:"fports"`
		Dnlink            any   `json:"dnlink"`
		FulfilledRequests []any `json:"fulfilled_requests"`
		CancelledRequests []any `json:"cancelled_requests"`
		File              any   `json:"file"`
		StreamRecords     any   `json:"stream_records"`
		PositionSolution  struct {
			Llh             []float64 `json:"llh"`
			Accuracy        float64   `json:"accuracy"`
			Ecef            []float64 `json:"ecef"`
			Gdop            float64   `json:"gdop"`
			CaptureTimeGps  float64   `json:"capture_time_gps"`
			CaptureTimeUtc  float64   `json:"capture_time_utc"`
			CaptureTimesGps []float64 `json:"capture_times_gps"`
			CaptureTimesUtc []float64 `json:"capture_times_utc"`
			Timestamp       float64   `json:"timestamp"`
			AlgorithmType   string    `json:"algorithm_type"`
		} `json:"position_solution"`
		Operation string `json:"operation"`
	} `json:"result"`
}

var _ decoder.UplinkFeatureTimestamp = &UplinkMsgResponse{}
var _ decoder.UplinkFeatureGNSS = &UplinkMsgResponse{}
var _ decoder.UplinkFeatureBuffered = &UplinkMsgResponse{}

func (p UplinkMsgResponse) GetTimestamp() *time.Time {
	var captureTs float64
	if p.Result.PositionSolution.AlgorithmType == "gnssng" {
		if p.Result.PositionSolution.CaptureTimeUtc != 0 {
			captureTs = p.Result.PositionSolution.CaptureTimeUtc
		}
	}

	if captureTs == 0 {
		return nil
	}

	seconds := int64(captureTs)
	nanoseconds := int64((captureTs - float64(seconds)) * 1e9)
	timestamp := time.Unix(seconds, nanoseconds)
	return &timestamp
}

func (p UplinkMsgResponse) GetLatitude() float64 {
	if len(p.Result.PositionSolution.Llh) > 0 {
		return p.Result.PositionSolution.Llh[0]
	}
	return 0
}

func (p UplinkMsgResponse) GetLongitude() float64 {
	if len(p.Result.PositionSolution.Llh) > 1 {
		return p.Result.PositionSolution.Llh[1]
	}
	return 0
}

func (p UplinkMsgResponse) GetAltitude() float64 {
	if len(p.Result.PositionSolution.Llh) > 2 {
		return p.Result.PositionSolution.Llh[2]
	}
	return 0
}

func (p UplinkMsgResponse) GetAccuracy() *float64 {
	return &p.Result.PositionSolution.Accuracy
}

func (p UplinkMsgResponse) GetTTF() *time.Duration {
	return nil
}

func (p UplinkMsgResponse) GetPDOP() *float64 {
	return &p.Result.PositionSolution.Gdop
}

func (p UplinkMsgResponse) GetSatellites() *uint8 {
	return nil
}

func (p UplinkMsgResponse) GetBufferLevel() *uint16 {
	return nil
}

func (p UplinkMsgResponse) IsBuffered() bool {
	if p.GetTimestamp() != nil && p.GetTimestamp().Before(time.Now().Add(-1*time.Minute)) {
		return true
	}
	return false
}
