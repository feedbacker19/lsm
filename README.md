# lsm

A small LSM-tree key-value store in Go.

## Plan
1. Memtable (sorted, in-memory)   <- memtable.go
2. Write-ahead log                <- wal.go
3. SSTable writer/reader          <- sstable.go
4. Compaction                     <- compaction.go
5. DB API: Put/Get/Delete         <- db.go
