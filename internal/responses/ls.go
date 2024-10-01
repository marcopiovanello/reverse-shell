package responses

type ResponseLS struct {
	Current string   `json:"current"`
	List    []string `json:"list"`
}
