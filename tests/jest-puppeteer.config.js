module.exports = {
    launch: {
      // dumpio: true,
      ignoreHTTPSErrors: true,
       slowMo: 50,
      headless: true,
     // devtools: true,
      defaultViewport: {
        width: 1300,
        height: 800,
      }
    },
     browser: 'chromium',
    browserContext: 'incognito',

  }
