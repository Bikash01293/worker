package work

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestCollector(t *testing.T) {
	jsonstr := WorkData{
		Name:     "Harry",
		City:     "Goa",
		WorkDays: 18,
		Salary:   205,
		Delay:    "3s",
	}

	jbyte, _ := json.Marshal(jsonstr)
	req, err := http.NewRequest("POST", "/collector", bytes.NewBuffer(jbyte))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	p := NewPrinter()
	const maxNumWorkers = 4
	d := NewDispatcher(p, maxNumWorkers, maxNumWorkers+1)
	handler := http.HandlerFunc(d.Collector)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	str := strings.TrimSpace(rr.Body.String())

	if str != string(jbyte) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			str, string(jbyte))
	}
}

func TestNewDispatcher(t *testing.T) {
	t.Run("testing dispatcher", func(t *testing.T) {
		const maxNumWorkers = 4
		p := NewPrinter()
		got := NewDispatcher(p, maxNumWorkers, maxNumWorkers+1)
		want := Dispatcher{
			sem:       make(chan struct{}, maxNumWorkers),
			jobBuffer: make(chan *WorkRequest, maxNumWorkers+1),
			worker:    p,
			wg:        sync.WaitGroup{},
		}
		// assert.Equal(t, got, want)

		if *got != want {
			t.Errorf("returned unexpected body: got %v want %v", *got, want)
		}

	})
}

// func TestStart(t *testing.T) {
// 	ctx, _ := context.WithCancel(context.Background())
// 	const maxWorkers = 4
// 	p := NewPrinter()
// 	d := NewDispatcher(p, maxWorkers, maxWorkers+1)
// 	got := d.Start(ctx)
// 	fmt.Println(got)
// }

// func TestLoop(t *testing.T) {
// 	ctx, _ := context.WithCancel(context.Background())
// 	const maxWorkers = 4
// 	p := NewPrinter()
// 	d := NewDispatcher(p, maxWorkers, maxWorkers+1)
// 	got := d.loop(ctx)
// }
