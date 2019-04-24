package util

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"os"

	"path"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/mitchellh/go-homedir"
)

func ToBytes(in interface{}) []byte {
	b, err := json.Marshal(in)
	if err != nil {
		panic("toBytes is unable to marshal input")
	}
	return b
}

func DumpJSON(descr string, in interface{}) {
	fmt.Printf("%s ------------------------- json dump start ---------------------------------------\n", descr)
	out, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	os.Stdout.Write(out)
	fmt.Printf("\n%s ------------------------- json dump end ---------------------------------------\n\n", descr)
}

func Dump(descr string, in interface{}) {
	fmt.Printf("%s ------------------------- dump start ---------------------------------------\n", descr)
	spew.Dump(in)
	fmt.Printf("%s -------------------------  dump end  ---------------------------------------\n\n", descr)
}

func SafeUnquote(in string) (string, error) {
	if strings.HasPrefix(in, "\"") && strings.HasSuffix(in, "\"") {
		q, err := strconv.Unquote(in)
		if err != nil {
			return "", err
		}

		return q, nil
	}

	return in, nil
}

//WaitForCondition is a testify Condition for timeout based testing
func WaitForCondition(d time.Duration, testFn func() bool) bool {
	if d < time.Second {
		panic("WaitForCondition: test duration to small")
	}

	test := time.Tick(500 * time.Millisecond)
	timeout := time.Tick(d)

	check := make(chan struct{}, 1)
	done := make(chan struct{}, 1)
	defer close(check)
	defer close(done)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-test:
				if testFn() {
					check <- struct{}{}
					return
				}
			}
		}
	}()

	for {
		select {
		case <-check:
			return true
		case <-timeout:
			done <- struct{}{}
			return false
		}
	}
}

func RandomizeBytes(in []byte) []byte {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(time.Now().Unix()))
	return append(in, bs...)
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixedRounded(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(num*output) / output
}

func JoinHome(p string) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return path.Join(home, p), nil
}
