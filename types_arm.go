package alsa

type swParams struct {
	TstampMode int32
	PeriodStep uint32
	SleepMin   uint32

	AvailMin         misc.Uframes
	XferAlign        misc.Uframes
	StartThreshold   misc.Uframes
	StopThreshold    misc.Uframes
	SilenceThreshold misc.Uframes
	SilenceSize      misc.Uframes
	Boundary         misc.Uframes

	Proto         pVersion
	TstampType    uint32
	Reserved      [56]byte
}
