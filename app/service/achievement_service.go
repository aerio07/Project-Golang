package service

import (
	"os"
	"path/filepath"
	"strings"

	"project_uas/app/model"
	"project_uas/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementService struct {
	Repo      repository.AchievementRepository
	MongoRepo repository.AchievementMongoRepository
}

func NewAchievementService(repo repository.AchievementRepository, mongoRepo repository.AchievementMongoRepository) *AchievementService {
	return &AchievementService{Repo: repo, MongoRepo: mongoRepo}
}

func getRoleAndUserID(c *fiber.Ctx) (role, userID string, ok bool) {
	u := c.Locals("user")
	if u == nil {
		return "", "", false
	}
	token, ok := u.(*jwt.Token)
	if !ok {
		return "", "", false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", false
	}
	r, rok := claims["role"].(string)
	s, sok := claims["sub"].(string)
	if !rok || !sok {
		return "", "", false
	}
	return r, s, true
}

//
// ===== LIST =====
//

func (s *AchievementService) GetAchievements(c *fiber.Ctx) error {
	role, userID, ok := getRoleAndUserID(c)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var (
		data interface{}
		err  error
	)

	switch role {
	case "Admin":
		data, err = s.Repo.GetAll()
	case "Mahasiswa":
		data, err = s.Repo.GetByStudent(userID)
	case "Dosen Wali":
		data, err = s.Repo.GetBySupervisor(userID)
	default:
		return c.SendStatus(fiber.StatusForbidden)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch achievements"})
	}
	return c.JSON(fiber.Map{"data": data})
}

//
// ===== CREATE (SRS: mongo_achievement_id NOT NULL) =====
//

type achievementUpsertReq struct {
	AchievementType string                 `json:"achievementType"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Details         map[string]interface{} `json:"details"`
	Tags            []string               `json:"tags"`
	Points          int                    `json:"points"`
}

func (s *AchievementService) CreateAchievement(c *fiber.Ctx) error {
	role, userID, ok := getRoleAndUserID(c)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if role != "Mahasiswa" {
		return c.SendStatus(fiber.StatusForbidden)
	}

	var body achievementUpsertReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid body"})
	}
	if body.AchievementType == "" || body.Title == "" || body.Description == "" {
		return c.Status(400).JSON(fiber.Map{"message": "achievementType, title, description are required"})
	}

	studentID, ok2, err := s.Repo.GetStudentIDByUserID(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to resolve student"})
	}
	if !ok2 {
		return c.Status(404).JSON(fiber.Map{"message": "student not found"})
	}

	refID := uuid.New().String()

	// 1) Mongo dulu
	doc := &model.AchievementMongo{
		AchievementRefID: refID,
		StudentID:        studentID,
		AchievementType:  body.AchievementType,
		Title:            body.Title,
		Description:      body.Description,
		Details:          body.Details,
		Tags:             body.Tags,
		Points:           body.Points,
	}

	mongoID, err := s.MongoRepo.Create(doc)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to create detail"})
	}

	// 2) Postgres insert dengan mongo id (NOT NULL aman)
	if err := s.Repo.CreateDraftWithMongo(refID, studentID, mongoID.Hex()); err != nil {
		_ = s.MongoRepo.Delete(mongoID) // rollback
		return c.Status(500).JSON(fiber.Map{"message": "failed to create reference"})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": fiber.Map{
			"id":     refID,
			"status": "draft",
			"detail": doc,
		},
	})
}

//
// ===== DETAIL =====
//

func (s *AchievementService) GetAchievementDetail(c *fiber.Ctx) error {
	refID := c.Params("id")

	role, userID, ok := getRoleAndUserID(c)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var mongoID *string
	var status string
	var okRef bool
	var err error

	switch role {
	case "Admin":
		mongoID, status, okRef, err = s.Repo.GetRefForDetailAdmin(refID)
	case "Mahasiswa":
		mongoID, status, okRef, err = s.Repo.GetRefForDetailStudent(refID, userID)
	case "Dosen Wali":
		mongoID, status, okRef, err = s.Repo.GetRefForDetailSupervisor(refID, userID)
	default:
		return c.SendStatus(fiber.StatusForbidden)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch reference"})
	}
	if !okRef || mongoID == nil {
		return c.Status(404).JSON(fiber.Map{"message": "achievement not found"})
	}

	oid, err := primitive.ObjectIDFromHex(*mongoID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "invalid mongo id"})
	}

	detail, err := s.MongoRepo.FindByID(oid)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "achievement detail not found"})
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"id":     refID,
			"status": status,
			"detail": detail,
		},
	})
}

//
// ===== UPDATE (DRAFT) =====
//

func (s *AchievementService) UpdateAchievement(c *fiber.Ctx) error {
	refID := c.Params("id")

	role, userID, ok := getRoleAndUserID(c)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if role != "Mahasiswa" {
		return c.SendStatus(fiber.StatusForbidden)
	}

	mongoID, status, okRef, err := s.Repo.GetRefForDetailStudent(refID, userID)
	if err != nil || !okRef || mongoID == nil {
		return c.Status(404).JSON(fiber.Map{"message": "achievement not found"})
	}
	if status != "draft" {
		return c.Status(422).JSON(fiber.Map{"message": "only draft can be updated"})
	}

	var body achievementUpsertReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid body"})
	}

	oid, err := primitive.ObjectIDFromHex(*mongoID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "invalid mongo id"})
	}

	update := map[string]interface{}{
		"achievementType": body.AchievementType,
		"title":           body.Title,
		"description":     body.Description,
		"details":         body.Details,
		"tags":            body.Tags,
		"points":          body.Points,
	}

	if err := s.MongoRepo.Update(oid, update); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to update achievement"})
	}

	return c.JSON(fiber.Map{"message": "achievement updated"})
}

//
// ===== DELETE (draft only) =====
//

func (s *AchievementService) DeleteAchievement(c *fiber.Ctx) error {
	refID := c.Params("id")

	role, userID, ok := getRoleAndUserID(c)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if role != "Mahasiswa" {
		return c.SendStatus(fiber.StatusForbidden)
	}

	mongoID, status, okRef, err := s.Repo.GetRefForDetailStudent(refID, userID)
	if err != nil || !okRef || mongoID == nil || status != "draft" {
		return c.Status(403).JSON(fiber.Map{"message": "achievement cannot be deleted"})
	}

	if err := s.Repo.SoftDelete(refID); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to delete achievement"})
	}

	// optional: tandai di mongo juga biar konsisten
	if oid, err := primitive.ObjectIDFromHex(*mongoID); err == nil {
		_ = s.MongoRepo.Update(oid, map[string]interface{}{"isDeleted": true})
	}

	return c.JSON(fiber.Map{"message": "achievement deleted"})
}

//
// ===== SUBMIT =====
//

func (s *AchievementService) SubmitAchievement(c *fiber.Ctx) error {
	role, userID, ok := getRoleAndUserID(c)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if role != "Mahasiswa" {
		return c.SendStatus(fiber.StatusForbidden)
	}

	if err := s.Repo.Submit(c.Params("id"), userID); err != nil {
		return c.Status(403).JSON(fiber.Map{"message": "cannot submit"})
	}
	return c.JSON(fiber.Map{"message": "achievement submitted"})
}

//
// ===== VERIFY / REJECT =====
//

func (s *AchievementService) VerifyAchievement(c *fiber.Ctx) error {
	role, userID, ok := getRoleAndUserID(c)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if role != "Dosen Wali" {
		return c.SendStatus(fiber.StatusForbidden)
	}

	if err := s.Repo.Verify(c.Params("id"), userID); err != nil {
		return c.Status(403).JSON(fiber.Map{"message": "cannot verify"})
	}
	return c.JSON(fiber.Map{"message": "achievement verified"})
}

func (s *AchievementService) RejectAchievement(c *fiber.Ctx) error {
	role, userID, ok := getRoleAndUserID(c)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if role != "Dosen Wali" {
		return c.SendStatus(fiber.StatusForbidden)
	}

	var body struct {
		Note string `json:"note"`
	}
	if err := c.BodyParser(&body); err != nil || strings.TrimSpace(body.Note) == "" {
		return c.Status(400).JSON(fiber.Map{"message": "note is required"})
	}

	if err := s.Repo.Reject(c.Params("id"), body.Note, userID); err != nil {
		return c.Status(403).JSON(fiber.Map{"message": "cannot reject"})
	}
	return c.JSON(fiber.Map{"message": "achievement rejected"})
}

//
// ===== HISTORY =====
//

func (s *AchievementService) GetAchievementHistory(c *fiber.Ctx) error {
	refID := c.Params("id")

	status, ok, err := s.Repo.GetStatusByID(refID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch achievement"})
	}
	if !ok {
		return c.Status(404).JSON(fiber.Map{"message": "achievement not found"})
	}

	history, err := s.Repo.GetImplicitHistory(refID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch history"})
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"currentStatus": status,
			"history":       history,
		},
	})
}

//
// ===== UPLOAD ATTACHMENT (DRAFT ONLY) =====
//

func (s *AchievementService) UploadAchievementAttachment(c *fiber.Ctx) error {
	refID := c.Params("id")

	role, userID, ok := getRoleAndUserID(c)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if role != "Mahasiswa" {
		return c.SendStatus(fiber.StatusForbidden)
	}

	mongoID, status, okRef, err := s.Repo.GetRefForDetailStudent(refID, userID)
	if err != nil || !okRef || mongoID == nil {
		return c.Status(404).JSON(fiber.Map{"message": "achievement not found"})
	}
	if status != "draft" {
		return c.Status(422).JSON(fiber.Map{"message": "only draft can upload attachments"})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "file is required"})
	}

	dir := filepath.Join("uploads", "achievements", refID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to prepare directory"})
	}

	filename := uuid.New().String() + "_" + strings.ReplaceAll(file.Filename, " ", "_")
	path := filepath.Join(dir, filename)

	if err := c.SaveFile(file, path); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to save file"})
	}

	fileURL := "/" + filepath.ToSlash(path)

	oid, err := primitive.ObjectIDFromHex(*mongoID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "invalid mongo id"})
	}

	attachment := model.Attachment{
		FileName: filename,
		FileURL:  fileURL,
		FileType: file.Header.Get("Content-Type"),
	}

	if err := s.MongoRepo.AddAttachment(oid, attachment); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to attach file"})
	}

	return c.Status(201).JSON(fiber.Map{"data": attachment})
}
