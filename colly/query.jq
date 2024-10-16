map({
    teams: (
        .teams
        | select(. != null)
        | map({
            #teamProfileLink,
            songs: (
                .songs
                | select(. != null)
                #| map(select(.songDownloadProcessed == null) | {
				| map({
                    songPageLink,
                    #songDownloadRaw: .songDownloadRaw,
					songTestStringArray: .songTestStringArray,
					songDownloadProcessed: .songDownloadProcessed,
                    #songTestString: .songTestString,
                })
            )
        })
    )
})


# map(select(.shortName == "boftt") | {fullName, listLink, infoLink, teamListLink, testString, testStringArray, teams: (.teams | select(. != null) | map({teamName, teamProfileLink, teamLeaderName, teamLeaderCountryCode, teamLeaderCountryFlag, teamMemberCount, teamReleasedWorksCount, teamDeclaredWorksCount, teamIsRecruiting, teamIsWithdrawn, teamIsDisqualified, teamIsWarned, teamUpdate, testString, testStringArray}))})

# map(select(.shortName == "boftt") | {fullName, listLink, infoLink, teamListLink, testString, testStringArray, teams: (.teams | select(. != null) | map (select(.teamMemberListIsCorrect == true)) | map({teamName, teamLeaderName, teamLeaderCountryCode, teamLeaderCountryFlag, teamMemberCount, teamReleasedWorksCount, teamDeclaredWorksCount, teamMemberListRaw, teamMemberListProcessed, testString, testStringArray}))})

#map(select(.shortName == "boftt") | {fullName, listLink, infoLink, teamListLink, testString, testStringArray, teams: (.teams | select(. != null) | map(select(.testStringArray | length != 2)) | map({teamName, teamLeaderName, teamLeaderCountryCode, teamLeaderCountryFlag, teamMemberCount, teamReleasedWorksCount, teamDeclaredWorksCount, testString, testStringArray}))})

#map(select(.hasModernList == true) | {fullName, listLink, testString, teams: (.teams | map(.teamName))})
#map(select(.hasModernList == true) | {fullName, listLink, testString, testStringArray })
#map(select(.shortName == "bofnt") | {fullName, listLink, testString, testStringArray })