import puppeteer from 'puppeteer';
import { writeFile } from 'fs';
import cliProgress from 'cli-progress';

(async () => {
  console.time('programExecution');
  // Create a new progress bar instance and use shades_classic theme
  const multibar = new cliProgress.MultiBar({
    clearOnComplete: false,
    hideCursor: true,
    format: ' {bar} | {filename} | {value}/{total}',
  }, cliProgress.Presets.shades_classic);
  const browser = await puppeteer.launch({
    headless: "new"
    /*
    headless: false,
    devtools: true
    */
  });

  const page = await browser.newPage();

  await page.setViewport({ width: 1920, height: 1080 });
  // await page.setViewport({ width: 1080, height: 1920 });

  page.setDefaultNavigationTimeout(60000);

  // Navigate to the Event page
  await page.goto('https://manbow.nothing.sh/event/event.cgi?action=List_def&event=142#186');

  const teams = new Map();

  const teamElements = await page.$$('.team_information');
  // console.log('Team Count:', teamElements.length);
  const numberOfTeamsToLimit = 1;
  // Slice the teamElements array to select a specific number of teams
  const limitedTeamElements = teamElements.slice(0, numberOfTeamsToLimit);
  const teamBar = multibar.create(limitedTeamElements.length, 0)

  let teamIndex = 1; // TODO is this needed?
  // Iterate through teams
  // for (const teamElement of teamElements.slice(0, -1)) { // for whatever reason the last element is just empty
  // Iterate through the limited teams
  for (const teamElement of limitedTeamElements) {
    teamElement.scrollIntoView();
    const teamInfo = await teamElement.$eval('.fancy-title :is(h2, h3) a', (link) => {
      const teamName = link.innerText.trim();
      const bannerImageSrc = link.querySelector('img') ? link.querySelector('img').erc : '';
      const teamPageLink = link.href;
      return { teamName, bannerImageSrc, teamPageLink };
    });
    teamBar.update({ filename: teamInfo.teamName });

    const emblemImageSrc = await teamElement.$eval('.header_emblem', (emblemElement) => {
      const dataBg = emblemElement.getAttribute('data-bg');
      const withoutPeriod = dataBg.substring(1);
      return withoutPeriod.length > 1 ? `https://manbow.nothing.sh/event${withoutPeriod}` : '';
    });
    teamInfo.emblemImageSrc = emblemImageSrc;
    teamInfo.teamImpression = await teamElement.$eval('#team_imp', (element) => element.innerText.trim());
    teamInfo.teamTotal = await teamElement.$eval('#team_total', (element) => element.innerText.trim());
    teamInfo.teamMedian = await teamElement.$eval('#team_med', (element) => element.innerText.trim());

    // console.log(`Processed team #${teamIndex}: ${teamInfo.teamName}`);

    const songElements = await teamElement.$$('.pricing-box.best-price');

    // Initialize an array to store song information for the current team
    const songs = new Map();

    if (!songElements) {
      console.log('No song elements found for this team.');
      continue;
    }

    let songIndex = 1;
    // Iterate through songs within the current team
    for (const songElement of songElements) {
      songElement.scrollIntoView();
      debugger;
      // Song Information and Points Information
      let songName = '';
      try {
        songName = await songElement.$eval('a', (a) => a.innerText.trim());
      } catch (error) {
        const text = await page.$eval('span#notready strong', (strongElement) => {
          return strongElement.textContent;
        });
        if (text === '- NO ENTRY -') {
          console.log('Skipping non entry for team', teamInfo.teamName);
          continue;
        }
      }
      try {
        const genreName = await songElement.$eval('h5', (h5) => h5.innerText.trim());
        const artistName = await songElement.$eval('.textOverflow:nth-child(3)', (textOverflow) => textOverflow.innerText.trim());
        const linkElement = await songElement.$('a');
        const songPageLink = linkElement ? await linkElement.getProperty('href').then(href => href.jsonValue()) : null;

        const pointsElements = await songElement.$$('xpath/ancestor::div[contains(@class, "col-sm-4")]');
        const spans = await pointsElements[0].$$('.bofu_meters span');
        const totalPoints = await spans[0].evaluate(span => span.innerText.replace('Total :', '').replace(' Point', '').trim());
        const medianPoints = await spans[1].evaluate(span => span.innerText.replace('Median :', '').replace(' Points', '').trim());

        const songInfo = {
          songName,
          genreName,
          artistName,
          songPageLink,
          totalPoints,
          medianPoints,
        };



        //BMS labels
        const bmsLabels = await songElement.$eval('.bmsinfo small', (labelElement) => {
          const labels = Array.from(labelElement.querySelectorAll('strong')).map((label) => label.innerText.trim());
          return labels;
        });
        songInfo.bmsLabels = bmsLabels;

        const songImpression = await songElement.$eval('.tleft.textOverflow', (impressionElement) => {
          const impressionCount = impressionElement.querySelector('span').textContent.trim();
          return impressionCount;
        });

        songInfo.songImpression = songImpression;

        const entryNumber = await songElement.$eval('.pricing-action span small', (updateElement) => {
          const entryText = updateElement.innerText.trim();
          const regex = /No.(\d+)/;
          const match = entryText.match(regex);

          let entryString = null;
          if (match) {
            entryString = match[1];
          }

          return { entryString };
        });
        songInfo.entryNumber = entryNumber.entryString;

        const updateInfo = await songElement.$eval('.pricing-action span small', (updateElement) => {
          const updateText = updateElement.innerText.trim();
          const regex = /update : (\d{4}\/\d{2}\/\d{2} \d{2}:\d{2})/;
          const match = updateText.match(regex);

          let updateDateString = null;
          if (match) {
            updateDateString = match[1];
          }

          return { updateDateString };
        });
        songInfo.updateDateTime = new Date(updateInfo.updateDateString);
        songInfo.scrapedDateTime = new Date();

        // Push the extracted song information to the songs array
        songs.set(songInfo.songName, songInfo);
        songIndex += 1;
        // console.log(`Processed Song #${songIndex}: ${songInfo.songName}`);
      } catch (error) {
        console.log('An Error occured:', error);
        console.log(teamInfo.teamName);
        // await page.waitForTimeout(60000);
      }
    }
    teamInfo.songs = songs;

    teams.set(teamInfo.teamName, teamInfo);
    teamIndex += 1;
    teamBar.increment();
  }
  // console.log(teams);


  // Now, you can access songPageLink within the existing songs map
  for (const [teamName, teamInfo] of teams.entries()) {

    // Navigate to the teamPageLink
    await page.goto(teamInfo.teamPageLink);
    const sectionElements = await page.$$('div.col_full.center.bottommargin-lg, div.col_half.center, div.col_half.col_last.center, div.col_full.center.bottommargin-lg, div.col_full.center.bottommargin-lg, div.col_half.center.nobottommargin, div.col_half.col_last.center.nobottommargin, div.post-grid.grid-container.post-masonry.clearfix, div.col_full.center.bottommargin-lg, div.col_one_third.bottommargin-lg.center, div.col_one_third.col_last.bottommargin-lg.center, div.col_full.bottommargin-lg, div.col_full.bottommargin-lg, div.col_half.bottommargin-lg, div.col_half.col_last.bottommargin-lg');

    // ghetto enums cause apparently javascript doesn't have em???
    const LEADER = 0;
    const TWITTER = 1;
    const WEBSITE = 2;
    const CONCEPT = 3;
    const BLANK_WORKS = 4;
    const WORKS = 5;
    const DECLARED = 6;
    const SONGS = 7;
    const BLANK = 8;
    const GENRE = 9;
    const SHARED = 10;
    const REASON = 11;
    const MEMBERS = 12;
    const COMMENT = 13;
    const REGIST = 14;
    const UPDATE = 15

    const leaderSection = await sectionElements[LEADER].$eval('p:nth-of-type(2)', (element) => {
      teamLeader = element.querySelector('big').innerText.trim();
      const teamLeaderCountry = element.querySelector('img').title;
      const textContent = element.textContent.trim();

      const teamLeaderLanguageMatch = textContent.match(/Language : ([^)]+)/);
      const teamLeaderLanguage = teamLeaderLanguageMatch ? teamLeaderLanguageMatch[1].trim() : '';

      return { teamLeader, teamLeaderCountry, teamLeaderLanguage };
    });
    teamInfo.teamLeader = leaderSection.teamLeader;
    teamInfo.teamLeaderCountry = leaderSection.teamLeaderCountry;
    teamInfo.teamLeaderLanguage = leaderSection.teamLeaderLanguage;

    const twitterSection = await sectionElements[TWITTER].$eval('p a', (element) => {
      const twitterLink = element.href;
      return { twitterLink };
    });
    teamInfo.twitterLink = twitterSection.twitterLink;

    const websiteSection = await sectionElements[WEBSITE].$eval('p a', (element) => {
      const websiteLink = element.href;
      return { websiteLink };
    });
    teamInfo.websiteLink = websiteSection.websiteLink;

    const conceptSection = await sectionElements[CONCEPT].$$eval('.col-md-3.center.bottommargin-lg', (elements) => {
      let concepts = [];
      let conceptImage = ''
      let conceptName = ''
      for (const element of elements) {
        conceptImage = element.querySelector('img') ? element.querySelector('img').src : '';
        const conceptName = element.querySelector('h3').textContent.trim();

        concepts.push({ conceptImage, conceptName });
      }
      return { concepts };
    });
    teamInfo.concepts = conceptSection.concepts;

    const worksSection = await sectionElements[WORKS].$eval('.counter', (element) => {
      const works = element.textContent.trim();
      return { works };
    });
    teamInfo.works = worksSection.works;

    const declaredWorksSection = await sectionElements[DECLARED].$eval('.counter', (element) => {
      const declaredWorks = element.textContent.trim();
      return { declaredWorks };
    });
    teamInfo.declaredWorks = declaredWorksSection.declaredWorks;

    const genreSection = await sectionElements[GENRE].$eval('p', (element) => {
      const textContent = element.textContent.trim();
      const genreMatch = textContent.match(/オリジナル \/ ([^)]+)/);
      const genre = genreMatch ? genreMatch[1].trim() : '';
      return { genre };
    });
    teamInfo.genre = genreSection.genre;

    const sharedSection = await sectionElements[SHARED].$eval('p', (element) => {
      const shared = element.textContent.trim();
      return { shared };
    });
    teamInfo.shared = sharedSection.shared;

    const reasonSection = await sectionElements[REASON].$eval('p', (element) => {
      const reason = element.textContent.trim();
      return { reason };
    });
    teamInfo.reason = reasonSection.reason;

    const memberSection = await sectionElements[MEMBERS].$$eval('p', (elements) => {
      const membersRaw = elements[0].textContent.trim();
      const memberCount = elements[1].textContent.trim().match(/[\d]+/);

      const membersProcessed = membersRaw.split(/[\n,/]/).map((member) => member.trim());

      return { membersRaw, memberCount, membersProcessed };
    });
    teamInfo.membersRaw = memberSection.membersRaw;
    teamInfo.memberCount = memberSection.memberCount;
    teamInfo.membersProcessed = memberSection.membersProcessed;

    for (const [songName, songInfo] of teamInfo.songs.entries()) {
      const songPageLink = songInfo.songPageLink;
      // Navigate to songPageLink
      await page.goto(songPageLink);

      // Use Puppeteer to extract the jacket source
      try {
        const jacketImageSrc = await page.$eval('.col_one_third.col_last.moreinfo-header.nobottommargin.hidden-xs.hidden-sm img', (imgElement) => {
          const withoutPeriod = imgElement.getAttribute('src').substring(1);
          return `https://manbow.nothing.sh/event${withoutPeriod}`;
        });
        songInfo.jacketImageSrc = jacketImageSrc
      } catch (error) {
        songInfo.jacketImageSrc = '';
      }

      try {
        // Use Puppeteer to extract the banner source
        const bannerImageElement = await page.$x("//div[contains(@style, 'upload')]/@style");
        if (bannerImageElement.length > 0) {
          const styleAttribute = await bannerImageElement[0].getProperty('textContent');
          // console.log('Style Attribute String:', styleAttribute.toString());
          // Use a regular expression to match URLs containing "upload"
          const uploadUrlMatch = (styleAttribute.toString()).match(/url\("([^"]*upload[^"]*)"\)/);
          // console.log('Upload Url Match 1', uploadUrlMatch[1]);
          songInfo.bannerImageSrc = `https://manbow.nothing.sh/event${uploadUrlMatch[1].substring(1)}`;
        } else {
          songInfo.bannerImageSrc = '';
        }
      } catch (error) {
        songInfo.bannerImageSrc = '';
      }


      // Extract youtube link.
      let youtubeLink = ''; // Initialize to a default value

      const iframeElement = await page.$('div.fluid-width-video-wrapper iframe');
      if (iframeElement) {
        youtubeLink = await page.$eval('div.fluid-width-video-wrapper iframe', (iframe) => {
          return iframe.getAttribute('src');
        });
      }
      songInfo.youtubeLink = youtubeLink;

      // Extract soundcloud link.
      debugger;
      try {
        // Wait for the iframe to load
        await page.waitForSelector('.m_audition iframe');

        // Get the iframe element
        const iframeElement = await page.$('.m_audition iframe');

        // Extract the src attribute of the iframe
        const soundcloudSrc = await page.evaluate(iframe => iframe.src, iframeElement);
        const soundcloudUrlSrc = new URL(soundcloudSrc);
        // Get the value of the 'url' parameter
        const urlParam = soundcloudUrlSrc.searchParams.get('url');
        // Extract the necessary part of the URL
        const soundcloudLink = new URL(urlParam).toString();
        songInfo.soundcloudLink = soundcloudLink;
      } catch (error) {
        songInfo.soundcloudLink = '';
      }


      // Extract only the linkUrls
      const linkUrls = await page.$$eval('blockquote p a', (elements) => {
        return elements.map((element) => element.getAttribute('href'));
      });

      // console.log('Link URLs:', linkUrls);

      // Extract all text within the <p> element separated by <br> tags
      const paragraphTexts = await page.$eval('p[style="font-size:75%"]', (element) => {
        const textWithEntities = element.innerHTML.split('<br>').map((text) => text.trim());

        // Define a mapping of character references to their corresponding characters
        const characterReferences = {
          '&lt;': '<',
          '&gt;': '>',
          '&quot;': '"',
          '&apos;': "'",
          '&amp;': '&',
          // Add more character references here as needed
        };

        // Replace character references in the text
        const textWithoutEntities = textWithEntities.map((text) => {
          for (const entity in characterReferences) {
            if (text.includes(entity)) {
              text = text.replace(new RegExp(entity, 'g'), characterReferences[entity]);
            }
          }
          return text;
        });

        return textWithoutEntities;
      });


      // console.log('Paragraph Text:', paragraphTexts);

      // Initialize the links array
      let links = [];

      // Handle inline link descriptions
      let inlineUrlDescs = [];
      for (const paragraphText of paragraphTexts) {
        let linkElement = {
          linkUrl: '',
          linkDesc: '',
        };

        for (const linkUrl of linkUrls) {
          if (paragraphText.includes(linkUrl)) {
            // Create a regular expression pattern to match the link pattern
            const linkPattern = new RegExp(`<a(.*?)</a>`, 'g');

            // Replace the link pattern with an empty string to remove it
            linkElement.linkDesc = paragraphText.replace(linkPattern, '');
            linkElement.linkUrl = linkUrl;
            break;
          }
        }
        if (linkElement.linkUrl == '') {
          linkElement.linkDesc = paragraphText;
        }

        if (linkElement.linkUrl !== '' || linkElement.linkDesc !== '') { // prevent blank linkElements
          inlineUrlDescs.push(linkElement)
        }
      }

      // handle link descriptions above the link
      const aboveUrlDescs = [];
      try {

        for (let i = 0; i < inlineUrlDescs.length; i++) {
          // match above descriptions to a link directly below

          if (
            i == 0 &&
            inlineUrlDescs[i].linkUrl === '' &&
            inlineUrlDescs[i].linkDesc !== '' &&
            inlineUrlDescs[i + 1].linkUrl !== '' &&
            inlineUrlDescs[i + 1].linkDesc === ''
          ) {
            const newUrl = inlineUrlDescs[i + 1].linkUrl;
            const newDesc = inlineUrlDescs[i].linkDesc;
            i++; // Increment i to skip the next element in the original array
            aboveUrlDescs.push({ linkUrl: newUrl, linkDesc: newDesc });
          } else if (
            i < inlineUrlDescs.length - 1 &&
            inlineUrlDescs[i].linkUrl === '' &&
            inlineUrlDescs[i].linkDesc !== '' &&
            inlineUrlDescs[i + 1].linkUrl !== '' &&
            inlineUrlDescs[i + 1].linkDesc === '' &&
            inlineUrlDescs[i - 1].linkUrl !== '' &&
            inlineUrlDescs[i - 1].linkDesc === ''
          ) {
            const newUrl = inlineUrlDescs[i + 1].linkUrl;
            const newDesc = inlineUrlDescs[i].linkDesc;
            i++; // Increment i to skip the next element in the original array
            aboveUrlDescs.push({ linkUrl: newUrl, linkDesc: newDesc });
          } else if (inlineUrlDescs[i].linkUrl !== '' || inlineUrlDescs[i].linkDesc !== '') {
            aboveUrlDescs.push(inlineUrlDescs[i]); // Keep the current element
          }
        }
      } catch (error) {
        console.log('No Link Found: ', songPageLink)
        // technically this could apply whatever text is there as a description with no url, but i don't have the patience for it atm
      }

      // handle multiline descs
      const multilineUrlDescs = [];
      let pendingUrl = '';
      let pendingDesc = '';

      for (const { linkUrl, linkDesc } of aboveUrlDescs) {
        // debugger;
        if (linkUrl) {
          if (pendingUrl) {
            multilineUrlDescs.push({ linkUrl: pendingUrl, linkDesc: pendingDesc });
            pendingUrl = '';
            pendingDesc = '';
          } else if (pendingDesc) {
            if (linkDesc) {
              multilineUrlDescs.push({ linkUrl: '', linkDesc: pendingDesc });
              pendingUrl = '';
              pendingDesc = '';
            } else {
              multilineUrlDescs.push({ linkUrl, linkDesc: pendingDesc });
              pendingUrl = '';
              pendingDesc = '';
              continue;
            }
          }
          if (linkDesc) {
            multilineUrlDescs.push({ linkUrl, linkDesc });
          }
        } else if (linkDesc) {
          pendingDesc = pendingDesc ? `${pendingDesc}\n${linkDesc}` : linkDesc;
        }
        if (linkUrl && !linkDesc) {
          multilineUrlDescs.push({ linkUrl, linkDesc });
        }
      }

      // Handle any pending items
      if (pendingUrl) {
        multilineUrlDescs.push({ linkUrl: pendingUrl, linkDesc: pendingDesc });
      } else if (pendingDesc) {
        multilineUrlDescs.push({ linkUrl: '', linkDesc: pendingDesc });
      }

      // Replace the original 'links' array with the modified 'newLinks' array
      links = multilineUrlDescs;
      songInfo.links = links;
      // console.log('Link Descriptions:', linkDescs);

      let bemuseLink = '';
      try {
        // Use Puppeteer to extract the Bemuse link
        bemuseLink = await page.$eval('.bmson-iframe-content iframe', (iframe) => {
          return iframe.getAttribute('src');
        });
        songInfo.bemuseLink = bemuseLink;

      } catch (error) {
        songInfo.bemuseLink = '';
      }
    }

  }
  // console.dir(teams, { depth: null });

  // write to file

  // Convert Map to JSON-friendly structure
  const teamsObject = {};

  teams.forEach((teamDetails, key) => {
    const teamDetailsObject = { ...teamDetails }; // Copy main team details

    // Handle nested songs Map
    if (teamDetails.songs instanceof Map) {
      teamDetailsObject.songs = Array.from(teamDetails.songs).reduce((acc, [songKey, songDetails]) => {
        acc.push({ songKey, ...songDetails });
        return acc;
      }, []);
    }

    teamsObject[key] = teamDetailsObject;
  });

  // Write teams data to a JSON file
  writeFile('./data/bof142.json', JSON.stringify(teamsObject, null, 2), (err) => {
    if (err) {
      console.error('Error writing to file', err);
    } else {
      console.log('Successfully wrote to file');
    }
  });




  // await page.waitForTimeout(120000);
  await browser.close();
  multibar.stop();
  console.timeEnd('programExecution'); // End the timer
})();