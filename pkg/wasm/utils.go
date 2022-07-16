/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package wasm

import (
	"reflect"
)

func filter(arr interface{}, predicate interface{}) interface{} {
	arrValue := reflect.ValueOf(arr)
	arrType := arrValue.Type()
	funcValue := reflect.ValueOf(predicate)
	resultSliceType := reflect.SliceOf(arrType.Elem())
	resultSlice := reflect.MakeSlice(resultSliceType, 0, 0)
	for i := 0; i < arrValue.Len(); i++ {
		elem := arrValue.Index(i)
		result := funcValue.Call([]reflect.Value{elem})[0].Interface().(bool)
		if result {
			resultSlice = reflect.Append(resultSlice, elem)
		}
	}
	return resultSlice.Interface()
}
