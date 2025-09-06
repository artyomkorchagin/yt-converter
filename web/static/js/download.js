function extractVideoID(url) {
    try {
      const urlObj = new URL(url);
      if (urlObj.hostname === 'youtu.be') {
        return urlObj.pathname.split('/')[1].split('?')[0];
      }
      if (urlObj.hostname.includes('youtube.com')) {
        const params = new URLSearchParams(urlObj.search);
        return params.get('v');
      }
    } catch (e) {
      if (/^[a-zA-Z0-9_-]{11}$/.test(url.trim())) {
        return url.trim();
      }
    }
    return null;
  }

  function downloadVideo() {
      const input = document.getElementById('urlInput').value;
      const videoID = extractVideoID(input);
      if (!videoID) {
          alert("Invalid YouTube URL");
          return;
      }

      const downloadUrl = `/api/v1/video/${encodeURIComponent(videoID)}`;
      const link = document.createElement('a');
      link.href = downloadUrl;
      link.click();
  }