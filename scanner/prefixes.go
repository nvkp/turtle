package scanner

// Prefixes returns all prefixes from the so far scanned content.
func (s *Scanner) Prefixes() map[string]string {
	return s.prefixes
}

// Base returns empty string in case there was no base specified
// in the so far scanned content and returns the string of the
// base prefix in case there was.
func (s *Scanner) Base() string {
	return s.base
}
