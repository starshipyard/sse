package sse

import (
	"encoding/json"
	"fmt"
	"io"

	"reflect"
	"strings"
)

type Event struct {
	Event string
	Id    string
	Retry uint
	Data  interface{}
}

func Encode(w io.Writer, event Event) error {
	writeId(w, event.Id)
	writeEvent(w, event.Event)
	writeRetry(w, event.Retry)
	return writeData(w, event.Data)
}

func writeId(w io.Writer, id string) {
	if len(id) > 0 {
		w.Write([]byte("id: "))
		w.Write([]byte(escape(id)))
		w.Write([]byte("\n"))
	}
}

func writeEvent(w io.Writer, event string) {
	if len(event) > 0 {
		w.Write([]byte("event: "))
		w.Write([]byte(escape(event)))
		w.Write([]byte("\n"))
	}
}

func writeRetry(w io.Writer, retry uint) {
	if retry > 0 {
		fmt.Fprintf(w, "retry: %d\n", retry)
	}
}

func writeData(w io.Writer, data interface{}) error {
	w.Write([]byte("data: "))
	switch typeOfData(data) {
	case reflect.Struct, reflect.Slice, reflect.Map:
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			return err
		}
		w.Write([]byte("\n"))
	default:
		text := fmt.Sprint(data)
		fmt.Fprint(w, escape(text), "\n\n")
	}
	return nil
}

func typeOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func escape(str string) string {
	str = strings.Replace(str, "\n", "\\n", -1)
	str = strings.Replace(str, "\r", "\\r", -1)
	return str
}
