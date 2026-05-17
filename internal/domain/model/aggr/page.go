package aggr

import (
	"fmt"
	"time"
)

// `Page` represents one page inside a chapter.
type Page struct {
	// `Id` is the page identifier.
	Id string

	// `ChapterId` is the owner chapter identifier.
	ChapterId string

	// `Index` is the chapter-scoped page order.
	Index int

	// `ImageKey` is the OSS object key of the page image.
	ImageKey *string

	// `ImageUploaded` marks whether the page image upload is completed.
	ImageUploaded bool

	// `TotalUnitCount` is the total unit count on this page.
	TotalUnitCount int

	// `TranslatedUnitCount` is the translated unit count on this page.
	TranslatedUnitCount int

	// `ProofreadUnitCount` is the proofread unit count on this page.
	ProofreadUnitCount int

	// `CreatedAt` is the creation timestamp.
	CreatedAt time.Time

	// `UpdatedAt` is the update timestamp.
	UpdatedAt time.Time
}

// `GenImageKey` returns the OSS object key for the page image with the given file extension.
func (p *Page) GenImageKey(ext string) string {
	return fmt.Sprintf("chapter_%s/page_%s.%s", p.ChapterId, p.Id, ext)
}

// `PageCre` holds the create payload for page insert.
type PageCre struct {
	// `Id` is the generated page identifier.
	Id string

	// `ChapterId` is the owner chapter identifier.
	ChapterId string

	// `Index` is the chapter-scoped page order.
	Index int

	// `ImageKey` is the reserved OSS object key.
	ImageKey *string
}

// `GenImageKey` returns the OSS object key for the page image with the given file extension.
func (c *PageCre) GenImageKey(ext string) string {
	return fmt.Sprintf("chapter_%s/page_%s.%s", c.ChapterId, c.Id, ext)
}
