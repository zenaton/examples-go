package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var SequentialWorkflow = workflow.New("SequentialWorkflow",
	func() (interface{}, error) {

		var a int
		tasks.A.New().Execute().Output(&a)

		if a == 0 {
			tasks.B.New().Execute()
		} else {
			tasks.C.New().Execute()
		}

		return nil, nil
	})

//var ModerationWorkflow = workflow.NewCustom("ModerationWorkflow", &Moderation{})
//
//type Moderation struct {
//	Review    review
//	User      user
//	Moderator moderator
//}
//
//func (m *Moderation) Init(review review, user user, moderator moderator) {
//	m.Review = review
//	m.User = user
//	m.Moderator = moderator
//}
//
//func (m *Moderation) Handle() (interface{}, error) {
//
//	SendReminderEmailToModerator.New(m.Review, m.Moderator).Execute()
//
//	var event ReviewModeratedEvent
//	task.Wait().ForEvent("ReviewModeratedEvent").Days(2).Execute().Output(&event)
//
//	if event != (ReviewModeratedEvent{}) && event.Rejected {
//		SendRejectionEmailToUser.New(m.User, m.Review).Execute()
//	} else {
//		ValidateReview.New(m.Review).Execute()
//	}
//
//	return nil, nil
//}
//
//func (m *Moderation) ID() string {
//	return m.Review.ID
//	reviewID := "bob"
//	ModerationWorkflow.WhereID(reviewID).Send("ReviewModeratedEvent", ReviewModeratedEvent{Rejected: true})
//	return ""
//}

//
//
//
//
//type ReviewModeratedEvent struct {
//	Rejected bool
//}
//type review struct {
//	ID string
//}
//type user struct{}
//type moderator struct{}
//
//var SendReminderEmailToModerator = task.NewCustom("SendReminderEmailToModerator", &sendReminderEmailToModerator{})
//var SendRejectionEmailToUser = task.NewCustom("SendReminderEmailToModerator", &sendReminderEmailToModerator{})
//var ValidateReview = task.NewCustom("SendReminderEmailToModerator", &sendReminderEmailToModerator{})
//
//type sendReminderEmailToModerator struct {
//	Review    review
//	Moderator moderator
//}
//
//func (s *sendReminderEmailToModerator) Init(review review, moderator moderator) {
//	s.Review = review
//	s.Moderator = moderator
//}
//
//func (s *sendReminderEmailToModerator) Handle() (interface{}, error) {
//	return nil, nil
//}
