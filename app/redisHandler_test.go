package main

import "testing"

func TestDeserialization(t *testing.T) {
	res := Deserialize("id:123,var1:20,var2:30")

	if res[0][0] != "id" || res[0][1] != "123" {
		t.Fatalf("Expected to be key:%v and value:%v, insted got key: %v, value:%v", "id", "123", res[0][0], res[0][1])
	}

	if res[1][0] != "var1" || res[1][1] != "20" {
		t.Fatalf("Expected to be key:%v and value:%v, insted got key: %v, value:%v", "var1", "20", res[1][0], res[1][1])
	}
	if res[2][0] != "var2" || res[2][1] != "30" {
		t.Fatalf("Expected to be key:%v and value:%v, insted got key: %v, value:%v", "var2", "30", res[2][0], res[2][1])
	}
}

func TestEmptyDeserialization(t *testing.T) {
	res := Deserialize("")

	if len(res) > 0 {
		t.Fatalf("Should be empty array, insted got:%v", res)
	}

}



func TestGetLastEntry(t *testing.T) {
	streamEntry := [][]string{{"id","1"},{"var1","10"},{"id","2"},{"var","2"}}

	res := GetLastStreamEntry(streamEntry)

	if len(res) != 2 {
		t.Fatalf("Expected length to be %v, got:%v", 2, len(res))
	}

	if res[0][1] != "2" {
		t.Fatalf("Expected id to be %v, got:%v", "2", res[0][1])
	}
}

func TestGetLastEntrySingle(t *testing.T) {
	streamEntry := [][]string{{"id","1-1"},{"foo","bar"}}

	res := GetLastStreamEntry(streamEntry)


	if res[0][1] != "1-1" {
		t.Fatalf("Expected id to be %v, got:%v", "1-1", res[0][1])
	}
}


func TestVerifyIntegrityNonEmpty(t *testing.T) {
	streamEntry := [][]string{{"id","0-0"},{"var1","10"},{"id","0-1"},{"var","2"}}

	ok, resp := VerifyStreamIntegrity(streamEntry, "0-3")

	if !ok && resp == "The ID specified in XADD is equal or smaller than the target stream top item" {
		t.Fatalf("Expected result to be true for entry:%v on testSet: %v, with message:%v, got false with message:%v:", "0-3",streamEntry, "The ID specified in XADD is equal or smaller than the target stream top item", resp)
	}

	ok, resp = VerifyStreamIntegrity(streamEntry, "0-1")

	if !ok && resp == "The ID specified in XADD is equal or smaller than the target stream top item" {
		t.Fatalf("Expected result to be true for entry:%v on testSet: %v, with message:%v, got false with message:%v", "0-1",streamEntry, "The ID specified in XADD is equal or smaller than the target stream top item", resp)
	}
}

func TestVerifyIntegrityEmptySet(t *testing.T) {
	streamEntry := [][]string{}

	ok, _ := VerifyStreamIntegrity(streamEntry, "0-1")

	if ok {
		t.Fatalf("Expected result to be true for entry:%v on testSet: %v", "0-3",streamEntry)
	}

	ok, resp := VerifyStreamIntegrity(streamEntry, "0-0")

	if !ok && resp == "The ID specified in XADD must be greater than 0-0" {
		t.Fatalf("Expected result to be true for entry:%v on testSet: %v, with message:%v, got false with message:%v", "0-1",streamEntry, "The ID specified in XADD must be greater than 0-0", resp)
	}
}