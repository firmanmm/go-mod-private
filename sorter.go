package gomodprivate

type stringSort []string

// Len is the number of elements in the collection.
func (s stringSort) Len() int {
	return len(s)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (s stringSort) Less(i int, j int) bool {
	return s[i] < s[j]
}

// Swap swaps the elements with indexes i and j.
func (s stringSort) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}

type sshCredentialSort []*SshCredential

// Len is the number of elements in the collection.
func (s sshCredentialSort) Len() int {
	return len(s)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (s sshCredentialSort) Less(i int, j int) bool {
	return s[i].Host < s[j].Host
}

// Swap swaps the elements with indexes i and j.
func (s sshCredentialSort) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}
