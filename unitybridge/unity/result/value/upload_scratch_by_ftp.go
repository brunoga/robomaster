package value

type UploadScratchByFTP struct {
	SkillIndex int    `json:"SkillIndex"`
	GUID       string `json:"guid"`
	Sign       string `json:"sign"`
	FilePath   string `json:"filePath"`
}
