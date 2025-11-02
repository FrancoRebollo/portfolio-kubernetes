package validators

import (
	"encoding/json"
	"fmt"

	"strconv"
	"strings"

	"github.com/FrancoRebollo/api-integration-svc/internal/domain"
	"github.com/gin-gonic/gin"
)

/*
* Función que permite validar si la query está vacía
 */
func ValidateEmptyQuery(parametros map[string][]string) error {
	queryErr := domain.HealthcheckError{
		Code:    domain.ErrCodeInvalidInput,
		Message: "la query debe estar vacía",
	}
	if len(parametros) != 0 {
		return &queryErr
	}
	return nil
}

/*
* Función que permite validar si el params está vacío
 */
func ValidateEmptyParams(parametros gin.Params) error {
	paramsErr := domain.HealthcheckError{
		Code:    domain.ErrCodeInvalidInput,
		Message: "el params debe estar vacío",
	}
	if len(parametros) != 0 {
		return &paramsErr
	}
	return nil
}

/*
* Función que permite validar si el body está vacía
 */
func ValidateEmptyBody(parametros []byte) error {
	bodyErr := domain.HealthcheckError{
		Code:    domain.ErrCodeInvalidInput,
		Message: "el body debe estar vacío",
	}
	if len(parametros) != 0 {
		return &bodyErr
	}
	return nil
}

/*
* Función que permite validar el contenido de los params
 */
func ValidateParams(params gin.Params, rules map[string][]string) error {
	paramsErr := domain.HealthcheckError{
		Code: domain.ErrCodeInvalidInput,
	}
	if len(params) == 0 {
		paramsErr.Message = "los parámetros no deben estar vacíos"
		return &paramsErr
	}

	for _, param := range params {
		paramName := param.Key
		paramValue := param.Value

		ruleValues, ok := rules[paramName]
		if !ok {
			paramsErr.Message = "no hay reglas definidas para el parámetro " + paramName
			return &paramsErr
		}

		for _, rule := range ruleValues {
			err := applyRuleParams(paramName, paramValue, rule)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

/*
* Función que permite validas las reglas de cada uno de los parámetros recibidos por params
 */
func applyRuleParams(paramName string, paramValue string, rule string) error {
	paramsErr := domain.HealthcheckError{
		Code: domain.ErrCodeInvalidInput,
	}
	parts := strings.Split(rule, ":")
	ruleName := parts[0]
	var ruleValue string
	if len(parts) > 1 {
		ruleValue = parts[1]
	}

	switch ruleName {
	case "required":
		if paramValue == "" {
			paramsErr.Message = "el parámetro " + paramName + " es requerido"
			return &paramsErr
		}
	case "string":
		if _, err := strconv.Atoi(paramValue); err == nil {
			paramsErr.Message = "el parámetro " + paramName + " debe ser un string"
			return &paramsErr
		}
	case "number":
		_, err := strconv.Atoi(paramValue)
		if err != nil {
			paramsErr.Message = "el parámetro " + paramName + " debe ser un número"
			return &paramsErr
		}
	case "enum":
		options := strings.Split(ruleValue, ",")
		found := false
		for _, option := range options {
			if paramValue == option {
				found = true
				break
			}
		}
		if !found {
			paramsErr.Message = "el parámetro " + paramName + " debe ser uno de los valores permitidos: " + ruleValue
			return &paramsErr
		}
	case "maxLength":
		maxLength, err := strconv.Atoi(ruleValue)
		if err != nil {
			paramsErr.Message = "valor de regla maxLength inválido"
			return &paramsErr
		}
		if len(paramValue) > maxLength {
			paramsErr.Message = "el parámetro " + paramName + " excede la longitud máxima permitida de " + strconv.Itoa(maxLength)
			return &paramsErr
		}
	case "maxValue":
		value, _ := strconv.Atoi(ruleValue)
		if len(paramValue) > value {
			paramsErr.Message = "el parámetro " + paramName + " excede el valor máximo permitido"
			return &paramsErr
		}
	default:
		paramsErr.Message = "regla desconocida: " + rule
		return &paramsErr
	}
	return nil
}

/*
* Función que permite validar el contenido de la query
 */
func ValidateQuery(parametros map[string][]string, rules map[string][]string) error {
	queryErr := domain.HealthcheckError{
		Code: domain.ErrCodeInvalidInput,
	}
	if len(parametros) == 0 {
		queryErr.Message = "la query no debe estar vacía"
		return &queryErr
	}

	for paramName, paramValues := range parametros {
		ruleValues, ok := rules[paramName]
		if !ok {
			queryErr.Message = "no hay reglas definidas para el parámetro " + paramName
			return &queryErr
		}

		for _, rule := range ruleValues {
			err := applyRuleQuery(paramName, paramValues[0], rule, rules)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

/*
* Función que permite validar las reglas de cada uno de los parametro de la query
 */
func applyRuleQuery(paramName, paramValue, rule string, rules map[string][]string) error {
	queryErr := domain.HealthcheckError{
		Code: domain.ErrCodeInvalidInput,
	}
	parts := strings.Split(rule, ":")
	ruleName := parts[0]
	var ruleValue string
	if len(parts) > 1 {
		ruleValue = parts[1]
	}

	switch ruleName {
	case "required":
		if paramValue == "" {
			queryErr.Message = "el parámetro " + paramName + " es requerido"
			return &queryErr
		}
	case "string":
		if _, err := strconv.Atoi(paramValue); err == nil {
			queryErr.Message = "el parámetro " + paramName + " debe ser un string"
			return &queryErr
		}
	case "number":
		_, err := strconv.Atoi(paramValue)
		if err != nil {
			queryErr.Message = "el parámetro " + paramName + " debe ser un número"
			return &queryErr
		}
	case "enum":
		options := strings.Split(ruleValue, ",")
		found := false
		for _, option := range options {
			if paramValue == option {
				found = true
				break
			}
		}
		if !found {
			queryErr.Message = "el parámetro " + paramName + " debe ser uno de los valores permitidos: " + ruleValue
			return &queryErr
		}
	case "maxLength":
		maxLength, err := strconv.Atoi(ruleValue)
		if err != nil {
			queryErr.Message = "valor de regla maxLength inválido"
			return &queryErr
		}
		if len(paramValue) > maxLength {
			queryErr.Message = "el parámetro " + paramName + " excede la longitud máxima permitida de " + strconv.Itoa(maxLength)
			return &queryErr
		}
	case "maxValue":
		value, _ := strconv.Atoi(ruleValue)
		if len(paramValue) > value {
			queryErr.Message = "el parámetro " + paramName + " excede el valor máximo permitido"
			return &queryErr
		}
	default:
		queryErr.Message = "regla desconocida: " + ruleName
		return &queryErr
	}
	return nil
}

/*
* función que permite validar el contenido del body
 */
func ValidateBody(body []byte, rules map[string]map[string]string) error {
	bodyErr := domain.HealthcheckError{
		Code: domain.ErrCodeInvalidInput,
	}

	var params map[string]interface{}
	err := json.Unmarshal(body, &params)
	if err != nil {
		if strings.Contains(err.Error(), "unexpected end of JSON input") {
			bodyErr.Message = "el body no debe estar vacío"
			return &bodyErr
		}
		return err
	}

	for key, value := range rules {
		if key == "" {
			for param, rule := range value {
				if err := applyRuleBody(param, params[param], rule); err != nil {
					return err
				}
			}
		} else {
			subParams, ok := params[key].(map[string]interface{})
			if !ok {
				bodyErr.Message = "el parámetro " + key + " no es un objeto"
				return &bodyErr
			}
			for param, rule := range value {
				if err := applyRuleBody(param, subParams[param], rule); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

/*
* Función que permite validar las reglas de cada uno de los parámetros del body
 */
func applyRuleBody(paramName string, paramValue interface{}, rule string) error {
	bodyErr := domain.HealthcheckError{
		Code: domain.ErrCodeInvalidInput,
	}

	rules := strings.Split(rule, "|")
	for _, singleRule := range rules {
		switch {
		case singleRule == "required":
			if isEmpty(paramValue) {
				bodyErr.Message = fmt.Sprintf("el parámetro %s es requerido", paramName)
				return &bodyErr
			}
		case singleRule == "string":
			switch paramValue.(type) {
			case string:
				// El valor es un string
			default:
				bodyErr.Message = "el parámetro " + paramName + " debe ser un string"
				return &bodyErr
			}
		case strings.HasPrefix(singleRule, "maxLengthNumber"):
			maxLength, err := strconv.Atoi(strings.Split(singleRule, ":")[1])
			if err != nil {
				bodyErr.Message = "valor de regla maxLength inválido"
				return &bodyErr
			}
			paramNumber, ok := paramValue.(float64)
			if !ok {
				bodyErr.Message = "el parámetro " + paramName + " no es un número"
				return &bodyErr
			}
			strValue := fmt.Sprintf("%v", int(paramNumber))
			if len(strValue) > maxLength {
				bodyErr.Message = "el parámetro " + paramName + " excede la longitud máxima permitida de " + strconv.Itoa(maxLength)
				return &bodyErr
			}
		case strings.HasPrefix(singleRule, "maxLength"):
			maxLength, err := strconv.Atoi(strings.Split(singleRule, ":")[1])
			if err != nil {
				bodyErr.Message = "valor de regla maxLength inválido"
				return &bodyErr
			}
			strValue := fmt.Sprintf("%v", paramValue)
			if len(strValue) > maxLength {
				bodyErr.Message = "el parámetro " + paramName + " excede la longitud máxima permitida de " + strconv.Itoa(maxLength)
				return &bodyErr
			}
		case singleRule == "number":
			switch paramValue.(type) {
			case float64:
				// El valor es un número
			default:
				bodyErr.Message = "el parámetro " + paramName + " debe ser un número"
				return &bodyErr
			}
		case strings.HasPrefix(singleRule, "maxValue"):
			maxValue, err := strconv.ParseFloat(strings.Split(singleRule, ":")[1], 64)
			if err != nil {
				bodyErr.Message = "valor de regla maxValue inválido"
				return &bodyErr
			}
			paramNumber, ok := paramValue.(float64)
			if !ok {
				bodyErr.Message = "el parámetro " + paramName + " no es un número"
				return &bodyErr
			}
			if paramNumber > maxValue {
				bodyErr.Message = "el parámetro " + paramName + " excede el valor máximo permitido de " + strconv.FormatFloat(maxValue, 'f', -1, 64)
				return &bodyErr
			}
		case strings.HasPrefix(singleRule, "decimalLength"):
			decimalLength, err := strconv.Atoi(strings.Split(singleRule, ":")[1])
			if err != nil {
				bodyErr.Message = "valor de regla decimalLength inválido"
				return &bodyErr
			}
			strValue := fmt.Sprintf("%v", paramValue)
			parts := strings.Split(strValue, ".")
			if len(parts) != 2 || len(parts[1]) != decimalLength {
				bodyErr.Message = "el parámetro " + paramName + " debe tener " + strconv.Itoa(decimalLength) + " decimales"
				return &bodyErr
			}
		case strings.HasPrefix(singleRule, "enum"):
			options := strings.Split(strings.TrimPrefix(singleRule, "enum:"), ",")
			// Verificar si el valor está dentro de las opciones permitidas
			valueStr := fmt.Sprintf("%v", paramValue)
			found := false
			for _, option := range options {
				if valueStr == option {
					found = true
					break
				}
			}
			if !found {
				bodyErr.Message = fmt.Sprintf("el valor para el parámetro %s no está entre las opciones permitidas: %v", paramName, options)
				return &bodyErr
			}
		default:
			bodyErr.Message = "regla desconocida: " + singleRule
			return &bodyErr
		}
	}
	return nil
}

/*
* Función que permite verificar si el contenido de un parámetro es nulo o vacío
 */
func isEmpty(value interface{}) bool {
	switch v := value.(type) {
	case string:
		return v == ""
	case int, int8, int16, int32, int64:
		return v == 0
	case uint, uint8, uint16, uint32, uint64:
		return v == 0
	case float32, float64:
		return v == 0
	case bool:
		return !v
	case []interface{}:
		return len(v) == 0
	case map[interface{}]interface{}:
		return len(v) == 0
	case nil:
		return true
	default:
		return false
	}
}
