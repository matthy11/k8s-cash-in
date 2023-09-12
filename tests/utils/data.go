package utils

type DataOptions struct {
	MaxInvalidItems int32
	MaxValidItems   int32
}

type ValidData struct {
	Data         map[string]interface{}
	ExpectedData map[string]interface{}
}

// TODO: replicate heypay-accounts-engine/tests logic
func GetAllPossibleInvalidAndValidCreationData(fullData map[string]interface{}, mandatoryFields []string, defaultValues map[string]interface{}, options DataOptions) ([]map[string]interface{}, []ValidData) {
	var invalidData []map[string]interface{}
	for i := 0; i < len(mandatoryFields)-1; i++ {
		var auxData = map[string]interface{}{}
		for j := 0; j < i; j++ {
			auxData[mandatoryFields[j]] = fullData[mandatoryFields[j]]
		}
		invalidData = append(invalidData, auxData)
	}
	return invalidData, []ValidData{{
		Data:         fullData,
		ExpectedData: nil,
	}}
}
