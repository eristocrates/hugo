package main

import "regexp"

var bofRegex = regexp.MustCompile(`BOF[: ]?[0-9A-Z]+`)
var bmsofRegex = regexp.MustCompile(`BMS OF FIGHTERS[: ]?[0-9A-Z]+`)
var descriptionType1Regex = regexp.MustCompile(`-[^-]+-`)
var titleType1Regex = regexp.MustCompile(`(THE BMS OF FIGHTERS[^-]+)`)
var titleType2Regex = regexp.MustCompile(`(BOF[^-]+)`)
var titleType3Regex = regexp.MustCompile(`(BMS OF FIGHTERS[^-]+)`)

var manbowEventUrlPrefix = "https://manbow.nothing.sh/event/"
