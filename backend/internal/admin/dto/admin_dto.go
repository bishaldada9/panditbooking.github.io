package dto

type DashboardMetrics struct {
	TotalUsers           int64   `json:"total_users"`
	TotalPandits         int64   `json:"total_pandits"`
	TotalBookings        int64   `json:"total_bookings"`
	TotalRevenue         float64 `json:"total_revenue"`
	PendingVerifications int64   `json:"pending_verifications"`
	ActiveBookings       int64   `json:"active_bookings"`
	FailedLogins         int64   `json:"failed_logins"`
	NewUsersToday        int64   `json:"new_users_today"`
}

type SuspendUserRequest struct {
	Reason string `json:"reason" validate:"required"`
}
