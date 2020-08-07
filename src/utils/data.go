package utils

type DataStats struct {
	CpuStats  []float64
	MemStats  []float64
	DiskStats [][]string
	NetStats  map[string][]float64
	FieldSet  string
}
