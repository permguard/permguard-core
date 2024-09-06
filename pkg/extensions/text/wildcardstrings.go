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

package text

import (
	"fmt"
	"regexp"
	"strings"
)

// wildcardChar is the character used to represent a wildcard in a wildcard string.
const wildcardChar = "*"

// WildcardString is a string that can contain wildcards.
type WildcardString string

// convertWildcardStringToRegexp converts a wildcard string to a regular expression pattern.
func convertWildcardStringToRegexp(wildcardString string) string {
	var pattern strings.Builder
	for i, literal := range strings.Split(wildcardString, wildcardChar) {
		if i > 0 {
			str := fmt.Sprintf(".%s", wildcardChar)
			pattern.WriteString(str)
		}
		pattern.WriteString(regexp.QuoteMeta(literal))
	}
	return pattern.String()
}

// compactWildcards removes consecutive wildcards from a wildcard string.
func compactWildcards(wildcardString string) string {
	return strings.ReplaceAll(wildcardString, fmt.Sprintf("%s%s", wildcardChar, wildcardChar), wildcardChar)
}

// wildcardMatch checks if a wildcard string matches a value.
func (a WildcardString) wildcardMatch(value string, sanitized bool) bool {
	var pattern string
	aStr := compactWildcards(string(a))
	valueStr := compactWildcards(value)
	pattern = convertWildcardStringToRegexp(aStr)
	sanitizedValue := valueStr
	pattern = fmt.Sprintf("^%s$", pattern)
	if sanitized {
		sanitizedValue = strings.ReplaceAll(valueStr, wildcardChar, "")
	}
	result, _ := regexp.MatchString(pattern, sanitizedValue)
	return result
}

// WildcardEqual checks if a wildcard string is equal to a value.
func (a WildcardString) WildcardEqual(value string) bool {
	aStr := compactWildcards(string(a))
	valueStr := compactWildcards(value)
	return aStr == valueStr
}

// WildcardInclude checks if a wildcard string includes a value.
func (a WildcardString) WildcardInclude(value string) bool {
	aStr := string(a)
	if a.WildcardEqual(value) {
		return false
	}
	aSanitizedMatch := a.wildcardMatch(value, false)
	vSanitizedMatch := WildcardString(value).wildcardMatch(aStr, false)
	if strings.ReplaceAll(aStr, wildcardChar, "") == strings.ReplaceAll(value, wildcardChar, "") {
		greater := strings.Count(aStr, wildcardChar) > strings.Count(value, wildcardChar)
		return greater && aSanitizedMatch && vSanitizedMatch
	}
	return aSanitizedMatch && !vSanitizedMatch
}
