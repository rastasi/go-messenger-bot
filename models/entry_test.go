package models

import (
	"testing"
)

func TestAnswerJSON(t *testing.T) {
	a := Answer{
		Object: "test",
	}

	result, err := a.JSON()
	if err != nil {
		t.Error(err)
	}

	expected := `{"object":"test","entry":null}`
	if string(result) != expected {
		t.Errorf("Result should be %s, instead of %s", expected, string(result))
	}
}
