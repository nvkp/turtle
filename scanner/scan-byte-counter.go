package scanner

import "bufio"

type scanByteCounter struct {
	BytesRead int
}

func (s *scanByteCounter) splitFunc() bufio.SplitFunc {
	return func(data []byte, atEOF bool) (int, []byte, error) {
		adv, tok, err := splitTurtle(data, atEOF)
		s.BytesRead += adv
		return adv, tok, err
	}
}
