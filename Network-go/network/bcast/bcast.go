package bcast

import (
	"TTK4145---project/Network-go/network/conn"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"time"
)

const bufSize = 1024

// Encodes received values from `chans` into type-tagged JSON, then broadcasts
// it on `port`
func Transmitter(ctx context.Context, port int, chans ...interface{}) {
	checkArgs(chans...)
	typeNames := make([]string, len(chans))
	selectCases := make([]reflect.SelectCase, len(typeNames))
	for i, ch := range chans {
		selectCases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
		typeNames[i] = reflect.TypeOf(ch).Elem().String()
	}

	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))
	for {
		select {
		case <-ctx.Done():
			fmt.Println("[Bcast] Transmitter shutting down")
			return
		default:
			chosen, value, ok := reflect.Select(selectCases)
			if !ok {
				continue
			}
			jsonstr, _ := json.Marshal(value.Interface())
			ttj, _ := json.Marshal(typeTaggedJSON{
				TypeId: typeNames[chosen],
				JSON:   jsonstr,
			})
			if len(ttj) > bufSize {
				panic(fmt.Sprintf("Message too long (%d bytes)", len(ttj)))
			}
			conn.WriteTo(ttj, addr)
		}
	}
}

// Matches type-tagged JSON received on `port` to element types of `chans`, then
// sends the decoded value on the corresponding channel
func Receiver(ctx context.Context, port int, chans ...interface{}) {
	checkArgs(chans...)
	chansMap := make(map[string]interface{})
	for _, ch := range chans {
		chansMap[reflect.TypeOf(ch).Elem().String()] = ch
	}

	var buf [bufSize]byte
	conn := conn.DialBroadcastUDP(port)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("[Bcast] Receiver shutting down")
			return
		default:
			conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
			n, _, err := conn.ReadFrom(buf[0:])
			if err != nil {
				continue
			}

			var ttj typeTaggedJSON
			json.Unmarshal(buf[0:n], &ttj)
			ch, ok := chansMap[ttj.TypeId]
			if !ok {
				continue
			}
			v := reflect.New(reflect.TypeOf(ch).Elem())
			json.Unmarshal(ttj.JSON, v.Interface())
			reflect.Select([]reflect.SelectCase{{
				Dir:  reflect.SelectSend,
				Chan: reflect.ValueOf(ch),
				Send: reflect.Indirect(v),
			}})
		}
	}
}

type typeTaggedJSON struct {
	TypeId string
	JSON   []byte
}

// Checks that args to Tx'er/Rx'er are valid:
//
//	All args must be channels
//	Element types of channels must be encodable with JSON
//	No element types are repeated
//
// Implementation note:
//   - Why there is no `isMarshalable()` function in encoding/json is a mystery,
//     so the tests on element type are hand-copied from `encoding/json/encode.go`
func checkArgs(chans ...interface{}) {
	n := 0
	for range chans {
		n++
	}
	elemTypes := make([]reflect.Type, n)

	for i, ch := range chans {
		// Must be a channel
		if reflect.ValueOf(ch).Kind() != reflect.Chan {
			panic(fmt.Sprintf(
				"Argument must be a channel, got '%s' instead (arg# %d)",
				reflect.TypeOf(ch).String(), i+1))
		}

		elemType := reflect.TypeOf(ch).Elem()

		// Element type must not be repeated
		for j, e := range elemTypes {
			if e == elemType {
				panic(fmt.Sprintf(
					"All channels must have mutually different element types, arg# %d and arg# %d both have element type '%s'",
					j+1, i+1, e.String()))
			}
		}
		elemTypes[i] = elemType

		// Element type must be encodable with JSON
		checkTypeRecursive(elemType, []int{i + 1})

	}
}

func checkTypeRecursive(val reflect.Type, offsets []int) {
	switch val.Kind() {
	case reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func, reflect.UnsafePointer:
		panic(fmt.Sprintf(
			"Channel element type must be supported by JSON, got '%s' instead (nested arg# %v)",
			val.String(), offsets))
	case reflect.Map:
		if val.Key().Kind() != reflect.String {
			panic(fmt.Sprintf(
				"Channel element type must be supported by JSON, got '%s' instead (map keys must be 'string') (nested arg# %v)",
				val.String(), offsets))
		}
		checkTypeRecursive(val.Elem(), offsets)
	case reflect.Array, reflect.Ptr, reflect.Slice:
		checkTypeRecursive(val.Elem(), offsets)
	case reflect.Struct:
		for idx := 0; idx < val.NumField(); idx++ {
			checkTypeRecursive(val.Field(idx).Type, append(offsets, idx+1))
		}
	}
}
