package modals

type Sharerequest struct {
	FileID      int64 `json:"file_id" binding:"required"`
	RecipientID int64 `json:"recipient_id" binding:"required"`
	CanEdit     bool  `json:"can_edit"`
	ExpiresIn   int64 `json:"expires_in" binding:"required,gt=0"`
}
