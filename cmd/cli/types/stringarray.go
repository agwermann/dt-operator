package types

// var (
// 	ErrUnmarshalNotSupported = errors.New("Unmarshal not supported")
// )

// // Custom type that accepts a String or Array of string when Unmarshalling JSON object
// // It converts always to a array of string
// type StringArray []string

// func (sa *StringArray) UnmarshalJSON(data []byte) error {
// 	var jsonObject interface{}
// 	err := json.Unmarshal(data, jsonObject)
// 	if err != nil {
// 		return err
// 	}
// 	switch obj := jsonObject.(type) {
// 	case string:
// 		*sa = []string{obj}
// 		return nil
// 	case []interface{}:
// 		array := make([]string, len(obj))
// 		for _, v := range obj {
// 			value, ok := v.(string)
// 			if !ok {
// 				return ErrUnmarshalNotSupported
// 			}
// 			array = append(array, value)
// 		}
// 		*sa = array
// 	}
// 	return ErrUnmarshalNotSupported
// }
