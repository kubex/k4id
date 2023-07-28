package k4id

import (
	"log"
	"testing"
)

func benchmarkIDGeneration(idGen Generator, b *testing.B) {
	// run the function b.N times
	for n := 0; n < b.N; n++ {
		idGen.New()
	}
}
func BenchmarkIDGeneration(b *testing.B) {
	benchmarkIDGeneration(globalIDHost, b)
}

func BenchmarkIDGenerationNano(b *testing.B) {
	idGen := DefaultGenerator()
	idGen.idLength = 16
	idGen.SetTimeSize(TimeGeneratorNano)
	benchmarkIDGeneration(idGen, b)
}

func BenchmarkIDGenerationMicro(b *testing.B) {
	idGen := DefaultGenerator()
	idGen.SetTimeSize(TimeGeneratorMicro)
	benchmarkIDGeneration(idGen, b)
}
func BenchmarkIDGenerationMilli(b *testing.B) {
	idGen := DefaultGenerator()
	idGen.SetTimeSize(TimeGeneratorMilli)
	benchmarkIDGeneration(idGen, b)
}

func BenchmarkIDGenerationSmall(b *testing.B) {
	idGen := DefaultGenerator()
	idGen.SetBaseLength(6)
	idGen.SetTimeSize(TimeGeneratorSecond)
	benchmarkIDGeneration(idGen, b)
}

func TestID(t *testing.T) {
	routines := 10
	iter := 10000
	test := routines * iter
	idStream := make(chan ID, test)

	for i := 0; i < routines; i++ {
		go func(generator Generator) {
			for i := 0; i < iter; i++ {
				gen := generator.New()
				idStream <- gen
				if i%500 == 0 {
					log.Println(gen.String())
					//	log.Println(gen.UUID())
				}
			}
		}(DefaultGenerator())
	}
	generated := map[string]bool{}
	lastProcess := 0
	for processed := 0; processed < test; processed++ {
		gen := <-idStream
		if _, found := generated[gen.String()]; found {
			t.Fatal("Duplicate ID generated ", gen)
		} else {
			generated[gen.String()] = true
		}
		lastProcess = processed
		if !FromString(gen.String()).IsValid() {
			t.Fatal("Invalid verification")
		}
		if FromUUID(gen.UUID()).UUID() != gen.UUID() {
			t.Fatal("Invalid uuid: ", gen.UUID(), " - ", gen.String())
		}
	}

	log.Println("Processed ", lastProcess+1, " IDs")

}

func TestIDUUID(t *testing.T) {
	h := DefaultGenerator()
	h.SetBaseLength(19)

	for i := 0; i < 10000; i++ {
		id := h.New()
		originalUUID := id.UUID()
		id2 := FromUUID(originalUUID)
		log.Println(id.UUID(), " : ", id2.UUID())
		log.Println(id.String(), " : ", id2.String())
		if id.String() != id2.String() {
			log.Fatal("UUID conversion failed")
		}
		if id.IsValid() && !id2.IsValid() {
			log.Fatal("UUID conversion failed checksum")
		}
	}
}

func TestUUIDImport(t *testing.T) {
	uuid := "690482d9-1dda-42b8-b6e1-8df9d16baf05"
	log.Println(FromUUID(uuid).String())
}
