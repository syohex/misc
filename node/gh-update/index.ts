import { Readable } from 'stream';
import * as tar from 'tar';

async function downloadAndExtract(url: string) {
  console.log(`Download: ${url}`);
  const res = await fetch(url, {
    redirect: 'follow',
  });
  const blob = await res.blob();
  const stream = Readable.from(blob.stream());

  stream.pipe(
    tar.t({
      onentry: entry => {
        console.log(entry.path);
      }
    })
  )
}

async function main() {
  const res = await fetch('https://api.github.com/repos/cli/cli/releases/latest');
  const json = await res.json();

  for (const asset of json.assets) {
    if (/linux_amd64\.tar\.gz$/.test(asset.name)) {
      const url = asset.browser_download_url;
      await downloadAndExtract(url);
    }
  }

  console.log(`latest version: ${json.tag_name}`)
}

main()
  .catch(err => {
    console.error(err);
  });