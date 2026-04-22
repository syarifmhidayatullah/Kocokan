package service

import (
	"errors"
	"math/rand"
	"time"

	"github.com/project/kocokan/internal/model"
	"github.com/project/kocokan/internal/repository"
)

type GroupService struct {
	groupRepo       *repository.GroupRepository
	participantRepo *repository.ParticipantRepository
	roundRepo       *repository.RoundRepository
}

func NewGroupService(gr *repository.GroupRepository, pr *repository.ParticipantRepository, rr *repository.RoundRepository) *GroupService {
	return &GroupService{gr, pr, rr}
}

func (s *GroupService) List(ownerID uint) ([]model.Group, error) {
	return s.groupRepo.ListByOwner(ownerID)
}

func (s *GroupService) Get(id, ownerID uint) (*model.Group, error) {
	return s.groupRepo.FindByID(id, ownerID)
}

func (s *GroupService) Create(ownerID uint, name, emoji, desc, periodType string, numParticipants, totalRounds int, prizeAmount int64) (*model.Group, error) {
	g := &model.Group{
		OwnerID:         ownerID,
		Name:            name,
		Emoji:           emoji,
		Description:     desc,
		NumParticipants: numParticipants,
		PeriodType:      periodType,
		PrizeAmount:     prizeAmount,
		TotalRounds:     totalRounds,
		IsActive:        true,
	}
	return g, s.groupRepo.Create(g)
}

func (s *GroupService) Update(id, ownerID uint, name, emoji, desc, periodType string, numParticipants, totalRounds int, prizeAmount int64) (*model.Group, error) {
	g, err := s.groupRepo.FindByID(id, ownerID)
	if err != nil {
		return nil, err
	}
	g.Name = name
	g.Emoji = emoji
	g.Description = desc
	g.PeriodType = periodType
	g.NumParticipants = numParticipants
	g.TotalRounds = totalRounds
	g.PrizeAmount = prizeAmount
	return g, s.groupRepo.Update(g)
}

func (s *GroupService) Delete(id, ownerID uint) error {
	return s.groupRepo.Delete(id, ownerID)
}

// Participant ops
func (s *GroupService) AddParticipant(groupID, ownerID uint, name, phone, notes string) (*model.Participant, error) {
	if _, err := s.groupRepo.FindByID(groupID, ownerID); err != nil {
		return nil, errors.New("grup tidak ditemukan")
	}
	p := &model.Participant{GroupID: groupID, Name: name, Phone: phone, Notes: notes}
	return p, s.participantRepo.Create(p)
}

func (s *GroupService) UpdateParticipant(id, groupID, ownerID uint, name, phone, notes string) (*model.Participant, error) {
	if _, err := s.groupRepo.FindByID(groupID, ownerID); err != nil {
		return nil, errors.New("grup tidak ditemukan")
	}
	ps, err := s.participantRepo.ListByGroup(groupID)
	if err != nil {
		return nil, err
	}
	var target *model.Participant
	for i := range ps {
		if ps[i].ID == id {
			target = &ps[i]
			break
		}
	}
	if target == nil {
		return nil, errors.New("peserta tidak ditemukan")
	}
	target.Name = name
	target.Phone = phone
	target.Notes = notes
	return target, s.participantRepo.Update(target)
}

func (s *GroupService) DeleteParticipant(id, groupID, ownerID uint) error {
	if _, err := s.groupRepo.FindByID(groupID, ownerID); err != nil {
		return errors.New("grup tidak ditemukan")
	}
	return s.participantRepo.Delete(id, groupID)
}

// Draw
func (s *GroupService) Draw(groupID, ownerID uint) (*model.Round, error) {
	g, err := s.groupRepo.FindByID(groupID, ownerID)
	if err != nil {
		return nil, errors.New("grup tidak ditemukan")
	}

	winnerIDs, err := s.roundRepo.WinnerIDs(groupID)
	if err != nil {
		return nil, err
	}
	wonSet := make(map[uint]bool)
	for _, id := range winnerIDs {
		wonSet[id] = true
	}

	eligible := []model.Participant{}
	for _, p := range g.Participants {
		if !wonSet[p.ID] {
			eligible = append(eligible, p)
		}
	}
	if len(eligible) == 0 {
		return nil, errors.New("semua peserta sudah pernah menang")
	}

	winner := eligible[rand.Intn(len(eligible))]
	now := time.Now()

	rounds, _ := s.roundRepo.ListByGroup(groupID)
	roundNum := len(rounds) + 1

	round := &model.Round{
		GroupID:     groupID,
		RoundNumber: roundNum,
		WinnerID:    &winner.ID,
		DrawnAt:     &now,
	}
	if err := s.roundRepo.Create(round); err != nil {
		return nil, err
	}
	round.Winner = &winner
	return round, nil
}

func (s *GroupService) UpdateWinner(roundID, groupID, ownerID uint, newWinnerID uint, notes string) (*model.Round, error) {
	if _, err := s.groupRepo.FindByID(groupID, ownerID); err != nil {
		return nil, errors.New("grup tidak ditemukan")
	}
	round, err := s.roundRepo.FindByID(roundID)
	if err != nil {
		return nil, errors.New("round tidak ditemukan")
	}
	if round.GroupID != groupID {
		return nil, errors.New("akses ditolak")
	}
	round.WinnerID = &newWinnerID
	round.Notes = notes
	return round, s.roundRepo.Save(round)
}
