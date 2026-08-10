package v251

// ORU_R01 represents an HL7 v2.5.1 Unsolicited Observation Message structure.
type ORU_R01 struct {
	MSH           any
	PatientResult []ORU_R01_PatientResult

	OBR any
	OBX any
}

// ORU_R01_PatientResult contains patient result information for an ORU_R01 message.
type ORU_R01_PatientResult struct {
	Patient *ORU_RO1_PatientResult_Patient
}

// ORU_RO1_PatientResult_Patient holds patient demographic and visit details.
type ORU_RO1_PatientResult_Patient struct {
	PID   any
	PD1   any
	NTE   any
	NK1   any
	Visit *ORU_RO1_PatientResult_Patient_Visit
}

// ORU_RO1_PatientResult_Patient_Visit contains patient visit segments.
type ORU_RO1_PatientResult_Patient_Visit struct {
	PV1 any
	PV2 any
}

// ORU_R01_PatientResult_OrderObservation contains order observation segments.
type ORU_R01_PatientResult_OrderObservation struct {
	ORC any
	OBR any
}

// ORU_R01_PatientResult_OrderObservation_Observation contains observation result segments and notes.
type ORU_R01_PatientResult_OrderObservation_Observation struct {
	OBX []any
	NTE []any
}
