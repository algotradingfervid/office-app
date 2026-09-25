package approvals

// Status is requests.status.
type Status string

// The statuses of a request (design spec §5.1). None is a request that does not exist yet.
const (
	None            Status = ""
	Pending         Status = "pending"
	Approved        Status = "approved"
	Rejected        Status = "rejected"
	Cancelled       Status = "cancelled"
	CancelRequested Status = "cancel_requested"
)

// Action is something done to a request; Next says whether it is allowed from a status.
type Action string

// The actions of the transition table (design spec §7).
const (
	Submit         Action = "submit"
	RecordOnBehalf Action = "recordOnBehalf"
	Forward        Action = "forward"
	ApproveFinal   Action = "approveFinal"
	Reject         Action = "reject"
	Cancel         Action = "cancel"
	Reassign       Action = "reassign"
	RequestCancel  Action = "requestCancel"
	HRCancel       Action = "hrCancel"
	ApproveCancel  Action = "approveCancel"
	DeclineCancel  Action = "declineCancel"
)

// StepAction is approval_steps.action: what one recorded step did (design spec §5.1).
type StepAction string

// The step actions an approval step can record.
const (
	StepSubmitted       StepAction = "submitted"
	StepForwarded       StepAction = "forwarded"
	StepApprovedFinal   StepAction = "approved_final"
	StepRejected        StepAction = "rejected"
	StepCancelled       StepAction = "cancelled"
	StepCancelRequested StepAction = "cancel_requested"
	StepCancelApproved  StepAction = "cancel_approved"
	StepCancelDeclined  StepAction = "cancel_declined"
	StepReassigned      StepAction = "reassigned"
	StepRecordedByHR    StepAction = "recorded_by_hr"
)

type move struct {
	from Status
	a    Action
}

// transitions is the design's §7 table. The "who" column is enforced by the service actions, not here.
var transitions = map[move]Status{
	{None, Submit}:                   Pending,         // requester
	{None, RecordOnBehalf}:           Approved,        // HR admin
	{Pending, Forward}:               Pending,         // current approver
	{Pending, ApproveFinal}:          Approved,        // current approver
	{Pending, Reject}:                Rejected,        // current approver
	{Pending, Cancel}:                Cancelled,       // requester
	{Pending, Reassign}:              Pending,         // HR admin
	{CancelRequested, Reassign}:      CancelRequested, // HR admin
	{Approved, RequestCancel}:        CancelRequested, // requester, only if the leave has not started
	{Approved, HRCancel}:             Cancelled,       // HR admin
	{CancelRequested, ApproveCancel}: Cancelled,       // current approver
	{CancelRequested, DeclineCancel}: Approved,        // current approver
}

// Next returns the status a request in status from moves to when a is done,
// and false when the table does not allow a from that status.
func Next(from Status, a Action) (Status, bool) {
	to, ok := transitions[move{from, a}]
	return to, ok
}
