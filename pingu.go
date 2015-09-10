package pingu

import (
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/hybridgroup/gobot/platforms/firmata/client"
	"github.com/tarm/goserial"
)

type request struct{}

// Tux is representing the waving pinguin device
type Tux struct {
	client *client.Client
	wave   chan request
	lock   *sync.Mutex
	port   string
	conn   io.ReadWriteCloser
	openSP func(port string) (io.ReadWriteCloser, error)
}

// NewTux creates a new instance of Tux and creates the connection to the arduino.
func NewTux(args ...interface{}) *Tux {
	t := &Tux{
		client: client.New(),
		port:   "",
		wave:   make(chan request),
		lock:   &sync.Mutex{},
		conn:   nil,
		openSP: func(port string) (io.ReadWriteCloser, error) {
			return serial.OpenPort(&serial.Config{Name: port, Baud: 57600})
		},
	}
	for _, arg := range args {
		switch arg.(type) {
		case string:
			t.port = arg.(string)
		case io.ReadWriteCloser:
			t.conn = arg.(io.ReadWriteCloser)
		}
	}

	return t
}

// Connect starts a connection to the board.
func (t *Tux) Connect() (errs []error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.conn == nil {
		sp, err := t.openSP(t.Port())
		if err != nil {
			return []error{err}
		}
		t.conn = sp
	}
	if err := t.client.Connect(t.conn); err != nil {
		return []error{err}
	}
	return
}

// Disconnect closes the io connection to the board
func (t *Tux) Disconnect() (err error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.client != nil {
		return t.client.Disconnect()
	}
	return nil
}

// Port returns the  FirmataAdaptors port
func (t *Tux) Port() string { return t.port }

// ServoWrite writes the 0-180 degree angle to the specified pin.
func (t *Tux) ServoWrite(pin string, angle byte) (err error) {
	p, err := strconv.Atoi(pin)
	if err != nil {
		return err
	}
	t.lock.Lock()
	defer t.lock.Unlock()
	err = t.client.AnalogWrite(p, int(angle))
	return
}

// AttachServo set pin mode to servo
func (t *Tux) AttachServo(pin string) (err error) {
	p, err := strconv.Atoi(pin)
	if err != nil {
		return err
	}
	err = t.client.SetPinMode(p, client.Servo)
	return
}

// DettachServo switch off servo
func (t *Tux) DettachServo(pin string) (err error) {
	p, err := strconv.Atoi(pin)
	if err != nil {
		return err
	}
	err = t.client.SetPinMode(p, client.Output)
	err = t.client.DigitalWrite(p, 0)
	return
}

// Wave penguin, wave!
func (t *Tux) Wave() {
	t.wave <- request{}
}

// Run the waving for a single arm.
func (t *Tux) Run(pin string) {
	t.Connect()
	defer t.Disconnect()
	for range t.wave {
		t.doTheWave(pin)
	}
}

func (t *Tux) doTheWave(pin string) {
	t.AttachServo(pin)
	t.ServoWrite(pin, 180)
	time.Sleep(400 * time.Millisecond)
	t.ServoWrite(pin, 120)
	time.Sleep(300 * time.Millisecond)
	t.ServoWrite(pin, 180)
	time.Sleep(300 * time.Millisecond)
	t.ServoWrite(pin, 120)
	time.Sleep(300 * time.Millisecond)
	t.ServoWrite(pin, 90)
	time.Sleep(200 * time.Millisecond)
	t.DettachServo(pin)
}
