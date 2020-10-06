package todoist

import (
	"log"
	"time"
)

// timeTrack measures the execution time of a method
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

// p<Type> functions are used to convert basic types to their pointer equivalents
func pString(s string) *string {
	return &s
}

func pInt64(i int64) *int64 {
	return &i
}
