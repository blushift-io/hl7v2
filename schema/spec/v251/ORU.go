package v251

type ORU_R01 struct {
	MSH           any
	PatientResult []ORU_R01_PatientResult

	OBR any
	OBX any
}

type ORU_R01_PatientResult struct {
	Patient *ORU_RO1_PatientResult_Patient
}

type ORU_RO1_PatientResult_Patient struct {
	PID   any
	PD1   any
	NTE   any
	NK1   any
	Visit *ORU_RO1_PatientResult_Patient_Visit
}

type ORU_RO1_PatientResult_Patient_Visit struct {
	PV1 any
	PV2 any
}

type ORU_R01_PatientResult_OrderObservation struct {
	ORC any
	OBR any
}

type ORU_R01_PatientResult_OrderObservation_Observation struct {
	OBX []any
	NTE []any
}
