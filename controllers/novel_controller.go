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

func AddNovel(c *fiber.Ctx) error {
	//expectedd body structure
	type NovelInput struct {
		Title         string `json:"title"`
		Summary       string `json:"summary"`
		Language      string `json:"language"`
		CoverImage    string `json:"coverImage"`
		AuthorIds     []uint `json:"authorIds"`
		GenreIds      []uint `json:"genreIds"`
		TagIds        []uint `json:"tagIds"`
		PublishedDate string `json:"publishedDate"` // optional, ISO8601 string
		PageCount     int    `json:"pageCount"`     // optional
		Status        string `json:"status"`        // optional
		Slug          string `json:"slug"`          // optional, can be generated
	}

	var input NovelInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	//validate required fields
	if strings.TrimSpace(input.Title) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title is required",
		})
	}
	//generate slug if not provided
	slug := input.Slug
	if slug == "" {
		slug = strings.ToLower(strings.ReplaceAll(input.Title, " ", "-"))
	}

	//parse published date if provided
	var publishedDate time.Time
	if input.PublishedDate != "" {
		var err error
		publishedDate, err = time.Parse(time.RFC3339, input.PublishedDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid published date format",
			})
		}
	}

	//fetch authors, genres, and tags
	var authors []models.Author
	if len(input.AuthorIds) > 0 {
		if err := db.DB.Find(&authors, input.AuthorIds).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch authors",
			})
		}
	}
	var genres []models.Genre
	if len(input.GenreIds) > 0 {
		if err := db.DB.Find(&genres, input.GenreIds).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch genres",
			})
		}
	}
	var tags []models.Tags
	if len(input.TagIds) > 0 {
		if err := db.DB.Find(&tags, input.TagIds).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch tags",
			})
		}
	}
	//create novel
	novel := models.Novel{
		Title:         input.Title,
		Slug:          slug,
		Summary:       input.Summary,
		PublishedDate: publishedDate,
		PageCount:     input.PageCount,
		Language:      input.Language,
		Status:        input.Status,
		CoverImageURL: input.CoverImage,
		Authors:       authors,
		Genres:        genres,
		Tags:          tags,
	}

	if err := db.DB.Create(&novel).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add novel to db",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Novel added successfully",
		"novel":   novel,
	})
}

func ViewNovels(c *fiber.Ctx) error {
	var novels []models.Novel
	if err := db.DB.Preload("Authors").Preload("Genres").Preload("Tags").Find(&novels).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch novels from db",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"novels": novels,
	})
}

func UpdateNovel(c *fiber.Ctx) error {

	slug := c.Params("slug")

	var novel models.Novel

	if err := db.DB.Preload("Authors").Preload("Genres").Preload("Tags").Where("slug =?", slug).First(&novel).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Novel not found",
		})
	}
	//parse request body
	type UpdateNovelInput struct {
		Title         *string `json:"title"`
		Summary       *string `json:"summary"`
		Language      *string `json:"language"`
		CoverImage    *string `json:"coverImage"`
		AuthorIds     *[]uint `json:"authorIds"`
		GenreIds      *[]uint `json:"genreIds"`
		TagIds        *[]uint `json:"tagIds"`
		PublishedDate *string `json:"publishedDate"`
		PageCount     *int    `json:"pageCount"`
		Status        *string `json:"status"`
		Slug          *string `json:"slug"`
	}

	var input UpdateNovelInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	//update fields if provided
	if input.Title != nil {
		novel.Title = *input.Title
	}
	if input.Summary != nil {
		novel.Summary = *input.Summary
	}
	if input.Language != nil {
		novel.Language = *input.Language
	}
	if input.CoverImage != nil {
		novel.CoverImageURL = *input.CoverImage
	}
	if input.PageCount != nil {
		novel.PageCount = *input.PageCount
	}
	if input.Status != nil {
		novel.Status = *input.Status
	}
	if input.Slug != nil && *input.Slug != "" {
		novel.Slug = *input.Slug
	}
	if input.PublishedDate != nil && *input.PublishedDate != "" {
		publishedDate, err := time.Parse(time.RFC3339, *input.PublishedDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid published date format",
			})
		}
		novel.PublishedDate = publishedDate
	}
	// Update associations if provided
	if input.AuthorIds != nil {
		var authors []models.Author
		if err := db.DB.Find(&authors, *input.AuthorIds).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch authors",
			})
		}
		if err := db.DB.Model(&novel).Association("Authors").Replace(authors); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update authors",
			})
		}
	}
	if input.GenreIds != nil {
		var genres []models.Genre
		if err := db.DB.Find(&genres, *input.GenreIds).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch genres",
			})
		}
		if err := db.DB.Model(&novel).Association("Genres").Replace(genres); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update genres",
			})
		}
	}
	if input.TagIds != nil {
		var tags []models.Tags
		if err := db.DB.Find(&tags, *input.TagIds).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch tags",
			})
		}
		if err := db.DB.Model(&novel).Association("Tags").Replace(tags); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update tags",
			})
		}
	}
	//save update novel 
	if err:= db.DB.Save(&novel).Error; err!=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update novel",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Novel updated successfully",
		"novel":   novel,
	})
}
