package Models

type GeneralResponse struct {
	Status bool `json:"status"`
	Code int `json:"code"`
	Message string `json:"message"`
}

type WelcomeResponse struct {
	GeneralResponse
	Uuid string `json:"uuid"`
}

type SchedulingResponse struct{
	GeneralResponse
}

type RecommendationResponse struct {
	GeneralResponse
}

type AnimeDetailResponse struct {
	 GeneralResponse
}