package bsgostuff_helpers

import (
	"fmt"

	"github.com/gosimple/slug"
)

type SlugChecker func(slug string) (bool, error)

func GenerateUniqueSlug(title string, checker SlugChecker) (string, error) {
	baseSlug := slug.Make(title)
	slugToCheck := baseSlug
	suffix := 0

	for {
		exists, err := checker(slugToCheck)
		if err != nil {
			return "", fmt.Errorf("slug check failed: %w", err)
		}

		if !exists {
			return slugToCheck, nil
		}

		suffix++
		slugToCheck = fmt.Sprintf("%s-%d", baseSlug, suffix)
	}
}
