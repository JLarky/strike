package action

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type ActionFunc func(ctx context.Context, args url.Values) (any, error)

type Action struct {
	Id     string
	Action ActionFunc
}

type ServerActions struct {
	actionMap map[string]Action
}

func NewServerActions() *ServerActions {
	return &ServerActions{}
}

func (s *ServerActions) Get(id string) Action {
	return s.actionMap[id]
}

func (s *ServerActions) GetOrFail(id string) Action {
	if val, ok := s.actionMap[id]; ok {
		return val
	}
	panic("action not found. make sure that you used serverActions.register")
}

func (s *ServerActions) Register(id string, action ActionFunc) {
	if s.actionMap == nil {
		s.actionMap = make(map[string]Action)
	}
	s.actionMap[id] = Action{
		Id:     id,
		Action: action,
	}
}

func (s *ServerActions) ConsumeForm(form url.Values) (Action, error) {
	var actionId string
	for k := range form {
		if strings.HasPrefix(k, "$ACTION_ID_") {
			actionId = strings.TrimPrefix(k, "$ACTION_ID_")
		}
	}
	if val, ok := s.actionMap[actionId]; ok {
		return val, nil
	} else {
		return Action{}, fmt.Errorf("action (%s) not found: %v", actionId, form)
	}
}

func (a Action) String() string {
	return fmt.Sprintf("Action(%s)", a.Id)
}

func (a Action) ToActionName() string {
	return "$ACTION_ID_" + a.Id
}

func (a Action) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"$strike": "action",
		"id":      a.Id,
	})
}
