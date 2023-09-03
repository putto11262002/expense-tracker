package services

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/putto11262002/expense-tracker/api/domains"
	"github.com/putto11262002/expense-tracker/api/repositories"
	"github.com/putto11262002/expense-tracker/api/utils"
)

type GroupService struct {
	groupRepository repositories.IGroupRepository
}

func NewGroupService(groupRepository repositories.IGroupRepository) *GroupService {
	return &GroupService{
		groupRepository: groupRepository,
	}

}

// IGroupService is an interface for managing groups and group members in the application.
type IGroupService interface {
	// CreateGroup creates a new group with the specified owner and name.
	CreateGroup(owner *domains.User, name string) (*domains.Group, error)

	// DeleteGroup deletes an existing group identified by its ID.
	// Note: This operation is not currently supported and returns an error.
	DeleteGroup(id uuid.UUID) error

	// AddMember adds a member to an existing group identified by its ID.
	AddMember(id uuid.UUID, user *domains.User) error
	// RemoveMember removes a member from an existing group identified by its ID.
	// Note: This operation is not currently supported and returns an error.
	RemoveMember(id uuid.UUID, user *domains.User) error

	GetGroupByID(id uuid.UUID) (*domains.Group, error)

	GetGroupsByUserID(id uuid.UUID) (*[]domains.Group, error)
}

func (s *GroupService) GetGroupByID(id uuid.UUID) (*domains.Group, error) {
	group, err := s.groupRepository.GetGroupByID(id)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (s *GroupService) GetGroupsByUserID(id uuid.UUID) (*[]domains.Group, error) {
	groups, err := s.groupRepository.GetGroups(id)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (s *GroupService) CreateGroup(owner *domains.User, name string) (*domains.Group, error) {
	group, err := s.groupRepository.CreateGroup(domains.NewGroup(name, *owner))
	if err != nil {
		return nil, err
	}
	// add owner as a member
	return group, nil
}

func (s *GroupService) AddMember(groupID uuid.UUID, user *domains.User) error {
	exist, err := s.groupRepository.GroupExistByID(groupID)
	if err != nil {
		return err
	}
	if !exist {
		return &utils.InvalidArgumentError{
			Message: "Invalid group id",
		}
	}

	exist, err = s.groupRepository.MemberExistByID(groupID, user.ID)

	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	if err := s.groupRepository.AddMember(groupID, user.ID); err != nil {
		return err
	}
	return nil
}

func (s *GroupService) RemoveMember(groupID uuid.UUID, user *domains.User) error {
	exist, err := s.groupRepository.GroupExistByID(groupID)
	if err != nil {
		return err
	}
	if !exist {
		return &utils.InvalidArgumentError{
			Message: "Invalid group id",
		}
	}

	exist, err = s.groupRepository.MemberExistByID(groupID, user.ID)

	if err != nil {
		return err
	}

	if exist {
		return &utils.InvalidArgumentError{
			Message: fmt.Sprintf("member does not exist in group %s", groupID),
		}
	}

	if err := s.groupRepository.RemoveMember(groupID, user.ID); err != nil {
		return err
	}

	// TODO notify other services that the member have been delete e.g. notify dept service to delete
	// NOT SURE if we really have to do that tho

	return nil

}

func (s *GroupService) DeleteGroup(id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
