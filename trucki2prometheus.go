package main

import (
	"os"
	"io"
	"fmt"
	"flag"
	"maps"
	"time"
	"runtime"
	"strings"
	"strconv"
	"net/http"
	"encoding/json"
	"net/http/httptest"
	
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	t2PromVersion 		string
	buildDate 			string
	goVersion 			string = runtime.Version()
	
	listenPort			string
	truckiAddress		string
	scrapeInterval		int64
	
	latestTruckiMetrics truckiMetrics
	metrics 			*promMetrics
)


func main() {
	fmt.Println("Hello from planet earth, this is Trucki2Prometheus")
	fmt.Printf("git-build-hash: %s; build-date: %s; go-version: %s\n", t2PromVersion, buildDate, goVersion)
	
	flag.StringVar(&listenPort, "p", "8080", "-p <t2p http listen port>; eg. -p 8080")
	flag.StringVar(&truckiAddress, "t", "", "-t <Trucki stick IP address or hostname>; eg. -t 192.168.178.58")
	flag.Int64Var(&scrapeInterval, "i", 5, "-i <scrape interval in seconds>; eg. -i 15")
	flag.Parse()
	
	if truckiAddress == "" {
		fmt.Println("No Trucki stick hostname set, please provide the Trucki stick IP address or hostname with the -t command flag")
		os.Exit(-1)
	}
	
	if scrapeInterval == 5 {
		fmt.Println("No custom scrape interval provided, defaulting to 5 seconds")
	
	} else if scrapeInterval < 1 {
		fmt.Println("Invalid scrape interval provided, please provide a number in seconds larger than zero!")
		os.Exit(-1)
	}
	
	//
	// Register Prometheus client registry & all the metrics
	// prior to starting to scrape the Truck stick & the HTTP server
	registry := registerPrometheusMetrics()
	
	//
	// Scrape the Trucki stick in a loop forever
	// on a background thread
	go func() {
		for {
			metrics, scrapeErr := scrapeTrucki()
			if scrapeErr != nil {
				fmt.Println("Failed to scrape Trucki stick:", scrapeErr)
			
			} else {
				latestTruckiMetrics = *metrics
			}
			
			time.Sleep(time.Duration(scrapeInterval) * time.Second)
		}
	}()
	
	http.Handle("/metrics", http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		rec := httptest.NewRecorder()
		promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}).ServeHTTP(rec, request)
		fmt.Printf("%s %s -> %d\n", request.Method, request.URL.Path, rec.Code)
		
		// Copy all the response headers back into the original response
		maps.Copy(response.Header(), rec.Header())
		
		// Write the HTTP response status code back into the original response 
		response.WriteHeader(rec.Code)
		
		// Write the HTTP response body back into the original response
		response.Write(rec.Body.Bytes())
		
	}))
	
	fmt.Printf("Starting trucki2prometheus exporter on :%s\n", listenPort)
	if listenErr := http.ListenAndServe(":"+listenPort, nil); listenErr != nil {
		fmt.Println("Failed to start HTTP server:", listenErr)
		os.Exit(-1)
	}
}

type truckiMetrics struct {
	VoltageGrid					float64	`json:"VGRID"`
	VoltageBattery 				float64	`json:"VBAT"`
	SetACPower 					int64	`json:"SETACPOWER"` 
	Temperature 				int64	`json:"TEMP"`
	PowerLimit 					int64	`json:"POWERLIMIT"`
	Sun2RoundTrip 				string	`json:"SUN2ROUNDTRIP"`
	Sun2RoundTripInt			int64
	Sun2Setpoint 				int64	`json:"SUN2SETPOINT"`
	Sun2PowerLimit 				int64	`json:"SUN2POWERLIMIT"`
	Sun3RoundTrip 				string	`json:"SUN3ROUNDTRIP"`
	Sun3RoundTripInt			int64
	Sun3Setpoint 				int64	`json:"SUN3SETPOINT"`
	Sun3PowerLimit 				int64	`json:"SUN3POWERLIMIT"`
	MeterReadout 				int64	`json:"METERREADOUT"`
	DayEnergy 					float64	`json:"DAYENERGY"`
	TotalEnergy 				float64	`json:"TOTALENERGY"`
	MeterDayEnergy 				float64	`json:"METERDAYENERGY"`
	ACPower 					string	`json:"ACPOWER"`
	ACPowerFloat				float64
	ACPowerSun2					string	`json:"ACPOWERSUN2"`
	ACPowerSun2Float			float64
	ACPowerSun3					string	`json:"ACPOWERSUN3"`
	ACPowerSun3Float			float64
	ZeroExportControlPower 		string	`json:"ZEPCPOWER"`
	ZeroExportControlPowerFloat float64
	MeterPower 					string	`json:"METERPOWER"`
	MeterPowerFloat				float64
	WiFiState 					string	`json:"WIFI"`
	WiFiStateInt				int64
	RSSI 						string	`json:"RSSI"`
	RSSIStateInt 				int64
}

func scrapeTrucki() (*truckiMetrics, error) {
	url := fmt.Sprintf("http://%s/jsonlive", truckiAddress)
	req, reqGenErr := http.NewRequest("GET", url, nil)
	if reqGenErr != nil {
		return nil, reqGenErr
	}
	
	httpClient := http.Client {
		Timeout: 15 * time.Second,
	}
	
	resp, reqErr := httpClient.Do(req)
	if reqErr != nil {
		return nil, reqErr
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Trucki stick scrape request returned with non 200 Ok HTTP status code. Response status-code: %d", resp.StatusCode)
	}
	
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}
	
	var truckiMetrics truckiMetrics
	jsonErr := json.Unmarshal(body, &truckiMetrics)
	if jsonErr != nil {
		return nil, jsonErr
	}
	
	//
	// Extract various values based on the strings being returned
	// in the JSON
	acPowerComponents := strings.Split(truckiMetrics.ACPower, " ")
	acPowerFloat, acPowerParseErr := strconv.ParseFloat(acPowerComponents[0], 64)
	if acPowerParseErr != nil {
		fmt.Println("Failed to extract number from JSON key 'ACPOWER':", acPowerParseErr)
		truckiMetrics.ACPowerFloat = -1.0
	
	} else {
		truckiMetrics.ACPowerFloat = acPowerFloat
	}
	
	sun2ACPowerComponents := strings.Split(truckiMetrics.ACPowerSun2, " ")
	sun2ACPowerFloat, sun2ACPowerParseErr := strconv.ParseFloat(sun2ACPowerComponents[0], 64)
	if sun2ACPowerParseErr != nil {
		fmt.Println("Failed to extract number from JSON key 'ACPOWERSUN2':", sun2ACPowerParseErr)
		truckiMetrics.ACPowerSun2Float = -1.0
	
	} else {
		truckiMetrics.ACPowerSun2Float = sun2ACPowerFloat
	}
	
	sun3ACPowerComponents := strings.Split(truckiMetrics.ACPowerSun3, " ")
	sun3ACPowerFloat, sun3ACPowerParseErr := strconv.ParseFloat(sun3ACPowerComponents[0], 64)
	if sun3ACPowerParseErr != nil {
		fmt.Println("Failed to extract number from JSON key 'ACPOWERSUN3':", sun3ACPowerParseErr)
		truckiMetrics.ACPowerSun3Float = -1.0
	
	} else {
		truckiMetrics.ACPowerSun3Float = sun3ACPowerFloat
	}
	
	zepcComponents := strings.Split(truckiMetrics.ZeroExportControlPower, " ")
	zepcACPowerFloat, zepcACPowerParseErr := strconv.ParseFloat(zepcComponents[0], 64)
	if zepcACPowerParseErr != nil {
		fmt.Println("Failed to extract number from JSON key 'ZEPCPOWER':", zepcACPowerParseErr)
		truckiMetrics.ZeroExportControlPowerFloat = -1.0
	
	} else {
		truckiMetrics.ZeroExportControlPowerFloat = zepcACPowerFloat
	}
	
	meterPowerComponents := strings.Split(truckiMetrics.MeterPower, " ")
	meterPowerFloat, meterPowerParseErr := strconv.ParseFloat(meterPowerComponents[0], 64)
	if meterPowerParseErr != nil {
		fmt.Println("Failed to extract number from JSON key 'METERPOWER':", meterPowerParseErr)
		truckiMetrics.MeterPowerFloat = -1.0
	
	} else {
		truckiMetrics.ZeroExportControlPowerFloat = meterPowerFloat
	}
	
	if truckiMetrics.WiFiState == "DISCONNECTED" {
		truckiMetrics.WiFiStateInt = 0
		
	} else if truckiMetrics.WiFiState == "CONNECTED" {
		truckiMetrics.WiFiStateInt = 1
	}
	
	if truckiMetrics.RSSI == "Unusable" {
		truckiMetrics.RSSIStateInt = 0
	
	} else if truckiMetrics.RSSI == "Not good" || truckiMetrics.RSSI == "Very Good" {
		truckiMetrics.RSSIStateInt = 1
	
	} else if truckiMetrics.RSSI == "Okay" {
		truckiMetrics.RSSIStateInt = 2
	
	} else if truckiMetrics.RSSI == "Very good" || truckiMetrics.RSSI == "Very Good" {
		truckiMetrics.RSSIStateInt = 3
	
	} else if truckiMetrics.RSSI == "Amazing" {
		truckiMetrics.RSSIStateInt = 4
	}
	
	if metrics != nil {
		metrics.voltageGrid.Set(truckiMetrics.VoltageGrid)
		metrics.voltageBattery.Set(truckiMetrics.VoltageBattery)
		metrics.setACPower.Set(float64(truckiMetrics.SetACPower))
		metrics.temperature.Set(float64(truckiMetrics.Temperature))
		metrics.powerLimit.Set(float64(truckiMetrics.PowerLimit))
		// metrics.sun2RoundTrip.Set(float64(truckiMetrics.Sun2RoundTripInt))
		metrics.sun2SetPoint.Set(float64(truckiMetrics.Sun2Setpoint))
		metrics.sun2PowerLimit.Set(float64(truckiMetrics.Sun2PowerLimit))
		// metrics.sun3RoundTrip.Set(float64(truckiMetrics.Sun3RoundTripInt))
		metrics.sun3SetPoint.Set(float64(truckiMetrics.Sun3Setpoint))
		metrics.sun3PowerLimit.Set(float64(truckiMetrics.Sun3PowerLimit))
		metrics.powerMeterReadout.Set(float64(truckiMetrics.MeterReadout))
		metrics.dayEnergyOutput.Set(truckiMetrics.DayEnergy)
		metrics.totalEnergyOutput.Set(truckiMetrics.TotalEnergy)
		metrics.powerMeterDayEnergy.Set(truckiMetrics.MeterDayEnergy)
		metrics.inverterACPowerOutput.Set(truckiMetrics.ACPowerFloat)
		metrics.sun2ACPowerOutput.Set(truckiMetrics.ACPowerSun2Float)
		metrics.sun3ACPowerOutput.Set(truckiMetrics.ACPowerSun3Float)
		metrics.zeroExportControlPower.Set(truckiMetrics.ZeroExportControlPowerFloat)
		metrics.powerMeterPower.Set(truckiMetrics.MeterPowerFloat)
		metrics.wifiState.Set(float64(truckiMetrics.WiFiStateInt))
		metrics.wifiRSSI.Set(float64(truckiMetrics.RSSIStateInt))
	}
	
	return &truckiMetrics, nil
}
