package helperstruct

type FilterQuery struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Query    string `json:"query"`   //search key word
	Filter   string `json:"filter"`  //to specify the column name
	SortBy   string `json:"sort_by"` //to specify column to set the sorting
	SortDesc bool   `json:"sort_desc"`
}

type ProjectManagement struct {
	IsManagementNeeded bool
	ModuleNumber       int
	ProjectId          string
}

type ModuleManagement struct {
	ModuleDetails     []string
	ManagementId      string
}
