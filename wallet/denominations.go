package wallet

const (
	// All units stem from Groth.
	Groth = 1
	Beam = 100000000
)

func GrothToBeam(groth uint64) float64 {
	return float64(groth) / float64(Beam)
}

func BeamToGroth(beam float64) uint64 {
	return uint64(beam * Beam)
}