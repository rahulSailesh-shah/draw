package service

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"strings"
// 	"time"

// 	"draw/internal/db/repo"
// 	"draw/internal/dto"
// 	"draw/pkg/config"
// 	"draw/pkg/inngest"
// 	"draw/pkg/livekit"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/credentials"
// 	"github.com/aws/aws-sdk-go-v2/service/s3"
// 	"github.com/google/uuid"

// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// type MeetingService interface {
// 	CreateMeeting(ctx context.Context, request dto.CreateMeetingRequest) (*dto.MeetingResponse, error)
// 	UpdateMeeting(ctx context.Context, request dto.UpdateMeetingRequest) (*dto.MeetingResponse, error)
// 	GetMeetings(ctx context.Context, request dto.GetMeetingsRequest) (*dto.PaginatedMeetingsResponse, error)
// 	GetMeeting(ctx context.Context, request dto.GetMeetingRequest) (*dto.MeetingResponse, error)
// 	DeleteMeeting(ctx context.Context, request dto.DeleteMeetingRequest) error
// 	StartMeeting(ctx context.Context, request dto.StartMeetingRequest) (string, error)
// 	GetPreSignedRecordingURL(ctx context.Context, request dto.GetPreSignedRecordingURLRequest) (string, error)
// }

// type meetingService struct {
// 	queries *repo.Queries
// 	inngest *inngest.Inngest
// 	db      *pgxpool.Pool

// 	// LiveKit configuration
// 	lkConfig     *config.LiveKitConfig
// 	geminiConfig *config.GeminiConfig
// 	awsConfig    *config.AWSConfig
// }

// func NewMeetingService(
// 	db *pgxpool.Pool,
// 	queries *repo.Queries,
// 	inngest *inngest.Inngest,
// 	lkConfig *config.LiveKitConfig,
// 	geminiConfig *config.GeminiConfig,
// 	awsConfig *config.AWSConfig,
// ) MeetingService {
// 	return &meetingService{
// 		db:           db,
// 		queries:      queries,
// 		lkConfig:     lkConfig,
// 		geminiConfig: geminiConfig,
// 		awsConfig:    awsConfig,
// 		inngest:      inngest,
// 	}
// }

// func (s *meetingService) CreateMeeting(ctx context.Context, request dto.CreateMeetingRequest) (*dto.MeetingResponse, error) {
// 	newMeeting, err := s.queries.CreateMeeting(ctx, repo.CreateMeetingParams{
// 		Name:    request.Name,
// 		UserID:  request.UserID,
// 		AgentID: request.AgentID,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return toMeetingResponse(newMeeting), nil
// }

// func (s *meetingService) UpdateMeeting(ctx context.Context, request dto.UpdateMeetingRequest) (*dto.MeetingResponse, error) {
// 	currentMeeting, err := s.queries.GetMeeting(ctx, repo.GetMeetingParams{
// 		ID:     request.ID,
// 		UserID: request.UserID,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("-- failed to get meeting --: %w", err)
// 	}

// 	if request.Name != "" {
// 		currentMeeting.Name = request.Name
// 	}

// 	if request.AgentID != uuid.Nil {
// 		currentMeeting.AgentID = request.AgentID
// 	}

// 	if request.Status != "" {
// 		currentMeeting.Status = request.Status
// 	}

// 	if request.StartTime != nil {
// 		currentMeeting.StartTime = request.StartTime
// 	}

// 	if request.EndTime != nil {
// 		currentMeeting.EndTime = request.EndTime
// 	}

// 	if request.TranscriptURL != nil {
// 		currentMeeting.TranscriptUrl = request.TranscriptURL
// 	}

// 	if request.RecordingURL != nil {
// 		currentMeeting.RecordingUrl = request.RecordingURL
// 	}

// 	if request.Summary != nil {
// 		currentMeeting.Summary = request.Summary
// 	}

// 	data, err := json.MarshalIndent(currentMeeting, "", "  ")
// 	if err != nil {
// 		return nil, fmt.Errorf("-- failed to marshal meeting --: %w", err)
// 	}
// 	fmt.Println(string(data))

// 	updatedMeeting, err := s.queries.UpdateMeeting(ctx, repo.UpdateMeetingParams{
// 		ID:            currentMeeting.ID,
// 		UserID:        currentMeeting.UserID,
// 		Name:          currentMeeting.Name,
// 		AgentID:       currentMeeting.AgentID,
// 		Status:        currentMeeting.Status,
// 		StartTime:     currentMeeting.StartTime,
// 		EndTime:       currentMeeting.EndTime,
// 		TranscriptUrl: currentMeeting.TranscriptUrl,
// 		RecordingUrl:  currentMeeting.RecordingUrl,
// 		Summary:       currentMeeting.Summary,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("-- failed to update meeting --: %w", err)
// 	}

// 	return toMeetingResponse(updatedMeeting), nil
// }

// func (s *meetingService) GetMeetings(ctx context.Context, request dto.GetMeetingsRequest) (*dto.PaginatedMeetingsResponse, error) {
// 	rows, err := s.queries.GetMeetings(ctx, repo.GetMeetingsParams{
// 		UserID:  request.UserID,
// 		Column2: request.Search,
// 		Limit:   request.Limit,
// 		Offset:  request.Offset,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	var totalCount int32
// 	if len(rows) > 0 {
// 		totalCount = int32(rows[0].TotalCount)
// 	}
// 	meetings := make([]dto.MeetingResponse, 0, len(rows))
// 	for _, row := range rows {
// 		meetings = append(meetings, dto.MeetingResponse{
// 			ID:        row.ID,
// 			UserID:    row.UserID,
// 			Name:      row.Name,
// 			AgentID:   row.AgentID,
// 			Status:    row.Status,
// 			StartTime: row.StartTime,
// 			EndTime:   row.EndTime,
// 			CreatedAt: row.CreatedAt,
// 			UpdatedAt: row.UpdatedAt,
// 			AgentDetails: &dto.AgentDetails{
// 				Name:         row.AgentName,
// 				Instructions: row.AgentInstructions,
// 			},
// 		})
// 	}

// 	currentPage := (request.Offset / request.Limit) + 1
// 	totalPages := (totalCount + request.Limit - 1) / request.Limit

// 	return &dto.PaginatedMeetingsResponse{
// 		Meetings:        meetings,
// 		HasNextPage:     currentPage < totalPages,
// 		HasPreviousPage: currentPage > 1,
// 		TotalCount:      totalCount,
// 		CurrentPage:     currentPage,
// 		TotalPages:      totalPages,
// 	}, nil
// }

// func (s *meetingService) DeleteMeeting(ctx context.Context, request dto.DeleteMeetingRequest) error {
// 	err := s.queries.DeleteMeeting(ctx, request.ID)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (s *meetingService) GetMeeting(ctx context.Context, request dto.GetMeetingRequest) (*dto.MeetingResponse, error) {
// 	meeting, err := s.queries.GetMeeting(ctx, repo.GetMeetingParams{
// 		ID:     request.ID,
// 		UserID: request.UserID,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return toMeetingAgentResponse(meeting), nil
// }

// func (s *meetingService) StartMeeting(ctx context.Context, request dto.StartMeetingRequest) (string, error) {
// 	meeting, err := s.queries.GetMeeting(ctx, repo.GetMeetingParams{
// 		ID:     request.ID,
// 		UserID: request.UserID,
// 	})
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get meeting: %w", err)
// 	}

// 	if meeting.Status != "upcoming" {
// 		return "", fmt.Errorf("meeting is not in upcoming state")
// 	}

// 	userDetails, err := s.queries.GetUserByID(ctx, request.UserID)
// 	if err != nil {
// 		return "", fmt.Errorf("user not found")
// 	}

// 	session := livekit.NewLiveKitSession(
// 		&meeting,
// 		&userDetails,
// 		s.lkConfig,
// 		s.geminiConfig,
// 		s.awsConfig,
// 		livekit.SessionCallbacks{
// 			// OnMeetingEnd: func(meetingID string, recordingURL string, transcriptURL string, err error) {
// 			// 	s.onMeetingEnd(meetingID, recordingURL, transcriptURL, err)
// 			// },
// 		},
// 	)

// 	if err := session.Start(); err != nil {
// 		return "", fmt.Errorf("failed to start session: %w", err)
// 	}
// 	startTime := time.Now()
// 	_, err = s.UpdateMeeting(ctx, dto.UpdateMeetingRequest{
// 		ID:        request.ID,
// 		UserID:    request.UserID,
// 		Status:    "active",
// 		StartTime: &startTime,
// 	})
// 	if err != nil {
// 		session.Stop()
// 		return "", fmt.Errorf("failed to update meeting: %w", err)
// 	}
// 	token, err := session.GenerateUserToken()
// 	if err != nil {
// 		session.Stop()
// 		return "", fmt.Errorf("failed to generate token: %w", err)
// 	}

// 	return token, nil
// }

// func (s *meetingService) GetPreSignedRecordingURL(ctx context.Context, request dto.GetPreSignedRecordingURLRequest) (string, error) {
// 	meeting, err := s.queries.GetMeeting(ctx, repo.GetMeetingParams{
// 		ID:     request.MeetingID,
// 		UserID: request.UserID,
// 	})
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get meeting: %w", err)
// 	}

// 	if meeting.Status != "completed" {
// 		return "", fmt.Errorf("meeting not completed yet")
// 	}

// 	// Select the appropriate URL based on fileType
// 	var s3URL *string
// 	switch request.FileType {
// 	case "recording":
// 		if meeting.RecordingUrl == nil {
// 			return "", fmt.Errorf("meeting has no recording URL")
// 		}
// 		s3URL = meeting.RecordingUrl
// 	case "transcript":
// 		if meeting.TranscriptUrl == nil {
// 			return "", fmt.Errorf("meeting has no transcript URL")
// 		}
// 		s3URL = meeting.TranscriptUrl
// 	default:
// 		return "", fmt.Errorf("invalid file type: must be 'recording' or 'transcript'")
// 	}

// 	awsCfg := aws.Config{
// 		Region:      s.awsConfig.Region,
// 		Credentials: credentials.NewStaticCredentialsProvider(s.awsConfig.AccessKey, s.awsConfig.SecretKey, ""),
// 	}
// 	s3Client := s3.NewFromConfig(awsCfg)
// 	presignClient := s3.NewPresignClient(s3Client)

// 	bucket, key, err := parseS3URL(*s3URL)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to parse S3 URL: %w", err)
// 	}

// 	presignReq, err := presignClient.PresignGetObject(context.TODO(),
// 		&s3.GetObjectInput{
// 			Bucket: &bucket,
// 			Key:    &key,
// 		},
// 		s3.WithPresignExpires(15*time.Minute),
// 	)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
// 	}

// 	return presignReq.URL, nil
// }

// func toMeetingAgentResponse(meeting repo.GetMeetingRow) *dto.MeetingResponse {
// 	return &dto.MeetingResponse{
// 		ID:            meeting.ID,
// 		Name:          meeting.Name,
// 		UserID:        meeting.UserID,
// 		AgentID:       meeting.AgentID,
// 		Status:        meeting.Status,
// 		CreatedAt:     meeting.CreatedAt,
// 		UpdatedAt:     meeting.UpdatedAt,
// 		StartTime:     meeting.StartTime,
// 		EndTime:       meeting.EndTime,
// 		TranscriptUrl: meeting.TranscriptUrl,
// 		RecordingUrl:  meeting.RecordingUrl,
// 		Summary:       meeting.Summary,
// 		AgentDetails: &dto.AgentDetails{
// 			Name:         meeting.AgentName,
// 			Instructions: meeting.AgentInstructions,
// 		},
// 	}
// }

// func toMeetingResponse(meeting repo.Meeting) *dto.MeetingResponse {
// 	return &dto.MeetingResponse{
// 		ID:            meeting.ID,
// 		Name:          meeting.Name,
// 		UserID:        meeting.UserID,
// 		AgentID:       meeting.AgentID,
// 		Status:        meeting.Status,
// 		TranscriptUrl: meeting.TranscriptUrl,
// 		RecordingUrl:  meeting.RecordingUrl,
// 		Summary:       meeting.Summary,
// 		CreatedAt:     meeting.CreatedAt,
// 		UpdatedAt:     meeting.UpdatedAt,
// 	}
// }

// func (s *meetingService) onMeetingEnd(meetingID string, recordingURL string, transcriptURL string, err error) {
// 	if err != nil {
// 		fmt.Printf("[ERROR] Meeting ended with errors: %v\n", err)
// 		return
// 	}
// 	fmt.Println("[-] Meeting ended, starting post-processing", "meetingID", meetingID, "recordingURL", recordingURL, "transcriptURL", transcriptURL)

// 	meetingUUID, err := uuid.Parse(meetingID)
// 	if err != nil {
// 		fmt.Printf("[ERROR] Failed to parse meeting ID: %v\n", err)
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	meeting, err := s.queries.GetMeetingByID(ctx, meetingUUID)
// 	if err != nil {
// 		fmt.Printf("[ERROR] Failed to get meeting for cleanup: %v\n", err)
// 		return
// 	}

// 	endTime := time.Now()
// 	updateRequest := dto.UpdateMeetingRequest{
// 		ID:      meetingUUID,
// 		UserID:  meeting.UserID,
// 		Status:  "completed",
// 		EndTime: &endTime,
// 	}

// 	if recordingURL != "" {
// 		updateRequest.RecordingURL = &recordingURL
// 	}
// 	if transcriptURL != "" {
// 		updateRequest.TranscriptURL = &transcriptURL
// 	}

// 	_, err = s.UpdateMeeting(ctx, updateRequest)
// 	if err != nil {
// 		fmt.Printf("[ERROR] Failed to update meeting on end: %v\n", err)
// 		return
// 	}

// 	fmt.Println("[-] Meeting cleanup completed successfully", "meetingID", meetingID)

// 	if err := s.inngest.PostProcessMeeting(ctx, meetingID, meeting.UserID); err != nil {
// 		fmt.Printf("[ERROR] Failed to trigger post-processing: %v\n", err)
// 	}
// }

// func parseS3URL(s3URL string) (bucket, key string, err error) {
// 	if !strings.HasPrefix(s3URL, "s3://") {
// 		return "", "", fmt.Errorf("invalid S3 URL format")
// 	}

// 	// Remove s3:// prefix
// 	path := strings.TrimPrefix(s3URL, "s3://")

// 	// Split by first /
// 	parts := strings.SplitN(path, "/", 2)
// 	if len(parts) != 2 {
// 		return "", "", fmt.Errorf("invalid S3 URL format")
// 	}

// 	return parts[0], parts[1], nil
// }
