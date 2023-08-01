package response

type (
	BookInfo struct {
		BookId     string
		Name       string
		WordNum    int
		StudyPlans []StudyPlan
	}
	StudyPlan struct {
		Time   string
		PlanId int64
		Num    int
		Status string
	}
)
