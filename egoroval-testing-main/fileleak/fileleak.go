//go:build !solution

package fileleak

import (
	"fmt"
	"log"
	"os"
	"reflect"
)

var path string = "/proc/self/fd"

type testingT interface {
	Errorf(msg string, args ...interface{})
	Cleanup(func())
}

func VerifyNone(t testingT) {
	start, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	sm := make(map[string]int)
	for _, f := range start {
		d, err := os.Readlink(path + "/" + f.Name())
		if err != nil {
			//log.Fatal(err)
			fmt.Println("")
		}
		_, ok := sm[d]
		if !ok {
			sm[d] = 1
			continue
		}
		sm[d]++
	}
	t.Cleanup(func() {
		end, err := os.ReadDir(path)
		if err != nil {
			log.Fatal(err)
			return
		}
		em := make(map[string]int)
		for _, f := range end {
			d, err := os.Readlink(path + "/" + f.Name())
			if err != nil {
				//log.Fatal(err)
				fmt.Println("")
			}
			_, ok := em[d]
			if !ok {
				em[d] = 1
				continue
			}
			em[d]++
		}
		if !reflect.DeepEqual(sm, em) {
			t.Errorf("Fileleak detected!")
		}
	})

}
