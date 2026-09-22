package profiling

import (
	"bytes"
	"runtime/pprof"
	"syscall/js"
)

var cpuBuf bytes.Buffer
var profiling bool

func StartCPUProfile() {
	if profiling {
		return
	}
	cpuBuf.Reset()
	pprof.StartCPUProfile(&cpuBuf)
	profiling = true
}

func StopCPUProfileAndDownload() {
	if !profiling {
		return
	}
	pprof.StopCPUProfile()
	profiling = false
	downloadBytes(cpuBuf.Bytes(), "cpu.prof")
}

func DumpHeapProfile() {
	var buf bytes.Buffer
	pprof.WriteHeapProfile(&buf)
	downloadBytes(buf.Bytes(), "heap.prof")
}

func downloadBytes(data []byte, filename string) {
	jsData := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(jsData, data)

	blob := js.Global().Get("Blob").New(
		js.Global().Get("Array").New(jsData),
		map[string]any{"type": "application/octet-stream"},
	)
	url := js.Global().Get("URL").Call("createObjectURL", blob)

	a := js.Global().Get("document").Call("createElement", "a")
	a.Set("href", url)
	a.Set("download", filename)
	a.Call("click")

	js.Global().Get("URL").Call("revokeObjectURL", url)
}
