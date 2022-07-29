package snowflake

type BitsAllocator struct {
	timestampLen int64
	workidLen    int64
	sequenceLen  int64
}

type Uidgenerator struct {
	timestamp int64
	workid    int64
	sequence  int64
}

func (b *BitsAllocator) Allocator(timestamp uint64, workid uint64, sequence uint64) uint64 {
	timestampShift := b.workidLen + b.sequenceLen
	workidShift := b.sequenceLen
	return uint64(timestamp<<uint64(timestampShift) | (workid << uint64(workidShift)) | sequence)

}

func (u *Uidgenerator) getUid(int64) {

}
