package notify

type NewIncident struct {
	// Required Fields:
	name      string `json:"name"`
	requester struct {
		email string `json:"email"`
	}
	Priority string `json:""`

	// Optional Fields:
	description string `json:""`
	due_at      string `json:""`
	assignee    struct {
		email string `json:"email"`
	}
	incidents struct {
		incident struct {
			number int `json:"number"`
		}
	}
	assets    string `json:""`
	problem   string `json:""`
	solutions struct {
		solution struct {
			number int `json:"number"`
		}
	}
	category struct {
		name string `json:"name"`
	}
}
