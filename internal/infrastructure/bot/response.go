package bot

type CreateChatInviteLinkResponse struct {
	CreatesJoinRequest bool `json:"creates_join_request"`
	Creator            struct {
		FirstName string `json:"first_name"`
		ID        int    `json:"id"`
		IsBot     bool   `json:"is_bot"`
		Username  string `json:"username"`
	} `json:"creator"`
	InviteLink  string `json:"invite_link"`
	IsPrimary   bool   `json:"is_primary"`
	IsRevoked   bool   `json:"is_revoked"`
	MemberLimit int    `json:"member_limit"`
}

type RevokeChatInviteLinkResponse struct {
	CreatesJoinRequest bool `json:"creates_join_request"`
	Creator            struct {
		FirstName string `json:"first_name"`
		ID        int    `json:"id"`
		IsBot     bool   `json:"is_bot"`
		Username  string `json:"username"`
	} `json:"creator"`
	ExpireDate  int    `json:"expire_date"`
	InviteLink  string `json:"invite_link"`
	IsPrimary   bool   `json:"is_primary"`
	IsRevoked   bool   `json:"is_revoked"`
	MemberLimit int    `json:"member_limit"`
}
