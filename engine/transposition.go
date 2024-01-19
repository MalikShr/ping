package engine

const (
	DefaultTableSize = 64
	SearchEntrySize  = 12
)

type TranspositionTable struct {
	Entries []SearchEntry
	Size    uint64
}

type SearchEntry struct {
	Hash uint64
	Best Move
}

func (tt *TranspositionTable) InitTransTable(sizeMB uint64) {
	size := (sizeMB * 1024 * 1024) / SearchEntrySize

	tt.Entries = make([]SearchEntry, size)
	tt.Size = size
}

func (tt *TranspositionTable) Store(hash uint64, move Move) {
	index := hash % tt.Size

	tt.Entries[index] = SearchEntry{
		Hash: hash,
		Best: move,
	}
}

func (tt *TranspositionTable) Probe(hash uint64) SearchEntry {
	index := hash % tt.Size

	return tt.Entries[index]
}

func (tt *TranspositionTable) Clear() {
	for i := uint64(0); i < tt.Size; i++ {
		tt.Entries[i] = *new(SearchEntry)
	}
}
