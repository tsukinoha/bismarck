package bismarck

import (
	"testing"
	"time"
)

type (
	subCmd0 struct {
		Desc    string
		field1  string `bismarck:""`
		field2  string `bismarck:"short=2"`
		field3  string `bismarck:"long=field3"`
		field4  string `bismarck:"desc=desc4"`
		field5  string `bismarck:"required"`
		field7  string `bismarck:"format=2006/01/02 15:04:05"`
		field8  string `bismarck:"handler=Handler"`
		field9  string `bismarck:"short=9,long=field9,desc= desc9, test ,required,format=Mon, 02 Jan 2006 15:04:05 MST,handler=Handler"`
		field10 string `bismarck:"short=a,long=field10,desc= desc10, test ,handler=Handler"`
	}
	subCmd1 struct {
		string_0  string    `bismarck:"short=s,long=string,desc=string desc"`
		bool_0    bool      `bismarck:"short=b,long=bool,desc=bool desc"`
		int_0     int       `bismarck:"short=i,long=int,desc=int desc"`
		int8_0    int8      `bismarck:"short=I,long=int8,desc=int8 desc"`
		int16_0   int16     `bismarck:"short=J,long=int16,desc=int16 desc"`
		int32_0   int32     `bismarck:"short=K,long=int32,desc=int32 desc"`
		int64_0   int64     `bismarck:"short=L,long=int64,desc=int64 desc"`
		uint_0    uint      `bismarck:"short=M,long=uint,desc=uint desc"`
		uint8_0   uint8     `bismarck:"short=U,long=uint8,desc=uint8 desc"`
		uint16_0  uint16    `bismarck:"short=V,long=uint16,desc=uint16 desc"`
		uint32_0  uint32    `bismarck:"short=W,long=uint32,desc=uint32 desc"`
		uint64_0  uint64    `bismarck:"short=X,long=uint64,desc=uint64 desc"`
		float32_0 float32   `bismarck:"short=f,long=float32,desc=float32 desc"`
		float64_0 float64   `bismarck:"short=F,long=float64,desc=float64 desc"`
		time_0    time.Time `bismarck:"short=t,long=time,desc=time desc"`
		unknown_0 complex64 `bismarck:"short=c,long=complex64,desc=complex64 desc"`
	}
	subCmd2 struct {
	}
)

func (cmd subCmd0) Init()           {}
func (cmd subCmd0) Validate() error { return nil }
func (cmd subCmd0) Run(args []string) error {
	return nil
}

func (cmd subCmd1) Init()           {}
func (cmd subCmd1) Validate() error { return nil }
func (cmd subCmd1) Run(args []string) error {
	return nil
}

func (cmd subCmd2) Init()           {}
func (cmd subCmd2) Validate() error { return nil }
func (cmd subCmd2) Run(args []string) error {
	return nil
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
	cleanup()
}

func setup() {
	println("* Setting up")
}

func cleanup() {
	println("* Cleaning up")
}
