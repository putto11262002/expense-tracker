package repositories

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/putto11262002/expense-tracker/api/domains"
	"gorm.io/gorm"
)

type IGroupRepository interface {
	CreateGroup(group *domains.Group) (*domains.Group, error)
	GetGroups(userID uuid.UUID) (*[]domains.Group, error)
	GroupExistByID(groupID uuid.UUID) (bool, error)
	MemberExistByID(groupID, userID uuid.UUID) (bool, error)
	AddMember(groupID uuid.UUID, userID uuid.UUID) error
	RemoveMember(groupID uuid.UUID, userID uuid.UUID) error
	SetOwner(userID uuid.UUID) error
	GetGroupByID(groupID uuid.UUID) (*domains.Group, error)
}

func (r *GroupRepository) GetGroupByID(id uuid.UUID) (group *domains.Group, err error) {
	if err = r.DB.Preload("Members").
		Preload("Owner").
		First(&group, "id = ?", id).
		Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return group, nil
}

func (r *GroupRepository) SetOwner(userID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

type GroupRepository struct {
	DB *gorm.DB
}

func (r *GroupRepository) GetGroups(userID uuid.UUID) (*[]domains.Group, error) {
	var groups []domains.Group
	if err := r.DB.Preload("Members").Preload("Owner").Where("id in (?)", r.DB.Table("user_groups").Where("user_id = ?", userID).Select("group_id")).Find(&groups).Error; err != nil {
		return nil, err
	}

	return &groups, nil
}

func (r *GroupRepository) GroupExistByID(groupID uuid.UUID) (bool, error) {
	var count int64
	if err := r.DB.Model(&domains.Group{}).Where("id = ?", groupID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GroupRepository) MemberExistByID(groupID, userID uuid.UUID) (bool, error) {
	var count int64
	if err := r.DB.Table("user_groups").Where("group_id = ? AND user_id =  ?", groupID, userID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GroupRepository) AddMember(groupID uuid.UUID, userID uuid.UUID) error {
	if err := r.DB.Table("user_groups").Create(map[string]interface{}{"user_id": userID, "group_id": groupID}).Error; err != nil {
		return err
	}
	return nil
}

func (r *GroupRepository) RemoveMember(groupID uuid.UUID, userID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{
		DB: db,
	}
}

func (r *GroupRepository) CreateGroup(group *domains.Group) (*domains.Group, error) {
	fmt.Println("Starting transaction")
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		var savedGroup = domains.Group{Name: group.Name, OwnerID: group.OwnerID}
		if err := tx.Create(&savedGroup).Error; err != nil {
			return err
		}

		group.ID = savedGroup.ID

		err := tx.Table("user_groups").Create(map[string]interface{}{"user_id": group.Owner.ID, "group_id": group.ID}).Error
		if err != nil {
			return fmt.Errorf("appending member to group: %w", err)
		}

		return nil

	})

	fmt.Printf("Finishing transaction: %v", err)
	if err != nil {
		return nil, err
	}

	return group, nil
}
