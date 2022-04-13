package utilx

import (
    "github.com/pkg/errors"
    "os"
    "strings"
)

func CheckEnvironmentForRendering() error {
    val := os.Getenv("LC_CTYPE")
    if !strings.Contains(val, "UTF-8") {
        return errors.Errorf("LC_CTYPE must bet set to `en_US.UTF-8` or similar")
    }
    val = os.Getenv("LC_ALL")
    if !strings.Contains(val, "UTF-8") {
        return errors.Errorf("LC_ALL must bet set to `en_US.UTF-8` or similar")
    }
    return nil
}
