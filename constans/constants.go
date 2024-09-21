package constans

const (
	WebDir     = "./web"
	Port       = "7540"
	DateFormat = "20060102"
	CountTasks = 15
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

type Response struct {
	ID    string `json:"id"`
	Error string `json:"error"`
}

var ErrorResponse struct {
	Error string `json:"error"`
}
