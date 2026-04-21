// -*- coding: utf-8 -*-
//
// © Copyright 2023 GSI Helmholtzzentrum für Schwerionenforschung
//
// This software is distributed under
// the terms of the GNU General Public Licence version 3 (GPL Version 3),
// copied verbatim in the file "LICENCE".

package sources

import (
	"testing"
)

func TestChangelogTarget(t *testing.T) {
	testBlock := `mdd.lustre-MDT0000.changelog_users=`
	expected := "lustre-MDT0000"

	result, err := regexCaptureChangelogTarget(testBlock)

	if err != nil {
		t.Fatal(err)
	}

	if expected != result {
		t.Fatalf("No changelog target found. Expected: %s, Got %s", expected, result)
	}

	testBlock = `mdd..changelog_users=`

	_, err = regexCaptureChangelogTarget(testBlock)

	if err == nil {
		t.Fatal("Expected error if not changelog target has been found")
	}
}

func TestChangelogCurrentIndex(t *testing.T) {
	testBlock := `mdd.lustre-MDT0000.changelog_users=
	current index: 34
	ID    index (idle seconds)
	cl1   0 (1725676)
	cl2   34 (28)`
	expected := float64(34)

	result, err := regexCaptureChangelogCurrentIndex(testBlock)

	if err != nil {
		t.Fatal(err)
	}

	if expected != result {
		t.Fatalf("Retrieved an unexpected value. Expected: %f, Got %f", expected, result)
	}

	testBlock = `mdd.lustre-MDT0000.changelog_users=
	current index: 0`
	expected = 0

	result, err = regexCaptureChangelogCurrentIndex(testBlock)

	if err != nil {
		t.Fatal(err)
	}

	if expected != result {
		t.Fatalf("Retrieved an unexpected value. Expected: %f, Got %f", expected, result)
	}

	testBlock = `mdd.lustre-MDT0000.changelog_users=
	ID    index (idle seconds)`

	_, err = regexCaptureChangelogCurrentIndex(testBlock)

	if err == nil {
		t.Fatal("Expected error if no current changelog index has been found")
	}
}

func TestChangelogUser(t *testing.T) {
	testBlock := `mdd.lustre-MDT0000.changelog_users=
	current index: 34
	ID    index (idle seconds)
	cl1   0 (1725676)
	cl2   34 (28)`

	result := regexCaptureChangelogUser(testBlock)

	if len(result) != 2 {
		t.Fatalf("Retrieved unexpected length of changelog reader. Expected: %d, Got: %d", 2, len(result))
	}

	expected := "cl1   0 (1725676)"
	matched := result[0][0]

	if expected != matched {
		t.Fatalf("Retrieved an unexpected value. Expected: %s, Got: %s", expected, matched)
	}

	expected = "cl2   34 (28)"
	matched = result[1][0]

	if expected != matched {
		t.Fatalf("Retrieved an unexpected value. Expected: %s, Got: %s", expected, matched)
	}
}

func TestHsmTarget(t *testing.T) {
	testBlock := `mdt.lustrefs-MDT0000.hsm.agents=`
	expected := "lustrefs-MDT0000"

	result, err := regexCaptureHsmTarget(testBlock)
	if err != nil {
		t.Fatal(err)
	}
	if expected != result {
		t.Fatalf("No HSM target found. Expected: %s, Got: %s", expected, result)
	}

	testBlock = `mdt..hsm.agents=`
	_, err = regexCaptureHsmTarget(testBlock)
	if err == nil {
		t.Fatal("Expected error when no HSM target is found")
	}
}

func TestHsmAgentRegex(t *testing.T) {
	line := `uuid=a1b1c2d3-e4f5-6789-abcd-ef0123456789 archive_id=1 requests=[current:3 ok:120 errors:2]`

	match := hsmAgentRegexPattern.FindStringSubmatch(line)
	if match == nil {
		t.Fatal("Expected HSM agent regex to match")
	}
	if match[1] != "a1b1c2d3-e4f5-6789-abcd-ef0123456789" {
		t.Fatalf("Unexpected uuid. Got: %s", match[1])
	}
	if match[2] != "1" {
		t.Fatalf("Unexpected archive_id. Got: %s", match[2])
	}
	if match[3] != "3" {
		t.Fatalf("Unexpected current. Got: %s", match[3])
	}
	if match[4] != "120" {
		t.Fatalf("Unexpected ok. Got: %s", match[4])
	}
	if match[5] != "2" {
		t.Fatalf("Unexpected errors. Got: %s", match[5])
	}

	// Test with archive_id=ANY
	line = `uuid=fedcba98-7654-3210-fedc-ba9876543210 archive_id=ANY requests=[current:0 ok:10 errors:1]`
	match = hsmAgentRegexPattern.FindStringSubmatch(line)
	if match == nil {
		t.Fatal("Expected HSM agent regex to match with archive_id=ANY")
	}
	if match[2] != "ANY" {
		t.Fatalf("Unexpected archive_id. Expected: ANY, Got: %s", match[2])
	}
}

func TestHsmActionRegex(t *testing.T) {
	line := `action=ARCHIVE archive#=1 fid=[0x200000403:0x1:0x0] compound/cookie=0/0 status=WAITING`

	match := hsmActionRegexPattern.FindStringSubmatch(line)
	if match == nil {
		t.Fatal("Expected HSM action regex to match")
	}
	if match[1] != "ARCHIVE" {
		t.Fatalf("Unexpected action. Got: %s", match[1])
	}
	if match[2] != "1" {
		t.Fatalf("Unexpected archive#. Got: %s", match[2])
	}
	if match[3] != "WAITING" {
		t.Fatalf("Unexpected status. Got: %s", match[3])
	}

	line = `action=RESTORE archive#=2 fid=[0x200000405:0x1:0x0] compound/cookie=4/1 status=STARTED`
	match = hsmActionRegexPattern.FindStringSubmatch(line)
	if match == nil {
		t.Fatal("Expected HSM action regex to match for RESTORE/STARTED")
	}
	if match[1] != "RESTORE" {
		t.Fatalf("Unexpected action. Got: %s", match[1])
	}
	if match[3] != "STARTED" {
		t.Fatalf("Unexpected status. Got: %s", match[3])
	}
}
