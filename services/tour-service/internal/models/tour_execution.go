package models   

import (
    "time"
)

type TourExecutionStatus string

const (
    ExecutionStarted   TourExecutionStatus = "STARTED"
    ExecutionCompleted TourExecutionStatus = "COMPLETED" 
    ExecutionAbandoned TourExecutionStatus = "ABANDONED"
)

type TourExecution struct {
    ID                 uint                `json:"id" gorm:"primaryKey"`
    TourID             uint                `json:"tourId"`
    TouristID          uint                `json:"touristId"`
    Status             TourExecutionStatus `json:"status"`
    StartTime          time.Time           `json:"startTime"`
    EndTime            *time.Time          `json:"endTime,omitempty"`
    LastActivity       time.Time           `json:"lastActivity"`
    CompletedKeyPoints []uint              `json:"completedKeyPoints" gorm:"type:integer[]"`
    StartingLatitude   float64             `json:"startingLatitude"`
    StartingLongitude  float64             `json:"startingLongitude"`
    CreatedAt          time.Time           `json:"createdAt"`
    UpdatedAt          time.Time           `json:"updatedAt"`
}

func (TourExecution) TableName() string { return "tour_executions" }