package example

import "time"

// @dym/wired:true
// @dym/table:order
type Order struct {
	ID           int       `json:"id"`
	ReID         int       `json:"re_id"`
	BerthID      int       `json:"berth_id"`
	BizOrderID   string    `json:"biz_order_id"`
	BigImgPath   string    `json:"big_img_path"`
	CarColor     int       `json:"car_color"`
	CarType      int       `json:"car_type"`
	CarModel     int       `json:"car_model"`
	CarBrand     int       `json:"car_brand"`
	PlateNumber  string    `json:"plate_number"`
	PlateColor   int       `json:"plate_color"`
	PlateType    int       `json:"plate_type"`
	PlatePicture string    `json:"plate_picture"`
	EnterSN      string    `json:"enter_sn"`
	ExitSN       string    `json:"exit_sn"`
	EnterTime    time.Time `json:"enter_time"`
	ExitTime     time.Time `json:"exit_time"`
	Remark       string    `json:"remark"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
