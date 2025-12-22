package services

import (
	"testing"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/domain/valueobjects"
)

func TestUserStateService_GetSession(t *testing.T) {
	svc := NewUserStateService()

	session1 := svc.GetSession(123)
	if session1 == nil {
		t.Fatal("GetSession() returned nil")
	}
	if session1.UserID != 123 {
		t.Errorf("UserID = %d, want 123", session1.UserID)
	}

	session2 := svc.GetSession(123)
	if session1 != session2 {
		t.Error("GetSession() should return same session for same userID")
	}

	session3 := svc.GetSession(456)
	if session1 == session3 {
		t.Error("Different users should have different sessions")
	}
}

func TestUserStateService_UpdateFilters(t *testing.T) {
	svc := NewUserStateService()

	year := 2010
	filters := &valueobjects.SearchFilters{Year: &year}

	svc.UpdateFilters(123, filters)

	session := svc.GetSession(123)
	if session.Filters.Year == nil || *session.Filters.Year != 2010 {
		t.Error("UpdateFilters() did not update filters correctly")
	}
}

func TestUserStateService_SetCurrentView(t *testing.T) {
	svc := NewUserStateService()

	svc.SetCurrentView(123, "filter_menu")

	session := svc.GetSession(123)
	if session.CurrentView != "filter_menu" {
		t.Errorf("CurrentView = %s, want filter_menu", session.CurrentView)
	}
}

func TestUserStateService_SetLastMessageID(t *testing.T) {
	svc := NewUserStateService()

	svc.SetLastMessageID(123, 999)

	session := svc.GetSession(123)
	if session.LastMessageID != 999 {
		t.Errorf("LastMessageID = %d, want 999", session.LastMessageID)
	}
}

func TestUserStateService_ResetFilters(t *testing.T) {
	svc := NewUserStateService()

	year := 2010
	svc.UpdateFilters(123, &valueobjects.SearchFilters{Year: &year})
	svc.ResetFilters(123)

	session := svc.GetSession(123)
	if session.Filters.Year != nil {
		t.Error("ResetFilters() did not clear filters")
	}
}

func TestUserStateService_ClearSession(t *testing.T) {
	svc := NewUserStateService()

	svc.GetSession(123)
	svc.ClearSession(123)

	session := svc.GetSession(123)
	if session.CurrentView != "" {
		t.Error("ClearSession() did not clear session properly")
	}
}

func TestUserStateService_GetCurrentPage(t *testing.T) {
	svc := NewUserStateService()

	// Default page should be 1
	page := svc.GetCurrentPage(123)
	if page != 1 {
		t.Errorf("GetCurrentPage() = %d, want 1", page)
	}

	svc.SetCurrentPage(123, 5)
	page = svc.GetCurrentPage(123)
	if page != 5 {
		t.Errorf("GetCurrentPage() = %d, want 5", page)
	}
}
