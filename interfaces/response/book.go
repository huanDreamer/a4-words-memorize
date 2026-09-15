package response

type (
	// IndexData 首页整页数据
	IndexData struct {
		Books      []BookInfo
		TotalBooks int
		TotalWords int
		TotalPlans int
		Now        string
	}
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
