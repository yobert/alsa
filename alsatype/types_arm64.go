package alsatype

type (
	Uframes uint64 // snd_pcm_uframes_t
	Sframes int64  // snd_pcm_sframes_t
)

type Timespec struct {
	Sec  int
	Nsec int
}

type SwParams struct {
	TstampMode int32
	PeriodStep uint32
	SleepMin   uint32

	AvailMin         Uframes
	XferAlign        Uframes
	StartThreshold   Uframes
	StopThreshold    Uframes
	SilenceThreshold Uframes
	SilenceSize      Uframes
	Boundary         Uframes

	Proto      PVersion
	TstampType uint32
	Reserved   [56]byte
}
type PVersion uint32
