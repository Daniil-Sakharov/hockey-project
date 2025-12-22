package entities

import "testing"

func TestNewUserSession(t *testing.T) {
	session := NewUserSession(12345)

	if session.UserID != 12345 {
		t.Errorf("UserID = %d, want 12345", session.UserID)
	}
	if session.ResultMsgIDs == nil {
		t.Error("ResultMsgIDs should not be nil")
	}
}

func TestUserSession_Reset(t *testing.T) {
	session := NewUserSession(12345)
	session.CurrentView = "filter_menu"
	session.CurrentPage = 5
	session.WaitingForInput = "fio"

	session.Reset()

	if session.CurrentView != "" {
		t.Errorf("CurrentView = %s, want empty", session.CurrentView)
	}
	if session.CurrentPage != 0 {
		t.Errorf("CurrentPage = %d, want 0", session.CurrentPage)
	}
	if session.WaitingForInput != "" {
		t.Errorf("WaitingForInput = %s, want empty", session.WaitingForInput)
	}
}

func TestUserSession_SetView(t *testing.T) {
	session := NewUserSession(12345)
	session.CurrentPage = 5

	session.SetView("search_results")

	if session.CurrentView != "search_results" {
		t.Errorf("CurrentView = %s, want search_results", session.CurrentView)
	}
	if session.CurrentPage != 0 {
		t.Errorf("CurrentPage = %d, want 0 (should reset)", session.CurrentPage)
	}
}

func TestUserSession_Pagination(t *testing.T) {
	session := NewUserSession(12345)

	session.NextPage()
	if session.CurrentPage != 1 {
		t.Errorf("CurrentPage = %d, want 1", session.CurrentPage)
	}

	session.NextPage()
	if session.CurrentPage != 2 {
		t.Errorf("CurrentPage = %d, want 2", session.CurrentPage)
	}

	session.PrevPage()
	if session.CurrentPage != 1 {
		t.Errorf("CurrentPage = %d, want 1", session.CurrentPage)
	}

	session.PrevPage()
	session.PrevPage() // Should not go below 0
	if session.CurrentPage != 0 {
		t.Errorf("CurrentPage = %d, want 0", session.CurrentPage)
	}
}

func TestUserSession_IsWaitingForInput(t *testing.T) {
	session := NewUserSession(12345)

	if session.IsWaitingForInput() {
		t.Error("IsWaitingForInput() should be false initially")
	}

	session.WaitingForInput = "fio"
	if !session.IsWaitingForInput() {
		t.Error("IsWaitingForInput() should be true")
	}
}
