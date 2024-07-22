import puppeteer from 'puppeteer';
import { writeFile } from 'fs';
import config from './puppeteerConfig.js'
import cliProgress from 'cli-progress';
import bunyan from 'bunyan';
import path from 'path'
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// $env:MODE='test'
// $env:MODE='production'
// TODO rehaul progress bar to track every action in a single function
// TODO revamp try-catch for error handling
// TODO use logging module to log errors to a file
// TODO scrape in parallel for better performance
// TODO rehaul program for main event list entry point, and scraping of all data for each event
// TODO seperate scraping by page, including larger event vs smaller event page types
// TODO collect regex matches for reusability
// TODO properly use StringJpDate on other arrays

// Determine if we are in test mode
const isTest = process.env.MODE === 'test';

let teams = new Map();
// Array to collect logs during script execution
let collectedLogs = [];


// Define the directory for logs
const logsDirectory = path.join(__dirname, 'logs');

// Define the log file path
const logFilePath = isTest
  ? path.join(logsDirectory, 'testEventLog.json')
  : path.join(logsDirectory, `event${config.latestEvent.eventKey}Log.json`);

const bunyanLogger = bunyan.createLogger({
  name: `event${config.latestEvent.fullName}`,  // Give your logger a name
  streams: [
    {
      type: 'file',
      // type: 'rotating-file',  // To handle log rotation
      path: logFilePath,      // Specify the log file path
      // period: '1d',           // Daily rotation
      // count: 3                // Keep 3 back copies
    }
  ],
  // Use Bunyan's serializers to ensure proper JSON formatting
  serializers: bunyan.stdSerializers,
  // Use a custom serializer to format log records as JSON with newline separators
  serializers: {
    bunyanRecord: function (rec) {
      return JSON.stringify(rec) + '\n';
    }
  }
});


// Function to log with line number
function log(level, message) {
  const err = new Error();
  // Extract file name and line number
  const stackLines = err.stack.split('\n').slice(1, -1);//.slice(2).map(line => line.trim());

  const logObject = {
    level: level,
    message: message,
    stack: stackLines,
    timestamp: new Date().toISOString()
  };

  // Push log object to collected logs array
  collectedLogs.push(logObject);

  // Log to Bunyan
  bunyanLogger[level](logObject);
}



function convertJpDate(dateString) {

  // Extract the year, month, day, hour, and minute using a regular expression
  const datePattern = /(\d{4})年(\d{1,2})月(\d{1,2})日\s(\d{1,2}):(\d{2})/;
  const match = dateString.match(datePattern);

  if (match) {
    const year = match[1];
    const month = String(parseInt(match[2], 10)).padStart(2, '0'); // Ensure month is 2 digits
    const day = String(parseInt(match[3], 10)).padStart(2, '0'); // Ensure day is 2 digits
    const hour = String(parseInt(match[4], 10)).padStart(2, '0'); // Ensure hour is 2 digits
    const minute = String(parseInt(match[5], 10)).padStart(2, '0'); // Ensure minute is 2 digits

    return `${year}-${month}-${day} ${hour}:${minute}`;
  }
  return '';
}


const StringConvertJpDate = convertJpDate.toString();

async function headerExists(elementHandle, headerXPath) {

  if (!elementHandle || !(await elementHandle.evaluate) || typeof (await elementHandle.evaluate) !== 'function') {
    // console.log(`Header element not found or 'evaluate' function missing for XPath: ${headerXPath}`);
    return false;
  }
  // Check if the handle represents an actual element
  const element = await elementHandle.asElement();
  if (!element) {
    // console.log(`Header element not found for XPath: ${headerXPath}`);
    return false;
  }
  // Check if the element has the required classes
  const hasRequiredClasses = await element.evaluate(el => {
    return el.classList.contains('spost') &&
      el.classList.contains('clearfix') &&
      el.classList.contains('nomarginbottom');
  });

  if (!hasRequiredClasses) {
    log('error', `Header element does not contain the required classes: ${headerXPath}`);
    // console.log(`Missing required Classes on Xpath: ${headerXPath} Classes: ${element.classList}`);
    return false;
  }

  // Check if the element's innerHTML is only &nbsp;
  const innerHTML = await element.evaluate(el => el.innerHTML.trim());
  if (!innerHTML) {
    return false;
  }
  if (!innerHTML.includes('<div')) {
    log('error', `Header element's innerHTML does not contain any <div> tags: ${headerXPath}`);
    // console.log(`InnerHTML didn't have <div> on Xpath: ${headerXPath} InnerHTML: ${innerHTML.trim()}`);
    return false;
  }
  // // console.log(`\n\n\nInnerHTML: ${innerHTML.trim()}\nXPath: ${headerXPath} considered a valid header element`);

  // console.log('\n');
  // console.log(`XPath considered valid: ${headerXPath}`);
  return true;
}

async function scrapeLongImpressions(page, headerXPathBase) {
  const longImpressions = [];
  let headerXPath = headerXPathBase;
  let XPathOffset = 0;
  let isReply = false;
  const headerOffset = 5;
  const replyOffset = 2;
  let headerElementHandle;
  let shouldLoop = true;
  let responseButton;

  while (shouldLoop) {
    // console.log(`\nXpath from top of loop: ${headerXPath}`);
    headerElementHandle = await page.evaluateHandle((XPath) => {
      const headerElement = document.evaluate(XPath, document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue;
      return headerElement;
    }, headerXPath);

    if (!(await headerExists(headerElementHandle, headerXPath))) {
      log('error', `Header element not found or is undefined for XPath: ${headerXPath}`);
      break;
    }


    // console.dir(headerElementHandle, { depth: null });

    // Extract pointsOverall from the <nobr> element
    const pointsOverallElementHandle = await headerElementHandle.$('nobr');
    const pointsOverall = pointsOverallElementHandle ? await (await pointsOverallElementHandle.getProperty('textContent')).jsonValue() : '';

    // Extract other elements similarly
    const topHeaderElementHandle = await headerElementHandle.$('div.entry-c div.entry-title');
    const topHeaderRawString = topHeaderElementHandle ? await (await topHeaderElementHandle.getProperty('textContent')).jsonValue() : '';
    const userNameRegex = /[\s]+([^\n]+?)[\s]*\n\t+/;
    const userNameMatch = topHeaderRawString.match(userNameRegex);
    const userName = userNameMatch ? userNameMatch[1].trim() : '';

    const pointsRegex = /([\w\s]+)\s*:\s*(\d+)\s*Pts\./g;
    const pointsRawString = topHeaderRawString.replace(userNameRegex, '');
    const pointBreakdownArray = [];
    let pointBreakdownMatch;

    while ((pointBreakdownMatch = pointsRegex.exec(pointsRawString)) !== null) {
      pointBreakdownArray.push({
        pointName: pointBreakdownMatch[1].trim(),
        pointValue: Number(pointBreakdownMatch[2])
      });
    }

    const countryCodeElementHandle = await headerElementHandle.$('div.spost.clearfix.nomarginbottom div.entry-c div.entry-title img.flag');
    const countryCode = countryCodeElementHandle ? await (await countryCodeElementHandle.getProperty('title')).jsonValue() : '';

    const countryFlagElementHandle = await headerElementHandle.$('div.spost.clearfix.nomarginbottom div.entry-c div.entry-title img');
    const countryFlag = countryFlagElementHandle ? await (await countryFlagElementHandle.getProperty('src')).jsonValue() : '';

    const jpDateTimeElementHandle = await headerElementHandle.$('div.spost.clearfix.nomarginbottom ul li');
    const jpDateTimeString = jpDateTimeElementHandle ? await (await jpDateTimeElementHandle.getProperty('textContent')).jsonValue() : null;
    const jpDateTime = jpDateTimeString ? convertJpDate(jpDateTimeString) : '';

    let responseButtonLinkHandle = await page.evaluateHandle((XPath) => {
      let responseButtonElement = document.evaluate(`${XPath}/following-sibling::div[3]/p/a`, document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue;

      if (responseButtonElement) {
      } else {
        responseButtonElement = document.evaluate(`${XPath}/following-sibling::div[5]/p/a`, document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue;
      }

      return responseButtonElement ? responseButtonElement.href : null;
    }, headerXPath);

    responseButton = responseButtonLinkHandle ? await responseButtonLinkHandle.jsonValue() : null;

    const commentSectionHandle = await page.evaluateHandle((XPath) => {
      const commentSectionElement = document.evaluate(`${XPath}/following-sibling::div[1]/p`, document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue;
      return commentSectionElement ? commentSectionElement.innerHTML : null;
    }, headerXPath);



    longImpressions.push({
      pointsOverall: Number(pointsOverall),
      userName: userName,
      countryCode: countryCode,
      countryFlag: countryFlag,
      pointBreakdown: pointBreakdownArray,
      jpDateTime: new Date(jpDateTime),
      commentSection: commentSectionHandle ? await commentSectionHandle.jsonValue() : null,
      responseButton: responseButton,
      isReply: isReply
    });

    // first check if there is a reply
    let testXPathOffset = XPathOffset + replyOffset;
    let testElementHandle = await page.evaluateHandle((XPath) => {
      const headerElement = document.evaluate(XPath, document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue;
      return headerElement;
    }, `${headerXPathBase}/following-sibling::div[${testXPathOffset}]`);

    if (await headerExists(testElementHandle, `${headerXPathBase}/following-sibling::div[${testXPathOffset}]`)) {
      // handle case where the current reply genuinely lack a response so the next impression is where a reply would be expected
      // isReply = longImpressions[longImpressions.length - 1].responseButton && responseButton ? true : false;
      isReply = pointBreakdownArray.length > 0 ? true : false;
    } else {
      // next check if there is another impression
      testXPathOffset = XPathOffset + headerOffset;
      testElementHandle = await page.evaluateHandle((XPath) => {
        const headerElement = document.evaluate(XPath, document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue;
        return headerElement;
      }, `${headerXPathBase}/following-sibling::div[${testXPathOffset}]`);

      if (await headerExists(testElementHandle, `${headerXPathBase}/following-sibling::div[${testXPathOffset}]`)) {
        isReply = false;
      } else {
        shouldLoop = false;
        break;
      }
    }
    XPathOffset = testXPathOffset;
    headerXPath = `${headerXPathBase}/following-sibling::div[${XPathOffset}]`; // Update headerXPath for next iteration

    // console.log(`Xpath from bottom of loop: ${headerXPath}`);
  }
  return longImpressions;
}

async function saveData() {
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
  const latestEventObject = {
    [config.latestEvent.eventKey]: {
      fullName: config.latestEvent.fullName,
      shortName: config.latestEvent.shortName,
      subTitle: config.latestEvent.subTitle,
      lastScrapeTime: new Date(),
      teams: teamsObject
    }
  };
  // Write teams data to a JSON file
  writeFile(`./data/event${config.latestEvent.eventKey}.json`, JSON.stringify(latestEventObject, null, 2), (err) => {
    if (err) {
      log('error', `Error writing to file: ${err}`);
    } else {
      log('info', 'Successfully wrote to file');
    }
  });

  // Write collected logs to file as a single JSON array
  writeFile(logFilePath, JSON.stringify(collectedLogs, null, 2), (err) => {
    if (err) {
      // console.log(`Error writing to logFile: ${err}`);
    }
  });

}

(async () => {
  console.time('programExecution');

  // TODO toggle in config
  // Create a new progress bar instance and use shades_classic theme
  const multibar = new cliProgress.MultiBar({
    clearOnComplete: false,
    hideCursor: true,
    format: ' {bar} {percentage}% | {filename}',
  }, cliProgress.Presets.shades_classic);

  const browser = await puppeteer.launch({
    headless: "new"
    /*
    headless: false,
    devtools: true
    */
  });

  const page = await browser.newPage();

  page.setDefaultTimeout(config.navigationTimeout);

  // TODO scrape media from event information page
  await page.goto(`https://manbow.nothing.sh/event/event.cgi?action=List_def&event=${config.latestEvent.eventKey}`);


  const allTeamElements = await page.$$('.team_information');

  // Apply the limit if in test mode
  const teamElements = config.numberOfTeamsToLimit
    ? allTeamElements.slice(0, config.numberOfTeamsToLimit)
    : allTeamElements.slice(0, -1); // for whatever reason the last element is just empty

  const eventPageTotal = teamElements.length * config.actions.eventPage
  const eventPageBar = multibar.create(eventPageTotal, 0)

  eventPageBar.update({ filename: `Event: ${config.latestEvent.fullName}` });

  let teamIndex = 1;
  for (const teamElement of teamElements) {
    const teamInfo = await teamElement.$eval('.fancy-title :is(h2, h3) a', (link) => {
      const teamName = link.innerText.trim();
      const bannerImageSrc = link.querySelector('img') ? link.querySelector('img').src : '';
      const teamPageLink = link.href;
      return { teamName, bannerImageSrc, teamPageLink };
    });
    eventPageBar.increment();

    const emblemImageSrc = await teamElement.$eval('.header_emblem', (emblemElement) => {
      const dataBg = emblemElement.getAttribute('data-bg');
      const withoutPeriod = dataBg.substring(1);
      return withoutPeriod.length > 1 ? `https://manbow.nothing.sh/event${withoutPeriod}` : '';
    });
    eventPageBar.increment();
    teamInfo.emblemImageSrc = emblemImageSrc;
    eventPageBar.increment();
    teamInfo.teamImpression = Number(await teamElement.$eval('#team_imp', (element) => element.innerText.trim()));
    eventPageBar.increment();
    teamInfo.teamTotal = Number(await teamElement.$eval('#team_total', (element) => element.innerText.trim()));
    eventPageBar.increment();
    teamInfo.teamMedian = Number(await teamElement.$eval('#team_med', (element) => element.innerText.trim()));
    eventPageBar.increment();


    // // console.log(`Processed team #${teamIndex}: ${teamInfo.teamName}`);

    const songElements = await teamElement.$$('.pricing-box.best-price');
    eventPageBar.increment();

    // Initialize an array to store song information for the current team
    const songs = new Map();

    if (!songElements) {
      // console.log('No song elements found for this team.');
      continue;
    }

    let songIndex = 1;
    // Iterate through songs within the current team
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
          // // console.log('Skipping non entry for team', teamInfo.teamName);
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
        const totalPoints = Number(await spans[0].evaluate(span => span.innerText.replace('Total :', '').replace(' Point', '').trim()));

        const medianPoints = Number(await spans[1].evaluate(span => span.innerText.replace('Median :', '').replace(' Points', '').trim()));


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


        // TODO scrape impression stars? maybe even more for impression visualizer

        const entryCompositionUpdateElement = await songElement.$eval('.pricing-action span small', (element) => {
          const text = element.innerText.trim();
          const entryRegex = /No.(\d+)/;
          const compositionRegex = /(Original|Copy|Arrange|Remix)/
          const updateRegex = /update : (\d{4}\/\d{2}\/\d{2} \d{2}:\d{2})/;

          const entryMatch = text.match(entryRegex);

          let entryString = null;
          if (entryMatch) {
            entryString = entryMatch[1];
          }

          const compositionMatch = text.match(compositionRegex);
          let compositionString = null;
          if (compositionMatch) {
            compositionString = compositionMatch[1];
          }

          const updateMatch = text.match(updateRegex);

          let updateDateString = null;
          if (updateMatch) {
            updateDateString = updateMatch[1];
          }

          return { entryString, compositionString, updateDateString };
        });

        songInfo.entryNumber = Number(entryCompositionUpdateElement.entryString);

        songInfo.compositionType = entryCompositionUpdateElement.compositionString;

        songInfo.updateDateTime = new Date(entryCompositionUpdateElement.updateDateString);

        songInfo.scrapedDateTime = new Date();


        // Push the extracted song information to the songs array
        songs.set(songInfo.songName, songInfo);

        songIndex += 1;
        // // console.log(`Processed Song #${songIndex}: ${songInfo.songName}`);
      } catch (error) {
        // console.log('An Error occured:', error);
        // console.log(teamInfo.teamName);
      }
    }
    teamInfo.songs = songs;
    eventPageBar.increment();

    const teamKey = teamInfo.teamPageLink.match(/team=([\d]+)/)[1].trim();
    teams.set(teamKey, teamInfo);
    eventPageBar.increment();
    teamIndex += 1;
  }
  // // console.log(teams);


  const teamPageTotal = teams.size * config.actions.teamPage;
  const teamPageBar = multibar.create(teamPageTotal, 0);
  // Assuming 'teams' is your Map object containing team information
  let totalSongs = 0;

  // Iterate over each entry in the 'teams' map
  for (const [teamName, teamInfo] of teams.entries()) {
    // Access the songs array from teamInfo and get its length
    const songCount = teamInfo.songs.size;

    // Add the song count to the total
    totalSongs += songCount;

  }

  // Now 'totalSongs' contains the count of all songs across all teams
  // // console.log(`Total songs across all teams: ${totalSongs}`);


  const songPageTotal = totalSongs * config.actions.songPage;
  const songPageBar = multibar.create(songPageTotal, 0);
  // Now, you can access songPageLink within the existing songs map
  for (const [teamName, teamInfo] of teams.entries()) {

    // Navigate to the teamPageLink
    await page.goto(teamInfo.teamPageLink);
    teamPageBar.update({ filename: `Teams: ${teamInfo.teamName}` });
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
      const teamLeaderCountryCode = element.querySelector('img').title;
      const teamLeaderCountryFlag = element.querySelector('img').src;
      const textContent = element.textContent.trim();

      const teamLeaderLanguageMatch = textContent.match(/Language : ([^)]+)/);
      const teamLeaderLanguage = teamLeaderLanguageMatch ? teamLeaderLanguageMatch[1].trim() : '';

      return { teamLeader, teamLeaderCountryCode, teamLeaderCountryFlag, teamLeaderLanguage };
    });
    teamPageBar.increment();
    teamInfo.teamLeader = leaderSection.teamLeader;
    teamPageBar.increment();
    teamInfo.teamLeaderCountryCode = leaderSection.teamLeaderCountryCode;
    teamInfo.teamLeaderCountryFlag = leaderSection.teamLeaderCountryFlag;
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
      const works = Number(element.textContent.trim());
      return { works };
    });
    teamPageBar.increment();
    teamInfo.works = worksSection.works;
    teamPageBar.increment();

    const declaredWorksSection = await sectionElements[DECLARED].$eval('.counter', (element) => {
      const declaredWorks = Number(element.textContent.trim());
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

    // TODO add logging to see if members processes = memberCount & update split accordingly
    const memberSection = await sectionElements[MEMBERS].$$eval('p', (elements) => {
      const membersRaw = elements[0].textContent.trim();
      const memberCount = Number(elements[1].textContent.trim().match(/[\d]+/));
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
      const comment = element.innerHTML;
      return { comment };
    });
    teamPageBar.increment();
    teamInfo.teamComment = commentSection.comment;
    teamPageBar.increment();

    const registSection = await sectionElements[REGIST].$eval('strong', (element) => {
      const regist = element.textContent.trim();
      return { regist };
    });
    teamPageBar.increment();
    teamInfo.teamRegist = new Date(registSection.regist);
    teamPageBar.increment();

    const updateSection = await sectionElements[UPDATE].$eval('strong', (element) => {
      const update = element.textContent.trim();
      return { update };
    });
    teamPageBar.increment();
    teamInfo.teamUpdate = new Date(updateSection.update);
    teamPageBar.increment();





    // song Page Scraping
    for (const [songName, songInfo] of teamInfo.songs.entries()) {
      const songPageLink = songInfo.songPageLink;
      songPageBar.update({ filename: `Songs: ${songName}` });
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
        // Use Puppeteer to select the div with the specified class and style attribute
        const styleElement = await page.$('.section.parallax.nomargin.notopborder');

        if (styleElement) {
          // Extract the style attribute value
          const styleAttribute = await page.evaluate(el => el.getAttribute('style'), styleElement);

          if (styleAttribute) {
            // Convert the style attribute to a string
            const styleString = styleAttribute.toString();

            // Use a regular expression to find all URLs within the style attribute
            const urlMatches = styleString.match(/url\("([^"]*upload[^"]*)"\)/);

          } else {
            songInfo.bannerImageSrc = urlMatches;
          }
        } else {
          songInfo.bannerImageSrc = '';
        }
      } catch (error) {
        songInfo.bannerImageSrc = '';
      }

      songPageBar.increment();

      const specialElement = await page.$('span.badge.rounded-pill');
      let isSpecial = false;
      let specialTitle = '';
      if (specialElement) {
        isSpecial = true;
        specialTitle = await specialElement.$eval('big', (element) => element.textContent.trim());
      }
      songInfo.isSpecial = isSpecial;
      songInfo.specialTitle = specialTitle;
      songPageBar.increment();


      const bpmLevelBgaElement = await page.$eval('.col_two_third.nobottommargin, .col_full.nobottommargin', (element) => {
        bpmMatches = element.textContent.match(/bpm : ([^\/]+)/);
        let bpm = '';
        let bpmAverage = '';
        let bpmLower = '';
        let bpmUpper = '';
        if (bpmMatches[1].split('～').length > 1) {
          bpmLower = bpmMatches[1].split('～')[0].trim();
          bpmUpper = bpmMatches[1].split('～')[1].trim();
          bpmAverage = String((Number(bpmUpper) + Number(bpmLower)) / 2)
          // TODO investigate bms file to see if i can evaluate the most common bpm
        } else {
          bpm = bpmMatches[1].trim();
        }

        levelMatches = element.textContent.match(/Level : ([^\/]+)/);
        let levelLower = '';
        let levelUpper = '';
        if (levelMatches[1].split('～').length > 1) {
          levelLower = levelMatches[1].split('～')[0].trim();
          levelUpper = levelMatches[1].split('～')[1].trim();
        }

        bgaStatus = element.textContent.match(/BGA : (.*)/)[1].trim();
        return { bpm, bpmLower, bpmUpper, bpmAverage, levelLower, levelUpper, bgaStatus };
      });

      songInfo.bpm = Number(bpmLevelBgaElement.bpm);
      songInfo.bpmLower = Number(bpmLevelBgaElement.bpmLower);
      songInfo.bpmUpper = Number(bpmLevelBgaElement.bpmUpper);
      songInfo.bpmAverage = Number(bpmLevelBgaElement.bpmAverage);
      songInfo.levelLower = Number(bpmLevelBgaElement.levelLower.replace("★x", "").trim());
      songInfo.levelUpper = Number(bpmLevelBgaElement.levelUpper.replace("★x", "").trim());
      songInfo.bgaStatus = bpmLevelBgaElement.bgaStatus;
      songPageBar.increment();
      songPageBar.increment();
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



      // Extract only the linkUrls
      const linkUrls = await page.$$eval('blockquote p a', (elements) => {
        return elements.map((element) => element.getAttribute('href'));
      });
      songPageBar.increment();

      const downloadSize = await page.$eval('blockquote footer', (element) => {
        const footerString = element.textContent.trim();
        sizeQuantity = footerString.match(/Total : ([\d]+)/)[1];
        sizeUnit = footerString.match(/Total : [\d]+ ([a-zA-Z]+)/)[1];
        return { sizeQuantity, sizeUnit };
      });
      songInfo.downloadSizeQuantity = Number(downloadSize.sizeQuantity);
      songInfo.downloadSizeUnit = downloadSize.sizeUnit;
      songPageBar.increment();
      // // console.log('Link URLs:', linkUrls);

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


      // // console.log('Paragraph Text:', paragraphTexts);

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
        // console.log('No Link Found: ', songPageLink)
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
      // // console.log('Link Descriptions:', linkDescs);



      const tags = await page.$$eval('div.bmsinfo2 span.label', (labels) => {
        return labels.map((label) => label.textContent);
      })
      songInfo.tags = tags;

      songPageBar.increment();


      // Extract soundcloud link.
      debugger;
      try {
        // Wait for the iframe to load
        await page.waitForSelector('.m_audition iframe', { timeout: config.waitForTimer });

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

      let commentText = ''
      try {
        // First attempt to select the <p> element using the first XPath selector
        commentText = await page.$eval(
          'div.col_full:nth-child(4) > div:nth-child(8) > p:nth-child(2)',
          pElement => pElement.textContent.trim()
        );

      } catch (error) {
        // console.error('First attempt failed. Trying the second XPath selector.');

        try {
          // Second attempt to select the <p> element using the second XPath selector
          commentText = await page.$eval(
            'div.col_full:nth-child(4) > div:nth-child(9) > p:nth-child(2)',
            pElement => pElement.innerHTML
          );

        } catch (error) {
          // console.error('Both attempts failed. Error:', error);
        }
      }

      songInfo.songComment = commentText;
      songPageBar.increment();

      const productionElement = await page.$eval('.seisakukankyo', (element) => {
        return element.innerHTML.trim();  // Use innerHTML directly and trim any extra spaces
      });
      songInfo.productionEnvironment = productionElement;
      songPageBar.increment();

      const registTime = await page.$$eval('table.table.nobottommargin tbody tr', (rows) => {
        for (const row of rows) {
          const th = row.querySelector('th');
          if (th && th.textContent.trim() === 'Regist Time') {
            const td = row.querySelector('td');
            return td ? td.textContent.trim() : null;
          }
        }
        return null;
      });

      if (registTime) {
        songInfo.registDateTime = new Date(registTime);
      }
      songPageBar.increment();

      songInfo.impressionLink = `${config.events.urlPrefix}action=PrintImpression&num=${songInfo.entryNumber}&event=${config.latestEvent.eventKey}`;
      songPageBar.increment();
      songInfo.shortImpressionLink = `${config.events.urlPrefix}action=More_def${songInfo.entryNumber}&event=${config.latestEvent.eventKey}#shotimprssion_form`;
      songPageBar.increment();

      songInfo.points = await page.$$eval('div.col-md-3.col-xs-6.nobottommargin.center.hideen-xs', (elements) => {
        let points = [];
        for (const element of elements) {
          const pointValue = Number(element.querySelector('div.counter').textContent.trim());
          const pointName = element.querySelector('h5').textContent.trim();
          points.push({ pointValue, pointName });
        }
        return points;
      });
      songPageBar.increment();

      songInfo.Vote = await page.$$eval('div.col_full div.col_full div.col_half.center.nobottommargin', (elements) => {
        let vote = [];
        for (const element of elements) {
          const voteValue = Number(element.querySelector('div.counter').textContent.trim());
          const voteName = element.querySelector('h5').textContent.trim();
          vote.push({ voteValue, voteName });
        }
        return vote;
      });
      songPageBar.increment();

      let lastVoteDateString = '';
      try {
        lastVoteDateString = await page.$eval('div.col_full div.col_full div.col_full.center small', (element) => {
          return element.textContent.trim();
        });
      } catch (error) {
        // console.log('An Error occured:', error);
      }
      songInfo.lastVoteDateTime = new Date(convertJpDate(lastVoteDateString));
      songPageBar.increment();
      // TODO implement my own version of 投票者一覧 after scraping impressions. maybe a version where mouseover shows a short impression

      const shortImpressions = await page.$$eval('div.col_full div.col_full div.col_full div.spost.clearfix.nomarginbottom', (elements) => {
        let impressions = [];
        let points = '';
        let userName = '';
        let countryCode = '';
        let countryFlag = '';
        let jpDateTime = '';
        let userId = '';
        let userPost = '';

        for (const element of elements) {
          const pointsElement = element.querySelector('.points_oneline');
          points = pointsElement ? Number(pointsElement.textContent.trim()) : '';
          userName = element.querySelector('strong').textContent.trim();
          countryCode = element.querySelector('.flag') ? element.querySelector('.flag').getAttribute('title').trim() : '';
          countryFlag = element.querySelector('img') ? element.querySelector('img').src : '';
          jpDateTime = element.querySelector('div.entry-title small') ? element.querySelector('div.entry-title small').textContent.trim() : '';
          userIdMatch = element.querySelector('div.entry-title small span') ? element.querySelector('div.entry-title small span').textContent.match(/\(([^)]+)\)/) : '';
          userId = userIdMatch ? userIdMatch[1].trim() : '';
          userPost = element.querySelector('.entry-title:last-child').textContent.trim();


          impressions.push({ points, userName, countryCode, countryFlag, jpDateTime, userId, userPost });
        }
        return impressions;
      });
      const updatedshortImpressions = shortImpressions.map(item => { // Extract the Japanese date/time string
        const jpDateTime = item.jpDateTime;

        // Convert using convertJpDate function
        const updatedJpDateTime = new Date(convertJpDate(jpDateTime));


        // Return a new object with updated jpDateTime
        return { ...item, jpDateTime: updatedJpDateTime };
      });
      songInfo.shortImpressions = updatedshortImpressions;
      songPageBar.increment();


      const headerXPath1 = `/html/body/div[1]/section/div[2]/div/div[4]/div[10]/div[25]`;
      const headerXPath2 = `/html/body/div[1]/section/div[2]/div/div[4]/div[9]/div[25]`;

      let longImpressions = await scrapeLongImpressions(page, headerXPath1);

      // Retry with headerXPath2 if longImpressions is empty
      if (longImpressions.length === 0) {
        longImpressions = await scrapeLongImpressions(page, headerXPath2);
      }

      // Assign to your object
      songInfo.longImpressions = longImpressions;
      songPageBar.increment();






    } // end song page scraping
  } // end team scraping
  // console.dir(teams, { depth: null });
  // console.dir(collectedLogs, { depth: null });

  // write to file
  saveData();



  // await page.waitForTimeout(120000);
  await browser.close();
  multibar.stop();
  console.timeEnd('programExecution'); // End the timer
})();