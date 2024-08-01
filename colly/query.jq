map(select(.hasModernList == true) | {fullName, listLink, infoLink, teamListLink, testString, testStringArray, teams: (.teams | select(. != null) | map({teamName, teamLeaderName, teamLeaderCountryCode, teamLeaderCountryFlag}))})


#map(select(.hasModernList == true) | {fullName, listLink, testString, teams: (.teams | map(.teamName))})
#map(select(.hasModernList == true) | {fullName, listLink, testString, testStringArray })
#map(select(.shortName == "bofnt") | {fullName, listLink, testString, testStringArray })