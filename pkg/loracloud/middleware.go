package loracloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type LoracloudMiddleware struct {
	accessToken string
	BaseUrl     string
}

func NewLoracloudMiddleware(accessToken string) LoracloudMiddleware {
	return LoracloudMiddleware{accessToken: accessToken, BaseUrl: "https://mgs.loracloud.com"}
}

func (m LoracloudMiddleware) post(url string, body []byte) (*http.Response, error) {
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
func (m LoracloudMiddleware) DeliverUplinkMessage(devEui string, uplinkMsg UplinkMsg) (*UplinkMsgResponse, error) {
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

	body := map[string]interface{}{
		"deveui": devEui,
		"uplink": uplinkMsg,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	response, err := m.post(url, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error sending request to loracloud: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		responseJson := map[string]interface{}{}
		err = json.NewDecoder(response.Body).Decode(&responseJson)
		if err != nil {
			return nil, fmt.Errorf("unexpected status code returned by loracloud: HTTP %v", response.StatusCode)
		}
		return nil, fmt.Errorf("unexpected status code returned by loracloud: HTTP %v, %v", response.StatusCode, responseJson)
	}

	uplinkResponse := UplinkMsgResponse{}
	err = json.NewDecoder(response.Body).Decode(&uplinkResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding loracloud response: %v", err)
	}

	// remove the '-' from the devEui
	uplinkResponse.Result.Deveui = strings.ReplaceAll(uplinkResponse.Result.Deveui, "-", "")

	return &uplinkResponse, nil
}

type UplinkMsg struct {
	MsgType                 string     `json:"msgtype"`
	FCount                  uint32     `json:"fcnt"`
	Port                    uint8      `json:"port"`
	Payload                 string     `json:"payload"` // HEX string with LoRaWAN message payload
	DR                      *uint8     `json:"dr"`
	Frequency               *uint32    `json:"freq"`
	Timestamp               *float64   `json:"timestamp"` // RX timestamp in seconds, UTC
	DNMTU                   *uint32    `json:"dn_mtu"`
	GNSSCaptureTime         *float64   `json:"gnss_capture_time"`
	GNSSCaptureTimeAccuracy *float64   `json:"gnss_capture_time_accuracy"`
	GNSSAssistPosition      *[]float64 `json:"gnss_assist_position"`
	GNSSAssistAltitude      *float64   `json:"gnss_assist_altitude"`
	GNSSUse2DSolver         *bool      `json:"gnss_use_2D_solver"`
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
			Requests []interface{} `json:"requests"`
			ID       int           `json:"id"`
			Updelay  int           `json:"updelay"`
			Upcount  int           `json:"upcount"`
		} `json:"pending_requests"`
		InfoFields struct {
			Rfu     interface{} `json:"rfu"`
			Temp    interface{} `json:"temp"`
			Charge  interface{} `json:"charge"`
			Deveui  interface{} `json:"deveui"`
			Region  interface{} `json:"region"`
			Rxtime  interface{} `json:"rxtime"`
			Signal  interface{} `json:"signal"`
			Status  interface{} `json:"status"`
			Uptime  interface{} `json:"uptime"`
			Adrmode interface{} `json:"adrmode"`
			Alcsync struct {
				Value struct {
					Time  int `json:"time"`
					Token int `json:"token"`
				} `json:"value"`
				Timestamp float64 `json:"timestamp"`
			} `json:"alcsync"`
			Chipeui interface{} `json:"chipeui"`
			Joineui interface{} `json:"joineui"`
			Session struct {
				Value     int     `json:"value"`
				Timestamp float64 `json:"timestamp"`
			} `json:"session"`
			Voltage  interface{} `json:"voltage"`
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
			Interval interface{} `json:"interval"`
			Rstcount struct {
				Value     int     `json:"value"`
				Timestamp float64 `json:"timestamp"`
			} `json:"rstcount"`
			Appstatus interface{} `json:"appstatus"`
			Streampar interface{} `json:"streampar"`
		} `json:"info_fields"`
		LogMessages []interface{} `json:"log_messages"`
		Fports      struct {
			Dmport     int `json:"dmport"`
			Gnssport   int `json:"gnssport"`
			Wifiport   int `json:"wifiport"`
			Fragport   int `json:"fragport"`
			Streamport int `json:"streamport"`
			Gnssngport int `json:"gnssngport"`
		} `json:"fports"`
		Dnlink            interface{}   `json:"dnlink"`
		FulfilledRequests []interface{} `json:"fulfilled_requests"`
		CancelledRequests []interface{} `json:"cancelled_requests"`
		File              interface{}   `json:"file"`
		StreamRecords     interface{}   `json:"stream_records"`
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
