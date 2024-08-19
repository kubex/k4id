package k4id

import (
	"log"
	"testing"
	"time"
)

func TestTimeGenerator(t *testing.T) {
	tests := []struct {
		name   string
		gen    TimeGenerator
		format string
	}{
		{"day", TimeGeneratorDay, time.DateOnly},
		{"hour", TimeGeneratorHour, "2006-01-02 15"},
		{"minute", TimeGeneratorMinute, "2006-01-02 15:04"},
		{"second", TimeGeneratorSecond, "2006-01-02 15:04:05"},
		{"milli", TimeGeneratorMilli, "2006-01-02 15:04:05.000"},
		{"micro", TimeGeneratorMicro, "2006-01-02 15:04:05.000000"},
		{"nano", TimeGeneratorNano, "2006-01-02 15:04:05.000000000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useTime := time.Date(2020, 1, 1, 3, 4, 5, 6, time.UTC)
			outID := tt.gen.Generate(useTime)
			log.Println(outID)
			parsed := tt.gen.Parse(outID)
			if parsed.Format(tt.format) != useTime.Format(tt.format) {
				t.Errorf("Time not equal for %s -- %s %s ", tt.name, parsed.Format(tt.format), useTime.Format(tt.format))
			}
		})
	}
}

func TestIDToTime(t *testing.T) {
	gen := DefaultGenerator()

	tests := []struct {
		name   string
		gen    TimeGenerator
		format string
	}{
		{"day", TimeGeneratorDay, time.DateOnly},
		{"hour", TimeGeneratorHour, "2006-01-02 15"},
		{"minute", TimeGeneratorMinute, "2006-01-02 15:04"},
		{"second", TimeGeneratorSecond, "2006-01-02 15:04:05"},
		{"milli", TimeGeneratorMilli, "2006-01-02 15:04:05.000"},
		{"micro", TimeGeneratorMicro, "2006-01-02 15:04:05.000000"},
		{"nano", TimeGeneratorNano, "2006-01-02 15:04:05.000000000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen.SetTimeSize(tt.gen)
			useTime := time.Now()
			gen.SetTime(useTime)
			newID := gen.New().String()
			parsed := gen.ExtractTime(newID)
			if parsed.Format(tt.format) != useTime.Format(tt.format) {
				t.Errorf("Time not equal for %s -- %s %s ", tt.name, parsed.Format(tt.format), useTime.Format(tt.format))
			}
		})
	}
}
