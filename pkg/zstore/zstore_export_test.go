package zstore

// MasterKeyForTest returns the master key slice for test assertions.
func (s *Store) MasterKeyForTest() []byte { return s.masterKey }

// SubKeysForTest returns the sub-key slices for test assertions.
func (s *Store) SubKeysForTest() [][]byte { return s.subKeys }
