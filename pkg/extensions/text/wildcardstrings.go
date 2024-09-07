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
	"fmt"
	"regexp"
	"strings"
)

// wildcardChar is the character used to represent a wildcard in a wildcard string.
const wildcardChar = "*"

// WildcardString is a string that can contain wildcards.
type WildcardString string

// convertWildcardStringToRegexp converts a wildcard string to a regular expression pattern.
// Replaces wildcard '*' with '.*' to match any characters and escapes the other parts.
func convertWildcardStringToRegexp(wildcardString string) string {
	var pattern strings.Builder
	for i, literal := range strings.Split(wildcardString, wildcardChar) {
		if i > 0 {
			pattern.WriteString(".*")
		}
		pattern.WriteString(regexp.QuoteMeta(literal))
	}
	return pattern.String()
}

// wildcardMatch checks if a wildcard string matches a value.
func (a WildcardString) wildcardMatch(value string, sanitized bool) bool {
	aStr := compactWildcards(string(a))
	valueStr := compactWildcards(value)
	pattern := convertWildcardStringToRegexp(aStr)
	pattern = fmt.Sprintf("^%s$", pattern)

	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	sanitizedValue := valueStr
	if sanitized {
		sanitizedValue = strings.ReplaceAll(valueStr, wildcardChar, "")
	}

	return re.MatchString(sanitizedValue)
}

// compactWildcards removes consecutive wildcards from a wildcard string.
// Reduces "**" to "*".
func compactWildcards(wildcardString string) string {
	return strings.ReplaceAll(wildcardString, wildcardChar+wildcardChar, wildcardChar)
}

// WildcardEqual checks if two wildcard strings are equal after compacting wildcards.
// It ensures that consecutive wildcards in either string are treated as a single wildcard.
func (a WildcardString) WildcardEqual(value string) bool {
	aStr := compactWildcards(string(a))
	valueStr := compactWildcards(value)
	return aStr == valueStr
}

// WildcardInclude checks if a wildcard string includes another value.
// It compares the compacted wildcard strings and performs a detailed match.
func (a WildcardString) WildcardInclude(value string) bool {
	aStr := string(a)
	valueStr := value

	if a.WildcardEqual(value) {
		return false
	}

	aSanitizedMatch := a.wildcardMatch(valueStr, false)
	vSanitizedMatch := WildcardString(valueStr).wildcardMatch(aStr, false)

	if strings.ReplaceAll(aStr, wildcardChar, "") == strings.ReplaceAll(valueStr, wildcardChar, "") {
		greater := strings.Count(aStr, wildcardChar) > strings.Count(valueStr, wildcardChar)
		return greater && aSanitizedMatch
	}
	return aSanitizedMatch && !vSanitizedMatch
}
