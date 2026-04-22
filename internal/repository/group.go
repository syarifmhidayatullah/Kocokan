package repository

import (
	"github.com/project/kocokan/internal/model"
	"gorm.io/gorm"
)

type GroupRepository struct{ db *gorm.DB }

func NewGroupRepository(db *gorm.DB) *GroupRepository { return &GroupRepository{db} }

func (r *GroupRepository) ListByOwner(ownerID uint) ([]model.Group, error) {
	var gs []model.Group
	err := r.db.Where("owner_id = ? AND is_active = true", ownerID).
		Preload("Participants").Preload("Rounds.Winner").
		Order("created_at DESC").Find(&gs).Error
	return gs, err
}

func (r *GroupRepository) FindByID(id, ownerID uint) (*model.Group, error) {
	var g model.Group
	err := r.db.Where("id = ? AND owner_id = ?", id, ownerID).
		Preload("Participants").Preload("Rounds.Winner").
		First(&g).Error
	return &g, err
}

func (r *GroupRepository) Create(g *model.Group) error {
	return r.db.Create(g).Error
}

func (r *GroupRepository) Update(g *model.Group) error {
	return r.db.Save(g).Error
}

func (r *GroupRepository) Delete(id, ownerID uint) error {
	return r.db.Model(&model.Group{}).Where("id = ? AND owner_id = ?", id, ownerID).
		Update("is_active", false).Error
}
