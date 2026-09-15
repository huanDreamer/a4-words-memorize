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
		WordTrans []string // 翻译，形如 "adj 杰出的"
		Phone     string   // 音标，形如 /'æksɛs/
		Rank      int      // 顺序
		Mark      int      // 标记
	}
)

type (
	StudyPage struct {
	}
)
