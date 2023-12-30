// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package corim

import (
	"fmt"
	"reflect"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/veraison/corim/comid"
)

var (
	em, emError = initCBOREncMode()
	dm, dmError = initCBORDecMode()
)

var (
	CoswidTag       = []byte{0xd9, 0x01, 0xf9} // 505()
	CoswidTagNumber = uint64(505)
	ComidTag        = []byte{0xd9, 0x01, 0xfa} // 506()
	ComidTagNumber  = uint64(506)
	CobomTagNumber  = uint64(508)

	corimTagsMap = map[uint64]interface{}{
		32: comid.TaggedURI(""), // #6.32 is URI, but a CoMID tag is #6.506
		// entity.go adds a tag for the nil value.
	}
)

func corimTags() cbor.TagSet {
	opts := cbor.TagOptions{
		EncTag: cbor.EncTagRequired,
		DecTag: cbor.DecTagRequired,
	}

	tags := cbor.NewTagSet()

	for tag, typ := range corimTagsMap {
		if err := tags.Add(opts, reflect.TypeOf(typ), tag); err != nil {
			panic(err)
		}
	}

	return tags
}

func initCBOREncMode() (en cbor.EncMode, err error) {
	encOpt := cbor.EncOptions{
		IndefLength: cbor.IndefLengthForbidden,
		TimeTag:     cbor.EncTagRequired,
	}
	// deeglaze: The only tags are for nil and URI. What about structural
	// restrictions on tags in nested structures? getEncodeFunc chooses the
	// encoding based on reflection. Encoding a struct depends on the _ field's cbor tag.
	// The tag can mean it's encoded as an array ("toarray"). A struct's field's tags can be
	// "omitempty", "keyasint", or the name (explicitly the first tag). Reflection can allow
	// a type to be a cbor.Tag{Number,Content}.
	return encOpt.EncModeWithTags(corimTags())
}

func initCBORDecMode() (dm cbor.DecMode, err error) {
	decOpt := cbor.DecOptions{
		IndefLength: cbor.IndefLengthForbidden,
		TimeTag:     cbor.DecTagRequired,
	}
	// deeglaze: Any unknown tag gets decoded to a cbor.Tag.
	return decOpt.DecModeWithTags(corimTags())
}

func registerCORIMTag(tag uint64, t interface{}) error {
	if _, exists := corimTagsMap[tag]; exists {
		return fmt.Errorf("tag %d is already registered", tag)
	}

	corimTagsMap[tag] = t

	var err error

	em, err = initCBOREncMode()
	if err != nil {
		return err
	}

	dm, err = initCBORDecMode()
	if err != nil {
		return err
	}

	return nil
}

func init() {
	if emError != nil {
		panic(emError)
	}
	if dmError != nil {
		panic(dmError)
	}
}
