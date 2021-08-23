package block

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aquasecurity/defsec/definition"
	"github.com/zclconf/go-cty/cty"
)

type Reference struct {
	blockType Type
	typeLabel string
	nameLabel string
	remainder []string
	key       cty.Value
}

func newReference(parts []string) (*Reference, error) {

	var ref Reference

	if len(parts) == 0 {
		return nil, fmt.Errorf("cannot create empty reference")
	}

	blockType, err := TypeFromRefName(parts[0])
	if err != nil {
		blockType = &TypeResource
	}

	ref.blockType = *blockType

	if ref.blockType.removeTypeInReference && parts[0] != blockType.name {
		ref.typeLabel = parts[0]
		if len(parts) > 1 {
			ref.nameLabel = parts[1]
		}
	} else {
		if len(parts) > 1 {
			ref.typeLabel = parts[1]
			if len(parts) > 2 {
				ref.nameLabel = parts[2]
			} else {
				ref.nameLabel = ref.typeLabel
				ref.typeLabel = ""
			}
		}
	}

	if strings.Contains(ref.nameLabel, "[") {
		bits := strings.Split(ref.nameLabel, "[")
		ref.nameLabel = bits[0]
		keyRaw := strings.ReplaceAll(bits[1][:strings.Index(bits[1], "]")], "\"", "")
		if i, err := strconv.Atoi(keyRaw); err == nil {
			ref.key = cty.NumberIntVal(int64(i))
		} else {
			ref.key = cty.StringVal(keyRaw)
		}
	}

	if len(parts) > 3 {
		ref.remainder = parts[3:]
	}

	return &ref, nil
}

func (r *Reference) BlockType() Type {
	return r.blockType
}

func (r *Reference) TypeLabel() string {
	return r.typeLabel
}

func (r *Reference) NameLabel() string {
	return r.nameLabel
}

func (r *Reference) String() string {

	base := fmt.Sprintf("%s.%s", r.typeLabel, r.nameLabel)

	if !r.blockType.removeTypeInReference {
		base = r.blockType.Name()
		if r.typeLabel != "" {
			base += "." + r.typeLabel
		}
		if r.nameLabel != "" {
			base += "." + r.nameLabel
		}
	}

	base += r.KeyBracketed()

	for _, rem := range r.remainder {
		base += "." + rem
	}

	return base
}

func (r *Reference) RefersTo(a definition.Reference) bool {
	other := a.(*Reference)

	fmt.Println(r.String() + " ~~ " + other.String())

	if r.BlockType() != other.BlockType() {
		return false
	}
	if r.TypeLabel() != other.TypeLabel() {
		return false
	}
	if r.NameLabel() != other.NameLabel() {
		return false
	}
	if (r.Key() != "" || other.Key() != "") && r.Key() != other.Key() {
		return false
	}
	return true
}

func (r *Reference) SetKey(key cty.Value) {
	r.key = key
}
func (r *Reference) KeyBracketed() string {
	switch r.key.Type() {
	case cty.Number:
		f := r.key.AsBigFloat()
		f64, _ := f.Float64()
		return fmt.Sprintf("[%d]", int(f64))
	case cty.String:
		return fmt.Sprintf("[%q]", r.key.AsString())
	default:
		return ""
	}
}
func (r *Reference) RawKey() cty.Value {
	return r.key
}

func (r *Reference) Key() string {
	switch r.key.Type() {
	case cty.Number:
		f := r.key.AsBigFloat()
		f64, _ := f.Float64()
		return fmt.Sprintf("%d", int(f64))
	case cty.String:
		return fmt.Sprintf("%s", r.key.AsString())
	default:
		return ""
	}
}
