const appApi = () => {
  if (!window?.go?.main?.App) {
    throw new Error('Wails backend API is unavailable');
  }
  return window.go.main.App;
};

export const getSettings = () => appApi().GetSettings();
export const saveSettings = (payload) => appApi().SaveSettings(payload);
export const analyzeFileByPath = (payload) => appApi().AnalyzeFileByPath(payload);
export const analyzeFileByContent = (payload) => appApi().AnalyzeFileByContent(payload);
export const summarizePayload = (payload) => appApi().SummarizePayload(payload);
export const getHistory = () => appApi().GetHistory();
export const pickFile = () => appApi().PickFile();
