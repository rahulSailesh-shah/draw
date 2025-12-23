package inngest

import (
	"context"
	"fmt"
	"time"

	"github.com/inngest/inngestgo"
)

func (i *Inngest) RegisterFunctions() error {
	err := i.postProcessMeeting()
	return err
}

type SessionTranscript struct {
	Segments []SessionTranscriptSegment `json:"segments"`
}

type SessionTranscriptSegment struct {
	Role      string    `json:"role"`      // "user" or "ai"
	Name      string    `json:"name"`      // Speaker's name
	Content   string    `json:"content"`   // Transcript text
	Timestamp time.Time `json:"timestamp"` // When the segment was captured
}

func (i *Inngest) PostProcessMeeting(ctx context.Context, meetingId string, userId string) error {
	fmt.Println("[--] Meeting post-processing event sent", "meetingID", meetingId)
	_, err := i.client.Send(ctx, inngestgo.Event{
		Name: "conversense/post-process-meeting",
		Data: map[string]any{
			"meetingId": meetingId,
			"userId":    userId,
		},
	})
	return err
}

func (i *Inngest) postProcessMeeting() error {
	type PostProcessEventData struct {
		MeetingId string `json:"meetingId"`
		UserID    string `json:"userId"`
	}

	_, err := inngestgo.CreateFunction(
		i.client,
		inngestgo.FunctionOpts{
			ID:   "post-process-meeting",
			Name: "Post Process Meeting",
		},
		inngestgo.EventTrigger("conversense/post-process-meeting", nil),
		func(ctx context.Context, input inngestgo.Input[PostProcessEventData]) (any, error) {
			return nil, nil
		},
	)
	return err
}

