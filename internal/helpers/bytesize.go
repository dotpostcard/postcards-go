package helpers

import (
	"fmt"
	"os"
)

const byteScaler = 1024

func SizeHuman(f *os.File) string {
	stat, err := f.Stat()
	if err != nil {
		return "unknown size"
	}

	s := stat.Size()
	if s < byteScaler {
		return fmt.Sprintf("%d B", s)
	}

	div, exp := int64(byteScaler), 0
	for n := s / byteScaler; n >= byteScaler; n /= byteScaler {
		div *= byteScaler
		exp++
	}

	return fmt.Sprintf("%.1f %ciB",
		float64(s)/float64(div), "KMGTPE"[exp])
}
