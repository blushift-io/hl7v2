package v21

import "time"

type EVN struct {
	EventTypeCode        string    `hl7:"EVN.1"`
	DateTimeOfEvent      time.Time `hl7:"EVN.2"`
	DateTimePlannedEvent time.Time `hl7:"EVN.3"`
	EventReasonCode      string    `hl7:"EVN.4"`
}

type DG1 struct {
	SetID                   string    `hl7:"DG1.1"`
	DiagnosisCodingMethod   string    `hl7:"DG1.2"`
	DiagnosisCode           string    `hl7:"DG1.3"`
	DiagnosisDescription    string    `hl7:"DG1.4"`
	DiagnosisDateTime       time.Time `hl7:"DG1.5"`
	DiagnosisDrgType        string    `hl7:"DG1.6"`
	MajorDiagnosticCategory string    `hl7:"DG1.7"`
	DiagnosticRelatedGroup  string    `hl7:"DG1.8"`
	DrgApprovalIndicator    string    `hl7:"DG1.9"`
	DrgGrouperReviewCode    string    `hl7:"DG1.10"`
	OutlierType             string    `hl7:"DG1.11"`
	OutlierDays             int       `hl7:"DG1.12"`
	OutlierCost             float64   `hl7:"DG1.13"`
	GrouperVersionAndType   string    `hl7:"DG1.14"`
}
