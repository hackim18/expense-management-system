package constants

const (
	MinExpenseAmount  int64 = 10000
	MaxExpenseAmount  int64 = 50000000
	ApprovalThreshold int64 = 1000000
)

const (
	ExpenseStatusAwaitingApproval = "awaiting_approval"
	ExpenseStatusApproved         = "approved"
	ExpenseStatusRejected         = "rejected"
	ExpenseStatusAutoApproved     = "auto_approved"
	ExpenseStatusCompleted        = "completed"
)

const (
	ApprovalStatusApproved = "approved"
	ApprovalStatusRejected = "rejected"
)
