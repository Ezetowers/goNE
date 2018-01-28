package network

import (
	"errors"
	"fmt"
)

type Mac struct {
	bytes []byte
}

const MacSize = 6

func NewMac(data []byte) *Mac {
	mac := &Mac{
		make([]byte, MacSize),
	}

	if len(data) < MacSize {
		Log.Errorf("[MAC] Invalid data received. Actual len: %v. Expected: %v", len(data), MacSize)
		return nil
	}

	copy(mac.bytes, data)
	return mac
}

func (self *Mac) String() string {
	return fmt.Sprintf("%.2x:%.2x:%.2x:%.2x:%.2x:%.2x",
		self.bytes[0],
		self.bytes[1],
		self.bytes[2],
		self.bytes[3],
		self.bytes[4],
		self.bytes[5])
}

func (self *Mac) Increase() error {
	overflow := true
	index := 5

	for overflow == true && index >= 0 {
		if self.bytes[index] == 255 {
			self.bytes[index] = 0
			overflow = true
		} else {
			self.bytes[index]++
			overflow = false
		}
		index--
	}

	if overflow != false {
		// The Mac could not be incremented
		self.bytes = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
		return errors.New("Mac cannot be increased. Max Mac was received")

	}

	return nil
}

/**
 * @brief Compare the IPMasks.
 *
 * @param      self Self
 * @param      rhs  Mac to compare
 *
 * @return     0  if self == rhs
 * 			   1  if self > rhs
 * 			   -1 if self < rhs
 */
func (self *Mac) Compare(rhs *Mac) int8 {
	for index := 0; index < MacSize; index++ {
		if self.bytes[index] != rhs.bytes[index] {
			if self.bytes[index] > rhs.bytes[index] {
				return 1
			} else {
				return -1
			}
		}
	}

	return 0
}
