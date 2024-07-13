import puppeteer from 'puppeteer';
import { writeFile } from 'fs';
import cliProgress from 'cli-progress';

const waitForTimer = 500;
const eventPageActions = 9;
const songElementActions = 17;
const teamPageActions = 32;
const songPageActions = 15;
(async () => {
  console.time('programExecution');
  // Create a new progress bar instance and use shades_classic theme
  const multibar = new cliProgress.MultiBar({
    clearOnComplete: false,
    hideCursor: true,
    format: ' {bar} {percentage}% | {filename} | {value}/{total}',
  }, cliProgress.Presets.shades_classic);
  const browser = await puppeteer.launch({
    headless: "new"
    /*
    headless: false,
    devtools: true
    */
  });

  const page = await browser.newPage();

  // await page.setViewport({ width: 1920, height: 1080 });
  // await page.setViewport({ width: 1080, height: 1920 });

  page.setDefaultTimeout(8000);

  await page.goto('https://manbow.nothing.sh/event/event.cgi?action=List_def&event=142#186');

  const teams = new Map();

  const teamElements = await page.$$('.team_information');
  // console.log('Team Count:', teamElements.length);
  const numberOfTeamsToLimit = 2;
  // Slice the teamElements array to select a specific number of teams
  const limitedTeamElements = teamElements.slice(0, numberOfTeamsToLimit);
  // TODO remember to change this back to teamElements.length
  const eventPageTotal = limitedTeamElements.length * eventPageActions
  const eventPageBar = multibar.create(eventPageTotal, 0)

  let teamIndex = 1;
  // for (const teamElement of teamElements.slice(0, -1)) { // for whatever reason the last element is just empty
  // Iterate through the limited teams
  for (const teamElement of limitedTeamElements) {
    const teamInfo = await teamElement.$eval('.fancy-title :is(h2, h3) a', (link) => {
      const teamName = link.innerText.trim();
      const bannerImageSrc = link.querySelector('img') ? link.querySelector('img').erc : '';
      const teamPageLink = link.href;
      return { teamName, bannerImageSrc, teamPageLink };
    });
    eventPageBar.update({ filename: `Event Page Team #${teamIndex}: ${teamInfo.teamName}` });
    eventPageBar.increment();

    const emblemImageSrc = await teamElement.$eval('.header_emblem', (emblemElement) => {
      const dataBg = emblemElement.getAttribute('data-bg');
      const withoutPeriod = dataBg.substring(1);
      return withoutPeriod.length > 1 ? `https://manbow.nothing.sh/event${withoutPeriod}` : '';
    });
    eventPageBar.increment();
    teamInfo.emblemImageSrc = emblemImageSrc;
    eventPageBar.increment();
    teamInfo.teamImpression = await teamElement.$eval('#team_imp', (element) => element.innerText.trim());
    eventPageBar.increment();
    teamInfo.teamTotal = await teamElement.$eval('#team_total', (element) => element.innerText.trim());
    eventPageBar.increment();
    teamInfo.teamMedian = await teamElement.$eval('#team_med', (element) => element.innerText.trim());
    eventPageBar.increment();


    // console.log(`Processed team #${teamIndex}: ${teamInfo.teamName}`);

    const songElements = await teamElement.$$('.pricing-box.best-price');
    eventPageBar.increment();

    // Initialize an array to store song information for the current team
    const songs = new Map();

    if (!songElements) {
      console.log('No song elements found for this team.');
      continue;
    }

    let songIndex = 1;
    // Iterate through songs within the current team
    const songElementTotal = songElements.length * songElementActions
    const songElementBar = multibar.create(songElementTotal, 0);
    for (const songElement of songElements) {
      // songElement.scrollIntoView();
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
      songElementBar.update({ filename: `Song Element Team #${teamIndex}: ${teamInfo.teamName}: Song #${songIndex}: ${songName}` });
      songElementBar.increment();
      try {
        const genreName = await songElement.$eval('h5', (h5) => h5.innerText.trim());
        songElementBar.increment();
        const artistName = await songElement.$eval('.textOverflow:nth-child(3)', (textOverflow) => textOverflow.innerText.trim());
        songElementBar.increment();
        const linkElement = await songElement.$('a');
        const songPageLink = linkElement ? await linkElement.getProperty('href').then(href => href.jsonValue()) : null;
        songElementBar.increment();

        const pointsElements = await songElement.$$('xpath/ancestor::div[contains(@class, "col-sm-4")]');
        songElementBar.increment();
        const spans = await pointsElements[0].$$('.bofu_meters span');
        const totalPoints = await spans[0].evaluate(span => span.innerText.replace('Total :', '').replace(' Point', '').trim());
        songElementBar.increment();
        const medianPoints = await spans[1].evaluate(span => span.innerText.replace('Median :', '').replace(' Points', '').trim());
        songElementBar.increment();

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
        songElementBar.increment();
        songInfo.bmsLabels = bmsLabels;
        songElementBar.increment();

        const songImpression = await songElement.$eval('.tleft.textOverflow', (impressionElement) => {
          const impressionCount = impressionElement.querySelector('span').textContent.trim();
          return impressionCount;
        });
        songElementBar.increment();

        songInfo.songImpression = songImpression;
        songElementBar.increment();

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
        songElementBar.increment();
        songInfo.entryNumber = entryNumber.entryString;
        songElementBar.increment();

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
        songElementBar.increment();
        songInfo.updateDateTime = new Date(updateInfo.updateDateString);
        songElementBar.increment();
        songInfo.scrapedDateTime = new Date();
        songElementBar.increment();

        // Push the extracted song information to the songs array
        songs.set(songInfo.songName, songInfo);
        songElementBar.increment();
        songIndex += 1;
        // console.log(`Processed Song #${songIndex}: ${songInfo.songName}`);
      } catch (error) {
        console.log('An Error occured:', error);
        console.log(teamInfo.teamName);
      }
    }
    teamInfo.songs = songs;
    eventPageBar.increment();

    teams.set(teamInfo.teamName, teamInfo);
    eventPageBar.increment();
    teamIndex += 1;
  }
  // console.log(teams);


  const teamPageTotal = teams.size * teamPageActions;
  const teamPageBar = multibar.create(teamPageTotal, 0);
  let songPageTotal = 0;
  const songPageBar = multibar.create(100, 0);
  // Now, you can access songPageLink within the existing songs map
  for (const [teamName, teamInfo] of teams.entries()) {

    songPageTotal = songPageTotal + (teamInfo.songs.size * songPageActions);
    songPageBar.setTotal(songPageTotal);
    // Navigate to the teamPageLink
    await page.goto(teamInfo.teamPageLink);
    teamPageBar.update({ filename: `Team Page: ${teamName}` });
    teamPageBar.increment();

    const sectionElements = await page.$$('div.col_full.center.bottommargin-lg, div.col_half.center, div.col_half.col_last.center, div.col_full.center.bottommargin-lg, div.col_full.center.bottommargin-lg, div.col_half.center.nobottommargin, div.col_half.col_last.center.nobottommargin, div.post-grid.grid-container.post-masonry.clearfix, div.col_full.center.bottommargin-lg, div.col_one_third.bottommargin-lg.center, div.col_one_third.col_last.bottommargin-lg.center, div.col_full.bottommargin-lg, div.col_full.bottommargin-lg, div.col_half.bottommargin-lg, div.col_half.col_last.bottommargin-lg');
    teamPageBar.increment();

    // ghetto enums cause apparently javascript doesn't have em???
    const LEADER = 0;
    const TWITTER = 1;
    const WEBSITE = 2;
    const CONCEPT = 3;
    // const BLANK_WORKS = 4;
    const WORKS = 5;
    const DECLARED = 6;
    // const SONGS = 7;
    // const BLANK = 8;
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
    teamPageBar.increment();
    teamInfo.teamLeader = leaderSection.teamLeader;
    teamPageBar.increment();
    teamInfo.teamLeaderCountry = leaderSection.teamLeaderCountry;
    teamPageBar.increment();
    teamInfo.teamLeaderLanguage = leaderSection.teamLeaderLanguage;
    teamPageBar.increment();

    const twitterSection = await sectionElements[TWITTER].$eval('p a', (element) => {
      const twitterLink = element.href;
      return { twitterLink };
    });
    teamPageBar.increment();
    teamInfo.twitterLink = twitterSection.twitterLink;
    teamPageBar.increment();

    const websiteSection = await sectionElements[WEBSITE].$eval('p a', (element) => {
      const websiteLink = element.href;
      return { websiteLink };
    });
    teamPageBar.increment();
    teamInfo.websiteLink = websiteSection.websiteLink;
    teamPageBar.increment();

    const conceptSection = await sectionElements[CONCEPT].$$eval('.col-md-3.center.bottommargin-lg', (elements) => {
      let concepts = [];

      for (const element of elements) {
        const conceptImage = element.querySelector('img') ? element.querySelector('img').src : '';
        const conceptName = element.querySelector('h3').textContent.trim();

        concepts.push({ conceptImage, conceptName });
      }
      return { concepts };
    });
    teamPageBar.increment();
    teamInfo.concepts = conceptSection.concepts;
    teamPageBar.increment();

    const worksSection = await sectionElements[WORKS].$eval('.counter', (element) => {
      const works = element.textContent.trim();
      return { works };
    });
    teamPageBar.increment();
    teamInfo.works = worksSection.works;
    teamPageBar.increment();

    const declaredWorksSection = await sectionElements[DECLARED].$eval('.counter', (element) => {
      const declaredWorks = element.textContent.trim();
      return { declaredWorks };
    });
    teamPageBar.increment();
    teamInfo.declaredWorks = declaredWorksSection.declaredWorks;
    teamPageBar.increment();

    const genreSection = await sectionElements[GENRE].$eval('p', (element) => {
      const textContent = element.textContent.trim();
      const genreMatch = textContent.match(/オリジナル \/ ([^)]+)/);
      const genre = genreMatch ? genreMatch[1].trim() : '';
      return { genre };
    });
    teamPageBar.increment();
    teamInfo.genre = genreSection.genre;
    teamPageBar.increment();

    const sharedSection = await sectionElements[SHARED].$eval('p', (element) => {
      const shared = element.textContent.trim();
      return { shared };
    });
    teamPageBar.increment();
    teamInfo.shared = sharedSection.shared;
    teamPageBar.increment();

    const reasonSection = await sectionElements[REASON].$eval('p', (element) => {
      const reason = element.textContent.trim();
      return { reason };
    });
    teamPageBar.increment();
    teamInfo.reason = reasonSection.reason;
    teamPageBar.increment();

    const memberSection = await sectionElements[MEMBERS].$$eval('p', (elements) => {
      const membersRaw = elements[0].textContent.trim();
      const memberCount = elements[1].textContent.trim().match(/[\d]+/);
      const membersProcessed = membersRaw.split(/[\n,/]/).map((member) => member.trim());
      return { membersRaw, memberCount, membersProcessed };
    });
    teamPageBar.increment();
    teamInfo.membersRaw = memberSection.membersRaw;
    teamPageBar.increment();
    teamInfo.memberCount = memberSection.memberCount;
    teamPageBar.increment();
    teamInfo.membersProcessed = memberSection.membersProcessed;
    teamPageBar.increment();

    const commentSection = await sectionElements[COMMENT].$eval('p', (element) => {
      const comment = element.textContent;
      return { comment };
    });
    teamPageBar.increment();
    teamInfo.comment = commentSection.comment;
    teamPageBar.increment();

    const registSection = await sectionElements[REGIST].$eval('strong', (element) => {
      const regist = element.textContent.trim();
      return { regist };
    });
    teamPageBar.increment();
    teamInfo.regist = new Date(registSection.regist);
    teamPageBar.increment();

    const updateSection = await sectionElements[UPDATE].$eval('strong', (element) => {
      const update = element.textContent.trim();
      return { update };
    });
    teamPageBar.increment();
    teamInfo.update = new Date(updateSection.update);
    teamPageBar.increment();





    // song Page Scraping
    for (const [songName, songInfo] of teamInfo.songs.entries()) {
      const songPageLink = songInfo.songPageLink;
      songPageBar.update({ filename: `Song Page: ${songName}` });
      // Navigate to songPageLink
      await page.goto(songPageLink);
      songPageBar.increment();

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
      songPageBar.increment();

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
      songPageBar.increment();


      // Extract youtube link.
      let youtubeLink = ''; // Initialize to a default value

      const iframeElement = await page.$('div.fluid-width-video-wrapper iframe');
      if (iframeElement) {
        youtubeLink = await page.$eval('div.fluid-width-video-wrapper iframe', (iframe) => {
          return iframe.getAttribute('src');
        });
      }
      songPageBar.increment();
      songInfo.youtubeLink = youtubeLink;
      songPageBar.increment();

      // Extract soundcloud link.
      debugger;
      try {
        // Wait for the iframe to load
        await page.waitForSelector('.m_audition iframe', {timeout: waitForTimer});

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
      songPageBar.increment();


      // Extract only the linkUrls
      const linkUrls = await page.$$eval('blockquote p a', (elements) => {
        return elements.map((element) => element.getAttribute('href'));
      });
      songPageBar.increment();

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
      songPageBar.increment();


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
      songPageBar.increment();

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
      songPageBar.increment();

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
      songPageBar.increment();

      // Handle any pending items
      if (pendingUrl) {
        multilineUrlDescs.push({ linkUrl: pendingUrl, linkDesc: pendingDesc });
      } else if (pendingDesc) {
        multilineUrlDescs.push({ linkUrl: '', linkDesc: pendingDesc });
      }
      songPageBar.increment();

      // Replace the original 'links' array with the modified 'newLinks' array
      links = multilineUrlDescs;
      songPageBar.increment();
      songInfo.links = links;
      songPageBar.increment();
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
      songPageBar.increment();
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
      // console.log('Successfully wrote to file');
    }
  });




  // await page.waitForTimeout(120000);
  await browser.close();
  multibar.stop();
  console.timeEnd('programExecution'); // End the timer
})();