package utils

import (
	"context"
	"heypay-cash-in-server/constants"
	. "heypay-cash-in-server/utils/types"
	"reflect"
	"strconv"
	"strings"
)

func ArrayIncludes(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func GenerateDigitVerifier(nationalID string) string {
	n, err := strconv.ParseInt(nationalID, 10, 0)
	if err != nil {
		return ""
	}
	s := 2
	t := 0
	for i := len(nationalID) - 1; i >= 0; i-- {
		t += s * (int(n) % 10)
		n /= 10
		s++
		if s > 7 {
			s = 2
		}
	}
	t %= 11
	t = 11 - t
	if t == 11 {
		return "0"
	} else if t == 10 {
		return "K"
	}
	return strconv.Itoa(t)
}

func GetCurrentStepCode(ctx context.Context) string {
	return GetValueFromContext(ctx, constants.ContextKeyCurrentStepCode)
}

func GetCurrentStepLabel(ctx context.Context) string {
	return GetValueFromContext(ctx, constants.ContextKeyCurrentStepLabel)
}

func GetIpAddress(ctx context.Context) string {
	return GetValueFromContext(ctx, constants.ContextKeyIPAddress)
}

func GetLoggerID(ctx context.Context) string {
	return GetValueFromContext(ctx, constants.ContextKeyLoggerID)
}

func GetValueFromContext(ctx context.Context, key constants.ContextKey) string {
	if ctx == nil {
		return "undefined"
	}
	value := ctx.Value(key)
	if value == nil {
		return "undefined"
	}
	return value.(string)
}

func IsUtil(file string) bool {
	return strings.Contains(file, "/utils/")
}

func MaskString(value string) string {
	return strings.Repeat("X", len(value)/2) + value[len(value)/2:]
}

func MaskStringPointer(value *string) string {
	if value == nil {
		return ""
	}
	return strings.Repeat("X", len(*value)/2) + (*value)[len(*value)/2:]
}

func MapStringToMapInterface(mapA map[string]string) Map {
	mapB := make(map[string]interface{}, len(mapA))
	for k, v := range mapA {
		mapB[k] = v
	}
	return mapB
}

func MaskMapKeys(mapObj Map, keys ...string) Map {
	maskedMapObj := Map{}
	for key := range mapObj {
		if ArrayIncludes(keys, key) && reflect.TypeOf(mapObj[key]).Kind() == reflect.String {
			maskedMapObj[key] = MaskString(mapObj[key].(string))
		} else {
			maskedMapObj[key] = mapObj[key]
		}
	}
	return maskedMapObj
}

func TruncateFileName(fileName string) string {
	lastIndex := strings.LastIndex(fileName, "/")
	if lastIndex >= 0 {
		fileName = (fileName)[lastIndex+1:]
	}
	return fileName
}

func TruncateFunctionName(functionName string) string {
	lastIndex := strings.LastIndex(functionName, "/")
	if lastIndex >= 0 {
		functionName = (functionName)[lastIndex+1:]
	}
	return functionName[strings.Index(functionName, ".")+1:]
}
