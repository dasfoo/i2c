package i2c

// Workaround for i2c-dev divergence not applying to cross-compilation subdir.

// #include "/usr/include/linux/i2c-dev.h"
import "C"

import (
	"fmt"
	"log"
	"os"
	"sync"
	"syscall"
)

const i2cLogFormat = "i2c %4s: [%#.2x:%#.2x] %v (err: %v)\n"

// Logger is a function to call whenever anything is read/written to Bus
type Logger func(string, ...interface{})

// Bus is an interface of I2C bus accessor
type Bus interface {
	ReadByteFromReg(byte, byte) (byte, error)
	ReadWordFromReg(byte, byte) (uint16, error)
	ReadSliceFromReg(byte, byte, []byte) (int, error)
	WriteSliceToReg(byte, byte, []byte) (int, error)
	WriteByteToReg(byte, byte, byte) error
	SetLogger(Logger)
	Close() error
}

type bus struct {
	file          *os.File
	opLock        sync.Mutex
	remoteAddress byte
	logger        Logger
}

// NewBus opens a Linux i2c bus file
func NewBus(id byte) (Bus, error) {
	file, err := os.OpenFile(
		fmt.Sprintf("/dev/i2c-%d", id),
		os.O_RDWR,
		os.ModeExclusive)
	if err != nil {
		return nil, err
	}
	return &bus{
		file:   file,
		logger: log.Printf,
	}, nil
}

func (b *bus) setRemoteAddress(addr byte) error {
	if addr != b.remoteAddress {
		if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(),
			C.I2C_SLAVE, uintptr(addr)); errno != 0 {
			return syscall.Errno(errno)
		}
		b.remoteAddress = addr
	}
	return nil
}

// ReadByteFromReg reads 1 byte from a register of a slave device
func (b *bus) ReadByteFromReg(addr, reg byte) (byte, error) {
	b.opLock.Lock()
	defer b.opLock.Unlock()
	if err := b.setRemoteAddress(addr); err != nil {
		return 0, err
	}
	value, err := C.i2c_smbus_read_byte_data(C.int(b.file.Fd()), C.__u8(reg))
	b.logger(i2cLogFormat, "recv", addr, reg, value, err)
	return byte(value), err
}

// ReadWordFromReg reads 2 bytes from a register of a slave device
func (b *bus) ReadWordFromReg(addr, reg byte) (uint16, error) {
	b.opLock.Lock()
	defer b.opLock.Unlock()
	if err := b.setRemoteAddress(addr); err != nil {
		return 0, err
	}
	value, err := C.i2c_smbus_read_word_data(C.int(b.file.Fd()), C.__u8(reg))
	b.logger(i2cLogFormat, "recv", addr, reg, value, err)
	leValue := uint16(value) // big endian yet
	return ((leValue >> 8) | (leValue << 8)), err
}

// ReadSliceFromReg reads an undefined number of bytes
func (b *bus) ReadSliceFromReg(addr, reg byte, value []byte) (int, error) {
	b.opLock.Lock()
	defer b.opLock.Unlock()
	if err := b.setRemoteAddress(addr); err != nil {
		return 0, err
	}
	if len(value) > C.I2C_SMBUS_BLOCK_MAX {
		// TODO: issue a warning or something.
	}
	size, err := C.i2c_smbus_read_i2c_block_data(C.int(b.file.Fd()),
		C.__u8(reg), C.__u8(len(value)), (*C.__u8)(&value[0]))
	b.logger(i2cLogFormat, "recv", addr, reg, value, err)
	return int(size), err
}

// WriteSliceToReg writes a defined number of bytes
func (b *bus) WriteSliceToReg(addr, reg byte, value []byte) (int, error) {
	b.opLock.Lock()
	defer b.opLock.Unlock()
	if err := b.setRemoteAddress(addr); err != nil {
		return 0, err
	}
	if len(value) > C.I2C_SMBUS_BLOCK_MAX {
		// TODO: issue a warning or something.
	}
	size, err := C.i2c_smbus_write_i2c_block_data(C.int(b.file.Fd()),
		C.__u8(reg), C.__u8(len(value)), (*C.__u8)(&value[0]))
	b.logger(i2cLogFormat, "send", addr, reg, value, err)
	return int(size), err
}

// WriteByteToReg writes 1 byte to a register of a slave device
func (b *bus) WriteByteToReg(addr, reg, value byte) error {
	b.opLock.Lock()
	defer b.opLock.Unlock()
	if err := b.setRemoteAddress(addr); err != nil {
		return err
	}
	_, err := C.i2c_smbus_write_byte_data(C.int(b.file.Fd()), C.__u8(reg),
		C.__u8(value))
	b.logger(i2cLogFormat, "send", addr, reg, value, err)
	return err
}

// Close frees any resources allocated for the bus
func (b *bus) Close() error {
	b.opLock.Lock()
	defer b.opLock.Unlock()
	return b.file.Close()
}

// SetLogger changes the bus logging function
func (b *bus) SetLogger(l Logger) {
	b.logger = l
}
