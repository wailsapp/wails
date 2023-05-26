const crowdin = require("@crowdin/crowdin-api-client");
const personalToken = process.env.CROWDIN_PERSONAL_TOKEN;
const projectId = "531392";
const branch = "v2";

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

async function getTranslationProgress() {
  let translationProgress = {};

  const client = initClient() || {};
  const { sourceFilesApi, translationStatusApi } = client;

  // do nothing if client failed to init
  if (!translationStatusApi) {
    return translationProgress;
  }

  const branchId = await sourceFilesApi
    .listProjectBranches(projectId)
    .then((res) => {
      for (const item of res.data) {
        if (item.data.name == branch) {
          return item.data.id;
        }
      }
    });

  await translationStatusApi
    .getBranchProgress(projectId, branchId)
    .then((res) => {
      for (const item of res.data) {
        translationProgress[item.data.languageId] = item.data.approvalProgress;
      }
    });

  return translationProgress;
}

module.exports = {
  getTranslationProgress,
};
