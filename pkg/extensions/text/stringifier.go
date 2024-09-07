// Copyright 2024 Nitro Agility S.r.l.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-License-Identifier: Apache-2.0

package text

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strings"
)

// stringifyObj	converts an object to a string.
func stringifyObj(obj any, exclude []string) string {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Array || val.Kind() == reflect.Slice {
		arrayString := []string{}
		for i := 0; i < val.Len(); i++ {
			arrayString = append(arrayString, fmt.Sprintf("#%s", stringifyObj(val.Index(i).Interface(), exclude)))
		}
		arrayBuilder := strings.Builder{}
		sort.Strings(arrayString)
		for _, item := range arrayString {
			arrayBuilder.WriteString(item)
		}
		return arrayBuilder.String()
	}
	return fmt.Sprintf("%v", obj)
}

// stringifyMap converts a map to a string.
func stringifyMap(obj any, exclude []string) string {
	if objMap, ok := obj.(map[string]any); ok {
		keys := make([]string, 0, len(objMap))
		for key := range objMap {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		builder := strings.Builder{}
		for _, key := range keys {
			if slices.Contains(exclude, key) {
				continue
			}
			value := objMap[key]
			if value != nil {
				if reflect.ValueOf(value).Kind() == reflect.Map {
					builder.WriteString(fmt.Sprintf("#%s#%s", key, stringifyMap(value, exclude)))
				} else {
					builder.WriteString(fmt.Sprintf("#%s#%s", key, stringifyObj(value, exclude)))
				}
			}
		}
		return builder.String()
	}
	return stringifyObj(obj, exclude)
}

// Stringify converts an object to a string.
func Stringify(obj any, exclude []string) (string, error) {
	if reflect.TypeOf(obj).Kind() == reflect.Array || reflect.TypeOf(obj).Kind() == reflect.Slice {
		return stringifyObj(obj, exclude), nil
	}

	var objMap map[string]any
	dataObj, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(dataObj, &objMap)
	if err != nil {
		return "", err
	}
	return stringifyMap(objMap, exclude), nil
}
