package network

import (
	"errors"
	"fmt"
	"net"
)

type IPAddress struct {
	ip net.IP
}

const IPV4_SIZE = 4
const IPV6_SIZE = 16

func NewIPAddress(data []byte) *IPAddress {
	ipSize := 0

	if len(data) == IPV4_SIZE {
		ipSize = IPV4_SIZE
	} else if len(data) == IPV6_SIZE {
		ipSize = IPV6_SIZE
	} else {
		Log.Errorf("[MAC] Invalid IP size received. Size received: %v. Valid sizes: (4,6)", len(data))
		return nil
	}

	ipAddr := &IPAddress{
		data[0:ipSize],
	}

	return ipAddr
}

/**
 * @return     Return the IP as a stream of bytes
 */
func (self *IPAddress) Raw() *net.IP {
	return &self.ip
}

/**
 * @return     String representing the IPAddress
 */
func (self *IPAddress) String() string {
	return fmt.Sprintf("%v", self.ip)
}

/**
 * @brief      Increase in one the IP.
 *
 * @return     An error is returned if the maximum length of the
 * 			   IP has been reached
 */
func (self *IPAddress) Increase() error {
	overflow := true
	index := len(self.ip) - 1

	for overflow == true && index >= 0 {
		if self.ip[index] == 255 {
			self.ip[index] = 0
			overflow = true
		} else {
			self.ip[index]++
			overflow = false
		}
		index--
	}

	if overflow != false {
		// The Mac could not be incremented
		for i := 0; i < len(self.ip); i++ {
			self.ip[i] = 0xff
		}
		return errors.New("IPAddress cannot be increased. Max IP was received")

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
 *             1  if self > rhs
 *             -1 if self < rhs
 */
func (self *IPAddress) Compare(rhs *IPAddress) (int8, error) {
	if len(self.ip) != len(rhs.ip) {
		return -2, errors.New("IPAdddresses doesn't have the same format. Cannot be compared.")
	}

	for index := 0; index < len(self.ip); index++ {
		if self.ip[index] != rhs.ip[index] {
			if self.ip[index] > rhs.ip[index] {
				return 1, nil
			} else {
				return -1, nil
			}
		}
	}

	return 0, nil
}
