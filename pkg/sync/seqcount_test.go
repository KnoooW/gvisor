// Copyright 2018 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sync

import (
	"reflect"
	"testing"
	"time"
)

func TestSeqCountWriteUncontended(t *testing.T) {
	var seq SeqCount
	seq.BeginWrite()
	seq.EndWrite()
}

func TestSeqCountReadUncontended(t *testing.T) {
	var seq SeqCount
	epoch := seq.BeginRead()
	if !seq.ReadOk(epoch) {
		t.Errorf("ReadOk: got false, wanted true")
	}
}

func TestSeqCountBeginReadAfterWrite(t *testing.T) {
	var seq SeqCount
	var data int32
	const want = 1
	seq.BeginWrite()
	data = want
	seq.EndWrite()
	epoch := seq.BeginRead()
	if data != want {
		t.Errorf("Reader: got %v, wanted %v", data, want)
	}
	if !seq.ReadOk(epoch) {
		t.Errorf("ReadOk: got false, wanted true")
	}
}

func TestSeqCountBeginReadDuringWrite(t *testing.T) {
	var seq SeqCount
	var data int
	const want = 1
	seq.BeginWrite()
	go func() {
		time.Sleep(time.Second)
		data = want
		seq.EndWrite()
	}()
	epoch := seq.BeginRead()
	if data != want {
		t.Errorf("Reader: got %v, wanted %v", data, want)
	}
	if !seq.ReadOk(epoch) {
		t.Errorf("ReadOk: got false, wanted true")
	}
}

func TestSeqCountReadOkAfterWrite(t *testing.T) {
	var seq SeqCount
	epoch := seq.BeginRead()
	seq.BeginWrite()
	seq.EndWrite()
	if seq.ReadOk(epoch) {
		t.Errorf("ReadOk: got true, wanted false")
	}
}

func TestSeqCountReadOkDuringWrite(t *testing.T) {
	var seq SeqCount
	epoch := seq.BeginRead()
	seq.BeginWrite()
	if seq.ReadOk(epoch) {
		t.Errorf("ReadOk: got true, wanted false")
	}
	seq.EndWrite()
}

func BenchmarkSeqCountWriteUncontended(b *testing.B) {
	var seq SeqCount
	for i := 0; i < b.N; i++ {
		seq.BeginWrite()
		seq.EndWrite()
	}
}

func BenchmarkSeqCountReadUncontended(b *testing.B) {
	var seq SeqCount
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			epoch := seq.BeginRead()
			if !seq.ReadOk(epoch) {
				b.Fatalf("ReadOk: got false, wanted true")
			}
		}
	})
}

func TestPointersInType(t *testing.T) {
	for _, test := range []struct {
		name string // used for both test and value name
		val  interface{}
		ptrs []string
	}{
		{
			name: "EmptyStruct",
			val:  struct{}{},
		},
		{
			name: "Int",
			val:  int(0),
		},
		{
			name: "MixedStruct",
			val: struct {
				b             bool
				I             int
				ExportedPtr   *struct{}
				unexportedPtr *struct{}
				arr           [2]int
				ptrArr        [2]*int
				nestedStruct  struct {
					nestedNonptr int
					nestedPtr    *int
				}
				structArr [1]struct {
					nonptr int
					ptr    *int
				}
			}{},
			ptrs: []string{
				"MixedStruct.ExportedPtr",
				"MixedStruct.unexportedPtr",
				"MixedStruct.ptrArr[]",
				"MixedStruct.nestedStruct.nestedPtr",
				"MixedStruct.structArr[].ptr",
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			typ := reflect.TypeOf(test.val)
			ptrs := PointersInType(typ, test.name)
			t.Logf("Found pointers: %v", ptrs)
			if (len(ptrs) != 0 || len(test.ptrs) != 0) && !reflect.DeepEqual(ptrs, test.ptrs) {
				t.Errorf("Got %v, wanted %v", ptrs, test.ptrs)
			}
		})
	}
}
