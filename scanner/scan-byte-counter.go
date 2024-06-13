package scanner

import "bufio"

type ScanByteCounter struct {
	BytesRead int
}

func (s *ScanByteCounter) SplitFunc() bufio.SplitFunc {
	return func(data []byte, atEOF bool) (int, []byte, error) {
		adv, tok, err := splitTurtle(data, atEOF)
		s.BytesRead += adv
		return adv, tok, err
	}
}
