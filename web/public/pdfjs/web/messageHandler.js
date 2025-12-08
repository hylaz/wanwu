window.addEventListener('message', message => {
  const { data } = message;
  const { type, value } = data;
  switch (type) {
    // 页面切换
    case 'SET_PAGE':
      window.PDFViewerApplication.page = value;
      break;
    // 高亮文本
    case 'HIGHT_LINE':
      setTimeout(() => {
        window.PDFViewerApplication.findBar.toggle();
        window.PDFViewerApplication.findBar.findField.value = value;
        window.PDFViewerApplication.findBar.highlightAll.checked = true;
        window.PDFViewerApplication.findBar.dispatchEvent('highlightallchange');
      }, 100);
      break;
    // 页面高度定位
    case 'SET_TOP':
      document.getElementById('viewerContainer').scrollTop =
        document.getElementById('viewerContainer').scrollTop +
        c_urlArray['top'] * 1;
    default:
      break;
  }
});
