// config.js
const config = {
  test: {
    url: 'https://example.com/test',
    navigationTimeout: 30000,
    waitForTimer: 500,
    numberOfTeamsToLimit: 2,
  },
  full: {
    url: 'https://example.com/full',
    navigationTimeout: 60000,
    waitForTimer: 500,
    numberOfTeamsToLimit: null,
  }
};

// Common variables
const commonConfig = {
  actions: {
    eventPage: 9,
    teamPage: 33,
    songPage: 31,
  },
  events: {
    urlPrefix: 'https://manbow.nothing.sh/event/event.cgi?',
    bofnt: {
      eventKey: 142,
      shortName: "bofnt",
      fullName: '[THE BMS OF FIGHTERS : NT]',
      subTitle: "-Twinkle Dream Traveler-",
    },
    bofet: {
      eventKey: 140,
      shortName: "bofet",
      fullName: '[THE BMS OF FIGHTERS : ET]',
      subTitle: "-Summer Dream Traveler-",
    },
  }
};

// Function to get the event with the latest eventKey
const getLatestEvent = (events) => {
  let latestEvent = null;
  for (const key in events) {
    if (events[key].eventKey && (!latestEvent || events[key].eventKey > latestEvent.eventKey)) {
      latestEvent = events[key];
    }
  }
  return latestEvent;
};

const latestEvent = getLatestEvent(commonConfig.events);

const mode = process.env.MODE || 'full'; // Default to 'full' if MODE is not set
export default { ...config[mode], ...commonConfig, latestEvent };
