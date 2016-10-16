package gen

import (
	"fmt"
	"strings"
)

/*
func (z *DemoType) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field

	// We treat empty fields as if they were Nil on the wire.
	var decodeMsgFieldOrder_ = []string{"Name", "BirthDay"}

	const maxFields_ = 2
*/
var templateDecodeMsg = `
	// -- templateDecodeMsg starts here--
    var totalEncodedFields_ uint32
	totalEncodedFields_, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	encodedFieldsLeft_ := totalEncodedFields_
	missingFieldsLeft_ := maxFields_ - totalEncodedFields_

	var nextMiss_ int32 = -1
	var found_ [maxFields_]bool
	var curField_ string

doneWithStruct_:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft_ > 0 || missingFieldsLeft_ > 0 {
        //fmt.Printf("encodedFieldsLeft: %%v, missingFieldsLeft: %%v, found: '%%v', fields: '%%#v'\n", encodedFieldsLeft_, missingFieldsLeft_, msgp.ShowFound(found_[:]), decodeMsgFieldOrder_)
		if encodedFieldsLeft_ > 0 {
			encodedFieldsLeft_--
			field, err = dc.ReadMapKeyPtr()
			if err != nil {
				return
			}
			curField_ = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss_ < 0 {
				// tell the reader to only give us Nils
				// until further notice.
				dc.PushAlwaysNil()
				nextMiss_ = 0
			}
			for nextMiss_ < maxFields_ && found_[nextMiss_] {
				nextMiss_++
			}
			if nextMiss_ == maxFields_ {
				// filled all the empty fields!
				break doneWithStruct_
			}
			missingFieldsLeft_--
			curField_ = decodeMsgFieldOrder_[nextMiss_]
		}
        //fmt.Printf("switching on curField: '%%v'\n", curField_)
		switch curField_ {
		// -- templateDecodeMsg ends here --
`

/*
		case "Name":
			found_[0] = true
			z.Name, err = dc.ReadString()
			if err != nil {
				return
			}
		case "BirthDay":
			found_[1] = true
			z.BirthDay, err = dc.ReadTime()
			if err != nil {
				return
			}
            default:
                err = dc.Skip()
                if err != nil {
                   return
                }
		} // end switch curField_
	} // end for
	if nextMiss_ != -1 {
		dc.PopAlwaysNil()
	}
	return
}
*/

func genDecodeMsgTemplate(n int) (template, nStr string) {
	nStr = fmt.Sprintf("%v%v", n, randIdent())
	return strings.Replace(templateDecodeMsg, `_`, nStr, -1), nStr
}

var templateUnmarshalMsg = `
	// -- templateUnmarshalMsg starts here--
    var totalEncodedFields_ uint32
    if !nbs.AlwaysNil {
	    totalEncodedFields_, bts, err = nbs.ReadMapHeaderBytes(bts)
	    if err != nil { 
          panic(err)
	  	  return
	    }
    }
	encodedFieldsLeft_ := totalEncodedFields_
	missingFieldsLeft_ := maxFields_ - totalEncodedFields_

	var nextMiss_ int32 = -1
	var found_ [maxFields_]bool
	var curField_ string

doneWithStruct_:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft_ > 0 || missingFieldsLeft_ > 0 {
        //fmt.Printf("encodedFieldsLeft: %%v, missingFieldsLeft: %%v, found: '%%v', fields: '%%#v'\n", encodedFieldsLeft_, missingFieldsLeft_, msgp.ShowFound(found_[:]), unmarshalMsgFieldOrder_)
		if encodedFieldsLeft_ > 0 {
			encodedFieldsLeft_--
			field, bts, err = nbs.ReadMapKeyZC(bts)
			if err != nil { 
                panic(err)
				return
			}
			curField_ = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss_ < 0 {
				// set bts to contain just mnil (0xc0)
				bts = nbs.PushAlwaysNil(bts)
				nextMiss_ = 0
			}
			for nextMiss_ < maxFields_ && found_[nextMiss_] {
				nextMiss_++
			}
			if nextMiss_ == maxFields_ {
				// filled all the empty fields!
				break doneWithStruct_
			}
			missingFieldsLeft_--
			curField_ = unmarshalMsgFieldOrder_[nextMiss_]
		}
        //fmt.Printf("switching on curField: '%%v'\n", curField_)
		switch curField_ {
		// -- templateUnmarshalMsg ends here --
`

func genUnmarshalMsgTemplate(n int) (template, nStr string) {
	nStr = fmt.Sprintf("%v%v", n, randIdent())
	return strings.Replace(templateUnmarshalMsg, `_`, nStr, -1), nStr
}
