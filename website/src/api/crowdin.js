const crowdin = require("@crowdin/crowdin-api-client");
const personalToken = process.env.CROWDIN_PERSONAL_TOKEN;
const projectId = 531392;

// initialization of crowdin client
const initClient = () => {
  if (!personalToken) {
    console.warn(
      "No crowding personal token, some features might not work as expected"
    );
    return;
  }

  return new crowdin.default({
    token: personalToken,
  });
};

const client = initClient() || {};

async function getTranslationProgress() {
  let translationProgress = {};
  const { translationStatusApi } = client;

  // do nothing if client failed to init
  if (!translationStatusApi) {
    return translationProgress;
  }

  await translationStatusApi.getProjectProgress(projectId).then((res) => {
    for (const item of res.data) {
      translationProgress[item.data.languageId] = item.data.approvalProgress;
    }
  });

  return translationProgress;
}

module.exports = {
  getTranslationProgress,
};
