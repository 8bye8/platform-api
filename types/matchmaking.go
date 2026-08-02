package types

import "fmt"

type SideToPlay string

const (
	White  SideToPlay = "white"
	Black  SideToPlay = "black"
	Random SideToPlay = "random"
)

func (s SideToPlay) IsValid() bool {
	switch s {
	case White, Black, Random:
		return true
	default:
		return false
	}
}

type TimeControl string

const (
	TimeControl1_0   TimeControl = "1_0"
	TimeControl2_1   TimeControl = "2_1"
	TimeControl3_0   TimeControl = "3_0"
	TimeControl3_2   TimeControl = "3_2"
	TimeControl5_0   TimeControl = "5_0"
	TimeControl10_0  TimeControl = "10_0"
	TimeControl15_0  TimeControl = "15_0"
	TimeControl15_10 TimeControl = "15_10"
	TimeControl30_0  TimeControl = "30_0"
	TimeControl60_60 TimeControl = "60_60"
)

func (t TimeControl) IsValid() bool {
	switch t {
	case TimeControl1_0,
		TimeControl2_1,
		TimeControl3_0,
		TimeControl3_2,
		TimeControl5_0,
		TimeControl10_0,
		TimeControl15_0,
		TimeControl15_10,
		TimeControl30_0,
		TimeControl60_60:
		return true
	default:
		return false
	}
}

type ChallengeParams struct {
	MinRating   int         `json:"min_rating"`
	MaxRating   int         `json:"max_rating"`
	Rating      int         `json:"rating,omitempty"`
	SideToPlay  SideToPlay  `json:"side_to_play"`
	TimeControl TimeControl `json:"time_control"`
}

func (p *ChallengeParams) Validate() error {
	if !p.SideToPlay.IsValid() {
		return fmt.Errorf("invalid side_to_play: %q", p.SideToPlay)
	}
	if !p.TimeControl.IsValid() {
		return fmt.Errorf("invalid time_control: %q", p.TimeControl)
	}
	if p.MinRating < p.Rating-200 || p.MinRating < 0 {
		return fmt.Errorf("min_rating must be non-negative and within 200 rating less than yours")
	}
	if p.MaxRating > p.Rating+200 {
		return fmt.Errorf("max_rating must be non-negative and within 200 rating more than yours")
	}
	if p.MaxRating != 0 && p.MinRating > p.MaxRating {
		return fmt.Errorf("min_rating cannot be greater than max_rating")
	}
	if p.MinRating > 5000 || p.MaxRating > 5000 || p.Rating > 5000 {
		return fmt.Errorf("fabricated rating")
	}
	return nil
}
