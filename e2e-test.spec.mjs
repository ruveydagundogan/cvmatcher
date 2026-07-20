import { chromium } from '@playwright/test';

(async () => {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();
  let scoreRequest = null;
  let scoreResponseBody = null;

  page.on('requestfinished', async (request) => {
    if (request.url().includes('/score') && request.method() === 'POST') {
      scoreRequest = request;
      try {
        scoreResponseBody = await request.response()?.text();
      } catch (e) {
        scoreResponseBody = `ERROR: ${e.message}`;
      }
    }
  });

  await page.goto('http://localhost:3000', { waitUntil: 'networkidle' });

  await page.waitForFunction(() => {
    const button = document.querySelector('button');
    return button && !button.disabled && button.textContent?.trim() === 'Ask AI';
  }, { timeout: 120000 });

  await page.fill('textarea', 'What is Flutter?');
  await page.click('button:has-text("Ask AI")');

  await page.waitForFunction(() => {
    const p = document.querySelector('div.min-h-24 p');
    return p && p.textContent && !p.textContent.includes('Response will appear here...') && !p.textContent.includes('Waiting for model response...');
  }, { timeout: 120000 });

  await page.waitForSelector('text=Decision Score', { timeout: 120000 });
  await page.waitForSelector('text=85 / 100', { timeout: 120000 });

  const displayedScore = await page.textContent('text=85 / 100');
  const responseText = await page.textContent('div.min-h-24 p');

  console.log('FOUND_SCORE_CARD=', displayedScore?.trim());
  console.log('FOUND_RESPONSE_TEXT=', responseText?.trim());
  console.log('SCORE_REQUEST_URL=', scoreRequest?.url() || 'none');
  console.log('SCORE_REQUEST_METHOD=', scoreRequest?.method() || 'none');
  console.log('SCORE_REQUEST_POST_DATA=', scoreRequest ? await scoreRequest.postData() : 'none');
  console.log('SCORE_RESPONSE_BODY=', scoreResponseBody);

  await browser.close();
})();
