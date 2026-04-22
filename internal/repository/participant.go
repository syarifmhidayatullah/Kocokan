package repository

import (
	"github.com/project/kocokan/internal/model"
	"gorm.io/gorm"
)

type ParticipantRepository struct{ db *gorm.DB }

func NewParticipantRepository(db *gorm.DB) *ParticipantRepository {
	return &ParticipantRepository{db}
}

func (r *ParticipantRepository) ListByGroup(groupID uint) ([]model.Participant, error) {
	var ps []model.Participant
	err := r.db.Where("group_id = ?", groupID).Order("id ASC").Find(&ps).Error
	return ps, err
}

func (r *ParticipantRepository) Create(p *model.Participant) error {
	return r.db.Create(p).Error
}

func (r *ParticipantRepository) Update(p *model.Participant) error {
	return r.db.Save(p).Error
}

func (r *ParticipantRepository) Delete(id, groupID uint) error {
	return r.db.Where("id = ? AND group_id = ?", id, groupID).Delete(&model.Participant{}).Error
}
