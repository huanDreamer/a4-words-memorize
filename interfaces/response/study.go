package response

type (
	StudyWordResp struct {
		Date   string
		BookId string
		PlanId int64
		Name   string
		Num    int
		Status string
		Words  []WordInfo
	}
	WordInfo struct {
		HeadWord  string
		WordTrans []string // 翻译
		Rank      int      // 顺序
		Mark      int      // 标记
	}
)

type (
	StudyPage struct {
	}
)
