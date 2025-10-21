package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)


type promMetrics struct {
	voltageGrid 			prometheus.Gauge
	voltageBattery 			prometheus.Gauge
	setACPower				prometheus.Gauge
	temperature 			prometheus.Gauge
	powerLimit 				prometheus.Gauge
	sun2RoundTrip 			prometheus.Gauge
	sun2SetPoint 			prometheus.Gauge
	sun2PowerLimit 			prometheus.Gauge
	sun3RoundTrip 			prometheus.Gauge
	sun3SetPoint 			prometheus.Gauge
	sun3PowerLimit 			prometheus.Gauge
	powerMeterReadout 		prometheus.Gauge
	dayEnergyOutput 		prometheus.Gauge
	totalEnergyOutput 		prometheus.Gauge
	powerMeterDayEnergy 	prometheus.Gauge
	inverterACPowerOutput 	prometheus.Gauge
	sun2ACPowerOutput 		prometheus.Gauge
	sun3ACPowerOutput 		prometheus.Gauge
	zeroExportControlPower 	prometheus.Gauge
	powerMeterPower 		prometheus.Gauge
	wifiState 				prometheus.Gauge
	wifiRSSI 				prometheus.Gauge
}


func registerPrometheusMetrics() *prometheus.Registry {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	
	metrics = &promMetrics {
		voltageGrid: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "trucki_voltage_grid",
			Help: "AC grid voltage (V) measured by the inverter",
		}),
		voltageBattery: prometheus.NewGauge(prometheus.GaugeOpts {
			Name: "trucki_voltage_battery",
			Help: "DC battery voltage (V) measured by the inverter",
		}),
		setACPower: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "trucki_set_ac_power",
			Help: "AC power output target in watts (W) set by Trucki stick",
		}),
		temperature: prometheus.NewGauge(prometheus.GaugeOpts {
			Name: "trucki_inverter_temperature",
			Help: "Temperature of the inverter in celcius (°C)",
		}),
		powerLimit: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "trucki_power_limit",
			Help: "Inverter AC output power limit in watts (W) set by Trucki stick",
		}),
		sun2RoundTrip: prometheus.NewGauge(prometheus.GaugeOpts {
			Name: "trucki_sun_2_round_trip",
			Help: "Latency or round trip time for a single packet to a second Lumentree (sun2) inverter in milliseconds (ms)",
		}),
		sun2SetPoint: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "trucki_sun_2_set_point",
			Help: "Set point in watts (W) to grid for a second Lumentree (sun2) inverter",
		}),
		sun2PowerLimit: prometheus.NewGauge(prometheus.GaugeOpts {
			Name: "trucki_sun_2_power_limit",
			Help: "Max power limit in watts (W) for a second Lumentree (sun2) inverter",
		}),
		sun3RoundTrip: prometheus.NewGauge(prometheus.GaugeOpts {
			Name: "trucki_sun_3_round_trip",
			Help: "Latency or round trip time for a single packet to a third Lumentree (sun3) inverter in milliseconds (ms)",
		}),
		sun3SetPoint: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "trucki_sun_3_set_point",
			Help: "Set point in watts (W) to grid for a third Lumentree (sun3) inverter",
		}),
		sun3PowerLimit: prometheus.NewGauge(prometheus.GaugeOpts {
			Name: "trucki_sun_3_power_limit",
			Help: "Max power limit in watts (W) for a third Lumentree (sun3) inverter",
		}),
		powerMeterReadout: prometheus.NewGauge(prometheus.GaugeOpts {
			Name: "trucki_power_meter_readout",
			Help: "Latency or round trip time for a single packet to the power meter in milliseconds (ms)",
		}),
		dayEnergyOutput: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "trucki_day_energy_grid_output",
			Help: "Energy output of the Lumentree Sun inverter to the grid based on calendar day borders (not last 24 hours!) in kilowatthours (kWh)",
		}),
		totalEnergyOutput: prometheus.NewGauge(prometheus.GaugeOpts {
			Name: "trucki_total_energy_grid_output",
			Help: "Total energy output of the Lumentree Sun inverter to the grid in kilowatthours (kWh)",
		}),
		powerMeterDayEnergy: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "trucki_power_meter_day_energy",
			Help: "Energy consumption measured by the power meter from the grid in kilowatthours (kWh)",
		}),
		inverterACPowerOutput: prometheus.NewGauge(prometheus.GaugeOpts {
			Name: "trucki_inverter_ac_power_output",
			Help: "AC power output to the grid by the inverter",
		}),
		sun2ACPowerOutput: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "trucki_sun_2_ac_power_output",
			Help: "AC power output from a second Lumentree (sun2) inverter in watts (W)",
		}),
		sun3ACPowerOutput: prometheus.NewGauge(prometheus.GaugeOpts {
			Name: "trucki_sun_3_ac_power_output",
			Help: "AC power output from a third Lumentree (sun3) inverter in watts (W)",
		}),
		zeroExportControlPower: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "trucki_zero_export_control_power",
			Help: "Calculated zero-export-power-control output power calculated by Trucki in watts (W)",
		}),
		powerMeterPower: prometheus.NewGauge(prometheus.GaugeOpts {
			Name: "trucki_power_meter_power",
			Help: "Power consumption measured by the power meter in watts (W)",
		}),
		wifiState: prometheus.NewGauge(prometheus.GaugeOpts {
			Name: "trucki_wifi_state",
			Help: "WiFi state of the Trucki stick. Many options, eg. 0. 'Disconnected'; 1. 'Connected'",
		}),
		wifiRSSI: prometheus.NewGauge(prometheus.GaugeOpts {
			Name: "trucki_wifi_rssi",
			Help: "WiFi RSSI (received signal streth indicator) of the Trucki stick. Possible values are 0. 'Bad'; 1. 'Not good'; 2. 'Okay'; 3. 'Fairly good'; 4. 'Very good'",
		}),
	}
	
	registry.MustRegister(metrics.voltageGrid)
	registry.MustRegister(metrics.voltageBattery)
	registry.MustRegister(metrics.setACPower)
	registry.MustRegister(metrics.temperature)
	registry.MustRegister(metrics.powerLimit)
	registry.MustRegister(metrics.sun2RoundTrip)
	registry.MustRegister(metrics.sun2SetPoint)
	registry.MustRegister(metrics.sun2PowerLimit)
	registry.MustRegister(metrics.sun3RoundTrip)
	registry.MustRegister(metrics.sun3SetPoint)
	registry.MustRegister(metrics.sun3PowerLimit)
	registry.MustRegister(metrics.powerMeterReadout)
	registry.MustRegister(metrics.dayEnergyOutput)
	registry.MustRegister(metrics.totalEnergyOutput)
	registry.MustRegister(metrics.powerMeterDayEnergy)
	registry.MustRegister(metrics.inverterACPowerOutput)
	registry.MustRegister(metrics.sun2ACPowerOutput)
	registry.MustRegister(metrics.sun3ACPowerOutput)
	registry.MustRegister(metrics.zeroExportControlPower)
	registry.MustRegister(metrics.powerMeterPower)
	registry.MustRegister(metrics.wifiState)
	registry.MustRegister(metrics.wifiRSSI)
	
	return registry
}
