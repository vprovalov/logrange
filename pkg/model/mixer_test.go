// Copyright 2018-2019 The logrange Authors
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

package model

import (
	"io"
	"reflect"
	"testing"
)

func tstGetFirst(ev1, ev2 LogEvent) bool {
	return true
}

func TestMixer(t *testing.T) {
	testMix(t, tstGetFirst,
		[]LogEvent{{1, []byte("1"), ""}, {2, []byte("2"), ""}},
		[]LogEvent{{3, []byte("1"), ""}, {4, []byte("2"), ""}},
		[]LogEvent{{1, []byte("1"), ""}, {2, []byte("2"), ""}, {3, []byte("1"), ""}, {4, []byte("2"), ""}})

	testMix(t, GetEarliest,
		[]LogEvent{{1, []byte("1"), ""}, {3, []byte("2"), ""}},
		[]LogEvent{{2, []byte("1"), ""}, {4, []byte("2"), ""}},
		[]LogEvent{{1, []byte("1"), ""}, {2, []byte("1"), ""}, {3, []byte([]byte("2")), ""}, {4, []byte("2"), ""}})

	testMix(t, GetEarliest,
		[]LogEvent{},
		[]LogEvent{{2, []byte("1"), ""}, {4, []byte("2"), ""}},
		[]LogEvent{{2, []byte("1"), ""}, {4, []byte("2"), ""}})

	testMix(t, GetEarliest,
		[]LogEvent{{1, []byte("1"), ""}, {3, []byte("2"), ""}},
		[]LogEvent{},
		[]LogEvent{{1, []byte("1"), ""}, {3, []byte("2"), ""}})

	testMix(t, GetEarliest,
		[]LogEvent{},
		[]LogEvent{},
		[]LogEvent{})
}

func testMix(t *testing.T, sf SelectF, i1, i2, res []LogEvent) {
	m := &Mixer{}
	it1 := &LogEventIterator{}
	it1.Wrap("t1", NewTestLogEventsWrapper(i1))
	it2 := &LogEventIterator{}
	it2.Wrap("t2", NewTestLogEventsWrapper(i2))
	m.Init(sf, it1, it2)
	testIt(t, m, res)
}

func testIt(t *testing.T, it Iterator, res []LogEvent) {
	idx := 0
	for {
		le, t1, err := it.Get(nil)
		if err == io.EOF {
			break
		}
		le2, t2, err := it.Get(nil)
		if !reflect.DeepEqual(le, le2) || t1 != t2 {
			t.Fatal("expecting le=", le, " to be equal to ", le2)
		}
		if !reflect.DeepEqual(le, res[idx]) {
			t.Fatal("expected ", res[idx], ", but got ", le)
		}
		idx++
		it.Next(nil)
	}
	if idx != len(res) {
		t.Fatal("Must be ", len(res), ", but idx=", idx, it)
	}

	it.SetBackward(true)
	for {
		idx--
		le, _, err := it.Get(nil)
		if err == io.EOF {
			break
		}
		if !reflect.DeepEqual(le, res[idx]) {
			t.Fatal("expected ", res[idx], ", but got ", le)
		}
		it.Next(nil)
	}

	if idx != -1 {
		t.Fatal("expecting idx==-1, but it is ", idx)
	}
}
