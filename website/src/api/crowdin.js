const crowdin = require("@crowdin/crowdin-api-client");

// initialization of crowdin client
const { translationStatusApi } = new crowdin.default({
  token: process.env.CROWDIN_PERSONAL_TOKEN,
});

async function getTranslationProgress() {
  let translationProgress = {};

  await translationStatusApi.getProjectProgress(531392).then((res) => {
    for (const item of res.data) {
      translationProgress[item.data.languageId] = item.data.approvalProgress;
    }
  });

  return translationProgress;
}

module.exports = {
  getTranslationProgress,
};
