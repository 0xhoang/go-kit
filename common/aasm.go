package common

const (
	LOG_STATUS_SUCCESS = "success"
	LOG_STATUS_FAILED  = "failed"
)

const ListenAddress = "listen_address"

const (
	AasmStateNew        = "new"
	AasmStateApproving  = "approving"
	AasmStateRejected   = "rejected"
	AasmStateSubmitted  = "submitted"
	AasmStateProcessing = "processing"
	AasmStateConfirming = "confirming"
	AasmStateSucceed    = "succeeded"
	AasmStateErrored    = "errored"
)

type StageStatusType int

const (
	StageStatus              StageStatusType = iota //0 new
	StageStatusSuccessStatus                        //1 //success
	StageStatusFailedStatus                         //2 //failed
	StageStatusRunningStatus                        //3 //running
)
