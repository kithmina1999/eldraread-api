package controllers

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kithmina1999/eldraread-api/config"
	"github.com/kithmina1999/eldraread-api/db"
	"github.com/kithmina1999/eldraread-api/models"
)

func AddGenre(c *fiber.Ctx) error {
	type GenreName struct {
		Name string `json:"name"`
	}

	var body GenreName

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	//convert genre name into lowercase
	name := strings.ToLower(body.Name)

	//check if genre already exist in db
	var existing models.Genre
	if err := db.DB.Where("name = ?", name).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Genre Already exists",
		})
	}

	//create new genre
	genre := models.Genre{
		Name: name,
	}

	if err := db.DB.Create(&genre).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add genre into db",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Genre created successfully",
		"genre":   genre,
	})

}

func ViewGenre(c *fiber.Ctx) error {
	var genres []models.Genre

	if err := db.DB.Find(&genres).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch genre from db",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"genres": genres,
	})
}

func AddTag(c *fiber.Ctx) error {
	type TagName struct {
		Name string `json:"name"`
	}

	var body TagName
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	lowerTag := strings.ToLower(body.Name)

	//check if already exists
	var existing models.Tags
	if err := db.DB.Where("name = ?", lowerTag).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Tag already exists",
		})
	}
	//create new tag
	tag := models.Tags{
		Name: lowerTag,
	}

	if err := db.DB.Create(&tag).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add tag into the db",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Tag created successfully",
		"tag":     tag,
	})
}

func ViewTag(c *fiber.Ctx) error {
	var tags []models.Tags

	if err := db.DB.Order("name ASC").Find(&tags).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tags form db",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"tags": tags,
	})
}

func DeleteTag(c *fiber.Ctx) error {
	id := c.Params("id")

	//check if the tag exists
	var tag models.Tags
	if err := db.DB.First(&tag, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tag not found",
		})
	}

	//delete tag
	if err := db.DB.Delete(&tag).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete tag",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"messsage": "Tag deleted successfully",
	})
}

func AddAuthor(c *fiber.Ctx) error {
	type AuthorInput struct {
		Name string `json:"name"`
		Bio  string `json:"bio"` // optional
	}

	var body AuthorInput
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	name := strings.TrimSpace(body.Name)
	if name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Author name is required",
		})
	}

	nameLower := strings.ToLower(name)

	//check for dublicated author
	var existing models.Author
	if err := db.DB.Where("LOWER(name) = ?", nameLower).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Author name has already included",
		})
	}

	//create new author
	author := models.Author{
		Name: name,
		Bio:  body.Bio,
	}

	if err := db.DB.Create(&author).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add author to db",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Author added successfully",
		"author":  author,
	})

}

func ViewAuthor(c *fiber.Ctx) error {
	var authors []models.Author

	if err := db.DB.Order("name ASC").Find(&authors).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch authors form db",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"authors": authors,
	})
}

func UploadCoverImage(c *fiber.Ctx) error {
	//get the uploaded file
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to get file",
		})
	}

	//open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open file",
		})
	}
	defer file.Close()

	//get firebase storage bucket
	bucket, err := config.FirebaseStorage()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get firebase storage",
		})
	}

	//generate unique filename
	filename := fmt.Sprintf("covers/%d_%s", time.Now().Unix(), fileHeader.Filename)

	//upload to bucket
	ctx := context.Background()
	writer := bucket.Object(filename).NewWriter(ctx)
	writer.ContentType = fileHeader.Header.Get("Content-Type")

	if _, err := io.Copy(writer, file); err != nil {
		writer.Close()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload",
		})
	}
	if err := writer.Close(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to close writer",
		})
	}
	// 6. Construct public URL (optional: adjust if you use signed URLs)
	publicURL := fmt.Sprintf(
		"https://firebasestorage.googleapis.com/v0/b/eldraread.firebasestorage.app/o/%s?alt=media",
		strings.ReplaceAll(filename, "/", "%2F"),
	)

	return c.JSON(fiber.Map{
		"message":   "Upload successful",
		"file_name": filename,
		"url":       publicURL,
	})
}
