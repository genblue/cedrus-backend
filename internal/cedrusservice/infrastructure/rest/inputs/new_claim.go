package inputs

type NewClaim struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	TreeCount uint   `json:"treeCount"`
}
