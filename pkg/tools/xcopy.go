package tools

import "github.com/jinzhu/copier"

func Copy(src, dst interface{}) error {
	return copier.CopyWithOption(dst, src, copier.Option{DeepCopy: true, IgnoreEmpty: true})
}
