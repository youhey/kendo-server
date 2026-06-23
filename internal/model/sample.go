package model

type ADXL struct {
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Z    float64 `json:"z"`
	Mag  float64 `json:"mag"`
	RMS  float64 `json:"rms"`
	Peak float64 `json:"peak"`
}

type Piezo struct {
	Raw  int64 `json:"raw"`
	Min  int64 `json:"min"`
	Max  int64 `json:"max"`
	Peak int64 `json:"peak"`
}

type Sample struct {
	NodeID     string `json:"node_id"`
	Seq        *int64 `json:"seq,omitempty"`
	MeasuredAt string `json:"measured_at"`
	ADXL       ADXL   `json:"adxl"`
	Piezo      Piezo  `json:"piezo"`
}

type SampleRecord struct {
	ID         int64  `json:"id"`
	NodeID     string `json:"node_id"`
	Seq        *int64 `json:"seq,omitempty"`
	MeasuredAt string `json:"measured_at"`
	ReceivedAt string `json:"received_at"`
	ADXL       ADXL   `json:"adxl"`
	Piezo      Piezo  `json:"piezo"`
	CreatedAt  string `json:"created_at"`
}
